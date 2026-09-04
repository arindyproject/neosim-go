package contracts

import (
	"context"
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
	CreateTag(ctx context.Context,m *models.Tag) error
	GetTagByID(ctx context.Context,id int64) (*models.Tag, error)
	ListTag(ctx context.Context,page, pageSize int, filter *dto.FilterTagRequest) ([]models.Tag, int64, error)
	UpdateTag(ctx context.Context,m *models.Tag) error
	DeleteTag(ctx context.Context,id int64, deletedBy int64) error
}

// TagService defines business logic operations for Tag.
// Diimplementasikan oleh struct 'service' yang sama dengan entitas utama.
type TagService interface {
	CreateTag(ctx context.Context,req *dto.CreateTagRequest, actor he.AuthContext) (*dto.TagResponse, error)
	GetTagByID(ctx context.Context,id int64, actor he.AuthContext) (*dto.TagResponse, error)
	ListTag(ctx context.Context,page, pageSize int, filter *dto.FilterTagRequest, actor he.AuthContext) ([]dto.TagResponse, int64, error)
	UpdateTag(ctx context.Context,id int64, req *dto.UpdateTagRequest, actor he.AuthContext) (*dto.TagResponse, error)
	DeleteTag(ctx context.Context,id int64, actor he.AuthContext) error
}
