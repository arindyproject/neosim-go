package repositories

import (
	"neosim_go/internal/modules/kepegawaian/lampiran/contracts"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewKepegawaianLampiranRepository membuat instance repository baru
func NewKepegawaianLampiranRepository(db *gorm.DB) contracts.Repository {
	return &repository{db: db}
}
