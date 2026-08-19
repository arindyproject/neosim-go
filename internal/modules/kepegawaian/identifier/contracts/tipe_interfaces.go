package contracts

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/identifier/dto"
	"neosim_go/internal/modules/kepegawaian/identifier/models"
	he "neosim_go/internal/shared/httputil"
)

// TipeRepository defines database operations for Tipe.
// Diimplementasikan oleh struct 'repository' yang sama dengan entitas utama
// sub-module ini (lihat repositories/repository.go) — TIDAK ADA struct baru.
// Method diberi suffix nama item agar tidak bentrok saat di-embed ke
// contracts.Repository.
type TipeRepository interface {
	CreateTipe(ctx context.Context, m *models.Tipe) error
	GetTipeByID(ctx context.Context, id int64) (*models.Tipe, error)
	GetTipeByCode(ctx context.Context, code string) (*models.Tipe, error)
	GetTipeByLabel(ctx context.Context, label string) (*models.Tipe, error)
	ListTipe(ctx context.Context, page, pageSize int, filter *dto.FilterTipeRequest) ([]models.Tipe, int64, error)
	UpdateTipe(ctx context.Context, m *models.Tipe) error
	DeleteTipe(ctx context.Context, id int64) error
}

// TipeService defines business logic operations for Tipe.
// Diimplementasikan oleh struct 'service' yang sama dengan entitas utama.
type TipeService interface {
	CreateTipe(ctx context.Context, req *dto.CreateTipeRequest, actor he.AuthContext) (*dto.TipeResponse, error)
	GetTipeByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.TipeResponse, error)
	GetTipeByCode(ctx context.Context, code string, actor he.AuthContext) (*dto.TipeResponse, error)
	GetTipeByLabel(ctx context.Context, label string, actor he.AuthContext) (*dto.TipeResponse, error)
	ListTipe(ctx context.Context, page, pageSize int, filter *dto.FilterTipeRequest, actor he.AuthContext) ([]dto.TipeResponse, int64, error)
	UpdateTipe(ctx context.Context, id int64, req *dto.UpdateTipeRequest, actor he.AuthContext) (*dto.TipeResponse, error)
	DeleteTipe(ctx context.Context, id int64, actor he.AuthContext) error
}
