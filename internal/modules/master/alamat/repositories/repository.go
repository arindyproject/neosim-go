package repositories

import (
	"neosim_go/internal/modules/master/alamat/contracts"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewMasterAlamatRepository membuat instance repository baru
func NewMasterAlamatRepository(db *gorm.DB) contracts.Repository {
	return &repository{db: db}
}
