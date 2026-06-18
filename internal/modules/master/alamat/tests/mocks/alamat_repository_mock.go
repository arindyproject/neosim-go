package mocks

import (
	"neosim_go/internal/modules/master/alamat/dto"
	"neosim_go/internal/modules/master/alamat/models"

	"github.com/stretchr/testify/mock"
)

// MasterAlamatRepositoryMock is a mock implementation of contracts.Repository
type MasterAlamatRepositoryMock struct {
	mock.Mock
}

// Negara ===========================================================================

func (m *MasterAlamatRepositoryMock) CreateNegara(item *models.MasterAlamatNegara) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) GetByIDNegara(id int64) (*models.MasterAlamatNegara, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterAlamatNegara), args.Error(1)
}

func (m *MasterAlamatRepositoryMock) ListNegara(page, pageSize int, filter *dto.FilterNegaraRequest) ([]models.MasterAlamatNegara, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.MasterAlamatNegara), args.Get(1).(int64), args.Error(2)
}

func (m *MasterAlamatRepositoryMock) UpdateNegara(item *models.MasterAlamatNegara) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) DeleteNegara(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

// Provinsi ==========================================================================

func (m *MasterAlamatRepositoryMock) CreateProvinsi(item *models.MasterAlamatProvinsi) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) GetByIDProvinsi(id int64) (*models.MasterAlamatProvinsi, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterAlamatProvinsi), args.Error(1)
}

func (m *MasterAlamatRepositoryMock) ListProvinsi(page, pageSize int, negaraID *int64, filter *dto.FilterProvinsiRequest) ([]models.MasterAlamatProvinsi, int64, error) {
	args := m.Called(page, pageSize, negaraID, filter)
	return args.Get(0).([]models.MasterAlamatProvinsi), args.Get(1).(int64), args.Error(2)
}

func (m *MasterAlamatRepositoryMock) UpdateProvinsi(item *models.MasterAlamatProvinsi) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) DeleteProvinsi(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) CountKotaByProvinsiID(provinsiID int64) (int64, error) {
	args := m.Called(provinsiID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MasterAlamatRepositoryMock) CountKecamatanByProvinsiID(provinsiID int64) (int64, error) {
	args := m.Called(provinsiID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MasterAlamatRepositoryMock) CountDesaByProvinsiID(provinsiID int64) (int64, error) {
	args := m.Called(provinsiID)
	return args.Get(0).(int64), args.Error(1)
}

// Kota/Kabupaten =====================================================================

func (m *MasterAlamatRepositoryMock) CreateKotaKabupaten(item *models.MasterAlamatKotaKabupaten) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) GetByIDKotaKabupaten(id int64) (*models.MasterAlamatKotaKabupaten, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterAlamatKotaKabupaten), args.Error(1)
}

func (m *MasterAlamatRepositoryMock) ListKotaKabupaten(page, pageSize int, provinsiID *int64, filter *dto.FilterKotaKabupatenRequest) ([]models.MasterAlamatKotaKabupaten, int64, error) {
	args := m.Called(page, pageSize, provinsiID, filter)
	return args.Get(0).([]models.MasterAlamatKotaKabupaten), args.Get(1).(int64), args.Error(2)
}

func (m *MasterAlamatRepositoryMock) UpdateKotaKabupaten(item *models.MasterAlamatKotaKabupaten) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) DeleteKotaKabupaten(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) CountKecamatanByKotaID(kotaID int64) (int64, error) {
	args := m.Called(kotaID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MasterAlamatRepositoryMock) CountDesaByKotaID(kotaID int64) (int64, error) {
	args := m.Called(kotaID)
	return args.Get(0).(int64), args.Error(1)
}

// Kecamatan ===========================================================================

func (m *MasterAlamatRepositoryMock) CreateKecamatan(item *models.MasterAlamatKecamatan) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) GetByIDKecamatan(id int64) (*models.MasterAlamatKecamatan, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterAlamatKecamatan), args.Error(1)
}

func (m *MasterAlamatRepositoryMock) ListKecamatan(page, pageSize int, kotaKabupatenID *int64, filter *dto.FilterKecamatanRequest) ([]models.MasterAlamatKecamatan, int64, error) {
	args := m.Called(page, pageSize, kotaKabupatenID, filter)
	return args.Get(0).([]models.MasterAlamatKecamatan), args.Get(1).(int64), args.Error(2)
}

func (m *MasterAlamatRepositoryMock) UpdateKecamatan(item *models.MasterAlamatKecamatan) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) DeleteKecamatan(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) CountDesaByKecamatanID(kecamatanID int64) (int64, error) {
	args := m.Called(kecamatanID)
	return args.Get(0).(int64), args.Error(1)
}

// Kelurahan/Desa ========================================================================

func (m *MasterAlamatRepositoryMock) CreateKelurahanDesa(item *models.MasterAlamatKelurahanDesa) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) GetByIDKelurahanDesa(id int64) (*models.MasterAlamatKelurahanDesa, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterAlamatKelurahanDesa), args.Error(1)
}

func (m *MasterAlamatRepositoryMock) ListKelurahanDesa(page, pageSize int, kecamatanID *int64, filter *dto.FilterKelurahanDesaRequest) ([]models.MasterAlamatKelurahanDesa, int64, error) {
	args := m.Called(page, pageSize, kecamatanID, filter)
	return args.Get(0).([]models.MasterAlamatKelurahanDesa), args.Get(1).(int64), args.Error(2)
}

func (m *MasterAlamatRepositoryMock) UpdateKelurahanDesa(item *models.MasterAlamatKelurahanDesa) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) DeleteKelurahanDesa(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
