package models

import (
	pegawai "neosim_go/internal/modules/kepegawaian/pegawai/models"
	"time"

	"gorm.io/gorm"
)

// KepegawaianKontak represents the kepegawaian_kontaks table in database
type KepegawaianKontak struct {
	ID        int64 `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	PegawaiID int64 `gorm:"column:pegawai_id;not null;index" json:"pegawai_id"`

	// Foreign key ke tabel kepegawaian_kontak_tipes
	TipeID int64 `gorm:"column:tipe_id;not null;index" json:"tipe_id"`

	// Nilai / nomor identifier
	Nilai string `gorm:"column:nilai;type:varchar(225);not null" json:"nilai"`

	// Apakah identifier ini yang utama untuk tipe tersebut (contoh: dokter punya 2 SIP)
	IsPrimary bool `gorm:"column:is_primary;default:false" json:"is_primary"`
	IsAktif   bool `gorm:"column:is_aktif;default:true" json:"is_aktif"`

	Description *string        `gorm:"column:description;type:text" json:"description"`
	CreatedBy   *int64         `gorm:"column:created_by" json:"created_by"`
	UpdatedBy   *int64         `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt   time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"`

	// Relasi
	Pegawai *pegawai.KepegawaianPegawai `gorm:"foreignKey:PegawaiID" json:"pegawai,omitempty"`
	Tipe    *Tipe                       `gorm:"foreignKey:TipeID" json:"tipe,omitempty"`
}

func (KepegawaianKontak) TableName() string {
	return "kepegawaian_kontaks"
}
