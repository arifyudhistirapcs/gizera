package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	fb "github.com/erp-sppg/backend/internal/firebase"
)

// FirebaseSyncService handles real-time data synchronization with Firebase
type FirebaseSyncService struct {
	client *db.Client
}

// NewFirebaseSyncService creates a new Firebase sync service
func NewFirebaseSyncService(firebaseApp *firebase.App) (*FirebaseSyncService, error) {
	ctx := context.Background()
	client, err := firebaseApp.Database(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal menginisialisasi Firebase Database client: %w", err)
	}

	return &FirebaseSyncService{
		client: client,
	}, nil
}

// PushUpdate pushes data to a Firebase path
func (s *FirebaseSyncService) PushUpdate(ctx context.Context, path string, data interface{}) error {
	ref := s.client.NewRef(path)
	if err := ref.Set(ctx, data); err != nil {
		return fmt.Errorf("gagal mengirim update ke Firebase path %s: %w", path, err)
	}
	return nil
}

// PushUpdateWithTimestamp pushes data with an updated_at timestamp
func (s *FirebaseSyncService) PushUpdateWithTimestamp(ctx context.Context, path string, data interface{}) error {
	// Add timestamp to data
	dataMap := make(map[string]interface{})
	
	// Convert data to map
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("gagal mengkonversi data: %w", err)
	}
	
	if err := json.Unmarshal(jsonData, &dataMap); err != nil {
		return fmt.Errorf("gagal mengkonversi data ke map: %w", err)
	}
	
	// Add timestamp
	dataMap["updated_at"] = time.Now().Unix()
	
	return s.PushUpdate(ctx, path, dataMap)
}

// DeletePath removes data from a Firebase path
func (s *FirebaseSyncService) DeletePath(ctx context.Context, path string) error {
	ref := s.client.NewRef(path)
	if err := ref.Delete(ctx); err != nil {
		return fmt.Errorf("gagal menghapus data dari Firebase path %s: %w", path, err)
	}
	return nil
}

// GetData retrieves data from a Firebase path
func (s *FirebaseSyncService) GetData(ctx context.Context, path string, result interface{}) error {
	ref := s.client.NewRef(path)
	if err := ref.Get(ctx, result); err != nil {
		return fmt.Errorf("gagal mengambil data dari Firebase path %s: %w", path, err)
	}
	return nil
}

// UpdateField updates a specific field in a Firebase path
func (s *FirebaseSyncService) UpdateField(ctx context.Context, path string, updates map[string]interface{}) error {
	ref := s.client.NewRef(path)
	if err := ref.Update(ctx, updates); err != nil {
		return fmt.Errorf("gagal mengupdate field di Firebase path %s: %w", path, err)
	}
	return nil
}

// PushKDSCookingUpdate pushes cooking status update to Firebase (tenant-aware)
func (s *FirebaseSyncService) PushKDSCookingUpdate(ctx context.Context, sppgID uint, date string, recipeID uint, data interface{}) error {
	path := fb.KDSCookingRecipePath(sppgID, date, recipeID)
	return s.PushUpdateWithTimestamp(ctx, path, data)
}

// PushKDSPackingUpdate pushes packing status update to Firebase (tenant-aware)
func (s *FirebaseSyncService) PushKDSPackingUpdate(ctx context.Context, sppgID uint, date string, schoolID uint, data interface{}) error {
	path := fb.KDSPackingSchoolPath(sppgID, date, schoolID)
	return s.PushUpdateWithTimestamp(ctx, path, data)
}

// PushDashboardUpdate pushes dashboard data to Firebase
func (s *FirebaseSyncService) PushDashboardUpdate(ctx context.Context, dashboardType string, data interface{}) error {
	path := fmt.Sprintf("/dashboard/%s", dashboardType)
	return s.PushUpdateWithTimestamp(ctx, path, data)
}

// PushDashboardKepalaSSPGUpdate pushes Kepala SPPG dashboard data to Firebase (tenant-aware)
func (s *FirebaseSyncService) PushDashboardKepalaSSPGUpdate(ctx context.Context, sppgID uint, data interface{}) error {
	path := fb.DashboardKepalaSSPGPath(sppgID)
	return s.PushUpdateWithTimestamp(ctx, path, data)
}

// PushDashboardKepalaYayasanUpdate pushes Kepala Yayasan dashboard data to Firebase
func (s *FirebaseSyncService) PushDashboardKepalaYayasanUpdate(ctx context.Context, yayasanID uint, data interface{}) error {
	path := fb.DashboardKepalaYayasanPath(yayasanID)
	return s.PushUpdateWithTimestamp(ctx, path, data)
}

// PushDashboardBGNUpdate pushes BGN dashboard data to Firebase
func (s *FirebaseSyncService) PushDashboardBGNUpdate(ctx context.Context, data interface{}) error {
	path := fb.DashboardBGNPath()
	return s.PushUpdateWithTimestamp(ctx, path, data)
}

// PushInventoryUpdate pushes inventory update to Firebase
func (s *FirebaseSyncService) PushInventoryUpdate(ctx context.Context, ingredientID uint, data interface{}) error {
	path := fmt.Sprintf("/inventory/%d", ingredientID)
	return s.PushUpdateWithTimestamp(ctx, path, data)
}

// PushDeliveryUpdate pushes delivery status update to Firebase
func (s *FirebaseSyncService) PushDeliveryUpdate(ctx context.Context, taskID uint, data interface{}) error {
	path := fmt.Sprintf("/delivery/%d", taskID)
	return s.PushUpdateWithTimestamp(ctx, path, data)
}

// PushMonitoringUpdate pushes monitoring delivery data to Firebase (tenant-aware)
func (s *FirebaseSyncService) PushMonitoringUpdate(ctx context.Context, sppgID uint, date string, recordID uint, data interface{}) error {
	path := fb.MonitoringDeliveryRecordPath(sppgID, date, recordID)
	return s.PushUpdateWithTimestamp(ctx, path, data)
}

// PushCleaningUpdate pushes cleaning record data to Firebase (tenant-aware)
func (s *FirebaseSyncService) PushCleaningUpdate(ctx context.Context, sppgID uint, cleaningID uint, data interface{}) error {
	path := fb.CleaningRecordPath(sppgID, cleaningID)
	return s.PushUpdateWithTimestamp(ctx, path, data)
}

// HandleConflict resolves conflicts by using server data (server wins strategy)
func (s *FirebaseSyncService) HandleConflict(ctx context.Context, path string, serverData interface{}) error {
	// Server data always wins in conflict resolution
	return s.PushUpdateWithTimestamp(ctx, path, serverData)
}

// ClearKDSData clears all KDS-related data from Firebase
// Note: This clears the top-level paths which includes all tenant data
func (s *FirebaseSyncService) ClearKDSData(ctx context.Context) error {
	// List of top-level paths to clear (clears all tenants)
	paths := []string{
		"/kds/cooking",
		"/kds/packing",
		"/delivery_tasks",
		"/delivery_records",
		"/monitoring",
		"/activity_tracker",
		"/cleaning",
	}

	for _, path := range paths {
		if err := s.DeletePath(ctx, path); err != nil {
			// Log error but continue with other paths
			fmt.Printf("Peringatan: gagal menghapus Firebase path %s: %v\n", path, err)
		}
	}

	return nil
}

// ClearKDSDataForSPPG clears KDS-related data for a specific SPPG tenant
func (s *FirebaseSyncService) ClearKDSDataForSPPG(ctx context.Context, sppgID uint) error {
	paths := []string{
		fmt.Sprintf("/kds/cooking/%d", sppgID),
		fmt.Sprintf("/kds/packing/%d", sppgID),
		fmt.Sprintf("/monitoring/%d", sppgID),
		fmt.Sprintf("/cleaning/%d", sppgID),
		fb.DashboardKepalaSSPGPath(sppgID),
	}

	for _, path := range paths {
		if err := s.DeletePath(ctx, path); err != nil {
			fmt.Printf("Peringatan: gagal menghapus Firebase path %s: %v\n", path, err)
		}
	}

	return nil
}
