package models

import (
	"time"

	pegawai "neosim_go/internal/modules/kepegawaian/pegawai/models"

	"gorm.io/gorm"
)

// KepegawaianIdentifier represents the kepegawaian_identifiers table in database
// FHIR R4: Practitioner.identifier[]
type KepegawaianIdentifier struct {
	ID        int64 `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	PegawaiID int64 `gorm:"column:pegawai_id;not null;index" json:"pegawai_id"`

	// Foreign key ke tabel kepegawaian_identifier_tipes
	TipeID int64 `gorm:"column:tipe_id;not null;index" json:"tipe_id"`

	// Nilai / nomor identifier
	Nilai string `gorm:"column:nilai;type:varchar(225);not null" json:"nilai"`

	// Tanggal terbit dan expired (wajib jika Tipe.HasExpiry = true)
	TanggalTerbit  *time.Time `gorm:"column:tanggal_terbit;type:date" json:"tanggal_terbit"`
	TanggalExpired *time.Time `gorm:"column:tanggal_expired;type:date" json:"tanggal_expired"`

	// Apakah identifier ini yang utama untuk tipe tersebut (contoh: dokter punya 2 SIP)
	IsPrimary bool `gorm:"column:is_primary;default:false" json:"is_primary"`
	IsAktif   bool `gorm:"column:is_aktif;default:true" json:"is_aktif"`

	// Audit fields
	CreatedBy *int64         `gorm:"column:created_by" json:"created_by"`
	UpdatedBy *int64         `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz;index" json:"deleted_at"`

	// Relasi
	Pegawai *pegawai.KepegawaianPegawai `gorm:"foreignKey:PegawaiID" json:"pegawai,omitempty"`
	Tipe    *Tipe                       `gorm:"foreignKey:TipeID" json:"tipe,omitempty"`
}

func (KepegawaianIdentifier) TableName() string {
	return "kepegawaian_identifiers"
}

// IsFHIRMappable mengecek apakah identifier ini bisa disync ke SATUSEHAT
func (ki *KepegawaianIdentifier) IsFHIRMappable() bool {
	return ki.Tipe != nil && ki.Tipe.FHIRSystem != nil && *ki.Tipe.FHIRSystem != ""
}

// IsExpired mengecek apakah identifier ini sudah expired
func (ki *KepegawaianIdentifier) IsExpired() bool {
	if ki.TanggalExpired == nil {
		return false
	}
	return time.Now().After(*ki.TanggalExpired)
}

// DaysUntilExpired mengembalikan sisa hari sebelum expired, -1 jika tidak ada expired date
func (ki *KepegawaianIdentifier) DaysUntilExpired() int {
	if ki.TanggalExpired == nil {
		return -1
	}
	return int(time.Until(*ki.TanggalExpired).Hours() / 24)
}
