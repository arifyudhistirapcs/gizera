package services

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/erp-sppg/backend/internal/models"
	"gorm.io/gorm"
)

var (
	ErrAttendanceNotFound      = errors.New("data absensi tidak ditemukan")
	ErrInvalidWiFi             = errors.New("anda harus terhubung ke Wi-Fi kantor untuk absen")
	ErrAlreadyCheckedIn        = errors.New("anda sudah melakukan check-in hari ini")
	ErrNotCheckedIn            = errors.New("anda belum melakukan check-in hari ini")
	ErrAlreadyCheckedOut       = errors.New("anda sudah melakukan check-out hari ini")
	ErrInvalidAttendanceData   = errors.New("data absensi tidak valid")
)

// startOfLocalDay returns the start of today in local timezone (midnight local time).
// time.Truncate(24h) truncates to UTC midnight which is wrong for non-UTC timezones.
func startOfLocalDay() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

// AttendanceService handles attendance operations
type AttendanceService struct {
	db              *gorm.DB
	employeeService *EmployeeService
}

// NewAttendanceService creates a new attendance service
func NewAttendanceService(db *gorm.DB, employeeService *EmployeeService) *AttendanceService {
	return &AttendanceService{
		db:              db,
		employeeService: employeeService,
	}
}

// ValidateWiFi validates if the provided SSID and BSSID match authorized networks
func (s *AttendanceService) ValidateWiFi(ssid, bssid string) (bool, error) {
	var wifiConfig models.WiFiConfig
	// Use column names explicitly to avoid GORM naming issues
	result := s.db.Where("ss_id = ? AND bss_id = ? AND is_active = ?", ssid, bssid, true).First(&wifiConfig)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}

	return true, nil
}

// ValidateIP validates if the provided IP address is in authorized range
func (s *AttendanceService) ValidateIP(ipAddress string) (bool, *models.WiFiConfig, error) {
	var configs []models.WiFiConfig
	result := s.db.Where("is_active = ?", true).Find(&configs)
	
	if result.Error != nil {
		return false, nil, result.Error
	}

	fmt.Printf("[ValidateIP] Checking IP: %s against %d active configs\n", ipAddress, len(configs))

	// Development mode: Allow localhost IPs (127.0.0.1, ::1, localhost)
	// In production, remove this or make it configurable
	if ipAddress == "127.0.0.1" || ipAddress == "::1" || ipAddress == "localhost" {
		// Return first active config for localhost testing
		if len(configs) > 0 {
			fmt.Printf("[ValidateIP] LOCALHOST DETECTED - Using config ID: %d, SSID: %s\n", configs[0].ID, configs[0].SSID)
			return true, &configs[0], nil
		}
	}

	for _, config := range configs {
		fmt.Printf("[ValidateIP] Checking config ID: %d, SSID: %s, IP Range: %s, Allowed IPs: %v\n", 
			config.ID, config.SSID, config.IPRange, config.AllowedIPs)
		
		// Check if IP is in the allowed IP range
		if config.IPRange != "" {
			if s.isIPInRange(ipAddress, config.IPRange) {
				fmt.Printf("[ValidateIP] IP MATCH via IP Range - Config ID: %d\n", config.ID)
				return true, &config, nil
			}
		}
		
		// Check if IP is in the specific allowed IPs list
		if len(config.AllowedIPs) > 0 {
			for _, allowedIP := range config.AllowedIPs {
				if ipAddress == allowedIP {
					fmt.Printf("[ValidateIP] IP MATCH via Allowed IPs - Config ID: %d\n", config.ID)
					return true, &config, nil
				}
			}
		}
	}

	fmt.Printf("[ValidateIP] NO MATCH - IP %s not authorized\n", ipAddress)
	return false, nil, nil
}

// isIPInRange checks if an IP address is within a CIDR range
func (s *AttendanceService) isIPInRange(ipAddress, cidr string) bool {
	// Simple IP range check
	// For production, use a proper CIDR library
	// This is a basic implementation
	
	// Parse CIDR (e.g., "192.168.1.0/24")
	parts := strings.Split(cidr, "/")
	if len(parts) != 2 {
		return false
	}
	
	networkIP := parts[0]
	
	// Simple check: if IP starts with the network prefix
	// For /24, check first 3 octets
	// For production, use proper CIDR calculation
	ipParts := strings.Split(ipAddress, ".")
	networkParts := strings.Split(networkIP, ".")
	
	if len(ipParts) != 4 || len(networkParts) != 4 {
		return false
	}
	
	// For /24 network, check first 3 octets
	if parts[1] == "24" {
		return ipParts[0] == networkParts[0] &&
			ipParts[1] == networkParts[1] &&
			ipParts[2] == networkParts[2]
	}
	
	// For /16 network, check first 2 octets
	if parts[1] == "16" {
		return ipParts[0] == networkParts[0] &&
			ipParts[1] == networkParts[1]
	}
	
	return false
}

// CheckIn records employee check-in with Wi-Fi validation
func (s *AttendanceService) CheckIn(employeeID uint, ssid, bssid string) (*models.Attendance, error) {
	// Skip WiFi validation if already validated by GPS (ssid starts with "GPS-")
	if !strings.HasPrefix(ssid, "GPS-") {
		// Validate Wi-Fi only for non-GPS check-ins
		isValid, err := s.ValidateWiFi(ssid, bssid)
		if err != nil {
			return nil, err
		}
		if !isValid {
			return nil, ErrInvalidWiFi
		}
	}

	// Check if employee exists and is active
	employee, err := s.employeeService.GetEmployeeByID(employeeID)
	if err != nil {
		return nil, err
	}
	if !employee.IsActive {
		return nil, errors.New("akun karyawan tidak aktif")
	}

	// Check if already checked in today
	today := startOfLocalDay()
	var existingAttendance models.Attendance
	result := s.db.Where("employee_id = ? AND date >= ? AND date < ?", 
		employeeID, today, today.Add(24*time.Hour)).First(&existingAttendance)
	
	if result.Error == nil {
		return nil, ErrAlreadyCheckedIn
	}

	// Create attendance record
	attendance := &models.Attendance{
		EmployeeID: employeeID,
		Date:       time.Now(),
		CheckIn:    time.Now(),
		SSID:       ssid,
		BSSID:      bssid,
		WorkHours:  0,
	}

	if err := s.db.Create(attendance).Error; err != nil {
		return nil, err
	}

	// Preload employee data
	if err := s.db.Preload("Employee").First(attendance, attendance.ID).Error; err != nil {
		return nil, err
	}

	return attendance, nil
}

// CheckOut records employee check-out and calculates work hours
func (s *AttendanceService) CheckOut(employeeID uint, ssid, bssid string) (*models.Attendance, error) {
	// Check-out doesn't require WiFi validation
	// Skip validation if SSID/BSSID are empty
	if ssid != "" && bssid != "" {
		// Validate Wi-Fi only if provided
		isValid, err := s.ValidateWiFi(ssid, bssid)
		if err != nil {
			return nil, err
		}
		if !isValid {
			return nil, ErrInvalidWiFi
		}
	}

	// Get today's attendance record
	today := startOfLocalDay()
	var attendance models.Attendance
	result := s.db.Where("employee_id = ? AND date >= ? AND date < ?", 
		employeeID, today, today.Add(24*time.Hour)).First(&attendance)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotCheckedIn
		}
		return nil, result.Error
	}

	// Check if already checked out
	if attendance.CheckOut != nil {
		return nil, ErrAlreadyCheckedOut
	}

	// Calculate work hours
	checkOutTime := time.Now()
	workHours := checkOutTime.Sub(attendance.CheckIn).Hours()

	// Update attendance record
	attendance.CheckOut = &checkOutTime
	attendance.WorkHours = workHours

	if err := s.db.Save(&attendance).Error; err != nil {
		return nil, err
	}

	// Preload employee data
	if err := s.db.Preload("Employee").First(&attendance, attendance.ID).Error; err != nil {
		return nil, err
	}

	return &attendance, nil
}

// GetAttendanceByID retrieves an attendance record by ID
func (s *AttendanceService) GetAttendanceByID(id uint) (*models.Attendance, error) {
	var attendance models.Attendance
	result := s.db.Preload("Employee").First(&attendance, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrAttendanceNotFound
		}
		return nil, result.Error
	}
	return &attendance, nil
}

// GetTodayAttendance retrieves today's attendance for an employee
func (s *AttendanceService) GetTodayAttendance(employeeID uint) (*models.Attendance, error) {
	today := startOfLocalDay()
	var attendance models.Attendance
	result := s.db.Preload("Employee").
		Where("employee_id = ? AND date >= ? AND date < ?", 
			employeeID, today, today.Add(24*time.Hour)).
		First(&attendance)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrAttendanceNotFound
		}
		return nil, result.Error
	}
	return &attendance, nil
}

// GetAttendanceByDateRange retrieves attendance records for a date range
func (s *AttendanceService) GetAttendanceByDateRange(employeeID *uint, startDate, endDate time.Time) ([]models.Attendance, error) {
	var attendances []models.Attendance
	query := s.db.Preload("Employee").
		Where("date >= ? AND date < ?", startDate, endDate.Add(24*time.Hour))

	if employeeID != nil {
		query = query.Where("employee_id = ?", *employeeID)
	}

	result := query.Order("date DESC, check_in DESC").Find(&attendances)
	if result.Error != nil {
		return nil, result.Error
	}

	return attendances, nil
}

// GetAttendanceReport generates an attendance report for a date range
func (s *AttendanceService) GetAttendanceReport(startDate, endDate time.Time) ([]map[string]interface{}, error) {
	var results []struct {
		EmployeeID   uint
		FullName     string
		Position     string
		TotalDays    int64
		TotalHours   float64
		AverageHours float64
	}

	// Use DATE() function to compare only the date part, ignoring time and timezone
	err := s.db.Model(&models.Attendance{}).
		Select(`
			attendances.employee_id,
			employees.full_name,
			employees.position,
			COUNT(*) as total_days,
			SUM(attendances.work_hours) as total_hours,
			AVG(attendances.work_hours) as average_hours
		`).
		Joins("JOIN employees ON employees.id = attendances.employee_id").
		Where("DATE(attendances.date) >= ? AND DATE(attendances.date) <= ?", 
			startDate.Format("2006-01-02"), 
			endDate.Format("2006-01-02")).
		Group("attendances.employee_id, employees.full_name, employees.position").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	fmt.Printf("[GetAttendanceReport] Query returned %d results for period %s to %s\n", 
		len(results), startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	// Convert to map for easier JSON serialization
	report := make([]map[string]interface{}, len(results))
	for i, r := range results {
		report[i] = map[string]interface{}{
			"employee_id":   r.EmployeeID,
			"full_name":     r.FullName,
			"position":      r.Position,
			"total_days":    r.TotalDays,
			"total_hours":   r.TotalHours,
			"average_hours": r.AverageHours,
		}
	}

	return report, nil
}

// GetAttendanceStats returns attendance statistics
func (s *AttendanceService) GetAttendanceStats(date time.Time) (map[string]interface{}, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var totalCheckedIn int64
	var totalCheckedOut int64
	var totalEmployees int64

	// Count total active employees
	if err := s.db.Model(&models.Employee{}).Where("is_active = ?", true).Count(&totalEmployees).Error; err != nil {
		return nil, err
	}

	// Count checked in today
	if err := s.db.Model(&models.Attendance{}).
		Where("date >= ? AND date < ?", startOfDay, endOfDay).
		Count(&totalCheckedIn).Error; err != nil {
		return nil, err
	}

	// Count checked out today
	if err := s.db.Model(&models.Attendance{}).
		Where("date >= ? AND date < ? AND check_out IS NOT NULL", startOfDay, endOfDay).
		Count(&totalCheckedOut).Error; err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"date":              date.Format("2006-01-02"),
		"total_employees":   totalEmployees,
		"checked_in":        totalCheckedIn,
		"checked_out":       totalCheckedOut,
		"not_checked_in":    totalEmployees - totalCheckedIn,
		"attendance_rate":   float64(totalCheckedIn) / float64(totalEmployees) * 100,
	}

	return stats, nil
}

// WiFi Configuration Management

// CreateWiFiConfig creates a new authorized Wi-Fi network
func (s *AttendanceService) CreateWiFiConfig(config *models.WiFiConfig) error {
	if config.SSID == "" || config.BSSID == "" {
		return fmt.Errorf("SSID dan BSSID tidak boleh kosong")
	}

	// Check for duplicate
	var existing models.WiFiConfig
	result := s.db.Where("ssid = ? AND bssid = ?", config.SSID, config.BSSID).First(&existing)
	if result.Error == nil {
		return fmt.Errorf("konfigurasi Wi-Fi sudah ada")
	}

	return s.db.Create(config).Error
}

// GetWiFiConfigByID retrieves a Wi-Fi config by ID
func (s *AttendanceService) GetWiFiConfigByID(id uint) (*models.WiFiConfig, error) {
	var config models.WiFiConfig
	result := s.db.First(&config, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("konfigurasi Wi-Fi tidak ditemukan")
		}
		return nil, result.Error
	}
	return &config, nil
}

// GetAllWiFiConfigs retrieves all Wi-Fi configurations
func (s *AttendanceService) GetAllWiFiConfigs(activeOnly bool) ([]models.WiFiConfig, error) {
	var configs []models.WiFiConfig
	query := s.db.Model(&models.WiFiConfig{})

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	result := query.Order("location ASC").Find(&configs)
	if result.Error != nil {
		return nil, result.Error
	}

	return configs, nil
}

// UpdateWiFiConfig updates a Wi-Fi configuration
func (s *AttendanceService) UpdateWiFiConfig(id uint, updates map[string]interface{}) (*models.WiFiConfig, error) {
	config, err := s.GetWiFiConfigByID(id)
	if err != nil {
		return nil, err
	}

	if err := s.db.Model(config).Updates(updates).Error; err != nil {
		return nil, err
	}

	return s.GetWiFiConfigByID(id)
}

// DeleteWiFiConfig deletes a Wi-Fi configuration
func (s *AttendanceService) DeleteWiFiConfig(id uint) error {
	config, err := s.GetWiFiConfigByID(id)
	if err != nil {
		return err
	}

	return s.db.Delete(config).Error
}

// ToggleWiFiConfigStatus toggles the active status of a Wi-Fi configuration
func (s *AttendanceService) ToggleWiFiConfigStatus(id uint) (*models.WiFiConfig, error) {
	config, err := s.GetWiFiConfigByID(id)
	if err != nil {
		return nil, err
	}

	config.IsActive = !config.IsActive
	if err := s.db.Save(config).Error; err != nil {
		return nil, err
	}

	return config, nil
}

// GPS Configuration Management

// ValidateGPS validates if the provided GPS coordinates are within authorized locations
func (s *AttendanceService) ValidateGPS(latitude, longitude float64) (bool, *models.GPSConfig, error) {
	var configs []models.GPSConfig
	result := s.db.Where("is_active = ?", true).Find(&configs)
	
	if result.Error != nil {
		return false, nil, result.Error
	}

	fmt.Printf("[ValidateGPS] Checking coordinates: lat=%f, lng=%f against %d active configs\n", 
		latitude, longitude, len(configs))

	for _, config := range configs {
		distance := s.calculateDistance(latitude, longitude, config.Latitude, config.Longitude)
		fmt.Printf("[ValidateGPS] Config ID: %d, Name: %s, Distance: %.2f meters, Radius: %d meters\n", 
			config.ID, config.Name, distance, config.Radius)
		
		if distance <= float64(config.Radius) {
			fmt.Printf("[ValidateGPS] GPS MATCH - Config ID: %d, Name: %s\n", config.ID, config.Name)
			return true, &config, nil
		}
	}

	fmt.Printf("[ValidateGPS] NO MATCH - Coordinates not within any authorized location\n")
	return false, nil, nil
}

// calculateDistance calculates the distance between two GPS coordinates using Haversine formula
func (s *AttendanceService) calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000 // Earth's radius in meters
	
	// Convert to radians
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180
	
	// Haversine formula
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	
	return earthRadius * c
}

// CreateGPSConfig creates a new authorized GPS location
func (s *AttendanceService) CreateGPSConfig(config *models.GPSConfig) error {
	if config.Name == "" {
		return fmt.Errorf("nama lokasi tidak boleh kosong")
	}
	if config.Radius <= 0 {
		config.Radius = 100 // Default 100 meters
	}

	return s.db.Create(config).Error
}

// GetGPSConfigByID retrieves a GPS config by ID
func (s *AttendanceService) GetGPSConfigByID(id uint) (*models.GPSConfig, error) {
	var config models.GPSConfig
	result := s.db.First(&config, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("konfigurasi GPS tidak ditemukan")
		}
		return nil, result.Error
	}
	return &config, nil
}

// GetAllGPSConfigs retrieves all GPS configurations
func (s *AttendanceService) GetAllGPSConfigs(activeOnly bool) ([]models.GPSConfig, error) {
	var configs []models.GPSConfig
	query := s.db.Model(&models.GPSConfig{})

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	result := query.Order("name ASC").Find(&configs)
	if result.Error != nil {
		return nil, result.Error
	}

	return configs, nil
}

// UpdateGPSConfig updates a GPS configuration
func (s *AttendanceService) UpdateGPSConfig(id uint, updates map[string]interface{}) (*models.GPSConfig, error) {
	config, err := s.GetGPSConfigByID(id)
	if err != nil {
		return nil, err
	}

	if err := s.db.Model(config).Updates(updates).Error; err != nil {
		return nil, err
	}

	return s.GetGPSConfigByID(id)
}

// DeleteGPSConfig deletes a GPS configuration
func (s *AttendanceService) DeleteGPSConfig(id uint) error {
	config, err := s.GetGPSConfigByID(id)
	if err != nil {
		return err
	}

	return s.db.Delete(config).Error
}

// ToggleGPSConfigStatus toggles the active status of a GPS configuration
func (s *AttendanceService) ToggleGPSConfigStatus(id uint) (*models.GPSConfig, error) {
	config, err := s.GetGPSConfigByID(id)
	if err != nil {
		return nil, err
	}

	config.IsActive = !config.IsActive
	if err := s.db.Save(config).Error; err != nil {
		return nil, err
	}

	return config, nil
}
