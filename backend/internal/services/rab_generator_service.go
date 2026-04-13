package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/erp-sppg/backend/internal/models"
	"gorm.io/gorm"
)

var (
	ErrMenuPlanNotReady = errors.New("menu plan belum disetujui")
	ErrRABAlreadyExists = errors.New("RAB sudah dibuat untuk menu plan ini")
	ErrNoMenuItems      = errors.New("menu plan tidak memiliki item menu")
)

// RABGeneratorService handles auto-generation of RAB from menu plan
type RABGeneratorService struct {
	db    *gorm.DB
	notif *NotificationService
}

// NewRABGeneratorService creates a new RAB generator service
func NewRABGeneratorService(db *gorm.DB, notif *NotificationService) *RABGeneratorService {
	return &RABGeneratorService{
		db:    db,
		notif: notif,
	}
}

// GenerateRABFromMenuPlan creates a RAB when a menu plan is approved.
// 1. Load menu plan with items, recipes, recipe items
// 2. Aggregate ingredients
// 3. Recommend supplier for each ingredient
// 4. Create RAB + RABItems in transaction
// 5. Send notification to kepala_sppg
func (s *RABGeneratorService) GenerateRABFromMenuPlan(menuPlanID uint, sppgID uint, yayasanID uint, createdBy uint) (*models.RAB, error) {
	// Validate menu plan exists and is approved
	var menuPlan models.MenuPlan
	if err := s.db.First(&menuPlan, menuPlanID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMenuPlanNotFound
		}
		return nil, fmt.Errorf("gagal mengambil menu plan: %w", err)
	}

	if menuPlan.Status != "approved" {
		return nil, ErrMenuPlanNotReady
	}

	// Check if RAB already exists for this menu plan
	var existingRAB models.RAB
	err := s.db.Where("menu_plan_id = ?", menuPlanID).First(&existingRAB).Error
	if err == nil {
		// If existing RAB is in revision_requested status, delete it and regenerate
		if existingRAB.Status == "revision_requested" {
			// Delete old RAB items first
			s.db.Where("rab_id = ?", existingRAB.ID).Delete(&models.RABItem{})
			// Delete old RAB
			s.db.Delete(&existingRAB)
		} else {
			return nil, ErrRABAlreadyExists
		}
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("gagal memeriksa RAB existing: %w", err)
	}

	// Aggregate ingredients from menu plan
	requirements, err := s.AggregateIngredients(menuPlanID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengagregasi kebutuhan bahan: %w", err)
	}

	if len(requirements) == 0 {
		return nil, ErrNoMenuItems
	}

	// Generate RAB number
	rabNumber, err := s.generateRABNumber()
	if err != nil {
		return nil, fmt.Errorf("gagal membuat nomor RAB: %w", err)
	}

	// Build RAB and items
	sppgIDPtr := &sppgID
	yayasanIDPtr := &yayasanID
	rab := &models.RAB{
		RABNumber:  rabNumber,
		MenuPlanID: menuPlanID,
		SPPGID:     sppgIDPtr,
		YayasanID:  yayasanIDPtr,
		Status:     "draft",
		CreatedBy:  createdBy,
	}

	var totalAmount float64
	var rabItems []models.RABItem

	for _, req := range requirements {
		supplierID, unitPrice := s.RecommendSupplier(req.IngredientID, yayasanID)
		subtotal := req.TotalQuantity * unitPrice

		item := models.RABItem{
			IngredientID:          req.IngredientID,
			Quantity:              req.TotalQuantity,
			Unit:                  req.Unit,
			UnitPrice:             unitPrice,
			Subtotal:              subtotal,
			RecommendedSupplierID: supplierID,
			Status:                "pending",
		}
		rabItems = append(rabItems, item)
		totalAmount += subtotal
	}

	rab.TotalAmount = totalAmount

	// Create RAB + items in transaction
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(rab).Error; err != nil {
			return fmt.Errorf("gagal membuat RAB: %w", err)
		}

		for i := range rabItems {
			rabItems[i].RABID = rab.ID
		}
		if err := tx.Create(&rabItems).Error; err != nil {
			return fmt.Errorf("gagal membuat RAB items: %w", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// Reload RAB with items
	s.db.Preload("Items.Ingredient").Preload("Items.RecommendedSupplier").First(rab, rab.ID)

	// Send notification to kepala_sppg (graceful degradation)
	if s.notif != nil {
		s.sendRABCreatedNotification(sppgID, rab.RABNumber)
	}

	return rab, nil
}

// sendRABCreatedNotification sends notification to kepala_sppg users of the SPPG
func (s *RABGeneratorService) sendRABCreatedNotification(sppgID uint, rabNumber string) {
	var users []models.User
	s.db.Where("sppg_id = ? AND role = ? AND is_active = ?", sppgID, "kepala_sppg", true).Find(&users)

	for _, user := range users {
		notification := &models.Notification{
			UserID:  user.ID,
			Type:    "rab_created",
			Title:   "RAB Baru Dibuat",
			Message: fmt.Sprintf("RAB %s telah dibuat otomatis dari menu plan. Silakan review dan setujui.", rabNumber),
			Link:    "/rab",
		}
		if err := s.notif.CreateNotification(context.Background(), notification); err != nil {
			fmt.Printf("Peringatan: gagal mengirim notifikasi RAB ke user %d: %v\n", user.ID, err)
		}
	}
}

// AggregateIngredients traverses MenuPlan → MenuItems → Recipe → RecipeItems → SemiFinishedGoods → SFGRecipeIngredients
// and returns aggregated ingredient requirements.
func (s *RABGeneratorService) AggregateIngredients(menuPlanID uint) ([]IngredientRequirement, error) {
	// Load menu plan with full recipe chain
	var menuItems []models.MenuItem
	if err := s.db.Where("menu_plan_id = ?", menuPlanID).
		Preload("Recipe.RecipeItems.SemiFinishedGoods.Recipe.Ingredients.Ingredient").
		Find(&menuItems).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil menu items: %w", err)
	}

	if len(menuItems) == 0 {
		return nil, nil
	}

	// Aggregate: ingredientID → total quantity
	aggregated := make(map[uint]*IngredientRequirement)

	for _, menuItem := range menuItems {
		portions := float64(menuItem.Portions)

		for _, recipeItem := range menuItem.Recipe.RecipeItems {
			sfg := recipeItem.SemiFinishedGoods

			// Quantity of semi-finished goods needed per menu item
			sfgQty := recipeItem.Quantity * portions

			// Traverse SFG recipe to get raw ingredients
			if sfg.Recipe != nil && len(sfg.Recipe.Ingredients) > 0 {
				yieldAmount := sfg.Recipe.YieldAmount
				if yieldAmount <= 0 {
					yieldAmount = 1 // prevent division by zero
				}

				for _, sfgIngredient := range sfg.Recipe.Ingredients {
					// Scale ingredient quantity based on how much SFG we need vs yield
					ingredientQty := (sfgIngredient.Quantity / yieldAmount) * sfgQty

					if existing, ok := aggregated[sfgIngredient.IngredientID]; ok {
						existing.TotalQuantity += ingredientQty
					} else {
						aggregated[sfgIngredient.IngredientID] = &IngredientRequirement{
							IngredientID:   sfgIngredient.IngredientID,
							IngredientName: sfgIngredient.Ingredient.Name,
							Unit:           sfgIngredient.Ingredient.Unit,
							TotalQuantity:  ingredientQty,
						}
					}
				}
			}
		}
	}

	// Convert map to slice
	result := make([]IngredientRequirement, 0, len(aggregated))
	for _, req := range aggregated {
		result = append(result, *req)
	}

	return result, nil
}

// RecommendSupplier finds the best supplier for an ingredient within a yayasan.
// Filter: is_available=true, supplier linked to yayasan, supplier is_active=true
// Sort: unit_price ASC, quality_rating DESC
// Returns: supplier_id (or nil), unit_price (or 0)
func (s *RABGeneratorService) RecommendSupplier(ingredientID uint, yayasanID uint) (*uint, float64) {
	var result struct {
		SupplierID uint
		UnitPrice  float64
	}

	err := s.db.Table("supplier_products sp").
		Select("sp.supplier_id, sp.unit_price").
		Joins("JOIN suppliers s ON sp.supplier_id = s.id").
		Joins("JOIN supplier_yayasans sy ON s.id = sy.supplier_id AND sy.yayasan_id = ?", yayasanID).
		Where("sp.ingredient_id = ? AND sp.is_available = ? AND s.is_active = ?", ingredientID, true, true).
		Order("sp.unit_price ASC, s.quality_rating DESC").
		Limit(1).
		Scan(&result).Error

	if err != nil || result.SupplierID == 0 {
		return nil, 0
	}

	return &result.SupplierID, result.UnitPrice
}

// generateRABNumber generates a unique RAB number with format RAB-YYYYMMDD-XXXX
func (s *RABGeneratorService) generateRABNumber() (string, error) {
	now := time.Now()
	datePrefix := now.Format("20060102")

	var count int64
	s.db.Model(&models.RAB{}).
		Where("rab_number LIKE ?", fmt.Sprintf("RAB-%s-%%", datePrefix)).
		Count(&count)

	rabNumber := fmt.Sprintf("RAB-%s-%04d", datePrefix, count+1)

	// Race condition protection
	var existing models.RAB
	err := s.db.Where("rab_number = ?", rabNumber).First(&existing).Error
	if err == nil {
		rabNumber = fmt.Sprintf("RAB-%s-%04d", datePrefix, count+2)
	}

	return rabNumber, nil
}
