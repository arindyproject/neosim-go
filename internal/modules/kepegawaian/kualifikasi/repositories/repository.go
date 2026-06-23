package repositories

import (
	"neosim_go/internal/modules/kepegawaian/kualifikasi/contracts"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewKepegawaianKualifikasiRepository membuat instance repository baru
func NewKepegawaianKualifikasiRepository(db *gorm.DB) contracts.Repository {
	return &repository{db: db}
}
