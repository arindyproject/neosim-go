package repositories

import (
	"neosim_go/internal/modules/kepegawaian/identifier/contracts"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewKepegawaianIdentifierRepository membuat instance repository baru
func NewKepegawaianIdentifierRepository(db *gorm.DB) contracts.Repository {
	return &repository{db: db}
}
