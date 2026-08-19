package contracts

import (
	"context"
	"neosim_go/internal/modules/master/alamat/dto"
	"neosim_go/internal/modules/master/alamat/models"
	he "neosim_go/internal/shared/httputil"
)

// Repository defines database operations
type Repository interface {
	// Negara-------------------------------------------------------
	CreateNegara(ctx context.Context, m *models.MasterAlamatNegara) error
	GetByIDNegara(ctx context.Context, id int64) (*models.MasterAlamatNegara, error)
	ListNegara(ctx context.Context, page, pageSize int, filter *dto.FilterNegaraRequest) ([]models.MasterAlamatNegara, int64, error)
	UpdateNegara(ctx context.Context, m *models.MasterAlamatNegara) error
	DeleteNegara(ctx context.Context, id int64) error
	ExistsNegaraByCode(ctx context.Context, code string, excludeID *int64) (bool, error)
	//--------------------------------------------------------------

	// Provinsi-----------------------------------------------------
	CreateProvinsi(ctx context.Context, m *models.MasterAlamatProvinsi) error
	GetByIDProvinsi(ctx context.Context, id int64) (*models.MasterAlamatProvinsi, error)
	ListProvinsi(ctx context.Context, page, pageSize int, negaraID *int64, filter *dto.FilterProvinsiRequest) ([]models.MasterAlamatProvinsi, int64, error)
	UpdateProvinsi(ctx context.Context, m *models.MasterAlamatProvinsi) error
	DeleteProvinsi(ctx context.Context, id int64) error
	ExistsProvinsiByCode(ctx context.Context, code string, excludeID *int64) (bool, error)
	//tambahan
	CountKotaByProvinsiID(ctx context.Context, provinsiID int64) (int64, error)
	CountKecamatanByProvinsiID(ctx context.Context, provinsiID int64) (int64, error)
	CountDesaByProvinsiID(ctx context.Context, provinsiID int64) (int64, error)
	//--------------------------------------------------------------

	// Kota/Kabupaten-----------------------------------------------
	CreateKotaKabupaten(ctx context.Context, m *models.MasterAlamatKotaKabupaten) error
	GetByIDKotaKabupaten(ctx context.Context, id int64) (*models.MasterAlamatKotaKabupaten, error)
	ListKotaKabupaten(ctx context.Context, page, pageSize int, provinsiID *int64, filter *dto.FilterKotaKabupatenRequest) ([]models.MasterAlamatKotaKabupaten, int64, error)
	UpdateKotaKabupaten(ctx context.Context, m *models.MasterAlamatKotaKabupaten) error
	DeleteKotaKabupaten(ctx context.Context, id int64) error
	ExistsKotaKabupatenByCode(ctx context.Context, code string, excludeID *int64) (bool, error)
	// tambahan
	CountKecamatanByKotaID(ctx context.Context, kotaID int64) (int64, error)
	CountDesaByKotaID(ctx context.Context, kotaID int64) (int64, error)
	//--------------------------------------------------------------

	// Kecamatan----------------------------------------------------
	CreateKecamatan(ctx context.Context, m *models.MasterAlamatKecamatan) error
	GetByIDKecamatan(ctx context.Context, id int64) (*models.MasterAlamatKecamatan, error)
	ListKecamatan(ctx context.Context, page, pageSize int, kotaKabupatenID *int64, filter *dto.FilterKecamatanRequest) ([]models.MasterAlamatKecamatan, int64, error)
	UpdateKecamatan(ctx context.Context, m *models.MasterAlamatKecamatan) error
	DeleteKecamatan(ctx context.Context, id int64) error
	ExistsKecamatanByCode(ctx context.Context, code string, excludeID *int64) (bool, error)
	// tambahan
	CountDesaByKecamatanID(ctx context.Context, kecamatanID int64) (int64, error)
	//--------------------------------------------------------------

	// Kelurahan/Desa-----------------------------------------------
	CreateKelurahanDesa(ctx context.Context, m *models.MasterAlamatKelurahanDesa) error
	GetByIDKelurahanDesa(ctx context.Context, id int64) (*models.MasterAlamatKelurahanDesa, error)
	ListKelurahanDesa(ctx context.Context, page, pageSize int, kecamatanID *int64, filter *dto.FilterKelurahanDesaRequest) ([]models.MasterAlamatKelurahanDesa, int64, error)
	UpdateKelurahanDesa(ctx context.Context, m *models.MasterAlamatKelurahanDesa) error
	DeleteKelurahanDesa(ctx context.Context, id int64) error
	ExistsKelurahanDesaByCode(ctx context.Context, code string, excludeID *int64) (bool, error)
	//--------------------------------------------------------------
}

// Service defines business logic operations
type Service interface {
	// Negara-------------------------------------------------------
	GetByIDNegara(ctx context.Context, id int64) (*dto.NegaraResponse, error)
	ListNegara(ctx context.Context, page, pageSize int, filter *dto.FilterNegaraRequest) ([]dto.NegaraResponse, int64, error)
	CreateNegara(ctx context.Context, req *dto.CreateNegaraRequest, actor he.AuthContext) (*dto.NegaraResponse, error)
	UpdateNegara(ctx context.Context, id int64, req *dto.UpdateNegaraRequest, actor he.AuthContext) (*dto.NegaraResponse, error)
	DeleteNegara(ctx context.Context, id int64, actor he.AuthContext) error
	// Negara-------------------------------------------------------

	// Provinsi-----------------------------------------------------
	GetByIDProvinsi(ctx context.Context, id int64) (*dto.ProvinsiDetailResponse, error)
	ListProvinsi(ctx context.Context, page, pageSize int, negaraID *int64, filter *dto.FilterProvinsiRequest) ([]dto.ProvinsiResponse, int64, error)
	CreateProvinsi(ctx context.Context, req *dto.CreateProvinsiRequest, actor he.AuthContext) (*dto.ProvinsiResponse, error)
	UpdateProvinsi(ctx context.Context, id int64, req *dto.UpdateProvinsiRequest, actor he.AuthContext) (*dto.ProvinsiResponse, error)
	DeleteProvinsi(ctx context.Context, id int64, actor he.AuthContext) error
	//--------------------------------------------------------------

	// Kota/Kabupaten-------------------------------------------------
	GetByIDKotaKabupaten(ctx context.Context, id int64) (*dto.KotaKabupatenDetailResponse, error)
	ListKotaKabupaten(ctx context.Context, page, pageSize int, provinsiID *int64, filter *dto.FilterKotaKabupatenRequest) ([]dto.KotaKabupatenResponse, int64, error)
	CreateKotaKabupaten(ctx context.Context, req *dto.CreateKotaKabupatenRequest, actor he.AuthContext) (*dto.KotaKabupatenResponse, error)
	UpdateKotaKabupaten(ctx context.Context, id int64, req *dto.UpdateKotaKabupatenRequest, actor he.AuthContext) (*dto.KotaKabupatenResponse, error)
	DeleteKotaKabupaten(ctx context.Context, id int64, actor he.AuthContext) error
	//------------------------------------------------------------------

	// Kecamatan--------------------------------------------------------
	GetByIDKecamatan(ctx context.Context, id int64) (*dto.KecamatanDetailResponse, error)
	ListKecamatan(ctx context.Context, page, pageSize int, kotaKabupatenID *int64, filter *dto.FilterKecamatanRequest) ([]dto.KecamatanResponse, int64, error)
	CreateKecamatan(ctx context.Context, req *dto.CreateKecamatanRequest, actor he.AuthContext) (*dto.KecamatanResponse, error)
	UpdateKecamatan(ctx context.Context, id int64, req *dto.UpdateKecamatanRequest, actor he.AuthContext) (*dto.KecamatanResponse, error)
	DeleteKecamatan(ctx context.Context, id int64, actor he.AuthContext) error
	//------------------------------------------------------------------

	// Kelurahan/Desa---------------------------------------------------
	GetByIDKelurahanDesa(ctx context.Context, id int64) (*dto.KelurahanDesaDetailResponse, error)
	ListKelurahanDesa(ctx context.Context, page, pageSize int, kecamatanID *int64, filter *dto.FilterKelurahanDesaRequest) ([]dto.KelurahanDesaResponse, int64, error)
	CreateKelurahanDesa(ctx context.Context, req *dto.CreateKelurahanDesaRequest, actor he.AuthContext) (*dto.KelurahanDesaResponse, error)
	UpdateKelurahanDesa(ctx context.Context, id int64, req *dto.UpdateKelurahanDesaRequest, actor he.AuthContext) (*dto.KelurahanDesaResponse, error)
	DeleteKelurahanDesa(ctx context.Context, id int64, actor he.AuthContext) error
	//------------------------------------------------------------------
}
