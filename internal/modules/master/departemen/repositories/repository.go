package repositories

import (
	"neosim_go/internal/modules/master/departemen/contracts"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewMasterDepartemenRepository membuat instance repository baru
func NewMasterDepartemenRepository(db *gorm.DB) contracts.Repository {
	return &repository{db: db}
}
