package handlers

import (
	"net/http"
	"time"

	"github.com/erp-sppg/backend/internal/models"
	"github.com/erp-sppg/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuditHandler handles audit trail endpoints
type AuditHandler struct {
	auditService *services.AuditTrailService
	db           *gorm.DB
}

// NewAuditHandler creates a new audit handler
func NewAuditHandler(db *gorm.DB) *AuditHandler {
	return &AuditHandler{
		auditService: services.NewAuditTrailService(db),
		db:           db,
	}
}

// GetAuditTrailRequest represents the query parameters for audit trail
type GetAuditTrailRequest struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
	UserID   uint   `form:"user_id"`
	Action   string `form:"action"`
	Entity   string `form:"entity"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Search    string `form:"search"`
}

// GetAuditTrail retrieves audit trail entries with filters.
// Scoping by role:
//   - superadmin: see all audit entries
//   - admin_bgn: see all audit entries (cross-tenant)
//   - kepala_yayasan: see audit entries where yayasan_id matches their yayasan,
//     or entries from SPPGs under their yayasan, or entries with no tenant context
//   - kepala_sppg: see audit entries where sppg_id matches their sppg
func (h *AuditHandler) GetAuditTrail(c *gin.Context) {
	var req GetAuditTrailRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Parameter tidak valid",
		})
		return
	}

	// Build filters
	filters := make(map[string]interface{})

	if req.UserID > 0 {
		filters["user_id"] = req.UserID
	}

	if req.Action != "" {
		filters["action"] = req.Action
	}

	if req.Entity != "" {
		filters["entity"] = req.Entity
	}

	if req.StartDate != "" {
		if startDate, err := time.Parse("2006-01-02", req.StartDate); err == nil {
			filters["start_date"] = startDate
		}
	}

	if req.EndDate != "" {
		if endDate, err := time.Parse("2006-01-02", req.EndDate); err == nil {
			endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			filters["end_date"] = endDate
		}
	}

	// Apply role-based scoping
	h.applyRoleScope(c, filters)

	// Calculate offset
	offset := (req.Page - 1) * req.PageSize

	// Get audit trail entries
	entries, total, err := h.auditService.GetAuditTrail(filters, req.PageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Gagal mengambil data audit trail",
		})
		return
	}

	// Transform entries for response
	var responseEntries []gin.H
	for _, entry := range entries {
		responseEntry := gin.H{
			"id":         entry.ID,
			"user_id":    entry.UserID,
			"timestamp":  entry.Timestamp,
			"action":     entry.Action,
			"entity":     entry.Entity,
			"entity_id":  entry.EntityID,
			"old_value":  entry.OldValue,
			"new_value":  entry.NewValue,
			"ip_address": entry.IPAddress,
			"sppg_id":    entry.SPPGID,
			"yayasan_id": entry.YayasanID,
			"level":      entry.Level,
			"user": gin.H{
				"id":        entry.User.ID,
				"nik":       entry.User.NIK,
				"full_name": entry.User.FullName,
				"email":     entry.User.Email,
				"role":      entry.User.Role,
			},
			"description": h.generateDescription(entry.Action, entry.Entity, entry.User.FullName),
		}
		responseEntries = append(responseEntries, responseEntry)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"data":       responseEntries,
		"total":      total,
		"page":       req.Page,
		"page_size":  req.PageSize,
		"total_pages": (total + int64(req.PageSize) - 1) / int64(req.PageSize),
	})
}

// applyRoleScope adds tenant-scoping filters based on the user's role.
func (h *AuditHandler) applyRoleScope(c *gin.Context, filters map[string]interface{}) {
	roleVal, exists := c.Get("user_role")
	if !exists {
		return
	}
	role, ok := roleVal.(string)
	if !ok {
		return
	}

	switch role {
	case "superadmin", "admin_bgn":
		// See all audit entries — no additional filter
		return

	case "kepala_yayasan":
		// Scope to entries from SPPGs under their yayasan (or entries with no tenant)
		yayasanIDVal, yExists := c.Get("yayasan_id")
		if !yExists {
			return
		}
		yayasanID, ok := yayasanIDVal.(uint)
		if !ok || yayasanID == 0 {
			return
		}
		// Get all SPPG IDs under this yayasan
		var sppgIDs []uint
		h.db.Model(&models.SPPG{}).Where("yayasan_id = ?", yayasanID).Pluck("id", &sppgIDs)
		filters["sppg_ids"] = sppgIDs

	case "kepala_sppg":
		// Scope to entries matching their sppg_id
		sppgIDVal, sExists := c.Get("sppg_id")
		if !sExists {
			return
		}
		sppgID, ok := sppgIDVal.(uint)
		if !ok || sppgID == 0 {
			return
		}
		filters["sppg_id"] = sppgID
	}
}

// generateDescription creates a human-readable description in Indonesian
func (h *AuditHandler) generateDescription(action, entity, userName string) string {
	actionMap := map[string]string{
		"create":               "membuat",
		"update":               "mengubah",
		"delete":               "menghapus",
		"login":                "masuk ke sistem",
		"logout":               "keluar dari sistem",
		"approve":              "menyetujui",
		"reject":               "menolak",
		"export":               "mengekspor",
		"cross_tenant_access":  "percobaan akses lintas tenant",
	}

	entityMap := map[string]string{
		"user":           "pengguna",
		"recipe":         "resep",
		"menu":           "menu",
		"supplier":       "supplier",
		"purchase_order": "purchase order",
		"inventory":      "inventori",
		"delivery_task":  "tugas pengiriman",
		"employee":       "karyawan",
		"asset":          "aset",
		"cash_flow":      "arus kas",
	}

	actionText := actionMap[action]
	if actionText == "" {
		actionText = action
	}

	entityText := entityMap[entity]
	if entityText == "" {
		entityText = entity
	}

	if action == "login" || action == "logout" {
		return userName + " " + actionText
	}

	if action == "cross_tenant_access" {
		return userName + " " + actionText + " pada " + entity
	}

	return userName + " " + actionText + " " + entityText
}

// GetAuditStats returns audit trail statistics
func (h *AuditHandler) GetAuditStats(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	filters := make(map[string]interface{})

	if startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filters["start_date"] = startDate
		}
	}

	if endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			filters["end_date"] = endDate
		}
	}

	// Apply role-based scoping
	h.applyRoleScope(c, filters)

	// Get total entries
	_, total, err := h.auditService.GetAuditTrail(filters, 1, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Gagal mengambil statistik audit trail",
		})
		return
	}

	// Get action breakdown
	actionStats := make(map[string]int64)
	actions := []string{"create", "update", "delete", "login", "logout", "approve", "reject", "export", "cross_tenant_access"}

	for _, action := range actions {
		actionFilters := make(map[string]interface{})
		for k, v := range filters {
			actionFilters[k] = v
		}
		actionFilters["action"] = action

		_, count, err := h.auditService.GetAuditTrail(actionFilters, 1, 0)
		if err == nil {
			actionStats[action] = count
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total_entries":    total,
			"action_breakdown": actionStats,
		},
	})
}
