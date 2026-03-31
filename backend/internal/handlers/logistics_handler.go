package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/erp-sppg/backend/internal/middleware"
	"github.com/erp-sppg/backend/internal/models"
	"github.com/erp-sppg/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LogisticsHandler handles logistics and distribution endpoints
type LogisticsHandler struct {
	db                     *gorm.DB
	schoolService          *services.SchoolService
	deliveryTaskService    *services.DeliveryTaskService
	epodService            *services.EPODService
	omprengTrackingService *services.OmprengTrackingService
}

// NewLogisticsHandler creates a new logistics handler
func NewLogisticsHandler(db *gorm.DB) *LogisticsHandler {
	return &LogisticsHandler{
		db:                     db,
		schoolService:          services.NewSchoolService(db),
		deliveryTaskService:    services.NewDeliveryTaskService(db),
		epodService:            services.NewEPODService(db),
		omprengTrackingService: services.NewOmprengTrackingService(db),
	}
}

// School Endpoints

// CreateSchoolRequest represents create school request
type CreateSchoolRequest struct {
	Name                 string  `json:"name" binding:"required"`
	Address              string  `json:"address"`
	Latitude             float64 `json:"latitude" binding:"required"`
	Longitude            float64 `json:"longitude" binding:"required"`
	ContactPerson        string  `json:"contact_person"`
	PhoneNumber          string  `json:"phone_number"`
	StudentCount         int     `json:"student_count" binding:"gte=0"`
	Category             string  `json:"category" binding:"required,oneof=SD SMP SMA"`
	StudentCountGrade13  int     `json:"student_count_grade_1_3" binding:"gte=0"`
	StudentCountGrade46  int     `json:"student_count_grade_4_6" binding:"gte=0"`
	StaffCount           int     `json:"staff_count" binding:"gte=0"`
	NPSN                 string  `json:"npsn"`
	PrincipalName        string  `json:"principal_name"`
	SchoolEmail          string  `json:"school_email"`
	SchoolPhone          string  `json:"school_phone"`
	CommitteeCount       int     `json:"committee_count" binding:"gte=0"`
	CooperationLetterURL string  `json:"cooperation_letter_url"`
}

// CreateSchool creates a new school
func (h *LogisticsHandler) CreateSchool(c *gin.Context) {
	var req CreateSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	school := &models.School{
		Name:                 req.Name,
		Address:              req.Address,
		Latitude:             req.Latitude,
		Longitude:            req.Longitude,
		ContactPerson:        req.ContactPerson,
		PhoneNumber:          req.PhoneNumber,
		StudentCount:         req.StudentCount,
		Category:             req.Category,
		StudentCountGrade13:  req.StudentCountGrade13,
		StudentCountGrade46:  req.StudentCountGrade46,
		StaffCount:           req.StaffCount,
		NPSN:                 req.NPSN,
		PrincipalName:        req.PrincipalName,
		SchoolEmail:          req.SchoolEmail,
		SchoolPhone:          req.SchoolPhone,
		CommitteeCount:       req.CommitteeCount,
		CooperationLetterURL: req.CooperationLetterURL,
	}

	// Auto-inject sppg_id for SPPG-level roles
	if sppgID, ok := middleware.GetTenantSPPGID(c); ok {
		school.SPPGID = &sppgID
	}

	scopedService := h.schoolService.WithDB(getTenantScopedDB(c, h.db))
	if err := scopedService.CreateSchool(school); err != nil {
		if err == services.ErrDuplicateSchool {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "DUPLICATE_SCHOOL",
				"message":    "Sekolah dengan nama yang sama sudah ada",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "CREATE_SCHOOL_ERROR",
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Sekolah berhasil dibuat",
		"school":  school,
	})
}

// GetSchool retrieves a school by ID
func (h *LogisticsHandler) GetSchool(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	scopedService := h.schoolService.WithDB(getTenantScopedDB(c, h.db))
	school, err := scopedService.GetSchoolByID(uint(id))
	if err != nil {
		if err == services.ErrSchoolNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "SCHOOL_NOT_FOUND",
				"message":    "Sekolah tidak ditemukan",
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
		"school":  school,
	})
}

// GetAllSchools retrieves all schools
func (h *LogisticsHandler) GetAllSchools(c *gin.Context) {
	activeOnly := c.DefaultQuery("active_only", "true") == "true"
	query := c.Query("q")

	scopedService := h.schoolService.WithDB(getTenantScopedDB(c, h.db))
	var schools []models.School
	var err error

	if query != "" {
		schools, err = scopedService.SearchSchools(query, activeOnly)
	} else {
		schools, err = scopedService.GetAllSchools(activeOnly)
	}

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
		"schools": schools,
	})
}

// UpdateSchool updates an existing school
func (h *LogisticsHandler) UpdateSchool(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	var req CreateSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	school := &models.School{
		Name:                 req.Name,
		Address:              req.Address,
		Latitude:             req.Latitude,
		Longitude:            req.Longitude,
		ContactPerson:        req.ContactPerson,
		PhoneNumber:          req.PhoneNumber,
		StudentCount:         req.StudentCount,
		Category:             req.Category,
		StudentCountGrade13:  req.StudentCountGrade13,
		StudentCountGrade46:  req.StudentCountGrade46,
		StaffCount:           req.StaffCount,
		NPSN:                 req.NPSN,
		PrincipalName:        req.PrincipalName,
		SchoolEmail:          req.SchoolEmail,
		SchoolPhone:          req.SchoolPhone,
		CommitteeCount:       req.CommitteeCount,
		CooperationLetterURL: req.CooperationLetterURL,
	}

	scopedService := h.schoolService.WithDB(getTenantScopedDB(c, h.db))
	if err := scopedService.UpdateSchool(uint(id), school); err != nil {
		if err == services.ErrSchoolNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "SCHOOL_NOT_FOUND",
				"message":    "Sekolah tidak ditemukan",
			})
			return
		}

		if err == services.ErrDuplicateSchool {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "DUPLICATE_SCHOOL",
				"message":    "Sekolah dengan nama yang sama sudah ada",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "UPDATE_SCHOOL_ERROR",
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Sekolah berhasil diperbarui",
	})
}

// DeleteSchool deletes a school
func (h *LogisticsHandler) DeleteSchool(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	scopedService := h.schoolService.WithDB(getTenantScopedDB(c, h.db))
	if err := scopedService.DeleteSchool(uint(id)); err != nil {
		if err == services.ErrSchoolNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "SCHOOL_NOT_FOUND",
				"message":    "Sekolah tidak ditemukan",
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
		"message": "Sekolah berhasil dihapus atau dinonaktifkan",
	})
}

// Delivery Task Endpoints

// CreateDeliveryTaskRequest represents create delivery task request
type CreateDeliveryTaskRequest struct {
	TaskDate         string                    `json:"task_date" binding:"required"`
	DeliveryTime     string                    `json:"delivery_time"`
	DeliveryRecordID uint                      `json:"delivery_record_id"` // Single record (legacy)
	DeliveryRecords  []DeliveryRecordRequest   `json:"delivery_records"`   // Multiple records (new)
	DriverID         uint                      `json:"driver_id" binding:"required"`
	// Legacy fields for backward compatibility
	SchoolID   uint                      `json:"school_id"`
	Portions   int                       `json:"portions"`
	RouteOrder int                       `json:"route_order"`
	MenuItems  []DeliveryMenuItemRequest `json:"menu_items"`
}

// DeliveryRecordRequest represents a delivery record with route order
type DeliveryRecordRequest struct {
	DeliveryRecordID uint `json:"delivery_record_id" binding:"required"`
	RouteOrder       int  `json:"route_order" binding:"required"`
}

// DeliveryMenuItemRequest represents delivery menu item request
type DeliveryMenuItemRequest struct {
	RecipeID uint `json:"recipe_id" binding:"required"`
	Portions int  `json:"portions" binding:"required,gt=0"`
}

// CreateDeliveryTask creates a new delivery task
func (h *LogisticsHandler) CreateDeliveryTask(c *gin.Context) {
	var req CreateDeliveryTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	// Parse task date
	taskDate, err := time.Parse("2006-01-02", req.TaskDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_DATE",
			"message":    "Format tanggal tidak valid (gunakan YYYY-MM-DD)",
		})
		return
	}

	// New format: using delivery_records array (multiple schools)
	if len(req.DeliveryRecords) > 0 {
		// Convert handler request to service format
		var serviceRecords []services.DeliveryRecordWithRoute
		for _, r := range req.DeliveryRecords {
			serviceRecords = append(serviceRecords, services.DeliveryRecordWithRoute{
				DeliveryRecordID: r.DeliveryRecordID,
				RouteOrder:       r.RouteOrder,
			})
		}
		
		tasks, err := h.deliveryTaskService.CreateDeliveryTasksFromRecords(taskDate, req.DriverID, serviceRecords)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "CREATE_TASK_ERROR",
				"message":    err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success":        true,
			"message":        "Tugas pengiriman berhasil dibuat",
			"delivery_tasks": tasks,
		})
		return
	}

	// Single record format: using delivery_record_id
	if req.DeliveryRecordID != 0 {
		task, err := h.deliveryTaskService.CreateDeliveryTaskFromRecord(taskDate, req.DriverID, req.DeliveryRecordID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "CREATE_TASK_ERROR",
				"message":    err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success":       true,
			"message":       "Tugas pengiriman berhasil dibuat",
			"delivery_task": task,
		})
		return
	}

	// Legacy format: using school_id, portions, menu_items
	task := &models.DeliveryTask{
		TaskDate:   taskDate,
		DriverID:   req.DriverID,
		SchoolID:   req.SchoolID,
		Portions:   req.Portions,
		RouteOrder: req.RouteOrder,
	}

	var menuItems []models.DeliveryMenuItem
	for _, item := range req.MenuItems {
		menuItems = append(menuItems, models.DeliveryMenuItem{
			RecipeID: item.RecipeID,
			Portions: item.Portions,
		})
	}

	if err := h.deliveryTaskService.CreateDeliveryTask(task, menuItems); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "CREATE_TASK_ERROR",
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":       true,
		"message":       "Tugas pengiriman berhasil dibuat",
		"delivery_task": task,
	})
}

// GetDeliveryTask retrieves a delivery task by ID
func (h *LogisticsHandler) GetDeliveryTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	scopedService := h.deliveryTaskService.WithDB(getTenantScopedDB(c, h.db))
	task, err := scopedService.GetDeliveryTaskByID(uint(id))
	if err != nil {
		if err == services.ErrDeliveryTaskNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "TASK_NOT_FOUND",
				"message":    "Tugas pengiriman tidak ditemukan",
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
		"success":       true,
		"delivery_task": task,
	})
}

// GetAllDeliveryTasks retrieves all delivery tasks with filters
func (h *LogisticsHandler) GetAllDeliveryTasks(c *gin.Context) {
	var driverID *uint
	if idStr := c.Query("driver_id"); idStr != "" {
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err == nil {
			uid := uint(id)
			driverID = &uid
		}
	}

	status := c.Query("status")

	var date *time.Time
	if dateStr := c.Query("date"); dateStr != "" {
		if d, err := time.Parse("2006-01-02", dateStr); err == nil {
			date = &d
		}
	}

	scopedService := h.deliveryTaskService.WithDB(getTenantScopedDB(c, h.db))
	tasks, err := scopedService.GetAllDeliveryTasks(driverID, status, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"delivery_tasks": tasks,
	})
}

// GetDriverTasksToday retrieves delivery tasks for a driver for today
func (h *LogisticsHandler) GetDriverTasksToday(c *gin.Context) {
	driverID, err := strconv.ParseUint(c.Param("driver_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID driver tidak valid",
		})
		return
	}

	scopedService := h.deliveryTaskService.WithDB(getTenantScopedDB(c, h.db))
	tasks, err := scopedService.GetDriverTasksForToday(uint(driverID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"delivery_tasks": tasks,
	})
}

// UpdateDeliveryTaskStatusRequest represents update status request
type UpdateDeliveryTaskStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending in_progress arrived received cancelled"`
}

// UpdateDeliveryTaskStatus updates the status of a delivery task
func (h *LogisticsHandler) UpdateDeliveryTaskStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	var req UpdateDeliveryTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	scopedService := h.deliveryTaskService.WithDB(getTenantScopedDB(c, h.db))
	if err := scopedService.UpdateDeliveryTaskStatus(uint(id), req.Status); err != nil {
		if err == services.ErrDeliveryTaskNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "TASK_NOT_FOUND",
				"message":    "Tugas pengiriman tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "UPDATE_STATUS_ERROR",
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Status tugas pengiriman berhasil diperbarui",
	})
}

// UpdateDeliveryTaskRequest represents update delivery task request
type UpdateDeliveryTaskRequest struct {
	TaskDate   string                    `json:"task_date"`
	DriverID   uint                      `json:"driver_id"`
	SchoolID   uint                      `json:"school_id"`
	Portions   int                       `json:"portions"`
	RouteOrder int                       `json:"route_order"`
	MenuItems  []DeliveryMenuItemRequest `json:"menu_items"`
}

// UpdateDeliveryTask updates an existing delivery task
func (h *LogisticsHandler) UpdateDeliveryTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	var req UpdateDeliveryTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	task := &models.DeliveryTask{}

	// Parse task date if provided
	if req.TaskDate != "" {
		taskDate, err := time.Parse("2006-01-02", req.TaskDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INVALID_DATE",
				"message":    "Format tanggal tidak valid (gunakan YYYY-MM-DD)",
			})
			return
		}
		task.TaskDate = taskDate
	}

	task.DriverID = req.DriverID
	task.SchoolID = req.SchoolID
	task.Portions = req.Portions
	task.RouteOrder = req.RouteOrder

	var menuItems []models.DeliveryMenuItem
	for _, item := range req.MenuItems {
		menuItems = append(menuItems, models.DeliveryMenuItem{
			RecipeID: item.RecipeID,
			Portions: item.Portions,
		})
	}

	if err := h.deliveryTaskService.UpdateDeliveryTask(uint(id), task, menuItems); err != nil {
		if err == services.ErrDeliveryTaskNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "TASK_NOT_FOUND",
				"message":    "Tugas pengiriman tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "UPDATE_TASK_ERROR",
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tugas pengiriman berhasil diperbarui",
	})
}

// DeleteDeliveryTask deletes a delivery task
func (h *LogisticsHandler) DeleteDeliveryTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	scopedService := h.deliveryTaskService.WithDB(getTenantScopedDB(c, h.db))
	if err := scopedService.DeleteDeliveryTask(uint(id)); err != nil {
		if err == services.ErrDeliveryTaskNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "TASK_NOT_FOUND",
				"message":    "Tugas pengiriman tidak ditemukan",
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
		"message": "Tugas pengiriman berhasil dihapus",
	})
}

// GetReadyOrders retrieves delivery records that are ready for delivery (status = siap_dikirim)
func (h *LogisticsHandler) GetReadyOrders(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "MISSING_DATE",
			"message":    "Parameter tanggal wajib diisi",
		})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_DATE",
			"message":    "Format tanggal tidak valid (gunakan YYYY-MM-DD)",
		})
		return
	}

	scopedService := h.deliveryTaskService.WithDB(getTenantScopedDB(c, h.db))
	orders, err := scopedService.GetReadyOrders(date)
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
		"orders":  orders,
	})
}

// GetAvailableDrivers retrieves drivers that are not assigned on the given date
func (h *LogisticsHandler) GetAvailableDrivers(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "MISSING_DATE",
			"message":    "Parameter tanggal wajib diisi",
		})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_DATE",
			"message":    "Format tanggal tidak valid (gunakan YYYY-MM-DD)",
		})
		return
	}

	scopedService := h.deliveryTaskService.WithDB(getTenantScopedDB(c, h.db))
	drivers, err := scopedService.GetAvailableDrivers(date)
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
		"drivers": drivers,
	})
}

// e-POD Endpoints

// CreateEPODRequest represents create e-POD request
type CreateEPODRequest struct {
	DeliveryTaskID uint    `json:"delivery_task_id" binding:"required"`
	Latitude       float64 `json:"latitude" binding:"required"`
	Longitude      float64 `json:"longitude" binding:"required"`
	RecipientName  string  `json:"recipient_name"`
	OmprengDropOff int     `json:"ompreng_drop_off" binding:"gte=0"`
	OmprengPickUp  int     `json:"ompreng_pick_up" binding:"gte=0"`
}
// GetEPODByDeliveryTask retrieves an e-POD by delivery task ID
func (h *LogisticsHandler) GetEPODByDeliveryTask(c *gin.Context) {
	// Support both delivery_task_id and delivery_record_id
	deliveryTaskIDStr := c.Query("delivery_task_id")
	deliveryRecordIDStr := c.Query("delivery_record_id")

	if deliveryTaskIDStr != "" {
		// Get by delivery_task_id
		deliveryTaskID, err := strconv.ParseUint(deliveryTaskIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INVALID_ID",
				"message":    "ID tugas pengiriman tidak valid",
			})
			return
		}

		epod, err := h.epodService.GetEPODByDeliveryTaskID(uint(deliveryTaskID))
		if err != nil {
			if err == services.ErrEPODNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"success":    false,
					"error_code": "EPOD_NOT_FOUND",
					"message":    "e-POD tidak ditemukan untuk tugas pengiriman ini",
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
			"epod":    epod,
		})
		return
	}

	if deliveryRecordIDStr != "" {
		// Get by delivery_record_id - find matching ePOD by school, driver, and date
		deliveryRecordID, err := strconv.ParseUint(deliveryRecordIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INVALID_ID",
				"message":    "ID delivery record tidak valid",
			})
			return
		}

		epod, err := h.epodService.GetEPODByDeliveryRecordID(uint(deliveryRecordID))
		if err != nil {
			if err == services.ErrEPODNotFound {
				c.JSON(http.StatusNotFound, gin.H{
					"success":    false,
					"error_code": "EPOD_NOT_FOUND",
					"message":    "e-POD tidak ditemukan untuk delivery record ini",
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
			"epod":    epod,
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"success":    false,
		"error_code": "MISSING_PARAMETER",
		"message":    "Parameter delivery_task_id atau delivery_record_id diperlukan",
	})
}



// CreateEPOD creates a new electronic proof of delivery
func (h *LogisticsHandler) CreateEPOD(c *gin.Context) {
	var req CreateEPODRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	epod := &models.ElectronicPOD{
		DeliveryTaskID: req.DeliveryTaskID,
		Latitude:       req.Latitude,
		Longitude:      req.Longitude,
		RecipientName:  req.RecipientName,
		OmprengDropOff: req.OmprengDropOff,
		OmprengPickUp:  req.OmprengPickUp,
	}

	if err := h.epodService.CreateEPOD(epod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "CREATE_EPOD_ERROR",
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "e-POD berhasil dibuat",
		"epod":    epod,
	})
}

// UploadEPODPhotoRequest represents upload photo request
type UploadEPODPhotoRequest struct {
	PhotoURL string `json:"photo_url"`
}

// UploadEPODPhoto uploads photo for an e-POD
func (h *LogisticsHandler) UploadEPODPhoto(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	// Check content type - support both JSON and multipart form
	contentType := c.GetHeader("Content-Type")
	var photoURL string

	if strings.Contains(contentType, "multipart/form-data") {
		// Handle file upload
		file, err := c.FormFile("photo")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "VALIDATION_ERROR",
				"message":    "File foto tidak ditemukan",
				"details":    err.Error(),
			})
			return
		}

		// Validate file size (max 5MB)
		if file.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "FILE_TOO_LARGE",
				"message":    "Ukuran file terlalu besar. Maksimal 5MB.",
			})
			return
		}

		// Generate unique filename
		ext := filepath.Ext(file.Filename)
		if ext == "" {
			ext = ".jpg"
		}
		filename := fmt.Sprintf("epod-photo-%d-%d%s", id, time.Now().UnixNano(), ext)
		uploadPath := filepath.Join("uploads", "epod", filename)

		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(uploadPath), 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success":    false,
				"error_code": "INTERNAL_ERROR",
				"message":    "Gagal membuat direktori upload",
			})
			return
		}

		// Save file
		if err := c.SaveUploadedFile(file, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success":    false,
				"error_code": "INTERNAL_ERROR",
				"message":    "Gagal menyimpan file",
			})
			return
		}

		// Generate URL
		photoURL = fmt.Sprintf("/uploads/epod/%s", filename)
	} else {
		// Handle JSON request
		var req UploadEPODPhotoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "VALIDATION_ERROR",
				"message":    "Data tidak valid",
				"details":    err.Error(),
			})
			return
		}
		photoURL = req.PhotoURL
	}

	if photoURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "URL foto tidak boleh kosong",
		})
		return
	}

	if err := h.epodService.UpdateEPODPhoto(uint(id), photoURL); err != nil {
		if err == services.ErrEPODNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "EPOD_NOT_FOUND",
				"message":    "e-POD tidak ditemukan",
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
		"success":   true,
		"message":   "Foto e-POD berhasil diunggah",
		"photo_url": photoURL,
	})
}

// UploadEPODSignatureRequest represents upload signature request
type UploadEPODSignatureRequest struct {
	SignatureURL string `json:"signature_url"`
}

// UploadEPODSignature uploads signature for an e-POD
func (h *LogisticsHandler) UploadEPODSignature(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	// Check content type - support both JSON and multipart form
	contentType := c.GetHeader("Content-Type")
	var signatureURL string

	if strings.Contains(contentType, "multipart/form-data") {
		// Handle file upload
		file, err := c.FormFile("signature")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "VALIDATION_ERROR",
				"message":    "File tanda tangan tidak ditemukan",
				"details":    err.Error(),
			})
			return
		}

		// Validate file size (max 2MB for signatures)
		if file.Size > 2*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "FILE_TOO_LARGE",
				"message":    "Ukuran file terlalu besar. Maksimal 2MB.",
			})
			return
		}

		// Generate unique filename
		ext := filepath.Ext(file.Filename)
		if ext == "" {
			ext = ".png"
		}
		filename := fmt.Sprintf("epod-signature-%d-%d%s", id, time.Now().UnixNano(), ext)
		uploadPath := filepath.Join("uploads", "epod", filename)

		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(uploadPath), 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success":    false,
				"error_code": "INTERNAL_ERROR",
				"message":    "Gagal membuat direktori upload",
			})
			return
		}

		// Save file
		if err := c.SaveUploadedFile(file, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success":    false,
				"error_code": "INTERNAL_ERROR",
				"message":    "Gagal menyimpan file",
			})
			return
		}

		// Generate URL
		signatureURL = fmt.Sprintf("/uploads/epod/%s", filename)
	} else {
		// Handle JSON request
		var req UploadEPODSignatureRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "VALIDATION_ERROR",
				"message":    "Data tidak valid",
				"details":    err.Error(),
			})
			return
		}
		signatureURL = req.SignatureURL
	}

	if signatureURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "URL tanda tangan tidak boleh kosong",
		})
		return
	}

	if err := h.epodService.UpdateEPODSignature(uint(id), signatureURL); err != nil {
		if err == services.ErrEPODNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "EPOD_NOT_FOUND",
				"message":    "e-POD tidak ditemukan",
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
		"success":       true,
		"message":       "Tanda tangan e-POD berhasil diunggah",
		"signature_url": signatureURL,
	})
}

// Ompreng Tracking Endpoints

// GetOmprengTracking retrieves ompreng tracking data
func (h *LogisticsHandler) GetOmprengTracking(c *gin.Context) {
	balances, err := h.omprengTrackingService.GetAllSchoolBalances()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"balances": balances,
	})
}

// RecordOmprengDropOffRequest represents drop-off request
type RecordOmprengDropOffRequest struct {
	SchoolID uint `json:"school_id" binding:"required"`
	Quantity int  `json:"quantity" binding:"required,gt=0"`
}

// RecordOmprengDropOff records ompreng drop-off at a school
func (h *LogisticsHandler) RecordOmprengDropOff(c *gin.Context) {
	var req RecordOmprengDropOffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, _ := c.Get("user_id")

	if err := h.omprengTrackingService.RecordOmprengMovement(req.SchoolID, req.Quantity, 0, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "RECORD_ERROR",
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Drop-off ompreng berhasil dicatat",
	})
}

// RecordOmprengPickUpRequest represents pick-up request
type RecordOmprengPickUpRequest struct {
	SchoolID uint `json:"school_id" binding:"required"`
	Quantity int  `json:"quantity" binding:"required,gt=0"`
}

// RecordOmprengPickUp records ompreng pick-up from a school
func (h *LogisticsHandler) RecordOmprengPickUp(c *gin.Context) {
	var req RecordOmprengPickUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, _ := c.Get("user_id")

	if err := h.omprengTrackingService.RecordOmprengMovement(req.SchoolID, 0, req.Quantity, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "RECORD_ERROR",
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pick-up ompreng berhasil dicatat",
	})
}

// GetOmprengReports generates ompreng circulation reports
func (h *LogisticsHandler) GetOmprengReports(c *gin.Context) {
	// Parse date range
	var startDate, endDate time.Time
	var err error

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		// Default to last 30 days
		endDate = time.Now()
		startDate = endDate.AddDate(0, 0, -30)
	} else {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INVALID_DATE",
				"message":    "Format tanggal mulai tidak valid (gunakan YYYY-MM-DD)",
			})
			return
		}

		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INVALID_DATE",
				"message":    "Format tanggal akhir tidak valid (gunakan YYYY-MM-DD)",
			})
			return
		}
	}

	report, err := h.omprengTrackingService.GenerateCirculationReport(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	// Get global inventory
	inventory, err := h.omprengTrackingService.GetGlobalInventory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	// Get missing ompreng
	missing, err := h.omprengTrackingService.GetMissingOmpreng()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"report":    report,
		"inventory": inventory,
		"missing":   missing,
	})
}

// UploadCooperationLetter handles cooperation letter file upload
func (h *LogisticsHandler) UploadCooperationLetter(c *gin.Context) {
	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "NO_FILE",
			"message":    "File tidak ditemukan",
		})
		return
	}

	// Validate file size (max 5MB for documents)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "FILE_TOO_LARGE",
			"message":    "Ukuran file maksimal 5MB",
		})
		return
	}

	// Validate file type (PDF, DOC, DOCX, JPG, PNG)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_FILE_TYPE",
			"message":    "Format file harus PDF, DOC, DOCX, JPG, JPEG, atau PNG",
		})
		return
	}

	// Create uploads directory if not exists
	uploadDir := "./uploads/cooperation-letters"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Gagal membuat direktori upload",
		})
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// Save file
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Gagal menyimpan file",
		})
		return
	}

	// Return URL
	fileURL := fmt.Sprintf("/uploads/cooperation-letters/%s", filename)
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "File berhasil diupload",
		"file_url": fileURL,
	})
}

// DeleteCooperationLetter deletes a cooperation letter file
func (h *LogisticsHandler) DeleteCooperationLetter(c *gin.Context) {
	fileURL := c.Query("file_url")
	if fileURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "MISSING_FILE_URL",
			"message":    "URL file tidak ditemukan",
		})
		return
	}

	// Extract filename from URL
	parts := strings.Split(fileURL, "/")
	if len(parts) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_FILE_URL",
			"message":    "URL file tidak valid",
		})
		return
	}

	filename := parts[len(parts)-1]
	filePath := filepath.Join("./uploads/cooperation-letters", filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"success":    false,
			"error_code": "FILE_NOT_FOUND",
			"message":    "File tidak ditemukan",
		})
		return
	}

	// Delete file
	if err := os.Remove(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Gagal menghapus file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File berhasil dihapus",
	})
}

// SchoolMapStat represents per-school delivery and review statistics for the map
type SchoolMapStat struct {
	SchoolID       uint    `json:"school_id"`
	PortionsSmall  int     `json:"portions_small"`
	PortionsLarge  int     `json:"portions_large"`
	TotalDelivered int     `json:"total_delivered"`
	TotalReviews   int     `json:"total_reviews"`
	AvgRating      float64 `json:"avg_rating"`
}

// GetSchoolMapStats returns per-school delivery and review stats (all-time, completed deliveries only)
// GET /api/v1/schools/map-stats
func (h *LogisticsHandler) GetSchoolMapStats(c *gin.Context) {
	scopedDB := getTenantScopedDB(c, h.db)

	// Aggregate delivery_records per school (only completed: current_status in received/completed stages)
	type deliveryStat struct {
		SchoolID      uint `gorm:"column:school_id"`
		PortionsSmall int  `gorm:"column:portions_small"`
		PortionsLarge int  `gorm:"column:portions_large"`
		TotalPortions int  `gorm:"column:total_portions"`
	}
	var deliveryStats []deliveryStat
	scopedDB.Model(&models.DeliveryRecord{}).
		Select("school_id, COALESCE(SUM(portions_small), 0) as portions_small, COALESCE(SUM(portions_large), 0) as portions_large, COALESCE(SUM(portions), 0) as total_portions").
		Where("current_stage >= 8"). // stage 8+ means delivered/received
		Group("school_id").
		Find(&deliveryStats)

	// Aggregate reviews per school
	type reviewStat struct {
		SchoolID     uint    `gorm:"column:school_id"`
		TotalReviews int     `gorm:"column:total_reviews"`
		AvgRating    float64 `gorm:"column:avg_rating"`
	}
	var reviewStats []reviewStat
	scopedDB.Model(&models.DeliveryReview{}).
		Select("school_id, COUNT(*) as total_reviews, COALESCE(AVG(overall_rating), 0) as avg_rating").
		Group("school_id").
		Find(&reviewStats)

	// Merge into map
	statsMap := make(map[uint]*SchoolMapStat)
	for _, d := range deliveryStats {
		statsMap[d.SchoolID] = &SchoolMapStat{
			SchoolID:       d.SchoolID,
			PortionsSmall:  d.PortionsSmall,
			PortionsLarge:  d.PortionsLarge,
			TotalDelivered: d.TotalPortions,
		}
	}
	for _, r := range reviewStats {
		if s, ok := statsMap[r.SchoolID]; ok {
			s.TotalReviews = r.TotalReviews
			s.AvgRating = r.AvgRating
		} else {
			statsMap[r.SchoolID] = &SchoolMapStat{
				SchoolID:     r.SchoolID,
				TotalReviews: r.TotalReviews,
				AvgRating:    r.AvgRating,
			}
		}
	}

	// Convert to slice
	result := make([]SchoolMapStat, 0, len(statsMap))
	for _, s := range statsMap {
		result = append(result, *s)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}
