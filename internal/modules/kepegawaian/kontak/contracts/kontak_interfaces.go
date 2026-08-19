package contracts

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/kontak/dto"
	"neosim_go/internal/modules/kepegawaian/kontak/models"
	he "neosim_go/internal/shared/httputil"
)

// KepegawaianKontakRepository defines database operations for KepegawaianKontak.
// Diimplementasikan oleh struct 'repository' (lihat repositories/repository.go
// & repositories/kontak_repository.go).
type KepegawaianKontakRepository interface {
	CreateKontak(ctx context.Context, m *models.KepegawaianKontak) error
	GetKontakByID(ctx context.Context, id int64) (*models.KepegawaianKontak, error)
	GetKontakByPegawaiID(ctx context.Context, pegawaiID int64) ([]models.KepegawaianKontak, error)
	ListKontak(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianKontakRequest) ([]models.KepegawaianKontak, int64, error)
	UpdateKontak(ctx context.Context, m *models.KepegawaianKontak) error
	DeleteKontak(ctx context.Context, id int64) error
}

// KepegawaianKontakService defines business logic operations for KepegawaianKontak.
// Diimplementasikan oleh struct 'service' (lihat services/service.go
// & services/kontak_service.go).
type KepegawaianKontakService interface {
	CreateKontak(ctx context.Context, req *dto.CreateKepegawaianKontakRequest, actor he.AuthContext) (*dto.KepegawaianKontakResponse, error)
	GetKontakByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.KepegawaianKontakResponse, error)
	GetKontakByPegawaiID(ctx context.Context, pegawaiID int64, actor he.AuthContext) ([]dto.KepegawaianKontakResponse, error)
	ListKontak(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianKontakRequest, actor he.AuthContext) ([]dto.KepegawaianKontakResponse, int64, error)
	UpdateKontak(ctx context.Context, id int64, req *dto.UpdateKepegawaianKontakRequest, actor he.AuthContext) (*dto.KepegawaianKontakResponse, error)
	DeleteKontak(ctx context.Context, id int64, actor he.AuthContext) error
}
