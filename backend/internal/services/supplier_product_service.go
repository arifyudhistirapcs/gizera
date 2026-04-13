package services

import (
	"errors"
	"fmt"

	"github.com/erp-sppg/backend/internal/models"
	"gorm.io/gorm"
)

var (
	ErrSupplierProductNotFound     = errors.New("produk supplier tidak ditemukan")
	ErrDuplicateSupplierProduct    = errors.New("produk supplier sudah ada untuk ingredient ini")
	ErrInvalidIngredientID         = errors.New("ingredient tidak valid")
	ErrSupplierProductUnauthorized = errors.New("tidak memiliki akses ke produk ini")
)

// SupplierProductService handles supplier product catalog business logic
type SupplierProductService struct {
	db *gorm.DB
}

// NewSupplierProductService creates a new supplier product service
func NewSupplierProductService(db *gorm.DB) *SupplierProductService {
	return &SupplierProductService{db: db}
}

// WithDB returns a new service instance with the given DB (for tenant scoping)
func (s *SupplierProductService) WithDB(db *gorm.DB) *SupplierProductService {
	return &SupplierProductService{db: db}
}

// CreateProduct validates ingredient_id, checks unique constraint, and creates the product
func (s *SupplierProductService) CreateProduct(product *models.SupplierProduct) error {
	// Validate ingredient exists
	var ingredient models.Ingredient
	if err := s.db.First(&ingredient, product.IngredientID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidIngredientID
		}
		return fmt.Errorf("gagal memvalidasi ingredient: %w", err)
	}

	// Check unique constraint (supplier_id, ingredient_id)
	var existing models.SupplierProduct
	err := s.db.Where("supplier_id = ? AND ingredient_id = ?", product.SupplierID, product.IngredientID).
		First(&existing).Error
	if err == nil {
		return ErrDuplicateSupplierProduct
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("gagal memeriksa duplikasi produk: %w", err)
	}

	// Set default availability
	product.IsAvailable = true

	if err := s.db.Create(product).Error; err != nil {
		return fmt.Errorf("gagal membuat produk supplier: %w", err)
	}

	return nil
}

// UpdateProduct validates ownership and updates the product fields
func (s *SupplierProductService) UpdateProduct(id uint, supplierID uint, updates *models.SupplierProduct) error {
	// Find existing product
	var product models.SupplierProduct
	if err := s.db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSupplierProductNotFound
		}
		return fmt.Errorf("gagal mengambil produk supplier: %w", err)
	}

	// Validate ownership
	if product.SupplierID != supplierID {
		return ErrSupplierProductUnauthorized
	}

	// If ingredient_id is being changed, validate it exists and check uniqueness
	if updates.IngredientID != 0 && updates.IngredientID != product.IngredientID {
		var ingredient models.Ingredient
		if err := s.db.First(&ingredient, updates.IngredientID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrInvalidIngredientID
			}
			return fmt.Errorf("gagal memvalidasi ingredient: %w", err)
		}

		// Check unique constraint for new ingredient_id
		var existing models.SupplierProduct
		err := s.db.Where("supplier_id = ? AND ingredient_id = ? AND id != ?", supplierID, updates.IngredientID, id).
			First(&existing).Error
		if err == nil {
			return ErrDuplicateSupplierProduct
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("gagal memeriksa duplikasi produk: %w", err)
		}
	}

	// Update fields
	updateMap := map[string]interface{}{
		"unit_price":     updates.UnitPrice,
		"min_order_qty":  updates.MinOrderQty,
		"is_available":   updates.IsAvailable,
		"stock_quantity": updates.StockQuantity,
	}
	if updates.IngredientID != 0 {
		updateMap["ingredient_id"] = updates.IngredientID
	}

	if err := s.db.Model(&models.SupplierProduct{}).Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return fmt.Errorf("gagal mengupdate produk supplier: %w", err)
	}

	return nil
}

// DeleteProduct hard deletes a product after validating ownership
func (s *SupplierProductService) DeleteProduct(id uint, supplierID uint) error {
	var product models.SupplierProduct
	if err := s.db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSupplierProductNotFound
		}
		return fmt.Errorf("gagal mengambil produk supplier: %w", err)
	}

	if product.SupplierID != supplierID {
		return ErrSupplierProductUnauthorized
	}

	if err := s.db.Delete(&models.SupplierProduct{}, id).Error; err != nil {
		return fmt.Errorf("gagal menghapus produk supplier: %w", err)
	}

	return nil
}

// GetProductsBySupplier returns all products for a given supplier, preloading Ingredient
func (s *SupplierProductService) GetProductsBySupplier(supplierID uint) ([]models.SupplierProduct, error) {
	var products []models.SupplierProduct
	if err := s.db.Preload("Ingredient").
		Where("supplier_id = ?", supplierID).
		Order("created_at DESC").
		Find(&products).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil produk supplier: %w", err)
	}
	return products, nil
}

// GetCatalogByYayasan returns available products from suppliers linked to a yayasan
func (s *SupplierProductService) GetCatalogByYayasan(yayasanID uint) ([]models.SupplierProduct, error) {
	var products []models.SupplierProduct
	if err := s.db.Preload("Supplier").Preload("Ingredient").
		Joins("JOIN suppliers ON suppliers.id = supplier_products.supplier_id").
		Joins("JOIN supplier_yayasans ON supplier_yayasans.supplier_id = suppliers.id").
		Where("supplier_yayasans.yayasan_id = ? AND supplier_products.is_available = ? AND suppliers.is_active = ?", yayasanID, true, true).
		Order("supplier_products.created_at DESC").
		Find(&products).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil katalog supplier: %w", err)
	}
	return products, nil
}

// ToggleAvailability updates is_available without deleting, after validating ownership
func (s *SupplierProductService) ToggleAvailability(id uint, supplierID uint, isAvailable bool) error {
	var product models.SupplierProduct
	if err := s.db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSupplierProductNotFound
		}
		return fmt.Errorf("gagal mengambil produk supplier: %w", err)
	}

	if product.SupplierID != supplierID {
		return ErrSupplierProductUnauthorized
	}

	if err := s.db.Model(&models.SupplierProduct{}).Where("id = ?", id).
		Update("is_available", isAvailable).Error; err != nil {
		return fmt.Errorf("gagal mengupdate ketersediaan produk: %w", err)
	}

	return nil
}
