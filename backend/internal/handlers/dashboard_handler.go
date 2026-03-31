package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	firebase "firebase.google.com/go/v4"
	fb "github.com/erp-sppg/backend/internal/firebase"
	"github.com/erp-sppg/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DashboardHandler handles dashboard endpoints
type DashboardHandler struct {
	dashboardService *services.DashboardService
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(db *gorm.DB, firebaseApp *firebase.App) (*DashboardHandler, error) {
	dashboardService, err := services.NewDashboardService(db, firebaseApp)
	if err != nil {
		return nil, err
	}

	return &DashboardHandler{
		dashboardService: dashboardService,
	}, nil
}

// GetKepalaSSPGDashboard retrieves operational dashboard for Kepala SPPG
// @Summary Get Kepala SPPG Dashboard
// @Description Retrieves operational dashboard with production status, delivery status, critical stock, and today's KPIs
// @Tags Dashboard
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/dashboard/kepala-sppg [get]
func (h *DashboardHandler) GetKepalaSSPGDashboard(c *gin.Context) {
	ctx := fb.InjectSPPGIDFromGin(c, context.Background())

	dashboard, err := h.dashboardService.GetKepalaSSPGDashboard(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan saat mengambil data dashboard",
			"details":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"dashboard": dashboard,
	})
}

// GetKepalaYayasanDashboard retrieves strategic dashboard for Kepala Yayasan
// @Summary Get Kepala Yayasan Dashboard
// @Description Retrieves strategic dashboard with budget absorption, nutrition distribution, supplier performance, and monthly trends
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD)" default(first day of current month)
// @Param end_date query string false "End date (YYYY-MM-DD)" default(today)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/dashboard/kepala-yayasan [get]
func (h *DashboardHandler) GetKepalaYayasanDashboard(c *gin.Context) {
	ctx := fb.InjectSPPGIDFromGin(c, context.Background())

	// Parse date range (default to current month)
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := now

	if startStr := c.Query("start_date"); startStr != "" {
		var err error
		startDate, err = time.Parse("2006-01-02", startStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INVALID_DATE",
				"message":    "Format start_date tidak valid (gunakan YYYY-MM-DD)",
			})
			return
		}
	}

	if endStr := c.Query("end_date"); endStr != "" {
		var err error
		endDate, err = time.Parse("2006-01-02", endStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INVALID_DATE",
				"message":    "Format end_date tidak valid (gunakan YYYY-MM-DD)",
			})
			return
		}
	}

	// Validate date range
	if endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_DATE_RANGE",
			"message":    "Tanggal akhir harus setelah tanggal awal",
		})
		return
	}

	dashboard, err := h.dashboardService.GetKepalaYayasanDashboard(ctx, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan saat mengambil data dashboard",
			"details":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"dashboard":  dashboard,
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
	})
}

// GetKPIs retrieves key performance indicators
// @Summary Get KPIs
// @Description Retrieves key performance indicators for today
// @Tags Dashboard
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/dashboard/kpi [get]
func (h *DashboardHandler) GetKPIs(c *gin.Context) {
	ctx := fb.InjectSPPGIDFromGin(c, context.Background())

	dashboard, err := h.dashboardService.GetKepalaSSPGDashboard(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan saat mengambil data KPI",
			"details":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"kpis":    dashboard.TodayKPIs,
	})
}

// SyncDashboardToFirebase syncs dashboard data to Firebase
// @Summary Sync Dashboard to Firebase
// @Description Manually triggers sync of dashboard data to Firebase for real-time updates
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param type query string true "Dashboard type (kepala_sppg or kepala_yayasan)"
// @Param start_date query string false "Start date for Kepala Yayasan dashboard (YYYY-MM-DD)"
// @Param end_date query string false "End date for Kepala Yayasan dashboard (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/dashboard/sync [post]
func (h *DashboardHandler) SyncDashboardToFirebase(c *gin.Context) {
	ctx := fb.InjectSPPGIDFromGin(c, context.Background())

	dashboardType := c.Query("type")
	if dashboardType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "MISSING_TYPE",
			"message":    "Parameter type diperlukan (kepala_sppg atau kepala_yayasan)",
		})
		return
	}

	var err error
	switch dashboardType {
	case "kepala_sppg":
		// Get sppg_id from context for tenant-aware Firebase path
		var sppgID uint
		if sppgIDVal, exists := c.Get("sppg_id"); exists {
			if id, ok := sppgIDVal.(uint); ok {
				sppgID = id
			}
		}
		err = h.dashboardService.SyncKepalaSSPGDashboardToFirebase(ctx, sppgID)

	case "kepala_yayasan":
		// Parse date range (default to current month)
		now := time.Now()
		startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		endDate := now

		if startStr := c.Query("start_date"); startStr != "" {
			startDate, err = time.Parse("2006-01-02", startStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success":    false,
					"error_code": "INVALID_DATE",
					"message":    "Format start_date tidak valid (gunakan YYYY-MM-DD)",
				})
				return
			}
		}

		if endStr := c.Query("end_date"); endStr != "" {
			endDate, err = time.Parse("2006-01-02", endStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success":    false,
					"error_code": "INVALID_DATE",
					"message":    "Format end_date tidak valid (gunakan YYYY-MM-DD)",
				})
				return
			}
		}

		// Get yayasan_id from context for tenant-aware Firebase path
		var yayasanID uint
		if yayasanIDVal, exists := c.Get("yayasan_id"); exists {
			if id, ok := yayasanIDVal.(uint); ok {
				yayasanID = id
			}
		}

		err = h.dashboardService.SyncKepalaYayasanDashboardToFirebase(ctx, yayasanID, startDate, endDate)

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_TYPE",
			"message":    "Tipe dashboard tidak valid (gunakan kepala_sppg atau kepala_yayasan)",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "SYNC_ERROR",
			"message":    "Gagal melakukan sinkronisasi ke Firebase",
			"details":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Dashboard berhasil disinkronkan ke Firebase",
	})
}

// ExportDashboardRequest represents export dashboard request
type ExportDashboardRequest struct {
	Type      string `json:"type" binding:"required,oneof=kepala_sppg kepala_yayasan"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Format    string `json:"format" binding:"required,oneof=json pdf excel"`
}

// ExportDashboard exports dashboard data
// @Summary Export Dashboard
// @Description Exports dashboard data in specified format (JSON, PDF, or Excel)
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param request body ExportDashboardRequest true "Export request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/dashboard/export [post]
func (h *DashboardHandler) ExportDashboard(c *gin.Context) {
	ctx := fb.InjectSPPGIDFromGin(c, context.Background())

	var req ExportDashboardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "VALIDATION_ERROR",
			"message":    "Data tidak valid",
			"details":    err.Error(),
		})
		return
	}

	// Parse dates for Kepala Yayasan dashboard
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := now

	if req.Type == "kepala_yayasan" {
		if req.StartDate != "" {
			var err error
			startDate, err = time.Parse("2006-01-02", req.StartDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success":    false,
					"error_code": "INVALID_DATE",
					"message":    "Format start_date tidak valid (gunakan YYYY-MM-DD)",
				})
				return
			}
		}

		if req.EndDate != "" {
			var err error
			endDate, err = time.Parse("2006-01-02", req.EndDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"success":    false,
					"error_code": "INVALID_DATE",
					"message":    "Format end_date tidak valid (gunakan YYYY-MM-DD)",
				})
				return
			}
		}
	}

	// Export data
	data, err := h.dashboardService.ExportDashboardData(ctx, req.Type, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "EXPORT_ERROR",
			"message":    "Gagal mengekspor data dashboard",
			"details":    err.Error(),
		})
		return
	}

	// Handle different export formats
	switch req.Format {
	case "json":
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    data,
		})

	case "pdf", "excel":
		// TODO: Implement PDF and Excel export
		// For now, return JSON with a message
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": fmt.Sprintf("Export format %s belum diimplementasikan, mengembalikan JSON", req.Format),
			"data":    data,
		})

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_FORMAT",
			"message":    "Format export tidak valid",
		})
	}
}

// ClearFirebaseKDSData clears all KDS-related data from Firebase
// @Summary Clear Firebase KDS Data
// @Description Clears all KDS, delivery, and monitoring data from Firebase
// @Tags Dashboard
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/dashboard/clear-firebase [post]
func (h *DashboardHandler) ClearFirebaseKDSData(c *gin.Context) {
	ctx := fb.InjectSPPGIDFromGin(c, context.Background())

	if err := h.dashboardService.ClearFirebaseKDSData(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan saat menghapus data Firebase",
			"details":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data KDS di Firebase berhasil dihapus",
	})
}

// =====================================================================
// AggregatedDashboardHandler — multi-tenant aggregated dashboards
// =====================================================================

// AggregatedDashboardHandler handles aggregated dashboard endpoints
// for Kepala Yayasan (cross-SPPG) and Admin BGN (national).
type AggregatedDashboardHandler struct {
	aggregatedService *services.AggregatedDashboardService
}

// NewAggregatedDashboardHandler creates a new AggregatedDashboardHandler.
func NewAggregatedDashboardHandler(aggregatedService *services.AggregatedDashboardService) *AggregatedDashboardHandler {
	return &AggregatedDashboardHandler{
		aggregatedService: aggregatedService,
	}
}

// parseDateRange parses start_date and end_date query params.
// Defaults to current month if not provided.
func parseDateRange(c *gin.Context) (time.Time, time.Time, bool) {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := now

	if s := c.Query("start_date"); s != "" {
		parsed, err := time.Parse("2006-01-02", s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INVALID_DATE",
				"message":    "Format start_date tidak valid (gunakan YYYY-MM-DD)",
			})
			return time.Time{}, time.Time{}, false
		}
		startDate = parsed
	}

	if s := c.Query("end_date"); s != "" {
		parsed, err := time.Parse("2006-01-02", s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":    false,
				"error_code": "INVALID_DATE",
				"message":    "Format end_date tidak valid (gunakan YYYY-MM-DD)",
			})
			return time.Time{}, time.Time{}, false
		}
		endDate = parsed
	}

	if endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_DATE_RANGE",
			"message":    "Tanggal akhir harus setelah tanggal awal",
		})
		return time.Time{}, time.Time{}, false
	}

	return startDate, endDate, true
}

// parseOptionalUintQuery parses an optional uint query parameter.
// Returns nil if the param is absent or empty.
func parseOptionalUintQuery(c *gin.Context, param string) (*uint, bool) {
	s := c.Query(param)
	if s == "" {
		return nil, true
	}
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "INVALID_PARAM",
			"message":    fmt.Sprintf("Parameter %s tidak valid", param),
		})
		return nil, false
	}
	u := uint(v)
	return &u, true
}

// GetAggregatedKepalaYayasanDashboard handles GET /api/v1/dashboard/kepala-yayasan
// Query params: start_date, end_date, sppg_id (optional drill-down), yayasan_id (optional, for superadmin/admin_bgn)
// For kepala_yayasan: yayasan_id is extracted from JWT context.
// For superadmin/admin_bgn: yayasan_id can be passed as query param.
func (h *AggregatedDashboardHandler) GetAggregatedKepalaYayasanDashboard(c *gin.Context) {
	startDate, endDate, ok := parseDateRange(c)
	if !ok {
		return
	}

	sppgID, ok := parseOptionalUintQuery(c, "sppg_id")
	if !ok {
		return
	}

	// Determine yayasan_id: from JWT context first, then from query param (for superadmin/admin_bgn)
	var yayasanID uint
	yayasanIDVal, exists := c.Get("yayasan_id")
	if exists && yayasanIDVal != nil {
		if id, ok := yayasanIDVal.(uint); ok && id > 0 {
			yayasanID = id
		}
	}

	// If not in context (superadmin/admin_bgn), try query param
	if yayasanID == 0 {
		if qp, ok := parseOptionalUintQuery(c, "yayasan_id"); ok && qp != nil {
			yayasanID = *qp
		}
	}

	// If still no yayasan_id, return error
	if yayasanID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "MISSING_YAYASAN",
			"message":    "Yayasan ID diperlukan. Gunakan query parameter ?yayasan_id=",
		})
		return
	}

	dashboard, err := h.aggregatedService.GetKepalaYayasanDashboard(yayasanID, startDate, endDate, sppgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan saat mengambil data dashboard",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"data":       dashboard,
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
	})
}

// GetAdminBGNDashboard handles GET /api/v1/dashboard/admin-bgn
// Query params: start_date, end_date, yayasan_id (optional), sppg_id (optional)
func (h *AggregatedDashboardHandler) GetAdminBGNDashboard(c *gin.Context) {
	startDate, endDate, ok := parseDateRange(c)
	if !ok {
		return
	}

	yayasanID, ok := parseOptionalUintQuery(c, "yayasan_id")
	if !ok {
		return
	}

	sppgID, ok := parseOptionalUintQuery(c, "sppg_id")
	if !ok {
		return
	}

	dashboard, err := h.aggregatedService.GetAdminBGNDashboard(startDate, endDate, yayasanID, sppgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan saat mengambil data dashboard",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"data":       dashboard,
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
	})
}

// ExportKepalaYayasanDashboard handles GET /api/v1/dashboard/kepala-yayasan/export
// Returns the same data as the dashboard endpoint (frontend handles CSV/PDF).
func (h *AggregatedDashboardHandler) ExportKepalaYayasanDashboard(c *gin.Context) {
	startDate, endDate, ok := parseDateRange(c)
	if !ok {
		return
	}

	sppgID, ok := parseOptionalUintQuery(c, "sppg_id")
	if !ok {
		return
	}

	// Determine yayasan_id: from JWT context first, then from query param
	var yayasanID uint
	yayasanIDVal, exists := c.Get("yayasan_id")
	if exists && yayasanIDVal != nil {
		if id, ok := yayasanIDVal.(uint); ok && id > 0 {
			yayasanID = id
		}
	}
	if yayasanID == 0 {
		if qp, ok := parseOptionalUintQuery(c, "yayasan_id"); ok && qp != nil {
			yayasanID = *qp
		}
	}
	if yayasanID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":    false,
			"error_code": "MISSING_YAYASAN",
			"message":    "Yayasan ID diperlukan. Gunakan query parameter ?yayasan_id=",
		})
		return
	}

	dashboard, err := h.aggregatedService.GetKepalaYayasanDashboard(yayasanID, startDate, endDate, sppgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan saat mengekspor data dashboard",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"data":       dashboard,
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
		"export":     true,
	})
}

// ExportAdminBGNDashboard handles GET /api/v1/dashboard/admin-bgn/export
// Returns the same data as the dashboard endpoint (frontend handles CSV/PDF).
func (h *AggregatedDashboardHandler) ExportAdminBGNDashboard(c *gin.Context) {
	startDate, endDate, ok := parseDateRange(c)
	if !ok {
		return
	}

	yayasanID, ok := parseOptionalUintQuery(c, "yayasan_id")
	if !ok {
		return
	}

	sppgID, ok := parseOptionalUintQuery(c, "sppg_id")
	if !ok {
		return
	}

	dashboard, err := h.aggregatedService.GetAdminBGNDashboard(startDate, endDate, yayasanID, sppgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":    false,
			"error_code": "INTERNAL_ERROR",
			"message":    "Terjadi kesalahan saat mengekspor data dashboard",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"data":       dashboard,
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
		"export":     true,
	})
}
