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
	role, _ := c.Get("user_role")
	sppgID, _ := c.Get("sppg_id")
	yayasanID, _ := c.Get("yayasan_id")

	roleStr, ok := role.(string)
	if !ok {
		return db.Where("1 = 0")
	}

	switch roleStr {
	case "superadmin", "admin_bgn":
		if filterSPPG := c.Query("sppg_id"); filterSPPG != "" {
			return db.Where("sppg_id = ?", filterSPPG)
		}
		if filterYayasan := c.Query("yayasan_id"); filterYayasan != "" {
			return db.Where("sppg_id IN (?)",
				db.Session(&gorm.Session{NewDB: true}).
					Table("sppgs").Select("id").Where("yayasan_id = ?", filterYayasan))
		}
		return db
	case "kepala_yayasan":
		return db.Where("sppg_id IN (?)",
			db.Session(&gorm.Session{NewDB: true}).
				Table("sppgs").Select("id").Where("yayasan_id = ?", yayasanID))
	default:
		return db.Where("sppg_id = ?", sppgID)
	}
}

// autoInjectSPPGID sets the SPPG ID on a model struct for SPPG-level roles.
// Returns the sppg_id and true if injection was performed, 0 and false otherwise.
// Handlers should call this before Create operations to auto-fill sppg_id.
func autoInjectSPPGID(c *gin.Context) (uint, bool) {
	return middleware.GetTenantSPPGID(c)
}
