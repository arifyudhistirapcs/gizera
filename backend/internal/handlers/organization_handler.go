package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/erp-sppg/backend/internal/models"
	"github.com/erp-sppg/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// OrganizationHandler handles Yayasan and SPPG endpoints
type OrganizationHandler struct {
	yayasanService *services.YayasanService
	sppgService    *services.SPPGService
}

// NewOrganizationHandler creates a new OrganizationHandler
func NewOrganizationHandler(yayasanService *services.YayasanService, sppgService *services.SPPGService) *OrganizationHandler {
	return &OrganizationHandler{
		yayasanService: yayasanService,
		sppgService:    sppgService,
	}
}

// --- Request types ---

// CreateYayasanRequest represents the request body for creating a Yayasan
type CreateYayasanRequest struct {
	Nama            string  `json:"nama" binding:"required"`
	Alamat          string  `json:"alamat"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	NomorTelepon    string  `json:"nomor_telepon"`
	Email           string  `json:"email"`
	PenanggungJawab string  `json:"penanggung_jawab"`
	NPWP            string  `json:"npwp"`
}

// UpdateYayasanRequest represents the request body for updating a Yayasan
type UpdateYayasanRequest struct {
	Nama            *string  `json:"nama"`
	Alamat          *string  `json:"alamat"`
	Latitude        *float64 `json:"latitude"`
	Longitude       *float64 `json:"longitude"`
	NomorTelepon    *string  `json:"nomor_telepon"`
	Email           *string  `json:"email"`
	PenanggungJawab *string  `json:"penanggung_jawab"`
	NPWP            *string  `json:"npwp"`
}

// SetStatusRequest represents the request body for activating/deactivating
type SetStatusRequest struct {
	IsActive bool `json:"is_active"`
}

// CreateSPPGRequest represents the request body for creating a SPPG
type CreateSPPGRequest struct {
	Nama         string  `json:"nama" binding:"required"`
	Alamat       string  `json:"alamat"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	NomorTelepon string  `json:"nomor_telepon"`
	Email        string  `json:"email"`
	YayasanID    uint    `json:"yayasan_id" binding:"required"`
}

// UpdateSPPGRequest represents the request body for updating a SPPG
type UpdateSPPGRequest struct {
	Nama         *string  `json:"nama"`
	Alamat       *string  `json:"alamat"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	NomorTelepon *string  `json:"nomor_telepon"`
	Email        *string  `json:"email"`
}

// TransferSPPGRequest represents the request body for transferring a SPPG
type TransferSPPGRequest struct {
	YayasanID uint `json:"yayasan_id" binding:"required"`
}

// --- Helper ---

// getUserID extracts user_id from gin context
func getUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	uid, ok := userID.(uint)
	return uid, ok
}

// parseIDParam parses a uint ID from a URL parameter
func parseIDParam(c *gin.Context, param string) (uint, bool) {
	idStr := c.Param(param)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return 0, false
	}
	return uint(id), true
}

// mapServiceError maps known service errors to appropriate HTTP responses
func mapServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrYayasanNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"success":    false,
			"error_code": "NOT_FOUND",
			"message":    err.Error(),
		})
	case errors.Is(err, services.ErrSPPGNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"success":    false,
			"error_code": "NOT_FOUND",
			"message":    err.Error(),
		})
	case errors.Is(err, services.ErrDuplicateYayasanKode),
		errors.Is(err, services.ErrDuplicateYayasanEmail),
		errors.Is(err, services.ErrDuplicateYayasanNPWP),
		errors.Is(err, services.ErrDuplicateSPPGKode),
		errors.Is(err, services.ErrDuplicateSPPGEmail):
		c.JSON(http.StatusConflict, gin.H{
			"success":    false,
			"error_code": "DUPLICATE",
			"message":    err.Error(),
		})
	case errors.Is(err, services.ErrInvalidYayasanID):
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_YAYASAN",
			"message":    err.Error(),
		})
	case errors.Is(err, services.ErrSPPGUnderInactiveYayasan),
		errors.Is(err, services.ErrYayasanInactive):
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "YAYASAN_INACTIVE",
			"message":    err.Error(),
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
	}
}

// =====================
// Yayasan Endpoints
// =====================

// CreateYayasan handles POST /api/v1/organizations/yayasan
func (h *OrganizationHandler) CreateYayasan(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success":    false,
			"error_code": "UNAUTHORIZED",
			"message":    "Autentikasi diperlukan",
		})
		return
	}

	var req CreateYayasanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	yayasan := &models.Yayasan{
		Nama:            req.Nama,
		Alamat:          req.Alamat,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		NomorTelepon:    req.NomorTelepon,
		Email:           req.Email,
		PenanggungJawab: req.PenanggungJawab,
		NPWP:            req.NPWP,
	}

	if err := h.yayasanService.Create(yayasan, userID); err != nil {
		mapServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Yayasan berhasil dibuat",
		"data":    yayasan,
	})
}

// GetAllYayasan handles GET /api/v1/organizations/yayasan
func (h *OrganizationHandler) GetAllYayasan(c *gin.Context) {
	yayasans, err := h.yayasanService.GetAll()
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
		"data":    yayasans,
	})
}

// GetYayasanByID handles GET /api/v1/organizations/yayasan/:id
func (h *OrganizationHandler) GetYayasanByID(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	yayasan, err := h.yayasanService.GetByID(id)
	if err != nil {
		mapServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    yayasan,
	})
}

// UpdateYayasan handles PUT /api/v1/organizations/yayasan/:id
func (h *OrganizationHandler) UpdateYayasan(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success":    false,
			"error_code": "UNAUTHORIZED",
			"message":    "Autentikasi diperlukan",
		})
		return
	}

	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req UpdateYayasanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	updates := make(map[string]interface{})
	if req.Nama != nil {
		updates["nama"] = *req.Nama
	}
	if req.Alamat != nil {
		updates["alamat"] = *req.Alamat
	}
	if req.Latitude != nil {
		updates["latitude"] = *req.Latitude
	}
	if req.Longitude != nil {
		updates["longitude"] = *req.Longitude
	}
	if req.NomorTelepon != nil {
		updates["nomor_telepon"] = *req.NomorTelepon
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.PenanggungJawab != nil {
		updates["penanggung_jawab"] = *req.PenanggungJawab
	}
	if req.NPWP != nil {
		updates["npwp"] = *req.NPWP
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Tidak ada data yang diperbarui",
		})
		return
	}

	updated, err := h.yayasanService.Update(id, updates, userID)
	if err != nil {
		mapServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Yayasan berhasil diperbarui",
		"data":    updated,
	})
}

// SetYayasanStatus handles PATCH /api/v1/organizations/yayasan/:id/status
func (h *OrganizationHandler) SetYayasanStatus(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success":    false,
			"error_code": "UNAUTHORIZED",
			"message":    "Autentikasi diperlukan",
		})
		return
	}

	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req SetStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	updated, err := h.yayasanService.SetStatus(id, req.IsActive, userID)
	if err != nil {
		mapServiceError(c, err)
		return
	}

	status := "diaktifkan"
	if !req.IsActive {
		status = "dinonaktifkan"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Yayasan berhasil " + status,
		"data":    updated,
	})
}

// =====================
// SPPG Endpoints
// =====================

// CreateSPPG handles POST /api/v1/organizations/sppg
func (h *OrganizationHandler) CreateSPPG(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success":    false,
			"error_code": "UNAUTHORIZED",
			"message":    "Autentikasi diperlukan",
		})
		return
	}

	var req CreateSPPGRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	sppg := &models.SPPG{
		Nama:         req.Nama,
		Alamat:       req.Alamat,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		NomorTelepon: req.NomorTelepon,
		Email:        req.Email,
		YayasanID:    req.YayasanID,
	}

	if err := h.sppgService.Create(sppg, userID); err != nil {
		mapServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "SPPG berhasil dibuat",
		"data":    sppg,
	})
}

// GetAllSPPG handles GET /api/v1/organizations/sppg
func (h *OrganizationHandler) GetAllSPPG(c *gin.Context) {
	sppgs, err := h.sppgService.GetAll()
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
		"data":    sppgs,
	})
}

// GetSPPGByID handles GET /api/v1/organizations/sppg/:id
func (h *OrganizationHandler) GetSPPGByID(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	sppg, err := h.sppgService.GetByID(id)
	if err != nil {
		mapServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    sppg,
	})
}

// UpdateSPPG handles PUT /api/v1/organizations/sppg/:id
func (h *OrganizationHandler) UpdateSPPG(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success":    false,
			"error_code": "UNAUTHORIZED",
			"message":    "Autentikasi diperlukan",
		})
		return
	}

	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req UpdateSPPGRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	updates := make(map[string]interface{})
	if req.Nama != nil {
		updates["nama"] = *req.Nama
	}
	if req.Alamat != nil {
		updates["alamat"] = *req.Alamat
	}
	if req.Latitude != nil {
		updates["latitude"] = *req.Latitude
	}
	if req.Longitude != nil {
		updates["longitude"] = *req.Longitude
	}
	if req.NomorTelepon != nil {
		updates["nomor_telepon"] = *req.NomorTelepon
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Tidak ada data yang diperbarui",
		})
		return
	}

	updated, err := h.sppgService.Update(id, updates, userID)
	if err != nil {
		mapServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "SPPG berhasil diperbarui",
		"data":    updated,
	})
}

// SetSPPGStatus handles PATCH /api/v1/organizations/sppg/:id/status
func (h *OrganizationHandler) SetSPPGStatus(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success":    false,
			"error_code": "UNAUTHORIZED",
			"message":    "Autentikasi diperlukan",
		})
		return
	}

	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req SetStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	updated, err := h.sppgService.SetStatus(id, req.IsActive, userID)
	if err != nil {
		mapServiceError(c, err)
		return
	}

	status := "diaktifkan"
	if !req.IsActive {
		status = "dinonaktifkan"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "SPPG berhasil " + status,
		"data":    updated,
	})
}

// TransferSPPG handles PUT /api/v1/organizations/sppg/:id/transfer
func (h *OrganizationHandler) TransferSPPG(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success":    false,
			"error_code": "UNAUTHORIZED",
			"message":    "Autentikasi diperlukan",
		})
		return
	}

	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req TransferSPPGRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	updated, err := h.sppgService.Transfer(id, req.YayasanID, userID)
	if err != nil {
		mapServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "SPPG berhasil dipindahkan ke Yayasan baru",
		"data":    updated,
	})
}
