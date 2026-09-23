package contracts

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	he "neosim_go/internal/shared/httputil"
)

// JobTitleRumpunProfesiRepository defines database operations for JobTitleRumpunProfesi.
// Diimplementasikan oleh struct 'repository' yang sama dengan entitas utama
// sub-module ini (lihat repositories/repository.go) — TIDAK ADA struct baru.
// Method diberi suffix nama item agar tidak bentrok saat di-embed ke
// contracts.Repository.
type JobTitleRumpunProfesiRepository interface {
	CheckJobTitleRumpunProfesi(ctx context.Context, id int64) (bool, error)
	CreateJobTitleRumpunProfesi(ctx context.Context, m *models.JobTitleRumpunProfesi) error
	GetJobTitleRumpunProfesiByID(ctx context.Context, id int64) (*models.JobTitleRumpunProfesi, error)
	GetJobTitleRumpunProfesiByCode(ctx context.Context, code string) (*models.JobTitleRumpunProfesi, error)
	GetJobTitleRumpunProfesiByLabel(ctx context.Context, label string) (*models.JobTitleRumpunProfesi, error)
	ListSelectJobTitleRumpunProfesi(ctx context.Context, search string) ([]models.JobTitleRumpunProfesi, error)
	ListJobTitleRumpunProfesi(ctx context.Context, page, pageSize int, filter *dto.FilterJobTitleRumpunProfesiRequest) ([]models.JobTitleRumpunProfesi, int64, error)
	UpdateJobTitleRumpunProfesi(ctx context.Context, m *models.JobTitleRumpunProfesi) error
	DeleteJobTitleRumpunProfesi(ctx context.Context, id int64, deletedBy int64) error
}

// JobTitleRumpunProfesiService defines business logic operations for JobTitleRumpunProfesi.
// Diimplementasikan oleh struct 'service' yang sama dengan entitas utama.
type JobTitleRumpunProfesiService interface {
	ListSelectJobTitleRumpunProfesi(ctx context.Context, search string, actor he.AuthContext) ([]dto.JobTitleRumpunProfesiSimpelResponse, error)
	CreateJobTitleRumpunProfesi(ctx context.Context, req *dto.CreateJobTitleRumpunProfesiRequest, actor he.AuthContext) (*dto.JobTitleRumpunProfesiResponse, error)
	GetJobTitleRumpunProfesiByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.JobTitleRumpunProfesiResponse, error)
	ListJobTitleRumpunProfesi(ctx context.Context, page, pageSize int, filter *dto.FilterJobTitleRumpunProfesiRequest, actor he.AuthContext) ([]dto.JobTitleRumpunProfesiResponse, int64, error)
	UpdateJobTitleRumpunProfesi(ctx context.Context, id int64, req *dto.UpdateJobTitleRumpunProfesiRequest, actor he.AuthContext) (*dto.JobTitleRumpunProfesiResponse, error)
	DeleteJobTitleRumpunProfesi(ctx context.Context, id int64, actor he.AuthContext) error
}
