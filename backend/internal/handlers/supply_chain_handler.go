package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/erp-sppg/backend/internal/middleware"
	"github.com/erp-sppg/backend/internal/models"
	"github.com/erp-sppg/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// uploadBaseDir is the base directory for file uploads
var uploadBaseDir = filepath.Join("uploads")

// SupplyChainHandler handles supply chain endpoints
type SupplyChainHandler struct {
	db                    *gorm.DB
	supplierService       *services.SupplierService
	purchaseOrderService  *services.PurchaseOrderService
	goodsReceiptService   *services.GoodsReceiptService
	inventoryService      *services.InventoryService
}

// NewSupplyChainHandler creates a new supply chain handler
func NewSupplyChainHandler(db *gorm.DB) *SupplyChainHandler {
	inventoryService := services.NewInventoryService(db)
	cashFlowService := services.NewCashFlowService(db)
	
	return &SupplyChainHandler{
		db:                   db,
		supplierService:      services.NewSupplierService(db),
		purchaseOrderService: services.NewPurchaseOrderService(db),
		goodsReceiptService:  services.NewGoodsReceiptService(db, inventoryService, cashFlowService),
		inventoryService:     inventoryService,
	}
}

// Supplier Endpoints

// CreateSupplierRequest represents create supplier request
type CreateSupplierRequest struct {
	Name            string  `json:"name" binding:"required"`
	ContactPerson   string  `json:"contact_person"`
	PhoneNumber     string  `json:"phone_number"`
	Email           string  `json:"email" binding:"omitempty,email"`
	Address         string  `json:"address"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	ProductCategory string  `json:"product_category"`
}

// CreateSupplier creates a new supplier
func (h *SupplyChainHandler) CreateSupplier(c *gin.Context) {
	var req CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	supplier := &models.Supplier{
		Name:            req.Name,
		ContactPerson:   req.ContactPerson,
		PhoneNumber:     req.PhoneNumber,
		Email:           req.Email,
		Address:         req.Address,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		ProductCategory: req.ProductCategory,
	}

	// Auto-inject sppg_id for SPPG-level roles
	if sppgID, ok := middleware.GetTenantSPPGID(c); ok {
		supplier.SPPGID = &sppgID
	}

	scopedService := h.supplierService.WithDB(getTenantScopedDB(c, h.db))
	if err := scopedService.CreateSupplier(supplier); err != nil {
		if err == services.ErrDuplicateSupplier {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "DUPLICATE_SUPPLIER",
				"message":    "Supplier dengan nama yang sama sudah ada",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	// If kepala_yayasan, also create SupplierYayasan junction record
	role, _ := c.Get("user_role")
	if roleStr, ok := role.(string); ok && roleStr == "kepala_yayasan" {
		if yayasanIDVal, yOk := c.Get("yayasan_id"); yOk {
			yayasanID := yayasanIDVal.(uint)
			supplierYayasan := &models.SupplierYayasan{
				SupplierID: supplier.ID,
				YayasanID:  yayasanID,
			}
			if err := h.db.Create(supplierYayasan).Error; err != nil {
				log.Printf("Warning: failed to create SupplierYayasan for supplier %d, yayasan %d: %v", supplier.ID, yayasanID, err)
			}
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":  true,
		"message":  "Supplier berhasil dibuat",
		"supplier": supplier,
	})
}

// GetSupplier retrieves a supplier by ID
func (h *SupplyChainHandler) GetSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	scopedService := h.supplierService.WithDB(getTenantScopedDB(c, h.db))
	supplier, err := scopedService.GetSupplierByID(uint(id))
	if err != nil {
		if err == services.ErrSupplierNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "SUPPLIER_NOT_FOUND",
				"message":    "Supplier tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"supplier": supplier,
	})
}

// GetAllSuppliers retrieves all suppliers
func (h *SupplyChainHandler) GetAllSuppliers(c *gin.Context) {
	// Support both "is_active" (frontend) and "active_only" (legacy)
	// isActiveFilter: "active" | "inactive" | "" (all)
	isActiveParam := c.Query("is_active")
	var isActiveFilter *bool
	switch isActiveParam {
	case "active":
		t := true
		isActiveFilter = &t
	case "inactive":
		f := false
		isActiveFilter = &f
	default:
		// legacy active_only param
		if c.Query("active_only") == "true" {
			t := true
			isActiveFilter = &t
		}
		// else nil = show all
	}

	// Support both "search" (frontend) and "q" (legacy)
	query := c.Query("search")
	if query == "" {
		query = c.Query("q")
	}
	productCategory := c.Query("product_category")

	// Check if kepala_yayasan — use JOIN on supplier_yayasans
	role, _ := c.Get("user_role")
	roleStr, _ := role.(string)

	if roleStr == "kepala_yayasan" {
		yayasanIDVal, yOk := c.Get("yayasan_id")
		if !yOk {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"error_code": "UNAUTHORIZED",
				"message":    "Yayasan ID tidak ditemukan",
			})
			return
		}
		yayasanID := yayasanIDVal.(uint)

		var suppliers []models.Supplier
		db := h.db.Joins("JOIN supplier_yayasans ON supplier_yayasans.supplier_id = suppliers.id").
			Where("supplier_yayasans.yayasan_id = ?", yayasanID)

		if isActiveFilter != nil {
			db = db.Where("suppliers.is_active = ?", *isActiveFilter)
		}
		if query != "" {
			db = db.Where("suppliers.name ILIKE ?", "%"+query+"%")
		}
		if productCategory != "" {
			db = db.Where("suppliers.product_category = ?", productCategory)
		}

		if err := db.Order("suppliers.name ASC").Find(&suppliers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success":    false,
				"error_code": "INTERNAL_ERROR",
				"message":    "Terjadi kesalahan pada server",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":   true,
			"suppliers": suppliers,
		})
		return
	}

	scopedService := h.supplierService.WithDB(getTenantScopedDB(c, h.db))
	suppliers, err := scopedService.FilterSuppliers(query, productCategory, isActiveFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"suppliers": suppliers,
	})
}

// UpdateSupplier updates an existing supplier
func (h *SupplyChainHandler) UpdateSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	var req CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	// If kepala_yayasan, verify supplier is linked to their yayasan
	role, _ := c.Get("user_role")
	if roleStr, ok := role.(string); ok && roleStr == "kepala_yayasan" {
		yayasanIDVal, yOk := c.Get("yayasan_id")
		if !yOk {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"error_code": "UNAUTHORIZED",
				"message":    "Yayasan ID tidak ditemukan",
			})
			return
		}
		yayasanID := yayasanIDVal.(uint)

		var count int64
		h.db.Model(&models.SupplierYayasan{}).
			Where("supplier_id = ? AND yayasan_id = ?", uint(id), yayasanID).
			Count(&count)
		if count == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"success":    false,
				"error_code": "FORBIDDEN",
				"message":    "Supplier tidak terhubung dengan yayasan Anda",
			})
			return
		}
	}

	supplier := &models.Supplier{
		Name:            req.Name,
		ContactPerson:   req.ContactPerson,
		PhoneNumber:     req.PhoneNumber,
		Email:           req.Email,
		Address:         req.Address,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		ProductCategory: req.ProductCategory,
	}

	scopedService := h.supplierService.WithDB(getTenantScopedDB(c, h.db))
	if err := scopedService.UpdateSupplier(uint(id), supplier); err != nil {
		if err == services.ErrSupplierNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "SUPPLIER_NOT_FOUND",
				"message":    "Supplier tidak ditemukan",
			})
			return
		}

		if err == services.ErrDuplicateSupplier {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "DUPLICATE_SUPPLIER",
				"message":    "Supplier dengan nama yang sama sudah ada",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Supplier berhasil diperbarui",
	})
}

// GetSupplierPerformance retrieves supplier performance metrics
func (h *SupplyChainHandler) GetSupplierPerformance(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	scopedService := h.supplierService.WithDB(getTenantScopedDB(c, h.db))
	performance, err := scopedService.GetSupplierPerformance(uint(id))
	if err != nil {
		if err == services.ErrSupplierNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "SUPPLIER_NOT_FOUND",
				"message":    "Supplier tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    performance,
	})
}

// GetSupplierStats retrieves supplier statistics
func (h *SupplyChainHandler) GetSupplierStats(c *gin.Context) {
	scopedService := h.supplierService.WithDB(getTenantScopedDB(c, h.db))
	stats, err := scopedService.GetSupplierStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"totalSuppliers":  stats.TotalSuppliers,
			"totalSpending":   stats.TotalSpending,
			"activeSuppliers": stats.ActiveSuppliers,
			"averageRating":   stats.AverageRating,
			"topSuppliers":    stats.TopSuppliers,
		},
	})
}

// Purchase Order Endpoints

// CreatePurchaseOrderRequest represents create PO request
type CreatePurchaseOrderRequest struct {
	SupplierID       uint                      `json:"supplier_id" binding:"required"`
	ExpectedDelivery string                    `json:"expected_delivery" binding:"required"`
	Items            []PurchaseOrderItemRequest `json:"items" binding:"required,min=1"`
}

// PurchaseOrderItemRequest represents PO item request
type PurchaseOrderItemRequest struct {
	IngredientID uint    `json:"ingredient_id" binding:"required"`
	Quantity     float64 `json:"quantity" binding:"required,gt=0"`
	UnitPrice    float64 `json:"unit_price" binding:"required,gte=0"`
}

// CreatePurchaseOrder creates a new purchase order
func (h *SupplyChainHandler) CreatePurchaseOrder(c *gin.Context) {
	var req CreatePurchaseOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	// Parse expected delivery date
	expectedDelivery, err := time.Parse("2006-01-02", req.ExpectedDelivery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_DATE",
			"message":    "Format tanggal tidak valid (gunakan YYYY-MM-DD)",
		})
		return
	}

	// Get user ID from context
	userID, _ := c.Get("user_id")

	po := &models.PurchaseOrder{
		SupplierID:       req.SupplierID,
		ExpectedDelivery: expectedDelivery,
	}

	var items []models.PurchaseOrderItem
	for _, item := range req.Items {
		items = append(items, models.PurchaseOrderItem{
			IngredientID: item.IngredientID,
			Quantity:     item.Quantity,
			UnitPrice:    item.UnitPrice,
		})
	}

	scopedService := h.purchaseOrderService.WithDB(getTenantScopedDB(c, h.db))
	if err := scopedService.CreatePurchaseOrder(po, items, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "CREATE_PO_ERROR",
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":        true,
		"message":        "Purchase order berhasil dibuat",
		"purchase_order": po,
	})
}

// GetPurchaseOrder retrieves a purchase order by ID
func (h *SupplyChainHandler) GetPurchaseOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	// Use fresh session to avoid tenant scope issues for supplier/yayasan roles
	role, _ := c.Get("user_role")
	roleStr, _ := role.(string)

	var po models.PurchaseOrder
	db := h.db.Session(&gorm.Session{NewDB: true})
	err = db.
		Preload("Supplier").
		Preload("POItems.Ingredient").
		Preload("Creator").
		Preload("Approver").
		Preload("RAB").
		Preload("TargetSPPG").
		Preload("Yayasan").
		First(&po, uint(id)).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "PO_NOT_FOUND",
				"message":    "Purchase order tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	// Verify access: supplier can only see their own POs
	if roleStr == "supplier" {
		supplierID, _ := middleware.GetSupplierID(c)
		if po.SupplierID != supplierID {
			c.JSON(http.StatusForbidden, gin.H{
				"success":    false,
				"error_code": "FORBIDDEN",
				"message":    "Anda tidak memiliki akses ke PO ini",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"purchase_order": po,
	})
}

// GetAllPurchaseOrders retrieves all purchase orders
func (h *SupplyChainHandler) GetAllPurchaseOrders(c *gin.Context) {
	status := c.Query("status")

	// For kepala_yayasan, use yayasan-aware query to include POs with yayasan_id
	role, _ := c.Get("user_role")
	roleStr, _ := role.(string)

	var pos []models.PurchaseOrder
	var err error

	if roleStr == "kepala_yayasan" {
		yayasanIDVal, _ := c.Get("yayasan_id")
		yayasanID, _ := yayasanIDVal.(uint)

		query := h.db.Session(&gorm.Session{NewDB: true}).
			Preload("Supplier").
			Preload("POItems.Ingredient").
			Preload("Creator").
			Preload("Approver").
			Preload("RAB").
			Preload("TargetSPPG").
			Where("yayasan_id = ? OR sppg_id IN (?)",
				yayasanID,
				h.db.Session(&gorm.Session{NewDB: true}).Table("sppgs").Select("id").Where("yayasan_id = ?", yayasanID))

		if status != "" {
			query = query.Where("status = ?", status)
		}
		err = query.Order("order_date DESC").Find(&pos).Error
	} else if roleStr == "supplier" {
		supplierID, _ := middleware.GetSupplierID(c)
		query := h.db.Session(&gorm.Session{NewDB: true}).
			Preload("Supplier").
			Preload("POItems.Ingredient").
			Preload("Creator").
			Preload("RAB").
			Preload("TargetSPPG").
			Preload("Yayasan").
			Where("supplier_id = ?", supplierID)

		if status != "" {
			query = query.Where("status = ?", status)
		}
		err = query.Order("order_date DESC").Find(&pos).Error
	} else {
		scopedService := h.purchaseOrderService.WithDB(getTenantScopedDB(c, h.db))
		pos, err = scopedService.GetAllPurchaseOrders(status)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"purchase_orders": pos,
	})
}

// UpdatePurchaseOrder updates an existing purchase order
func (h *SupplyChainHandler) UpdatePurchaseOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	var req CreatePurchaseOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	// Parse expected delivery date
	expectedDelivery, err := time.Parse("2006-01-02", req.ExpectedDelivery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_DATE",
			"message":    "Format tanggal tidak valid (gunakan YYYY-MM-DD)",
		})
		return
	}

	po := &models.PurchaseOrder{
		SupplierID:       req.SupplierID,
		ExpectedDelivery: expectedDelivery,
	}

	var items []models.PurchaseOrderItem
	for _, item := range req.Items {
		items = append(items, models.PurchaseOrderItem{
			IngredientID: item.IngredientID,
			Quantity:     item.Quantity,
			UnitPrice:    item.UnitPrice,
		})
	}

	scopedService := h.purchaseOrderService.WithDB(getTenantScopedDB(c, h.db))
	if err := scopedService.UpdatePurchaseOrder(uint(id), po, items); err != nil {
		if err == services.ErrPONotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "PO_NOT_FOUND",
				"message":    "Purchase order tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "UPDATE_PO_ERROR",
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Purchase order berhasil diperbarui",
	})
}

// ApprovePurchaseOrder approves a purchase order
func (h *SupplyChainHandler) ApprovePurchaseOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	// Get user ID from context
	userID, _ := c.Get("user_id")

	scopedService := h.purchaseOrderService.WithDB(getTenantScopedDB(c, h.db))
	if err := scopedService.ApprovePurchaseOrder(uint(id), userID.(uint)); err != nil {
		if err == services.ErrPONotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "PO_NOT_FOUND",
				"message":    "Purchase order tidak ditemukan",
			})
			return
		}

		if err == services.ErrPOAlreadyApproved {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "PO_ALREADY_APPROVED",
				"message":    "Purchase order sudah disetujui",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "APPROVE_PO_ERROR",
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Purchase order berhasil disetujui",
	})
}

// CreateBatchPORequest represents batch PO creation from RAB request
type CreateBatchPORequest struct {
	RABID            uint   `json:"rab_id" binding:"required"`
	ExpectedDelivery string `json:"expected_delivery" binding:"required"`
}

// CreateBatchPOFromRAB creates multiple POs from an approved RAB, grouped by supplier
func (h *SupplyChainHandler) CreateBatchPOFromRAB(c *gin.Context) {
	var req CreateBatchPORequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	// Parse expected delivery date
	expectedDelivery, err := time.Parse("2006-01-02", req.ExpectedDelivery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_DATE",
			"message":    "Format tanggal tidak valid (gunakan YYYY-MM-DD)",
		})
		return
	}

	// Get user ID from context
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success":    false,
			"error_code": "UNAUTHORIZED",
			"message":    "User tidak terautentikasi",
		})
		return
	}

	// Get yayasan_id from context
	yayasanID := getYayasanIDFromContext(c)

	// Load RAB to get target SPPG ID
	var rab models.RAB
	baseDB := h.db.Session(&gorm.Session{NewDB: true})
	if err := baseDB.First(&rab, req.RABID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success":    false,
			"error_code": "NOT_FOUND",
			"message":    "RAB tidak ditemukan",
		})
		return
	}

	// Determine target SPPG ID from RAB
	var targetSPPGID uint
	if rab.SPPGID != nil {
		targetSPPGID = *rab.SPPGID
	}

	// If yayasan_id is 0 (superadmin), use RAB's yayasan_id
	if yayasanID == 0 && rab.YayasanID != nil {
		yayasanID = *rab.YayasanID
	}

	result, err := h.purchaseOrderService.CreatePurchaseOrdersFromRAB(
		req.RABID, yayasanID, targetSPPGID, expectedDelivery, userID,
	)
	if err != nil {
		switch err {
		case services.ErrRABNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "NOT_FOUND",
				"message":    "RAB tidak ditemukan",
			})
		case services.ErrRABNotApprovedYayasan:
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "RAB_NOT_APPROVED",
				"message":    "RAB belum disetujui yayasan",
			})
		case services.ErrNoPendingRABItems:
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "NO_PENDING_ITEMS",
				"message":    "Tidak ada item RAB yang belum memiliki PO",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success":    false,
				"error_code": "INTERNAL_ERROR",
				"message":    "Gagal membuat batch PO",
				"details":    err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": fmt.Sprintf("Berhasil membuat %d purchase order dari RAB", result.Created),
		"data":    result,
	})
}

// Goods Receipt Endpoints

// CreateGoodsReceiptRequest represents create GRN request
type CreateGoodsReceiptRequest struct {
	POID          uint                      `json:"po_id" binding:"required"`
	Notes         string                    `json:"notes"`
	QualityRating float64                   `json:"quality_rating" binding:"gte=0,lte=5"` // 0-5 rating
	Items         []GoodsReceiptItemRequest `json:"items" binding:"required,min=1"`
}

// GoodsReceiptItemRequest represents GRN item request
type GoodsReceiptItemRequest struct {
	IngredientID     uint    `json:"ingredient_id" binding:"required"`
	ReceivedQuantity float64 `json:"received_quantity" binding:"required,gte=0"`
	ExpiryDate       *string `json:"expiry_date"`
}

// CreateGoodsReceipt creates a new goods receipt
func (h *SupplyChainHandler) CreateGoodsReceipt(c *gin.Context) {
	var req CreateGoodsReceiptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, _ := c.Get("user_id")

	// Get SPPG ID from tenant context for the GRN record
	sppgIDVal, _ := c.Get("sppg_id")
	var sppgIDPtr *uint
	if sppgID, ok := sppgIDVal.(uint); ok && sppgID > 0 {
		sppgIDPtr = &sppgID
	}

	grn := &models.GoodsReceipt{
		POID:          req.POID,
		SPPGID:        sppgIDPtr,
		Notes:         req.Notes,
		QualityRating: req.QualityRating,
	}

	var items []models.GoodsReceiptItem
	for _, item := range req.Items {
		grnItem := models.GoodsReceiptItem{
			IngredientID:     item.IngredientID,
			ReceivedQuantity: item.ReceivedQuantity,
		}

		// Parse expiry date if provided
		if item.ExpiryDate != nil && *item.ExpiryDate != "" {
			expiryDate, err := time.Parse("2006-01-02", *item.ExpiryDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success":    false,
					"error_code": "INVALID_DATE",
					"message":    "Format tanggal kadaluarsa tidak valid (gunakan YYYY-MM-DD)",
				})
				return
			}
			grnItem.ExpiryDate = &expiryDate
		}

		items = append(items, grnItem)
	}

	scopedService := h.goodsReceiptService.WithDB(getTenantScopedDB(c, h.db))
	if err := scopedService.CreateGoodsReceipt(grn, items, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "CREATE_GRN_ERROR",
			"message":    err.Error(),
		})
		return
	}

	// Update supplier metrics (rating and on-time delivery) after GRN is created
	// This is done asynchronously to not block the response
	go func() {
		// Get PO to find supplier
		var po models.PurchaseOrder
		if err := h.db.Preload("Supplier").First(&po, req.POID).Error; err == nil {
			// Update supplier performance metrics
			_, _ = h.supplierService.GetSupplierPerformance(po.SupplierID)
		}
	}()

	c.JSON(http.StatusCreated, gin.H{
		"success":       true,
		"message":       "Goods receipt berhasil dibuat dan metrik supplier diperbarui",
		"goods_receipt": grn,
	})
}

// GetAllGoodsReceipts retrieves all goods receipts
func (h *SupplyChainHandler) GetAllGoodsReceipts(c *gin.Context) {
	scopedService := h.goodsReceiptService.WithDB(getTenantScopedDB(c, h.db))
	grns, err := scopedService.GetAllGoodsReceipts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	// Add received_by_name field for frontend
	var response []map[string]interface{}
	for _, grn := range grns {
		grnMap := map[string]interface{}{
			"id":              grn.ID,
			"grn_number":      grn.GRNNumber,
			"po_id":           grn.POID,
			"receipt_date":    grn.ReceiptDate,
			"invoice_photo":   grn.InvoicePhoto,
			"received_by":     grn.ReceivedBy,
			"notes":           grn.Notes,
			"quality_rating":  grn.QualityRating,
			"created_at":      grn.CreatedAt,
			"purchase_order":  grn.PurchaseOrder,
			"receiver":        grn.Receiver,
			"grn_items":       grn.GRNItems,
			"received_by_name": grn.Receiver.FullName,
		}
		response = append(response, grnMap)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"goods_receipts": response,
	})
}

// GetGoodsReceipt retrieves a goods receipt by ID
func (h *SupplyChainHandler) GetGoodsReceipt(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	scopedService := h.goodsReceiptService.WithDB(getTenantScopedDB(c, h.db))
	grn, err := scopedService.GetGoodsReceiptByID(uint(id))
	if err != nil {
		if err == services.ErrGRNNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "GRN_NOT_FOUND",
				"message":    "Goods receipt tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	// Add received_by_name field for frontend
	response := map[string]interface{}{
		"id":              grn.ID,
		"grn_number":      grn.GRNNumber,
		"po_id":           grn.POID,
		"receipt_date":    grn.ReceiptDate,
		"invoice_photo":   grn.InvoicePhoto,
		"received_by":     grn.ReceivedBy,
		"notes":           grn.Notes,
		"quality_rating":  grn.QualityRating,
		"created_at":      grn.CreatedAt,
		"purchase_order":  grn.PurchaseOrder,
		"receiver":        grn.Receiver,
		"grn_items":       grn.GRNItems,
		"received_by_name": grn.Receiver.FullName,
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"goods_receipt": response,
	})
}

// UploadInvoicePhoto uploads invoice photo for a goods receipt
func (h *SupplyChainHandler) UploadInvoicePhoto(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	// Get file from form
	file, err := c.FormFile("invoice_photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "NO_FILE",
			"message":    "File foto invoice tidak ditemukan",
		})
		return
	}

	// Save file
	filename := fmt.Sprintf("invoice_%d_%d%s", id, time.Now().Unix(), filepath.Ext(file.Filename))
	invoiceDir := filepath.Join(uploadBaseDir, "invoices")
	savePath := filepath.Join(invoiceDir, filename)
	
	// Create directory if not exists
	if err := os.MkdirAll(invoiceDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "UPLOAD_ERROR",
			"message":    "Gagal membuat direktori upload",
		})
		return
	}

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "UPLOAD_ERROR",
			"message":    "Gagal menyimpan file",
		})
		return
	}

	// Generate URL
	photoURL := fmt.Sprintf("/uploads/invoices/%s", filename)

	if err := h.goodsReceiptService.UpdateInvoicePhoto(uint(id), photoURL); err != nil {
		if err == services.ErrGRNNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "GRN_NOT_FOUND",
				"message":    "Goods receipt tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "Foto invoice berhasil diunggah",
		"photo_url": photoURL,
	})
}

// Inventory Endpoints

// GetInventory retrieves all inventory items
func (h *SupplyChainHandler) GetInventory(c *gin.Context) {
	scopedService := h.inventoryService.WithDB(getTenantScopedDB(c, h.db))
	items, err := scopedService.GetAllInventory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"inventory_items": items,
	})
}

// GetInventoryByIngredient retrieves inventory for a specific ingredient
func (h *SupplyChainHandler) GetInventoryByIngredient(c *gin.Context) {
	ingredientID, err := strconv.ParseUint(c.Param("ingredient_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_INGREDIENT_ID",
			"message":    "ID bahan tidak valid",
		})
		return
	}

	scopedService := h.inventoryService.WithDB(getTenantScopedDB(c, h.db))
	inventory, err := scopedService.GetInventoryByIngredient(uint(ingredientID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    inventory,
	})
}

// GetInventoryAlerts retrieves low stock alerts
func (h *SupplyChainHandler) GetInventoryAlerts(c *gin.Context) {
	scopedService := h.inventoryService.WithDB(getTenantScopedDB(c, h.db))
	alerts, err := scopedService.CheckLowStock()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"alerts":  alerts,
	})
}

// GetInventoryMovements retrieves inventory movements
func (h *SupplyChainHandler) GetInventoryMovements(c *gin.Context) {
	var ingredientID *uint
	if idStr := c.Query("ingredient_id"); idStr != "" {
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err == nil {
			uid := uint(id)
			ingredientID = &uid
		}
	}

	movementType := c.Query("movement_type")

	var startDate, endDate *time.Time
	if startStr := c.Query("start_date"); startStr != "" {
		if sd, err := time.Parse("2006-01-02", startStr); err == nil {
			startDate = &sd
		}
	}
	if endStr := c.Query("end_date"); endStr != "" {
		if ed, err := time.Parse("2006-01-02", endStr); err == nil {
			endDate = &ed
		}
	}

	scopedService := h.inventoryService.WithDB(getTenantScopedDB(c, h.db))
	movements, err := scopedService.GetMovements(ingredientID, movementType, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"movements": movements,
	})
}

// InitializeInventory initializes inventory records for all ingredients that don't have one
func (h *SupplyChainHandler) InitializeInventory(c *gin.Context) {
	scopedService := h.inventoryService.WithDB(getTenantScopedDB(c, h.db))
	if err := scopedService.InitializeInventoryForAllIngredients(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INITIALIZE_ERROR",
			"message":    "Gagal menginisialisasi inventory: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Inventory berhasil diinisialisasi untuk semua bahan",
	})
}

// InitializeInventoryItem initializes inventory record for a specific ingredient
func (h *SupplyChainHandler) InitializeInventoryItem(c *gin.Context) {
	ingredientID, err := strconv.ParseUint(c.Param("ingredient_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID bahan tidak valid",
		})
		return
	}

	scopedService := h.inventoryService.WithDB(getTenantScopedDB(c, h.db))
	if err := scopedService.InitializeInventoryForIngredient(uint(ingredientID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INITIALIZE_ERROR",
			"message":    "Gagal menginisialisasi inventory: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Bahan berhasil ditambahkan ke inventory",
	})
}

// === PO Confirmation/Revision by Supplier ===

// ConfirmPOBySupplier handles POST /purchase-orders/:id/confirm
// Supplier confirms PO as-is (pending → confirmed)
func (h *SupplyChainHandler) ConfirmPOBySupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	supplierID, ok := middleware.GetSupplierID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success":    false,
			"error_code": "UNAUTHORIZED",
			"message":    "Supplier ID tidak ditemukan",
		})
		return
	}

	if err := h.purchaseOrderService.ConfirmBySupplier(uint(id), supplierID); err != nil {
		switch err {
		case services.ErrPONotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "PO_NOT_FOUND",
				"message":    "Purchase order tidak ditemukan",
			})
		case services.ErrSupplierMismatch:
			c.JSON(http.StatusForbidden, gin.H{
				"success":    false,
				"error_code": "FORBIDDEN",
				"message":    "Anda tidak memiliki akses ke PO ini",
			})
		case services.ErrInvalidPOStatus:
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INVALID_STATUS",
				"message":    "PO hanya dapat dikonfirmasi jika statusnya pending",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success":    false,
				"error_code": "INTERNAL_ERROR",
				"message":    "Terjadi kesalahan pada server",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Purchase order berhasil dikonfirmasi oleh supplier",
	})
}

// MarkPOAsShipping handles POST /purchase-orders/:id/shipping
// Supplier marks PO as shipping (approved → shipping)
func (h *SupplyChainHandler) MarkPOAsShipping(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	supplierID, ok := middleware.GetSupplierID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success":    false,
			"error_code": "UNAUTHORIZED",
			"message":    "Supplier ID tidak ditemukan",
		})
		return
	}

	if err := h.purchaseOrderService.MarkAsShipping(uint(id), supplierID); err != nil {
		switch err {
		case services.ErrPONotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "PO_NOT_FOUND",
				"message":    "Purchase order tidak ditemukan",
			})
		case services.ErrSupplierMismatch:
			c.JSON(http.StatusForbidden, gin.H{
				"success":    false,
				"error_code": "FORBIDDEN",
				"message":    "Anda tidak memiliki akses ke PO ini",
			})
		case services.ErrInvalidPOStatus:
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INVALID_STATUS",
				"message":    "PO hanya dapat dikirim jika statusnya sudah disetujui",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success":    false,
				"error_code": "INTERNAL_ERROR",
				"message":    "Terjadi kesalahan pada server",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Status PO berhasil diubah ke sedang dikirim",
	})
}

// SupplierRevisionRequest represents supplier revision request body
type SupplierRevisionRequest struct {
	Notes string                     `json:"notes" binding:"required"`
	Items []PurchaseOrderItemRequest `json:"items" binding:"required,min=1"`
}

// RequestPORevisionBySupplier handles POST /purchase-orders/:id/request-revision
// Supplier requests changes on a PO (pending → revision_by_supplier)
func (h *SupplyChainHandler) RequestPORevisionBySupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	supplierID, ok := middleware.GetSupplierID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success":    false,
			"error_code": "UNAUTHORIZED",
			"message":    "Supplier ID tidak ditemukan",
		})
		return
	}

	var req SupplierRevisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	var items []models.PurchaseOrderItem
	for _, item := range req.Items {
		items = append(items, models.PurchaseOrderItem{
			IngredientID: item.IngredientID,
			Quantity:     item.Quantity,
			UnitPrice:    item.UnitPrice,
		})
	}

	if err := h.purchaseOrderService.RequestRevisionBySupplier(uint(id), supplierID, items, req.Notes); err != nil {
		switch err {
		case services.ErrPONotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "PO_NOT_FOUND",
				"message":    "Purchase order tidak ditemukan",
			})
		case services.ErrSupplierMismatch:
			c.JSON(http.StatusForbidden, gin.H{
				"success":    false,
				"error_code": "FORBIDDEN",
				"message":    "Anda tidak memiliki akses ke PO ini",
			})
		case services.ErrInvalidPOStatus:
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INVALID_STATUS",
				"message":    "PO hanya dapat direvisi jika statusnya pending",
			})
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "REVISION_ERROR",
				"message":    err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Permintaan revisi PO berhasil dikirim",
	})
}

// AcceptSupplierRevision handles POST /purchase-orders/:id/accept-revision
// Yayasan accepts supplier's revision (revision_by_supplier → confirmed)
func (h *SupplyChainHandler) AcceptSupplierRevision(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	if err := h.purchaseOrderService.AcceptSupplierRevision(uint(id)); err != nil {
		switch err {
		case services.ErrPONotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "PO_NOT_FOUND",
				"message":    "Purchase order tidak ditemukan",
			})
		case services.ErrInvalidPOStatus:
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INVALID_STATUS",
				"message":    "PO hanya dapat diterima jika statusnya revision_by_supplier",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success":    false,
				"error_code": "INTERNAL_ERROR",
				"message":    "Terjadi kesalahan pada server",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Revisi supplier berhasil diterima",
	})
}

// RevisePOByYayasan handles POST /purchase-orders/:id/revise
// Yayasan revises PO and sends back to supplier (revision_by_supplier → pending)
func (h *SupplyChainHandler) RevisePOByYayasan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	var req struct {
		Items []PurchaseOrderItemRequest `json:"items" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	var items []models.PurchaseOrderItem
	for _, item := range req.Items {
		items = append(items, models.PurchaseOrderItem{
			IngredientID: item.IngredientID,
			Quantity:     item.Quantity,
			UnitPrice:    item.UnitPrice,
		})
	}

	if err := h.purchaseOrderService.RevisePOByYayasan(uint(id), items); err != nil {
		switch err {
		case services.ErrPONotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "PO_NOT_FOUND",
				"message":    "Purchase order tidak ditemukan",
			})
		case services.ErrInvalidPOStatus:
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INVALID_STATUS",
				"message":    "PO hanya dapat direvisi jika statusnya revision_by_supplier",
			})
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "REVISE_PO_ERROR",
				"message":    err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "PO berhasil direvisi dan dikirim kembali ke supplier",
	})
}
