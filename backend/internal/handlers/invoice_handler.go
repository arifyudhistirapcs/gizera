package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/erp-sppg/backend/internal/middleware"
	"github.com/erp-sppg/backend/internal/models"
	"github.com/erp-sppg/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// InvoiceHandler handles invoice and payment endpoints
type InvoiceHandler struct {
	invoiceService *services.InvoiceService
}

// NewInvoiceHandler creates a new invoice handler
func NewInvoiceHandler(invoiceService *services.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{
		invoiceService: invoiceService,
	}
}

// GetInvoices handles GET /invoices — list invoices scoped by role
func (h *InvoiceHandler) GetInvoices(c *gin.Context) {
	role, _ := c.Get("user_role")
	roleStr, _ := role.(string)

	var invoices []models.Invoice
	var err error

	switch roleStr {
	case "supplier":
		supplierID, ok := middleware.GetSupplierID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"error_code": "UNAUTHORIZED",
				"message":    "Supplier ID tidak ditemukan",
			})
			return
		}
		invoices, err = h.invoiceService.GetInvoicesBySupplier(supplierID)

	case "kepala_yayasan":
		yayasanIDVal, ok := c.Get("yayasan_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success":    false,
				"error_code": "UNAUTHORIZED",
				"message":    "Yayasan ID tidak ditemukan",
			})
			return
		}
		yayasanID := yayasanIDVal.(uint)
		invoices, err = h.invoiceService.GetInvoicesByYayasan(yayasanID)

	default:
		c.JSON(http.StatusForbidden, gin.H{
			"success":    false,
			"error_code": "FORBIDDEN",
			"message":    "Anda tidak memiliki izin untuk mengakses resource ini",
		})
		return
	}

	if err != nil {
		log.Printf("GetInvoices error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    invoices,
	})
}

// CreateInvoiceRequest represents the request body for creating an invoice
type CreateInvoiceRequest struct {
	POID    uint    `json:"po_id" binding:"required"`
	Amount  float64 `json:"amount" binding:"required,gt=0"`
	DueDate string  `json:"due_date" binding:"required"`
}

// CreateInvoice handles POST /invoices — supplier creates invoice
func (h *InvoiceHandler) CreateInvoice(c *gin.Context) {
	supplierID, ok := middleware.GetSupplierID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success":    false,
			"error_code": "UNAUTHORIZED",
			"message":    "Supplier ID tidak ditemukan",
		})
		return
	}

	var req CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_DATE",
			"message":    "Format tanggal tidak valid (gunakan YYYY-MM-DD)",
		})
		return
	}

	invoice := &models.Invoice{
		POID:       req.POID,
		SupplierID: supplierID,
		Amount:     req.Amount,
		DueDate:    dueDate,
	}

	if err := h.invoiceService.CreateInvoice(invoice); err != nil {
		if err == services.ErrPONotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "PO_NOT_FOUND",
				"message":    "Purchase order tidak ditemukan",
			})
			return
		}
		if err == services.ErrGRNNotCompleted {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "GRN_NOT_COMPLETED",
				"message":    "GRN belum selesai untuk PO ini",
			})
			return
		}
		if err == services.ErrInvoiceAmountMismatch {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INVOICE_AMOUNT_MISMATCH",
				"message":    "Amount invoice tidak sesuai dengan total PO",
			})
			return
		}
		if err == services.ErrDuplicateInvoice {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "DUPLICATE_INVOICE",
				"message":    "Invoice sudah dibuat untuk PO ini",
			})
			return
		}
		log.Printf("CreateInvoice error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Invoice berhasil dibuat",
		"data":    invoice,
	})
}

// GetInvoiceDetail handles GET /invoices/:id
func (h *InvoiceHandler) GetInvoiceDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	invoice, err := h.invoiceService.GetInvoiceByID(uint(id))
	if err != nil {
		if err == services.ErrInvoiceNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "INVOICE_NOT_FOUND",
				"message":    "Invoice tidak ditemukan",
			})
			return
		}
		log.Printf("GetInvoiceDetail error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    invoice,
	})
}

// PayInvoiceRequest represents the request body for paying an invoice
type PayInvoiceRequest struct {
	PaymentDate   string  `json:"payment_date" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	PaymentMethod string  `json:"payment_method" binding:"required"`
}

// PayInvoice handles POST /invoices/:id/pay — kepala_yayasan pays invoice
func (h *InvoiceHandler) PayInvoice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	var req PayInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	paymentDate, err := time.Parse("2006-01-02", req.PaymentDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_DATE",
			"message":    "Format tanggal tidak valid (gunakan YYYY-MM-DD)",
		})
		return
	}

	userID, _ := c.Get("user_id")

	payment := &models.Payment{
		PaymentDate:   paymentDate,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		PaidBy:        userID.(uint),
	}

	if err := h.invoiceService.ProcessPayment(uint(id), payment); err != nil {
		if err == services.ErrInvoiceNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "INVOICE_NOT_FOUND",
				"message":    "Invoice tidak ditemukan",
			})
			return
		}
		if err == services.ErrInvoiceAlreadyPaid {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INVOICE_ALREADY_PAID",
				"message":    "Invoice sudah dibayar",
			})
			return
		}
		log.Printf("PayInvoice error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pembayaran berhasil dicatat",
	})
}

// UploadPaymentProof handles POST /invoices/:id/upload-proof
func (h *InvoiceHandler) UploadPaymentProof(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_ID",
			"message":    "ID tidak valid",
		})
		return
	}

	// Verify invoice exists
	invoice, err := h.invoiceService.GetInvoiceByID(uint(id))
	if err != nil {
		if err == services.ErrInvoiceNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success":    false,
				"error_code": "INVOICE_NOT_FOUND",
				"message":    "Invoice tidak ditemukan",
			})
			return
		}
		log.Printf("UploadPaymentProof error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan pada server",
		})
		return
	}

	// Get file from form
	file, err := c.FormFile("proof")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "NO_FILE",
			"message":    "File bukti pembayaran tidak ditemukan",
		})
		return
	}

	// Save file
	filename := fmt.Sprintf("payment_proof_%d_%d%s", id, time.Now().Unix(), filepath.Ext(file.Filename))
	proofDir := filepath.Join("uploads", "payment-proofs")

	if err := os.MkdirAll(proofDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "UPLOAD_ERROR",
			"message":    "Gagal membuat direktori upload",
		})
		return
	}

	savePath := filepath.Join(proofDir, filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "UPLOAD_ERROR",
			"message":    "Gagal menyimpan file",
		})
		return
	}

	proofURL := fmt.Sprintf("/uploads/payment-proofs/%s", filename)

	// Update payment proof_url if payment exists
	if invoice.Payment != nil {
		h.invoiceService.UpdatePaymentProof(invoice.Payment.ID, proofURL)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "Bukti pembayaran berhasil diunggah",
		"proof_url": proofURL,
	})
}
