package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/erp-sppg/backend/internal/models"
	"github.com/erp-sppg/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupProvisioningTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&models.Yayasan{}, &models.SPPG{}, &models.User{}, &models.AuditTrail{})
	require.NoError(t, err)
	return db
}

func seedProvisioningData(t *testing.T, db *gorm.DB) (*models.Yayasan, *models.SPPG, *models.User) {
	y := models.Yayasan{Kode: "YYS-0001", Nama: "Yayasan Test", IsActive: true}
	require.NoError(t, db.Create(&y).Error)
	s := models.SPPG{Kode: "SPPG-0001", Nama: "SPPG Test", YayasanID: y.ID, IsActive: true}
	require.NoError(t, db.Create(&s).Error)

	// Create a superadmin user in DB so getRequester can load it
	sa := models.User{
		NIK: "SA-001", Email: "sa@test.com", PasswordHash: "hash",
		FullName: "Super Admin", Role: "superadmin", IsActive: true,
	}
	require.NoError(t, db.Create(&sa).Error)
	return &y, &s, &sa
}

func setupProvisioningRouter(db *gorm.DB) (*gin.Engine, *ProvisioningHandler) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := services.NewProvisioningService(db, nil)
	h := NewProvisioningHandler(svc, db)

	users := r.Group("/api/v1/users")
	{
		users.POST("", h.CreateUser)
		users.GET("", h.GetUsers)
		users.GET("/:id", h.GetUserByID)
		users.PUT("/:id", h.UpdateUser)
		users.PATCH("/:id/status", h.SetUserStatus)
	}
	return r, h
}

// withUserID is a test middleware that sets user_id in context
func withUserID(userID uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	}
}

func setupProvisioningRouterWithAuth(db *gorm.DB, userID uint) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(withUserID(userID))

	svc := services.NewProvisioningService(db, nil)
	h := NewProvisioningHandler(svc, db)

	users := r.Group("/api/v1/users")
	{
		users.POST("", h.CreateUser)
		users.GET("", h.GetUsers)
		users.GET("/:id", h.GetUserByID)
		users.PUT("/:id", h.UpdateUser)
		users.PATCH("/:id/status", h.SetUserStatus)
	}
	return r
}

// TestCreateUser_Success tests POST /api/v1/users returns 201 with valid data
func TestCreateUser_Success(t *testing.T) {
	db := setupProvisioningTestDB(t)
	_, sppg, sa := seedProvisioningData(t, db)
	r := setupProvisioningRouterWithAuth(db, sa.ID)

	body, _ := json.Marshal(map[string]interface{}{
		"nik":       "NIK-001",
		"email":     "chef@test.com",
		"password":  "password123",
		"full_name": "Chef Test",
		"role":      "chef",
		"sppg_id":   sppg.ID,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
	assert.Equal(t, "Pengguna berhasil dibuat", resp["message"])

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "chef", data["role"])
	assert.Equal(t, "NIK-001", data["nik"])
}

// TestCreateUser_Forbidden tests POST /api/v1/users returns 403 for admin_bgn
func TestCreateUser_Forbidden(t *testing.T) {
	db := setupProvisioningTestDB(t)
	seedProvisioningData(t, db)

	// Create admin_bgn user in DB
	bgn := models.User{
		NIK: "BGN-001", Email: "bgn@test.com", PasswordHash: "hash",
		FullName: "Admin BGN", Role: "admin_bgn", IsActive: true,
	}
	require.NoError(t, db.Create(&bgn).Error)

	r := setupProvisioningRouterWithAuth(db, bgn.ID)

	body, _ := json.Marshal(map[string]interface{}{
		"nik": "NIK-x", "email": "x@test.com", "password": "password123",
		"full_name": "X", "role": "chef",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// TestGetUsers_Success tests GET /api/v1/users returns user list
func TestGetUsers_Success(t *testing.T) {
	db := setupProvisioningTestDB(t)
	_, sppg, sa := seedProvisioningData(t, db)
	r := setupProvisioningRouterWithAuth(db, sa.ID)

	// Create a user first
	svc := services.NewProvisioningService(db, nil)
	svc.CreateUser(&services.CreateUserRequest{
		NIK: "NIK-u1", Email: "u1@test.com", Password: "password123",
		FullName: "User 1", Role: "chef", SPPGID: &sppg.ID,
	}, sa)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))

	data := resp["data"].([]interface{})
	// superadmin + created user
	assert.GreaterOrEqual(t, len(data), 1)
}

// TestGetUserByID_Success tests GET /api/v1/users/:id returns user detail
func TestGetUserByID_Success(t *testing.T) {
	db := setupProvisioningTestDB(t)
	_, sppg, sa := seedProvisioningData(t, db)
	r := setupProvisioningRouterWithAuth(db, sa.ID)

	svc := services.NewProvisioningService(db, nil)
	created, _ := svc.CreateUser(&services.CreateUserRequest{
		NIK: "NIK-u2", Email: "u2@test.com", Password: "password123",
		FullName: "User 2", Role: "chef", SPPGID: &sppg.ID,
	}, sa)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/"+uintToStr(created.ID), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "User 2", data["full_name"])
}

// TestGetUserByID_NotFound tests GET /api/v1/users/:id returns 404 for non-existent user
func TestGetUserByID_NotFound(t *testing.T) {
	db := setupProvisioningTestDB(t)
	_, _, sa := seedProvisioningData(t, db)
	r := setupProvisioningRouterWithAuth(db, sa.ID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/9999", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestUpdateUser_Success tests PUT /api/v1/users/:id updates user data
func TestUpdateUser_Success(t *testing.T) {
	db := setupProvisioningTestDB(t)
	_, sppg, sa := seedProvisioningData(t, db)
	r := setupProvisioningRouterWithAuth(db, sa.ID)

	svc := services.NewProvisioningService(db, nil)
	created, _ := svc.CreateUser(&services.CreateUserRequest{
		NIK: "NIK-u3", Email: "u3@test.com", Password: "password123",
		FullName: "User 3", Role: "chef", SPPGID: &sppg.ID,
	}, sa)

	body, _ := json.Marshal(map[string]interface{}{
		"full_name": "User 3 Updated",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/users/"+uintToStr(created.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "User 3 Updated", data["full_name"])
}

// TestSetUserStatus_Success tests PATCH /api/v1/users/:id/status toggles user status
func TestSetUserStatus_Success(t *testing.T) {
	db := setupProvisioningTestDB(t)
	_, sppg, sa := seedProvisioningData(t, db)
	r := setupProvisioningRouterWithAuth(db, sa.ID)

	svc := services.NewProvisioningService(db, nil)
	created, _ := svc.CreateUser(&services.CreateUserRequest{
		NIK: "NIK-u4", Email: "u4@test.com", Password: "password123",
		FullName: "User 4", Role: "chef", SPPGID: &sppg.ID,
	}, sa)

	// Deactivate
	body, _ := json.Marshal(map[string]interface{}{"is_active": false})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/api/v1/users/"+uintToStr(created.ID)+"/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp["message"], "dinonaktifkan")
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, false, data["is_active"])
}

// TestCreateUser_Unauthorized tests POST /api/v1/users returns 401 without auth
func TestCreateUser_Unauthorized(t *testing.T) {
	db := setupProvisioningTestDB(t)
	seedProvisioningData(t, db)

	// Router WITHOUT auth middleware (no user_id in context)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := services.NewProvisioningService(db, nil)
	h := NewProvisioningHandler(svc, db)
	r.POST("/api/v1/users", h.CreateUser)

	body, _ := json.Marshal(map[string]interface{}{
		"nik": "NIK-x", "email": "x@test.com", "password": "password123",
		"full_name": "X", "role": "chef",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// uintToStr converts uint to string for URL building
func uintToStr(id uint) string {
	return fmt.Sprintf("%d", id)
}
