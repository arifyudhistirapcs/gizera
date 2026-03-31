package handlers

import (
	"github.com/erp-sppg/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// getTenantScopedDB returns a GORM DB instance with tenant scoping applied.
// This should be used by handlers when creating scoped service instances
// for querying operational data.
//
// The returned DB will automatically filter data based on the user's role:
//   - SPPG-level roles: WHERE sppg_id = ?
//   - kepala_yayasan: WHERE sppg_id IN (SELECT id FROM sppgs WHERE yayasan_id = ?)
//   - admin_bgn/superadmin: no filter (with optional query param filtering)
func getTenantScopedDB(c *gin.Context, db *gorm.DB) *gorm.DB {
	// Apply tenant scope and return a cloneable DB so each query in the service
	// gets a fresh statement (prevents GORM session state accumulation).
	scopeFn := middleware.TenantScope(c)
	// Eagerly apply the scope to get the scoped DB, then wrap in a new session
	// so the service gets a fresh clone for each query.
	scopedDB := scopeFn(db.Session(&gorm.Session{NewDB: true}))
	// Return with clone=2 so each query clones the statement (preserving the scope's WHERE clauses)
	return scopedDB.Session(&gorm.Session{})
}

// autoInjectSPPGID sets the SPPG ID on a model struct for SPPG-level roles.
// Returns the sppg_id and true if injection was performed, 0 and false otherwise.
// Handlers should call this before Create operations to auto-fill sppg_id.
func autoInjectSPPGID(c *gin.Context) (uint, bool) {
	return middleware.GetTenantSPPGID(c)
}
