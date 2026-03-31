package services

import (
	"encoding/json"
	"time"

	"github.com/erp-sppg/backend/internal/models"
	"gorm.io/gorm"
)

// AuditTrailService handles audit trail operations
type AuditTrailService struct {
	db *gorm.DB
}

// NewAuditTrailService creates a new audit trail service
func NewAuditTrailService(db *gorm.DB) *AuditTrailService {
	return &AuditTrailService{
		db: db,
	}
}

// RecordAction records a user action in the audit trail
func (s *AuditTrailService) RecordAction(userID uint, action, entity, entityID string, oldValue, newValue interface{}, ipAddress string) error {
	// Convert values to JSON strings
	oldJSON, err := json.Marshal(oldValue)
	if err != nil {
		oldJSON = []byte("{}")
	}

	newJSON, err := json.Marshal(newValue)
	if err != nil {
		newJSON = []byte("{}")
	}

	auditEntry := models.AuditTrail{
		UserID:    userID,
		Timestamp: time.Now(),
		Action:    action,
		Entity:    entity,
		EntityID:  entityID,
		OldValue:  string(oldJSON),
		NewValue:  string(newJSON),
		IPAddress: ipAddress,
		Level:     "info",
	}

	return s.db.Create(&auditEntry).Error
}

// RecordActionWithTenant records a user action in the audit trail with tenant context.
func (s *AuditTrailService) RecordActionWithTenant(userID uint, action, entity, entityID string, oldValue, newValue interface{}, ipAddress string, sppgID, yayasanID *uint) error {
	oldJSON, err := json.Marshal(oldValue)
	if err != nil {
		oldJSON = []byte("{}")
	}

	newJSON, err := json.Marshal(newValue)
	if err != nil {
		newJSON = []byte("{}")
	}

	auditEntry := models.AuditTrail{
		UserID:    userID,
		Timestamp: time.Now(),
		Action:    action,
		Entity:    entity,
		EntityID:  entityID,
		OldValue:  string(oldJSON),
		NewValue:  string(newJSON),
		IPAddress: ipAddress,
		SPPGID:    sppgID,
		YayasanID: yayasanID,
		Level:     "info",
	}

	return s.db.Create(&auditEntry).Error
}

// RecordLogin records a user login action
func (s *AuditTrailService) RecordLogin(userID uint, ipAddress string) error {
	return s.RecordAction(userID, "login", "user", "", nil, map[string]interface{}{
		"timestamp": time.Now(),
		"ip":        ipAddress,
	}, ipAddress)
}

// RecordLogout records a user logout action
func (s *AuditTrailService) RecordLogout(userID uint, ipAddress string) error {
	return s.RecordAction(userID, "logout", "user", "", nil, map[string]interface{}{
		"timestamp": time.Now(),
		"ip":        ipAddress,
	}, ipAddress)
}

// GetAuditTrail retrieves audit trail entries with filters
func (s *AuditTrailService) GetAuditTrail(filters map[string]interface{}, limit, offset int) ([]models.AuditTrail, int64, error) {
	var entries []models.AuditTrail
	var total int64

	query := s.db.Model(&models.AuditTrail{}).Preload("User")

	// Apply filters
	if userID, ok := filters["user_id"].(uint); ok && userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	if action, ok := filters["action"].(string); ok && action != "" {
		query = query.Where("action = ?", action)
	}

	if entity, ok := filters["entity"].(string); ok && entity != "" {
		query = query.Where("entity = ?", entity)
	}

	if startDate, ok := filters["start_date"].(time.Time); ok {
		query = query.Where("timestamp >= ?", startDate)
	}

	if endDate, ok := filters["end_date"].(time.Time); ok {
		query = query.Where("timestamp <= ?", endDate)
	}

	// Tenant scoping filters
	if sppgID, ok := filters["sppg_id"].(uint); ok && sppgID > 0 {
		query = query.Where("sppg_id = ?", sppgID)
	}

	if yayasanID, ok := filters["yayasan_id"].(uint); ok && yayasanID > 0 {
		query = query.Where("yayasan_id = ?", yayasanID)
	}

	// Filter by sppg_ids (for kepala_yayasan scope)
	if sppgIDs, ok := filters["sppg_ids"].([]uint); ok && len(sppgIDs) > 0 {
		query = query.Where("sppg_id IN ? OR sppg_id IS NULL", sppgIDs)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.Order("timestamp DESC").Limit(limit).Offset(offset).Find(&entries).Error; err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}
