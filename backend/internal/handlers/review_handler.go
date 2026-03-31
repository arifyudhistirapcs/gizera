package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/erp-sppg/backend/internal/models"
	"github.com/erp-sppg/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ReviewHandler handles delivery review endpoints
type ReviewHandler struct {
	db            *gorm.DB
	reviewService *services.ReviewService
}

// NewReviewHandler creates a new review handler
func NewReviewHandler(db *gorm.DB) *ReviewHandler {
	return &ReviewHandler{
		db:            db,
		reviewService: services.NewReviewService(db),
	}
}

// CreateReviewRequest represents create review request
type CreateReviewRequest struct {
	DeliveryRecordID uint   `json:"delivery_record_id" binding:"required"`
	SchoolID         uint   `json:"school_id" binding:"required"`
	ReviewerName     string `json:"reviewer_name"`
	ReviewerRole     string `json:"reviewer_role"`
	
	// Menu Ratings
	RatingFoodTaste       int `json:"rating_food_taste" binding:"required,min=1,max=5"`
	RatingFoodCleanliness int `json:"rating_food_cleanliness" binding:"required,min=1,max=5"`
	RatingMenuAccuracy    int `json:"rating_menu_accuracy" binding:"required,min=1,max=5"`
	RatingPortionSize     int `json:"rating_portion_size" binding:"required,min=1,max=5"`
	RatingMenuVariety     int `json:"rating_menu_variety" binding:"required,min=1,max=5"`
	
	// Service Ratings
	RatingDeliveryTime       int `json:"rating_delivery_time" binding:"required,min=1,max=5"`
	RatingDriverAttitude     int `json:"rating_driver_attitude" binding:"required,min=1,max=5"`
	RatingFoodCondition      int `json:"rating_food_condition" binding:"required,min=1,max=5"`
	RatingDriverTidiness     int `json:"rating_driver_tidiness" binding:"required,min=1,max=5"`
	RatingServiceConsistency int `json:"rating_service_consistency" binding:"required,min=1,max=5"`
	
	Comments string `json:"comments"`
	PhotoURL string `json:"photo_url"`
}

// CreateReview creates a new delivery review
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	// Handle base64 photo - save to file
	photoURL := ""
	if req.PhotoURL != "" && strings.HasPrefix(req.PhotoURL, "data:image") {
		// Extract base64 data
		parts := strings.Split(req.PhotoURL, ",")
		if len(parts) == 2 {
			// Decode base64
			imageData, err := base64.StdEncoding.DecodeString(parts[1])
			if err == nil {
				// Determine file extension
				ext := ".png"
				if strings.Contains(parts[0], "jpeg") || strings.Contains(parts[0], "jpg") {
					ext = ".jpg"
				}

				// Create uploads directory
				uploadDir := "uploads/reviews"
				if err := os.MkdirAll(uploadDir, 0755); err == nil {
					// Generate filename
					filename := fmt.Sprintf("review_%d_%d%s", req.DeliveryRecordID, time.Now().Unix(), ext)
					filePath := filepath.Join(uploadDir, filename)

					// Save file
					if err := os.WriteFile(filePath, imageData, 0644); err == nil {
						photoURL = "/" + filePath
					}
				}
			}
		}
	}

	review := &models.DeliveryReview{
		DeliveryRecordID:         req.DeliveryRecordID,
		SchoolID:                 req.SchoolID,
		ReviewerName:             req.ReviewerName,
		ReviewerRole:             req.ReviewerRole,
		RatingFoodTaste:          req.RatingFoodTaste,
		RatingFoodCleanliness:    req.RatingFoodCleanliness,
		RatingMenuAccuracy:       req.RatingMenuAccuracy,
		RatingPortionSize:        req.RatingPortionSize,
		RatingMenuVariety:        req.RatingMenuVariety,
		RatingDeliveryTime:       req.RatingDeliveryTime,
		RatingDriverAttitude:     req.RatingDriverAttitude,
		RatingFoodCondition:      req.RatingFoodCondition,
		RatingDriverTidiness:     req.RatingDriverTidiness,
		RatingServiceConsistency: req.RatingServiceConsistency,
		Comments:                 req.Comments,
		PhotoURL:                 photoURL,
	}

	scopedService := h.reviewService.WithDB(getTenantScopedDB(c, h.db))
	if err := scopedService.CreateReview(review); err != nil {
		if err == services.ErrReviewAlreadyExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "REVIEW_EXISTS",
				"message":    "Ulasan sudah ada untuk pengiriman ini",
			})
			return
		}
		if err == services.ErrInvalidRating {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INVALID_RATING",
				"message":    "Rating harus antara 1-5",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Gagal menyimpan ulasan",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Ulasan berhasil disimpan",
		"review":  review,
	})
}

// GetReview retrieves a review by ID
func (h *ReviewHandler) GetReview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	scopedService := h.reviewService.WithDB(getTenantScopedDB(c, h.db))
	review, err := scopedService.GetReviewByID(uint(id))
	if err != nil {
		if err == services.ErrReviewNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "REVIEW_NOT_FOUND",
				"message":    "Ulasan tidak ditemukan",
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
		"review":  review,
	})
}

// GetReviewByDeliveryRecord retrieves a review by delivery record ID
func (h *ReviewHandler) GetReviewByDeliveryRecord(c *gin.Context) {
	deliveryRecordID, err := strconv.ParseUint(c.Query("delivery_record_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID delivery record tidak valid",
		})
		return
	}

	scopedService := h.reviewService.WithDB(getTenantScopedDB(c, h.db))
	review, err := scopedService.GetReviewByDeliveryRecordID(uint(deliveryRecordID))
	if err != nil {
		if err == services.ErrReviewNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "REVIEW_NOT_FOUND",
				"message":    "Ulasan tidak ditemukan",
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
		"review":  review,
	})
}

// GetAllReviews retrieves all reviews with filters
func (h *ReviewHandler) GetAllReviews(c *gin.Context) {
	var schoolID *uint
	if idStr := c.Query("school_id"); idStr != "" {
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err == nil {
			uid := uint(id)
			schoolID = &uid
		}
	}

	var startDate, endDate *time.Time
	if dateStr := c.Query("start_date"); dateStr != "" {
		if d, err := time.Parse("2006-01-02", dateStr); err == nil {
			startDate = &d
		}
	}
	if dateStr := c.Query("end_date"); dateStr != "" {
		if d, err := time.Parse("2006-01-02", dateStr); err == nil {
			endDate = &d
		}
	}

	limit := 20
	offset := 0
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	scopedService := h.reviewService.WithDB(getTenantScopedDB(c, h.db))
	reviews, total, err := scopedService.GetAllReviews(schoolID, startDate, endDate, limit, offset)
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
		"reviews": reviews,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}

// GetReviewSummary retrieves aggregated review statistics
func (h *ReviewHandler) GetReviewSummary(c *gin.Context) {
	var schoolID *uint
	if idStr := c.Query("school_id"); idStr != "" {
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err == nil {
			uid := uint(id)
			schoolID = &uid
		}
	}

	var startDate, endDate *time.Time
	if dateStr := c.Query("start_date"); dateStr != "" {
		if d, err := time.Parse("2006-01-02", dateStr); err == nil {
			startDate = &d
		}
	}
	if dateStr := c.Query("end_date"); dateStr != "" {
		if d, err := time.Parse("2006-01-02", dateStr); err == nil {
			endDate = &d
		}
	}

	scopedService := h.reviewService.WithDB(getTenantScopedDB(c, h.db))
	summary, err := scopedService.GetReviewSummary(schoolID, startDate, endDate)
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
		"summary": summary,
	})
}

// CheckReviewExists checks if a review exists for a delivery record
func (h *ReviewHandler) CheckReviewExists(c *gin.Context) {
	deliveryRecordID, err := strconv.ParseUint(c.Query("delivery_record_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID delivery record tidak valid",
		})
		return
	}

	scopedService := h.reviewService.WithDB(getTenantScopedDB(c, h.db))
	exists, err := scopedService.HasReviewForDeliveryRecord(uint(deliveryRecordID))
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
		"exists":  exists,
	})
}
