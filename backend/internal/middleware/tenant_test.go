package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	return db
}

// --- TenantMiddleware Tests ---

func TestTenantMiddleware_FailClosed_NoRole(t *testing.T) {
	db := setupTestDB(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/recipes", nil)

	// No user_role set — should fail closed
	handler := TenantMiddleware(db)
	handler(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
	if !c.IsAborted() {
		t.Error("expected request to be aborted")
	}
}

func TestTenantMiddleware_FailClosed_EmptyRole(t *testing.T) {
	db := setupTestDB(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/recipes", nil)
	c.Set("user_role", "")

	handler := TenantMiddleware(db)
	handler(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestTenantMiddleware_FailClosed_UnknownRole(t *testing.T) {
	db := setupTestDB(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/recipes", nil)
	c.Set("user_role", "unknown_role")

	handler := TenantMiddleware(db)
	handler(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestTenantMiddleware_Superadmin_Passes(t *testing.T) {
	db := setupTestDB(t)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/recipes", nil)
	c.Set("user_role", "superadmin")

	passed := false
	r.Use(func(c *gin.Context) {
		c.Set("user_role", "superadmin")
	})
	r.Use(TenantMiddleware(db))
	r.GET("/api/v1/recipes", func(c *gin.Context) {
		passed = true
		_, exists := c.Get("tenant_db")
		if !exists {
			t.Error("expected tenant_db to be set")
		}
		c.Status(http.StatusOK)
	})

	w = httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/recipes", nil)
	r.ServeHTTP(w, req)

	if !passed {
		t.Error("expected handler to be called for superadmin")
	}
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestTenantMiddleware_SPPGRole_NoSPPGID_Fails(t *testing.T) {
	db := setupTestDB(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/recipes", nil)
	c.Set("user_role", "chef")
	// No sppg_id set

	handler := TenantMiddleware(db)
	handler(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestTenantMiddleware_SPPGRole_WithSPPGID_Passes(t *testing.T) {
	db := setupTestDB(t)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	passed := false
	r.Use(func(c *gin.Context) {
		c.Set("user_role", "chef")
		c.Set("sppg_id", uint(1))
	})
	r.Use(TenantMiddleware(db))
	r.GET("/api/v1/recipes", func(c *gin.Context) {
		passed = true
		_, exists := c.Get("tenant_db")
		if !exists {
			t.Error("expected tenant_db to be set")
		}
		c.Status(http.StatusOK)
	})

	w = httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/recipes", nil)
	r.ServeHTTP(w, req)

	if !passed {
		t.Error("expected handler to be called for chef with sppg_id")
	}
}

func TestTenantMiddleware_KepalaYayasan_NoYayasanID_Fails(t *testing.T) {
	db := setupTestDB(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/recipes", nil)
	c.Set("user_role", "kepala_yayasan")
	// No yayasan_id set

	handler := TenantMiddleware(db)
	handler(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestTenantMiddleware_KepalaYayasan_WithYayasanID_Passes(t *testing.T) {
	db := setupTestDB(t)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	passed := false
	r.Use(func(c *gin.Context) {
		c.Set("user_role", "kepala_yayasan")
		c.Set("yayasan_id", uint(1))
	})
	r.Use(TenantMiddleware(db))
	r.GET("/api/v1/recipes", func(c *gin.Context) {
		passed = true
		c.Status(http.StatusOK)
	})

	w = httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/recipes", nil)
	r.ServeHTTP(w, req)

	if !passed {
		t.Error("expected handler to be called for kepala_yayasan with yayasan_id")
	}
}

// --- Read-Only Enforcement Tests ---

func TestTenantMiddleware_KepalaYayasan_WriteOperational_Rejected(t *testing.T) {
	db := setupTestDB(t)

	methods := []string{"POST", "PUT", "PATCH", "DELETE"}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)

			r.Use(func(c *gin.Context) {
				c.Set("user_role", "kepala_yayasan")
				c.Set("yayasan_id", uint(1))
			})
			r.Use(TenantMiddleware(db))
			r.Handle(method, "/api/v1/recipes", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			w = httptest.NewRecorder()
			req := httptest.NewRequest(method, "/api/v1/recipes", nil)
			r.ServeHTTP(w, req)

			if w.Code != http.StatusForbidden {
				t.Errorf("expected 403 for %s on operational endpoint, got %d", method, w.Code)
			}
		})
	}
}

func TestTenantMiddleware_AdminBGN_WriteOperational_Rejected(t *testing.T) {
	db := setupTestDB(t)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.Use(func(c *gin.Context) {
		c.Set("user_role", "admin_bgn")
	})
	r.Use(TenantMiddleware(db))
	r.POST("/api/v1/recipes", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w = httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/recipes", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 for admin_bgn POST on operational endpoint, got %d", w.Code)
	}
}

func TestTenantMiddleware_AdminBGN_ReadOperational_Allowed(t *testing.T) {
	db := setupTestDB(t)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	passed := false
	r.Use(func(c *gin.Context) {
		c.Set("user_role", "admin_bgn")
	})
	r.Use(TenantMiddleware(db))
	r.GET("/api/v1/recipes", func(c *gin.Context) {
		passed = true
		c.Status(http.StatusOK)
	})

	w = httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/recipes", nil)
	r.ServeHTTP(w, req)

	if !passed {
		t.Error("expected GET to be allowed for admin_bgn on operational endpoint")
	}
}

func TestTenantMiddleware_Superadmin_WriteOperational_Allowed(t *testing.T) {
	db := setupTestDB(t)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	passed := false
	r.Use(func(c *gin.Context) {
		c.Set("user_role", "superadmin")
	})
	r.Use(TenantMiddleware(db))
	r.POST("/api/v1/recipes", func(c *gin.Context) {
		passed = true
		c.Status(http.StatusOK)
	})

	w = httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/recipes", nil)
	r.ServeHTTP(w, req)

	if !passed {
		t.Error("expected POST to be allowed for superadmin on operational endpoint")
	}
}

func TestTenantMiddleware_KepalaYayasan_WriteNonOperational_Allowed(t *testing.T) {
	db := setupTestDB(t)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	passed := false
	r.Use(func(c *gin.Context) {
		c.Set("user_role", "kepala_yayasan")
		c.Set("yayasan_id", uint(1))
	})
	r.Use(TenantMiddleware(db))
	// /api/v1/users is NOT in operationalPathPrefixes — should be allowed
	r.POST("/api/v1/users", func(c *gin.Context) {
		passed = true
		c.Status(http.StatusOK)
	})

	w = httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/users", nil)
	r.ServeHTTP(w, req)

	if !passed {
		t.Error("expected POST to be allowed for kepala_yayasan on non-operational endpoint")
	}
}

// --- TenantScope Tests ---

func TestTenantScope_SPPGRole_FiltersBySPPGID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/recipes", nil)
	c.Set("user_role", "chef")
	c.Set("sppg_id", uint(42))

	scope := TenantScope(c)
	db := setupTestDB(t)
	stmt := db.Session(&gorm.Session{DryRun: true}).Scopes(scope).Find(&struct{}{}).Statement

	sql := stmt.SQL.String()
	if !containsSubstring(sql, "sppg_id") {
		t.Errorf("expected SQL to contain sppg_id filter, got: %s", sql)
	}
}

func TestTenantScope_Superadmin_NoFilter(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/recipes", nil)
	c.Set("user_role", "superadmin")

	scope := TenantScope(c)
	db := setupTestDB(t)
	stmt := db.Session(&gorm.Session{DryRun: true}).Scopes(scope).Find(&struct{}{}).Statement

	sql := stmt.SQL.String()
	if containsSubstring(sql, "sppg_id") {
		t.Errorf("expected no sppg_id filter for superadmin, got: %s", sql)
	}
}

func TestTenantScope_AdminBGN_WithSPPGQueryParam(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/recipes?sppg_id=5", nil)
	c.Set("user_role", "admin_bgn")

	scope := TenantScope(c)
	db := setupTestDB(t)
	stmt := db.Session(&gorm.Session{DryRun: true}).Scopes(scope).Find(&struct{}{}).Statement

	sql := stmt.SQL.String()
	if !containsSubstring(sql, "sppg_id") {
		t.Errorf("expected SQL to contain sppg_id filter from query param, got: %s", sql)
	}
}

func TestTenantScope_KepalaYayasan_FiltersSubquery(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/recipes", nil)
	c.Set("user_role", "kepala_yayasan")
	c.Set("yayasan_id", uint(3))

	scope := TenantScope(c)
	db := setupTestDB(t)
	stmt := db.Session(&gorm.Session{DryRun: true}).Scopes(scope).Find(&struct{}{}).Statement

	sql := stmt.SQL.String()
	if !containsSubstring(sql, "sppg_id IN") {
		t.Errorf("expected SQL to contain sppg_id IN subquery, got: %s", sql)
	}
}

func TestTenantScope_NoRole_FailClosed(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/recipes", nil)
	// No user_role set

	scope := TenantScope(c)
	db := setupTestDB(t)
	stmt := db.Session(&gorm.Session{DryRun: true}).Scopes(scope).Find(&struct{}{}).Statement

	sql := stmt.SQL.String()
	if !containsSubstring(sql, "1 = 0") {
		t.Errorf("expected fail-closed SQL (1 = 0), got: %s", sql)
	}
}

// --- IsSPPGLevelRole Tests ---

func TestIsSPPGLevelRole(t *testing.T) {
	sppgRoles := []string{"kepala_sppg", "akuntan", "ahli_gizi", "pengadaan", "chef", "packing", "driver", "asisten_lapangan", "kebersihan"}
	for _, role := range sppgRoles {
		if !IsSPPGLevelRole(role) {
			t.Errorf("expected %s to be SPPG-level role", role)
		}
	}

	nonSPPGRoles := []string{"superadmin", "admin_bgn", "kepala_yayasan", "unknown"}
	for _, role := range nonSPPGRoles {
		if IsSPPGLevelRole(role) {
			t.Errorf("expected %s to NOT be SPPG-level role", role)
		}
	}
}

// --- GetTenantSPPGID Tests ---

func TestGetTenantSPPGID_SPPGRole(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_role", "chef")
	c.Set("sppg_id", uint(42))

	id, ok := GetTenantSPPGID(c)
	if !ok {
		t.Error("expected ok=true for chef with sppg_id")
	}
	if id != 42 {
		t.Errorf("expected sppg_id=42, got %d", id)
	}
}

func TestGetTenantSPPGID_NonSPPGRole(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_role", "superadmin")

	_, ok := GetTenantSPPGID(c)
	if ok {
		t.Error("expected ok=false for superadmin")
	}
}

// --- Helper ---

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
