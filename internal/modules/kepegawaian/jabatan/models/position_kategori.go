package models

import (
	"time"

	"gorm.io/gorm"
)

// PositionKategori represents the kepegawaian_jabatan_position_kategoris table in database
type PositionKategori struct {
	ID        int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Code      string         `gorm:"column:code;type:varchar(100);not null;unique" json:"code"`
	Label     string         `gorm:"column:label;type:varchar(255);not null" json:"label"`
	FHIRCode  *string        `gorm:"column:fhir_code;type:text" json:"fhir_code"`
	Point     *float64       `gorm:"column:point;type:decimal(10,2);default:0.00" json:"point"`
	CreatedBy *int64         `gorm:"column:created_by" json:"created_by"`
	UpdatedBy *int64         `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"`
}

func (PositionKategori) TableName() string {
	return "kepegawaian_jabatan_position_kategoris"
}
