package handlers

import (
	"errors"
	"net/http"

	"github.com/erp-sppg/backend/internal/models"
	"github.com/erp-sppg/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProvisioningHandler handles user provisioning endpoints
type ProvisioningHandler struct {
	provisioningService *services.ProvisioningService
	db                  *gorm.DB
}

// NewProvisioningHandler creates a new ProvisioningHandler
func NewProvisioningHandler(provisioningService *services.ProvisioningService, db *gorm.DB) *ProvisioningHandler {
	return &ProvisioningHandler{
		provisioningService: provisioningService,
		db:                  db,
	}
}

// getRequester loads the full User object (with SPPG/Yayasan relations) from the DB
// using the user_id stored in the Gin context by the auth middleware.
func (h *ProvisioningHandler) getRequester(c *gin.Context) (*models.User, bool) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success":    false,
			"error_code": "UNAUTHORIZED",
			"message":    "Autentikasi diperlukan",
		})
		return nil, false
	}

	var user models.User
	if err := h.db.Preload("SPPG").Preload("Yayasan").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success":    false,
			"error_code": "UNAUTHORIZED",
			"message":    "Pengguna tidak ditemukan",
		})
		return nil, false
	}

	return &user, true
}

// mapProvisioningError maps provisioning service errors to HTTP responses
func mapProvisioningError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrProvisioningForbidden):
		c.JSON(http.StatusForbidden, gin.H{
			"success":    false,
			"error_code": "FORBIDDEN",
			"message":    err.Error(),
		})
	case errors.Is(err, services.ErrRoleNotAllowed):
		c.JSON(http.StatusForbidden, gin.H{
			"success":    false,
			"error_code": "ROLE_NOT_ALLOWED",
			"message":    err.Error(),
		})
	case errors.Is(err, services.ErrSPPGNotInYayasan):
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "SPPG_NOT_IN_YAYASAN",
			"message":    err.Error(),
		})
	case errors.Is(err, services.ErrSPPGRequired):
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    err.Error(),
		})
	case errors.Is(err, services.ErrYayasanRequired):
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    err.Error(),
		})
	case errors.Is(err, services.ErrPasswordRequired):
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    err.Error(),
		})
	case errors.Is(err, services.ErrDuplicateNIK):
		c.JSON(http.StatusConflict, gin.H{
			"success":    false,
			"error_code": "DUPLICATE",
			"message":    err.Error(),
		})
	case errors.Is(err, services.ErrDuplicateEmail):
		c.JSON(http.StatusConflict, gin.H{
			"success":    false,
			"error_code": "DUPLICATE",
			"message":    err.Error(),
		})
	case errors.Is(err, services.ErrProvisioningUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"success":    false,
			"error_code": "NOT_FOUND",
			"message":    err.Error(),
		})
	case errors.Is(err, services.ErrSPPGNotFound):
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "SPPG_NOT_FOUND",
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

// CreateUser handles POST /api/v1/users
func (h *ProvisioningHandler) CreateUser(c *gin.Context) {
	requester, ok := h.getRequester(c)
	if !ok {
		return
	}

	var req services.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	user, err := h.provisioningService.CreateUser(&req, requester)
	if err != nil {
		mapProvisioningError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Pengguna berhasil dibuat",
		"data":    user,
	})
}

// GetUsers handles GET /api/v1/users
func (h *ProvisioningHandler) GetUsers(c *gin.Context) {
	requester, ok := h.getRequester(c)
	if !ok {
		return
	}

	users, err := h.provisioningService.GetUsers(requester)
	if err != nil {
		mapProvisioningError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
	})
}

// GetUserByID handles GET /api/v1/users/:id
func (h *ProvisioningHandler) GetUserByID(c *gin.Context) {
	requester, ok := h.getRequester(c)
	if !ok {
		return
	}

	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	user, err := h.provisioningService.GetUserByID(id, requester)
	if err != nil {
		mapProvisioningError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

// UpdateUser handles PUT /api/v1/users/:id
func (h *ProvisioningHandler) UpdateUser(c *gin.Context) {
	requester, ok := h.getRequester(c)
	if !ok {
		return
	}

	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req services.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	user, err := h.provisioningService.UpdateUser(id, &req, requester)
	if err != nil {
		mapProvisioningError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pengguna berhasil diperbarui",
		"data":    user,
	})
}

// SetUserStatus handles PATCH /api/v1/users/:id/status
func (h *ProvisioningHandler) SetUserStatus(c *gin.Context) {
	requester, ok := h.getRequester(c)
	if !ok {
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

	user, err := h.provisioningService.SetUserStatus(id, req.IsActive, requester)
	if err != nil {
		mapProvisioningError(c, err)
		return
	}

	status := "diaktifkan"
	if !req.IsActive {
		status = "dinonaktifkan"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pengguna berhasil " + status,
		"data":    user,
	})
}
