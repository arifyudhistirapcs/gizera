package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// sppgLevelRoles defines all roles that operate at the SPPG level
var sppgLevelRoles = map[string]bool{
	"kepala_sppg":      true,
	"akuntan":          true,
	"ahli_gizi":        true,
	"pengadaan":        true,
	"chef":             true,
	"packing":          true,
	"driver":           true,
	"asisten_lapangan": true,
	"kebersihan":       true,
}

// operationalPathPrefixes defines API path prefixes for operational data
// These are the endpoints where kepala_yayasan and admin_bgn should be read-only
var operationalPathPrefixes = []string{
	"/api/v1/recipes",
	"/api/v1/ingredients",
	"/api/v1/semi-finished",
	"/api/v1/menu-plans",
	"/api/v1/suppliers",
	"/api/v1/purchase-orders",
	"/api/v1/goods-receipts",
	"/api/v1/inventory",
	"/api/v1/stok-opname",
	"/api/v1/schools",
	"/api/v1/delivery-tasks",
	"/api/v1/pickup-tasks",
	"/api/v1/epod",
	"/api/v1/reviews",
	"/api/v1/ompreng",
	"/api/v1/employees",
	"/api/v1/attendance",
	"/api/v1/wifi-config",
	"/api/v1/gps-config",
	"/api/v1/assets",
	"/api/v1/cash-flow",
	"/api/v1/kds",
	"/api/v1/cleaning",
	"/api/v1/monitoring",
	"/api/v1/activity-tracker",
	"/api/v1/notifications",
}

// writeMethods defines HTTP methods that modify data
var writeMethods = map[string]bool{
	"POST":   true,
	"PUT":    true,
	"PATCH":  true,
	"DELETE": true,
}

// IsSPPGLevelRole returns true if the given role is an SPPG-level operational role
func IsSPPGLevelRole(role string) bool {
	return sppgLevelRoles[role]
}

// isOperationalEndpoint checks if the request path is an operational data endpoint
func isOperationalEndpoint(path string) bool {
	for _, prefix := range operationalPathPrefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

// TenantMiddleware extracts tenant context from Gin context (set by JWTAuth)
// and configures a tenant-scoped GORM DB instance in the request context.
//
// Behavior:
//   - SPPG-level roles: require sppg_id, set scoped DB with WHERE sppg_id = ?
//   - kepala_yayasan: require yayasan_id, set scoped DB with WHERE sppg_id IN (SELECT id FROM sppgs WHERE yayasan_id = ?)
//   - admin_bgn/superadmin: set unscoped DB (with optional query param filtering)
//   - Fail-closed: reject request if tenant context extraction fails for SPPG roles
//   - Read-only enforcement: reject write operations on operational data for kepala_yayasan and admin_bgn
func TenantMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"error_code": "UNAUTHORIZED",
				"message":    "Konteks tenant tidak dapat diekstrak: peran pengguna tidak ditemukan",
			})
			c.Abort()
			return
		}

		role, ok := roleVal.(string)
		if !ok || role == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"error_code": "UNAUTHORIZED",
				"message":    "Konteks tenant tidak dapat diekstrak: peran pengguna tidak valid",
			})
			c.Abort()
			return
		}

		switch {
		case role == "superadmin":
			// Superadmin: full access, no tenant filter, optional query param filtering
			c.Set("tenant_db", db)
			c.Next()

		case role == "admin_bgn":
			// Admin BGN: read-only on operational data, optional query param filtering
			if err := enforceReadOnly(c, role); err != nil {
				return
			}
			c.Set("tenant_db", db)
			c.Next()

		case role == "kepala_yayasan":
			// Kepala Yayasan: read-only on operational data, scoped to yayasan's SPPGs
			if err := enforceReadOnly(c, role); err != nil {
				return
			}
			yayasanIDVal, yExists := c.Get("yayasan_id")
			if !yExists {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success":    false,
					"error_code": "UNAUTHORIZED",
					"message":    "Konteks tenant tidak dapat diekstrak: yayasan_id tidak ditemukan",
				})
				c.Abort()
				return
			}
			yayasanID, ok := yayasanIDVal.(uint)
			if !ok || yayasanID == 0 {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success":    false,
					"error_code": "UNAUTHORIZED",
					"message":    "Konteks tenant tidak dapat diekstrak: yayasan_id tidak valid",
				})
				c.Abort()
				return
			}
			c.Set("tenant_db", db)
			c.Next()

		case IsSPPGLevelRole(role):
			// SPPG-level roles: require sppg_id, strict tenant isolation
			sppgIDVal, sExists := c.Get("sppg_id")
			if !sExists {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success":    false,
					"error_code": "UNAUTHORIZED",
					"message":    "Konteks tenant tidak dapat diekstrak: sppg_id tidak ditemukan",
				})
				c.Abort()
				return
			}
			sppgID, ok := sppgIDVal.(uint)
			if !ok || sppgID == 0 {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success":    false,
					"error_code": "UNAUTHORIZED",
					"message":    "Konteks tenant tidak dapat diekstrak: sppg_id tidak valid",
				})
				c.Abort()
				return
			}
			c.Set("tenant_db", db)
			c.Next()

		default:
			// Fail-closed: unknown role, reject request
			c.JSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"error_code": "UNAUTHORIZED",
				"message":    "Konteks tenant tidak dapat diekstrak: peran tidak dikenali",
			})
			c.Abort()
			return
		}
	}
}

// TenantScope returns a GORM scope function that filters data based on the
// tenant context of the current request. This should be used by handlers/services
// when querying operational data.
//
// Filtering behavior:
//   - superadmin/admin_bgn: no filter (bypass), with optional query param filtering
//   - kepala_yayasan: WHERE sppg_id IN (SELECT id FROM sppgs WHERE yayasan_id = ?)
//   - SPPG-level roles: WHERE sppg_id = ?
func TenantScope(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	role, _ := c.Get("user_role")
	sppgID, _ := c.Get("sppg_id")
	yayasanID, _ := c.Get("yayasan_id")

	return func(db *gorm.DB) *gorm.DB {
		roleStr, ok := role.(string)
		if !ok {
			// Fail-closed: if role cannot be determined, return impossible condition
			return db.Where("1 = 0")
		}

		switch roleStr {
		case "superadmin", "admin_bgn":
			// No tenant filter — access all data
			// But check query params for optional filtering
			if filterSPPG := c.Query("sppg_id"); filterSPPG != "" {
				return db.Where("sppg_id = ?", filterSPPG)
			}
			if filterYayasan := c.Query("yayasan_id"); filterYayasan != "" {
				// Sub-query: sppg_id IN (SELECT id FROM sppgs WHERE yayasan_id = ?)
				return db.Where("sppg_id IN (?)",
					db.Session(&gorm.Session{NewDB: true}).
						Table("sppgs").Select("id").Where("yayasan_id = ?", filterYayasan))
			}
			return db

		case "kepala_yayasan":
			// Filter based on all SPPGs under the Yayasan
			return db.Where("sppg_id IN (?)",
				db.Session(&gorm.Session{NewDB: true}).
					Table("sppgs").Select("id").Where("yayasan_id = ?", yayasanID))

		default:
			// SPPG-level roles — filter by sppg_id
			return db.Where("sppg_id = ?", sppgID)
		}
	}
}

// AutoInjectSPPGID returns a GORM callback scope that automatically sets sppg_id
// on new records for SPPG-level roles. This should be used as a GORM scope
// or called before Create operations.
func AutoInjectSPPGID(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		roleVal, _ := c.Get("user_role")
		role, ok := roleVal.(string)
		if !ok {
			return db
		}

		if IsSPPGLevelRole(role) {
			sppgIDVal, exists := c.Get("sppg_id")
			if exists {
				if sppgID, ok := sppgIDVal.(uint); ok && sppgID > 0 {
					// Use Statement.SetColumn to inject sppg_id on create
					db.Statement.SetColumn("sppg_id", sppgID)
				}
			}
		}
		return db
	}
}

// GetTenantSPPGID extracts the sppg_id from the Gin context for SPPG-level roles.
// Returns the sppg_id and true if found, or 0 and false otherwise.
// This is a helper for handlers that need to auto-inject sppg_id on INSERT.
func GetTenantSPPGID(c *gin.Context) (uint, bool) {
	roleVal, _ := c.Get("user_role")
	role, ok := roleVal.(string)
	if !ok {
		return 0, false
	}

	if IsSPPGLevelRole(role) {
		sppgIDVal, exists := c.Get("sppg_id")
		if exists {
			if sppgID, ok := sppgIDVal.(uint); ok && sppgID > 0 {
				return sppgID, true
			}
		}
	}
	return 0, false
}

// enforceReadOnly checks if the current request is a write operation on an
// operational data endpoint and rejects it for read-only roles (kepala_yayasan, admin_bgn).
// Returns an error (non-nil) if the request was rejected, nil if allowed.
func enforceReadOnly(c *gin.Context, role string) error {
	method := c.Request.Method
	path := c.Request.URL.Path

	if writeMethods[method] && isOperationalEndpoint(path) {
		c.JSON(http.StatusForbidden, gin.H{
			"success":    false,
			"error_code": "FORBIDDEN",
			"message":    "Peran " + role + " tidak diizinkan melakukan operasi tulis pada data operasional",
		})
		c.Abort()
		return errReadOnlyViolation
	}
	return nil
}

// errReadOnlyViolation is a sentinel error for read-only enforcement
var errReadOnlyViolation = &readOnlyError{}

type readOnlyError struct{}

func (e *readOnlyError) Error() string {
	return "read-only violation: write operation on operational data not allowed"
}
