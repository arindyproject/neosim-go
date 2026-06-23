package repositories

import (
	"neosim_go/internal/modules/kepegawaian/jabatan/contracts"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewKepegawaianJabatanRepository membuat instance repository baru
func NewKepegawaianJabatanRepository(db *gorm.DB) contracts.Repository {
	return &repository{db: db}
}
