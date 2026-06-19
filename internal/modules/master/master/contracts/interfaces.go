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
	// Pekerjaan----------------------------------------------------
	CreatePekerjaan(m *models.MasterPekerjaan) error
	GetByIDPekerjaan(id int64) (*models.MasterPekerjaan, error)
	GetByNamePekerjaan(name string) (*models.MasterPekerjaan, error)
	ListPekerjaan(page, pageSize int, filter *dto.FilterMasterPekerjaanRequest) ([]models.MasterPekerjaan, int64, error)
	UpdatePekerjaan(m *models.MasterPekerjaan) error
	DeletePekerjaan(id int64) error
	//--------------------------------------------------------------

	// Pendidikan---------------------------------------------------
	CreatePendidikan(m *models.MasterPendidikan) error
	GetByIDPendidikan(id int64) (*models.MasterPendidikan, error)
	GetByNamePendidikan(name string) (*models.MasterPendidikan, error)
	ListPendidikan(page, pageSize int, filter *dto.FilterMasterPendidikanRequest) ([]models.MasterPendidikan, int64, error)
	UpdatePendidikan(m *models.MasterPendidikan) error
	DeletePendidikan(id int64) error
	//---------------------------------------------------------------

	// Agama---------------------------------------------------------
	CreateAgama(m *models.MasterAgama) error
	GetByIDAgama(id int64) (*models.MasterAgama, error)
	GetByNameAgama(name string) (*models.MasterAgama, error)
	ListAgama(page, pageSize int, filter *dto.FilterMasterAgamaRequest) ([]models.MasterAgama, int64, error)
	UpdateAgama(m *models.MasterAgama) error
	DeleteAgama(id int64) error
	//----------------------------------------------------------------

	// StatusPernikahan-----------------------------------------------
	CreateStatusPernikahan(m *models.MasterStatusPernikahan) error
	GetByIDStatusPernikahan(id int64) (*models.MasterStatusPernikahan, error)
	GetByNameStatusPernikahan(name string) (*models.MasterStatusPernikahan, error)
	ListStatusPernikahan(page, pageSize int, filter *dto.FilterMasterStatusPernikahanRequest) ([]models.MasterStatusPernikahan, int64, error)
	UpdateStatusPernikahan(m *models.MasterStatusPernikahan) error
	DeleteStatusPernikahan(id int64) error
	//----------------------------------------------------------------
}

// Service defines business logic operations
type Service interface {
	//Create(req *dto.CreateMasterRequest, createdBy *int64, actor AuthContext) (*dto.MasterResponse, error)
	//GetByID(id int64, actor AuthContext) (*dto.MasterResponse, error)
	//List(page, pageSize int, filter *dto.FilterMasterRequest, actor AuthContext) ([]dto.MasterResponse, int64, error)
	//Update(id int64, req *dto.UpdateMasterRequest, updatedBy *int64, actor AuthContext) (*dto.MasterResponse, error)
	//Delete(id int64, actor AuthContext) error
}
