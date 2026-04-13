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
	ErrInvoiceNotFound       = errors.New("invoice tidak ditemukan")
	ErrInvoiceAlreadyPaid    = errors.New("invoice sudah dibayar")
	ErrInvoiceAmountMismatch = errors.New("amount invoice tidak sesuai dengan PO")
	ErrGRNNotCompleted       = errors.New("GRN belum selesai untuk PO ini")
	ErrDuplicateInvoice      = errors.New("invoice sudah dibuat untuk PO ini")
)

// InvoiceService handles invoice and payment business logic
type InvoiceService struct {
	db       *gorm.DB
	cashFlow *CashFlowService
	notif    *NotificationService
}

// NewInvoiceService creates a new invoice service
func NewInvoiceService(db *gorm.DB, cashFlow *CashFlowService, notif *NotificationService) *InvoiceService {
	return &InvoiceService{
		db:       db,
		cashFlow: cashFlow,
		notif:    notif,
	}
}

// WithDB returns a new service instance with the given DB
func (s *InvoiceService) WithDB(db *gorm.DB) *InvoiceService {
	return &InvoiceService{
		db:       db,
		cashFlow: s.cashFlow,
		notif:    s.notif,
	}
}

// CreateInvoice creates a new invoice after GRN is complete
func (s *InvoiceService) CreateInvoice(invoice *models.Invoice) error {
	// Validate PO exists and has status "received"
	var po models.PurchaseOrder
	if err := s.db.First(&po, invoice.POID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPONotFound
		}
		return fmt.Errorf("gagal mengambil PO: %w", err)
	}
	if po.Status != "received" {
		return ErrGRNNotCompleted
	}

	// Validate PO has a GRN
	var grnCount int64
	s.db.Model(&models.GoodsReceipt{}).Where("po_id = ?", invoice.POID).Count(&grnCount)
	if grnCount == 0 {
		return ErrGRNNotCompleted
	}

	// Validate amount equals PO total_amount
	if invoice.Amount != po.TotalAmount {
		return ErrInvoiceAmountMismatch
	}

	// Check no duplicate invoice for same PO
	var existingCount int64
	s.db.Model(&models.Invoice{}).Where("po_id = ?", invoice.POID).Count(&existingCount)
	if existingCount > 0 {
		return ErrDuplicateInvoice
	}

	// Generate invoice number
	invoiceNumber, err := s.generateInvoiceNumber()
	if err != nil {
		return fmt.Errorf("gagal generate invoice number: %w", err)
	}

	// Set invoice fields
	invoice.InvoiceNumber = invoiceNumber
	invoice.Status = "pending"
	invoice.SupplierID = po.SupplierID
	if po.YayasanID != nil {
		invoice.YayasanID = *po.YayasanID
	}

	if err := s.db.Create(invoice).Error; err != nil {
		return fmt.Errorf("gagal membuat invoice: %w", err)
	}

	// Send notification to kepala_yayasan (graceful degradation)
	if s.notif != nil && invoice.YayasanID > 0 {
		go s.sendInvoiceNotificationToYayasan(invoice.YayasanID, invoice.InvoiceNumber)
	}

	return nil
}

// ProcessPayment processes payment for an invoice
func (s *InvoiceService) ProcessPayment(invoiceID uint, payment *models.Payment) error {
	// Validate invoice exists and has status "pending"
	var invoice models.Invoice
	if err := s.db.First(&invoice, invoiceID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvoiceNotFound
		}
		return fmt.Errorf("gagal mengambil invoice: %w", err)
	}

	if invoice.Status == "paid" {
		return ErrInvoiceAlreadyPaid
	}

	if invoice.Status != "pending" {
		return ErrInvoiceNotFound
	}

	// No duplicate payment check
	var existingPaymentCount int64
	s.db.Model(&models.Payment{}).Where("invoice_id = ?", invoiceID).Count(&existingPaymentCount)
	if existingPaymentCount > 0 {
		return ErrInvoiceAlreadyPaid
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Create Payment record
		payment.InvoiceID = invoiceID
		if err := tx.Create(payment).Error; err != nil {
			return fmt.Errorf("gagal membuat payment: %w", err)
		}

		// Update invoice status to "paid"
		if err := tx.Model(&models.Invoice{}).Where("id = ?", invoiceID).Updates(map[string]interface{}{
			"status":     "paid",
			"updated_at": time.Now(),
		}).Error; err != nil {
			return fmt.Errorf("gagal mengupdate status invoice: %w", err)
		}

		// Create CashFlowEntry
		if s.cashFlow != nil {
			cashFlowEntry := &models.CashFlowEntry{
				TransactionID: fmt.Sprintf("PAY-%s", invoice.InvoiceNumber),
				Date:          payment.PaymentDate,
				Category:      "pengadaan",
				Type:          "expense",
				Amount:        payment.Amount,
				Description:   fmt.Sprintf("Pembayaran invoice %s ke supplier", invoice.InvoiceNumber),
				Reference:     invoice.InvoiceNumber,
				YayasanID:     &invoice.YayasanID,
				CreatedBy:     payment.PaidBy,
			}
			if err := s.cashFlow.CreateCashFlowEntryWithTx(tx, cashFlowEntry); err != nil {
				return fmt.Errorf("gagal membuat cash flow entry: %w", err)
			}
		}

		// Send notification to supplier (graceful degradation)
		if s.notif != nil {
			go s.sendPaymentNotificationToSupplier(invoice.SupplierID, invoice.InvoiceNumber, payment.Amount)
		}

		return nil
	})
}

// GetInvoiceByID retrieves an invoice by ID with preloaded relations
func (s *InvoiceService) GetInvoiceByID(id uint) (*models.Invoice, error) {
	var invoice models.Invoice
	err := s.db.Preload("PurchaseOrder.Supplier").
		Preload("Payment").
		First(&invoice, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvoiceNotFound
		}
		return nil, fmt.Errorf("gagal mengambil invoice: %w", err)
	}

	return &invoice, nil
}

// GetInvoicesBySupplier retrieves invoices filtered by supplier_id
func (s *InvoiceService) GetInvoicesBySupplier(supplierID uint) ([]models.Invoice, error) {
	var invoices []models.Invoice
	err := s.db.Preload("PurchaseOrder.Supplier").
		Preload("Payment").
		Where("supplier_id = ?", supplierID).
		Order("created_at DESC").
		Find(&invoices).Error
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil invoices supplier: %w", err)
	}
	return invoices, nil
}

// GetInvoicesByYayasan retrieves invoices filtered by yayasan_id
func (s *InvoiceService) GetInvoicesByYayasan(yayasanID uint) ([]models.Invoice, error) {
	var invoices []models.Invoice
	err := s.db.Preload("PurchaseOrder.Supplier").
		Preload("Payment").
		Where("yayasan_id = ?", yayasanID).
		Order("created_at DESC").
		Find(&invoices).Error
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil invoices yayasan: %w", err)
	}
	return invoices, nil
}

// GetSupplierPaymentHistory retrieves payments received by a supplier
func (s *InvoiceService) GetSupplierPaymentHistory(supplierID uint) ([]models.Payment, error) {
	var payments []models.Payment
	err := s.db.Joins("JOIN invoices ON invoices.id = payments.invoice_id").
		Preload("Invoice").
		Preload("Payer").
		Where("invoices.supplier_id = ?", supplierID).
		Order("payments.payment_date DESC").
		Find(&payments).Error
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil riwayat pembayaran supplier: %w", err)
	}
	return payments, nil
}

// generateInvoiceNumber generates a unique invoice number with format INV-YYYYMMDD-XXXX
func (s *InvoiceService) generateInvoiceNumber() (string, error) {
	now := time.Now()
	datePrefix := now.Format("20060102")

	var count int64
	s.db.Model(&models.Invoice{}).
		Where("invoice_number LIKE ?", fmt.Sprintf("INV-%s-%%", datePrefix)).
		Count(&count)

	invoiceNumber := fmt.Sprintf("INV-%s-%04d", datePrefix, count+1)

	// Race condition protection
	var existing models.Invoice
	err := s.db.Where("invoice_number = ?", invoiceNumber).First(&existing).Error
	if err == nil {
		invoiceNumber = fmt.Sprintf("INV-%s-%04d", datePrefix, count+2)
	}

	return invoiceNumber, nil
}

// UpdatePaymentProof updates the proof_url on a payment record
func (s *InvoiceService) UpdatePaymentProof(paymentID uint, proofURL string) error {
	if err := s.db.Model(&models.Payment{}).Where("id = ?", paymentID).
		Update("proof_url", proofURL).Error; err != nil {
		return fmt.Errorf("gagal mengupdate bukti pembayaran: %w", err)
	}
	return nil
}

// sendInvoiceNotificationToYayasan sends notification to kepala_yayasan about new invoice
func (s *InvoiceService) sendInvoiceNotificationToYayasan(yayasanID uint, invoiceNumber string) {
	var users []models.User
	s.db.Where("yayasan_id = ? AND role = ? AND is_active = ?", yayasanID, "kepala_yayasan", true).Find(&users)

	for _, user := range users {
		notification := &models.Notification{
			UserID:  user.ID,
			Type:    "invoice_created",
			Title:   "Invoice Baru dari Supplier",
			Message: fmt.Sprintf("Invoice %s telah dibuat dan menunggu pembayaran.", invoiceNumber),
			Link:    "/invoices",
		}
		if err := s.notif.CreateNotification(context.Background(), notification); err != nil {
			fmt.Printf("Peringatan: gagal mengirim notifikasi invoice ke kepala_yayasan %d: %v\n", user.ID, err)
		}
	}
}

// sendPaymentNotificationToSupplier sends notification to supplier about payment
func (s *InvoiceService) sendPaymentNotificationToSupplier(supplierID uint, invoiceNumber string, amount float64) {
	var users []models.User
	s.db.Where("supplier_id = ? AND role = ? AND is_active = ?", supplierID, "supplier", true).Find(&users)

	for _, user := range users {
		notification := &models.Notification{
			UserID:  user.ID,
			Type:    "payment_received",
			Title:   "Pembayaran Diterima",
			Message: fmt.Sprintf("Pembayaran untuk invoice %s sebesar Rp %.0f telah dilakukan.", invoiceNumber, amount),
			Link:    "/invoices",
		}
		if err := s.notif.CreateNotification(context.Background(), notification); err != nil {
			fmt.Printf("Peringatan: gagal mengirim notifikasi pembayaran ke supplier %d: %v\n", user.ID, err)
		}
	}
}
