package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/erp-sppg/backend/internal/models"
	"gorm.io/gorm"
)

var (
	ErrGRNNotFound       = errors.New("goods receipt tidak ditemukan")
	ErrGRNValidation     = errors.New("validasi goods receipt gagal")
	ErrPONotApproved     = errors.New("purchase order belum disetujui")
	ErrPOAlreadyReceived = errors.New("purchase order sudah diterima")
	ErrPOAlreadyHasGRN   = errors.New("PO sudah memiliki GRN")
)

// GoodsReceiptService handles goods receipt business logic
type GoodsReceiptService struct {
	db               *gorm.DB
	inventoryService *InventoryService
	cashFlowService  *CashFlowService
	rabService       *RABService
}

// NewGoodsReceiptService creates a new goods receipt service
func NewGoodsReceiptService(db *gorm.DB, inventoryService *InventoryService, cashFlowService *CashFlowService) *GoodsReceiptService {
	return &GoodsReceiptService{
		db:               db,
		inventoryService: inventoryService,
		cashFlowService:  cashFlowService,
		rabService:       NewRABService(db),
	}
}

// QuantityDiscrepancy represents a discrepancy between ordered and received quantities
type QuantityDiscrepancy struct {
	IngredientID      uint    `json:"ingredient_id"`
	IngredientName    string  `json:"ingredient_name"`
	OrderedQuantity   float64 `json:"ordered_quantity"`
	ReceivedQuantity  float64 `json:"received_quantity"`
	Difference        float64 `json:"difference"`
	DifferencePercent float64 `json:"difference_percent"`
}

// CreateGoodsReceipt creates a new goods receipt and updates inventory
func (s *GoodsReceiptService) CreateGoodsReceipt(grn *models.GoodsReceipt, items []models.GoodsReceiptItem, userID uint) error {
	// Validate PO exists and is shipping (or approved for legacy data)
	var po models.PurchaseOrder
	if err := s.db.Session(&gorm.Session{NewDB: true}).Preload("POItems.Ingredient").First(&po, grn.POID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("purchase order tidak ditemukan")
		}
		return err
	}

	if po.Status != "shipping" && po.Status != "approved" {
		return ErrPONotApproved
	}

	if po.Status == "received" {
		return ErrPOAlreadyReceived
	}

	// 1:1 PO-GRN validation: check if PO already has a GRN
	var existingGRNCount int64
	s.db.Session(&gorm.Session{NewDB: true}).Model(&models.GoodsReceipt{}).Where("po_id = ?", grn.POID).Count(&existingGRNCount)
	if existingGRNCount > 0 {
		return ErrPOAlreadyHasGRN
	}

	// Validate items
	if len(items) == 0 {
		return errors.New("goods receipt harus memiliki minimal 1 item")
	}

	// Validate items match PO items
	poItemsMap := make(map[uint]*models.PurchaseOrderItem)
	for i := range po.POItems {
		poItemsMap[po.POItems[i].IngredientID] = &po.POItems[i]
	}

	for i := range items {
		poItem, exists := poItemsMap[items[i].IngredientID]
		if !exists {
			return fmt.Errorf("bahan baku dengan ID %d tidak ada dalam purchase order", items[i].IngredientID)
		}
		items[i].OrderedQuantity = poItem.Quantity
	}

	// Generate GRN number
	grnNumber, err := s.generateGRNNumber()
	if err != nil {
		return err
	}

	// Set GRN fields
	grn.GRNNumber = grnNumber
	grn.ReceivedBy = userID
	grn.ReceiptDate = time.Now()

	// Set sppg_id from PO's target_sppg_id or sppg_id
	if po.TargetSPPGID != nil {
		grn.SPPGID = po.TargetSPPGID
	} else if po.SPPGID != nil {
		grn.SPPGID = po.SPPGID
	}

	// Create GRN in transaction (use fresh session to avoid tenant scope issues)
	return s.db.Session(&gorm.Session{NewDB: true}).Transaction(func(tx *gorm.DB) error {
		// Create GRN
		if err := tx.Create(grn).Error; err != nil {
			return err
		}

		// Create GRN items
		for i := range items {
			items[i].GRNID = grn.ID
		}
		if err := tx.Create(&items).Error; err != nil {
			return err
		}

		// Update PO status to received
		if err := tx.Model(&models.PurchaseOrder{}).Where("id = ?", grn.POID).Updates(map[string]interface{}{
			"status":     "received",
			"updated_at": time.Now(),
		}).Error; err != nil {
			return err
		}

		// Update inventory for each item
		for _, item := range items {
			if err := s.inventoryService.UpdateStockWithTx(tx, item.IngredientID, item.ReceivedQuantity, "in", grn.GRNNumber, userID, ""); err != nil {
				return err
			}
		}

		// Update RABItem.grn_id: find RABItems linked to this PO and update
		if po.RABID != nil {
			if err := tx.Model(&models.RABItem{}).
				Where("po_id = ? AND grn_id IS NULL", grn.POID).
				Updates(map[string]interface{}{
					"grn_id": grn.ID,
					"status": "grn_received",
				}).Error; err != nil {
				return fmt.Errorf("gagal mengupdate RAB items dengan grn_id: %w", err)
			}

			// Auto-complete RAB: check if all items are received
			rabSvc := s.rabService.WithDB(tx)
			if err := rabSvc.CheckAndCompleteRAB(*po.RABID); err != nil {
				return fmt.Errorf("gagal cek auto-complete RAB: %w", err)
			}
		}

		// Update supplier quality_rating average
		if grn.QualityRating > 0 {
			if err := s.updateSupplierQualityRating(tx, po.SupplierID); err != nil {
				return fmt.Errorf("gagal mengupdate quality rating supplier: %w", err)
			}
		}

		// Create cash flow entry
		if s.cashFlowService != nil {
			cashFlowEntry := &models.CashFlowEntry{
				Date:        grn.ReceiptDate,
				Category:    "bahan_baku",
				Type:        "expense",
				Amount:      po.TotalAmount,
				Description: fmt.Sprintf("Pembelian bahan baku dari %s (PO: %s)", po.Supplier.Name, po.PONumber),
				Reference:   grn.GRNNumber,
				CreatedBy:   userID,
			}
			if err := s.cashFlowService.CreateCashFlowEntryWithTx(tx, cashFlowEntry); err != nil {
				return err
			}
		}

		return nil
	})
}

// updateSupplierQualityRating calculates and updates the average quality_rating for a supplier
func (s *GoodsReceiptService) updateSupplierQualityRating(tx *gorm.DB, supplierID uint) error {
	var avgRating float64
	err := tx.Model(&models.GoodsReceipt{}).
		Joins("JOIN purchase_orders ON purchase_orders.id = goods_receipts.po_id").
		Where("purchase_orders.supplier_id = ? AND goods_receipts.quality_rating > 0", supplierID).
		Select("COALESCE(AVG(goods_receipts.quality_rating), 0)").
		Scan(&avgRating).Error
	if err != nil {
		return err
	}

	return tx.Model(&models.Supplier{}).Where("id = ?", supplierID).
		Update("quality_rating", avgRating).Error
}

// GetGoodsReceiptByID retrieves a goods receipt by ID with related data
func (s *GoodsReceiptService) GetGoodsReceiptByID(id uint) (*models.GoodsReceipt, error) {
	var grn models.GoodsReceipt
	err := s.db.Preload("PurchaseOrder.Supplier").
		Preload("PurchaseOrder.POItems.Ingredient").
		Preload("GRNItems.Ingredient").
		Preload("Receiver").
		First(&grn, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGRNNotFound
		}
		return nil, err
	}

	return &grn, nil
}

// GetAllGoodsReceipts retrieves all goods receipts
func (s *GoodsReceiptService) GetAllGoodsReceipts() ([]models.GoodsReceipt, error) {
	var grns []models.GoodsReceipt
	err := s.db.Preload("PurchaseOrder.Supplier").
		Preload("GRNItems.Ingredient").
		Preload("Receiver").
		Order("receipt_date DESC").
		Find(&grns).Error
	return grns, err
}

// GetGoodsReceiptsByDateRange retrieves goods receipts within a date range
func (s *GoodsReceiptService) GetGoodsReceiptsByDateRange(startDate, endDate time.Time) ([]models.GoodsReceipt, error) {
	var grns []models.GoodsReceipt
	err := s.db.Preload("PurchaseOrder.Supplier").
		Preload("GRNItems.Ingredient").
		Preload("Receiver").
		Where("receipt_date BETWEEN ? AND ?", startDate, endDate).
		Order("receipt_date DESC").
		Find(&grns).Error
	return grns, err
}

// GetGoodsReceiptsByPO retrieves goods receipt for a specific purchase order
func (s *GoodsReceiptService) GetGoodsReceiptsByPO(poID uint) (*models.GoodsReceipt, error) {
	var grn models.GoodsReceipt
	err := s.db.Preload("PurchaseOrder.Supplier").
		Preload("GRNItems.Ingredient").
		Preload("Receiver").
		Where("po_id = ?", poID).
		First(&grn).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGRNNotFound
		}
		return nil, err
	}

	return &grn, nil
}

// CheckQuantityDiscrepancies checks for discrepancies between ordered and received quantities
func (s *GoodsReceiptService) CheckQuantityDiscrepancies(grnID uint) ([]QuantityDiscrepancy, error) {
	grn, err := s.GetGoodsReceiptByID(grnID)
	if err != nil {
		return nil, err
	}

	var discrepancies []QuantityDiscrepancy
	for _, item := range grn.GRNItems {
		if item.OrderedQuantity != item.ReceivedQuantity {
			diff := item.ReceivedQuantity - item.OrderedQuantity
			diffPercent := 0.0
			if item.OrderedQuantity > 0 {
				diffPercent = (diff / item.OrderedQuantity) * 100
			}

			discrepancies = append(discrepancies, QuantityDiscrepancy{
				IngredientID:      item.IngredientID,
				IngredientName:    item.Ingredient.Name,
				OrderedQuantity:   item.OrderedQuantity,
				ReceivedQuantity:  item.ReceivedQuantity,
				Difference:        diff,
				DifferencePercent: diffPercent,
			})
		}
	}

	return discrepancies, nil
}

// UpdateInvoicePhoto updates the invoice photo URL for a goods receipt
func (s *GoodsReceiptService) UpdateInvoicePhoto(grnID uint, photoURL string) error {
	result := s.db.Model(&models.GoodsReceipt{}).Where("id = ?", grnID).Update("invoice_photo", photoURL)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrGRNNotFound
	}
	return nil
}

// generateGRNNumber generates a unique GRN number
func (s *GoodsReceiptService) generateGRNNumber() (string, error) {
	// Format: GRN-YYYYMMDD-XXXX
	now := time.Now()
	datePrefix := now.Format("20060102")

	baseDB := s.db.Session(&gorm.Session{NewDB: true})

	// Count GRNs created today
	var count int64
	baseDB.Model(&models.GoodsReceipt{}).
		Where("grn_number LIKE ?", fmt.Sprintf("GRN-%s-%%", datePrefix)).
		Count(&count)

	// Generate GRN number
	grnNumber := fmt.Sprintf("GRN-%s-%04d", datePrefix, count+1)

	// Check if it already exists (race condition protection)
	var existing models.GoodsReceipt
	err := baseDB.Where("grn_number = ?", grnNumber).First(&existing).Error
	if err == nil {
		// If exists, try with incremented number
		grnNumber = fmt.Sprintf("GRN-%s-%04d", datePrefix, count+2)
	}

	return grnNumber, nil
}
