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
	CreateKualifikasi(ctx context.Context,m *models.KepegawaianKualifikasi) error
	GetKualifikasiByID(ctx context.Context,id int64) (*models.KepegawaianKualifikasi, error)
	ListKualifikasi(ctx context.Context,page, pageSize int, filter *dto.FilterKepegawaianKualifikasiRequest) ([]models.KepegawaianKualifikasi, int64, error)
	UpdateKualifikasi(ctx context.Context,m *models.KepegawaianKualifikasi) error
	DeleteKualifikasi(ctx context.Context,id int64) error
}

// KepegawaianKualifikasiService defines business logic operations for KepegawaianKualifikasi.
// Diimplementasikan oleh struct 'service' (lihat services/service.go
// & services/kualifikasi_service.go).
type KepegawaianKualifikasiService interface {
	CreateKualifikasi(ctx context.Context,req *dto.CreateKepegawaianKualifikasiRequest, actor he.AuthContext) (*dto.KepegawaianKualifikasiResponse, error)
	GetKualifikasiByID(ctx context.Context,id int64, actor he.AuthContext) (*dto.KepegawaianKualifikasiResponse, error)
	ListKualifikasi(ctx context.Context,page, pageSize int, filter *dto.FilterKepegawaianKualifikasiRequest, actor he.AuthContext) ([]dto.KepegawaianKualifikasiResponse, int64, error)
	UpdateKualifikasi(ctx context.Context,id int64, req *dto.UpdateKepegawaianKualifikasiRequest, actor he.AuthContext) (*dto.KepegawaianKualifikasiResponse, error)
	DeleteKualifikasi(ctx context.Context,id int64, actor he.AuthContext) error
}
