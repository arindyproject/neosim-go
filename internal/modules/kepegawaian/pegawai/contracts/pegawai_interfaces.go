package contracts

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/pegawai/dto"
	"neosim_go/internal/modules/kepegawaian/pegawai/models"
	he "neosim_go/internal/shared/httputil"
)

// KepegawaianPegawaiRepository defines database operations for KepegawaianPegawai.
// Diimplementasikan oleh struct 'repository' (lihat repositories/repository.go
// & repositories/pegawai_repository.go).
type KepegawaianPegawaiRepository interface {
	CreatePegawai(ctx context.Context, m *models.KepegawaianPegawai) error
	GetPegawaiByID(ctx context.Context, id int64) (*models.KepegawaianPegawai, error)
	ListPegawai(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianPegawaiRequest) ([]models.KepegawaianPegawai, int64, error)
	UpdatePegawai(ctx context.Context, m *models.KepegawaianPegawai) error
	DeletePegawai(ctx context.Context, id int64, deletedBy int64) error
}

// KepegawaianPegawaiService defines business logic operations for KepegawaianPegawai.
// Diimplementasikan oleh struct 'service' (lihat services/service.go
// & services/pegawai_service.go).
type KepegawaianPegawaiService interface {
	CreatePegawai(ctx context.Context, req *dto.CreateKepegawaianPegawaiRequest, actor he.AuthContext) (*dto.KepegawaianPegawaiResponse, error)
	GetPegawaiByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.KepegawaianPegawaiResponse, error)
	ListPegawai(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianPegawaiRequest, actor he.AuthContext) ([]dto.KepegawaianPegawaiResponse, int64, error)
	UpdatePegawai(ctx context.Context, id int64, req *dto.UpdateKepegawaianPegawaiRequest, actor he.AuthContext) (*dto.KepegawaianPegawaiResponse, error)
	DeletePegawai(ctx context.Context, id int64, actor he.AuthContext) error
}
