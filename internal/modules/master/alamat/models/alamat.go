package models

import (
	"time"

	"gorm.io/gorm"
)

// -----------------------------------------------------------------------------------------------------------
// Negara
type MasterAlamatNegara struct {
	ID          int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Code        string         `gorm:"column:code;type:varchar(5);not null;unique" json:"code"` // Cth: ID, SG
	Name        string         `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description *string        `gorm:"column:description;type:text" json:"description"`
	CreatedBy   *int64         `gorm:"column:created_by" json:"created_by"`
	UpdatedBy   *int64         `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt   time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"`
}

func (MasterAlamatNegara) TableName() string {
	return "master_alamat_negara"
} // ---------------------------------------------------------------------------------------------------------

// -----------------------------------------------------------------------------------------------------------
// Provinsi
type MasterAlamatProvinsi struct {
	ID        int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	NegaraID  int64          `gorm:"column:negara_id;not null" json:"negara_id"`
	Code      string         `gorm:"column:code;type:varchar(10);not null;unique" json:"code"` // Cth: 35
	Name      string         `gorm:"column:name;type:varchar(255);not null" json:"name"`
	CreatedBy *int64         `gorm:"column:created_by" json:"created_by"`
	UpdatedBy *int64         `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"`

	// Ubah bagian ini agar GORM membuat nama foreign key yang lebih pendek
	Negara MasterAlamatNegara `gorm:"foreignKey:NegaraID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"negara,omitempty"`
}

func (MasterAlamatProvinsi) TableName() string {
	return "master_alamat_provinsi"
} // ---------------------------------------------------------------------------------------------------------

// -----------------------------------------------------------------------------------------------------------
// Kota/Kabupaten
type MasterAlamatKotaKabupaten struct {
	ID         int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	ProvinsiID int64          `gorm:"column:provinsi_id;not null" json:"provinsi_id"`
	Code       string         `gorm:"column:code;type:varchar(10);not null;unique" json:"code"` // Cth: 35.21
	Name       string         `gorm:"column:name;type:varchar(255);not null" json:"name"`
	CreatedBy  *int64         `gorm:"column:created_by" json:"created_by"`
	UpdatedBy  *int64         `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt  time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"`

	Provinsi MasterAlamatProvinsi `gorm:"foreignKey:ProvinsiID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"provinsi,omitempty"`
}

func (MasterAlamatKotaKabupaten) TableName() string {
	return "master_alamat_kota_kabupaten"
} // ---------------------------------------------------------------------------------------------------------

// -----------------------------------------------------------------------------------------------------------
// Kecamatan
type MasterAlamatKecamatan struct {
	ID              int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	KotaKabupatenID int64          `gorm:"column:kota_kabupaten_id;not null" json:"kota_kabupaten_id"`
	Code            string         `gorm:"column:code;type:varchar(10);not null;unique" json:"code"` // Cth: 35.21.01
	Name            string         `gorm:"column:name;type:varchar(255);not null" json:"name"`
	CreatedBy       *int64         `gorm:"column:created_by" json:"created_by"`
	UpdatedBy       *int64         `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt       time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"`

	KotaKabupaten MasterAlamatKotaKabupaten `gorm:"foreignKey:KotaKabupatenID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"kota_kabupaten,omitempty"`
}

func (MasterAlamatKecamatan) TableName() string {
	return "master_alamat_kecamatan"
} // ---------------------------------------------------------------------------------------------------------

// -----------------------------------------------------------------------------------------------------------
// Desa
type MasterAlamatKelurahanDesa struct {
	ID          int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	KecamatanID int64          `gorm:"column:kecamatan_id;not null" json:"kecamatan_id"`
	Code        string         `gorm:"column:code;type:varchar(15);not null;unique" json:"code"` // Cth: 35.21.01.2001
	Name        string         `gorm:"column:name;type:varchar(255);not null" json:"name"`
	PostalCode  *string        `gorm:"column:postal_code;type:varchar(10)" json:"postal_code"`
	CreatedBy   *int64         `gorm:"column:created_by" json:"created_by"`
	UpdatedBy   *int64         `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt   time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"`

	Kecamatan MasterAlamatKecamatan `gorm:"foreignKey:KecamatanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"kecamatan,omitempty"`
}

func (MasterAlamatKelurahanDesa) TableName() string {
	return "master_alamat_kelurahan_desa"
} // ---------------------------------------------------------------------------------------------------------
