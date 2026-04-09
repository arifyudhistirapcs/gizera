package services

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/erp-sppg/backend/internal/models"
	"gorm.io/gorm"
)

// Risk assessment service errors
var (
	ErrRASPPGNotFound         = errors.New("SPPG tidak ditemukan atau tidak berada di bawah Yayasan Anda")
	ErrRAFormNotFound         = errors.New("formulir risk assessment tidak ditemukan")
	ErrRAFormAlreadySubmitted = errors.New("formulir yang sudah disubmit tidak dapat diubah")
	ErrRAInvalidScore         = errors.New("compliance score harus antara 1 dan 5")
	ErrRAIncompleteScores     = errors.New("semua item checklist harus memiliki skor sebelum submit")
	ErrRASnapshotError        = errors.New("gagal mengambil data operasional SPPG")
	ErrSOPCategoryNotFound    = errors.New("kategori SOP tidak ditemukan")
	ErrSOPChecklistItemNotFound = errors.New("item checklist SOP tidak ditemukan")
	ErrRAFormNotDraft         = errors.New("hanya formulir draft yang dapat dihapus")
)

// FormFilter defines filters for querying risk assessment forms
type FormFilter struct {
	YayasanID uint
	SPPGID    uint
	Status    string
	RiskLevel string
	DateFrom  *time.Time
	DateTo    *time.Time
	Page      int
	PageSize  int
}

// UpdateItemRequest represents a request to update a single risk assessment item
type UpdateItemRequest struct {
	ItemID          uint   `json:"item_id"`
	ComplianceScore *int   `json:"compliance_score"`
	Catatan         string `json:"catatan"`
	EvidenceURL     string `json:"evidence_url"`
}

// RiskAssessmentStats holds aggregated stats for a single SPPG
type RiskAssessmentStats struct {
	SPPGID       uint     `json:"sppg_id"`
	SPPGNama     string   `json:"sppg_nama"`
	TotalAudits  int      `json:"total_audits"`
	AverageScore *float64 `json:"average_score"`
	LatestLevel  *string  `json:"latest_level"`
}

// CreateSOPCategoryInput represents input for creating a new SOP category
type CreateSOPCategoryInput struct {
	Nama      string `json:"nama" validate:"required"`
	Deskripsi string `json:"deskripsi"`
}

// UpdateSOPCategoryInput represents input for updating an SOP category
type UpdateSOPCategoryInput struct {
	Nama      string `json:"nama"`
	Deskripsi string `json:"deskripsi"`
	Urutan    *int   `json:"urutan"`
}

// CreateSOPChecklistItemInput represents input for creating a new SOP checklist item
type CreateSOPChecklistItemInput struct {
	SOPCategoryID uint   `json:"sop_category_id" validate:"required"`
	Nama          string `json:"nama" validate:"required"`
	Deskripsi     string `json:"deskripsi"`
}

// UpdateSOPChecklistItemInput represents input for updating an SOP checklist item
type UpdateSOPChecklistItemInput struct {
	Nama      string `json:"nama"`
	Deskripsi string `json:"deskripsi"`
	Urutan    *int   `json:"urutan"`
}

// RiskAssessmentService handles risk assessment business logic
type RiskAssessmentService struct {
	db              *gorm.DB
	snapshotService *SnapshotService
}

// NewRiskAssessmentService creates a new RiskAssessmentService
func NewRiskAssessmentService(db *gorm.DB, snapshotService *SnapshotService) *RiskAssessmentService {
	return &RiskAssessmentService{
		db:              db,
		snapshotService: snapshotService,
	}
}

// WithDB returns a new RiskAssessmentService with the given DB (for tenant scoping)
func (s *RiskAssessmentService) WithDB(db *gorm.DB) *RiskAssessmentService {
	return &RiskAssessmentService{
		db:              db,
		snapshotService: s.snapshotService,
	}
}

// CreateForm creates a new risk assessment form for a given SPPG.
// It validates SPPG ownership, fetches active checklist items, creates the form
// with snapshot item/category names, and captures an operational snapshot.
func (s *RiskAssessmentService) CreateForm(sppgID, yayasanID, createdByUserID uint) (*models.RiskAssessmentForm, error) {
	// Validate SPPG belongs to the user's Yayasan
	var sppg models.SPPG
	if err := s.db.Where("id = ? AND yayasan_id = ?", sppgID, yayasanID).First(&sppg).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRASPPGNotFound
		}
		return nil, err
	}

	// Fetch all active checklist items with their categories
	var checklistItems []models.SOPChecklistItem
	if err := s.db.
		Preload("SOPCategory").
		Where("is_active = ?", true).
		Order("sop_category_id ASC, urutan ASC").
		Find(&checklistItems).Error; err != nil {
		return nil, err
	}

	// Capture operational snapshot
	var snapshot *models.SPPGOperationalSnapshot
	if s.snapshotService != nil {
		var err error
		snapshot, err = s.snapshotService.CaptureSnapshot(sppgID)
		if err != nil {
			return nil, ErrRASnapshotError
		}
	}

	// Create form with items in a transaction
	form := &models.RiskAssessmentForm{
		SPPGID:          sppgID,
		YayasanID:       yayasanID,
		CreatedByUserID: createdByUserID,
		Status:          "draft",
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Create the form
		if err := tx.Create(form).Error; err != nil {
			return err
		}

		// Create assessment items with snapshot names
		for _, item := range checklistItems {
			assessmentItem := models.RiskAssessmentItem{
				FormID:             form.ID,
				SOPChecklistItemID: item.ID,
				SOPCategoryID:      item.SOPCategoryID,
				ItemNama:           item.Nama,
				CategoryNama:       item.SOPCategory.Nama,
			}
			if err := tx.Create(&assessmentItem).Error; err != nil {
				return err
			}
		}

		// Save snapshot linked to this form
		if snapshot != nil {
			snapshot.FormID = form.ID
			if err := tx.Create(snapshot).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Reload form with all associations
	return s.GetForm(form.ID, yayasanID)
}

// GetForm retrieves a single risk assessment form with all associations.
// Applies tenant filter via yayasanID.
func (s *RiskAssessmentService) GetForm(formID, yayasanID uint) (*models.RiskAssessmentForm, error) {
	var form models.RiskAssessmentForm

	query := s.db.
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Order("sop_category_id ASC, id ASC")
		}).
		Preload("CategoryScores").
		Preload("Snapshot").
		Preload("SPPG").
		Preload("CreatedByUser")

	// Apply tenant filter: yayasanID 0 means superadmin (no filter)
	if yayasanID > 0 {
		query = query.Where("yayasan_id = ?", yayasanID)
	}

	if err := query.First(&form, formID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRAFormNotFound
		}
		return nil, err
	}

	return &form, nil
}

// ListForms returns a paginated list of risk assessment forms with tenant filtering.
// Results are ordered by created_at DESC (newest first).
func (s *RiskAssessmentService) ListForms(filter FormFilter) ([]models.RiskAssessmentForm, int64, error) {
	var forms []models.RiskAssessmentForm
	var totalCount int64

	query := s.db.Model(&models.RiskAssessmentForm{})

	// Tenant filter
	if filter.YayasanID > 0 {
		query = query.Where("yayasan_id = ?", filter.YayasanID)
	}

	// SPPG filter
	if filter.SPPGID > 0 {
		query = query.Where("sppg_id = ?", filter.SPPGID)
	}

	// Status filter
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// Risk level filter
	if filter.RiskLevel != "" {
		query = query.Where("risk_level = ?", filter.RiskLevel)
	}

	// Date range filters
	if filter.DateFrom != nil {
		query = query.Where("created_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("created_at <= ?", *filter.DateTo)
	}

	// Get total count before pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Pagination defaults
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// Fetch with ordering and preloads
	err := query.
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Preload("SPPG").
		Preload("CreatedByUser").
		Find(&forms).Error

	if err != nil {
		return nil, 0, err
	}

	return forms, totalCount, nil
}

// UpdateItemEvidence updates the evidence_url for a specific item in a form.
func (s *RiskAssessmentService) UpdateItemEvidence(formID, itemID uint, evidenceURL string) error {
	result := s.db.Model(&models.RiskAssessmentItem{}).
		Where("id = ? AND form_id = ?", itemID, formID).
		Update("evidence_url", evidenceURL)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRAFormNotFound
	}
	return nil
}

// UpdateDraft updates items in a draft risk assessment form.
// Validates that the form is still in draft status and that scores are in range 1-5.
func (s *RiskAssessmentService) UpdateDraft(formID, yayasanID uint, items []UpdateItemRequest) error {
	// Load the form with tenant filter
	var form models.RiskAssessmentForm
	query := s.db.Where("id = ?", formID)
	if yayasanID > 0 {
		query = query.Where("yayasan_id = ?", yayasanID)
	}
	if err := query.First(&form).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRAFormNotFound
		}
		return err
	}

	// Validate form is still draft
	if form.Status != "draft" {
		return ErrRAFormAlreadySubmitted
	}

	// Validate scores
	for _, item := range items {
		if item.ComplianceScore != nil {
			score := *item.ComplianceScore
			if score < 1 || score > 5 {
				return ErrRAInvalidScore
			}
		}
	}

	// Update items in a transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			updates := map[string]interface{}{
				"catatan": item.Catatan,
			}
			// Only update evidence_url if explicitly provided (non-empty)
			if item.EvidenceURL != "" {
				updates["evidence_url"] = item.EvidenceURL
			}
			if item.ComplianceScore != nil {
				updates["compliance_score"] = *item.ComplianceScore
			}

			result := tx.Model(&models.RiskAssessmentItem{}).
				Where("id = ? AND form_id = ?", item.ItemID, formID).
				Updates(updates)

			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return fmt.Errorf("item dengan ID %d tidak ditemukan dalam formulir ini", item.ItemID)
			}
		}
		return nil
	})
}

// GetStats returns aggregated statistics per SPPG: total audits, average score, and latest risk level.
func (s *RiskAssessmentService) GetStats(sppgIDs []uint) ([]RiskAssessmentStats, error) {
	if len(sppgIDs) == 0 {
		return []RiskAssessmentStats{}, nil
	}

	type statsRow struct {
		SPPGID       uint
		SPPGNama     string
		TotalAudits  int
		AverageScore *float64
	}

	var rows []statsRow
	err := s.db.
		Table("risk_assessment_forms").
		Select(`
			risk_assessment_forms.sppg_id,
			sppgs.nama as sppg_nama,
			COUNT(*) as total_audits,
			AVG(risk_assessment_forms.overall_risk_score) as average_score
		`).
		Joins("JOIN sppgs ON sppgs.id = risk_assessment_forms.sppg_id").
		Where("risk_assessment_forms.sppg_id IN ? AND risk_assessment_forms.status = ?", sppgIDs, "submitted").
		Group("risk_assessment_forms.sppg_id, sppgs.nama").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	// Build result with latest risk level per SPPG
	results := make([]RiskAssessmentStats, 0, len(rows))
	for _, row := range rows {
		stat := RiskAssessmentStats{
			SPPGID:       row.SPPGID,
			SPPGNama:     row.SPPGNama,
			TotalAudits:  row.TotalAudits,
			AverageScore: row.AverageScore,
		}

		// Get latest risk level (most recent submitted form)
		var latestForm models.RiskAssessmentForm
		if err := s.db.
			Where("sppg_id = ? AND status = ?", row.SPPGID, "submitted").
			Order("submitted_at DESC").
			First(&latestForm).Error; err == nil {
			stat.LatestLevel = latestForm.RiskLevel
		}

		results = append(results, stat)
	}

	return results, nil
}

// determineRiskLevel maps a score to a risk level string.
// "rendah" for [4.0, 5.0], "sedang" for [2.5, 3.9], "tinggi" for [1.0, 2.4]
func (s *RiskAssessmentService) determineRiskLevel(score float64) string {
	if score >= 4.0 {
		return "rendah"
	}
	if score >= 2.5 {
		return "sedang"
	}
	return "tinggi"
}

// SubmitForm validates, calculates scores, applies penalties, and submits a risk assessment form.
func (s *RiskAssessmentService) SubmitForm(formID, yayasanID uint) (*models.RiskAssessmentForm, error) {
	// Load form with items and snapshot
	var form models.RiskAssessmentForm
	query := s.db.
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Order("sop_category_id ASC, id ASC")
		}).
		Preload("Snapshot").
		Where("id = ?", formID)

	if yayasanID > 0 {
		query = query.Where("yayasan_id = ?", yayasanID)
	}

	if err := query.First(&form).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRAFormNotFound
		}
		return nil, err
	}

	// Validate form is still draft
	if form.Status != "draft" {
		return nil, ErrRAFormAlreadySubmitted
	}

	// Validate all items have non-nil ComplianceScore
	var unscoredItems []string
	for _, item := range form.Items {
		if item.ComplianceScore == nil {
			unscoredItems = append(unscoredItems, item.ItemNama)
		}
	}
	if len(unscoredItems) > 0 {
		return nil, fmt.Errorf("%w: %s", ErrRAIncompleteScores, fmt.Sprintf("item belum dinilai: %s", joinStrings(unscoredItems)))
	}

	// Calculate category scores
	categoryScores := s.calculateCategoryScores(form.Items)

	// Calculate overall SOP compliance score (before penalty)
	sopScore := s.calculateOverallScore(form.Items)

	// Apply operational penalty from snapshot
	finalScore := s.applyOperationalPenalty(sopScore, form.Snapshot)

	// Determine risk level
	riskLevel := s.determineRiskLevel(finalScore)

	// Set form fields
	now := time.Now()
	form.Status = "submitted"
	form.SubmittedAt = &now
	form.OverallRiskScore = &finalScore
	form.RiskLevel = &riskLevel

	// Save everything in a transaction
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Update form status, scores, and timestamp
		if err := tx.Model(&form).Updates(map[string]interface{}{
			"status":             form.Status,
			"submitted_at":       form.SubmittedAt,
			"overall_risk_score": form.OverallRiskScore,
			"risk_level":         form.RiskLevel,
		}).Error; err != nil {
			return err
		}

		// Save category scores
		for i := range categoryScores {
			categoryScores[i].FormID = form.ID
			if err := tx.Create(&categoryScores[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Reload form with all associations
	return s.GetForm(form.ID, yayasanID)
}

// calculateCategoryScores groups items by SOPCategoryID, computes the average
// compliance score per category, and determines the risk level for each category.
func (s *RiskAssessmentService) calculateCategoryScores(items []models.RiskAssessmentItem) []models.RiskAssessmentCategoryScore {
	type categoryAccum struct {
		categoryNama string
		totalScore   float64
		count        int
	}

	categoryMap := make(map[uint]*categoryAccum)

	for _, item := range items {
		if item.ComplianceScore == nil {
			continue
		}
		accum, exists := categoryMap[item.SOPCategoryID]
		if !exists {
			accum = &categoryAccum{categoryNama: item.CategoryNama}
			categoryMap[item.SOPCategoryID] = accum
		}
		accum.totalScore += float64(*item.ComplianceScore)
		accum.count++
	}

	scores := make([]models.RiskAssessmentCategoryScore, 0, len(categoryMap))
	for catID, accum := range categoryMap {
		avg := roundScore(accum.totalScore / float64(accum.count))
		scores = append(scores, models.RiskAssessmentCategoryScore{
			SOPCategoryID: catID,
			CategoryNama:  accum.categoryNama,
			AverageScore:  avg,
			RiskLevel:     s.determineRiskLevel(avg),
			ItemCount:     accum.count,
		})
	}

	return scores
}

// calculateOverallScore computes the arithmetic mean of all compliance scores.
func (s *RiskAssessmentService) calculateOverallScore(items []models.RiskAssessmentItem) float64 {
	var total float64
	var count int

	for _, item := range items {
		if item.ComplianceScore == nil {
			continue
		}
		total += float64(*item.ComplianceScore)
		count++
	}

	if count == 0 {
		return 0
	}

	return roundScore(total / float64(count))
}

// applyOperationalPenalty applies a -0.5 penalty per trigger from the operational snapshot.
// Triggers: review rating < 3.0, budget absorption < 50%, on-time delivery < 70%.
// Maximum total penalty is -1.5. The result is clamped to [1.0, 5.0].
func (s *RiskAssessmentService) applyOperationalPenalty(sopScore float64, snapshot *models.SPPGOperationalSnapshot) float64 {
	if snapshot == nil {
		return sopScore
	}

	var penalty float64

	if snapshot.AverageOverallRating < 3.0 {
		penalty += 0.5
	}
	if snapshot.BudgetAbsorptionRate < 50.0 {
		penalty += 0.5
	}
	if snapshot.OnTimeDeliveryRate < 70.0 {
		penalty += 0.5
	}

	// Cap penalty at 1.5
	if penalty > 1.5 {
		penalty = 1.5
	}

	finalScore := sopScore - penalty

	// Clamp to [1.0, 5.0]
	if finalScore < 1.0 {
		finalScore = 1.0
	}
	if finalScore > 5.0 {
		finalScore = 5.0
	}

	return roundScore(finalScore)
}

// joinStrings joins a slice of strings with ", " separator.
func joinStrings(strs []string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += ", "
		}
		result += s
	}
	return result
}

// roundScore rounds a float64 to 2 decimal places.
func roundScore(val float64) float64 {
	return math.Round(val*100) / 100
}

// --- SOP Category CRUD ---

// DeleteForm deletes a draft risk assessment form and all its related data.
func (s *RiskAssessmentService) DeleteForm(formID, yayasanID uint) error {
	var form models.RiskAssessmentForm
	query := s.db.Where("id = ?", formID)
	if yayasanID > 0 {
		query = query.Where("yayasan_id = ?", yayasanID)
	}
	if err := query.First(&form).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRAFormNotFound
		}
		return err
	}
	if form.Status != "draft" {
		return ErrRAFormNotDraft
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		tx.Where("form_id = ?", formID).Delete(&models.SPPGOperationalSnapshot{})
		tx.Where("form_id = ?", formID).Delete(&models.RiskAssessmentCategoryScore{})
		tx.Where("form_id = ?", formID).Delete(&models.RiskAssessmentItem{})
		return tx.Delete(&form).Error
	})
}

// GetSPPGsByYayasan returns SPPGs belonging to the given Yayasan.
// If yayasanID is 0 (superadmin), returns all active SPPGs.
func (s *RiskAssessmentService) GetSPPGsByYayasan(yayasanID uint) ([]models.SPPG, error) {
	var sppgs []models.SPPG
	query := s.db.Where("is_active = ?", true)
	if yayasanID > 0 {
		query = query.Where("yayasan_id = ?", yayasanID)
	}
	if err := query.Order("nama ASC").Find(&sppgs).Error; err != nil {
		return nil, err
	}
	return sppgs, nil
}

// GetSOPCategories returns all SOP categories ordered by Urutan.
func (s *RiskAssessmentService) GetSOPCategories() ([]models.SOPCategory, error) {
	var categories []models.SOPCategory
	if err := s.db.Order("urutan ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// CreateSOPCategory creates a new SOP category with auto-assigned Urutan.
func (s *RiskAssessmentService) CreateSOPCategory(input CreateSOPCategoryInput) (*models.SOPCategory, error) {
	// Auto-assign urutan: max existing + 1
	var maxUrutan int
	s.db.Model(&models.SOPCategory{}).Select("COALESCE(MAX(urutan), 0)").Scan(&maxUrutan)

	category := models.SOPCategory{
		Nama:      input.Nama,
		Deskripsi: input.Deskripsi,
		Urutan:    maxUrutan + 1,
		IsActive:  true,
	}

	if err := s.db.Create(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// UpdateSOPCategory updates an existing SOP category's nama, deskripsi, and/or urutan.
func (s *RiskAssessmentService) UpdateSOPCategory(id uint, input UpdateSOPCategoryInput) (*models.SOPCategory, error) {
	var category models.SOPCategory
	if err := s.db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSOPCategoryNotFound
		}
		return nil, err
	}

	updates := map[string]interface{}{}
	if input.Nama != "" {
		updates["nama"] = input.Nama
	}
	if input.Deskripsi != "" {
		updates["deskripsi"] = input.Deskripsi
	}
	if input.Urutan != nil {
		updates["urutan"] = *input.Urutan
	}

	if len(updates) > 0 {
		if err := s.db.Model(&category).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	// Reload
	if err := s.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// --- SOP Checklist Item CRUD ---

// GetSOPChecklistItems returns checklist items, optionally filtered by category, ordered by Urutan.
func (s *RiskAssessmentService) GetSOPChecklistItems(categoryID *uint) ([]models.SOPChecklistItem, error) {
	var items []models.SOPChecklistItem
	query := s.db.Preload("SOPCategory").Order("sop_category_id ASC, urutan ASC")

	if categoryID != nil && *categoryID > 0 {
		query = query.Where("sop_category_id = ?", *categoryID)
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// CreateSOPChecklistItem creates a new checklist item with auto-assigned Urutan (max existing + 1 in category).
func (s *RiskAssessmentService) CreateSOPChecklistItem(input CreateSOPChecklistItemInput) (*models.SOPChecklistItem, error) {
	// Verify category exists
	var category models.SOPCategory
	if err := s.db.First(&category, input.SOPCategoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSOPCategoryNotFound
		}
		return nil, err
	}

	// Auto-assign urutan: max existing in this category + 1
	var maxUrutan int
	s.db.Model(&models.SOPChecklistItem{}).
		Where("sop_category_id = ?", input.SOPCategoryID).
		Select("COALESCE(MAX(urutan), 0)").
		Scan(&maxUrutan)

	item := models.SOPChecklistItem{
		SOPCategoryID: input.SOPCategoryID,
		Nama:          input.Nama,
		Deskripsi:     input.Deskripsi,
		Urutan:        maxUrutan + 1,
		IsActive:      true,
	}

	if err := s.db.Create(&item).Error; err != nil {
		return nil, err
	}

	// Reload with category
	if err := s.db.Preload("SOPCategory").First(&item, item.ID).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// UpdateSOPChecklistItem updates a checklist item without affecting existing forms.
func (s *RiskAssessmentService) UpdateSOPChecklistItem(id uint, input UpdateSOPChecklistItemInput) (*models.SOPChecklistItem, error) {
	var item models.SOPChecklistItem
	if err := s.db.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSOPChecklistItemNotFound
		}
		return nil, err
	}

	updates := map[string]interface{}{}
	if input.Nama != "" {
		updates["nama"] = input.Nama
	}
	if input.Deskripsi != "" {
		updates["deskripsi"] = input.Deskripsi
	}
	if input.Urutan != nil {
		updates["urutan"] = *input.Urutan
	}

	if len(updates) > 0 {
		if err := s.db.Model(&item).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	// Reload with category
	if err := s.db.Preload("SOPCategory").First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// SetSOPChecklistItemStatus activates or deactivates a checklist item.
func (s *RiskAssessmentService) SetSOPChecklistItemStatus(id uint, isActive bool) (*models.SOPChecklistItem, error) {
	var item models.SOPChecklistItem
	if err := s.db.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSOPChecklistItemNotFound
		}
		return nil, err
	}

	if err := s.db.Model(&item).Update("is_active", isActive).Error; err != nil {
		return nil, err
	}

	// Reload with category
	if err := s.db.Preload("SOPCategory").First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
