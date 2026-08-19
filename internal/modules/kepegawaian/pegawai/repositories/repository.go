package repositories

import (
	"neosim_go/internal/modules/kepegawaian/pegawai/contracts"
	"gorm.io/gorm"
)

// repository adalah satu-satunya struct repository untuk sub-module ini.
// Item baru (mode add-item) TIDAK membuat struct repository baru — method
// CRUD-nya ditempelkan langsung ke struct ini di file terpisah
// (mis. repositories/tag_repository.go), sehingga satu instance struct ini
// otomatis memenuhi contracts.Repository maupun interface item (TagRepository, dst).
type repository struct {
	db *gorm.DB
}

// NewKepegawaianPegawaiRepository membuat instance repository baru
func NewKepegawaianPegawaiRepository(db *gorm.DB) contracts.Repository {
	return &repository{db: db}
}
