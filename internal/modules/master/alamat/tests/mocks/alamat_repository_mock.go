package mocks

import (
	"context"
	"neosim_go/internal/modules/master/alamat/dto"
	"neosim_go/internal/modules/master/alamat/models"

	"github.com/stretchr/testify/mock"
)

// MasterAlamatRepositoryMock is a mock implementation of contracts.Repository
type MasterAlamatRepositoryMock struct {
	mock.Mock
}

// Negara ===========================================================================

func (m *MasterAlamatRepositoryMock) CreateNegara(ctx context.Context, item *models.MasterAlamatNegara) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) GetByIDNegara(ctx context.Context, id int64) (*models.MasterAlamatNegara, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterAlamatNegara), args.Error(1)
}

func (m *MasterAlamatRepositoryMock) ListNegara(ctx context.Context, page, pageSize int, filter *dto.FilterNegaraRequest) ([]models.MasterAlamatNegara, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.MasterAlamatNegara), args.Get(1).(int64), args.Error(2)
}

func (m *MasterAlamatRepositoryMock) UpdateNegara(ctx context.Context, item *models.MasterAlamatNegara) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) DeleteNegara(ctx context.Context, id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) ExistsNegaraByCode(ctx context.Context, code string, excludeID *int64) (bool, error) {
	args := m.Called(code, excludeID)
	return args.Bool(0), args.Error(1)
}

// Provinsi ==========================================================================

func (m *MasterAlamatRepositoryMock) CreateProvinsi(ctx context.Context, item *models.MasterAlamatProvinsi) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) GetByIDProvinsi(ctx context.Context, id int64) (*models.MasterAlamatProvinsi, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterAlamatProvinsi), args.Error(1)
}

func (m *MasterAlamatRepositoryMock) ListProvinsi(ctx context.Context, page, pageSize int, negaraID *int64, filter *dto.FilterProvinsiRequest) ([]models.MasterAlamatProvinsi, int64, error) {
	args := m.Called(page, pageSize, negaraID, filter)
	return args.Get(0).([]models.MasterAlamatProvinsi), args.Get(1).(int64), args.Error(2)
}

func (m *MasterAlamatRepositoryMock) UpdateProvinsi(ctx context.Context, item *models.MasterAlamatProvinsi) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) DeleteProvinsi(ctx context.Context, id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) CountKotaByProvinsiID(ctx context.Context, provinsiID int64) (int64, error) {
	args := m.Called(provinsiID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MasterAlamatRepositoryMock) CountKecamatanByProvinsiID(ctx context.Context, provinsiID int64) (int64, error) {
	args := m.Called(provinsiID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MasterAlamatRepositoryMock) CountDesaByProvinsiID(ctx context.Context, provinsiID int64) (int64, error) {
	args := m.Called(provinsiID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MasterAlamatRepositoryMock) ExistsProvinsiByCode(ctx context.Context, code string, excludeID *int64) (bool, error) {
	args := m.Called(code, excludeID)
	return args.Bool(0), args.Error(1)
}

// Kota/Kabupaten =====================================================================

func (m *MasterAlamatRepositoryMock) CreateKotaKabupaten(ctx context.Context, item *models.MasterAlamatKotaKabupaten) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) GetByIDKotaKabupaten(ctx context.Context, id int64) (*models.MasterAlamatKotaKabupaten, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterAlamatKotaKabupaten), args.Error(1)
}

func (m *MasterAlamatRepositoryMock) ListKotaKabupaten(ctx context.Context, page, pageSize int, provinsiID *int64, filter *dto.FilterKotaKabupatenRequest) ([]models.MasterAlamatKotaKabupaten, int64, error) {
	args := m.Called(page, pageSize, provinsiID, filter)
	return args.Get(0).([]models.MasterAlamatKotaKabupaten), args.Get(1).(int64), args.Error(2)
}

func (m *MasterAlamatRepositoryMock) UpdateKotaKabupaten(ctx context.Context, item *models.MasterAlamatKotaKabupaten) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) DeleteKotaKabupaten(ctx context.Context, id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) CountKecamatanByKotaID(ctx context.Context, kotaID int64) (int64, error) {
	args := m.Called(kotaID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MasterAlamatRepositoryMock) CountDesaByKotaID(ctx context.Context, kotaID int64) (int64, error) {
	args := m.Called(kotaID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MasterAlamatRepositoryMock) ExistsKotaKabupatenByCode(ctx context.Context, code string, excludeID *int64) (bool, error) {
	args := m.Called(code, excludeID)
	return args.Bool(0), args.Error(1)
}

// Kecamatan ===========================================================================

func (m *MasterAlamatRepositoryMock) CreateKecamatan(ctx context.Context, item *models.MasterAlamatKecamatan) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) GetByIDKecamatan(ctx context.Context, id int64) (*models.MasterAlamatKecamatan, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterAlamatKecamatan), args.Error(1)
}

func (m *MasterAlamatRepositoryMock) ListKecamatan(ctx context.Context, page, pageSize int, kotaKabupatenID *int64, filter *dto.FilterKecamatanRequest) ([]models.MasterAlamatKecamatan, int64, error) {
	args := m.Called(page, pageSize, kotaKabupatenID, filter)
	return args.Get(0).([]models.MasterAlamatKecamatan), args.Get(1).(int64), args.Error(2)
}

func (m *MasterAlamatRepositoryMock) UpdateKecamatan(ctx context.Context, item *models.MasterAlamatKecamatan) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) DeleteKecamatan(ctx context.Context, id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) CountDesaByKecamatanID(ctx context.Context, kecamatanID int64) (int64, error) {
	args := m.Called(kecamatanID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MasterAlamatRepositoryMock) ExistsKecamatanByCode(ctx context.Context, code string, excludeID *int64) (bool, error) {
	args := m.Called(code, excludeID)
	return args.Bool(0), args.Error(1)
}

// Kelurahan/Desa ========================================================================

func (m *MasterAlamatRepositoryMock) CreateKelurahanDesa(ctx context.Context, item *models.MasterAlamatKelurahanDesa) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) GetByIDKelurahanDesa(ctx context.Context, id int64) (*models.MasterAlamatKelurahanDesa, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterAlamatKelurahanDesa), args.Error(1)
}

func (m *MasterAlamatRepositoryMock) ListKelurahanDesa(ctx context.Context, page, pageSize int, kecamatanID *int64, filter *dto.FilterKelurahanDesaRequest) ([]models.MasterAlamatKelurahanDesa, int64, error) {
	args := m.Called(page, pageSize, kecamatanID, filter)
	return args.Get(0).([]models.MasterAlamatKelurahanDesa), args.Get(1).(int64), args.Error(2)
}

func (m *MasterAlamatRepositoryMock) UpdateKelurahanDesa(ctx context.Context, item *models.MasterAlamatKelurahanDesa) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) DeleteKelurahanDesa(ctx context.Context, id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) ExistsKelurahanDesaByCode(ctx context.Context, code string, excludeID *int64) (bool, error) {
	args := m.Called(code, excludeID)
	return args.Bool(0), args.Error(1)
}
