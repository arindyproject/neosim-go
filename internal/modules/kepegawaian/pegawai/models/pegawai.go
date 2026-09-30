package models

import (
	"time"

	masterModels "neosim_go/internal/modules/master/master/models"

	"gorm.io/gorm"
)

// KepegawaianPegawai represents the kepegawaian_pegawais table in database
// FHIR R4: Practitioner
type KepegawaianPegawai struct {
	ID     int64  `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	UserID *int64 `gorm:"column:user_id;uniqueIndex:uq_kepegawaian_pegawais_user,where:user_id IS NOT NULL AND deleted_at IS NULL" json:"user_id"`

	// Identitas
	NIK          string  `gorm:"column:nik;type:varchar(16);not null;uniqueIndex:uq_kepegawaian_pegawais_nik,where:deleted_at IS NULL" json:"nik"`
	IHSNumber    *string `gorm:"column:ihs_number;type:varchar(50);uniqueIndex:uq_kepegawaian_pegawais_ihs,where:ihs_number IS NOT NULL AND deleted_at IS NULL" json:"ihs_number"`
	NomorPegawai string  `gorm:"column:nomor_pegawai;type:varchar(30);not null;uniqueIndex:uq_kepegawaian_pegawais_nomor,where:deleted_at IS NULL" json:"nomor_pegawai"`
	NamaLengkap  string  `gorm:"column:nama_lengkap;type:varchar(150);not null;index:idx_kepegawaian_pegawais_nama" json:"nama_lengkap"`

	// Biodata
	JenisKelaminID     int64     `gorm:"column:jenis_kelamin_id;not null;index:idx_kepegawaian_pegawais_jk" json:"jenis_kelamin_id"`
	TanggalLahir       time.Time `gorm:"column:tanggal_lahir;type:date;not null" json:"tanggal_lahir"`
	TempatLahir        string    `gorm:"column:tempat_lahir;type:varchar(100);not null" json:"tempat_lahir"`
	GolonganDarahID    *int64    `gorm:"column:golongan_darah_id;index:idx_kepegawaian_pegawais_goldar" json:"golongan_darah_id"`
	AgamaID            int64     `gorm:"column:agama_id;not null;index:idx_kepegawaian_pegawais_agama" json:"agama_id"`
	StatusPernikahanID int64     `gorm:"column:status_pernikahan_id;not null;index:idx_kepegawaian_pegawais_pernikahan" json:"status_pernikahan_id"`
	Kewarganegaraan    *string   `gorm:"column:kewarganegaraan;type:varchar(3)" json:"kewarganegaraan"` // WNI | WNA, nullable

	// Kepegawaian
	TanggalMasuk  time.Time  `gorm:"column:tanggal_masuk;type:date;not null" json:"tanggal_masuk"`
	TanggalKeluar *time.Time `gorm:"column:tanggal_keluar;type:date" json:"tanggal_keluar"`
	JenisID       int64      `gorm:"column:jenis_id;not null;index:idx_kepegawaian_pegawais_jenis" json:"jenis_id"`
	StatusID      int64      `gorm:"column:status_id;not null;index:idx_kepegawaian_pegawais_status" json:"status_id"`
	FotoURL       *string    `gorm:"column:foto_url;type:varchar(500)" json:"foto_url"`
	IsAktif       bool       `gorm:"column:is_aktif;not null;default:true" json:"is_aktif"`

	// Audit
	CreatedBy *int64         `gorm:"column:created_by" json:"created_by"`
	UpdatedBy *int64         `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"`

	// Relasi master modul pegawai (package models yang sama)
	Jenis  *Jenis  `gorm:"foreignKey:JenisID;references:ID" json:"jenis,omitempty"`
	Status *Status `gorm:"foreignKey:StatusID;references:ID" json:"status,omitempty"`

	// Relasi master modul master/master
	JenisKelamin     *masterModels.MasterJenisKelamin     `gorm:"foreignKey:JenisKelaminID;references:ID" json:"jenis_kelamin,omitempty"`
	GolonganDarah    *masterModels.MasterGolonganDarah    `gorm:"foreignKey:GolonganDarahID;references:ID" json:"golongan_darah,omitempty"`
	Agama            *masterModels.MasterAgama            `gorm:"foreignKey:AgamaID;references:ID" json:"agama,omitempty"`
	StatusPernikahan *masterModels.MasterStatusPernikahan `gorm:"foreignKey:StatusPernikahanID;references:ID" json:"status_pernikahan,omitempty"`

	// Relasi ke tabel turunan
	//Identifiers []KepegawaianIdentifier  `gorm:"foreignKey:PegawaiID" json:"identifiers,omitempty"`
	//Pendidikan  []KepegawaianPendidikan  `gorm:"foreignKey:PegawaiID" json:"pendidikan,omitempty"`
	//Kualifikasi []KepegawaianKualifikasi `gorm:"foreignKey:PegawaiID" json:"kualifikasi,omitempty"`
	//Kontak      []KepegawaianKontak      `gorm:"foreignKey:PegawaiID" json:"kontak,omitempty"`
	//Alamat      []KepegawaianAlamat      `gorm:"foreignKey:PegawaiID" json:"alamat,omitempty"`
	//Jabatan     []KepegawaianJabatan     `gorm:"foreignKey:PegawaiID" json:"jabatan,omitempty"`
}

func (KepegawaianPegawai) TableName() string {
	return "kepegawaian_pegawais"
}
