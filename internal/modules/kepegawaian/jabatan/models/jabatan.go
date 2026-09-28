package models

import (
	pegawai "neosim_go/internal/modules/kepegawaian/pegawai/models"
	"time"

	"gorm.io/gorm"
)

// KepegawaianJabatan represents the kepegawaian_jabatans table in database.
// FHIR R4: PractitionerRole
//
// Tabel historis penugasan pegawai di suatu unit dan jabatan. Saat pegawai
// mutasi, baris lama di-close (tanggal_selesai diisi) dan baris baru dibuat.
// Satu pegawai bisa punya beberapa jabatan aktif di unit berbeda.
type KepegawaianJabatan struct {
	ID int64 `gorm:"primaryKey;autoIncrement;column:id" json:"id"`

	// FK lintas module: disimpan sebagai int64 polos (bukan GORM association)
	// untuk menghindari import silang antar module di modular monolith.
	PegawaiID    int64 `gorm:"column:pegawai_id;not null;index" json:"pegawai_id"`
	DepartmentID int64 `gorm:"column:department_id;not null;index" json:"department_id"` // FK -> departments (unit tempat bertugas)

	// FK ke master di module jabatan yang sama.
	PositionID       int64  `gorm:"column:position_id;not null;index" json:"position_id"`    // FK -> positions (jabatan struktural)
	JobTitleID       int64  `gorm:"column:job_title_id;not null;index" json:"job_title_id"`  // FK -> job_titles (jabatan fungsional)
	SpecializationID *int64 `gorm:"column:specialization_id;index" json:"specialization_id"` // FK -> specializations (nullable, hanya dokter/profesi tertentu)

	IsPrimary bool `gorm:"column:is_primary;not null" json:"is_primary"`

	TanggalMulai   time.Time  `gorm:"column:tanggal_mulai;type:date;not null" json:"tanggal_mulai"`
	TanggalSelesai *time.Time `gorm:"column:tanggal_selesai;type:date" json:"tanggal_selesai"` // NULL = masih aktif di jabatan tersebut

	NomorSK   *string    `gorm:"column:nomor_sk;type:varchar(100)" json:"nomor_sk"`
	TanggalSK *time.Time `gorm:"column:tanggal_sk;type:date" json:"tanggal_sk"`

	IsAktif bool `gorm:"column:is_aktif;not null" json:"is_aktif"`

	CreatedBy *int64         `gorm:"column:created_by" json:"created_by"`
	UpdatedBy *int64         `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz;index" json:"deleted_at"`

	// Relations (hanya yang berada di module jabatan yang sama)
	Pegawai        *pegawai.KepegawaianPegawai `gorm:"foreignKey:PegawaiID" json:"pegawai,omitempty"`
	Position       *Position                   `gorm:"foreignKey:PositionID" json:"position,omitempty"`
	JobTitle       *JobTitle                   `gorm:"foreignKey:JobTitleID" json:"job_title,omitempty"`
	Specialization *Specialization             `gorm:"foreignKey:SpecializationID" json:"specialization,omitempty"`
}

func (KepegawaianJabatan) TableName() string {
	return "kepegawaian_jabatans"
}
