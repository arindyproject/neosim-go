package contracts

import (
	"neosim_go/internal/modules/kepegawaian/lampiran/dto"
	"neosim_go/internal/modules/kepegawaian/lampiran/models"
	he "neosim_go/internal/shared/httputil"
)

// Repository defines database operations
type Repository interface {
	Create(m *models.KepegawaianLampiran) error
	GetByID(id int64) (*models.KepegawaianLampiran, error)
	List(page, pageSize int, filter *dto.FilterKepegawaianLampiranRequest) ([]models.KepegawaianLampiran, int64, error)
	Update(m *models.KepegawaianLampiran) error
	Delete(id int64) error
}

// Service defines business logic operations
type Service interface {
	Create(req *dto.CreateKepegawaianLampiranRequest, createdBy *int64, actor he.AuthContext) (*dto.KepegawaianLampiranResponse, error)
	GetByID(id int64, actor he.AuthContext) (*dto.KepegawaianLampiranResponse, error)
	List(page, pageSize int, filter *dto.FilterKepegawaianLampiranRequest, actor he.AuthContext) ([]dto.KepegawaianLampiranResponse, int64, error)
	Update(id int64, req *dto.UpdateKepegawaianLampiranRequest, updatedBy *int64, actor he.AuthContext) (*dto.KepegawaianLampiranResponse, error)
	Delete(id int64, actor he.AuthContext) error
}
