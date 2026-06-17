package contracts

import (
	"neosim_go/internal/modules/master/master/dto"
	"neosim_go/internal/modules/master/master/models"
)

// AuthContext berisi informasi user yang sedang login untuk authorization
type AuthContext struct {
	UserID       int64
	IsSuperadmin bool
}

// Repository defines database operations
type Repository interface {
	Create(m *models.Master) error
	GetByID(id int64) (*models.Master, error)
	List(page, pageSize int, filter *dto.FilterMasterRequest) ([]models.Master, int64, error)
	Update(m *models.Master) error
	Delete(id int64) error
}

// Service defines business logic operations
type Service interface {
	Create(req *dto.CreateMasterRequest, createdBy *int64, actor AuthContext) (*dto.MasterResponse, error)
	GetByID(id int64 , actor AuthContext) (*dto.MasterResponse, error)
	List(page, pageSize int,filter *dto.FilterMasterRequest, actor AuthContext) ([]dto.MasterResponse, int64, error)
	Update(id int64, req *dto.UpdateMasterRequest, updatedBy *int64, actor AuthContext) (*dto.MasterResponse, error)
	Delete(id int64, actor AuthContext) error
}
