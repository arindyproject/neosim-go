package contracts

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	he "neosim_go/internal/shared/httputil"
)

// KepegawaianJabatanRepository defines database operations for KepegawaianJabatan.
// Diimplementasikan oleh struct 'repository' (lihat repositories/repository.go
// & repositories/jabatan_repository.go).
type KepegawaianJabatanRepository interface {
	CreateJabatan(ctx context.Context,m *models.KepegawaianJabatan) error
	GetJabatanByID(ctx context.Context,id int64) (*models.KepegawaianJabatan, error)
	ListJabatan(ctx context.Context,page, pageSize int, filter *dto.FilterKepegawaianJabatanRequest) ([]models.KepegawaianJabatan, int64, error)
	UpdateJabatan(ctx context.Context,m *models.KepegawaianJabatan) error
	DeleteJabatan(ctx context.Context,id int64, deletedBy int64) error
}

// KepegawaianJabatanService defines business logic operations for KepegawaianJabatan.
// Diimplementasikan oleh struct 'service' (lihat services/service.go
// & services/jabatan_service.go).
type KepegawaianJabatanService interface {
	CreateJabatan(ctx context.Context,req *dto.CreateKepegawaianJabatanRequest, actor he.AuthContext) (*dto.KepegawaianJabatanResponse, error)
	GetJabatanByID(ctx context.Context,id int64, actor he.AuthContext) (*dto.KepegawaianJabatanResponse, error)
	ListJabatan(ctx context.Context,page, pageSize int, filter *dto.FilterKepegawaianJabatanRequest, actor he.AuthContext) ([]dto.KepegawaianJabatanResponse, int64, error)
	UpdateJabatan(ctx context.Context,id int64, req *dto.UpdateKepegawaianJabatanRequest, actor he.AuthContext) (*dto.KepegawaianJabatanResponse, error)
	DeleteJabatan(ctx context.Context,id int64, actor he.AuthContext) error
}
