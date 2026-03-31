package services

import (
	"time"

	"github.com/erp-sppg/backend/internal/models"
	"gorm.io/gorm"
)

// SnapshotService captures SPPG operational data snapshots for risk assessment forms.
// It leverages AggregatedDashboardService for review, financial, delivery, and production
// metrics, and queries inventory and HR data directly.
type SnapshotService struct {
	db                *gorm.DB
	aggregatedService *AggregatedDashboardService
}

// NewSnapshotService creates a new SnapshotService.
func NewSnapshotService(db *gorm.DB, aggregatedService *AggregatedDashboardService) *SnapshotService {
	return &SnapshotService{
		db:                db,
		aggregatedService: aggregatedService,
	}
}

// CaptureSnapshot captures a point-in-time snapshot of operational metrics for the given SPPG.
// The snapshot period covers from the first day of the current month until now.
func (s *SnapshotService) CaptureSnapshot(sppgID uint) (*models.SPPGOperationalSnapshot, error) {
	now := time.Now()
	periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	periodEnd := now

	sppgIDs := []uint{sppgID}

	// Review metrics via AggregatedDashboardService
	review := s.aggregatedService.getAggregatedReview(sppgIDs, periodStart, periodEnd)

	// Financial metrics via AggregatedDashboardService
	financial := s.aggregatedService.getAggregatedFinancial(sppgIDs, periodStart, periodEnd)

	// Delivery metrics via AggregatedDashboardService
	delivery := s.aggregatedService.getAggregatedDelivery(sppgIDs, periodStart, periodEnd)

	// Production metrics via AggregatedDashboardService
	production := s.aggregatedService.getAggregatedProduction(sppgIDs, periodStart, periodEnd)

	// Inventory metrics — query directly
	totalInventory, criticalStock := s.getInventoryMetrics(sppgID)

	// HR metrics — query directly
	activeEmployees, attendanceRate := s.getHRMetrics(sppgID, periodStart, periodEnd)

	// Compute income for the period
	totalIncome := s.getTotalIncome(sppgID, periodStart, periodEnd)

	snapshot := &models.SPPGOperationalSnapshot{
		// Review
		AverageOverallRating: review.AverageOverall,
		AverageMenuRating:    review.AverageMenuRating,
		AverageServiceRating: review.AverageServiceRating,
		TotalReviews:         review.TotalReviews,

		// Financial
		TotalIncome:          totalIncome,
		TotalExpense:         financial.TotalSpent,
		BudgetTarget:         financial.TotalBudget,
		BudgetAbsorptionRate: financial.AbsorptionRate,

		// Delivery
		TotalDeliveries:        delivery.TotalDeliveries,
		CompletedDeliveries:    delivery.CompletedDeliveries,
		DeliveryCompletionRate: delivery.CompletionRate,
		OnTimeDeliveryRate:     delivery.OnTimeRate,

		// Production
		TotalPortionsProduced:    production.TotalPortions,
		ProductionCompletionRate: production.CompletionRate,

		// Inventory
		TotalInventoryItems: totalInventory,
		CriticalStockItems:  criticalStock,

		// HR
		TotalActiveEmployees: activeEmployees,
		AttendanceRate:       attendanceRate,

		// Metadata
		SnapshotPeriodStart: periodStart,
		SnapshotPeriodEnd:   periodEnd,
		CapturedAt:          now,
	}

	return snapshot, nil
}

// getInventoryMetrics returns the total inventory item count and the number of items
// with quantity below their min_threshold for the given SPPG.
func (s *SnapshotService) getInventoryMetrics(sppgID uint) (totalItems int, criticalItems int) {
	var total int64
	s.db.Model(&models.InventoryItem{}).
		Where("sppg_id = ?", sppgID).
		Count(&total)

	var critical int64
	s.db.Model(&models.InventoryItem{}).
		Where("sppg_id = ? AND quantity < min_threshold", sppgID).
		Count(&critical)

	return int(total), int(critical)
}

// getHRMetrics returns the active employee count and attendance rate (percentage)
// for the given SPPG within the specified period.
func (s *SnapshotService) getHRMetrics(sppgID uint, periodStart, periodEnd time.Time) (activeEmployees int, attendanceRate float64) {
	var employeeCount int64
	s.db.Model(&models.Employee{}).
		Where("sppg_id = ? AND is_active = ?", sppgID, true).
		Count(&employeeCount)

	if employeeCount == 0 {
		return 0, 0
	}

	// Calculate attendance rate: (attendance records / (active employees * working days)) * 100
	var attendanceCount int64
	s.db.Model(&models.Attendance{}).
		Where("sppg_id = ? AND date >= ? AND date <= ?", sppgID, periodStart, periodEnd).
		Count(&attendanceCount)

	// Count distinct working days in the period that have any attendance
	var workingDays int64
	s.db.Model(&models.Attendance{}).
		Where("sppg_id = ? AND date >= ? AND date <= ?", sppgID, periodStart, periodEnd).
		Distinct("date").
		Count(&workingDays)

	if workingDays == 0 {
		return int(employeeCount), 0
	}

	// Attendance rate = actual attendance records / (employees * working days) * 100
	expectedAttendance := float64(employeeCount) * float64(workingDays)
	rate := float64(attendanceCount) / expectedAttendance * 100
	if rate > 100 {
		rate = 100
	}

	return int(employeeCount), rate
}

// getTotalIncome returns the total income amount for the given SPPG within the period.
func (s *SnapshotService) getTotalIncome(sppgID uint, start, end time.Time) float64 {
	var total float64
	s.db.Model(&models.CashFlowEntry{}).
		Where("sppg_id = ? AND type = ? AND date BETWEEN ? AND ?", sppgID, "income", start, end).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total)
	return total
}
