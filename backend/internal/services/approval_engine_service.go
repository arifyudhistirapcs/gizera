package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/erp-sppg/backend/internal/models"
	"gorm.io/gorm"
)

var (
	ErrRABNotFound      = errors.New("RAB tidak ditemukan")
	ErrRABInvalidStatus = errors.New("status RAB tidak valid untuk operasi ini")
	ErrRABNotEditable   = errors.New("RAB tidak dapat diedit")
	ErrRABNotApprovable = errors.New("RAB tidak dapat disetujui dari status saat ini")
)

// ApprovalEngineService handles multi-level RAB approval workflow
type ApprovalEngineService struct {
	db    *gorm.DB
	notif *NotificationService
}

// NewApprovalEngineService creates a new approval engine service
func NewApprovalEngineService(db *gorm.DB, notif *NotificationService) *ApprovalEngineService {
	return &ApprovalEngineService{
		db:    db,
		notif: notif,
	}
}

// ApproveByKepalaSPPG transitions RAB from draft → approved_sppg
func (s *ApprovalEngineService) ApproveByKepalaSPPG(rabID uint, userID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var rab models.RAB
		if err := tx.First(&rab, rabID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrRABNotFound
			}
			return fmt.Errorf("gagal mengambil RAB: %w", err)
		}

		if rab.Status != "draft" {
			return ErrRABNotApprovable
		}

		now := time.Now()
		if err := tx.Model(&models.RAB{}).Where("id = ?", rabID).Updates(map[string]interface{}{
			"status":           "approved_sppg",
			"approved_by_sppg": userID,
			"approved_at_sppg": now,
		}).Error; err != nil {
			return fmt.Errorf("gagal mengupdate status RAB: %w", err)
		}

		// Create audit trail
		if err := s.createAuditTrailTx(tx, rabID, userID, "draft", "approved_sppg", "Disetujui oleh Kepala SPPG"); err != nil {
			return err
		}

		// Send notification to kepala_yayasan (outside tx, graceful)
		if s.notif != nil && rab.YayasanID != nil {
			go s.sendApprovalNotificationToYayasan(*rab.YayasanID, rab.RABNumber, "approved_sppg")
		}

		return nil
	})
}

// ApproveByKepalaYayasan transitions RAB from approved_sppg → approved_yayasan
func (s *ApprovalEngineService) ApproveByKepalaYayasan(rabID uint, userID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var rab models.RAB
		if err := tx.First(&rab, rabID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrRABNotFound
			}
			return fmt.Errorf("gagal mengambil RAB: %w", err)
		}

		if rab.Status != "approved_sppg" {
			return ErrRABNotApprovable
		}

		now := time.Now()
		if err := tx.Model(&models.RAB{}).Where("id = ?", rabID).Updates(map[string]interface{}{
			"status":              "approved_yayasan",
			"approved_by_yayasan": userID,
			"approved_at_yayasan": now,
		}).Error; err != nil {
			return fmt.Errorf("gagal mengupdate status RAB: %w", err)
		}

		// Create audit trail
		if err := s.createAuditTrailTx(tx, rabID, userID, "approved_sppg", "approved_yayasan", "Disetujui oleh Kepala Yayasan"); err != nil {
			return err
		}

		// Send notification to kepala_sppg
		if s.notif != nil && rab.SPPGID != nil {
			go s.sendApprovalNotificationToSPPG(*rab.SPPGID, rab.RABNumber, "approved_yayasan")
		}

		return nil
	})
}

// RejectByKepalaYayasan transitions RAB from approved_sppg → revision_requested
func (s *ApprovalEngineService) RejectByKepalaYayasan(rabID uint, userID uint, notes string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var rab models.RAB
		if err := tx.First(&rab, rabID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrRABNotFound
			}
			return fmt.Errorf("gagal mengambil RAB: %w", err)
		}

		if rab.Status != "approved_sppg" {
			return ErrRABInvalidStatus
		}

		if err := tx.Model(&models.RAB{}).Where("id = ?", rabID).Updates(map[string]interface{}{
			"status":         "revision_requested",
			"revision_notes": notes,
		}).Error; err != nil {
			return fmt.Errorf("gagal mengupdate status RAB: %w", err)
		}

		// Create audit trail
		if err := s.createAuditTrailTx(tx, rabID, userID, "approved_sppg", "revision_requested", notes); err != nil {
			return err
		}

		// Reset menu plan back to draft so SPPG can revise it
		if err := tx.Model(&models.MenuPlan{}).Where("id = ?", rab.MenuPlanID).Updates(map[string]interface{}{
			"status":      "draft",
			"approved_by": nil,
			"approved_at": nil,
		}).Error; err != nil {
			return fmt.Errorf("gagal mereset menu plan ke draft: %w", err)
		}

		// Send notification to kepala_sppg
		if s.notif != nil && rab.SPPGID != nil {
			go s.sendApprovalNotificationToSPPG(*rab.SPPGID, rab.RABNumber, "revision_requested")
		}

		return nil
	})
}

// ResubmitRAB transitions RAB from revision_requested → draft
func (s *ApprovalEngineService) ResubmitRAB(rabID uint, userID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var rab models.RAB
		if err := tx.First(&rab, rabID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrRABNotFound
			}
			return fmt.Errorf("gagal mengambil RAB: %w", err)
		}

		if rab.Status != "revision_requested" {
			return ErrRABInvalidStatus
		}

		if err := tx.Model(&models.RAB{}).Where("id = ?", rabID).Updates(map[string]interface{}{
			"status":         "draft",
			"revision_notes": "",
		}).Error; err != nil {
			return fmt.Errorf("gagal mengupdate status RAB: %w", err)
		}

		// Create audit trail
		if err := s.createAuditTrailTx(tx, rabID, userID, "revision_requested", "draft", "RAB disubmit ulang setelah revisi"); err != nil {
			return err
		}

		return nil
	})
}

// createAuditTrail creates an AuditTrail record for RAB status change (non-transactional)
func (s *ApprovalEngineService) createAuditTrail(rabID uint, userID uint, fromStatus, toStatus, notes string) error {
	return s.createAuditTrailTx(s.db, rabID, userID, fromStatus, toStatus, notes)
}

// createAuditTrailTx creates an AuditTrail record within a transaction
func (s *ApprovalEngineService) createAuditTrailTx(tx *gorm.DB, rabID uint, userID uint, fromStatus, toStatus, notes string) error {
	audit := &models.AuditTrail{
		UserID:    userID,
		Timestamp: time.Now(),
		Action:    "status_change",
		Entity:    "rab",
		EntityID:  fmt.Sprintf("%d", rabID),
		OldValue:  fromStatus,
		NewValue:  toStatus,
		Level:     "info",
	}

	if notes != "" {
		audit.NewValue = fmt.Sprintf("%s (catatan: %s)", toStatus, notes)
	}

	if err := tx.Create(audit).Error; err != nil {
		return fmt.Errorf("gagal membuat audit trail: %w", err)
	}

	return nil
}

// sendApprovalNotificationToYayasan sends notification to kepala_yayasan users
func (s *ApprovalEngineService) sendApprovalNotificationToYayasan(yayasanID uint, rabNumber string, newStatus string) {
	var users []models.User
	s.db.Where("yayasan_id = ? AND role = ? AND is_active = ?", yayasanID, "kepala_yayasan", true).Find(&users)

	title := "RAB Menunggu Persetujuan"
	message := fmt.Sprintf("RAB %s telah disetujui oleh Kepala SPPG dan menunggu persetujuan Anda.", rabNumber)

	for _, user := range users {
		notification := &models.Notification{
			UserID:  user.ID,
			Type:    "rab_approval",
			Title:   title,
			Message: message,
			Link:    "/rab",
		}
		if err := s.notif.CreateNotification(context.Background(), notification); err != nil {
			fmt.Printf("Peringatan: gagal mengirim notifikasi ke kepala_yayasan %d: %v\n", user.ID, err)
		}
	}
}

// sendApprovalNotificationToSPPG sends notification to kepala_sppg users
func (s *ApprovalEngineService) sendApprovalNotificationToSPPG(sppgID uint, rabNumber string, newStatus string) {
	var users []models.User
	s.db.Where("sppg_id = ? AND role = ? AND is_active = ?", sppgID, "kepala_sppg", true).Find(&users)

	var title, message string
	switch newStatus {
	case "approved_yayasan":
		title = "RAB Disetujui Yayasan"
		message = fmt.Sprintf("RAB %s telah disetujui oleh Kepala Yayasan. PO dapat dibuat.", rabNumber)
	case "revision_requested":
		title = "RAB Perlu Revisi"
		message = fmt.Sprintf("RAB %s ditolak oleh Kepala Yayasan. Silakan revisi dan submit ulang.", rabNumber)
	default:
		title = "Update Status RAB"
		message = fmt.Sprintf("Status RAB %s berubah menjadi %s.", rabNumber, newStatus)
	}

	for _, user := range users {
		notification := &models.Notification{
			UserID:  user.ID,
			Type:    "rab_approval",
			Title:   title,
			Message: message,
			Link:    "/rab",
		}
		if err := s.notif.CreateNotification(context.Background(), notification); err != nil {
			fmt.Printf("Peringatan: gagal mengirim notifikasi ke kepala_sppg %d: %v\n", user.ID, err)
		}
	}
}
