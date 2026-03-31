package models

import (
	"time"
)

// KepalaYayasanAggregatedDashboard — dashboard agregasi lintas SPPG untuk Kepala Yayasan
type KepalaYayasanAggregatedDashboard struct {
	YayasanID            uint                  `json:"yayasan_id"`
	YayasanNama          string                `json:"yayasan_nama"`
	TotalSPPG            int                   `json:"total_sppg"`
	SPPGSummaries        []SPPGSummary         `json:"sppg_summaries"`
	AggregatedProduction *AggregatedProduction `json:"aggregated_production"`
	AggregatedDelivery   *AggregatedDelivery   `json:"aggregated_delivery"`
	AggregatedFinancial  *AggregatedFinancial  `json:"aggregated_financial"`
	AggregatedReview     *AggregatedReview     `json:"aggregated_review"`
	UpdatedAt            time.Time             `json:"updated_at"`
}

// SPPGSummary — ringkasan performa per SPPG
type SPPGSummary struct {
	SPPGID              uint    `json:"sppg_id"`
	SPPGNama            string  `json:"sppg_nama"`
	SPPGKode            string  `json:"sppg_kode"`
	TotalPortions       int     `json:"total_portions"`
	DeliveryRate        float64 `json:"delivery_rate"`
	BudgetAbsorption    float64 `json:"budget_absorption"`
	AverageReviewRating float64 `json:"average_review_rating"`
}

// AggregatedProduction — metrik produksi agregat
type AggregatedProduction struct {
	TotalPortions    int     `json:"total_portions"`
	CompletionRate   float64 `json:"completion_rate"`
	TotalRecipes     int     `json:"total_recipes"`
	RecipesCompleted int     `json:"recipes_completed"`
}

// AggregatedDelivery — metrik pengiriman agregat
type AggregatedDelivery struct {
	TotalDeliveries     int     `json:"total_deliveries"`
	CompletedDeliveries int     `json:"completed_deliveries"`
	OnTimeRate          float64 `json:"on_time_rate"`
	CompletionRate      float64 `json:"completion_rate"`
}

// AggregatedFinancial — metrik keuangan agregat
type AggregatedFinancial struct {
	TotalBudget    float64 `json:"total_budget"`
	TotalSpent     float64 `json:"total_spent"`
	AbsorptionRate float64 `json:"absorption_rate"`
}

// AggregatedReview — metrik ulasan agregat
type AggregatedReview struct {
	TotalReviews         int     `json:"total_reviews"`
	AverageOverall       float64 `json:"average_overall"`
	AverageMenuRating    float64 `json:"average_menu_rating"`
	AverageServiceRating float64 `json:"average_service_rating"`
}

// AdminBGNDashboard — dashboard agregasi nasional untuk Admin BGN
type AdminBGNDashboard struct {
	TotalYayasan         int                   `json:"total_yayasan"`
	TotalSPPG            int                   `json:"total_sppg"`
	YayasanSummaries     []YayasanSummary      `json:"yayasan_summaries"`
	SPPGSummaries        []SPPGSummary         `json:"sppg_summaries"`
	AggregatedProduction *AggregatedProduction `json:"aggregated_production"`
	AggregatedDelivery   *AggregatedDelivery   `json:"aggregated_delivery"`
	AggregatedFinancial  *AggregatedFinancial  `json:"aggregated_financial"`
	AggregatedReview     *AggregatedReview     `json:"aggregated_review"`
	UpdatedAt            time.Time             `json:"updated_at"`
}

// YayasanSummary — ringkasan performa per Yayasan
type YayasanSummary struct {
	YayasanID           uint    `json:"yayasan_id"`
	YayasanNama         string  `json:"yayasan_nama"`
	YayasanKode         string  `json:"yayasan_kode"`
	TotalSPPG           int     `json:"total_sppg"`
	TotalPortions       int     `json:"total_portions"`
	TotalSpent          float64 `json:"total_spent"`
	AverageReviewRating float64 `json:"average_review_rating"`
}
