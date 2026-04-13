package services

import (
	"errors"
	"fmt"

	"github.com/erp-sppg/backend/internal/models"
	"gorm.io/gorm"
)

// RABSummary holds aggregated RAB tracking data
type RABSummary struct {
	TotalItems    int     `json:"total_items"`
	ItemsWithPO   int     `json:"items_with_po"`
	ItemsReceived int     `json:"items_received"`
	TotalBudget   float64 `json:"total_budget"`
	TotalSpent    float64 `json:"total_spent"`
}

// POTrackingItem represents a PO linked to a RAB with GRN status
type POTrackingItem struct {
	POID         uint    `json:"po_id"`
	PONumber     string  `json:"po_number"`
	SupplierName string  `json:"supplier_name"`
	TotalAmount  float64 `json:"total_amount"`
	Status       string  `json:"status"`
	HasGRN       bool    `json:"has_grn"`
	GRNID        *uint   `json:"grn_id"`
}

// RABComparisonItem represents planned vs actual per ingredient
type RABComparisonItem struct {
	IngredientID   uint    `json:"ingredient_id"`
	IngredientName string  `json:"ingredient_name"`
	PlannedQty     float64 `json:"planned_qty"`
	PlannedAmount  float64 `json:"planned_amount"`
	ActualQty      float64 `json:"actual_qty"`
	ActualAmount   float64 `json:"actual_amount"`
}

// RABService handles RAB CRUD and tracking
type RABService struct {
	db *gorm.DB
}

// NewRABService creates a new RAB service
func NewRABService(db *gorm.DB) *RABService {
	return &RABService{db: db}
}

// WithDB returns a new service instance with the given DB
func (s *RABService) WithDB(db *gorm.DB) *RABService {
	return &RABService{db: db}
}

// GetRABByID retrieves a RAB by ID with full preloads
func (s *RABService) GetRABByID(id uint) (*models.RAB, error) {
	var rab models.RAB
	db := s.db.Session(&gorm.Session{NewDB: true})
	err := db.
		Preload("Items.Ingredient").
		Preload("Items.RecommendedSupplier").
		Preload("Items.PurchaseOrder").
		Preload("Items.GoodsReceipt").
		Preload("MenuPlan").
		Preload("SPPG").
		Preload("SPPGApprover").
		Preload("YayasanApprover").
		First(&rab, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRABNotFound
		}
		return nil, fmt.Errorf("gagal mengambil RAB: %w", err)
	}

	// Manually load Creator since GORM preload may conflict with MenuPlan.Creator
	if rab.CreatedBy > 0 {
		var creator models.User
		if err := db.First(&creator, rab.CreatedBy).Error; err == nil {
			rab.Creator = creator
		}
	}

	return &rab, nil
}

// GetRABList retrieves RABs filtered by sppg_id or yayasan_id
func (s *RABService) GetRABList(sppgID *uint, yayasanID *uint) ([]models.RAB, error) {
	var rabs []models.RAB
	query := s.db.Session(&gorm.Session{NewDB: true}).Preload("MenuPlan").Preload("SPPG").Preload("Creator").Preload("Items")

	if sppgID != nil {
		query = query.Where("sppg_id = ?", *sppgID)
	}
	if yayasanID != nil {
		query = query.Where("yayasan_id = ?", *yayasanID)
	}

	if err := query.Order("created_at DESC").Find(&rabs).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil daftar RAB: %w", err)
	}

	return rabs, nil
}

// UpdateRAB validates status is editable, updates items, and recalculates total_amount
func (s *RABService) UpdateRAB(rabID uint, items []models.RABItem) error {
	var rab models.RAB
	if err := s.db.First(&rab, rabID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRABNotFound
		}
		return fmt.Errorf("gagal mengambil RAB: %w", err)
	}

	// Only allow edits on draft or revision_requested
	if rab.Status != "draft" && rab.Status != "revision_requested" {
		return ErrRABNotEditable
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete existing items
		if err := tx.Where("rab_id = ?", rabID).Delete(&models.RABItem{}).Error; err != nil {
			return fmt.Errorf("gagal menghapus RAB items lama: %w", err)
		}

		// Create new items and calculate total
		var totalAmount float64
		for i := range items {
			items[i].RABID = rabID
			items[i].Subtotal = items[i].Quantity * items[i].UnitPrice
			if items[i].Status == "" {
				items[i].Status = "pending"
			}
			totalAmount += items[i].Subtotal
		}

		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return fmt.Errorf("gagal membuat RAB items baru: %w", err)
			}
		}

		// Update total amount
		if err := tx.Model(&models.RAB{}).Where("id = ?", rabID).
			Update("total_amount", totalAmount).Error; err != nil {
			return fmt.Errorf("gagal mengupdate total RAB: %w", err)
		}

		return nil
	})
}

// GetRABSummary calculates summary metrics for a RAB
func (s *RABService) GetRABSummary(rabID uint) (*RABSummary, error) {
	var rab models.RAB
	if err := s.db.Preload("Items").First(&rab, rabID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRABNotFound
		}
		return nil, fmt.Errorf("gagal mengambil RAB: %w", err)
	}

	summary := &RABSummary{
		TotalItems:  len(rab.Items),
		TotalBudget: rab.TotalAmount,
	}

	for _, item := range rab.Items {
		if item.POID != nil {
			summary.ItemsWithPO++
		}
		if item.GRNID != nil {
			summary.ItemsReceived++
		}
	}

	// Calculate total spent from actual GRN data
	// Sum subtotals of PO items that have been received
	var totalSpent float64
	s.db.Table("rab_items").
		Select("COALESCE(SUM(rab_items.subtotal), 0)").
		Where("rab_items.rab_id = ? AND rab_items.grn_id IS NOT NULL", rabID).
		Scan(&totalSpent)
	summary.TotalSpent = totalSpent

	return summary, nil
}

// GetPOTracking returns POs linked to a RAB with GRN status
func (s *RABService) GetPOTracking(rabID uint) ([]POTrackingItem, error) {
	// Verify RAB exists
	var rab models.RAB
	if err := s.db.First(&rab, rabID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRABNotFound
		}
		return nil, fmt.Errorf("gagal mengambil RAB: %w", err)
	}

	// Find distinct POs linked to this RAB's items
	var pos []models.PurchaseOrder
	if err := s.db.Preload("Supplier").
		Where("rab_id = ?", rabID).
		Order("created_at DESC").
		Find(&pos).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil PO tracking: %w", err)
	}

	result := make([]POTrackingItem, 0, len(pos))
	for _, po := range pos {
		item := POTrackingItem{
			POID:         po.ID,
			PONumber:     po.PONumber,
			SupplierName: po.Supplier.Name,
			TotalAmount:  po.TotalAmount,
			Status:       po.Status,
		}

		// Check if GRN exists for this PO
		var grn models.GoodsReceipt
		err := s.db.Where("po_id = ?", po.ID).First(&grn).Error
		if err == nil {
			item.HasGRN = true
			item.GRNID = &grn.ID
		}

		result = append(result, item)
	}

	return result, nil
}

// GetRABComparison returns planned vs actual per ingredient
func (s *RABService) GetRABComparison(rabID uint) ([]RABComparisonItem, error) {
	var rab models.RAB
	if err := s.db.Preload("Items.Ingredient").First(&rab, rabID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRABNotFound
		}
		return nil, fmt.Errorf("gagal mengambil RAB: %w", err)
	}

	result := make([]RABComparisonItem, 0, len(rab.Items))
	for _, item := range rab.Items {
		comparison := RABComparisonItem{
			IngredientID:   item.IngredientID,
			IngredientName: item.Ingredient.Name,
			PlannedQty:     item.Quantity,
			PlannedAmount:  item.Subtotal,
		}

		// Get actual data from GRN if available
		if item.GRNID != nil {
			var grnItem models.GoodsReceiptItem
			err := s.db.Where("grn_id = ? AND ingredient_id = ?", *item.GRNID, item.IngredientID).
				First(&grnItem).Error
			if err == nil {
				comparison.ActualQty = grnItem.ReceivedQuantity
				// Actual amount = received qty * unit price from PO item
				if item.POID != nil {
					var poItem models.PurchaseOrderItem
					err := s.db.Where("po_id = ? AND ingredient_id = ?", *item.POID, item.IngredientID).
						First(&poItem).Error
					if err == nil {
						comparison.ActualAmount = grnItem.ReceivedQuantity * poItem.UnitPrice
					}
				}
			}
		}

		result = append(result, comparison)
	}

	return result, nil
}

// CheckAndCompleteRAB checks if all RABItems have grn_id and auto-completes the RAB
func (s *RABService) CheckAndCompleteRAB(rabID uint) error {
	var rab models.RAB
	if err := s.db.Preload("Items").First(&rab, rabID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRABNotFound
		}
		return fmt.Errorf("gagal mengambil RAB: %w", err)
	}

	// Only check if RAB is in approved_yayasan status
	if rab.Status != "approved_yayasan" {
		return nil
	}

	// Check if all items have grn_id
	if len(rab.Items) == 0 {
		return nil
	}

	allReceived := true
	for _, item := range rab.Items {
		if item.GRNID == nil {
			allReceived = false
			break
		}
	}

	if allReceived {
		if err := s.db.Model(&models.RAB{}).Where("id = ?", rabID).
			Update("status", "completed").Error; err != nil {
			return fmt.Errorf("gagal menyelesaikan RAB: %w", err)
		}
	}

	return nil
}
