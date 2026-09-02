package contracts

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/kualifikasi/dto"
	"neosim_go/internal/modules/kepegawaian/kualifikasi/models"
	he "neosim_go/internal/shared/httputil"
)

// KepegawaianKualifikasiRepository defines database operations for KepegawaianKualifikasi.
// Diimplementasikan oleh struct 'repository' (lihat repositories/repository.go
// & repositories/kualifikasi_repository.go).
type KepegawaianKualifikasiRepository interface {
	CreateKualifikasi(ctx context.Context, m *models.KepegawaianKualifikasi) error
	GetKualifikasiByID(ctx context.Context, id int64) (*models.KepegawaianKualifikasi, error)
	GetKualifikasiByPegawaiID(ctx context.Context, pegawaiID int64, page, pageSize int) ([]models.KepegawaianKualifikasi, int64, error)
	GetKualifikasiByPegawaiIDAndTipe(ctx context.Context, pegawaiID, tipeID int64) ([]models.KepegawaianKualifikasi, error)
	GetKualifikasiByTipe(ctx context.Context, tipeID int64, page, pageSize int) ([]models.KepegawaianKualifikasi, int64, error)
	GetExpiringSoonKualifikasi(ctx context.Context, days int, page, pageSize int) ([]models.KepegawaianKualifikasi, int64, error)
	GetExpiredKualifikasi(ctx context.Context, page, pageSize int) ([]models.KepegawaianKualifikasi, int64, error)
	ListKualifikasi(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianKualifikasiRequest) ([]models.KepegawaianKualifikasi, int64, error)
	UpdateKualifikasi(ctx context.Context, m *models.KepegawaianKualifikasi) error
	DeleteKualifikasi(ctx context.Context, id int64, deletedBy int64) error
	ExistsByNomorSertifikatAndTipe(ctx context.Context, tipeID int64, NomorSertifikat string, excludeID int64) (bool, error)
}

// KepegawaianKualifikasiService defines business logic operations for KepegawaianKualifikasi.
// Diimplementasikan oleh struct 'service' (lihat services/service.go
// & services/kualifikasi_service.go).
type KepegawaianKualifikasiService interface {
	CreateKualifikasi(ctx context.Context, req *dto.CreateKepegawaianKualifikasiRequest, actor he.AuthContext) (*dto.KepegawaianKualifikasiResponse, error)
	GetKualifikasiByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.KepegawaianKualifikasiResponse, error)
	ListByPegawai(ctx context.Context, pegawaiID int64, page, pageSize int, actor he.AuthContext) ([]dto.KepegawaianKualifikasiResponse, int64, error)
	ListKualifikasi(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianKualifikasiRequest, actor he.AuthContext) ([]dto.KepegawaianKualifikasiResponse, int64, error)
	UpdateKualifikasi(ctx context.Context, id int64, req *dto.UpdateKepegawaianKualifikasiRequest, actor he.AuthContext) (*dto.KepegawaianKualifikasiResponse, error)
	DeleteKualifikasi(ctx context.Context, id int64, actor he.AuthContext) error
	GetExpiringSoonIdentifier(ctx context.Context, days int, page, pageSize int, actor he.AuthContext) ([]dto.KepegawaianKualifikasiResponse, int64, error)
	GetExpiredKualifikasi(ctx context.Context, page, pageSize int, actor he.AuthContext) ([]dto.KepegawaianKualifikasiResponse, int64, error)
}
