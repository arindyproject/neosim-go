package contracts

import (
	"neosim_go/internal/modules/master/master/dto"
	"neosim_go/internal/modules/master/master/models"
	he "neosim_go/internal/shared/httputil"
)

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

	// GolonganDarah--------------------------------------------------
	CreateGolonganDarah(m *models.MasterGolonganDarah) error
	GetByIDGolonganDarah(id int64) (*models.MasterGolonganDarah, error)
	GetByNameGolonganDarah(name string) (*models.MasterGolonganDarah, error)
	ListGolonganDarah(page, pageSize int, filter *dto.FilterMasterGolonganDarahRequest) ([]models.MasterGolonganDarah, int64, error)
	UpdateGolonganDarah(m *models.MasterGolonganDarah) error
	DeleteGolonganDarah(id int64) error
	//----------------------------------------------------------------

	// Suku-----------------------------------------------------------
	CreateSuku(m *models.MasterSuku) error
	GetByIDSuku(id int64) (*models.MasterSuku, error)
	GetByNameSuku(name string) (*models.MasterSuku, error)
	ListSuku(page, pageSize int, filter *dto.FilterMasterSukuRequest) ([]models.MasterSuku, int64, error)
	UpdateSuku(m *models.MasterSuku) error
	DeleteSuku(id int64) error
	//----------------------------------------------------------------

	// JenisKelamin---------------------------------------------------
	CreateJenisKelamin(m *models.MasterJenisKelamin) error
	GetByIDJenisKelamin(id int64) (*models.MasterJenisKelamin, error)
	GetByNameJenisKelamin(name string) (*models.MasterJenisKelamin, error)
	ListJenisKelamin(page, pageSize int, filter *dto.FilterMasterJenisKelaminRequest) ([]models.MasterJenisKelamin, int64, error)
	UpdateJenisKelamin(m *models.MasterJenisKelamin) error
	DeleteJenisKelamin(id int64) error
	//----------------------------------------------------------------
}

// Service defines business logic operations
type Service interface {
	// Pekerjaan----------------------------------------------------
	GetByIDPekerjaan(id int64) (*dto.MasterPekerjaanResponse, error)
	ListPekerjaan(page, pageSize int, filter *dto.FilterMasterPekerjaanRequest) ([]dto.MasterPekerjaanResponse, int64, error)
	CreatePekerjaan(req *dto.CreateMasterPekerjaanRequest, actor he.AuthContext) (*dto.MasterPekerjaanResponse, error)
	UpdatePekerjaan(id int64, req *dto.UpdateMasterPekerjaanRequest, actor he.AuthContext) (*dto.MasterPekerjaanResponse, error)
	DeletePekerjaan(id int64, actor he.AuthContext) error
	//--------------------------------------------------------------

	// Pendidikan---------------------------------------------------
	GetByIDPendidikan(id int64) (*dto.MasterPendidikanResponse, error)
	ListPendidikan(page, pageSize int, filter *dto.FilterMasterPendidikanRequest) ([]dto.MasterPendidikanResponse, int64, error)
	CreatePendidikan(req *dto.CreateMasterPendidikanRequest, actor he.AuthContext) (*dto.MasterPendidikanResponse, error)
	UpdatePendidikan(id int64, req *dto.UpdateMasterPendidikanRequest, actor he.AuthContext) (*dto.MasterPendidikanResponse, error)
	DeletePendidikan(id int64, actor he.AuthContext) error
	//--------------------------------------------------------------

	// Agama--------------------------------------------------------
	GetByIDAgama(id int64) (*dto.MasterAgamaResponse, error)
	ListAgama(page, pageSize int, filter *dto.FilterMasterAgamaRequest) ([]dto.MasterAgamaResponse, int64, error)
	CreateAgama(req *dto.CreateMasterAgamaRequest, actor he.AuthContext) (*dto.MasterAgamaResponse, error)
	UpdateAgama(id int64, req *dto.UpdateMasterAgamaRequest, actor he.AuthContext) (*dto.MasterAgamaResponse, error)
	DeleteAgama(id int64, actor he.AuthContext) error
	//--------------------------------------------------------------

	// StatusPernikahan-----------------------------------------------
	GetByIDStatusPernikahan(id int64) (*dto.MasterStatusPernikahanResponse, error)
	ListStatusPernikahan(page, pageSize int, filter *dto.FilterMasterStatusPernikahanRequest) ([]dto.MasterStatusPernikahanResponse, int64, error)
	CreateStatusPernikahan(req *dto.CreateMasterStatusPernikahanRequest, actor he.AuthContext) (*dto.MasterStatusPernikahanResponse, error)
	UpdateStatusPernikahan(id int64, req *dto.UpdateMasterStatusPernikahanRequest, actor he.AuthContext) (*dto.MasterStatusPernikahanResponse, error)
	DeleteStatusPernikahan(id int64, actor he.AuthContext) error
	//--------------------------------------------------------------
}
