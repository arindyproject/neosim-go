package contracts

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	he "neosim_go/internal/shared/httputil"
)

// JobTitleRepository defines database operations for JobTitle.
// Diimplementasikan oleh struct 'repository' yang sama dengan entitas utama
// sub-module ini (lihat repositories/repository.go) — TIDAK ADA struct baru.
// Method diberi suffix nama item agar tidak bentrok saat di-embed ke
// contracts.Repository.
type JobTitleRepository interface {
	// Create
	CreateJobTitle(ctx context.Context, m *models.JobTitle) error

	// Update
	UpdateJobTitle(ctx context.Context, m *models.JobTitle) error

	// Delete (soft delete)
	DeleteJobTitle(ctx context.Context, id int64, deletedBy int64) error

	// Find
	GetJobTitleByID(ctx context.Context, id int64) (*models.JobTitle, error)
	ListJobTitle(ctx context.Context, page, pageSize int, filter *dto.FilterJobTitleRequest) ([]models.JobTitle, int64, error)
	ListSelectJobTitle(ctx context.Context, search string) ([]models.JobTitle, error)

	// Exists
	ExistsJobTitleByID(ctx context.Context, id int64) (bool, error)
	ExistsByCode(ctx context.Context, code string, excludeID int64) (bool, error)
}

// JobTitleService defines business logic operations for JobTitle.
// Diimplementasikan oleh struct 'service' yang sama dengan entitas utama.
type JobTitleService interface {
	CreateJobTitle(ctx context.Context, req *dto.CreateJobTitleRequest, actor he.AuthContext) (*dto.JobTitleResponse, error)
	GetJobTitleByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.JobTitleResponse, error)
	ListSelectJobTitle(ctx context.Context, search string, actor he.AuthContext) ([]dto.JobTitleSimpelResponse, error)
	ListJobTitle(ctx context.Context, page, pageSize int, filter *dto.FilterJobTitleRequest, actor he.AuthContext) ([]dto.JobTitleResponse, int64, error)
	UpdateJobTitle(ctx context.Context, id int64, req *dto.UpdateJobTitleRequest, actor he.AuthContext) (*dto.JobTitleResponse, error)
	DeleteJobTitle(ctx context.Context, id int64, actor he.AuthContext) error
}
