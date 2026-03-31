package aggregated_dashboard_test

import (
	"testing"
	"time"

	"github.com/erp-sppg/backend/internal/models"
	"github.com/erp-sppg/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAggDashboardTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(
		&models.Yayasan{},
		&models.SPPG{},
		&models.User{},
		&models.MenuPlan{},
		&models.MenuItem{},
		&models.Recipe{},
		&models.DeliveryRecord{},
		&models.DeliveryReview{},
		&models.CashFlowEntry{},
		&models.BudgetTarget{},
		&models.School{},
	)
	require.NoError(t, err)
	return db
}

// seedTestData creates two Yayasans, each with one SPPG, and operational data.
func seedTestData(t *testing.T, db *gorm.DB) (yayasan1, yayasan2 models.Yayasan, sppg1, sppg2 models.SPPG) {
	yayasan1 = models.Yayasan{Kode: "YYS-0001", Nama: "Yayasan Alpha", Email: "alpha@yayasan.id", NPWP: "NPWP-001", IsActive: true}
	yayasan2 = models.Yayasan{Kode: "YYS-0002", Nama: "Yayasan Beta", Email: "beta@yayasan.id", NPWP: "NPWP-002", IsActive: true}
	require.NoError(t, db.Create(&yayasan1).Error)
	require.NoError(t, db.Create(&yayasan2).Error)

	sppg1 = models.SPPG{Kode: "SPPG-0001", Nama: "SPPG Satu", Email: "sppg1@test.id", YayasanID: yayasan1.ID, IsActive: true}
	sppg2 = models.SPPG{Kode: "SPPG-0002", Nama: "SPPG Dua", Email: "sppg2@test.id", YayasanID: yayasan2.ID, IsActive: true}
	require.NoError(t, db.Create(&sppg1).Error)
	require.NoError(t, db.Create(&sppg2).Error)

	// User for created_by
	user := models.User{NIK: "1234", Email: "u@test.com", PasswordHash: "x", FullName: "Test", Role: "chef", SPPGID: &sppg1.ID}
	require.NoError(t, db.Create(&user).Error)

	// Recipe
	recipe := models.Recipe{Name: "Nasi Goreng", SPPGID: &sppg1.ID, CreatedBy: user.ID, IsActive: true}
	require.NoError(t, db.Create(&recipe).Error)

	baseDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)

	// MenuPlan + MenuItem for SPPG1
	mp1 := models.MenuPlan{SPPGID: &sppg1.ID, WeekStart: baseDate, WeekEnd: baseDate.AddDate(0, 0, 6), Status: "approved", CreatedBy: user.ID}
	require.NoError(t, db.Create(&mp1).Error)
	mi1 := models.MenuItem{MenuPlanID: mp1.ID, Date: baseDate, RecipeID: recipe.ID, Portions: 100}
	require.NoError(t, db.Create(&mi1).Error)

	// MenuPlan + MenuItem for SPPG2
	mp2 := models.MenuPlan{SPPGID: &sppg2.ID, WeekStart: baseDate, WeekEnd: baseDate.AddDate(0, 0, 6), Status: "approved", CreatedBy: user.ID}
	require.NoError(t, db.Create(&mp2).Error)
	mi2 := models.MenuItem{MenuPlanID: mp2.ID, Date: baseDate, RecipeID: recipe.ID, Portions: 200}
	require.NoError(t, db.Create(&mi2).Error)

	// School
	school := models.School{Name: "SD Test", SPPGID: &sppg1.ID, Latitude: -6.2, Longitude: 106.8, Category: "SD", IsActive: true}
	require.NoError(t, db.Create(&school).Error)

	// DeliveryRecords for SPPG1: 2 total, 1 received
	dr1 := models.DeliveryRecord{SPPGID: &sppg1.ID, DeliveryDate: baseDate, SchoolID: school.ID, MenuItemID: mi1.ID, Portions: 50, CurrentStatus: "received", CurrentStage: 8}
	dr2 := models.DeliveryRecord{SPPGID: &sppg1.ID, DeliveryDate: baseDate, SchoolID: school.ID, MenuItemID: mi1.ID, Portions: 50, CurrentStatus: "in_progress", CurrentStage: 3}
	require.NoError(t, db.Create(&dr1).Error)
	require.NoError(t, db.Create(&dr2).Error)

	// DeliveryRecords for SPPG2: 1 total, 1 received
	dr3 := models.DeliveryRecord{SPPGID: &sppg2.ID, DeliveryDate: baseDate, SchoolID: school.ID, MenuItemID: mi2.ID, Portions: 200, CurrentStatus: "received", CurrentStage: 8}
	require.NoError(t, db.Create(&dr3).Error)

	// DeliveryReviews
	rev1 := models.DeliveryReview{SPPGID: &sppg1.ID, DeliveryRecordID: dr1.ID, SchoolID: school.ID, OverallRating: 4.0, AverageMenuRating: 4.0, AverageServiceRating: 4.0,
		RatingFoodTaste: 4, RatingFoodCleanliness: 4, RatingMenuAccuracy: 4, RatingPortionSize: 4, RatingMenuVariety: 4,
		RatingDeliveryTime: 4, RatingDriverAttitude: 4, RatingFoodCondition: 4, RatingDriverTidiness: 4, RatingServiceConsistency: 4,
		CreatedAt: baseDate}
	rev2 := models.DeliveryReview{SPPGID: &sppg2.ID, DeliveryRecordID: dr3.ID, SchoolID: school.ID, OverallRating: 5.0, AverageMenuRating: 5.0, AverageServiceRating: 5.0,
		RatingFoodTaste: 5, RatingFoodCleanliness: 5, RatingMenuAccuracy: 5, RatingPortionSize: 5, RatingMenuVariety: 5,
		RatingDeliveryTime: 5, RatingDriverAttitude: 5, RatingFoodCondition: 5, RatingDriverTidiness: 5, RatingServiceConsistency: 5,
		CreatedAt: baseDate}
	require.NoError(t, db.Create(&rev1).Error)
	require.NoError(t, db.Create(&rev2).Error)

	// CashFlowEntries
	cf1 := models.CashFlowEntry{SPPGID: &sppg1.ID, TransactionID: "TX-001", Date: baseDate, Category: "bahan_baku", Type: "expense", Amount: 1000000, CreatedBy: user.ID}
	cf2 := models.CashFlowEntry{SPPGID: &sppg2.ID, TransactionID: "TX-002", Date: baseDate, Category: "bahan_baku", Type: "expense", Amount: 2000000, CreatedBy: user.ID}
	require.NoError(t, db.Create(&cf1).Error)
	require.NoError(t, db.Create(&cf2).Error)

	// BudgetTargets
	bt1 := models.BudgetTarget{SPPGID: &sppg1.ID, Year: 2025, Month: 1, Category: "bahan_baku", Target: 5000000, Actual: 1000000}
	bt2 := models.BudgetTarget{SPPGID: &sppg2.ID, Year: 2025, Month: 1, Category: "bahan_baku", Target: 10000000, Actual: 2000000}
	require.NoError(t, db.Create(&bt1).Error)
	require.NoError(t, db.Create(&bt2).Error)

	return
}

func TestGetKepalaYayasanDashboard(t *testing.T) {
	db := setupAggDashboardTestDB(t)
	yayasan1, _, sppg1, _ := seedTestData(t, db)
	svc := services.NewAggregatedDashboardService(db)

	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)

	t.Run("aggregates all SPPGs under Yayasan", func(t *testing.T) {
		dash, err := svc.GetKepalaYayasanDashboard(yayasan1.ID, start, end, nil)
		require.NoError(t, err)

		assert.Equal(t, yayasan1.ID, dash.YayasanID)
		assert.Equal(t, "Yayasan Alpha", dash.YayasanNama)
		assert.Equal(t, 1, dash.TotalSPPG)
		assert.Len(t, dash.SPPGSummaries, 1)
		assert.Equal(t, sppg1.ID, dash.SPPGSummaries[0].SPPGID)
		assert.Equal(t, 100, dash.SPPGSummaries[0].TotalPortions)

		// Production
		assert.Equal(t, 100, dash.AggregatedProduction.TotalPortions)

		// Delivery: 2 total, 1 received
		assert.Equal(t, 2, dash.AggregatedDelivery.TotalDeliveries)
		assert.Equal(t, 1, dash.AggregatedDelivery.CompletedDeliveries)
		assert.InDelta(t, 50.0, dash.AggregatedDelivery.CompletionRate, 0.1)

		// Financial
		assert.Equal(t, 5000000.0, dash.AggregatedFinancial.TotalBudget)
		assert.Equal(t, 1000000.0, dash.AggregatedFinancial.TotalSpent)

		// Review
		assert.Equal(t, 1, dash.AggregatedReview.TotalReviews)
		assert.InDelta(t, 4.0, dash.AggregatedReview.AverageOverall, 0.01)
	})

	t.Run("drill-down to specific SPPG", func(t *testing.T) {
		dash, err := svc.GetKepalaYayasanDashboard(yayasan1.ID, start, end, &sppg1.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, dash.TotalSPPG)
		assert.Equal(t, 100, dash.AggregatedProduction.TotalPortions)
	})

	t.Run("returns empty metrics for non-existent Yayasan", func(t *testing.T) {
		_, err := svc.GetKepalaYayasanDashboard(9999, start, end, nil)
		assert.Error(t, err) // record not found
	})
}

func TestGetAdminBGNDashboard(t *testing.T) {
	db := setupAggDashboardTestDB(t)
	yayasan1, _, sppg1, _ := seedTestData(t, db)
	svc := services.NewAggregatedDashboardService(db)

	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)

	t.Run("national aggregation", func(t *testing.T) {
		dash, err := svc.GetAdminBGNDashboard(start, end, nil, nil)
		require.NoError(t, err)

		assert.Equal(t, 2, dash.TotalYayasan)
		assert.Equal(t, 2, dash.TotalSPPG)
		assert.Len(t, dash.YayasanSummaries, 2)

		// Total portions across both SPPGs: 100 + 200
		assert.Equal(t, 300, dash.AggregatedProduction.TotalPortions)

		// Deliveries: 3 total, 2 received
		assert.Equal(t, 3, dash.AggregatedDelivery.TotalDeliveries)
		assert.Equal(t, 2, dash.AggregatedDelivery.CompletedDeliveries)

		// Financial: budget 15M, spent 3M
		assert.Equal(t, 15000000.0, dash.AggregatedFinancial.TotalBudget)
		assert.Equal(t, 3000000.0, dash.AggregatedFinancial.TotalSpent)

		// Reviews: 2 total, avg (4+5)/2 = 4.5
		assert.Equal(t, 2, dash.AggregatedReview.TotalReviews)
		assert.InDelta(t, 4.5, dash.AggregatedReview.AverageOverall, 0.01)
	})

	t.Run("filter by Yayasan", func(t *testing.T) {
		dash, err := svc.GetAdminBGNDashboard(start, end, &yayasan1.ID, nil)
		require.NoError(t, err)

		// Only SPPG1 data
		assert.Equal(t, 100, dash.AggregatedProduction.TotalPortions)
		assert.Equal(t, 2, dash.AggregatedDelivery.TotalDeliveries)
	})

	t.Run("filter by SPPG", func(t *testing.T) {
		dash, err := svc.GetAdminBGNDashboard(start, end, nil, &sppg1.ID)
		require.NoError(t, err)

		assert.Equal(t, 100, dash.AggregatedProduction.TotalPortions)
		assert.Equal(t, 1, dash.AggregatedReview.TotalReviews)
	})

	t.Run("date range filtering excludes out-of-range data", func(t *testing.T) {
		// Use a range that doesn't include any data
		farStart := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		farEnd := time.Date(2026, 12, 31, 23, 59, 59, 0, time.UTC)
		dash, err := svc.GetAdminBGNDashboard(farStart, farEnd, nil, nil)
		require.NoError(t, err)

		assert.Equal(t, 0, dash.AggregatedProduction.TotalPortions)
		assert.Equal(t, 0, dash.AggregatedDelivery.TotalDeliveries)
		assert.Equal(t, 0, dash.AggregatedReview.TotalReviews)
	})
}
