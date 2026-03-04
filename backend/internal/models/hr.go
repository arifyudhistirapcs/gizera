package models

import (
	"time"
)

// Employee represents an employee in the organization
type Employee struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"uniqueIndex;not null" json:"user_id"` // links to User table
	NIK         string    `gorm:"uniqueIndex;size:20;not null" json:"nik" validate:"required"`
	FullName    string    `gorm:"size:100;not null;index" json:"full_name" validate:"required"`
	Email       string    `gorm:"uniqueIndex;size:100;not null" json:"email" validate:"required,email"`
	PhoneNumber string    `gorm:"size:20" json:"phone_number"`
	Position    string    `gorm:"size:100;index" json:"position"`
	JoinDate    time.Time `gorm:"not null" json:"join_date"`
	IsActive    bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// Attendance represents employee attendance records
type Attendance struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	EmployeeID uint       `gorm:"index;not null" json:"employee_id"`
	Date       time.Time  `gorm:"index;not null" json:"date"`
	CheckIn    time.Time  `gorm:"not null" json:"check_in"`
	CheckOut   *time.Time `json:"check_out"`
	WorkHours  float64    `gorm:"default:0" json:"work_hours"`
	SSID       string     `gorm:"column:ss_id;size:100" json:"ssid"`
	BSSID      string     `gorm:"column:bss_id;size:100" json:"bssid"`
	CreatedAt  time.Time  `json:"created_at"`
	Employee   Employee   `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
}

// WiFiConfig represents authorized Wi-Fi networks for attendance
type WiFiConfig struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	SSID       string    `gorm:"column:ss_id;size:100;not null;index" json:"ssid" validate:"required"`
	BSSID      string    `gorm:"column:bss_id;size:100;not null;index" json:"bssid" validate:"required"`
	Location   string    `gorm:"size:200" json:"location"`
	IPRange    string    `gorm:"size:100" json:"ip_range"`           // e.g., "192.168.1.0/24"
	AllowedIPs []string  `gorm:"type:text[]" json:"allowed_ips"`     // Specific IPs allowed
	IsActive   bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TableName specifies the table name for WiFiConfig
func (WiFiConfig) TableName() string {
	return "wi_fi_configs"
}

// GPSConfig represents authorized GPS locations for attendance
type GPSConfig struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name" validate:"required"`
	Latitude    float64   `gorm:"not null" json:"latitude" validate:"required"`
	Longitude   float64   `gorm:"not null" json:"longitude" validate:"required"`
	Radius      int       `gorm:"not null;default:100" json:"radius"` // Radius in meters
	Address     string    `gorm:"size:500" json:"address"`
	Description string    `gorm:"size:500" json:"description"`
	IsActive    bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName specifies the table name for GPSConfig
func (GPSConfig) TableName() string {
	return "gps_configs"
}
