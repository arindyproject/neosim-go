package models

import (
	pegawai "neosim_go/internal/modules/kepegawaian/pegawai/models"
	"time"

	"gorm.io/gorm"
)

// KepegawaianKualifikasi represents the kepegawaian_kualifikasis table in database
type KepegawaianKualifikasi struct {
	ID        int64 `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	PegawaiID int64 `gorm:"column:pegawai_id;not null;index" json:"pegawai_id"`

	// Foreign key ke tabel kepegawaian_kualifikasi_tipes
	TipeID int64 `gorm:"column:tipe_id;not null;index" json:"tipe_id"`

	Nama          string `gorm:"column:nama;type:varchar(255);not null" json:"nama"`
	Penyelenggara string `gorm:"column:penyelenggara;type:varchar(255);not null" json:"penyelenggara"`

	NomorSertifikat *string `gorm:"column:nomor_sertifikat;type:varchar(255)" json:"nomor_sertifikat"`

	TanggalTerbit  *time.Time `gorm:"column:tanggal_terbit;type:date" json:"tanggal_terbit"`
	TanggalExpired *time.Time `gorm:"column:tanggal_expired;type:date" json:"tanggal_expired"`

	IsAktif bool `gorm:"column:is_aktif;default:true" json:"is_aktif"`

	FhirCode   *string `gorm:"column:fhir_code;type:varchar(255)" json:"fhir_code"`
	FhirSystem *string `gorm:"column:fhir_system;type:varchar(255)" json:"fhir_system"`

	CreatedBy *int64         `gorm:"column:created_by" json:"created_by"`
	UpdatedBy *int64         `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"`

	// Relasi
	Pegawai *pegawai.KepegawaianPegawai `gorm:"foreignKey:PegawaiID" json:"pegawai,omitempty"`
	Tipe    *Tipe                       `gorm:"foreignKey:TipeID" json:"tipe,omitempty"`
}

func (KepegawaianKualifikasi) TableName() string {
	return "kepegawaian_kualifikasis"
}

// IsExpired mengecek apakah identifier ini sudah expired
func (ki *KepegawaianKualifikasi) IsExpired() bool {
	if ki.TanggalExpired == nil {
		return false
	}
	return time.Now().After(*ki.TanggalExpired)
}

// DaysUntilExpired mengembalikan sisa hari sebelum expired, -1 jika tidak ada expired date
func (ki *KepegawaianKualifikasi) DaysUntilExpired() int {
	if ki.TanggalExpired == nil {
		return -1
	}
	return int(time.Until(*ki.TanggalExpired).Hours() / 24)
}
