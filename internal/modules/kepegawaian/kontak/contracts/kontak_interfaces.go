package contracts

import (
	"neosim_go/internal/modules/kepegawaian/kontak/dto"
	"neosim_go/internal/modules/kepegawaian/kontak/models"
	he "neosim_go/internal/shared/httputil"
)

// KepegawaianKontakRepository defines database operations for KepegawaianKontak.
// Diimplementasikan oleh struct 'repository' (lihat repositories/repository.go
// & repositories/kontak_repository.go).
type KepegawaianKontakRepository interface {
	CreateKontak(m *models.KepegawaianKontak) error
	GetKontakByID(id int64) (*models.KepegawaianKontak, error)
	ListKontak(page, pageSize int, filter *dto.FilterKepegawaianKontakRequest) ([]models.KepegawaianKontak, int64, error)
	UpdateKontak(m *models.KepegawaianKontak) error
	DeleteKontak(id int64) error
}

// KepegawaianKontakService defines business logic operations for KepegawaianKontak.
// Diimplementasikan oleh struct 'service' (lihat services/service.go
// & services/kontak_service.go).
type KepegawaianKontakService interface {
	CreateKontak(req *dto.CreateKepegawaianKontakRequest, actor he.AuthContext) (*dto.KepegawaianKontakResponse, error)
	GetKontakByID(id int64, actor he.AuthContext) (*dto.KepegawaianKontakResponse, error)
	ListKontak(page, pageSize int, filter *dto.FilterKepegawaianKontakRequest, actor he.AuthContext) ([]dto.KepegawaianKontakResponse, int64, error)
	UpdateKontak(id int64, req *dto.UpdateKepegawaianKontakRequest, actor he.AuthContext) (*dto.KepegawaianKontakResponse, error)
	DeleteKontak(id int64, actor he.AuthContext) error
}
