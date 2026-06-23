package contracts

import (
	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	he "neosim_go/internal/shared/httputil"
)

// Repository defines database operations
type Repository interface {
	Create(m *models.KepegawaianJabatan) error
	GetByID(id int64) (*models.KepegawaianJabatan, error)
	List(page, pageSize int, filter *dto.FilterKepegawaianJabatanRequest) ([]models.KepegawaianJabatan, int64, error)
	Update(m *models.KepegawaianJabatan) error
	Delete(id int64) error
}

// Service defines business logic operations
type Service interface {
	Create(req *dto.CreateKepegawaianJabatanRequest, createdBy *int64, actor he.AuthContext) (*dto.KepegawaianJabatanResponse, error)
	GetByID(id int64, actor he.AuthContext) (*dto.KepegawaianJabatanResponse, error)
	List(page, pageSize int, filter *dto.FilterKepegawaianJabatanRequest, actor he.AuthContext) ([]dto.KepegawaianJabatanResponse, int64, error)
	Update(id int64, req *dto.UpdateKepegawaianJabatanRequest, updatedBy *int64, actor he.AuthContext) (*dto.KepegawaianJabatanResponse, error)
	Delete(id int64, actor he.AuthContext) error
}
