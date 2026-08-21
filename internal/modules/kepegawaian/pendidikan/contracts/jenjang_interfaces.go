package contracts

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/pendidikan/dto"
	"neosim_go/internal/modules/kepegawaian/pendidikan/models"
	he "neosim_go/internal/shared/httputil"
)

// JenjangRepository defines database operations for Jenjang.
// Diimplementasikan oleh struct 'repository' yang sama dengan entitas utama
// sub-module ini (lihat repositories/repository.go) — TIDAK ADA struct baru.
// Method diberi suffix nama item agar tidak bentrok saat di-embed ke
// contracts.Repository.
type JenjangRepository interface {
	CreateJenjang(ctx context.Context,m *models.Jenjang) error
	GetJenjangByID(ctx context.Context,id int64) (*models.Jenjang, error)
	ListJenjang(ctx context.Context,page, pageSize int, filter *dto.FilterJenjangRequest) ([]models.Jenjang, int64, error)
	UpdateJenjang(ctx context.Context,m *models.Jenjang) error
	DeleteJenjang(ctx context.Context,id int64) error
}

// JenjangService defines business logic operations for Jenjang.
// Diimplementasikan oleh struct 'service' yang sama dengan entitas utama.
type JenjangService interface {
	CreateJenjang(ctx context.Context,req *dto.CreateJenjangRequest, actor he.AuthContext) (*dto.JenjangResponse, error)
	GetJenjangByID(ctx context.Context,id int64, actor he.AuthContext) (*dto.JenjangResponse, error)
	ListJenjang(ctx context.Context,page, pageSize int, filter *dto.FilterJenjangRequest, actor he.AuthContext) ([]dto.JenjangResponse, int64, error)
	UpdateJenjang(ctx context.Context,id int64, req *dto.UpdateJenjangRequest, actor he.AuthContext) (*dto.JenjangResponse, error)
	DeleteJenjang(ctx context.Context,id int64, actor he.AuthContext) error
}
