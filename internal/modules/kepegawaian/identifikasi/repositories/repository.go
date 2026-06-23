package repositories

import (
	"neosim_go/internal/modules/kepegawaian/identifikasi/contracts"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewKepegawaianIdentifikasiRepository membuat instance repository baru
func NewKepegawaianIdentifikasiRepository(db *gorm.DB) contracts.Repository {
	return &repository{db: db}
}
