package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"github.com/erp-sppg/backend/internal/config"
	"google.golang.org/api/option"
)

// Context keys for tenant-aware Firebase operations
type contextKey string

const (
	// SPPGIDKey is the context key for SPPG ID
	SPPGIDKey contextKey = "firebase_sppg_id"
	// YayasanIDKey is the context key for Yayasan ID
	YayasanIDKey contextKey = "firebase_yayasan_id"
)

// WithSPPGID returns a new context with the given sppg_id set.
func WithSPPGID(ctx context.Context, sppgID uint) context.Context {
	return context.WithValue(ctx, SPPGIDKey, sppgID)
}

// GetSPPGID extracts sppg_id from context. Returns 0 if not set.
func GetSPPGID(ctx context.Context) uint {
	if v, ok := ctx.Value(SPPGIDKey).(uint); ok {
		return v
	}
	return 0
}

// WithYayasanID returns a new context with the given yayasan_id set.
func WithYayasanID(ctx context.Context, yayasanID uint) context.Context {
	return context.WithValue(ctx, YayasanIDKey, yayasanID)
}

// GetYayasanID extracts yayasan_id from context. Returns 0 if not set.
func GetYayasanID(ctx context.Context) uint {
	if v, ok := ctx.Value(YayasanIDKey).(uint); ok {
		return v
	}
	return 0
}

func Initialize(cfg *config.Config) (*firebase.App, error) {
	ctx := context.Background()

	opt := option.WithCredentialsFile(cfg.FirebaseCredentialsPath)
	
	firebaseConfig := &firebase.Config{
		DatabaseURL:   cfg.FirebaseDatabaseURL,
		StorageBucket: cfg.StorageBucket,
	}

	app, err := firebase.NewApp(ctx, firebaseConfig, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase: %w", err)
	}

	return app, nil
}

// --- Tenant-Aware Firebase Path Helpers ---

// KDSCookingPath returns the tenant-aware Firebase path for KDS cooking data.
// Format: /kds/cooking/{sppg_id}/{date}/{recipe_id}
func KDSCookingPath(sppgID uint, date string) string {
	return fmt.Sprintf("/kds/cooking/%d/%s", sppgID, date)
}

// KDSCookingRecipePath returns the tenant-aware Firebase path for a specific recipe in cooking.
// Format: /kds/cooking/{sppg_id}/{date}/{recipe_id}
func KDSCookingRecipePath(sppgID uint, date string, recipeID uint) string {
	return fmt.Sprintf("/kds/cooking/%d/%s/%d", sppgID, date, recipeID)
}

// KDSPackingPath returns the tenant-aware Firebase path for KDS packing data.
// Format: /kds/packing/{sppg_id}/{date}
func KDSPackingPath(sppgID uint, date string) string {
	return fmt.Sprintf("/kds/packing/%d/%s", sppgID, date)
}

// KDSPackingSchoolPath returns the tenant-aware Firebase path for a specific school in packing.
// Format: /kds/packing/{sppg_id}/{date}/{school_id}
func KDSPackingSchoolPath(sppgID uint, date string, schoolID uint) string {
	return fmt.Sprintf("/kds/packing/%d/%s/%d", sppgID, date, schoolID)
}

// DashboardKepalaSSPGPath returns the tenant-aware Firebase path for Kepala SPPG dashboard.
// Format: /dashboard/kepala_sppg/{sppg_id}
func DashboardKepalaSSPGPath(sppgID uint) string {
	return fmt.Sprintf("/dashboard/kepala_sppg/%d", sppgID)
}

// DashboardKepalaYayasanPath returns the Firebase path for Kepala Yayasan dashboard.
// Format: /dashboard/kepala_yayasan/{yayasan_id}
func DashboardKepalaYayasanPath(yayasanID uint) string {
	return fmt.Sprintf("/dashboard/kepala_yayasan/%d", yayasanID)
}

// DashboardBGNPath returns the Firebase path for BGN dashboard.
// Format: /dashboard/bgn
func DashboardBGNPath() string {
	return "/dashboard/bgn"
}

// MonitoringDeliveriesPath returns the tenant-aware Firebase path for monitoring deliveries.
// Format: /monitoring/{sppg_id}/{date}
func MonitoringDeliveriesPath(sppgID uint, date string) string {
	return fmt.Sprintf("/monitoring/%d/%s", sppgID, date)
}

// MonitoringDeliveryRecordPath returns the tenant-aware Firebase path for a specific delivery record.
// Format: /monitoring/{sppg_id}/{date}/record_{id}
func MonitoringDeliveryRecordPath(sppgID uint, date string, recordID uint) string {
	return fmt.Sprintf("/monitoring/%d/%s/record_%d", sppgID, date, recordID)
}

// CleaningPendingPath returns the tenant-aware Firebase path for pending cleaning records.
// Format: /cleaning/{sppg_id}/pending
func CleaningPendingPath(sppgID uint) string {
	return fmt.Sprintf("/cleaning/%d/pending", sppgID)
}

// CleaningRecordPath returns the tenant-aware Firebase path for a specific cleaning record.
// Format: /cleaning/{sppg_id}/pending/cleaning_{id}
func CleaningRecordPath(sppgID uint, cleaningID uint) string {
	return fmt.Sprintf("/cleaning/%d/pending/cleaning_%d", sppgID, cleaningID)
}

// InjectSPPGIDFromGin extracts sppg_id from Gin context and returns a new
// context.Context with the sppg_id set. This bridges the Gin context to
// the standard Go context used by services.
func InjectSPPGIDFromGin(c interface{ Get(string) (interface{}, bool) }, ctx context.Context) context.Context {
	if sppgIDVal, exists := c.Get("sppg_id"); exists {
		if sppgID, ok := sppgIDVal.(uint); ok {
			ctx = WithSPPGID(ctx, sppgID)
		}
	}
	if yayasanIDVal, exists := c.Get("yayasan_id"); exists {
		if yayasanID, ok := yayasanIDVal.(uint); ok {
			ctx = WithYayasanID(ctx, yayasanID)
		}
	}
	return ctx
}
