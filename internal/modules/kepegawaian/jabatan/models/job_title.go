package models

import (
	"time"

	"gorm.io/gorm"
)

// JobTitle represents the kepegawaian_jabatan_job_titles table in database
type JobTitle struct {
	ID          int64   `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Code        string  `gorm:"column:code;type:varchar(20);not null;unique" json:"code"`
	Label       string  `gorm:"column:label;type:varchar(255);not null" json:"label"`
	Description *string `gorm:"column:description;type:text" json:"description"`

	// Kategori: medis | paramedis | penunjang_medis | non_medis.
	// Dinormalisasi sebagai FK ke JobTitleKategori (bukan varchar enum),
	// konsisten dengan PositionKategori dan lookup table lain di modul ini.
	KategoriID int64             `gorm:"column:kategori_id;not null;index" json:"kategori_id"`
	Kategori   *JobTitleKategori `gorm:"foreignKey:KategoriID" json:"kategori,omitempty"`

	Point *float64 `gorm:"column:point;type:decimal(10,2);default:0.00" json:"point"`

	// Rumpun profesi: dokter | perawat | bidan | apoteker | tenaga_kesehatan_lain
	// | non_kesehatan — pengelompokan lebih spesifik dari Kategori, selaras
	// rumpun SDM Kesehatan Kemenkes. Opsional karena tidak semua job title
	// (mis. staf non-kesehatan) punya rumpun profesi yang jelas.
	RumpunProfesiID *int64                 `gorm:"column:rumpun_profesi_id;index" json:"rumpun_profesi_id"`
	RumpunProfesi   *JobTitleRumpunProfesi `gorm:"foreignKey:RumpunProfesiID" json:"rumpun_profesi,omitempty"`

	// Dipakai sistem untuk memvalidasi kelengkapan kepegawaian_identifiers
	// sebelum pegawai bisa ditugaskan ke jabatan klinis.
	MemerlukanSTR bool `gorm:"column:memerlukan_str;not null;default:false" json:"memerlukan_str"`
	MemerlukanSIP bool `gorm:"column:memerlukan_sip;not null;default:false" json:"memerlukan_sip"`

	// Syarat pendidikan minimum (SD..S3/Sp2), dicek terhadap
	// kepegawaian_pendidikan.jenjang. Tidak dinormalisasi ke FK karena
	// jenjang di kepegawaian_pendidikan sendiri juga varchar, bukan lookup table.
	JenjangMin *string `gorm:"column:jenjang_min;type:varchar(20)" json:"jenjang_min"`

	// Kode & sistem terminologi Jenis SDM Kesehatan SATUSEHAT milik JobTitle
	// itu sendiri, dipetakan saat sync PractitionerRole.code — beda dari
	// FHIRCode milik JobTitleKategori/JobTitleRumpunProfesi yang mewakili
	// kode klasifikasi kategori/rumpun-nya sendiri.
	FHIRCode   *string `gorm:"column:fhir_code;type:varchar(50)" json:"fhir_code"`
	FHIRSystem *string `gorm:"column:fhir_system;type:varchar(200)" json:"fhir_system"`

	IsAktif bool `gorm:"column:is_aktif;not null;default:true" json:"is_aktif"`

	CreatedBy *int64         `gorm:"column:created_by" json:"created_by"`
	UpdatedBy *int64         `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"`
}

func (JobTitle) TableName() string {
	return "kepegawaian_jabatan_job_titles"
}
