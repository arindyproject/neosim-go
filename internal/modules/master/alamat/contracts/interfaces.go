package contracts

import (
	"neosim_go/internal/modules/master/alamat/dto"
	"neosim_go/internal/modules/master/alamat/models"
	he "neosim_go/internal/shared/httputil"
)

// Repository defines database operations
type Repository interface {
	// Negara-------------------------------------------------------
	CreateNegara(m *models.MasterAlamatNegara) error
	GetByIDNegara(id int64) (*models.MasterAlamatNegara, error)
	ListNegara(page, pageSize int, filter *dto.FilterNegaraRequest) ([]models.MasterAlamatNegara, int64, error)
	UpdateNegara(m *models.MasterAlamatNegara) error
	DeleteNegara(id int64) error
	ExistsNegaraByCode(code string, excludeID *int64) (bool, error)
	//--------------------------------------------------------------

	// Provinsi-----------------------------------------------------
	CreateProvinsi(m *models.MasterAlamatProvinsi) error
	GetByIDProvinsi(id int64) (*models.MasterAlamatProvinsi, error)
	ListProvinsi(page, pageSize int, negaraID *int64, filter *dto.FilterProvinsiRequest) ([]models.MasterAlamatProvinsi, int64, error)
	UpdateProvinsi(m *models.MasterAlamatProvinsi) error
	DeleteProvinsi(id int64) error
	ExistsProvinsiByCode(code string, excludeID *int64) (bool, error)
	//tambahan
	CountKotaByProvinsiID(provinsiID int64) (int64, error)
	CountKecamatanByProvinsiID(provinsiID int64) (int64, error)
	CountDesaByProvinsiID(provinsiID int64) (int64, error)
	//--------------------------------------------------------------

	// Kota/Kabupaten-----------------------------------------------
	CreateKotaKabupaten(m *models.MasterAlamatKotaKabupaten) error
	GetByIDKotaKabupaten(id int64) (*models.MasterAlamatKotaKabupaten, error)
	ListKotaKabupaten(page, pageSize int, provinsiID *int64, filter *dto.FilterKotaKabupatenRequest) ([]models.MasterAlamatKotaKabupaten, int64, error)
	UpdateKotaKabupaten(m *models.MasterAlamatKotaKabupaten) error
	DeleteKotaKabupaten(id int64) error
	ExistsKotaKabupatenByCode(code string, excludeID *int64) (bool, error)
	// tambahan
	CountKecamatanByKotaID(kotaID int64) (int64, error)
	CountDesaByKotaID(kotaID int64) (int64, error)
	//--------------------------------------------------------------

	// Kecamatan----------------------------------------------------
	CreateKecamatan(m *models.MasterAlamatKecamatan) error
	GetByIDKecamatan(id int64) (*models.MasterAlamatKecamatan, error)
	ListKecamatan(page, pageSize int, kotaKabupatenID *int64, filter *dto.FilterKecamatanRequest) ([]models.MasterAlamatKecamatan, int64, error)
	UpdateKecamatan(m *models.MasterAlamatKecamatan) error
	DeleteKecamatan(id int64) error
	ExistsKecamatanByCode(code string, excludeID *int64) (bool, error)
	// tambahan
	CountDesaByKecamatanID(kecamatanID int64) (int64, error)
	//--------------------------------------------------------------

	// Kelurahan/Desa-----------------------------------------------
	CreateKelurahanDesa(m *models.MasterAlamatKelurahanDesa) error
	GetByIDKelurahanDesa(id int64) (*models.MasterAlamatKelurahanDesa, error)
	ListKelurahanDesa(page, pageSize int, kecamatanID *int64, filter *dto.FilterKelurahanDesaRequest) ([]models.MasterAlamatKelurahanDesa, int64, error)
	UpdateKelurahanDesa(m *models.MasterAlamatKelurahanDesa) error
	DeleteKelurahanDesa(id int64) error
	ExistsKelurahanDesaByCode(code string, excludeID *int64) (bool, error)
	//--------------------------------------------------------------
}

// Service defines business logic operations
type Service interface {
	// Negara-------------------------------------------------------
	GetByIDNegara(id int64) (*dto.NegaraResponse, error)
	ListNegara(page, pageSize int, filter *dto.FilterNegaraRequest) ([]dto.NegaraResponse, int64, error)
	CreateNegara(req *dto.CreateNegaraRequest, actor he.AuthContext) (*dto.NegaraResponse, error)
	UpdateNegara(id int64, req *dto.UpdateNegaraRequest, actor he.AuthContext) (*dto.NegaraResponse, error)
	DeleteNegara(id int64, actor he.AuthContext) error
	// Negara-------------------------------------------------------

	// Provinsi-----------------------------------------------------
	GetByIDProvinsi(id int64) (*dto.ProvinsiDetailResponse, error)
	ListProvinsi(page, pageSize int, negaraID *int64, filter *dto.FilterProvinsiRequest) ([]dto.ProvinsiResponse, int64, error)
	CreateProvinsi(req *dto.CreateProvinsiRequest, actor he.AuthContext) (*dto.ProvinsiResponse, error)
	UpdateProvinsi(id int64, req *dto.UpdateProvinsiRequest, actor he.AuthContext) (*dto.ProvinsiResponse, error)
	DeleteProvinsi(id int64, actor he.AuthContext) error
	//--------------------------------------------------------------

	// Kota/Kabupaten-------------------------------------------------
	GetByIDKotaKabupaten(id int64) (*dto.KotaKabupatenDetailResponse, error)
	ListKotaKabupaten(page, pageSize int, provinsiID *int64, filter *dto.FilterKotaKabupatenRequest) ([]dto.KotaKabupatenResponse, int64, error)
	CreateKotaKabupaten(req *dto.CreateKotaKabupatenRequest, actor he.AuthContext) (*dto.KotaKabupatenResponse, error)
	UpdateKotaKabupaten(id int64, req *dto.UpdateKotaKabupatenRequest, actor he.AuthContext) (*dto.KotaKabupatenResponse, error)
	DeleteKotaKabupaten(id int64, actor he.AuthContext) error
	//------------------------------------------------------------------

	// Kecamatan--------------------------------------------------------
	GetByIDKecamatan(id int64) (*dto.KecamatanDetailResponse, error)
	ListKecamatan(page, pageSize int, kotaKabupatenID *int64, filter *dto.FilterKecamatanRequest) ([]dto.KecamatanResponse, int64, error)
	CreateKecamatan(req *dto.CreateKecamatanRequest, actor he.AuthContext) (*dto.KecamatanResponse, error)
	UpdateKecamatan(id int64, req *dto.UpdateKecamatanRequest, actor he.AuthContext) (*dto.KecamatanResponse, error)
	DeleteKecamatan(id int64, actor he.AuthContext) error
	//------------------------------------------------------------------

	// Kelurahan/Desa---------------------------------------------------
	GetByIDKelurahanDesa(id int64) (*dto.KelurahanDesaDetailResponse, error)
	ListKelurahanDesa(page, pageSize int, kecamatanID *int64, filter *dto.FilterKelurahanDesaRequest) ([]dto.KelurahanDesaResponse, int64, error)
	CreateKelurahanDesa(req *dto.CreateKelurahanDesaRequest, actor he.AuthContext) (*dto.KelurahanDesaResponse, error)
	UpdateKelurahanDesa(id int64, req *dto.UpdateKelurahanDesaRequest, actor he.AuthContext) (*dto.KelurahanDesaResponse, error)
	DeleteKelurahanDesa(id int64, actor he.AuthContext) error
	//------------------------------------------------------------------
}
