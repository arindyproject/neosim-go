package contracts

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	he "neosim_go/internal/shared/httputil"
)

// SpecializationRepository defines database operations for Specialization.
// Diimplementasikan oleh struct 'repository' yang sama dengan entitas utama
// sub-module ini (lihat repositories/repository.go) — TIDAK ADA struct baru.
// Method diberi suffix nama item agar tidak bentrok saat di-embed ke
// contracts.Repository.
type SpecializationRepository interface {
	CreateSpecialization(ctx context.Context,m *models.Specialization) error
	GetSpecializationByID(ctx context.Context,id int64) (*models.Specialization, error)
	ListSpecialization(ctx context.Context,page, pageSize int, filter *dto.FilterSpecializationRequest) ([]models.Specialization, int64, error)
	UpdateSpecialization(ctx context.Context,m *models.Specialization) error
	DeleteSpecialization(ctx context.Context,id int64, deletedBy int64) error
}

// SpecializationService defines business logic operations for Specialization.
// Diimplementasikan oleh struct 'service' yang sama dengan entitas utama.
type SpecializationService interface {
	CreateSpecialization(ctx context.Context,req *dto.CreateSpecializationRequest, actor he.AuthContext) (*dto.SpecializationResponse, error)
	GetSpecializationByID(ctx context.Context,id int64, actor he.AuthContext) (*dto.SpecializationResponse, error)
	ListSpecialization(ctx context.Context,page, pageSize int, filter *dto.FilterSpecializationRequest, actor he.AuthContext) ([]dto.SpecializationResponse, int64, error)
	UpdateSpecialization(ctx context.Context,id int64, req *dto.UpdateSpecializationRequest, actor he.AuthContext) (*dto.SpecializationResponse, error)
	DeleteSpecialization(ctx context.Context,id int64, actor he.AuthContext) error
}
