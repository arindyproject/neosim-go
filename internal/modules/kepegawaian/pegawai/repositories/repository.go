package repositories

import (
	"neosim_go/internal/modules/kepegawaian/pegawai/contracts"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewKepegawaianPegawaiRepository membuat instance repository baru
func NewKepegawaianPegawaiRepository(db *gorm.DB) contracts.Repository {
	return &repository{db: db}
}
