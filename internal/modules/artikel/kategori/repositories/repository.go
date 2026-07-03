package repositories

import (
	"neosim_go/internal/modules/artikel/kategori/contracts"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewArtikelKategoriRepository membuat instance repository baru
func NewArtikelKategoriRepository(db *gorm.DB) contracts.Repository {
	return &repository{db: db}
}
