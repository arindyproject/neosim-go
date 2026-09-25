package contracts

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	he "neosim_go/internal/shared/httputil"
)

// SpecializationKategoriRepository defines database operations for SpecializationKategori.
// Diimplementasikan oleh struct 'repository' yang sama dengan entitas utama
// sub-module ini (lihat repositories/repository.go) — TIDAK ADA struct baru.
// Method diberi suffix nama item agar tidak bentrok saat di-embed ke
// contracts.Repository.
type SpecializationKategoriRepository interface {
	CreateSpecializationKategori(ctx context.Context,m *models.SpecializationKategori) error
	GetSpecializationKategoriByID(ctx context.Context,id int64) (*models.SpecializationKategori, error)
	ListSpecializationKategori(ctx context.Context,page, pageSize int, filter *dto.FilterSpecializationKategoriRequest) ([]models.SpecializationKategori, int64, error)
	UpdateSpecializationKategori(ctx context.Context,m *models.SpecializationKategori) error
	DeleteSpecializationKategori(ctx context.Context,id int64, deletedBy int64) error
}

// SpecializationKategoriService defines business logic operations for SpecializationKategori.
// Diimplementasikan oleh struct 'service' yang sama dengan entitas utama.
type SpecializationKategoriService interface {
	CreateSpecializationKategori(ctx context.Context,req *dto.CreateSpecializationKategoriRequest, actor he.AuthContext) (*dto.SpecializationKategoriResponse, error)
	GetSpecializationKategoriByID(ctx context.Context,id int64, actor he.AuthContext) (*dto.SpecializationKategoriResponse, error)
	ListSpecializationKategori(ctx context.Context,page, pageSize int, filter *dto.FilterSpecializationKategoriRequest, actor he.AuthContext) ([]dto.SpecializationKategoriResponse, int64, error)
	UpdateSpecializationKategori(ctx context.Context,id int64, req *dto.UpdateSpecializationKategoriRequest, actor he.AuthContext) (*dto.SpecializationKategoriResponse, error)
	DeleteSpecializationKategori(ctx context.Context,id int64, actor he.AuthContext) error
}
