package contracts

import (
	"neosim_go/internal/modules/kepegawaian/identifikasi/dto"
	"neosim_go/internal/modules/kepegawaian/identifikasi/models"
	he "neosim_go/internal/shared/httputil"
)

// Repository defines database operations
type Repository interface {
	Create(m *models.KepegawaianIdentifikasi) error
	GetByID(id int64) (*models.KepegawaianIdentifikasi, error)
	List(page, pageSize int, filter *dto.FilterKepegawaianIdentifikasiRequest) ([]models.KepegawaianIdentifikasi, int64, error)
	Update(m *models.KepegawaianIdentifikasi) error
	Delete(id int64) error
}

// Service defines business logic operations
type Service interface {
	Create(req *dto.CreateKepegawaianIdentifikasiRequest, createdBy *int64, actor he.AuthContext) (*dto.KepegawaianIdentifikasiResponse, error)
	GetByID(id int64, actor he.AuthContext) (*dto.KepegawaianIdentifikasiResponse, error)
	List(page, pageSize int, filter *dto.FilterKepegawaianIdentifikasiRequest, actor he.AuthContext) ([]dto.KepegawaianIdentifikasiResponse, int64, error)
	Update(id int64, req *dto.UpdateKepegawaianIdentifikasiRequest, updatedBy *int64, actor he.AuthContext) (*dto.KepegawaianIdentifikasiResponse, error)
	Delete(id int64, actor he.AuthContext) error
}
