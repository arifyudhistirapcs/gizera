package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/erp-sppg/backend/internal/models"
	"gorm.io/gorm"
)

// Organization service errors
var (
	ErrYayasanNotFound      = errors.New("yayasan tidak ditemukan")
	ErrYayasanInactive      = errors.New("yayasan tidak aktif")
	ErrDuplicateYayasanKode = errors.New("kode yayasan sudah terdaftar")
	ErrDuplicateYayasanEmail = errors.New("email yayasan sudah terdaftar")
	ErrDuplicateYayasanNPWP = errors.New("NPWP yayasan sudah terdaftar")
	ErrSPPGNotFound         = errors.New("SPPG tidak ditemukan")
	ErrDuplicateSPPGKode    = errors.New("kode SPPG sudah terdaftar")
	ErrDuplicateSPPGEmail   = errors.New("email SPPG sudah terdaftar")
	ErrInvalidYayasanID     = errors.New("yayasan_id wajib diisi dan harus valid")
	ErrSPPGUnderInactiveYayasan = errors.New("tidak dapat membuat SPPG di bawah yayasan yang nonaktif")
)

// YayasanService handles Yayasan CRUD operations
type YayasanService struct {
	db           *gorm.DB
	auditService *AuditTrailService
}

// NewYayasanService creates a new YayasanService
func NewYayasanService(db *gorm.DB, auditService *AuditTrailService) *YayasanService {
	return &YayasanService{
		db:           db,
		auditService: auditService,
	}
}

// generateYayasanKode generates the next unique code in format "YYS-XXXX"
func (s *YayasanService) generateYayasanKode() (string, error) {
	var lastYayasan models.Yayasan
	err := s.db.Order("kode DESC").First(&lastYayasan).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "YYS-0001", nil
		}
		return "", err
	}

	// Parse the number from the last code
	parts := strings.Split(lastYayasan.Kode, "-")
	if len(parts) != 2 {
		return "YYS-0001", nil
	}

	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return "YYS-0001", nil
	}

	return fmt.Sprintf("YYS-%04d", num+1), nil
}

// validateYayasanUniqueness checks that kode, email, and NPWP are unique
func (s *YayasanService) validateYayasanUniqueness(yayasan *models.Yayasan, excludeID uint) error {
	var existing models.Yayasan

	// Check kode uniqueness
	q := s.db.Where("kode = ?", yayasan.Kode)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	if err := q.First(&existing).Error; err == nil {
		return ErrDuplicateYayasanKode
	}

	// Check email uniqueness (only if email is provided)
	if yayasan.Email != "" {
		q = s.db.Where("email = ?", yayasan.Email)
		if excludeID > 0 {
			q = q.Where("id != ?", excludeID)
		}
		if err := q.First(&existing).Error; err == nil {
			return ErrDuplicateYayasanEmail
		}
	}

	// Check NPWP uniqueness (only if NPWP is provided)
	if yayasan.NPWP != "" {
		q = s.db.Where("npwp = ?", yayasan.NPWP)
		if excludeID > 0 {
			q = q.Where("id != ?", excludeID)
		}
		if err := q.First(&existing).Error; err == nil {
			return ErrDuplicateYayasanNPWP
		}
	}

	return nil
}

// Create creates a new Yayasan with auto-generated kode
func (s *YayasanService) Create(yayasan *models.Yayasan, userID uint) error {
	// Auto-generate kode
	kode, err := s.generateYayasanKode()
	if err != nil {
		return fmt.Errorf("gagal generate kode yayasan: %w", err)
	}
	yayasan.Kode = kode
	yayasan.IsActive = true

	// Validate uniqueness
	if err := s.validateYayasanUniqueness(yayasan, 0); err != nil {
		return err
	}

	if err := s.db.Create(yayasan).Error; err != nil {
		return err
	}

	// Record audit trail
	if s.auditService != nil {
		s.auditService.RecordAction(userID, "create", "yayasan", fmt.Sprintf("%d", yayasan.ID), nil, yayasan, "")
	}

	return nil
}

// GetAll returns all Yayasan records
func (s *YayasanService) GetAll() ([]models.Yayasan, error) {
	var yayasans []models.Yayasan
	if err := s.db.Preload("SPPGs").Order("id ASC").Find(&yayasans).Error; err != nil {
		return nil, err
	}
	return yayasans, nil
}

// GetByID returns a Yayasan by ID
func (s *YayasanService) GetByID(id uint) (*models.Yayasan, error) {
	var yayasan models.Yayasan
	if err := s.db.Preload("SPPGs").First(&yayasan, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrYayasanNotFound
		}
		return nil, err
	}
	return &yayasan, nil
}

// Update updates a Yayasan's data
func (s *YayasanService) Update(id uint, updates map[string]interface{}, userID uint) (*models.Yayasan, error) {
	// Get existing record
	existing, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	oldValue := *existing

	// Validate uniqueness if email or NPWP is being updated
	tempYayasan := &models.Yayasan{
		Kode:  existing.Kode,
		Email: existing.Email,
		NPWP:  existing.NPWP,
	}
	if email, ok := updates["email"].(string); ok {
		tempYayasan.Email = email
	}
	if npwp, ok := updates["npwp"].(string); ok {
		tempYayasan.NPWP = npwp
	}
	if err := s.validateYayasanUniqueness(tempYayasan, id); err != nil {
		return nil, err
	}

	// Remove fields that should not be updated directly
	delete(updates, "id")
	delete(updates, "kode")
	delete(updates, "created_at")

	updates["updated_at"] = time.Now()

	if err := s.db.Model(&models.Yayasan{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Reload updated record
	updated, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Record audit trail
	if s.auditService != nil {
		s.auditService.RecordAction(userID, "update", "yayasan", fmt.Sprintf("%d", id), oldValue, updated, "")
	}

	return updated, nil
}

// SetStatus activates or deactivates a Yayasan
func (s *YayasanService) SetStatus(id uint, isActive bool, userID uint) (*models.Yayasan, error) {
	existing, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	oldValue := *existing

	if err := s.db.Model(&models.Yayasan{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_active":  isActive,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return nil, err
	}

	updated, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	action := "activate"
	if !isActive {
		action = "deactivate"
	}

	// Record audit trail
	if s.auditService != nil {
		s.auditService.RecordAction(userID, action, "yayasan", fmt.Sprintf("%d", id), oldValue, updated, "")
	}

	return updated, nil
}

// SPPGService handles SPPG CRUD operations
type SPPGService struct {
	db             *gorm.DB
	auditService   *AuditTrailService
	yayasanService *YayasanService
}

// NewSPPGService creates a new SPPGService
func NewSPPGService(db *gorm.DB, auditService *AuditTrailService, yayasanService *YayasanService) *SPPGService {
	return &SPPGService{
		db:             db,
		auditService:   auditService,
		yayasanService: yayasanService,
	}
}

// generateSPPGKode generates the next unique code in format "SPPG-XXXX"
func (s *SPPGService) generateSPPGKode() (string, error) {
	var lastSPPG models.SPPG
	err := s.db.Order("kode DESC").First(&lastSPPG).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "SPPG-0001", nil
		}
		return "", err
	}

	parts := strings.Split(lastSPPG.Kode, "-")
	if len(parts) != 2 {
		return "SPPG-0001", nil
	}

	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return "SPPG-0001", nil
	}

	return fmt.Sprintf("SPPG-%04d", num+1), nil
}

// validateSPPGUniqueness checks that kode and email are unique
func (s *SPPGService) validateSPPGUniqueness(sppg *models.SPPG, excludeID uint) error {
	var existing models.SPPG

	// Check kode uniqueness
	q := s.db.Where("kode = ?", sppg.Kode)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	if err := q.First(&existing).Error; err == nil {
		return ErrDuplicateSPPGKode
	}

	// Check email uniqueness (only if email is provided)
	if sppg.Email != "" {
		q = s.db.Where("email = ?", sppg.Email)
		if excludeID > 0 {
			q = q.Where("id != ?", excludeID)
		}
		if err := q.First(&existing).Error; err == nil {
			return ErrDuplicateSPPGEmail
		}
	}

	return nil
}

// Create creates a new SPPG with auto-generated kode
func (s *SPPGService) Create(sppg *models.SPPG, userID uint) error {
	// Validate YayasanID is provided
	if sppg.YayasanID == 0 {
		return ErrInvalidYayasanID
	}

	// Validate Yayasan exists and is active
	yayasan, err := s.yayasanService.GetByID(sppg.YayasanID)
	if err != nil {
		return ErrInvalidYayasanID
	}
	if !yayasan.IsActive {
		return ErrSPPGUnderInactiveYayasan
	}

	// Auto-generate kode
	kode, err := s.generateSPPGKode()
	if err != nil {
		return fmt.Errorf("gagal generate kode SPPG: %w", err)
	}
	sppg.Kode = kode
	sppg.IsActive = true

	// Validate uniqueness
	if err := s.validateSPPGUniqueness(sppg, 0); err != nil {
		return err
	}

	if err := s.db.Create(sppg).Error; err != nil {
		return err
	}

	// Record audit trail
	if s.auditService != nil {
		s.auditService.RecordAction(userID, "create", "sppg", fmt.Sprintf("%d", sppg.ID), nil, sppg, "")
	}

	return nil
}

// GetAll returns all SPPG records
func (s *SPPGService) GetAll() ([]models.SPPG, error) {
	var sppgs []models.SPPG
	if err := s.db.Preload("Yayasan").Order("id ASC").Find(&sppgs).Error; err != nil {
		return nil, err
	}
	return sppgs, nil
}

// GetByID returns a SPPG by ID
func (s *SPPGService) GetByID(id uint) (*models.SPPG, error) {
	var sppg models.SPPG
	if err := s.db.Preload("Yayasan").First(&sppg, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSPPGNotFound
		}
		return nil, err
	}
	return &sppg, nil
}

// Update updates a SPPG's data
func (s *SPPGService) Update(id uint, updates map[string]interface{}, userID uint) (*models.SPPG, error) {
	existing, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	oldValue := *existing

	// Validate email uniqueness if being updated
	tempSPPG := &models.SPPG{
		Kode:  existing.Kode,
		Email: existing.Email,
	}
	if email, ok := updates["email"].(string); ok {
		tempSPPG.Email = email
	}
	if err := s.validateSPPGUniqueness(tempSPPG, id); err != nil {
		return nil, err
	}

	// If yayasan_id is being changed, validate the new Yayasan
	if yayasanID, ok := updates["yayasan_id"]; ok {
		var yID uint
		switch v := yayasanID.(type) {
		case uint:
			yID = v
		case float64:
			yID = uint(v)
		}
		if yID > 0 {
			yayasan, err := s.yayasanService.GetByID(yID)
			if err != nil {
				return nil, ErrInvalidYayasanID
			}
			if !yayasan.IsActive {
				return nil, ErrSPPGUnderInactiveYayasan
			}
		}
	}

	// Remove fields that should not be updated directly
	delete(updates, "id")
	delete(updates, "kode")
	delete(updates, "created_at")

	updates["updated_at"] = time.Now()

	if err := s.db.Model(&models.SPPG{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}

	updated, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Record audit trail
	if s.auditService != nil {
		s.auditService.RecordAction(userID, "update", "sppg", fmt.Sprintf("%d", id), oldValue, updated, "")
	}

	return updated, nil
}

// SetStatus activates or deactivates a SPPG
func (s *SPPGService) SetStatus(id uint, isActive bool, userID uint) (*models.SPPG, error) {
	existing, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	oldValue := *existing

	if err := s.db.Model(&models.SPPG{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_active":  isActive,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return nil, err
	}

	updated, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	action := "activate"
	if !isActive {
		action = "deactivate"
	}

	if s.auditService != nil {
		s.auditService.RecordAction(userID, action, "sppg", fmt.Sprintf("%d", id), oldValue, updated, "")
	}

	return updated, nil
}

// Transfer moves a SPPG from one Yayasan to another
func (s *SPPGService) Transfer(id uint, newYayasanID uint, userID uint) (*models.SPPG, error) {
	existing, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Validate new Yayasan exists and is active
	newYayasan, err := s.yayasanService.GetByID(newYayasanID)
	if err != nil {
		return nil, ErrInvalidYayasanID
	}
	if !newYayasan.IsActive {
		return nil, ErrYayasanInactive
	}

	oldValue := *existing

	if err := s.db.Model(&models.SPPG{}).Where("id = ?", id).Updates(map[string]interface{}{
		"yayasan_id": newYayasanID,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return nil, err
	}

	updated, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Record audit trail for transfer
	if s.auditService != nil {
		s.auditService.RecordAction(userID, "transfer", "sppg", fmt.Sprintf("%d", id), oldValue, updated, "")
	}

	return updated, nil
}
