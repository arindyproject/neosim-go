package contracts

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/pegawai/dto"
	"neosim_go/internal/modules/kepegawaian/pegawai/models"
	he "neosim_go/internal/shared/httputil"
)

// StatusRepository defines database operations for Status.
// Diimplementasikan oleh struct 'repository' yang sama dengan entitas utama
// sub-module ini (lihat repositories/repository.go) — TIDAK ADA struct baru.
// Method diberi suffix nama item agar tidak bentrok saat di-embed ke
// contracts.Repository.
type StatusRepository interface {
	CreateStatus(ctx context.Context, m *models.Status) error
	GetStatusByID(ctx context.Context, id int64) (*models.Status, error)
	GetStatusByCode(ctx context.Context, code string) (*models.Status, error)
	GetStatusByLabel(ctx context.Context, label string) (*models.Status, error)
	ListSelectStatus(ctx context.Context, search string) ([]models.Status, error)
	ListStatus(ctx context.Context, page, pageSize int, filter *dto.FilterStatusRequest) ([]models.Status, int64, error)
	UpdateStatus(ctx context.Context, m *models.Status) error
	DeleteStatus(ctx context.Context, id int64, deletedBy int64) error
}

// StatusService defines business logic operations for Status.
// Diimplementasikan oleh struct 'service' yang sama dengan entitas utama.
type StatusService interface {
	CreateStatus(ctx context.Context, req *dto.CreateStatusRequest, actor he.AuthContext) (*dto.StatusResponse, error)
	GetStatusByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.StatusResponse, error)
	GetStatusByCode(ctx context.Context, code string, actor he.AuthContext) (*dto.StatusResponse, error)
	GetStatusByLabel(ctx context.Context, label string, actor he.AuthContext) (*dto.StatusResponse, error)
	ListSelectStatus(ctx context.Context, search string, actor he.AuthContext) ([]dto.StatusSelectResponse, error)
	ListStatus(ctx context.Context, page, pageSize int, filter *dto.FilterStatusRequest, actor he.AuthContext) ([]dto.StatusResponse, int64, error)
	UpdateStatus(ctx context.Context, id int64, req *dto.UpdateStatusRequest, actor he.AuthContext) (*dto.StatusResponse, error)
	DeleteStatus(ctx context.Context, id int64, actor he.AuthContext) error
}
