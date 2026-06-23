package contracts

import (
	"neosim_go/internal/modules/master/departemen/dto"
	"neosim_go/internal/modules/master/departemen/models"
	he "neosim_go/internal/shared/httputil"
)

// Repository defines database operations
type Repository interface {
	Create(m *models.MasterDepartemen) error
	GetByID(id int64) (*models.MasterDepartemen, error)
	List(page, pageSize int, filter *dto.FilterMasterDepartemenRequest) ([]models.MasterDepartemen, int64, error)
	Update(m *models.MasterDepartemen) error
	Delete(id int64) error
}

// Service defines business logic operations
type Service interface {
	Create(req *dto.CreateMasterDepartemenRequest, createdBy *int64, actor he.AuthContext) (*dto.MasterDepartemenResponse, error)
	GetByID(id int64, actor he.AuthContext) (*dto.MasterDepartemenResponse, error)
	List(page, pageSize int, filter *dto.FilterMasterDepartemenRequest, actor he.AuthContext) ([]dto.MasterDepartemenResponse, int64, error)
	Update(id int64, req *dto.UpdateMasterDepartemenRequest, updatedBy *int64, actor he.AuthContext) (*dto.MasterDepartemenResponse, error)
	Delete(id int64, actor he.AuthContext) error
}
