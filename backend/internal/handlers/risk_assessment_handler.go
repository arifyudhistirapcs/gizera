package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/erp-sppg/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RiskAssessmentHandler handles risk assessment endpoints
type RiskAssessmentHandler struct {
	service *services.RiskAssessmentService
}

// NewRiskAssessmentHandler creates a new RiskAssessmentHandler
func NewRiskAssessmentHandler(service *services.RiskAssessmentService) *RiskAssessmentHandler {
	return &RiskAssessmentHandler{
		service: service,
	}
}

// --- Helper: extract yayasan_id from context ---

// getYayasanIDFromContext extracts yayasan_id from the gin context.
// For kepala_yayasan, it comes from JWT claims. For superadmin, returns 0 (no filter).
func getYayasanIDFromContext(c *gin.Context) uint {
	yayasanIDVal, exists := c.Get("yayasan_id")
	if !exists {
		return 0
	}
	yayasanID, ok := yayasanIDVal.(uint)
	if !ok {
		return 0
	}
	return yayasanID
}

// getUserIDFromContext extracts user_id from the gin context.
func getUserIDFromContext(c *gin.Context) (uint, bool) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	userID, ok := userIDVal.(uint)
	return userID, ok
}

// --- SOP Category Endpoints ---

// GetSOPCategories returns all SOP categories ordered by urutan
func (h *RiskAssessmentHandler) GetSOPCategories(c *gin.Context) {
	categories, err := h.service.GetSOPCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan saat mengambil kategori SOP",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    categories,
		"message": "Berhasil mengambil daftar kategori SOP",
	})
}

// CreateSOPCategoryRequest represents the request body for creating a SOP category
type CreateSOPCategoryRequest struct {
	Nama      string `json:"nama" binding:"required"`
	Deskripsi string `json:"deskripsi"`
}

// CreateSOPCategory creates a new SOP category
func (h *RiskAssessmentHandler) CreateSOPCategory(c *gin.Context) {
	var req CreateSOPCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
		})
		return
	}

	category, err := h.service.CreateSOPCategory(services.CreateSOPCategoryInput{
		Nama:      req.Nama,
		Deskripsi: req.Deskripsi,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Gagal membuat kategori SOP",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    category,
		"message": "Berhasil membuat kategori SOP",
	})
}

// UpdateSOPCategoryRequest represents the request body for updating a SOP category
type UpdateSOPCategoryRequest struct {
	Nama      string `json:"nama"`
	Deskripsi string `json:"deskripsi"`
	Urutan    *int   `json:"urutan"`
}

// UpdateSOPCategory updates an existing SOP category
func (h *RiskAssessmentHandler) UpdateSOPCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	var req UpdateSOPCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
		})
		return
	}

	category, err := h.service.UpdateSOPCategory(uint(id), services.UpdateSOPCategoryInput{
		Nama:      req.Nama,
		Deskripsi: req.Deskripsi,
		Urutan:    req.Urutan,
	})
	if err != nil {
		if err == services.ErrSOPCategoryNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "NOT_FOUND",
				"message":    "Kategori SOP tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Gagal memperbarui kategori SOP",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    category,
		"message": "Berhasil memperbarui kategori SOP",
	})
}

// --- SOP Checklist Item Endpoints ---

// GetSOPChecklistItems returns checklist items, optionally filtered by category
func (h *RiskAssessmentHandler) GetSOPChecklistItems(c *gin.Context) {
	var categoryID *uint
	if catIDStr := c.Query("category_id"); catIDStr != "" {
		catID, err := strconv.ParseUint(catIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INVALID_ID",
				"message":    "Category ID tidak valid",
			})
			return
		}
		uid := uint(catID)
		categoryID = &uid
	}

	items, err := h.service.GetSOPChecklistItems(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan saat mengambil checklist items",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    items,
		"message": "Berhasil mengambil daftar checklist items",
	})
}

// CreateSOPChecklistItemRequest represents the request body for creating a checklist item
type CreateSOPChecklistItemRequest struct {
	SOPCategoryID uint   `json:"sop_category_id" binding:"required"`
	Nama          string `json:"nama" binding:"required"`
	Deskripsi     string `json:"deskripsi"`
}

// CreateSOPChecklistItem creates a new SOP checklist item
func (h *RiskAssessmentHandler) CreateSOPChecklistItem(c *gin.Context) {
	var req CreateSOPChecklistItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
		})
		return
	}

	item, err := h.service.CreateSOPChecklistItem(services.CreateSOPChecklistItemInput{
		SOPCategoryID: req.SOPCategoryID,
		Nama:          req.Nama,
		Deskripsi:     req.Deskripsi,
	})
	if err != nil {
		if err == services.ErrSOPCategoryNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "NOT_FOUND",
				"message":    "Kategori SOP tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Gagal membuat checklist item",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    item,
		"message": "Berhasil membuat checklist item",
	})
}

// UpdateSOPChecklistItemRequest represents the request body for updating a checklist item
type UpdateSOPChecklistItemRequest struct {
	Nama      string `json:"nama"`
	Deskripsi string `json:"deskripsi"`
	Urutan    *int   `json:"urutan"`
}

// UpdateSOPChecklistItem updates an existing SOP checklist item
func (h *RiskAssessmentHandler) UpdateSOPChecklistItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	var req UpdateSOPChecklistItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
		})
		return
	}

	item, err := h.service.UpdateSOPChecklistItem(uint(id), services.UpdateSOPChecklistItemInput{
		Nama:      req.Nama,
		Deskripsi: req.Deskripsi,
		Urutan:    req.Urutan,
	})
	if err != nil {
		if err == services.ErrSOPChecklistItemNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "NOT_FOUND",
				"message":    "Checklist item tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Gagal memperbarui checklist item",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    item,
		"message": "Berhasil memperbarui checklist item",
	})
}

// SetSOPChecklistItemStatusRequest represents the request body for setting item status
type SetSOPChecklistItemStatusRequest struct {
	IsActive bool `json:"is_active"`
}

// SetSOPChecklistItemStatus activates or deactivates a checklist item
func (h *RiskAssessmentHandler) SetSOPChecklistItemStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	var req SetSOPChecklistItemStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
		})
		return
	}

	item, err := h.service.SetSOPChecklistItemStatus(uint(id), req.IsActive)
	if err != nil {
		if err == services.ErrSOPChecklistItemNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "NOT_FOUND",
				"message":    "Checklist item tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Gagal mengubah status checklist item",
		})
		return
	}

	statusMsg := "diaktifkan"
	if !req.IsActive {
		statusMsg = "dinonaktifkan"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    item,
		"message": fmt.Sprintf("Checklist item berhasil %s", statusMsg),
	})
}

// --- Form Endpoints ---

// CreateRiskAssessmentFormRequest represents the request body for creating a risk assessment form
type CreateRiskAssessmentFormRequest struct {
	SPPGID uint `json:"sppg_id" binding:"required"`
}

// CreateForm creates a new risk assessment form for a given SPPG
func (h *RiskAssessmentHandler) CreateForm(c *gin.Context) {
	var req CreateRiskAssessmentFormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid: sppg_id wajib diisi",
		})
		return
	}

	yayasanID := getYayasanIDFromContext(c)
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success":    false,
			"error_code": "UNAUTHORIZED",
			"message":    "User ID tidak ditemukan",
		})
		return
	}

	form, err := h.service.CreateForm(req.SPPGID, yayasanID, userID)
	if err != nil {
		if err == services.ErrRASPPGNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "SPPG_NOT_FOUND",
				"message":    err.Error(),
			})
			return
		}
		if err == services.ErrRASnapshotError {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success":    false,
				"error_code": "SNAPSHOT_ERROR",
				"message":    err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Gagal membuat formulir risk assessment",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    form,
		"message": "Berhasil membuat formulir risk assessment",
	})
}

// GetForms returns a paginated list of risk assessment forms
func (h *RiskAssessmentHandler) GetForms(c *gin.Context) {
	yayasanID := getYayasanIDFromContext(c)

	filter := services.FormFilter{
		YayasanID: yayasanID,
	}

	// Parse optional filters
	if sppgIDStr := c.Query("sppg_id"); sppgIDStr != "" {
		if sppgID, err := strconv.ParseUint(sppgIDStr, 10, 32); err == nil {
			filter.SPPGID = uint(sppgID)
		}
	}

	if status := c.Query("status"); status != "" {
		filter.Status = status
	}

	if riskLevel := c.Query("risk_level"); riskLevel != "" {
		filter.RiskLevel = riskLevel
	}

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		if t, err := time.Parse("2006-01-02", dateFrom); err == nil {
			filter.DateFrom = &t
		}
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		if t, err := time.Parse("2006-01-02", dateTo); err == nil {
			filter.DateTo = &t
		}
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			filter.Page = page
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
			filter.PageSize = pageSize
		}
	}

	forms, totalCount, err := h.service.ListForms(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan saat mengambil daftar formulir",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    forms,
		"message": "Berhasil mengambil daftar formulir risk assessment",
		"meta": gin.H{
			"total": totalCount,
			"page":  filter.Page,
			"size":  filter.PageSize,
		},
	})
}

// GetForm returns a single risk assessment form with all associations
func (h *RiskAssessmentHandler) GetForm(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	yayasanID := getYayasanIDFromContext(c)

	form, err := h.service.GetForm(uint(id), yayasanID)
	if err != nil {
		if err == services.ErrRAFormNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "FORM_NOT_FOUND",
				"message":    "Formulir risk assessment tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan saat mengambil formulir",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    form,
		"message": "Berhasil mengambil detail formulir risk assessment",
	})
}

// UpdateDraftRequest represents the request body for updating a draft form
type UpdateDraftRequest struct {
	Items []services.UpdateItemRequest `json:"items" binding:"required"`
}

// UpdateDraft updates items in a draft risk assessment form
func (h *RiskAssessmentHandler) UpdateDraft(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	var req UpdateDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
		})
		return
	}

	yayasanID := getYayasanIDFromContext(c)

	if err := h.service.UpdateDraft(uint(id), yayasanID, req.Items); err != nil {
		if err == services.ErrRAFormNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "FORM_NOT_FOUND",
				"message":    "Formulir risk assessment tidak ditemukan",
			})
			return
		}
		if err == services.ErrRAFormAlreadySubmitted {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "FORM_ALREADY_SUBMITTED",
				"message":    err.Error(),
			})
			return
		}
		if err == services.ErrRAInvalidScore {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INVALID_SCORE",
				"message":    err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "UPDATE_ERROR",
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil memperbarui draft formulir",
	})
}

// SubmitForm validates, calculates scores, and submits a risk assessment form
func (h *RiskAssessmentHandler) SubmitForm(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	yayasanID := getYayasanIDFromContext(c)

	form, err := h.service.SubmitForm(uint(id), yayasanID)
	if err != nil {
		if err == services.ErrRAFormNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "FORM_NOT_FOUND",
				"message":    "Formulir risk assessment tidak ditemukan",
			})
			return
		}
		if err == services.ErrRAFormAlreadySubmitted {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "FORM_ALREADY_SUBMITTED",
				"message":    err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), services.ErrRAIncompleteScores.Error()) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INCOMPLETE_SCORES",
				"message":    err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Gagal mengirimkan formulir risk assessment",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    form,
		"message": "Berhasil mengirimkan formulir risk assessment",
	})
}

// UploadEvidence handles file upload for evidence photos
func (h *RiskAssessmentHandler) UploadEvidence(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	yayasanID := getYayasanIDFromContext(c)

	// Verify form exists and belongs to tenant
	form, err := h.service.GetForm(uint(id), yayasanID)
	if err != nil {
		if err == services.ErrRAFormNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "FORM_NOT_FOUND",
				"message":    "Formulir risk assessment tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan saat memverifikasi formulir",
		})
		return
	}

	// Prevent upload to submitted forms
	if form.Status != "draft" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "FORM_ALREADY_SUBMITTED",
			"message":    "Tidak dapat mengunggah evidence ke formulir yang sudah disubmit",
		})
		return
	}

	// Handle file upload
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "UPLOAD_FAILED",
			"message":    "File foto tidak ditemukan",
		})
		return
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "FILE_TOO_LARGE",
			"message":    "Ukuran file terlalu besar. Maksimal 5MB.",
		})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	filename := fmt.Sprintf("evidence-%d-%s%s", id, uuid.New().String(), ext)
	uploadDir := filepath.Join("uploads", "risk-assessment")
	uploadPath := filepath.Join(uploadDir, filename)

	// Ensure directory exists
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Gagal membuat direktori upload",
		})
		return
	}

	// Save file
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "UPLOAD_FAILED",
			"message":    "Gagal menyimpan file evidence",
		})
		return
	}

	// Generate URL
	evidenceURL := fmt.Sprintf("/uploads/risk-assessment/%s", filename)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"evidence_url": evidenceURL,
		},
		"message": "Berhasil mengunggah foto evidence",
	})
}

// DeleteForm deletes a draft risk assessment form
func (h *RiskAssessmentHandler) DeleteForm(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error_code": "INVALID_ID", "message": "ID tidak valid"})
		return
	}
	yayasanID := getYayasanIDFromContext(c)
	if err := h.service.DeleteForm(uint(id), yayasanID); err != nil {
		if err == services.ErrRAFormNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error_code": "FORM_NOT_FOUND", "message": "Formulir tidak ditemukan"})
			return
		}
		if err == services.ErrRAFormNotDraft {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error_code": "FORM_NOT_DRAFT", "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error_code": "INTERNAL_ERROR", "message": "Gagal menghapus formulir"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Formulir draft berhasil dihapus"})
}

// GetSPPGList returns SPPGs under the user's Yayasan for risk assessment form creation
func (h *RiskAssessmentHandler) GetSPPGList(c *gin.Context) {
	yayasanID := getYayasanIDFromContext(c)

	sppgs, err := h.service.GetSPPGsByYayasan(yayasanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Gagal mengambil daftar SPPG",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    sppgs,
		"message": "Berhasil mengambil daftar SPPG",
	})
}

// GetStats returns aggregated statistics per SPPG
func (h *RiskAssessmentHandler) GetStats(c *gin.Context) {
	// Parse sppg_ids from query parameter (comma-separated)
	sppgIDsStr := c.Query("sppg_ids")
	if sppgIDsStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Parameter sppg_ids wajib diisi (comma-separated)",
		})
		return
	}

	var sppgIDs []uint
	for _, idStr := range strings.Split(sppgIDsStr, ",") {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INVALID_ID",
				"message":    fmt.Sprintf("SPPG ID tidak valid: %s", idStr),
			})
			return
		}
		sppgIDs = append(sppgIDs, uint(id))
	}

	stats, err := h.service.GetStats(sppgIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan saat mengambil statistik",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
		"message": "Berhasil mengambil statistik risk assessment",
	})
}
