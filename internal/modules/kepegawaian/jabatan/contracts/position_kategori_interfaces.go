package contracts

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	he "neosim_go/internal/shared/httputil"
)

// PositionKategoriRepository defines database operations for PositionKategori.
// Diimplementasikan oleh struct 'repository' yang sama dengan entitas utama
// sub-module ini (lihat repositories/repository.go) — TIDAK ADA struct baru.
// Method diberi suffix nama item agar tidak bentrok saat di-embed ke
// contracts.Repository.
type PositionKategoriRepository interface {
	CreatePositionKategori(ctx context.Context, m *models.PositionKategori) error
	GetPositionKategoriByID(ctx context.Context, id int64) (*models.PositionKategori, error)
	GetPositionKategoriByCode(ctx context.Context, code string) (*models.PositionKategori, error)
	GetPositionKategoriByLabel(ctx context.Context, label string) (*models.PositionKategori, error)
	ListSelectPositionKategori(ctx context.Context, search string) ([]models.PositionKategori, error)
	ListPositionKategori(ctx context.Context, page, pageSize int, filter *dto.FilterPositionKategoriRequest) ([]models.PositionKategori, int64, error)
	UpdatePositionKategori(ctx context.Context, m *models.PositionKategori) error
	DeletePositionKategori(ctx context.Context, id int64, deletedBy int64) error
}

// PositionKategoriService defines business logic operations for PositionKategori.
// Diimplementasikan oleh struct 'service' yang sama dengan entitas utama.
type PositionKategoriService interface {
	CreatePositionKategori(ctx context.Context, req *dto.CreatePositionKategoriRequest, actor he.AuthContext) (*dto.PositionKategoriResponse, error)
	GetPositionKategoriByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.PositionKategoriResponse, error)
	GetPositionKategoriByCode(ctx context.Context, code string, actor he.AuthContext) (*dto.PositionKategoriResponse, error)
	GetPositionKategoriByLabel(ctx context.Context, label string, actor he.AuthContext) (*dto.PositionKategoriResponse, error)
	ListSelectPositionKategori(ctx context.Context, search string, actor he.AuthContext) ([]dto.PositionKategoriSelectResponse, error)
	ListPositionKategori(ctx context.Context, page, pageSize int, filter *dto.FilterPositionKategoriRequest, actor he.AuthContext) ([]dto.PositionKategoriResponse, int64, error)
	UpdatePositionKategori(ctx context.Context, id int64, req *dto.UpdatePositionKategoriRequest, actor he.AuthContext) (*dto.PositionKategoriResponse, error)
	DeletePositionKategori(ctx context.Context, id int64, actor he.AuthContext) error
}
