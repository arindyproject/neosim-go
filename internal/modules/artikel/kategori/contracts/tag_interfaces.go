package contracts

import (
	"neosim_go/internal/modules/artikel/kategori/dto"
	"neosim_go/internal/modules/artikel/kategori/models"
	he "neosim_go/internal/shared/httputil"
)

// TagRepository defines database operations for Tag.
// Diimplementasikan oleh struct 'repository' yang sama dengan entitas utama
// sub-module ini (lihat repositories/repository.go) — TIDAK ADA struct baru.
// Method diberi suffix nama item agar tidak bentrok saat di-embed ke
// contracts.Repository.
type TagRepository interface {
	CreateTag(m *models.Tag) error
	GetTagByID(id int64) (*models.Tag, error)
	ListTag(page, pageSize int, filter *dto.FilterTagRequest) ([]models.Tag, int64, error)
	UpdateTag(m *models.Tag) error
	DeleteTag(id int64) error
}

// TagService defines business logic operations for Tag.
// Diimplementasikan oleh struct 'service' yang sama dengan entitas utama.
type TagService interface {
	CreateTag(req *dto.CreateTagRequest, actor he.AuthContext) (*dto.TagResponse, error)
	GetTagByID(id int64, actor he.AuthContext) (*dto.TagResponse, error)
	ListTag(page, pageSize int, filter *dto.FilterTagRequest, actor he.AuthContext) ([]dto.TagResponse, int64, error)
	UpdateTag(id int64, req *dto.UpdateTagRequest, actor he.AuthContext) (*dto.TagResponse, error)
	DeleteTag(id int64, actor he.AuthContext) error
}
