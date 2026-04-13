package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/erp-sppg/backend/internal/middleware"
	"github.com/erp-sppg/backend/internal/models"
	"github.com/erp-sppg/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SupplierPortalHandler handles supplier portal endpoints
type SupplierPortalHandler struct {
	db                     *gorm.DB
	supplierProductService *services.SupplierProductService
	invoiceService         *services.InvoiceService
	purchaseOrderService   *services.PurchaseOrderService
}

// NewSupplierPortalHandler creates a new supplier portal handler
func NewSupplierPortalHandler(
	db *gorm.DB,
	supplierProductService *services.SupplierProductService,
	invoiceService *services.InvoiceService,
	purchaseOrderService *services.PurchaseOrderService,
) *SupplierPortalHandler {
	return &SupplierPortalHandler{
		db:                     db,
		supplierProductService: supplierProductService,
		invoiceService:         invoiceService,
		purchaseOrderService:   purchaseOrderService,
	}
}

// GetSupplierDashboard handles GET /supplier/dashboard
func (h *SupplierPortalHandler) GetSupplierDashboard(c *gin.Context) {
	supplierID, ok := middleware.GetSupplierID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success":    false,
			"error_code": "UNAUTHORIZED",
			"message":    "Supplier ID tidak ditemukan",
		})
		return
	}

	// Count active POs (pending + approved)
	var activePOCount int64
	h.db.Model(&models.PurchaseOrder{}).
		Where("supplier_id = ? AND status IN ?", supplierID, []string{"pending", "approved"}).
		Count(&activePOCount)

	// Count completed POs (received)
	var completedPOCount int64
	h.db.Model(&models.PurchaseOrder{}).
		Where("supplier_id = ? AND status = ?", supplierID, "received").
		Count(&completedPOCount)

	// Count pending invoices
	var pendingInvoiceCount int64
	h.db.Model(&models.Invoice{}).
		Where("supplier_id = ? AND status = ?", supplierID, "pending").
		Count(&pendingInvoiceCount)

	// Sum payments received
	var totalPaymentsReceived float64
	h.db.Table("payments").
		Joins("JOIN invoices ON invoices.id = payments.invoice_id").
		Where("invoices.supplier_id = ?", supplierID).
		Select("COALESCE(SUM(payments.amount), 0)").
		Scan(&totalPaymentsReceived)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total_po_aktif":        activePOCount,
			"total_po_selesai":      completedPOCount,
			"invoice_pending":       pendingInvoiceCount,
			"pembayaran_diterima":   totalPaymentsReceived,
		},
	})
}

// GetSupplierPayments handles GET /supplier/payments
func (h *SupplierPortalHandler) GetSupplierPayments(c *gin.Context) {
	supplierID, ok := middleware.GetSupplierID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success":    false,
			"error_code": "UNAUTHORIZED",
			"message":    "Supplier ID tidak ditemukan",
		})
		return
	}

	payments, err := h.invoiceService.GetSupplierPaymentHistory(supplierID)
	if err != nil {
		log.Printf("GetSupplierPayments error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    payments,
	})
}

// GetSupplierProducts handles GET /supplier-products
func (h *SupplierPortalHandler) GetSupplierProducts(c *gin.Context) {
	role, _ := c.Get("user_role")
	roleStr, _ := role.(string)

	switch roleStr {
	case "supplier":
		supplierID, ok := middleware.GetSupplierID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"error_code": "UNAUTHORIZED",
				"message":    "Supplier ID tidak ditemukan",
			})
			return
		}

		products, err := h.supplierProductService.GetProductsBySupplier(supplierID)
		if err != nil {
			log.Printf("GetSupplierProducts error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success":    false,
				"error_code": "INTERNAL_ERROR",
				"message":    "Terjadi kesalahan pada server",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    products,
		})

	case "kepala_yayasan":
		yayasanIDVal, ok := c.Get("yayasan_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"error_code": "UNAUTHORIZED",
				"message":    "Yayasan ID tidak ditemukan",
			})
			return
		}
		yayasanID := yayasanIDVal.(uint)

		products, err := h.supplierProductService.GetCatalogByYayasan(yayasanID)
		if err != nil {
			log.Printf("GetSupplierProducts (yayasan) error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success":    false,
				"error_code": "INTERNAL_ERROR",
				"message":    "Terjadi kesalahan pada server",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    products,
		})

	default:
		c.JSON(http.StatusForbidden, gin.H{
			"success":    false,
			"error_code": "FORBIDDEN",
			"message":    "Anda tidak memiliki izin untuk mengakses resource ini",
		})
	}
}

// CreateSupplierProductRequest represents the request body for creating a supplier product
type CreateSupplierProductRequest struct {
	IngredientID  uint    `json:"ingredient_id" binding:"required"`
	UnitPrice     float64 `json:"unit_price" binding:"required,gte=0"`
	MinOrderQty   float64 `json:"min_order_qty" binding:"gte=0"`
	StockQuantity float64 `json:"stock_quantity" binding:"gte=0"`
}

// CreateSupplierProduct handles POST /supplier-products
func (h *SupplierPortalHandler) CreateSupplierProduct(c *gin.Context) {
	supplierID, ok := middleware.GetSupplierID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success":    false,
			"error_code": "UNAUTHORIZED",
			"message":    "Supplier ID tidak ditemukan",
		})
		return
	}

	var req CreateSupplierProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	product := &models.SupplierProduct{
		SupplierID:    supplierID,
		IngredientID:  req.IngredientID,
		UnitPrice:     req.UnitPrice,
		MinOrderQty:   req.MinOrderQty,
		StockQuantity: req.StockQuantity,
	}

	if err := h.supplierProductService.CreateProduct(product); err != nil {
		if err == services.ErrDuplicateSupplierProduct {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "DUPLICATE_SUPPLIER_PRODUCT",
				"message":    "Produk supplier sudah ada untuk ingredient ini",
			})
			return
		}
		if err == services.ErrInvalidIngredientID {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INVALID_INGREDIENT",
				"message":    "Ingredient tidak valid",
			})
			return
		}
		log.Printf("CreateSupplierProduct error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Produk supplier berhasil dibuat",
		"data":    product,
	})
}

// UpdateSupplierProductRequest represents the request body for updating a supplier product
type UpdateSupplierProductRequest struct {
	IngredientID  uint    `json:"ingredient_id"`
	UnitPrice     float64 `json:"unit_price" binding:"gte=0"`
	MinOrderQty   float64 `json:"min_order_qty" binding:"gte=0"`
	IsAvailable   bool    `json:"is_available"`
	StockQuantity float64 `json:"stock_quantity" binding:"gte=0"`
}

// UpdateSupplierProduct handles PUT /supplier-products/:id
func (h *SupplierPortalHandler) UpdateSupplierProduct(c *gin.Context) {
	supplierID, ok := middleware.GetSupplierID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success":    false,
			"error_code": "UNAUTHORIZED",
			"message":    "Supplier ID tidak ditemukan",
		})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	var req UpdateSupplierProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	updates := &models.SupplierProduct{
		IngredientID:  req.IngredientID,
		UnitPrice:     req.UnitPrice,
		MinOrderQty:   req.MinOrderQty,
		IsAvailable:   req.IsAvailable,
		StockQuantity: req.StockQuantity,
	}

	if err := h.supplierProductService.UpdateProduct(uint(id), supplierID, updates); err != nil {
		if err == services.ErrSupplierProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "SUPPLIER_PRODUCT_NOT_FOUND",
				"message":    "Produk supplier tidak ditemukan",
			})
			return
		}
		if err == services.ErrSupplierProductUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{
				"success":    false,
				"error_code": "FORBIDDEN",
				"message":    "Tidak memiliki akses ke produk ini",
			})
			return
		}
		if err == services.ErrDuplicateSupplierProduct {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "DUPLICATE_SUPPLIER_PRODUCT",
				"message":    "Produk supplier sudah ada untuk ingredient ini",
			})
			return
		}
		if err == services.ErrInvalidIngredientID {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INVALID_INGREDIENT",
				"message":    "Ingredient tidak valid",
			})
			return
		}
		log.Printf("UpdateSupplierProduct error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Produk supplier berhasil diperbarui",
	})
}

// DeleteSupplierProduct handles DELETE /supplier-products/:id
func (h *SupplierPortalHandler) DeleteSupplierProduct(c *gin.Context) {
	supplierID, ok := middleware.GetSupplierID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success":    false,
			"error_code": "UNAUTHORIZED",
			"message":    "Supplier ID tidak ditemukan",
		})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	if err := h.supplierProductService.DeleteProduct(uint(id), supplierID); err != nil {
		if err == services.ErrSupplierProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "SUPPLIER_PRODUCT_NOT_FOUND",
				"message":    "Produk supplier tidak ditemukan",
			})
			return
		}
		if err == services.ErrSupplierProductUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{
				"success":    false,
				"error_code": "FORBIDDEN",
				"message":    "Tidak memiliki akses ke produk ini",
			})
			return
		}
		log.Printf("DeleteSupplierProduct error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Produk supplier berhasil dihapus",
	})
}
