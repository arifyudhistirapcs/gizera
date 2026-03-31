package models

import (
	"time"
)

// SOPCategory merepresentasikan kategori utama dalam SOP Dapur MBG
type SOPCategory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Nama      string    `gorm:"size:200;not null;uniqueIndex" json:"nama" validate:"required"`
	Deskripsi string    `gorm:"type:text" json:"deskripsi"`
	Urutan    int       `gorm:"not null;index" json:"urutan"`
	IsActive  bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for SOPCategory
func (SOPCategory) TableName() string {
	return "sop_categories"
}

// SOPChecklistItem merepresentasikan item individual dalam SOP yang dinilai
type SOPChecklistItem struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	SOPCategoryID uint        `gorm:"index;not null" json:"sop_category_id"`
	Nama          string      `gorm:"size:500;not null" json:"nama" validate:"required"`
	Deskripsi     string      `gorm:"type:text" json:"deskripsi"`
	Urutan        int         `gorm:"not null;index" json:"urutan"`
	IsActive      bool        `gorm:"default:true;index" json:"is_active"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	SOPCategory   SOPCategory `gorm:"foreignKey:SOPCategoryID" json:"sop_category,omitempty"`
}

// TableName specifies the table name for SOPChecklistItem
func (SOPChecklistItem) TableName() string {
	return "sop_checklist_items"
}

// RiskAssessmentForm merepresentasikan formulir audit risk assessment
type RiskAssessmentForm struct {
	ID               uint                          `gorm:"primaryKey" json:"id"`
	SPPGID           uint                          `gorm:"index;not null" json:"sppg_id"`
	YayasanID        uint                          `gorm:"index;not null" json:"yayasan_id"`
	CreatedByUserID  uint                          `gorm:"index;not null" json:"created_by_user_id"`
	Status           string                        `gorm:"size:20;not null;index;default:'draft'" json:"status"` // draft, submitted
	OverallRiskScore *float64                      `json:"overall_risk_score"`
	RiskLevel        *string                       `gorm:"size:20" json:"risk_level"` // rendah, sedang, tinggi
	SubmittedAt      *time.Time                    `gorm:"index" json:"submitted_at"`
	CreatedAt        time.Time                     `gorm:"index" json:"created_at"`
	UpdatedAt        time.Time                     `json:"updated_at"`
	SPPG             SPPG                          `gorm:"foreignKey:SPPGID" json:"sppg,omitempty"`
	Yayasan          Yayasan                       `gorm:"foreignKey:YayasanID" json:"yayasan,omitempty"`
	CreatedByUser    User                          `gorm:"foreignKey:CreatedByUserID" json:"created_by_user,omitempty"`
	Items            []RiskAssessmentItem          `gorm:"foreignKey:FormID" json:"items,omitempty"`
	CategoryScores   []RiskAssessmentCategoryScore `gorm:"foreignKey:FormID" json:"category_scores,omitempty"`
	Snapshot         *SPPGOperationalSnapshot      `gorm:"foreignKey:FormID" json:"snapshot,omitempty"`
}

// TableName specifies the table name for RiskAssessmentForm
func (RiskAssessmentForm) TableName() string {
	return "risk_assessment_forms"
}

// RiskAssessmentItem merepresentasikan penilaian per checklist item dalam form
type RiskAssessmentItem struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	FormID             uint      `gorm:"index;not null" json:"form_id"`
	SOPChecklistItemID uint      `gorm:"index;not null" json:"sop_checklist_item_id"`
	SOPCategoryID      uint      `gorm:"index;not null" json:"sop_category_id"`
	ItemNama           string    `gorm:"size:500;not null" json:"item_nama"`           // Snapshot nama item saat form dibuat
	CategoryNama       string    `gorm:"size:200;not null" json:"category_nama"`       // Snapshot nama kategori
	ComplianceScore    *int      `json:"compliance_score" validate:"omitempty,min=1,max=5"` // 1-5, nil = belum dinilai
	Catatan            string    `gorm:"type:text" json:"catatan"`
	EvidenceURL        string    `gorm:"size:500" json:"evidence_url"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// TableName specifies the table name for RiskAssessmentItem
func (RiskAssessmentItem) TableName() string {
	return "risk_assessment_items"
}

// RiskAssessmentCategoryScore menyimpan skor rata-rata per kategori SOP
type RiskAssessmentCategoryScore struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	FormID        uint    `gorm:"index;not null" json:"form_id"`
	SOPCategoryID uint    `gorm:"index;not null" json:"sop_category_id"`
	CategoryNama  string  `gorm:"size:200;not null" json:"category_nama"`
	AverageScore  float64 `gorm:"not null" json:"average_score"`
	RiskLevel     string  `gorm:"size:20;not null" json:"risk_level"` // rendah, sedang, tinggi
	ItemCount     int     `gorm:"not null" json:"item_count"`
}

// TableName specifies the table name for RiskAssessmentCategoryScore
func (RiskAssessmentCategoryScore) TableName() string {
	return "risk_assessment_category_scores"
}

// SPPGOperationalSnapshot menyimpan snapshot data operasional SPPG saat audit
type SPPGOperationalSnapshot struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	FormID uint `gorm:"uniqueIndex;not null" json:"form_id"`

	// Review metrics
	AverageOverallRating float64 `json:"average_overall_rating"`
	AverageMenuRating    float64 `json:"average_menu_rating"`
	AverageServiceRating float64 `json:"average_service_rating"`
	TotalReviews         int     `json:"total_reviews"`

	// Financial metrics
	TotalIncome          float64 `json:"total_income"`
	TotalExpense         float64 `json:"total_expense"`
	BudgetTarget         float64 `json:"budget_target"`
	BudgetAbsorptionRate float64 `json:"budget_absorption_rate"` // percentage

	// Delivery metrics
	TotalDeliveries        int     `json:"total_deliveries"`
	CompletedDeliveries    int     `json:"completed_deliveries"`
	DeliveryCompletionRate float64 `json:"delivery_completion_rate"` // percentage
	OnTimeDeliveryRate     float64 `json:"on_time_delivery_rate"`    // percentage

	// Production metrics
	TotalPortionsProduced    int     `json:"total_portions_produced"`
	ProductionCompletionRate float64 `json:"production_completion_rate"` // percentage

	// Inventory metrics
	TotalInventoryItems int `json:"total_inventory_items"`
	CriticalStockItems  int `json:"critical_stock_items"` // items below min_threshold

	// HR metrics
	TotalActiveEmployees int     `json:"total_active_employees"`
	AttendanceRate       float64 `json:"attendance_rate"` // percentage bulan berjalan

	// Metadata
	SnapshotPeriodStart time.Time `json:"snapshot_period_start"` // awal bulan berjalan
	SnapshotPeriodEnd   time.Time `json:"snapshot_period_end"`   // tanggal capture
	CapturedAt          time.Time `json:"captured_at"`
}

// TableName specifies the table name for SPPGOperationalSnapshot
func (SPPGOperationalSnapshot) TableName() string {
	return "sppg_operational_snapshots"
}
