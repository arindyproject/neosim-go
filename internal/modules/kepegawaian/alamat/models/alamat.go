package models

import (
	pegawai "neosim_go/internal/modules/kepegawaian/pegawai/models"
	"time"

	"gorm.io/gorm"
)

// KepegawaianAlamat represents the kepegawaian_alamats table in database
type KepegawaianAlamat struct {
	ID        int64 `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	PegawaiID int64 `gorm:"column:pegawai_id;not null;index" json:"pegawai_id"`

	// Foreign key ke tabel kepegawaian_kontak_tipes
	TipeID int64 `gorm:"column:tipe_id;not null;index" json:"tipe_id"`

	Jalan   string  `gorm:"column:jalan;type:text;not null" json:"jalan"`
	RT      *string `gorm:"column:rt;type:varchar(20);" json:"rt"`
	RW      *string `gorm:"column:rw;type:varchar(20);" json:"rw"`
	KodePos *string `gorm:"column:kode_pos;type:varchar(20);" json:"kode_pos"`

	// koneksi ke alamat
	NegaraID        *int64 `gorm:"column:negara_id" json:"negara_id"`
	ProvinsiID      *int64 `gorm:"column:provinsi_id" json:"provinsi_id"`
	KotaKabupatenID *int64 `gorm:"column:kota_kabupaten_id" json:"kota_kabupaten_id"`
	KecamatanID     *int64 `gorm:"column:kecamatan_id" json:"kecamatan_id"`
	KelurahanDesaID *int64 `gorm:"column:kelurahan_desa_id" json:"kelurahan_desa_id"`

	// Apakah identifier ini yang utama untuk tipe tersebut (contoh: dokter punya 2 SIP)
	IsPrimary bool `gorm:"column:is_primary;default:false" json:"is_primary"`

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

func (KepegawaianAlamat) TableName() string {
	return "kepegawaian_alamats"
}
