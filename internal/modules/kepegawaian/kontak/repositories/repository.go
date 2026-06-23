package repositories

import (
	"neosim_go/internal/modules/kepegawaian/kontak/contracts"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewKepegawaianKontakRepository membuat instance repository baru
func NewKepegawaianKontakRepository(db *gorm.DB) contracts.Repository {
	return &repository{db: db}
}
