package contracts

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/pendidikan/dto"
	"neosim_go/internal/modules/kepegawaian/pendidikan/models"
	he "neosim_go/internal/shared/httputil"
)

// KepegawaianPendidikanRepository defines database operations for KepegawaianPendidikan.
// Diimplementasikan oleh struct 'repository' (lihat repositories/repository.go
// & repositories/pendidikan_repository.go).
type KepegawaianPendidikanRepository interface {
	CreatePendidikan(ctx context.Context, m *models.KepegawaianPendidikan) error
	GetPendidikanByID(ctx context.Context, id int64) (*models.KepegawaianPendidikan, error)
	GetPendidikanByPegawaiID(ctx context.Context, pegawaiID int64, page, pageSize int) ([]models.KepegawaianPendidikan, int64, error)
	GetByPegawaiIDAndTipe(ctx context.Context, pegawaiID, jenjangID int64) ([]models.KepegawaianPendidikan, error)
	ListPendidikan(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianPendidikanRequest) ([]models.KepegawaianPendidikan, int64, error)
	UpdatePendidikan(ctx context.Context, m *models.KepegawaianPendidikan) error
	DeletePendidikan(ctx context.Context, id int64) error

	ExistsPendidikanByID(ctx context.Context, id int64) (bool, error)
	ExistsByNomorIjazah(ctx context.Context, jenjangID int64, nomorIjazah string, excludeID int64) (bool, error)
	ExistsByNomorIjazahOnly(ctx context.Context, nomorIjazah string, excludeID int64) (bool, error)
}

// KepegawaianPendidikanService defines business logic operations for KepegawaianPendidikan.
// Diimplementasikan oleh struct 'service' (lihat services/service.go
// & services/pendidikan_service.go).
type KepegawaianPendidikanService interface {
	CreatePendidikan(ctx context.Context, req *dto.CreateKepegawaianPendidikanRequest, actor he.AuthContext) (*dto.KepegawaianPendidikanResponse, error)
	GetPendidikanByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.KepegawaianPendidikanResponse, error)
	ListPendidikan(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianPendidikanRequest, actor he.AuthContext) ([]dto.KepegawaianPendidikanResponse, int64, error)
	ListPendidikanByPegawai(ctx context.Context, pegawaiID int64, page, pageSize int, actor he.AuthContext) ([]dto.KepegawaianPendidikanResponse, int64, error)
	UpdatePendidikan(ctx context.Context, id int64, req *dto.UpdateKepegawaianPendidikanRequest, actor he.AuthContext) (*dto.KepegawaianPendidikanResponse, error)
	DeletePendidikan(ctx context.Context, id int64, actor he.AuthContext) error
}
