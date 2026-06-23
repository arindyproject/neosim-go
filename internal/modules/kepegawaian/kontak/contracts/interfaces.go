package contracts

import (
	"neosim_go/internal/modules/kepegawaian/kontak/dto"
	"neosim_go/internal/modules/kepegawaian/kontak/models"
	he "neosim_go/internal/shared/httputil"
)

// Repository defines database operations
type Repository interface {
	Create(m *models.KepegawaianKontak) error
	GetByID(id int64) (*models.KepegawaianKontak, error)
	List(page, pageSize int, filter *dto.FilterKepegawaianKontakRequest) ([]models.KepegawaianKontak, int64, error)
	Update(m *models.KepegawaianKontak) error
	Delete(id int64) error
}

// Service defines business logic operations
type Service interface {
	Create(req *dto.CreateKepegawaianKontakRequest, createdBy *int64, actor he.AuthContext) (*dto.KepegawaianKontakResponse, error)
	GetByID(id int64, actor he.AuthContext) (*dto.KepegawaianKontakResponse, error)
	List(page, pageSize int, filter *dto.FilterKepegawaianKontakRequest, actor he.AuthContext) ([]dto.KepegawaianKontakResponse, int64, error)
	Update(id int64, req *dto.UpdateKepegawaianKontakRequest, updatedBy *int64, actor he.AuthContext) (*dto.KepegawaianKontakResponse, error)
	Delete(id int64, actor he.AuthContext) error
}
