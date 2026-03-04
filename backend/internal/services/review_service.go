package services

import (
	"errors"
	"time"

	"github.com/erp-sppg/backend/internal/models"
	"gorm.io/gorm"
)

var (
	ErrReviewNotFound      = errors.New("ulasan tidak ditemukan")
	ErrReviewAlreadyExists = errors.New("ulasan sudah ada untuk pengiriman ini")
	ErrInvalidRating       = errors.New("rating harus antara 1-5")
)

// ReviewService handles delivery review business logic
type ReviewService struct {
	db *gorm.DB
}

// NewReviewService creates a new review service
func NewReviewService(db *gorm.DB) *ReviewService {
	return &ReviewService{db: db}
}

// CreateReview creates a new delivery review
func (s *ReviewService) CreateReview(review *models.DeliveryReview) error {
	// Validate ratings (1-5)
	ratings := []int{
		review.RatingFoodTaste,
		review.RatingFoodCleanliness,
		review.RatingMenuAccuracy,
		review.RatingPortionSize,
		review.RatingMenuVariety,
		review.RatingDeliveryTime,
		review.RatingDriverAttitude,
		review.RatingFoodCondition,
		review.RatingDriverTidiness,
		review.RatingServiceConsistency,
	}
	
	for _, r := range ratings {
		if r < 1 || r > 5 {
			return ErrInvalidRating
		}
	}

	// Check if review already exists
	var existing models.DeliveryReview
	err := s.db.Where("delivery_record_id = ?", review.DeliveryRecordID).First(&existing).Error
	if err == nil {
		return ErrReviewAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Calculate averages
	review.CalculateAverages()
	review.CreatedAt = time.Now()
	review.UpdatedAt = time.Now()

	return s.db.Create(review).Error
}

// GetReviewByDeliveryRecordID retrieves a review by delivery record ID
func (s *ReviewService) GetReviewByDeliveryRecordID(deliveryRecordID uint) (*models.DeliveryReview, error) {
	var review models.DeliveryReview
	err := s.db.Preload("School").
		Preload("DeliveryRecord").
		Where("delivery_record_id = ?", deliveryRecordID).
		First(&review).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReviewNotFound
		}
		return nil, err
	}

	return &review, nil
}

// GetReviewByID retrieves a review by ID
func (s *ReviewService) GetReviewByID(id uint) (*models.DeliveryReview, error) {
	var review models.DeliveryReview
	err := s.db.Preload("School").
		Preload("DeliveryRecord").
		First(&review, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReviewNotFound
		}
		return nil, err
	}

	return &review, nil
}

// GetAllReviews retrieves all reviews with optional filters
func (s *ReviewService) GetAllReviews(schoolID *uint, startDate, endDate *time.Time, limit, offset int) ([]models.DeliveryReview, int64, error) {
	var reviews []models.DeliveryReview
	var total int64

	query := s.db.Model(&models.DeliveryReview{}).
		Preload("School").
		Preload("DeliveryRecord")

	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}

	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}

	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&reviews).Error
	return reviews, total, err
}

// ReviewSummary represents aggregated review statistics
type ReviewSummary struct {
	TotalReviews         int64   `json:"total_reviews"`
	AverageOverallRating float64 `json:"average_overall_rating"`
	AverageMenuRating    float64 `json:"average_menu_rating"`
	AverageServiceRating float64 `json:"average_service_rating"`
	
	// Individual rating averages
	AvgFoodTaste       float64 `json:"avg_food_taste"`
	AvgFoodCleanliness float64 `json:"avg_food_cleanliness"`
	AvgMenuAccuracy    float64 `json:"avg_menu_accuracy"`
	AvgPortionSize     float64 `json:"avg_portion_size"`
	AvgMenuVariety     float64 `json:"avg_menu_variety"`
	AvgDeliveryTime    float64 `json:"avg_delivery_time"`
	AvgDriverAttitude  float64 `json:"avg_driver_attitude"`
	AvgFoodCondition   float64 `json:"avg_food_condition"`
	AvgDriverTidiness  float64 `json:"avg_driver_tidiness"`
	AvgServiceConsistency float64 `json:"avg_service_consistency"`
	
	// Rating distribution
	RatingDistribution map[int]int `json:"rating_distribution"` // 1-5 stars count
}

// GetReviewSummary retrieves aggregated review statistics
func (s *ReviewService) GetReviewSummary(schoolID *uint, startDate, endDate *time.Time) (*ReviewSummary, error) {
	summary := &ReviewSummary{
		RatingDistribution: make(map[int]int),
	}

	query := s.db.Model(&models.DeliveryReview{})

	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}

	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}

	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}

	// Get aggregated stats
	var result struct {
		Count              int64
		AvgOverall         float64
		AvgMenu            float64
		AvgService         float64
		AvgFoodTaste       float64
		AvgFoodCleanliness float64
		AvgMenuAccuracy    float64
		AvgPortionSize     float64
		AvgMenuVariety     float64
		AvgDeliveryTime    float64
		AvgDriverAttitude  float64
		AvgFoodCondition   float64
		AvgDriverTidiness  float64
		AvgServiceConsistency float64
	}

	err := query.Select(`
		COUNT(*) as count,
		AVG(overall_rating) as avg_overall,
		AVG(average_menu_rating) as avg_menu,
		AVG(average_service_rating) as avg_service,
		AVG(rating_food_taste) as avg_food_taste,
		AVG(rating_food_cleanliness) as avg_food_cleanliness,
		AVG(rating_menu_accuracy) as avg_menu_accuracy,
		AVG(rating_portion_size) as avg_portion_size,
		AVG(rating_menu_variety) as avg_menu_variety,
		AVG(rating_delivery_time) as avg_delivery_time,
		AVG(rating_driver_attitude) as avg_driver_attitude,
		AVG(rating_food_condition) as avg_food_condition,
		AVG(rating_driver_tidiness) as avg_driver_tidiness,
		AVG(rating_service_consistency) as avg_service_consistency
	`).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	summary.TotalReviews = result.Count
	summary.AverageOverallRating = result.AvgOverall
	summary.AverageMenuRating = result.AvgMenu
	summary.AverageServiceRating = result.AvgService
	summary.AvgFoodTaste = result.AvgFoodTaste
	summary.AvgFoodCleanliness = result.AvgFoodCleanliness
	summary.AvgMenuAccuracy = result.AvgMenuAccuracy
	summary.AvgPortionSize = result.AvgPortionSize
	summary.AvgMenuVariety = result.AvgMenuVariety
	summary.AvgDeliveryTime = result.AvgDeliveryTime
	summary.AvgDriverAttitude = result.AvgDriverAttitude
	summary.AvgFoodCondition = result.AvgFoodCondition
	summary.AvgDriverTidiness = result.AvgDriverTidiness
	summary.AvgServiceConsistency = result.AvgServiceConsistency

	// Get rating distribution
	var distributions []struct {
		Rating int
		Count  int
	}

	err = s.db.Model(&models.DeliveryReview{}).
		Select("ROUND(overall_rating) as rating, COUNT(*) as count").
		Group("ROUND(overall_rating)").
		Scan(&distributions).Error

	if err != nil {
		return nil, err
	}

	for _, d := range distributions {
		summary.RatingDistribution[d.Rating] = d.Count
	}

	return summary, nil
}

// HasReviewForDeliveryRecord checks if a review exists for a delivery record
func (s *ReviewService) HasReviewForDeliveryRecord(deliveryRecordID uint) (bool, error) {
	var count int64
	err := s.db.Model(&models.DeliveryReview{}).
		Where("delivery_record_id = ?", deliveryRecordID).
		Count(&count).Error
	return count > 0, err
}
