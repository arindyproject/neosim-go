package models

import (
	"time"

	"gorm.io/gorm"
)

// MasterDepartemen represents the master_departemens table in database
type MasterDepartemen struct {
	ID           int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Code         string         `gorm:"column:code;type:varchar(255);not null" json:"code"`
	Name         string         `gorm:"column:name;type:varchar(255);not null" json:"name"`
	SystemModule string         `gorm:"column:system_module;type:varchar(255);not null" json:"system_module"`
	Description  *string        `gorm:"column:description;type:text" json:"description"`
	CreatedBy    *int64         `gorm:"column:created_by" json:"created_by"`
	UpdatedBy    *int64         `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt    time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"`
}

func (MasterDepartemen) TableName() string {
	return "master_departemens"
}
