package models

import (
	"time"

	"gorm.io/gorm"
)

// Specialization represents the kepegawaian_jabatan_specializations table in database.
// FHIR R4: PractitionerRole.specialty (KFA/SNOMED CT sesuai Terminologi SATUSEHAT)
type Specialization struct {
	ID int64 `gorm:"primaryKey;autoIncrement;column:id" json:"id"`

	Code                string  `gorm:"column:code;type:varchar(20);not null" json:"code"`
	Label               string  `gorm:"column:label;type:varchar(150);not null" json:"label"`
	JobTitleID          *int64  `gorm:"column:job_title_id" json:"job_title_id"`
	KategoriID          *int64  `gorm:"column:kategori_id" json:"kategori_id"`      // FK -> specialization_kategoris (bedah | non_bedah | penunjang)
	Gelar               *string `gorm:"column:gelar;type:varchar(20)" json:"gelar"` // mis. Sp.A, Sp.B, Sp.OG, Sp.PD
	LamaPendidikanTahun *int16  `gorm:"column:lama_pendidikan_tahun;type:smallint" json:"lama_pendidikan_tahun"`
	FHIRCode            *string `gorm:"column:fhir_code;type:varchar(50)" json:"fhir_code"`
	FHIRSystem          *string `gorm:"column:fhir_system;type:varchar(200)" json:"fhir_system"`
	IsAktif             bool    `gorm:"column:is_aktif;not null;default:true" json:"is_aktif"`

	CreatedBy *int64         `gorm:"column:created_by" json:"created_by"`
	UpdatedBy *int64         `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz;index" json:"deleted_at"`

	// Relations
	JobTitle *JobTitle               `gorm:"foreignKey:JobTitleID" json:"job_title,omitempty"`
	Kategori *SpecializationKategori `gorm:"foreignKey:KategoriID" json:"kategori,omitempty"`
}

func (Specialization) TableName() string {
	return "kepegawaian_jabatan_specializations"
}
