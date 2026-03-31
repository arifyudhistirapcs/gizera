package services

import (
	"time"

	"github.com/erp-sppg/backend/internal/models"
	"gorm.io/gorm"
)

// AggregatedDashboardService provides aggregated dashboard data
// for Kepala Yayasan (cross-SPPG) and Admin BGN (national).
type AggregatedDashboardService struct {
	db *gorm.DB
}

// NewAggregatedDashboardService creates a new AggregatedDashboardService.
func NewAggregatedDashboardService(db *gorm.DB) *AggregatedDashboardService {
	return &AggregatedDashboardService{db: db}
}

// GetKepalaYayasanDashboard returns an aggregated dashboard for a Kepala Yayasan.
// It aggregates production, delivery, financial, and review metrics from all SPPGs
// under the given Yayasan. An optional sppgID drills down to a single SPPG.
func (s *AggregatedDashboardService) GetKepalaYayasanDashboard(
	yayasanID uint,
	startDate, endDate time.Time,
	sppgID *uint,
) (*models.KepalaYayasanAggregatedDashboard, error) {
	// Fetch Yayasan info
	var yayasan models.Yayasan
	if err := s.db.First(&yayasan, yayasanID).Error; err != nil {
		return nil, err
	}

	// Fetch SPPGs under this Yayasan
	var sppgs []models.SPPG
	sppgQuery := s.db.Where("yayasan_id = ? AND is_active = ?", yayasanID, true)
	if sppgID != nil {
		sppgQuery = sppgQuery.Where("id = ?", *sppgID)
	}
	if err := sppgQuery.Find(&sppgs).Error; err != nil {
		return nil, err
	}

	sppgIDs := make([]uint, len(sppgs))
	for i, sp := range sppgs {
		sppgIDs[i] = sp.ID
	}

	dashboard := &models.KepalaYayasanAggregatedDashboard{
		YayasanID:   yayasan.ID,
		YayasanNama: yayasan.Nama,
		TotalSPPG:   len(sppgs),
		UpdatedAt:   time.Now(),
	}

	if len(sppgIDs) == 0 {
		dashboard.AggregatedProduction = &models.AggregatedProduction{}
		dashboard.AggregatedDelivery = &models.AggregatedDelivery{}
		dashboard.AggregatedFinancial = &models.AggregatedFinancial{}
		dashboard.AggregatedReview = &models.AggregatedReview{}
		return dashboard, nil
	}

	// Build per-SPPG summaries
	summaries := make([]models.SPPGSummary, 0, len(sppgs))
	for _, sp := range sppgs {
		ids := []uint{sp.ID}
		summary := models.SPPGSummary{
			SPPGID:   sp.ID,
			SPPGNama: sp.Nama,
			SPPGKode: sp.Kode,
		}
		summary.TotalPortions = s.getTotalPortions(ids, startDate, endDate)
		summary.DeliveryRate = s.getDeliveryRate(ids, startDate, endDate)
		summary.BudgetAbsorption = s.getBudgetAbsorption(ids, startDate, endDate)
		summary.AverageReviewRating = s.getAverageReviewRating(ids, startDate, endDate)
		summaries = append(summaries, summary)
	}
	dashboard.SPPGSummaries = summaries

	// Aggregated metrics across all selected SPPGs
	dashboard.AggregatedProduction = s.getAggregatedProduction(sppgIDs, startDate, endDate)
	dashboard.AggregatedDelivery = s.getAggregatedDelivery(sppgIDs, startDate, endDate)
	dashboard.AggregatedFinancial = s.getAggregatedFinancial(sppgIDs, startDate, endDate)
	dashboard.AggregatedReview = s.getAggregatedReview(sppgIDs, startDate, endDate)

	return dashboard, nil
}

// GetAdminBGNDashboard returns a national-level aggregated dashboard for Admin BGN.
// It aggregates metrics from all SPPGs across all Yayasans.
// Optional yayasanID and sppgID allow drill-down filtering.
func (s *AggregatedDashboardService) GetAdminBGNDashboard(
	startDate, endDate time.Time,
	yayasanID *uint,
	sppgID *uint,
) (*models.AdminBGNDashboard, error) {
	// Count totals
	var totalYayasan int64
	s.db.Model(&models.Yayasan{}).Where("is_active = ?", true).Count(&totalYayasan)

	var totalSPPG int64
	s.db.Model(&models.SPPG{}).Where("is_active = ?", true).Count(&totalSPPG)

	// Determine which SPPGs to include based on filters
	sppgIDs, err := s.resolveSPPGIDs(yayasanID, sppgID)
	if err != nil {
		return nil, err
	}

	dashboard := &models.AdminBGNDashboard{
		TotalYayasan: int(totalYayasan),
		TotalSPPG:    int(totalSPPG),
		UpdatedAt:    time.Now(),
	}

	// Build Yayasan summaries
	var yayasans []models.Yayasan
	yayasanQuery := s.db.Where("is_active = ?", true)
	if yayasanID != nil {
		yayasanQuery = yayasanQuery.Where("id = ?", *yayasanID)
	}
	if err := yayasanQuery.Preload("SPPGs", "is_active = ?", true).Find(&yayasans).Error; err != nil {
		return nil, err
	}

	yayasanSummaries := make([]models.YayasanSummary, 0, len(yayasans))
	for _, y := range yayasans {
		yIDs := make([]uint, 0, len(y.SPPGs))
		for _, sp := range y.SPPGs {
			yIDs = append(yIDs, sp.ID)
		}
		summary := models.YayasanSummary{
			YayasanID:   y.ID,
			YayasanNama: y.Nama,
			YayasanKode: y.Kode,
			TotalSPPG:   len(y.SPPGs),
		}
		if len(yIDs) > 0 {
			summary.TotalPortions = s.getTotalPortions(yIDs, startDate, endDate)
			summary.TotalSpent = s.getTotalSpent(yIDs, startDate, endDate)
			summary.AverageReviewRating = s.getAverageReviewRating(yIDs, startDate, endDate)
		}
		yayasanSummaries = append(yayasanSummaries, summary)
	}
	dashboard.YayasanSummaries = yayasanSummaries

	// Build per-SPPG summaries
	var allSPPGs []models.SPPG
	sppgQuery := s.db.Where("is_active = ?", true)
	if yayasanID != nil {
		sppgQuery = sppgQuery.Where("yayasan_id = ?", *yayasanID)
	}
	if sppgID != nil {
		sppgQuery = sppgQuery.Where("id = ?", *sppgID)
	}
	sppgQuery.Find(&allSPPGs)

	sppgSummaries := make([]models.SPPGSummary, 0, len(allSPPGs))
	for _, sp := range allSPPGs {
		ids := []uint{sp.ID}
		sppgSummaries = append(sppgSummaries, models.SPPGSummary{
			SPPGID:              sp.ID,
			SPPGNama:            sp.Nama,
			SPPGKode:            sp.Kode,
			TotalPortions:       s.getTotalPortions(ids, startDate, endDate),
			DeliveryRate:        s.getDeliveryRate(ids, startDate, endDate),
			BudgetAbsorption:    s.getBudgetAbsorption(ids, startDate, endDate),
			AverageReviewRating: s.getAverageReviewRating(ids, startDate, endDate),
		})
	}
	dashboard.SPPGSummaries = sppgSummaries

	if len(sppgIDs) == 0 {
		dashboard.AggregatedProduction = &models.AggregatedProduction{}
		dashboard.AggregatedDelivery = &models.AggregatedDelivery{}
		dashboard.AggregatedFinancial = &models.AggregatedFinancial{}
		dashboard.AggregatedReview = &models.AggregatedReview{}
		return dashboard, nil
	}

	// Aggregated metrics
	dashboard.AggregatedProduction = s.getAggregatedProduction(sppgIDs, startDate, endDate)
	dashboard.AggregatedDelivery = s.getAggregatedDelivery(sppgIDs, startDate, endDate)
	dashboard.AggregatedFinancial = s.getAggregatedFinancial(sppgIDs, startDate, endDate)
	dashboard.AggregatedReview = s.getAggregatedReview(sppgIDs, startDate, endDate)

	return dashboard, nil
}

// resolveSPPGIDs returns the list of SPPG IDs to include based on optional filters.
func (s *AggregatedDashboardService) resolveSPPGIDs(yayasanID *uint, sppgID *uint) ([]uint, error) {
	if sppgID != nil {
		// Drill-down to a single SPPG
		return []uint{*sppgID}, nil
	}

	query := s.db.Model(&models.SPPG{}).Where("is_active = ?", true)
	if yayasanID != nil {
		query = query.Where("yayasan_id = ?", *yayasanID)
	}

	var ids []uint
	if err := query.Pluck("id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

// ---------------------------------------------------------------------------
// Per-SPPG summary helpers
// ---------------------------------------------------------------------------

// getTotalPortions returns the sum of portions from menu_items within the date range
// for the given SPPG IDs.
func (s *AggregatedDashboardService) getTotalPortions(sppgIDs []uint, start, end time.Time) int {
	var total int64
	s.db.Model(&models.MenuItem{}).
		Joins("JOIN menu_plans ON menu_plans.id = menu_items.menu_plan_id").
		Where("menu_plans.sppg_id IN ?", sppgIDs).
		Where("menu_items.date BETWEEN ? AND ?", start, end).
		Select("COALESCE(SUM(menu_items.portions), 0)").
		Scan(&total)
	return int(total)
}

// getDeliveryRate returns the completion rate of deliveries (completed / total).
func (s *AggregatedDashboardService) getDeliveryRate(sppgIDs []uint, start, end time.Time) float64 {
	var total int64
	var completed int64
	s.db.Model(&models.DeliveryRecord{}).
		Where("sppg_id IN ? AND delivery_date BETWEEN ? AND ?", sppgIDs, start, end).
		Count(&total)
	s.db.Model(&models.DeliveryRecord{}).
		Where("sppg_id IN ? AND delivery_date BETWEEN ? AND ?", sppgIDs, start, end).
		Where("current_status = ?", "received").
		Count(&completed)
	if total == 0 {
		return 0
	}
	return float64(completed) / float64(total) * 100
}

// getBudgetAbsorption returns (actual / target) * 100 for budget targets.
func (s *AggregatedDashboardService) getBudgetAbsorption(sppgIDs []uint, start, end time.Time) float64 {
	type result struct {
		TotalTarget float64
		TotalActual float64
	}
	var r result
	s.db.Model(&models.BudgetTarget{}).
		Where("sppg_id IN ?", sppgIDs).
		Where("(year > ? OR (year = ? AND month >= ?)) AND (year < ? OR (year = ? AND month <= ?))",
			start.Year(), start.Year(), int(start.Month()),
			end.Year(), end.Year(), int(end.Month())).
		Select("COALESCE(SUM(target), 0) as total_target, COALESCE(SUM(actual), 0) as total_actual").
		Scan(&r)
	if r.TotalTarget == 0 {
		return 0
	}
	return r.TotalActual / r.TotalTarget * 100
}

// getAverageReviewRating returns the average overall_rating from delivery reviews.
func (s *AggregatedDashboardService) getAverageReviewRating(sppgIDs []uint, start, end time.Time) float64 {
	var avg float64
	s.db.Model(&models.DeliveryReview{}).
		Where("sppg_id IN ? AND created_at BETWEEN ? AND ?", sppgIDs, start, end).
		Select("COALESCE(AVG(overall_rating), 0)").
		Scan(&avg)
	return avg
}

// getTotalSpent returns the total expense amount from cash flow entries.
func (s *AggregatedDashboardService) getTotalSpent(sppgIDs []uint, start, end time.Time) float64 {
	var total float64
	s.db.Model(&models.CashFlowEntry{}).
		Where("sppg_id IN ? AND type = ? AND date BETWEEN ? AND ?", sppgIDs, "expense", start, end).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total)
	return total
}

// ---------------------------------------------------------------------------
// Aggregated metric builders
// ---------------------------------------------------------------------------

// getAggregatedProduction builds production metrics across the given SPPGs.
func (s *AggregatedDashboardService) getAggregatedProduction(sppgIDs []uint, start, end time.Time) *models.AggregatedProduction {
	totalPortions := s.getTotalPortions(sppgIDs, start, end)

	var totalRecipes int64
	s.db.Model(&models.MenuItem{}).
		Joins("JOIN menu_plans ON menu_plans.id = menu_items.menu_plan_id").
		Where("menu_plans.sppg_id IN ?", sppgIDs).
		Where("menu_items.date BETWEEN ? AND ?", start, end).
		Distinct("menu_items.recipe_id").
		Count(&totalRecipes)

	// Recipes that have at least one delivery record with status "received"
	var recipesCompleted int64
	s.db.Model(&models.MenuItem{}).
		Joins("JOIN menu_plans ON menu_plans.id = menu_items.menu_plan_id").
		Joins("JOIN delivery_records ON delivery_records.menu_item_id = menu_items.id").
		Where("menu_plans.sppg_id IN ?", sppgIDs).
		Where("menu_items.date BETWEEN ? AND ?", start, end).
		Where("delivery_records.current_status = ?", "received").
		Distinct("menu_items.recipe_id").
		Count(&recipesCompleted)

	completionRate := float64(0)
	if totalRecipes > 0 {
		completionRate = float64(recipesCompleted) / float64(totalRecipes) * 100
	}

	return &models.AggregatedProduction{
		TotalPortions:    totalPortions,
		CompletionRate:   completionRate,
		TotalRecipes:     int(totalRecipes),
		RecipesCompleted: int(recipesCompleted),
	}
}

// getAggregatedDelivery builds delivery metrics across the given SPPGs.
func (s *AggregatedDashboardService) getAggregatedDelivery(sppgIDs []uint, start, end time.Time) *models.AggregatedDelivery {
	var totalDeliveries int64
	var completedDeliveries int64

	s.db.Model(&models.DeliveryRecord{}).
		Where("sppg_id IN ? AND delivery_date BETWEEN ? AND ?", sppgIDs, start, end).
		Count(&totalDeliveries)

	s.db.Model(&models.DeliveryRecord{}).
		Where("sppg_id IN ? AND delivery_date BETWEEN ? AND ?", sppgIDs, start, end).
		Where("current_status = ?", "received").
		Count(&completedDeliveries)

	completionRate := float64(0)
	if totalDeliveries > 0 {
		completionRate = float64(completedDeliveries) / float64(totalDeliveries) * 100
	}

	// On-time rate: deliveries that reached "received" status
	// (simplified — same as completion rate since we don't track expected delivery time)
	onTimeRate := completionRate

	return &models.AggregatedDelivery{
		TotalDeliveries:     int(totalDeliveries),
		CompletedDeliveries: int(completedDeliveries),
		OnTimeRate:          onTimeRate,
		CompletionRate:      completionRate,
	}
}

// getAggregatedFinancial builds financial metrics across the given SPPGs.
func (s *AggregatedDashboardService) getAggregatedFinancial(sppgIDs []uint, start, end time.Time) *models.AggregatedFinancial {
	type budgetResult struct {
		TotalTarget float64
		TotalActual float64
	}
	var br budgetResult
	s.db.Model(&models.BudgetTarget{}).
		Where("sppg_id IN ?", sppgIDs).
		Where("(year > ? OR (year = ? AND month >= ?)) AND (year < ? OR (year = ? AND month <= ?))",
			start.Year(), start.Year(), int(start.Month()),
			end.Year(), end.Year(), int(end.Month())).
		Select("COALESCE(SUM(target), 0) as total_target, COALESCE(SUM(actual), 0) as total_actual").
		Scan(&br)

	totalSpent := s.getTotalSpent(sppgIDs, start, end)

	absorptionRate := float64(0)
	if br.TotalTarget > 0 {
		absorptionRate = totalSpent / br.TotalTarget * 100
	}

	return &models.AggregatedFinancial{
		TotalBudget:    br.TotalTarget,
		TotalSpent:     totalSpent,
		AbsorptionRate: absorptionRate,
	}
}

// getAggregatedReview builds review metrics across the given SPPGs.
func (s *AggregatedDashboardService) getAggregatedReview(sppgIDs []uint, start, end time.Time) *models.AggregatedReview {
	type reviewResult struct {
		TotalReviews         int64
		AverageOverall       float64
		AverageMenuRating    float64
		AverageServiceRating float64
	}
	var rr reviewResult
	s.db.Model(&models.DeliveryReview{}).
		Where("sppg_id IN ? AND created_at BETWEEN ? AND ?", sppgIDs, start, end).
		Select(`
			COUNT(*) as total_reviews,
			COALESCE(AVG(overall_rating), 0) as average_overall,
			COALESCE(AVG(average_menu_rating), 0) as average_menu_rating,
			COALESCE(AVG(average_service_rating), 0) as average_service_rating
		`).
		Scan(&rr)

	return &models.AggregatedReview{
		TotalReviews:         int(rr.TotalReviews),
		AverageOverall:       rr.AverageOverall,
		AverageMenuRating:    rr.AverageMenuRating,
		AverageServiceRating: rr.AverageServiceRating,
	}
}
