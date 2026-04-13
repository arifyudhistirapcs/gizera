package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/erp-sppg/backend/internal/models"
	"gorm.io/gorm"
)

var (
	ErrPONotFound           = errors.New("purchase order tidak ditemukan")
	ErrPOValidation         = errors.New("validasi purchase order gagal")
	ErrPOAlreadyApproved    = errors.New("purchase order sudah disetujui")
	ErrPOCancelled          = errors.New("purchase order sudah dibatalkan")
	ErrInvalidPOStatus      = errors.New("status purchase order tidak valid")
	ErrSupplierNotLinked    = errors.New("supplier tidak terhubung dengan yayasan")
	ErrRABNotApprovedYayasan = errors.New("RAB belum disetujui yayasan")
	ErrNoPendingRABItems    = errors.New("tidak ada item RAB yang belum memiliki PO")
	ErrSupplierMismatch     = errors.New("supplier tidak sesuai dengan PO ini")
)

// BatchPOResult represents the result of batch PO creation from RAB
type BatchPOResult struct {
	POs     []models.PurchaseOrder `json:"pos"`
	Created int                    `json:"created"`
	Skipped int                    `json:"skipped"`
}

// PurchaseOrderService handles purchase order business logic
type PurchaseOrderService struct {
	db *gorm.DB
}

// NewPurchaseOrderService creates a new purchase order service
func NewPurchaseOrderService(db *gorm.DB) *PurchaseOrderService {
	return &PurchaseOrderService{
		db: db,
	}
}

// CreatePurchaseOrder creates a new purchase order
func (s *PurchaseOrderService) CreatePurchaseOrder(po *models.PurchaseOrder, items []models.PurchaseOrderItem, userID uint) error {
	// Validate supplier exists and is active
	var supplier models.Supplier
	if err := s.db.First(&supplier, po.SupplierID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("supplier tidak ditemukan")
		}
		return err
	}
	if !supplier.IsActive {
		return errors.New("supplier tidak aktif")
	}

	// RAB validation: if RABID is set, validate RAB exists and has status "approved_yayasan"
	if po.RABID != nil {
		var rab models.RAB
		if err := s.db.First(&rab, *po.RABID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrRABNotApprovedYayasan
			}
			return fmt.Errorf("gagal mengambil RAB: %w", err)
		}
		if rab.Status != "approved_yayasan" {
			return ErrRABNotApprovedYayasan
		}
	}

	// Supplier-yayasan validation: if YayasanID is set, validate supplier is linked
	if po.YayasanID != nil {
		var count int64
		s.db.Model(&models.SupplierYayasan{}).
			Where("supplier_id = ? AND yayasan_id = ?", po.SupplierID, *po.YayasanID).
			Count(&count)
		if count == 0 {
			return ErrSupplierNotLinked
		}
	}

	// Validate items
	if len(items) == 0 {
		return errors.New("purchase order harus memiliki minimal 1 item")
	}

	// Calculate total amount and validate items
	var totalAmount float64
	for i := range items {
		// Validate ingredient exists
		var ingredient models.Ingredient
		if err := s.db.First(&ingredient, items[i].IngredientID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("bahan baku dengan ID %d tidak ditemukan", items[i].IngredientID)
			}
			return err
		}

		// Calculate subtotal
		items[i].Subtotal = items[i].Quantity * items[i].UnitPrice
		totalAmount += items[i].Subtotal
	}

	// Generate PO number
	poNumber, err := s.generatePONumber()
	if err != nil {
		return err
	}

	// Set PO fields
	po.PONumber = poNumber
	po.Status = "pending"
	po.TotalAmount = totalAmount
	po.CreatedBy = userID
	po.OrderDate = time.Now()

	// Create PO in transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Create PO
		if err := tx.Create(po).Error; err != nil {
			return err
		}

		// Create PO items
		for i := range items {
			items[i].POID = po.ID
		}
		if err := tx.Create(&items).Error; err != nil {
			return err
		}

		// Update RABItems with po_id if RABID is set
		if po.RABID != nil {
			ingredientIDs := make([]uint, 0, len(items))
			for _, item := range items {
				ingredientIDs = append(ingredientIDs, item.IngredientID)
			}
			if len(ingredientIDs) > 0 {
				if err := tx.Model(&models.RABItem{}).
					Where("rab_id = ? AND ingredient_id IN ? AND po_id IS NULL", *po.RABID, ingredientIDs).
					Updates(map[string]interface{}{
						"po_id":  po.ID,
						"status": "po_created",
					}).Error; err != nil {
					return fmt.Errorf("gagal mengupdate RAB items: %w", err)
				}
			}
		}

		return nil
	})
}

// GetPurchaseOrderByID retrieves a purchase order by ID with related data
func (s *PurchaseOrderService) GetPurchaseOrderByID(id uint) (*models.PurchaseOrder, error) {
	var po models.PurchaseOrder
	err := s.db.Preload("Supplier").
		Preload("POItems.Ingredient").
		Preload("Creator").
		Preload("Approver").
		First(&po, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPONotFound
		}
		return nil, err
	}

	return &po, nil
}

// GetAllPurchaseOrders retrieves all purchase orders with optional status filter
func (s *PurchaseOrderService) GetAllPurchaseOrders(status string) ([]models.PurchaseOrder, error) {
	var pos []models.PurchaseOrder
	query := s.db.Preload("Supplier").
		Preload("POItems.Ingredient").
		Preload("Creator").
		Preload("Approver")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("order_date DESC").Find(&pos).Error
	return pos, err
}

// UpdatePurchaseOrder updates an existing purchase order (only if pending)
func (s *PurchaseOrderService) UpdatePurchaseOrder(id uint, updates *models.PurchaseOrder, items []models.PurchaseOrderItem) error {
	// Get existing PO
	existingPO, err := s.GetPurchaseOrderByID(id)
	if err != nil {
		return err
	}

	// Only allow updates if status is pending
	if existingPO.Status != "pending" {
		return errors.New("hanya purchase order dengan status pending yang dapat diubah")
	}

	// Validate supplier
	var supplier models.Supplier
	if err := s.db.First(&supplier, updates.SupplierID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("supplier tidak ditemukan")
		}
		return err
	}
	if !supplier.IsActive {
		return errors.New("supplier tidak aktif")
	}

	// Validate items
	if len(items) == 0 {
		return errors.New("purchase order harus memiliki minimal 1 item")
	}

	// Calculate total amount
	var totalAmount float64
	for i := range items {
		// Validate ingredient exists
		var ingredient models.Ingredient
		if err := s.db.First(&ingredient, items[i].IngredientID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("bahan baku dengan ID %d tidak ditemukan", items[i].IngredientID)
			}
			return err
		}

		// Calculate subtotal
		items[i].Subtotal = items[i].Quantity * items[i].UnitPrice
		totalAmount += items[i].Subtotal
	}

	// Update PO in transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete old PO items
		if err := tx.Where("po_id = ?", id).Delete(&models.PurchaseOrderItem{}).Error; err != nil {
			return err
		}

		// Update PO
		if err := tx.Model(&models.PurchaseOrder{}).Where("id = ?", id).Updates(map[string]interface{}{
			"supplier_id":       updates.SupplierID,
			"expected_delivery": updates.ExpectedDelivery,
			"total_amount":      totalAmount,
			"updated_at":        time.Now(),
		}).Error; err != nil {
			return err
		}

		// Create new PO items
		for i := range items {
			items[i].POID = id
		}
		if err := tx.Create(&items).Error; err != nil {
			return err
		}

		return nil
	})
}

// ApprovePurchaseOrder approves a purchase order (confirmed → approved)
func (s *PurchaseOrderService) ApprovePurchaseOrder(id uint, approverID uint) error {
	// Get existing PO
	po, err := s.GetPurchaseOrderByID(id)
	if err != nil {
		return err
	}

	// Check if already approved
	if po.Status == "approved" {
		return ErrPOAlreadyApproved
	}

	// Check if cancelled
	if po.Status == "cancelled" {
		return ErrPOCancelled
	}

	// Only allow approval if status is pending (kepala_sppg approves pending POs)
	if po.Status != "pending" {
		return ErrInvalidPOStatus
	}

	// Update status to approved
	now := time.Now()
	return s.db.Model(&models.PurchaseOrder{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      "approved",
		"approved_by": approverID,
		"approved_at": now,
		"updated_at":  now,
	}).Error
}

// CancelPurchaseOrder cancels a purchase order
func (s *PurchaseOrderService) CancelPurchaseOrder(id uint) error {
	// Get existing PO
	po, err := s.GetPurchaseOrderByID(id)
	if err != nil {
		return err
	}

	// Check if already received
	if po.Status == "received" {
		return errors.New("purchase order yang sudah diterima tidak dapat dibatalkan")
	}

	// Update status to cancelled
	return s.db.Model(&models.PurchaseOrder{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     "cancelled",
		"updated_at": time.Now(),
	}).Error
}

// GetPendingPurchaseOrders retrieves all pending purchase orders
func (s *PurchaseOrderService) GetPendingPurchaseOrders() ([]models.PurchaseOrder, error) {
	return s.GetAllPurchaseOrders("pending")
}

// GetApprovedPurchaseOrders retrieves all approved purchase orders
func (s *PurchaseOrderService) GetApprovedPurchaseOrders() ([]models.PurchaseOrder, error) {
	return s.GetAllPurchaseOrders("approved")
}

// GetPurchaseOrdersBySupplier retrieves all purchase orders for a specific supplier
func (s *PurchaseOrderService) GetPurchaseOrdersBySupplier(supplierID uint) ([]models.PurchaseOrder, error) {
	var pos []models.PurchaseOrder
	err := s.db.Preload("Supplier").
		Preload("POItems.Ingredient").
		Preload("Creator").
		Preload("Approver").
		Where("supplier_id = ?", supplierID).
		Order("order_date DESC").
		Find(&pos).Error
	return pos, err
}

// GetPurchaseOrdersByDateRange retrieves purchase orders within a date range
func (s *PurchaseOrderService) GetPurchaseOrdersByDateRange(startDate, endDate time.Time) ([]models.PurchaseOrder, error) {
	var pos []models.PurchaseOrder
	err := s.db.Preload("Supplier").
		Preload("POItems.Ingredient").
		Preload("Creator").
		Preload("Approver").
		Where("order_date BETWEEN ? AND ?", startDate, endDate).
		Order("order_date DESC").
		Find(&pos).Error
	return pos, err
}

// CreatePurchaseOrdersFromRAB creates POs grouped by supplier from an approved RAB
func (s *PurchaseOrderService) CreatePurchaseOrdersFromRAB(rabID uint, yayasanID uint, targetSPPGID uint, expectedDelivery time.Time, createdBy uint) (*BatchPOResult, error) {
	// Use NewDB session to avoid tenant scope interference inside the transaction
	baseDB := s.db.Session(&gorm.Session{NewDB: true})

	// 1. Load RAB and validate status
	var rab models.RAB
	if err := baseDB.First(&rab, rabID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRABNotFound
		}
		return nil, fmt.Errorf("gagal mengambil RAB: %w", err)
	}
	if rab.Status != "approved_yayasan" {
		return nil, ErrRABNotApprovedYayasan
	}

	// 2. Load RAB items that don't have po_id yet (pending items)
	var rabItems []models.RABItem
	if err := baseDB.
		Where("rab_id = ? AND po_id IS NULL AND status = ?", rabID, "pending").
		Preload("Ingredient").
		Find(&rabItems).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil item RAB: %w", err)
	}
	if len(rabItems) == 0 {
		return nil, ErrNoPendingRABItems
	}

	// 3. Group items by recommended_supplier_id
	supplierGroups := make(map[uint][]models.RABItem)
	skipped := 0
	for _, item := range rabItems {
		if item.RecommendedSupplierID == nil {
			skipped++
			continue
		}
		supplierGroups[*item.RecommendedSupplierID] = append(supplierGroups[*item.RecommendedSupplierID], item)
	}

	if len(supplierGroups) == 0 {
		return &BatchPOResult{
			POs:     []models.PurchaseOrder{},
			Created: 0,
			Skipped: skipped,
		}, nil
	}

	// 4. Create POs inside a transaction
	var createdPOs []models.PurchaseOrder

	err := baseDB.Transaction(func(tx *gorm.DB) error {
		for supplierID, items := range supplierGroups {
			// Validate supplier is active
			var supplier models.Supplier
			if err := tx.First(&supplier, supplierID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					skipped += len(items)
					continue
				}
				return fmt.Errorf("gagal mengambil supplier %d: %w", supplierID, err)
			}
			if !supplier.IsActive {
				skipped += len(items)
				continue
			}

			// Generate PO number
			poNumber, err := s.generatePONumberInTx(tx)
			if err != nil {
				return fmt.Errorf("gagal generate nomor PO: %w", err)
			}

			// Calculate total and build PO items
			var totalAmount float64
			var poItems []models.PurchaseOrderItem
			for _, rabItem := range items {
				subtotal := rabItem.Quantity * rabItem.UnitPrice
				totalAmount += subtotal
				poItems = append(poItems, models.PurchaseOrderItem{
					IngredientID: rabItem.IngredientID,
					Quantity:     rabItem.Quantity,
					UnitPrice:    rabItem.UnitPrice,
					Subtotal:     subtotal,
				})
			}

			// Create PO
			po := models.PurchaseOrder{
				PONumber:         poNumber,
				SupplierID:       supplierID,
				OrderDate:        time.Now(),
				ExpectedDelivery: expectedDelivery,
				Status:           "pending",
				TotalAmount:      totalAmount,
				CreatedBy:        createdBy,
				YayasanID:        &yayasanID,
				RABID:            &rabID,
				TargetSPPGID:     &targetSPPGID,
			}
			if err := tx.Create(&po).Error; err != nil {
				return fmt.Errorf("gagal membuat PO untuk supplier %d: %w", supplierID, err)
			}

			// Create PO items
			for i := range poItems {
				poItems[i].POID = po.ID
			}
			if err := tx.Create(&poItems).Error; err != nil {
				return fmt.Errorf("gagal membuat item PO: %w", err)
			}

			// Update RAB items with po_id and status
			rabItemIDs := make([]uint, 0, len(items))
			for _, item := range items {
				rabItemIDs = append(rabItemIDs, item.ID)
			}
			if err := tx.Model(&models.RABItem{}).
				Where("id IN ?", rabItemIDs).
				Updates(map[string]interface{}{
					"po_id":  po.ID,
					"status": "po_created",
				}).Error; err != nil {
				return fmt.Errorf("gagal mengupdate RAB items: %w", err)
			}

			po.Supplier = supplier
			po.POItems = poItems
			createdPOs = append(createdPOs, po)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &BatchPOResult{
		POs:     createdPOs,
		Created: len(createdPOs),
		Skipped: skipped,
	}, nil
}

// MarkAsShipping allows a supplier to mark an approved PO as shipping (approved → shipping)
func (s *PurchaseOrderService) MarkAsShipping(poID uint, supplierID uint) error {
	baseDB := s.db.Session(&gorm.Session{NewDB: true})

	var po models.PurchaseOrder
	if err := baseDB.First(&po, poID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPONotFound
		}
		return fmt.Errorf("gagal mengambil PO: %w", err)
	}

	if po.SupplierID != supplierID {
		return ErrSupplierMismatch
	}

	// Allow shipping from approved or confirmed (legacy data)
	if po.Status != "approved" && po.Status != "confirmed" {
		return ErrInvalidPOStatus
	}

	return baseDB.Model(&models.PurchaseOrder{}).Where("id = ?", poID).Updates(map[string]interface{}{
		"status":     "shipping",
		"updated_at": time.Now(),
	}).Error
}

// ConfirmBySupplier allows a supplier to confirm a PO as-is (pending → approved)
func (s *PurchaseOrderService) ConfirmBySupplier(poID uint, supplierID uint) error {
	baseDB := s.db.Session(&gorm.Session{NewDB: true})

	var po models.PurchaseOrder
	if err := baseDB.First(&po, poID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPONotFound
		}
		return fmt.Errorf("gagal mengambil PO: %w", err)
	}

	if po.SupplierID != supplierID {
		return ErrSupplierMismatch
	}

	if po.Status != "pending" {
		return ErrInvalidPOStatus
	}

	now := time.Now()
	return baseDB.Model(&models.PurchaseOrder{}).Where("id = ?", poID).Updates(map[string]interface{}{
		"status":      "approved",
		"approved_at": &now,
		"updated_at":  now,
	}).Error
}

// RequestRevisionBySupplier allows a supplier to request changes on a PO (pending → revision_by_supplier)
func (s *PurchaseOrderService) RequestRevisionBySupplier(poID uint, supplierID uint, items []models.PurchaseOrderItem, notes string) error {
	baseDB := s.db.Session(&gorm.Session{NewDB: true})

	var po models.PurchaseOrder
	if err := baseDB.First(&po, poID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPONotFound
		}
		return fmt.Errorf("gagal mengambil PO: %w", err)
	}

	if po.SupplierID != supplierID {
		return ErrSupplierMismatch
	}

	if po.Status != "pending" {
		return ErrInvalidPOStatus
	}

	if len(items) == 0 {
		return errors.New("item revisi harus memiliki minimal 1 item")
	}

	// Validate items and calculate total
	var totalAmount float64
	for i := range items {
		var ingredient models.Ingredient
		if err := baseDB.First(&ingredient, items[i].IngredientID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("bahan baku dengan ID %d tidak ditemukan", items[i].IngredientID)
			}
			return fmt.Errorf("gagal mengambil bahan baku: %w", err)
		}
		items[i].Subtotal = items[i].Quantity * items[i].UnitPrice
		totalAmount += items[i].Subtotal
	}

	return baseDB.Transaction(func(tx *gorm.DB) error {
		// Delete old PO items
		if err := tx.Where("po_id = ?", poID).Delete(&models.PurchaseOrderItem{}).Error; err != nil {
			return fmt.Errorf("gagal menghapus item PO lama: %w", err)
		}

		// Create new PO items
		for i := range items {
			items[i].POID = poID
		}
		if err := tx.Create(&items).Error; err != nil {
			return fmt.Errorf("gagal membuat item PO baru: %w", err)
		}

		// Update PO status, total, and revision notes
		if err := tx.Model(&models.PurchaseOrder{}).Where("id = ?", poID).Updates(map[string]interface{}{
			"status":                  "revision_by_supplier",
			"total_amount":            totalAmount,
			"supplier_revision_notes": notes,
			"updated_at":              time.Now(),
		}).Error; err != nil {
			return fmt.Errorf("gagal mengupdate PO: %w", err)
		}

		return nil
	})
}

// AcceptSupplierRevision allows yayasan to accept supplier's revision (revision_by_supplier → confirmed)
func (s *PurchaseOrderService) AcceptSupplierRevision(poID uint) error {
	baseDB := s.db.Session(&gorm.Session{NewDB: true})

	var po models.PurchaseOrder
	if err := baseDB.First(&po, poID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPONotFound
		}
		return fmt.Errorf("gagal mengambil PO: %w", err)
	}

	if po.Status != "revision_by_supplier" {
		return ErrInvalidPOStatus
	}

	now := time.Now()
	return baseDB.Model(&models.PurchaseOrder{}).Where("id = ?", poID).Updates(map[string]interface{}{
		"status":      "approved",
		"approved_at": &now,
		"updated_at":  now,
	}).Error
}

// RevisePOByYayasan allows yayasan to revise PO and send back to supplier (revision_by_supplier → pending)
func (s *PurchaseOrderService) RevisePOByYayasan(poID uint, items []models.PurchaseOrderItem) error {
	baseDB := s.db.Session(&gorm.Session{NewDB: true})

	var po models.PurchaseOrder
	if err := baseDB.First(&po, poID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPONotFound
		}
		return fmt.Errorf("gagal mengambil PO: %w", err)
	}

	if po.Status != "revision_by_supplier" {
		return ErrInvalidPOStatus
	}

	if len(items) == 0 {
		return errors.New("item revisi harus memiliki minimal 1 item")
	}

	// Validate items and calculate total
	var totalAmount float64
	for i := range items {
		var ingredient models.Ingredient
		if err := baseDB.First(&ingredient, items[i].IngredientID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("bahan baku dengan ID %d tidak ditemukan", items[i].IngredientID)
			}
			return fmt.Errorf("gagal mengambil bahan baku: %w", err)
		}
		items[i].Subtotal = items[i].Quantity * items[i].UnitPrice
		totalAmount += items[i].Subtotal
	}

	return baseDB.Transaction(func(tx *gorm.DB) error {
		// Delete old PO items
		if err := tx.Where("po_id = ?", poID).Delete(&models.PurchaseOrderItem{}).Error; err != nil {
			return fmt.Errorf("gagal menghapus item PO lama: %w", err)
		}

		// Create new PO items
		for i := range items {
			items[i].POID = poID
		}
		if err := tx.Create(&items).Error; err != nil {
			return fmt.Errorf("gagal membuat item PO baru: %w", err)
		}

		// Update PO: clear revision notes, update total, set status back to pending
		if err := tx.Model(&models.PurchaseOrder{}).Where("id = ?", poID).Updates(map[string]interface{}{
			"status":                  "pending",
			"total_amount":            totalAmount,
			"supplier_revision_notes": "",
			"updated_at":              time.Now(),
		}).Error; err != nil {
			return fmt.Errorf("gagal mengupdate PO: %w", err)
		}

		return nil
	})
}

// generatePONumber generates a unique PO number
func (s *PurchaseOrderService) generatePONumber() (string, error) {
	return s.generatePONumberInTx(s.db)
}

// generatePONumberInTx generates a unique PO number within a transaction
func (s *PurchaseOrderService) generatePONumberInTx(tx *gorm.DB) (string, error) {
	// Format: PO-YYYYMMDD-XXXX
	now := time.Now()
	datePrefix := now.Format("20060102")

	// Count POs created today
	var count int64
	tx.Model(&models.PurchaseOrder{}).
		Where("po_number LIKE ?", fmt.Sprintf("PO-%s-%%", datePrefix)).
		Count(&count)

	// Generate PO number
	poNumber := fmt.Sprintf("PO-%s-%04d", datePrefix, count+1)

	// Check if it already exists (race condition protection)
	var existing models.PurchaseOrder
	err := tx.Where("po_number = ?", poNumber).First(&existing).Error
	if err == nil {
		// If exists, try with incremented number
		poNumber = fmt.Sprintf("PO-%s-%04d", datePrefix, count+2)
	}

	return poNumber, nil
}
