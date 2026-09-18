package contracts

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	he "neosim_go/internal/shared/httputil"
)

// PositionRepository defines database operations for Position.
// Diimplementasikan oleh struct 'repository' yang sama dengan entitas utama
// sub-module ini (lihat repositories/repository.go) — TIDAK ADA struct baru.
// Method diberi suffix nama item agar tidak bentrok saat di-embed ke
// contracts.Repository.
type PositionRepository interface {
	CreatePosition(ctx context.Context,m *models.Position) error
	GetPositionByID(ctx context.Context,id int64) (*models.Position, error)
	ListPosition(ctx context.Context,page, pageSize int, filter *dto.FilterPositionRequest) ([]models.Position, int64, error)
	UpdatePosition(ctx context.Context,m *models.Position) error
	DeletePosition(ctx context.Context,id int64, deletedBy int64) error
}

// PositionService defines business logic operations for Position.
// Diimplementasikan oleh struct 'service' yang sama dengan entitas utama.
type PositionService interface {
	CreatePosition(ctx context.Context,req *dto.CreatePositionRequest, actor he.AuthContext) (*dto.PositionResponse, error)
	GetPositionByID(ctx context.Context,id int64, actor he.AuthContext) (*dto.PositionResponse, error)
	ListPosition(ctx context.Context,page, pageSize int, filter *dto.FilterPositionRequest, actor he.AuthContext) ([]dto.PositionResponse, int64, error)
	UpdatePosition(ctx context.Context,id int64, req *dto.UpdatePositionRequest, actor he.AuthContext) (*dto.PositionResponse, error)
	DeletePosition(ctx context.Context,id int64, actor he.AuthContext) error
}
