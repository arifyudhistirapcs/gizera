package services

import (
	"errors"
	"time"

	"github.com/erp-sppg/backend/internal/models"
	"gorm.io/gorm"
)

var (
	ErrSupplierNotFound   = errors.New("supplier tidak ditemukan")
	ErrSupplierValidation = errors.New("validasi supplier gagal")
	ErrDuplicateSupplier  = errors.New("supplier dengan nama yang sama sudah ada")
)

// SupplierService handles supplier business logic
type SupplierService struct {
	db *gorm.DB
}

// NewSupplierService creates a new supplier service
func NewSupplierService(db *gorm.DB) *SupplierService {
	return &SupplierService{
		db: db,
	}
}

// CreateSupplier creates a new supplier
func (s *SupplierService) CreateSupplier(supplier *models.Supplier) error {
	// Check for duplicate name, scoped by sppg_id if present
	// Use NewDB session to avoid tenant middleware adding duplicate sppg_id conditions
	var existing models.Supplier
	dupQuery := s.db.Session(&gorm.Session{NewDB: true}).Where("name = ?", supplier.Name)
	if supplier.SPPGID != nil {
		dupQuery = dupQuery.Where("sppg_id = ?", *supplier.SPPGID)
	} else {
		dupQuery = dupQuery.Where("sppg_id IS NULL")
	}
	err := dupQuery.First(&existing).Error
	if err == nil {
		// Record found → duplicate
		return ErrDuplicateSupplier
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Unexpected DB error
		return err
	}
	// err == gorm.ErrRecordNotFound → name is available, proceed to create

	// Set defaults
	supplier.IsActive = true
	supplier.OnTimeDelivery = 0
	supplier.QualityRating = 0

	// Use unscoped DB for Create to avoid tenant scope interfering with GORM's post-insert SELECT
	return s.db.Session(&gorm.Session{NewDB: true}).Create(supplier).Error
}

// GetSupplierByID retrieves a supplier by ID
func (s *SupplierService) GetSupplierByID(id uint) (*models.Supplier, error) {
	var supplier models.Supplier
	err := s.db.First(&supplier, id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSupplierNotFound
		}
		return nil, err
	}

	return &supplier, nil
}

// GetAllSuppliers retrieves all suppliers
func (s *SupplierService) GetAllSuppliers(activeOnly bool) ([]models.Supplier, error) {
	var suppliers []models.Supplier
	query := s.db.Model(&models.Supplier{})

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	err := query.Order("name ASC").Find(&suppliers).Error
	return suppliers, err
}

// FilterSuppliers returns suppliers filtered by search, productCategory, and optional is_active.
// isActive nil = all, true = active only, false = inactive only
func (s *SupplierService) FilterSuppliers(search string, productCategory string, isActive *bool) ([]models.Supplier, error) {
	var suppliers []models.Supplier
	query := s.db.Model(&models.Supplier{})

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if productCategory != "" {
		query = query.Where("product_category = ?", productCategory)
	}

	err := query.Order("name ASC").Find(&suppliers).Error
	return suppliers, err
}

// UpdateSupplier updates an existing supplier
func (s *SupplierService) UpdateSupplier(id uint, updates *models.Supplier) error {
	// Check if supplier exists
	_, err := s.GetSupplierByID(id)
	if err != nil {
		return err
	}

	// Check for duplicate name (excluding current supplier)
	// Use NewDB session to avoid tenant middleware adding duplicate sppg_id conditions
	var existing models.Supplier
	dupQuery := s.db.Session(&gorm.Session{NewDB: true}).Where("name = ? AND id != ?", updates.Name, id)
	if updates.SPPGID != nil {
		dupQuery = dupQuery.Where("sppg_id = ?", *updates.SPPGID)
	}
	err = dupQuery.First(&existing).Error
	if err == nil {
		return ErrDuplicateSupplier
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Update supplier — use NewDB to avoid tenant middleware injecting duplicate FROM clause
	return s.db.Session(&gorm.Session{NewDB: true}).Model(&models.Supplier{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":             updates.Name,
		"contact_person":   updates.ContactPerson,
		"phone_number":     updates.PhoneNumber,
		"email":            updates.Email,
		"address":          updates.Address,
		"latitude":         updates.Latitude,
		"longitude":        updates.Longitude,
		"product_category": updates.ProductCategory,
		"updated_at":       time.Now(),
	}).Error
}

// DeactivateSupplier marks a supplier as inactive
func (s *SupplierService) DeactivateSupplier(id uint) error {
	result := s.db.Model(&models.Supplier{}).Where("id = ?", id).Update("is_active", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSupplierNotFound
	}
	return nil
}

// ActivateSupplier marks a supplier as active
func (s *SupplierService) ActivateSupplier(id uint) error {
	result := s.db.Model(&models.Supplier{}).Where("id = ?", id).Update("is_active", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSupplierNotFound
	}
	return nil
}

// SupplierPerformance represents performance metrics for a supplier
type SupplierPerformance struct {
	SupplierID       uint                `json:"supplier_id"`
	SupplierName     string              `json:"supplier_name"`
	TotalOrders      int                 `json:"total_orders"`
	CompletedOrders  int                 `json:"completed_orders"`
	OnTimeDeliveries int                 `json:"on_time_deliveries"`
	OnTimeRate       float64             `json:"on_time_rate"`
	QualityRating    float64             `json:"quality_rating"`
	TotalAmount      float64             `json:"total_amount"`
	LastOrderDate    *time.Time          `json:"last_order_date"`
	Transactions     []TransactionSummary `json:"transactions"`
}

// TransactionSummary represents a summary of a purchase order transaction
type TransactionSummary struct {
	PONumber  string    `json:"po_number"`
	OrderDate time.Time `json:"order_date"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
}

// GetSupplierPerformance retrieves performance metrics for a supplier
func (s *SupplierService) GetSupplierPerformance(id uint) (*SupplierPerformance, error) {
	supplier, err := s.GetSupplierByID(id)
	if err != nil {
		return nil, err
	}

	performance := &SupplierPerformance{
		SupplierID:   supplier.ID,
		SupplierName: supplier.Name,
	}

	// Count total orders
	var totalOrders int64
	s.db.Model(&models.PurchaseOrder{}).Where("supplier_id = ?", id).Count(&totalOrders)
	performance.TotalOrders = int(totalOrders)

	// Count completed orders (received status)
	var completedOrders int64
	s.db.Model(&models.PurchaseOrder{}).
		Where("supplier_id = ? AND status = ?", id, "received").
		Count(&completedOrders)
	performance.CompletedOrders = int(completedOrders)

	// Calculate on-time deliveries
	// An order is on-time if it was received before or on the expected delivery date
	var onTimeCount int64
	s.db.Model(&models.PurchaseOrder{}).
		Joins("LEFT JOIN goods_receipts ON goods_receipts.po_id = purchase_orders.id").
		Where("purchase_orders.supplier_id = ? AND purchase_orders.status = ? AND goods_receipts.receipt_date <= purchase_orders.expected_delivery", id, "received").
		Count(&onTimeCount)
	performance.OnTimeDeliveries = int(onTimeCount)

	// Calculate on-time rate
	if performance.CompletedOrders > 0 {
		performance.OnTimeRate = float64(performance.OnTimeDeliveries) / float64(performance.CompletedOrders) * 100
	}

	// Calculate average quality rating from all GRNs (not from supplier table)
	var avgRating float64
	s.db.Model(&models.GoodsReceipt{}).
		Joins("JOIN purchase_orders ON purchase_orders.id = goods_receipts.po_id").
		Where("purchase_orders.supplier_id = ? AND goods_receipts.quality_rating > 0", id).
		Select("COALESCE(AVG(goods_receipts.quality_rating), 0)").
		Scan(&avgRating)
	performance.QualityRating = avgRating

	// Calculate total amount
	var totalAmount float64
	s.db.Model(&models.PurchaseOrder{}).
		Where("supplier_id = ? AND status IN ?", id, []string{"approved", "received"}).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&totalAmount)
	performance.TotalAmount = totalAmount

	// Get last order date
	var lastOrder models.PurchaseOrder
	err = s.db.Where("supplier_id = ?", id).
		Order("order_date DESC").
		First(&lastOrder).Error
	if err == nil {
		performance.LastOrderDate = &lastOrder.OrderDate
	}

	// Update supplier's on-time delivery rate and quality rating in supplier table
	s.db.Model(&models.Supplier{}).Where("id = ?", id).Updates(map[string]interface{}{
		"on_time_delivery": performance.OnTimeRate,
		"quality_rating":   performance.QualityRating,
	})

	// Get recent transactions (last 10 purchase orders)
	var purchaseOrders []models.PurchaseOrder
	err = s.db.Where("supplier_id = ?", id).
		Order("order_date DESC").
		Limit(10).
		Find(&purchaseOrders).Error
	
	if err == nil {
		performance.Transactions = make([]TransactionSummary, 0, len(purchaseOrders))
		for _, po := range purchaseOrders {
			performance.Transactions = append(performance.Transactions, TransactionSummary{
				PONumber:  po.PONumber,
				OrderDate: po.OrderDate,
				Amount:    po.TotalAmount,
				Status:    po.Status,
			})
		}
	}

	return performance, nil
}

// UpdateSupplierRating updates the quality rating for a supplier
func (s *SupplierService) UpdateSupplierRating(id uint, rating float64, notes string) error {
	if rating < 1 || rating > 5 {
		return errors.New("rating harus antara 1 dan 5")
	}

	_, err := s.GetSupplierByID(id)
	if err != nil {
		return err
	}

	// Update the rating (simple average for now)
	// In a full implementation, we might want to track individual ratings
	return s.db.Model(&models.Supplier{}).Where("id = ?", id).Update("quality_rating", rating).Error
}

// SearchSuppliers searches suppliers by name or product category
func (s *SupplierService) SearchSuppliers(query string, productCategory string, activeOnly bool) ([]models.Supplier, error) {
	var suppliers []models.Supplier
	db := s.db.Model(&models.Supplier{})

	if activeOnly {
		db = db.Where("is_active = ?", true)
	}

	if query != "" {
		db = db.Where("name LIKE ?", "%"+query+"%")
	}

	if productCategory != "" {
		db = db.Where("product_category = ?", productCategory)
	}

	err := db.Order("name ASC").Find(&suppliers).Error
	return suppliers, err
}
// SupplierStats represents overall supplier statistics
type SupplierStats struct {
	TotalSuppliers   int              `json:"total_suppliers"`
	TotalSpending    float64          `json:"total_spending"`
	ActiveSuppliers  int              `json:"active_suppliers"`
	AverageRating    float64          `json:"average_rating"`
	TopSuppliers     []TopSupplier    `json:"top_suppliers"`
}

// TopSupplier represents a top performing supplier
type TopSupplier struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	TotalOrders int     `json:"total_orders"`
	TotalAmount float64 `json:"total_amount"`
}

// GetSupplierStats retrieves overall supplier statistics
func (s *SupplierService) GetSupplierStats() (*SupplierStats, error) {
	stats := &SupplierStats{}

	// Get total suppliers count
	var totalSuppliers int64
	if err := s.db.Session(&gorm.Session{}).Model(&models.Supplier{}).Count(&totalSuppliers).Error; err != nil {
		return nil, err
	}
	stats.TotalSuppliers = int(totalSuppliers)

	// Get active suppliers count
	var activeSuppliers int64
	if err := s.db.Session(&gorm.Session{}).Model(&models.Supplier{}).Where("is_active = ?", true).Count(&activeSuppliers).Error; err != nil {
		return nil, err
	}
	stats.ActiveSuppliers = int(activeSuppliers)

	// Get average quality rating
	var avgRating float64
	if err := s.db.Session(&gorm.Session{}).Model(&models.Supplier{}).
		Where("quality_rating > 0").
		Select("COALESCE(AVG(quality_rating), 0)").
		Scan(&avgRating).Error; err != nil {
		return nil, err
	}
	stats.AverageRating = avgRating

	// Get total spending from purchase orders (approved and received)
	var totalSpending float64
	if err := s.db.Session(&gorm.Session{}).Model(&models.PurchaseOrder{}).
		Where("status IN ?", []string{"approved", "received"}).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&totalSpending).Error; err != nil {
		return nil, err
	}
	stats.TotalSpending = totalSpending

	// Get top 3 suppliers by order count and total amount
	type SupplierOrderStats struct {
		SupplierID   uint
		SupplierName string
		OrderCount   int
		TotalAmount  float64
	}

	var topSupplierStats []SupplierOrderStats
	// Use raw SQL to avoid tenant scope ambiguity on JOIN queries
	if err := s.db.Session(&gorm.Session{NewDB: true}).Raw(`
		SELECT suppliers.id as supplier_id, suppliers.name as supplier_name, 
			COUNT(purchase_orders.id) as order_count, 
			COALESCE(SUM(purchase_orders.total_amount), 0) as total_amount
		FROM purchase_orders
		JOIN suppliers ON suppliers.id = purchase_orders.supplier_id
		WHERE purchase_orders.status IN ('approved', 'received')
		GROUP BY suppliers.id, suppliers.name
		ORDER BY total_amount DESC, order_count DESC
		LIMIT 3
	`).Scan(&topSupplierStats).Error; err != nil {
		topSupplierStats = nil
	}

	// Convert to TopSupplier format
	stats.TopSuppliers = make([]TopSupplier, 0, len(topSupplierStats))
	for _, ts := range topSupplierStats {
		stats.TopSuppliers = append(stats.TopSuppliers, TopSupplier{
			ID:          ts.SupplierID,
			Name:        ts.SupplierName,
			TotalOrders: ts.OrderCount,
			TotalAmount: ts.TotalAmount,
		})
	}

	return stats, nil
}

