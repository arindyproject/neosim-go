package contracts

import (
	"neosim_go/internal/modules/kepegawaian/pegawai/dto"
	"neosim_go/internal/modules/kepegawaian/pegawai/models"
	he "neosim_go/internal/shared/httputil"
)

// Repository defines database operations
type Repository interface {
	Create(m *models.KepegawaianPegawai) error
	GetByID(id int64) (*models.KepegawaianPegawai, error)
	List(page, pageSize int, filter *dto.FilterKepegawaianPegawaiRequest) ([]models.KepegawaianPegawai, int64, error)
	Update(m *models.KepegawaianPegawai) error
	Delete(id int64) error
}

// Service defines business logic operations
type Service interface {
	Create(req *dto.CreateKepegawaianPegawaiRequest, createdBy *int64, actor he.AuthContext) (*dto.KepegawaianPegawaiResponse, error)
	GetByID(id int64, actor he.AuthContext) (*dto.KepegawaianPegawaiResponse, error)
	List(page, pageSize int, filter *dto.FilterKepegawaianPegawaiRequest, actor he.AuthContext) ([]dto.KepegawaianPegawaiResponse, int64, error)
	Update(id int64, req *dto.UpdateKepegawaianPegawaiRequest, updatedBy *int64, actor he.AuthContext) (*dto.KepegawaianPegawaiResponse, error)
	Delete(id int64, actor he.AuthContext) error
}
