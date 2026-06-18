package models

import (
	"time"

	"gorm.io/gorm"
)

// Master represents the masters table in database
type Master struct {
	ID          int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name        string         `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description *string        `gorm:"column:description;type:text" json:"description"`
	CreatedBy   *int64         `gorm:"column:created_by" json:"created_by"`
	UpdatedBy   *int64         `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt   time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"`
}

func (Master) TableName() string {
	return "masters"
}

// Pekerjaan
// ----------------------------------------------------------------------------------------
type MasterPekerjaan struct {
	ID           int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name         string         `gorm:"column:name;type:varchar(255);not null;index:idx_name_active,unique,where:deleted_at IS NULL" json:"name"`
	KodeKemenkes *string        `gorm:"column:kode_kemenkes;type:varchar(50)" json:"kode_kemenkes"`
	Description  *string        `gorm:"column:description;type:text" json:"description"`
	CreatedBy    *int64         `gorm:"column:created_by" json:"created_by"`
	UpdatedBy    *int64         `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt    time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"`
}

func (MasterPekerjaan) TableName() string {
	return "master_pekerjaan"
} //---------------------------------------------------------------------------------------

// Pendidikan
// ----------------------------------------------------------------------------------------
type MasterPendidikan struct {
	ID           int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name         string         `gorm:"column:name;type:varchar(255);not null;index:idx_pendidikan_name_active,unique,where:deleted_at IS NULL" json:"name"`
	KodeKemenkes *string        `gorm:"column:kode_kemenkes;type:varchar(50)" json:"kode_kemenkes"`
	Description  *string        `gorm:"column:description;type:text" json:"description"`
	CreatedBy    *int64         `gorm:"column:created_by" json:"created_by"`
	UpdatedBy    *int64         `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt    time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"`
}

func (MasterPendidikan) TableName() string {
	return "master_pendidikan"
} //---------------------------------------------------------------------------------------

// StatusPernikahan
// ----------------------------------------------------------------------------------------
type MasterStatusPernikahan struct {
	ID           int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name         string         `gorm:"column:name;type:varchar(255);not null;index:idx_marital_name_active,unique,where:deleted_at IS NULL" json:"name"`
	KodeKemenkes *string        `gorm:"column:kode_kemenkes;type:varchar(50)" json:"kode_kemenkes"` // Mapping kode HL7 / SatuSehat
	Description  *string        `gorm:"column:description;type:text" json:"description"`
	CreatedBy    *int64         `gorm:"column:created_by" json:"created_by"`
	UpdatedBy    *int64         `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt    time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"`
}

func (MasterStatusPernikahan) TableName() string {
	return "master_status_pernikahan"
} //---------------------------------------------------------------------------------------

// Agama
// ----------------------------------------------------------------------------------------
type MasterAgama struct {
	ID           int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name         string         `gorm:"column:name;type:varchar(255);not null;index:idx_agama_name_active,unique,where:deleted_at IS NULL" json:"name"`
	KodeKemenkes *string        `gorm:"column:kode_kemenkes;type:varchar(50)" json:"kode_kemenkes"` // Mapping kode SatuSehat / SIRS
	Description  *string        `gorm:"column:description;type:text" json:"description"`
	CreatedBy    *int64         `gorm:"column:created_by" json:"created_by"`
	UpdatedBy    *int64         `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt    time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"`
}

func (MasterAgama) TableName() string {
	return "master_agama"
} //---------------------------------------------------------------------------------------
