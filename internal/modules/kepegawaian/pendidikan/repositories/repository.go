package repositories

import (
	"neosim_go/internal/modules/kepegawaian/pendidikan/contracts"
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

// NewKepegawaianPendidikanRepository membuat instance repository baru
func NewKepegawaianPendidikanRepository(db *gorm.DB) contracts.Repository {
	return &repository{db: db}
}
