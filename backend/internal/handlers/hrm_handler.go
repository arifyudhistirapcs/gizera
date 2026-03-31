package handlers

import (
	"bytes"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/erp-sppg/backend/internal/models"
	"github.com/erp-sppg/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HRMHandler handles Human Resource Management endpoints
type HRMHandler struct {
	db                *gorm.DB
	employeeService   *services.EmployeeService
	attendanceService *services.AttendanceService
	auditService      *services.AuditTrailService
}

// NewHRMHandler creates a new HRM handler
func NewHRMHandler(db *gorm.DB, authService *services.AuthService) *HRMHandler {
	employeeService := services.NewEmployeeService(db, authService)
	return &HRMHandler{
		db:                db,
		employeeService:   employeeService,
		attendanceService: services.NewAttendanceService(db, employeeService),
		auditService:      services.NewAuditTrailService(db),
	}
}

// Employee Endpoints

// CreateEmployeeRequest represents create employee request
type CreateEmployeeRequest struct {
	NIK         string `json:"nik" binding:"required"`
	FullName    string `json:"full_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number"`
	Position    string `json:"position" binding:"required"`
	Role        string `json:"role" binding:"required"`
	JoinDate    string `json:"join_date" binding:"required"`
	Password    string `json:"password" binding:"omitempty,min=6"`
	IsActive    bool   `json:"is_active"`
}

// CreateEmployee creates a new employee
func (h *HRMHandler) CreateEmployee(c *gin.Context) {
	var req CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[CREATE EMPLOYEE] Validation error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid: " + err.Error(),
		})
		return
	}

	// Parse join_date from string to time.Time
	joinDate, err := time.Parse("2006-01-02", req.JoinDate)
	if err != nil {
		log.Printf("[CREATE EMPLOYEE] Date parse error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_DATE",
			"message":    "Format tanggal tidak valid",
		})
		return
	}

	employee := &models.Employee{
		NIK:         req.NIK,
		FullName:    req.FullName,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Position:    req.Position,
		JoinDate:    joinDate,
		IsActive:    req.IsActive,
	}

	// Create employee with custom password if provided
	var user *models.User
	var password string
	
	if req.Password != "" {
		user, err = h.employeeService.CreateEmployeeWithPassword(employee, req.Role, req.Password)
		password = req.Password
	} else {
		user, password, err = h.employeeService.CreateEmployee(employee, req.Role)
	}
	if err != nil {
		if err == services.ErrDuplicateNIK {
			c.JSON(http.StatusConflict, gin.H{
				"success":    false,
				"error_code": "DUPLICATE_NIK",
				"message":    "NIK sudah terdaftar",
			})
			return
		}
		if err == services.ErrDuplicateEmail {
			c.JSON(http.StatusConflict, gin.H{
				"success":    false,
				"error_code": "DUPLICATE_EMAIL",
				"message":    "Email sudah terdaftar",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	// Record in audit trail
	userID, _ := c.Get("user_id")
	h.auditService.RecordAction(userID.(uint), "create", "employee", strconv.Itoa(int(employee.ID)), "", "", c.ClientIP())

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Karyawan berhasil dibuat",
		"data": gin.H{
			"employee": employee,
			"user":     user,
			"credentials": gin.H{
				"password": password,
			},
		},
	})
}

// GetEmployees retrieves all employees
func (h *HRMHandler) GetEmployees(c *gin.Context) {
	isActiveStr := c.Query("is_active")
	position := c.Query("position")

	var isActive *bool
	if isActiveStr != "" {
		val := isActiveStr == "true"
		isActive = &val
	}

	scopedService := h.employeeService.WithDB(getTenantScopedDB(c, h.db))
	employees, err := scopedService.GetAllEmployees(isActive, position)
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
		"data":    employees,
	})
}

// GetEmployeeByID retrieves an employee by ID
func (h *HRMHandler) GetEmployeeByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	scopedService := h.employeeService.WithDB(getTenantScopedDB(c, h.db))
	employee, err := scopedService.GetEmployeeByID(uint(id))
	if err != nil {
		if err == services.ErrEmployeeNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "EMPLOYEE_NOT_FOUND",
				"message":    "Karyawan tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    employee,
	})
}

// UpdateEmployee updates an employee
func (h *HRMHandler) UpdateEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
		})
		return
	}

	scopedService := h.employeeService.WithDB(getTenantScopedDB(c, h.db))
	employee, err := scopedService.UpdateEmployee(uint(id), updates)
	if err != nil {
		if err == services.ErrEmployeeNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "EMPLOYEE_NOT_FOUND",
				"message":    "Karyawan tidak ditemukan",
			})
			return
		}
		if err == services.ErrDuplicateNIK {
			c.JSON(http.StatusConflict, gin.H{
				"success":    false,
				"error_code": "DUPLICATE_NIK",
				"message":    "NIK sudah terdaftar",
			})
			return
		}
		if err == services.ErrDuplicateEmail {
			c.JSON(http.StatusConflict, gin.H{
				"success":    false,
				"error_code": "DUPLICATE_EMAIL",
				"message":    "Email sudah terdaftar",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	// Record in audit trail
	userID, _ := c.Get("user_id")
	h.auditService.RecordAction(userID.(uint), "update", "employee", strconv.Itoa(int(id)), "", "", c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Karyawan berhasil diperbarui",
		"data":    employee,
	})
}

// DeactivateEmployee deactivates an employee
func (h *HRMHandler) DeactivateEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	scopedService := h.employeeService.WithDB(getTenantScopedDB(c, h.db))
	if err := scopedService.DeactivateEmployee(uint(id)); err != nil {
		if err == services.ErrEmployeeNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "EMPLOYEE_NOT_FOUND",
				"message":    "Karyawan tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	// Record in audit trail
	userID, _ := c.Get("user_id")
	h.auditService.RecordAction(userID.(uint), "deactivate", "employee", strconv.Itoa(int(id)), "", "", c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Karyawan berhasil dinonaktifkan",
	})
}

// GetEmployeeStats retrieves employee statistics
func (h *HRMHandler) GetEmployeeStats(c *gin.Context) {
	scopedService := h.employeeService.WithDB(getTenantScopedDB(c, h.db))
	stats, err := scopedService.GetEmployeeStats()
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
		"data":    stats,
	})
}

// Attendance Endpoints

// CheckInRequest represents check-in request
type CheckInRequest struct {
	SSID      string   `json:"ssid"`
	BSSID     string   `json:"bssid"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

// CheckIn records employee check-in
func (h *HRMHandler) CheckIn(c *gin.Context) {
	var req CheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
		})
		return
	}

	// Get employee ID from user ID
	userID, _ := c.Get("user_id")
	employee, err := h.employeeService.GetEmployeeByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success":    false,
			"error_code": "EMPLOYEE_NOT_FOUND",
			"message":    "Data karyawan tidak ditemukan",
		})
		return
	}

	// Get client IP address
	clientIP := c.ClientIP()
	
	// Log untuk debugging
	log.Printf("[CHECK-IN] Employee ID: %d, Client IP: %s, SSID: %s, Lat: %v, Lng: %v", 
		employee.ID, clientIP, req.SSID, req.Latitude, req.Longitude)
	
	var ssid, bssid string
	var validationMethod string
	var validationDetails map[string]interface{}
	var isValidated bool
	
	// GPS validation only (for testing)
	if req.Latitude == nil || req.Longitude == nil {
		log.Printf("[CHECK-IN] GPS coordinates not provided")
		c.JSON(http.StatusForbidden, gin.H{
			"success":    false,
			"error_code": "GPS_REQUIRED",
			"message":    "Koordinat GPS diperlukan untuk check-in. Pastikan GPS aktif di perangkat Anda.",
		})
		return
	}
	
	log.Printf("[CHECK-IN] Trying GPS validation - Lat: %f, Lng: %f", *req.Latitude, *req.Longitude)
	isValidGPS, gpsConfig, err := h.attendanceService.ValidateGPS(*req.Latitude, *req.Longitude)
	if err != nil {
		log.Printf("[CHECK-IN] GPS validation error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan saat validasi GPS",
		})
		return
	}
	
	if isValidGPS && gpsConfig != nil {
		log.Printf("[CHECK-IN] GPS validation SUCCESS - Location: %s, Radius: %d meters", 
			gpsConfig.Name, gpsConfig.Radius)
		ssid = "GPS-" + gpsConfig.Name
		bssid = "GPS"
		validationMethod = "gps_validation"
		validationDetails = map[string]interface{}{
			"location_name": gpsConfig.Name,
			"latitude":      *req.Latitude,
			"longitude":     *req.Longitude,
			"radius":        gpsConfig.Radius,
		}
		isValidated = true
	} else {
		log.Printf("[CHECK-IN] GPS validation FAILED - Not within any authorized location")
	}
	
	// If GPS validation failed, return error
	if !isValidated {
		log.Printf("[CHECK-IN] GPS validation FAILED - Employee ID: %d", employee.ID)
		c.JSON(http.StatusForbidden, gin.H{
			"success":    false,
			"error_code": "INVALID_LOCATION",
			"message":    "Anda tidak berada di area kantor yang terdaftar. Pastikan GPS aktif dan Anda berada dalam radius lokasi kantor.",
			"details": map[string]interface{}{
				"latitude":  *req.Latitude,
				"longitude": *req.Longitude,
			},
		})
		return
	}

	attendance, err := h.attendanceService.CheckIn(employee.ID, ssid, bssid)
	if err != nil {
		if err == services.ErrInvalidWiFi {
			c.JSON(http.StatusForbidden, gin.H{
				"success":    false,
				"error_code": "INVALID_WIFI",
				"message":    "Anda harus terhubung ke Wi-Fi kantor untuk absen",
			})
			return
		}
		if err == services.ErrAlreadyCheckedIn {
			log.Printf("[CHECK-IN] Already checked in - Employee ID: %d", employee.ID)
			c.JSON(http.StatusConflict, gin.H{
				"success":    false,
				"error_code": "ALREADY_CHECKED_IN",
				"message":    "Anda sudah melakukan check-in hari ini",
			})
			return
		}
		log.Printf("[CHECK-IN] Error creating attendance: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	// Record in audit trail
	h.auditService.RecordAction(userID.(uint), "check_in", "attendance", strconv.Itoa(int(attendance.ID)), "", "", clientIP)

	log.Printf("[CHECK-IN] SUCCESS - Attendance ID: %d, Employee ID: %d, Method: %s", 
		attendance.ID, employee.ID, validationMethod)

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"message":       "Check-in berhasil",
		"data":          attendance,
		"validated_by":  map[string]interface{}{
			"method":  validationMethod,
			"details": validationDetails,
		},
	})
}

// CheckOut records employee check-out
func (h *HRMHandler) CheckOut(c *gin.Context) {
	// Check-out doesn't require WiFi/IP validation
	// Employee can check-out from anywhere
	
	// Get employee ID from user ID
	userID, _ := c.Get("user_id")
	employee, err := h.employeeService.GetEmployeeByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success":    false,
			"error_code": "EMPLOYEE_NOT_FOUND",
			"message":    "Data karyawan tidak ditemukan",
		})
		return
	}

	// No WiFi validation needed for check-out
	// Use empty strings for SSID and BSSID
	attendance, err := h.attendanceService.CheckOut(employee.ID, "", "")
	if err != nil {
		if err == services.ErrNotCheckedIn {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "NOT_CHECKED_IN",
				"message":    "Anda belum melakukan check-in hari ini",
			})
			return
		}
		if err == services.ErrAlreadyCheckedOut {
			c.JSON(http.StatusConflict, gin.H{
				"success":    false,
				"error_code": "ALREADY_CHECKED_OUT",
				"message":    "Anda sudah melakukan check-out hari ini",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	// Record in audit trail
	h.auditService.RecordAction(userID.(uint), "check_out", "attendance", strconv.Itoa(int(attendance.ID)), "", "", c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Check-out berhasil",
		"data":    attendance,
	})
}

// ValidateWiFi validates Wi-Fi connection
func (h *HRMHandler) ValidateWiFi(c *gin.Context) {
	var req CheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
		})
		return
	}

	isValid, err := h.attendanceService.ValidateWiFi(req.SSID, req.BSSID)
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
		"valid":   isValid,
	})
}

// GetTodayAttendance retrieves today's attendance for current user
func (h *HRMHandler) GetTodayAttendance(c *gin.Context) {
	userID, _ := c.Get("user_id")
	employee, err := h.employeeService.GetEmployeeByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success":    false,
			"error_code": "EMPLOYEE_NOT_FOUND",
			"message":    "Data karyawan tidak ditemukan",
		})
		return
	}

	attendance, err := h.attendanceService.GetTodayAttendance(employee.ID)
	if err != nil {
		if err == services.ErrAttendanceNotFound {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    nil,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    attendance,
	})
}

// GetAttendanceReport retrieves attendance report
func (h *HRMHandler) GetAttendanceReport(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Tanggal mulai dan tanggal akhir harus diisi",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_DATE",
			"message":    "Format tanggal tidak valid (gunakan YYYY-MM-DD)",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_DATE",
			"message":    "Format tanggal tidak valid (gunakan YYYY-MM-DD)",
		})
		return
	}

	report, err := h.attendanceService.GetAttendanceReport(startDate, endDate)
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
		"data":    report,
	})
}

// GetAttendanceStats retrieves attendance statistics
func (h *HRMHandler) GetAttendanceStats(c *gin.Context) {
	dateStr := c.Query("date")
	var date time.Time
	var err error

	if dateStr == "" {
		date = time.Now()
	} else {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INVALID_DATE",
				"message":    "Format tanggal tidak valid (gunakan YYYY-MM-DD)",
			})
			return
		}
	}

	stats, err := h.attendanceService.GetAttendanceStats(date)
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
		"data":    stats,
	})
}

// WiFi Configuration Endpoints

// CreateWiFiConfigRequest represents create Wi-Fi config request
type CreateWiFiConfigRequest struct {
	SSID     string `json:"ssid" binding:"required"`
	BSSID    string `json:"bssid" binding:"required"`
	Location string `json:"location"`
}

// CreateWiFiConfig creates a new Wi-Fi configuration
func (h *HRMHandler) CreateWiFiConfig(c *gin.Context) {
	var req CreateWiFiConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
		})
		return
	}

	config := &models.WiFiConfig{
		SSID:     req.SSID,
		BSSID:    req.BSSID,
		Location: req.Location,
		IsActive: true,
	}

	if err := h.attendanceService.CreateWiFiConfig(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    err.Error(),
		})
		return
	}

	// Record in audit trail
	userID, _ := c.Get("user_id")
	h.auditService.RecordAction(userID.(uint), "create", "wifi_config", strconv.Itoa(int(config.ID)), "", "", c.ClientIP())

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Konfigurasi Wi-Fi berhasil dibuat",
		"data":    config,
	})
}

// GetWiFiConfigs retrieves all Wi-Fi configurations
func (h *HRMHandler) GetWiFiConfigs(c *gin.Context) {
	activeOnlyStr := c.Query("active_only")
	activeOnly := activeOnlyStr == "true"

	configs, err := h.attendanceService.GetAllWiFiConfigs(activeOnly)
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
		"data":    configs,
	})
}

// UpdateWiFiConfig updates a Wi-Fi configuration
func (h *HRMHandler) UpdateWiFiConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
		})
		return
	}

	config, err := h.attendanceService.UpdateWiFiConfig(uint(id), updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    err.Error(),
		})
		return
	}

	// Record in audit trail
	userID, _ := c.Get("user_id")
	h.auditService.RecordAction(userID.(uint), "update", "wifi_config", strconv.Itoa(int(id)), "", "", c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Konfigurasi Wi-Fi berhasil diperbarui",
		"data":    config,
	})
}

// DeleteWiFiConfig deletes a Wi-Fi configuration
func (h *HRMHandler) DeleteWiFiConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	if err := h.attendanceService.DeleteWiFiConfig(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    err.Error(),
		})
		return
	}

	// Record in audit trail
	userID, _ := c.Get("user_id")
	h.auditService.RecordAction(userID.(uint), "delete", "wifi_config", strconv.Itoa(int(id)), "", "", c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Konfigurasi Wi-Fi berhasil dihapus",
	})
}
// ExportAttendanceReport exports attendance report to Excel or PDF
func (h *HRMHandler) ExportAttendanceReport(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	format := c.Query("format") // "excel" or "pdf"

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Tanggal mulai dan tanggal akhir harus diisi",
		})
		return
	}

	if format != "excel" && format != "pdf" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Format harus 'excel' atau 'pdf'",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_DATE",
			"message":    "Format tanggal tidak valid (gunakan YYYY-MM-DD)",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_DATE",
			"message":    "Format tanggal tidak valid (gunakan YYYY-MM-DD)",
		})
		return
	}

	// Get attendance report data
	report, err := h.attendanceService.GetAttendanceReport(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	// Get current user for generated by field
	userID, _ := c.Get("user_id")
	user, _ := h.employeeService.GetEmployeeByUserID(userID.(uint))
	generatedBy := "System"
	if user != nil {
		generatedBy = user.FullName
	}

	// Prepare export data
	exportService := services.NewExportService("Sistem ERP SPPG")
	
	headers := []string{
		"Nama Karyawan",
		"Posisi",
		"Total Hari",
		"Total Jam",
		"Rata-rata Jam/Hari",
	}

	rows := make([][]string, len(report))
	for i, item := range report {
		totalHours := "0.0"
		averageHours := "0.0"
		
		if val, ok := item["total_hours"].(float64); ok {
			totalHours = strconv.FormatFloat(val, 'f', 1, 64)
		}
		if val, ok := item["average_hours"].(float64); ok {
			averageHours = strconv.FormatFloat(val, 'f', 1, 64)
		}

		rows[i] = []string{
			item["full_name"].(string),
			item["position"].(string),
			strconv.Itoa(int(item["total_days"].(int64))),
			totalHours + " jam",
			averageHours + " jam",
		}
	}

	exportData := &services.ExportData{
		Title:       "Laporan Absensi Karyawan",
		Headers:     headers,
		Rows:        rows,
		DateRange:   startDate.Format("02/01/2006") + " - " + endDate.Format("02/01/2006"),
		GeneratedBy: generatedBy,
	}

	var buffer *bytes.Buffer
	var contentType string
	var filename string

	if format == "excel" {
		buffer, err = exportService.ExportToExcel(exportData)
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		filename = "laporan-absensi-" + startDate.Format("2006-01-02") + "-" + endDate.Format("2006-01-02") + ".xlsx"
	} else {
		buffer, err = exportService.ExportToPDF(exportData)
		contentType = "application/pdf"
		filename = "laporan-absensi-" + startDate.Format("2006-01-02") + "-" + endDate.Format("2006-01-02") + ".pdf"
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "EXPORT_ERROR",
			"message":    "Gagal mengekspor laporan: " + err.Error(),
		})
		return
	}

	// Record in audit trail
	h.auditService.RecordAction(userID.(uint), "export", "attendance_report", 
		format, "", "Export periode: "+exportData.DateRange, c.ClientIP())

	// Set headers and return file
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Length", strconv.Itoa(buffer.Len()))
	
	c.Data(http.StatusOK, contentType, buffer.Bytes())
}
// GetAttendanceByDateRange retrieves attendance records for a specific employee and date range
func (h *HRMHandler) GetAttendanceByDateRange(c *gin.Context) {
	employeeIDStr := c.Query("employee_id")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Tanggal mulai dan tanggal akhir harus diisi",
		})
		return
	}

	var empID *uint
	
	// If employee_id not provided, use current user's employee_id
	if employeeIDStr == "" {
		userID, _ := c.Get("user_id")
		employee, err := h.employeeService.GetEmployeeByUserID(userID.(uint))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "EMPLOYEE_NOT_FOUND",
				"message":    "Data karyawan tidak ditemukan",
			})
			return
		}
		empID = &employee.ID
	} else {
		employeeID, err := strconv.ParseUint(employeeIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INVALID_EMPLOYEE_ID",
				"message":    "Employee ID tidak valid",
			})
			return
		}
		id := uint(employeeID)
		empID = &id
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_DATE",
			"message":    "Format tanggal tidak valid (gunakan YYYY-MM-DD)",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_DATE",
			"message":    "Format tanggal tidak valid (gunakan YYYY-MM-DD)",
		})
		return
	}

	attendances, err := h.attendanceService.GetAttendanceByDateRange(empID, startDate, endDate)
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
		"data":    attendances,
	})
}


// GPS Configuration Endpoints

// CreateGPSConfigRequest represents create GPS config request
type CreateGPSConfigRequest struct {
	Name        string  `json:"name" binding:"required"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	Radius      int     `json:"radius"`
	Address     string  `json:"address"`
	Description string  `json:"description"`
}

// CreateGPSConfig creates a new GPS configuration
func (h *HRMHandler) CreateGPSConfig(c *gin.Context) {
	var req CreateGPSConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
		})
		return
	}

	config := &models.GPSConfig{
		Name:        req.Name,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Radius:      req.Radius,
		Address:     req.Address,
		Description: req.Description,
		IsActive:    true,
	}

	if config.Radius <= 0 {
		config.Radius = 100 // Default 100 meters
	}

	if err := h.attendanceService.CreateGPSConfig(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    err.Error(),
		})
		return
	}

	// Record in audit trail
	userID, _ := c.Get("user_id")
	h.auditService.RecordAction(userID.(uint), "create", "gps_config", strconv.Itoa(int(config.ID)), "", "", c.ClientIP())

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Konfigurasi GPS berhasil dibuat",
		"data":    config,
	})
}

// GetGPSConfigs retrieves all GPS configurations
func (h *HRMHandler) GetGPSConfigs(c *gin.Context) {
	activeOnlyStr := c.Query("active_only")
	activeOnly := activeOnlyStr == "true"

	configs, err := h.attendanceService.GetAllGPSConfigs(activeOnly)
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
		"data":    configs,
	})
}

// UpdateGPSConfig updates a GPS configuration
func (h *HRMHandler) UpdateGPSConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
		})
		return
	}

	config, err := h.attendanceService.UpdateGPSConfig(uint(id), updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    err.Error(),
		})
		return
	}

	// Record in audit trail
	userID, _ := c.Get("user_id")
	h.auditService.RecordAction(userID.(uint), "update", "gps_config", strconv.Itoa(int(id)), "", "", c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Konfigurasi GPS berhasil diperbarui",
		"data":    config,
	})
}

// DeleteGPSConfig deletes a GPS configuration
func (h *HRMHandler) DeleteGPSConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	if err := h.attendanceService.DeleteGPSConfig(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    err.Error(),
		})
		return
	}

	// Record in audit trail
	userID, _ := c.Get("user_id")
	h.auditService.RecordAction(userID.(uint), "delete", "gps_config", strconv.Itoa(int(id)), "", "", c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Konfigurasi GPS berhasil dihapus",
	})
}
