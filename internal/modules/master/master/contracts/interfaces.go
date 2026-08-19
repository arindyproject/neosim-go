package contracts

import (
	"context"
	"neosim_go/internal/modules/master/master/dto"
	"neosim_go/internal/modules/master/master/models"
	he "neosim_go/internal/shared/httputil"
)

// Repository defines database operations
type Repository interface {
	// Pekerjaan----------------------------------------------------
	CreatePekerjaan(ctx context.Context, m *models.MasterPekerjaan) error
	GetByIDPekerjaan(ctx context.Context, id int64) (*models.MasterPekerjaan, error)
	GetByNamePekerjaan(ctx context.Context, name string) (*models.MasterPekerjaan, error)
	ListPekerjaan(ctx context.Context, page, pageSize int, filter *dto.FilterMasterPekerjaanRequest) ([]models.MasterPekerjaan, int64, error)
	UpdatePekerjaan(ctx context.Context, m *models.MasterPekerjaan) error
	DeletePekerjaan(ctx context.Context, id int64) error
	//--------------------------------------------------------------

	// Pendidikan---------------------------------------------------
	CreatePendidikan(ctx context.Context, m *models.MasterPendidikan) error
	GetByIDPendidikan(ctx context.Context, id int64) (*models.MasterPendidikan, error)
	GetByNamePendidikan(ctx context.Context, name string) (*models.MasterPendidikan, error)
	ListPendidikan(ctx context.Context, page, pageSize int, filter *dto.FilterMasterPendidikanRequest) ([]models.MasterPendidikan, int64, error)
	UpdatePendidikan(ctx context.Context, m *models.MasterPendidikan) error
	DeletePendidikan(ctx context.Context, id int64) error
	//---------------------------------------------------------------

	// Agama---------------------------------------------------------
	CreateAgama(ctx context.Context, m *models.MasterAgama) error
	GetByIDAgama(ctx context.Context, id int64) (*models.MasterAgama, error)
	GetByNameAgama(ctx context.Context, name string) (*models.MasterAgama, error)
	ListAgama(ctx context.Context, page, pageSize int, filter *dto.FilterMasterAgamaRequest) ([]models.MasterAgama, int64, error)
	UpdateAgama(ctx context.Context, m *models.MasterAgama) error
	DeleteAgama(ctx context.Context, id int64) error
	//----------------------------------------------------------------

	// StatusPernikahan-----------------------------------------------
	CreateStatusPernikahan(ctx context.Context, m *models.MasterStatusPernikahan) error
	GetByIDStatusPernikahan(ctx context.Context, id int64) (*models.MasterStatusPernikahan, error)
	GetByNameStatusPernikahan(ctx context.Context, name string) (*models.MasterStatusPernikahan, error)
	ListStatusPernikahan(ctx context.Context, page, pageSize int, filter *dto.FilterMasterStatusPernikahanRequest) ([]models.MasterStatusPernikahan, int64, error)
	UpdateStatusPernikahan(ctx context.Context, m *models.MasterStatusPernikahan) error
	DeleteStatusPernikahan(ctx context.Context, id int64) error
	//----------------------------------------------------------------

	// GolonganDarah--------------------------------------------------
	CreateGolonganDarah(ctx context.Context, m *models.MasterGolonganDarah) error
	GetByIDGolonganDarah(ctx context.Context, id int64) (*models.MasterGolonganDarah, error)
	GetByNameGolonganDarah(ctx context.Context, name string) (*models.MasterGolonganDarah, error)
	ListGolonganDarah(ctx context.Context, page, pageSize int, filter *dto.FilterMasterGolonganDarahRequest) ([]models.MasterGolonganDarah, int64, error)
	UpdateGolonganDarah(ctx context.Context, m *models.MasterGolonganDarah) error
	DeleteGolonganDarah(ctx context.Context, id int64) error
	//----------------------------------------------------------------

	// Suku-----------------------------------------------------------
	CreateSuku(ctx context.Context, m *models.MasterSuku) error
	GetByIDSuku(ctx context.Context, id int64) (*models.MasterSuku, error)
	GetByNameSuku(ctx context.Context, name string) (*models.MasterSuku, error)
	ListSuku(ctx context.Context, page, pageSize int, filter *dto.FilterMasterSukuRequest) ([]models.MasterSuku, int64, error)
	UpdateSuku(ctx context.Context, m *models.MasterSuku) error
	DeleteSuku(ctx context.Context, id int64) error
	//----------------------------------------------------------------

	// JenisKelamin---------------------------------------------------
	CreateJenisKelamin(ctx context.Context, m *models.MasterJenisKelamin) error
	GetByIDJenisKelamin(ctx context.Context, id int64) (*models.MasterJenisKelamin, error)
	GetByNameJenisKelamin(ctx context.Context, name string) (*models.MasterJenisKelamin, error)
	ListJenisKelamin(ctx context.Context, page, pageSize int, filter *dto.FilterMasterJenisKelaminRequest) ([]models.MasterJenisKelamin, int64, error)
	UpdateJenisKelamin(ctx context.Context, m *models.MasterJenisKelamin) error
	DeleteJenisKelamin(ctx context.Context, id int64) error
	//----------------------------------------------------------------
}

// Service defines business logic operations
type Service interface {
	// Pekerjaan----------------------------------------------------
	GetByIDPekerjaan(ctx context.Context, id int64) (*dto.MasterPekerjaanResponse, error)
	ListPekerjaan(ctx context.Context, page, pageSize int, filter *dto.FilterMasterPekerjaanRequest) ([]dto.MasterPekerjaanResponse, int64, error)
	CreatePekerjaan(ctx context.Context, req *dto.CreateMasterPekerjaanRequest, actor he.AuthContext) (*dto.MasterPekerjaanResponse, error)
	UpdatePekerjaan(ctx context.Context, id int64, req *dto.UpdateMasterPekerjaanRequest, actor he.AuthContext) (*dto.MasterPekerjaanResponse, error)
	DeletePekerjaan(ctx context.Context, id int64, actor he.AuthContext) error
	//--------------------------------------------------------------

	// Pendidikan---------------------------------------------------
	GetByIDPendidikan(ctx context.Context, id int64) (*dto.MasterPendidikanResponse, error)
	ListPendidikan(ctx context.Context, page, pageSize int, filter *dto.FilterMasterPendidikanRequest) ([]dto.MasterPendidikanResponse, int64, error)
	CreatePendidikan(ctx context.Context, req *dto.CreateMasterPendidikanRequest, actor he.AuthContext) (*dto.MasterPendidikanResponse, error)
	UpdatePendidikan(ctx context.Context, id int64, req *dto.UpdateMasterPendidikanRequest, actor he.AuthContext) (*dto.MasterPendidikanResponse, error)
	DeletePendidikan(ctx context.Context, id int64, actor he.AuthContext) error
	//--------------------------------------------------------------

	// Agama--------------------------------------------------------
	GetByIDAgama(ctx context.Context, id int64) (*dto.MasterAgamaResponse, error)
	ListAgama(ctx context.Context, page, pageSize int, filter *dto.FilterMasterAgamaRequest) ([]dto.MasterAgamaResponse, int64, error)
	CreateAgama(ctx context.Context, req *dto.CreateMasterAgamaRequest, actor he.AuthContext) (*dto.MasterAgamaResponse, error)
	UpdateAgama(ctx context.Context, id int64, req *dto.UpdateMasterAgamaRequest, actor he.AuthContext) (*dto.MasterAgamaResponse, error)
	DeleteAgama(ctx context.Context, id int64, actor he.AuthContext) error
	//--------------------------------------------------------------

	// StatusPernikahan---------------------------------------------
	GetByIDStatusPernikahan(ictx context.Context, d int64) (*dto.MasterStatusPernikahanResponse, error)
	ListStatusPernikahan(ctx context.Context, page, pageSize int, filter *dto.FilterMasterStatusPernikahanRequest) ([]dto.MasterStatusPernikahanResponse, int64, error)
	CreateStatusPernikahan(ctx context.Context, req *dto.CreateMasterStatusPernikahanRequest, actor he.AuthContext) (*dto.MasterStatusPernikahanResponse, error)
	UpdateStatusPernikahan(ctx context.Context, id int64, req *dto.UpdateMasterStatusPernikahanRequest, actor he.AuthContext) (*dto.MasterStatusPernikahanResponse, error)
	DeleteStatusPernikahan(ctx context.Context, id int64, actor he.AuthContext) error
	//--------------------------------------------------------------

	// JenisKelamin-------------------------------------------------
	GetByIDJenisKelamin(ctx context.Context, id int64) (*dto.MasterJenisKelaminResponse, error)
	ListJenisKelamin(ctx context.Context, page, pageSize int, filter *dto.FilterMasterJenisKelaminRequest) ([]dto.MasterJenisKelaminResponse, int64, error)
	CreateJenisKelamin(ctx context.Context, req *dto.CreateMasterJenisKelaminRequest, actor he.AuthContext) (*dto.MasterJenisKelaminResponse, error)
	UpdateJenisKelamin(ctx context.Context, id int64, req *dto.UpdateMasterJenisKelaminRequest, actor he.AuthContext) (*dto.MasterJenisKelaminResponse, error)
	DeleteJenisKelamin(ctx context.Context, id int64, actor he.AuthContext) error
	//--------------------------------------------------------------

	// Suku---------------------------------------------------------
	GetByIDSuku(ctx context.Context, id int64) (*dto.MasterSukuResponse, error)
	ListSuku(ctx context.Context, page, pageSize int, filter *dto.FilterMasterSukuRequest) ([]dto.MasterSukuResponse, int64, error)
	CreateSuku(ctx context.Context, req *dto.CreateMasterSukuRequest, actor he.AuthContext) (*dto.MasterSukuResponse, error)
	UpdateSuku(ctx context.Context, id int64, req *dto.UpdateMasterSukuRequest, actor he.AuthContext) (*dto.MasterSukuResponse, error)
	DeleteSuku(ctx context.Context, id int64, actor he.AuthContext) error
	//--------------------------------------------------------------

	// GolonganDarah------------------------------------------------
	GetByIDGolonganDarah(ctx context.Context, id int64) (*dto.MasterGolonganDarahResponse, error)
	ListGolonganDarah(ctx context.Context, page, pageSize int, filter *dto.FilterMasterGolonganDarahRequest) ([]dto.MasterGolonganDarahResponse, int64, error)
	CreateGolonganDarah(ctx context.Context, req *dto.CreateMasterGolonganDarahRequest, actor he.AuthContext) (*dto.MasterGolonganDarahResponse, error)
	UpdateGolonganDarah(ctx context.Context, id int64, req *dto.UpdateMasterGolonganDarahRequest, actor he.AuthContext) (*dto.MasterGolonganDarahResponse, error)
	DeleteGolonganDarah(ctx context.Context, id int64, actor he.AuthContext) error
	//--------------------------------------------------------------
}
