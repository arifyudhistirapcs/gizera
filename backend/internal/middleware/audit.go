package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"strings"
	"time"

	"github.com/erp-sppg/backend/internal/models"
	"github.com/erp-sppg/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuditTrail middleware records all create/update/delete actions
// with tenant context (sppg_id, yayasan_id) from the authenticated user.
func AuditTrail(db *gorm.DB) gin.HandlerFunc {
	auditService := services.NewAuditTrailService(db)

	return func(c *gin.Context) {
		// Only audit CUD operations
		method := c.Request.Method
		if method != "POST" && method != "PUT" && method != "PATCH" && method != "DELETE" {
			c.Next()
			return
		}

		// Get user ID from context
		userID, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}

		// Extract tenant context from Gin context (set by JWT auth middleware)
		sppgID := extractUintFromContext(c, "sppg_id")
		yayasanID := extractUintFromContext(c, "yayasan_id")

		// Read request body for old/new values
		var requestBody map[string]interface{}
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			json.Unmarshal(bodyBytes, &requestBody)
		}

		// Determine action type
		action := ""
		switch method {
		case "POST":
			action = "create"
		case "PUT", "PATCH":
			action = "update"
		case "DELETE":
			action = "delete"
		}

		// Extract entity from path
		path := c.Request.URL.Path
		parts := strings.Split(strings.Trim(path, "/"), "/")
		entity := ""
		entityID := ""

		// Try to extract entity name from path (e.g., /api/v1/recipes -> recipes)
		if len(parts) >= 3 {
			entity = parts[2]
		}

		// Try to extract entity ID from path (e.g., /api/v1/recipes/123 -> 123)
		if len(parts) >= 4 {
			entityID = parts[3]
		}

		// Get client IP
		ipAddress := c.ClientIP()

		// Process request
		c.Next()

		// Only record if request was successful (2xx status)
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			var oldValue, newValue interface{}

			switch action {
			case "create":
				newValue = requestBody
			case "update":
				newValue = requestBody
			case "delete":
				oldValue = requestBody
			}

			// Convert values to JSON strings
			oldJSON, err := json.Marshal(oldValue)
			if err != nil {
				oldJSON = []byte("{}")
			}
			newJSON, err := json.Marshal(newValue)
			if err != nil {
				newJSON = []byte("{}")
			}

			auditEntry := models.AuditTrail{
				UserID:    userID.(uint),
				Timestamp: time.Now(),
				Action:    action,
				Entity:    entity,
				EntityID:  entityID,
				OldValue:  string(oldJSON),
				NewValue:  string(newJSON),
				IPAddress: ipAddress,
				SPPGID:    sppgID,
				YayasanID: yayasanID,
				Level:     "info",
			}

			// Ignore errors to not affect main request
			db.Create(&auditEntry)
		}

		// Check if this was a cross-tenant access attempt (blocked by TenantMiddleware)
		// TenantMiddleware aborts with 401/403 for cross-tenant violations
		if c.IsAborted() && (c.Writer.Status() == 401 || c.Writer.Status() == 403) {
			RecordCrossTenantWarning(auditService, db, userID.(uint), sppgID, yayasanID, ipAddress, c.Request.URL.Path, method)
		}
	}
}

// RecordCrossTenantWarning logs a warning-level audit entry for cross-tenant access attempts.
func RecordCrossTenantWarning(auditService *services.AuditTrailService, db *gorm.DB, userID uint, sppgID, yayasanID *uint, ipAddress, path, method string) {
	log.Printf("[AUDIT WARNING] Cross-tenant access attempt by user %d on %s %s", userID, method, path)

	warningEntry := models.AuditTrail{
		UserID:    userID,
		Timestamp: time.Now(),
		Action:    "cross_tenant_access",
		Entity:    path,
		EntityID:  "",
		OldValue:  "",
		NewValue:  `{"method":"` + method + `","path":"` + path + `"}`,
		IPAddress: ipAddress,
		SPPGID:    sppgID,
		YayasanID: yayasanID,
		Level:     "warning",
	}

	if err := db.Create(&warningEntry).Error; err != nil {
		log.Printf("[AUDIT WARNING] Failed to record cross-tenant access warning: %v", err)
	}
}

// extractUintFromContext extracts a *uint value from the Gin context.
// Returns nil if the key doesn't exist or the value is zero.
func extractUintFromContext(c *gin.Context, key string) *uint {
	val, exists := c.Get(key)
	if !exists {
		return nil
	}
	switch v := val.(type) {
	case uint:
		if v == 0 {
			return nil
		}
		return &v
	case *uint:
		return v
	default:
		return nil
	}
}
