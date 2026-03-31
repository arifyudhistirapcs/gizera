package models

import (
	"time"
)

// Yayasan merepresentasikan yayasan/lembaga yang mengelola SPPG
type Yayasan struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Kode            string    `gorm:"uniqueIndex;size:20;not null" json:"kode"`
	Nama            string    `gorm:"size:200;not null;index" json:"nama" validate:"required"`
	Alamat          string    `gorm:"type:text" json:"alamat"`
	Latitude        float64   `gorm:"default:0" json:"latitude"`
	Longitude       float64   `gorm:"default:0" json:"longitude"`
	NomorTelepon    string    `gorm:"size:20" json:"nomor_telepon"`
	Email           string    `gorm:"uniqueIndex;size:100" json:"email" validate:"omitempty,email"`
	PenanggungJawab string    `gorm:"size:100" json:"penanggung_jawab"`
	NPWP            string    `gorm:"uniqueIndex;size:30" json:"npwp"`
	IsActive        bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	SPPGs           []SPPG    `gorm:"foreignKey:YayasanID" json:"sppgs,omitempty"`
}

// SPPG merepresentasikan Satuan Pelayanan Pemenuhan Gizi (tenant)
type SPPG struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Kode         string    `gorm:"uniqueIndex;size:20;not null" json:"kode"`
	Nama         string    `gorm:"size:200;not null;index" json:"nama" validate:"required"`
	Alamat       string    `gorm:"type:text" json:"alamat"`
	Latitude     float64   `gorm:"default:0" json:"latitude"`
	Longitude    float64   `gorm:"default:0" json:"longitude"`
	NomorTelepon string    `gorm:"size:20" json:"nomor_telepon"`
	Email        string    `gorm:"uniqueIndex;size:100" json:"email" validate:"omitempty,email"`
	YayasanID    uint      `gorm:"index;not null" json:"yayasan_id"`
	IsActive     bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Yayasan      Yayasan   `gorm:"foreignKey:YayasanID" json:"yayasan,omitempty"`
}
