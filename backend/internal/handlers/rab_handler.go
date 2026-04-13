package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/erp-sppg/backend/internal/models"
	"github.com/erp-sppg/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// RABHandler handles RAB (Rencana Anggaran Belanja) endpoints
type RABHandler struct {
	rabService     *services.RABService
	rabGenerator   *services.RABGeneratorService
	approvalEngine *services.ApprovalEngineService
}

// NewRABHandler creates a new RAB handler
func NewRABHandler(rabService *services.RABService, rabGenerator *services.RABGeneratorService, approvalEngine *services.ApprovalEngineService) *RABHandler {
	return &RABHandler{
		rabService:     rabService,
		rabGenerator:   rabGenerator,
		approvalEngine: approvalEngine,
	}
}

// GetRABList handles GET /rab — list RABs scoped by role
func (h *RABHandler) GetRABList(c *gin.Context) {
	role, _ := c.Get("user_role")
	roleStr, _ := role.(string)

	var sppgID *uint
	var yayasanID *uint

	switch roleStr {
	case "kepala_yayasan":
		if yID, ok := c.Get("yayasan_id"); ok {
			id := yID.(uint)
			yayasanID = &id
		}
	case "superadmin", "admin_bgn":
		// Optional query param filtering
		if qYayasan := c.Query("yayasan_id"); qYayasan != "" {
			if id, err := strconv.ParseUint(qYayasan, 10, 32); err == nil {
				uid := uint(id)
				yayasanID = &uid
			}
		}
		if qSPPG := c.Query("sppg_id"); qSPPG != "" {
			if id, err := strconv.ParseUint(qSPPG, 10, 32); err == nil {
				uid := uint(id)
				sppgID = &uid
			}
		}
	default:
		// SPPG-level roles
		if sID, ok := c.Get("sppg_id"); ok {
			id := sID.(uint)
			sppgID = &id
		}
	}

	rabs, err := h.rabService.GetRABList(sppgID, yayasanID)
	if err != nil {
		log.Printf("GetRABList error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    rabs,
	})
}

// GetRABDetail handles GET /rab/:id — get RAB by ID with full preloads
func (h *RABHandler) GetRABDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	rab, err := h.rabService.GetRABByID(uint(id))
	if err != nil {
		if err == services.ErrRABNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "RAB_NOT_FOUND",
				"message":    "RAB tidak ditemukan",
			})
			return
		}
		log.Printf("GetRABDetail error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    rab,
	})
}

// UpdateRABRequest represents the request body for updating a RAB
type UpdateRABRequest struct {
	Items []UpdateRABItemRequest `json:"items" binding:"required,min=1"`
}

// UpdateRABItemRequest represents a single RAB item in the update request
type UpdateRABItemRequest struct {
	IngredientID          uint    `json:"ingredient_id" binding:"required"`
	Quantity              float64 `json:"quantity" binding:"required,gt=0"`
	Unit                  string  `json:"unit" binding:"required"`
	UnitPrice             float64 `json:"unit_price" binding:"gte=0"`
	RecommendedSupplierID *uint   `json:"recommended_supplier_id"`
}

// UpdateRAB handles PUT /rab/:id — update RAB items
func (h *RABHandler) UpdateRAB(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	var req UpdateRABRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	// Transform request to model items
	var items []models.RABItem
	for _, item := range req.Items {
		items = append(items, models.RABItem{
			IngredientID:          item.IngredientID,
			Quantity:              item.Quantity,
			Unit:                  item.Unit,
			UnitPrice:             item.UnitPrice,
			RecommendedSupplierID: item.RecommendedSupplierID,
		})
	}

	if err := h.rabService.UpdateRAB(uint(id), items); err != nil {
		if err == services.ErrRABNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "RAB_NOT_FOUND",
				"message":    "RAB tidak ditemukan",
			})
			return
		}
		if err == services.ErrRABNotEditable {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "RAB_NOT_EDITABLE",
				"message":    "RAB tidak dapat diedit pada status saat ini",
			})
			return
		}
		log.Printf("UpdateRAB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "RAB berhasil diperbarui",
	})
}

// ApproveSPPG handles POST /rab/:id/approve-sppg
func (h *RABHandler) ApproveSPPG(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	userID, _ := c.Get("user_id")

	if err := h.approvalEngine.ApproveByKepalaSPPG(uint(id), userID.(uint)); err != nil {
		if err == services.ErrRABNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "RAB_NOT_FOUND",
				"message":    "RAB tidak ditemukan",
			})
			return
		}
		if err == services.ErrRABNotApprovable {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "RAB_NOT_APPROVABLE",
				"message":    "RAB tidak dapat disetujui dari status saat ini",
			})
			return
		}
		log.Printf("ApproveSPPG error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "RAB berhasil disetujui oleh Kepala SPPG",
	})
}

// ApproveYayasan handles POST /rab/:id/approve-yayasan
func (h *RABHandler) ApproveYayasan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	userID, _ := c.Get("user_id")

	if err := h.approvalEngine.ApproveByKepalaYayasan(uint(id), userID.(uint)); err != nil {
		if err == services.ErrRABNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "RAB_NOT_FOUND",
				"message":    "RAB tidak ditemukan",
			})
			return
		}
		if err == services.ErrRABNotApprovable {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "RAB_NOT_APPROVABLE",
				"message":    "RAB tidak dapat disetujui dari status saat ini",
			})
			return
		}
		log.Printf("ApproveYayasan error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "RAB berhasil disetujui oleh Kepala Yayasan",
	})
}

// RejectRABRequest represents the request body for rejecting a RAB
type RejectRABRequest struct {
	Notes string `json:"notes" binding:"required"`
}

// RejectRAB handles POST /rab/:id/reject
func (h *RABHandler) RejectRAB(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	var req RejectRABRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	userID, _ := c.Get("user_id")

	if err := h.approvalEngine.RejectByKepalaYayasan(uint(id), userID.(uint), req.Notes); err != nil {
		if err == services.ErrRABNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "RAB_NOT_FOUND",
				"message":    "RAB tidak ditemukan",
			})
			return
		}
		if err == services.ErrRABInvalidStatus {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "RAB_INVALID_STATUS",
				"message":    "Status RAB tidak valid untuk operasi ini",
			})
			return
		}
		log.Printf("RejectRAB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "RAB berhasil ditolak",
	})
}

// ResubmitRAB handles POST /rab/:id/resubmit
func (h *RABHandler) ResubmitRAB(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	userID, _ := c.Get("user_id")

	if err := h.approvalEngine.ResubmitRAB(uint(id), userID.(uint)); err != nil {
		if err == services.ErrRABNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "RAB_NOT_FOUND",
				"message":    "RAB tidak ditemukan",
			})
			return
		}
		if err == services.ErrRABInvalidStatus {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "RAB_INVALID_STATUS",
				"message":    "Status RAB tidak valid untuk operasi ini",
			})
			return
		}
		log.Printf("ResubmitRAB error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "RAB berhasil disubmit ulang",
	})
}

// GetRABComparison handles GET /rab/:id/comparison
func (h *RABHandler) GetRABComparison(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	comparison, err := h.rabService.GetRABComparison(uint(id))
	if err != nil {
		if err == services.ErrRABNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "RAB_NOT_FOUND",
				"message":    "RAB tidak ditemukan",
			})
			return
		}
		log.Printf("GetRABComparison error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    comparison,
	})
}

// GetPOTracking handles GET /rab/:id/po-tracking
func (h *RABHandler) GetPOTracking(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	tracking, err := h.rabService.GetPOTracking(uint(id))
	if err != nil {
		if err == services.ErrRABNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "RAB_NOT_FOUND",
				"message":    "RAB tidak ditemukan",
			})
			return
		}
		log.Printf("GetPOTracking error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tracking,
	})
}
