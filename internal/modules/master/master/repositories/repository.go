package repositories

import (
	"neosim_go/internal/modules/master/master/contracts"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewMasterRepository membuat instance repository baru
func NewMasterRepository(db *gorm.DB) contracts.Repository {
	return &repository{db: db}
}
