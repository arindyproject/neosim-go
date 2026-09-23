package contracts

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	he "neosim_go/internal/shared/httputil"
)

// JobTitleKategoriRepository defines database operations for JobTitleKategori.
// Diimplementasikan oleh struct 'repository' yang sama dengan entitas utama
// sub-module ini (lihat repositories/repository.go) — TIDAK ADA struct baru.
// Method diberi suffix nama item agar tidak bentrok saat di-embed ke
// contracts.Repository.
type JobTitleKategoriRepository interface {
	CheckJobTitleKategori(ctx context.Context, id int64) (bool, error)
	CreateJobTitleKategori(ctx context.Context, m *models.JobTitleKategori) error
	GetJobTitleKategoriByID(ctx context.Context, id int64) (*models.JobTitleKategori, error)
	GetJobTitleKategoriByCode(ctx context.Context, code string) (*models.JobTitleKategori, error)
	GetJobTitleKategoriByLabel(ctx context.Context, label string) (*models.JobTitleKategori, error)
	ListSelectTipe(ctx context.Context, search string) ([]models.JobTitleKategori, error)
	ListJobTitleKategori(ctx context.Context, page, pageSize int, filter *dto.FilterJobTitleKategoriRequest) ([]models.JobTitleKategori, int64, error)
	UpdateJobTitleKategori(ctx context.Context, m *models.JobTitleKategori) error
	DeleteJobTitleKategori(ctx context.Context, id int64, deletedBy int64) error
}

// JobTitleKategoriService defines business logic operations for JobTitleKategori.
// Diimplementasikan oleh struct 'service' yang sama dengan entitas utama.
type JobTitleKategoriService interface {
	ListSelectJobTitleKategori(ctx context.Context, search string, actor he.AuthContext) ([]dto.JobTitleKategoriSimpelResponse, error)
	CreateJobTitleKategori(ctx context.Context, req *dto.CreateJobTitleKategoriRequest, actor he.AuthContext) (*dto.JobTitleKategoriResponse, error)
	GetJobTitleKategoriByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.JobTitleKategoriResponse, error)
	ListJobTitleKategori(ctx context.Context, page, pageSize int, filter *dto.FilterJobTitleKategoriRequest, actor he.AuthContext) ([]dto.JobTitleKategoriResponse, int64, error)
	UpdateJobTitleKategori(ctx context.Context, id int64, req *dto.UpdateJobTitleKategoriRequest, actor he.AuthContext) (*dto.JobTitleKategoriResponse, error)
	DeleteJobTitleKategori(ctx context.Context, id int64, actor he.AuthContext) error
}
