package models

import (
	"time"

	"gorm.io/gorm"
)

// Jenjang represents the kepegawaian_pendidikan_jenjangs table in database
type Jenjang struct {
	ID         int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Code       string         `gorm:"column:code;type:varchar(100);not null;uniqueIndex" json:"code"`
	Label      string         `gorm:"column:label;type:varchar(255);not null;uniqueIndex" json:"label"`
	FHIRSystem *string        `gorm:"column:fhir_system;type:text" json:"fhir_system"`
	CreatedBy  *int64         `gorm:"column:created_by" json:"created_by"`
	UpdatedBy  *int64         `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt  time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"`
}

func (Jenjang) TableName() string {
	return "kepegawaian_pendidikan_jenjangs"
}
