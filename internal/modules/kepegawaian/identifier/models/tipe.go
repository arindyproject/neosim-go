package models

import (
	"time"

	"gorm.io/gorm"
)

// Tipe represents the kepegawaian_identifier_tipes table in database
type Tipe struct {
	ID          int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Code        string         `gorm:"column:code;type:varchar(100);not null;unique" json:"code"`
	Label       string         `gorm:"column:label;type:varchar(255);not null" json:"label"`
	Penerbit    *string        `gorm:"column:penerbit;type:varchar(255)" json:"penerbit"`
	FHIRSystem  *string        `gorm:"column:fhir_system;type:varchar(255)" json:"fhir_system"`
	HasExpiry   bool           `gorm:"column:has_expiry;type:boolean;not null;default:false" json:"has_expiry"`
	IsNakes     bool           `gorm:"column:is_nakes;type:boolean;not null;default:false" json:"is_nakes"`
	IsRequired  bool           `gorm:"column:is_required;type:boolean;not null;default:false" json:"is_required"`
	Description *string        `gorm:"column:description;type:text" json:"description"`
	CreatedBy   *int64         `gorm:"column:created_by" json:"created_by"`
	UpdatedBy   *int64         `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt   time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz;index" json:"deleted_at"`
}

func (Tipe) TableName() string {
	return "kepegawaian_identifier_tipes"
}
