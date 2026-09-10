package contracts

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/alamat/dto"
	"neosim_go/internal/modules/kepegawaian/alamat/models"
	he "neosim_go/internal/shared/httputil"
)

// KepegawaianAlamatRepository defines database operations for KepegawaianAlamat.
// Diimplementasikan oleh struct 'repository' (lihat repositories/repository.go
// & repositories/alamat_repository.go).
type KepegawaianAlamatRepository interface {
	CreateAlamat(ctx context.Context,m *models.KepegawaianAlamat) error
	GetAlamatByID(ctx context.Context,id int64) (*models.KepegawaianAlamat, error)
	ListAlamat(ctx context.Context,page, pageSize int, filter *dto.FilterKepegawaianAlamatRequest) ([]models.KepegawaianAlamat, int64, error)
	UpdateAlamat(ctx context.Context,m *models.KepegawaianAlamat) error
	DeleteAlamat(ctx context.Context,id int64, deletedBy int64) error
}

// KepegawaianAlamatService defines business logic operations for KepegawaianAlamat.
// Diimplementasikan oleh struct 'service' (lihat services/service.go
// & services/alamat_service.go).
type KepegawaianAlamatService interface {
	CreateAlamat(ctx context.Context,req *dto.CreateKepegawaianAlamatRequest, actor he.AuthContext) (*dto.KepegawaianAlamatResponse, error)
	GetAlamatByID(ctx context.Context,id int64, actor he.AuthContext) (*dto.KepegawaianAlamatResponse, error)
	ListAlamat(ctx context.Context,page, pageSize int, filter *dto.FilterKepegawaianAlamatRequest, actor he.AuthContext) ([]dto.KepegawaianAlamatResponse, int64, error)
	UpdateAlamat(ctx context.Context,id int64, req *dto.UpdateKepegawaianAlamatRequest, actor he.AuthContext) (*dto.KepegawaianAlamatResponse, error)
	DeleteAlamat(ctx context.Context,id int64, actor he.AuthContext) error
}
