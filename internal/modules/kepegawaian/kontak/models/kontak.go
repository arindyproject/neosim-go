package models

import (
	"time"

	"gorm.io/gorm"
)

// KepegawaianKontak represents the kepegawaian_kontaks table in database
type KepegawaianKontak struct {
	ID          int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name        string         `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description *string        `gorm:"column:description;type:text" json:"description"`
	CreatedBy   *int64         `gorm:"column:created_by" json:"created_by"`
	UpdatedBy   *int64         `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt   time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz" json:"deleted_at"`
}

func (KepegawaianKontak) TableName() string {
	return "kepegawaian_kontaks"
}
