package contracts

import (
	"neosim_go/internal/modules/kepegawaian/kualifikasi/dto"
	"neosim_go/internal/modules/kepegawaian/kualifikasi/models"
	he "neosim_go/internal/shared/httputil"
)

// Repository defines database operations
type Repository interface {
	Create(m *models.KepegawaianKualifikasi) error
	GetByID(id int64) (*models.KepegawaianKualifikasi, error)
	List(page, pageSize int, filter *dto.FilterKepegawaianKualifikasiRequest) ([]models.KepegawaianKualifikasi, int64, error)
	Update(m *models.KepegawaianKualifikasi) error
	Delete(id int64) error
}

// Service defines business logic operations
type Service interface {
	Create(req *dto.CreateKepegawaianKualifikasiRequest, createdBy *int64, actor he.AuthContext) (*dto.KepegawaianKualifikasiResponse, error)
	GetByID(id int64, actor he.AuthContext) (*dto.KepegawaianKualifikasiResponse, error)
	List(page, pageSize int, filter *dto.FilterKepegawaianKualifikasiRequest, actor he.AuthContext) ([]dto.KepegawaianKualifikasiResponse, int64, error)
	Update(id int64, req *dto.UpdateKepegawaianKualifikasiRequest, updatedBy *int64, actor he.AuthContext) (*dto.KepegawaianKualifikasiResponse, error)
	Delete(id int64, actor he.AuthContext) error
}
