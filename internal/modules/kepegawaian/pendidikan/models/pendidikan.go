package models

import (
	"time"

	pegawai "neosim_go/internal/modules/kepegawaian/pegawai/models"

	"gorm.io/gorm"
)

// KepegawaianPendidikan represents the kepegawaian_pendidikans table in database
type KepegawaianPendidikan struct {
	ID        int64 `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	PegawaiID int64 `gorm:"column:pegawai_id;not null;index" json:"pegawai_id"`

	// Foreign key ke tabel kepegawaian_pendidikan_jenjangs
	JenjangID int64 `gorm:"column:jenjang_id;not null;index" json:"jenjang_id"`

	NamaInstitusi   string  `gorm:"column:nama_institusi;not null;type:varchar(225)" json:"nama_institusi"`
	NomorIjazah     *string `gorm:"column:nomor_ijazah;type:varchar(255);uniqueIndex" json:"nomor_ijazah"`
	BidangStudi     *string `gorm:"column:bidang_studi;type:varchar(225)" json:"bidang_studi"`
	AlamatInstitusi *string `gorm:"column:alamat_institusi;type:text" json:"alamat_institusi"`

	NilaiAkhir *string `gorm:"column:nilai_akhir;type:varchar(225);" json:"nilai_akhir"`

	TanggalMasuk *time.Time `gorm:"column:tanggal_masuk;type:date" json:"tanggal_masuk"`
	TanggalLulus *time.Time `gorm:"column:tanggal_lulus;type:date" json:"tanggal_lulus"`

	FHIRCode   *string `gorm:"column:fhir_code;type:varchar(255)" json:"fhir_code"`
	FHIRSystem *string `gorm:"column:fhir_system;type:varchar(255)" json:"fhir_system"`

	CreatedBy *int64         `gorm:"column:created_by" json:"created_by"`
	UpdatedBy *int64         `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"`

	Pegawai *pegawai.KepegawaianPegawai `gorm:"foreignKey:PegawaiID" json:"pegawai,omitempty"`
	Jenjang *Jenjang                    `gorm:"foreignKey:JenjangID" json:"jenjang,omitempty"`
}

func (KepegawaianPendidikan) TableName() string {
	return "kepegawaian_pendidikans"
}
