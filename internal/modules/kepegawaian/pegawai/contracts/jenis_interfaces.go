package contracts

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/pegawai/dto"
	"neosim_go/internal/modules/kepegawaian/pegawai/models"
	he "neosim_go/internal/shared/httputil"
)

// JenisRepository defines database operations for Jenis.
// Diimplementasikan oleh struct 'repository' yang sama dengan entitas utama
// sub-module ini (lihat repositories/repository.go) — TIDAK ADA struct baru.
// Method diberi suffix nama item agar tidak bentrok saat di-embed ke
// contracts.Repository.
type JenisRepository interface {
	CreateJenis(ctx context.Context, m *models.Jenis) error
	GetJenisByCode(ctx context.Context, code string) (*models.Jenis, error)
	GetJenisByLabel(ctx context.Context, label string) (*models.Jenis, error)
	ListSelectJenis(ctx context.Context, search string) ([]models.Jenis, error)
	GetJenisByID(ctx context.Context, id int64) (*models.Jenis, error)
	ListJenis(ctx context.Context, page, pageSize int, filter *dto.FilterJenisRequest) ([]models.Jenis, int64, error)
	UpdateJenis(ctx context.Context, m *models.Jenis) error
	DeleteJenis(ctx context.Context, id int64, deletedBy int64) error
}

// JenisService defines business logic operations for Jenis.
// Diimplementasikan oleh struct 'service' yang sama dengan entitas utama.
type JenisService interface {
	CreateJenis(ctx context.Context, req *dto.CreateJenisRequest, actor he.AuthContext) (*dto.JenisResponse, error)
	GetJenisByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.JenisResponse, error)
	GetJenisByCode(ctx context.Context, code string, actor he.AuthContext) (*dto.JenisResponse, error)
	GetJenisByLabel(ctx context.Context, label string, actor he.AuthContext) (*dto.JenisResponse, error)
	ListSelectJenis(ctx context.Context, search string, actor he.AuthContext) ([]dto.JenisSelectResponse, error)
	ListJenis(ctx context.Context, page, pageSize int, filter *dto.FilterJenisRequest, actor he.AuthContext) ([]dto.JenisResponse, int64, error)
	UpdateJenis(ctx context.Context, id int64, req *dto.UpdateJenisRequest, actor he.AuthContext) (*dto.JenisResponse, error)
	DeleteJenis(ctx context.Context, id int64, actor he.AuthContext) error
}
