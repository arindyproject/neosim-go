package mocks

import (
	"context"
	"neosim_go/internal/modules/master/master/dto"
	"neosim_go/internal/modules/master/master/models"

	"github.com/stretchr/testify/mock"
)

type MasterRepositoryMock struct {
	mock.Mock
}

// Pekerjaan ============================================================

func (m *MasterRepositoryMock) CreatePekerjaan(ctx context.Context, item *models.MasterPekerjaan) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterRepositoryMock) GetByIDPekerjaan(ctx context.Context, id int64) (*models.MasterPekerjaan, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterPekerjaan), args.Error(1)
}

func (m *MasterRepositoryMock) GetByNamePekerjaan(ctx context.Context, name string) (*models.MasterPekerjaan, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterPekerjaan), args.Error(1)
}

func (m *MasterRepositoryMock) ListPekerjaan(ctx context.Context, page, pageSize int, filter *dto.FilterMasterPekerjaanRequest) ([]models.MasterPekerjaan, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.MasterPekerjaan), args.Get(1).(int64), args.Error(2)
}

func (m *MasterRepositoryMock) UpdatePekerjaan(ctx context.Context, item *models.MasterPekerjaan) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterRepositoryMock) DeletePekerjaan(ctx context.Context, id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

// Pendidikan =========================================================

func (m *MasterRepositoryMock) CreatePendidikan(ctx context.Context, item *models.MasterPendidikan) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterRepositoryMock) GetByIDPendidikan(ctx context.Context, id int64) (*models.MasterPendidikan, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterPendidikan), args.Error(1)
}

func (m *MasterRepositoryMock) GetByNamePendidikan(ctx context.Context, name string) (*models.MasterPendidikan, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterPendidikan), args.Error(1)
}

func (m *MasterRepositoryMock) ListPendidikan(ctx context.Context, page, pageSize int, filter *dto.FilterMasterPendidikanRequest) ([]models.MasterPendidikan, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.MasterPendidikan), args.Get(1).(int64), args.Error(2)
}

func (m *MasterRepositoryMock) UpdatePendidikan(ctx context.Context, item *models.MasterPendidikan) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterRepositoryMock) DeletePendidikan(ctx context.Context, id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

// Agama ================================================================

func (m *MasterRepositoryMock) CreateAgama(ctx context.Context, item *models.MasterAgama) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterRepositoryMock) GetByIDAgama(ctx context.Context, id int64) (*models.MasterAgama, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterAgama), args.Error(1)
}

func (m *MasterRepositoryMock) GetByNameAgama(ctx context.Context, name string) (*models.MasterAgama, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterAgama), args.Error(1)
}

func (m *MasterRepositoryMock) ListAgama(ctx context.Context, page, pageSize int, filter *dto.FilterMasterAgamaRequest) ([]models.MasterAgama, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.MasterAgama), args.Get(1).(int64), args.Error(2)
}

func (m *MasterRepositoryMock) UpdateAgama(ctx context.Context, item *models.MasterAgama) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterRepositoryMock) DeleteAgama(ctx context.Context, id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

// Status Pernikahan =================================================

func (m *MasterRepositoryMock) CreateStatusPernikahan(ctx context.Context, item *models.MasterStatusPernikahan) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterRepositoryMock) GetByIDStatusPernikahan(ctx context.Context, id int64) (*models.MasterStatusPernikahan, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterStatusPernikahan), args.Error(1)
}

func (m *MasterRepositoryMock) GetByNameStatusPernikahan(ctx context.Context, name string) (*models.MasterStatusPernikahan, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterStatusPernikahan), args.Error(1)
}

func (m *MasterRepositoryMock) ListStatusPernikahan(ctx context.Context, page, pageSize int, filter *dto.FilterMasterStatusPernikahanRequest) ([]models.MasterStatusPernikahan, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.MasterStatusPernikahan), args.Get(1).(int64), args.Error(2)
}

func (m *MasterRepositoryMock) UpdateStatusPernikahan(ctx context.Context, item *models.MasterStatusPernikahan) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterRepositoryMock) DeleteStatusPernikahan(ctx context.Context, id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

// Suku ================================================================

func (m *MasterRepositoryMock) CreateSuku(ctx context.Context, item *models.MasterSuku) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterRepositoryMock) GetByIDSuku(ctx context.Context, id int64) (*models.MasterSuku, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterSuku), args.Error(1)
}

func (m *MasterRepositoryMock) GetByNameSuku(ctx context.Context, name string) (*models.MasterSuku, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterSuku), args.Error(1)
}

func (m *MasterRepositoryMock) ListSuku(ctx context.Context, page, pageSize int, filter *dto.FilterMasterSukuRequest) ([]models.MasterSuku, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.MasterSuku), args.Get(1).(int64), args.Error(2)
}

func (m *MasterRepositoryMock) UpdateSuku(ctx context.Context, item *models.MasterSuku) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterRepositoryMock) DeleteSuku(ctx context.Context, id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

// Golongan Darah ================================================================

func (m *MasterRepositoryMock) CreateGolonganDarah(ctx context.Context, item *models.MasterGolonganDarah) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterRepositoryMock) GetByIDGolonganDarah(ctx context.Context, id int64) (*models.MasterGolonganDarah, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterGolonganDarah), args.Error(1)
}

func (m *MasterRepositoryMock) GetByNameGolonganDarah(ctx context.Context, name string) (*models.MasterGolonganDarah, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterGolonganDarah), args.Error(1)
}

func (m *MasterRepositoryMock) ListGolonganDarah(ctx context.Context, page, pageSize int, filter *dto.FilterMasterGolonganDarahRequest) ([]models.MasterGolonganDarah, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.MasterGolonganDarah), args.Get(1).(int64), args.Error(2)
}

func (m *MasterRepositoryMock) UpdateGolonganDarah(ctx context.Context, item *models.MasterGolonganDarah) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterRepositoryMock) DeleteGolonganDarah(ctx context.Context, id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

// Jenis Kelamin ================================================================

func (m *MasterRepositoryMock) CreateJenisKelamin(ctx context.Context, item *models.MasterJenisKelamin) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterRepositoryMock) GetByIDJenisKelamin(ctx context.Context, id int64) (*models.MasterJenisKelamin, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterJenisKelamin), args.Error(1)
}

func (m *MasterRepositoryMock) GetByNameJenisKelamin(ctx context.Context, name string) (*models.MasterJenisKelamin, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterJenisKelamin), args.Error(1)
}

func (m *MasterRepositoryMock) ListJenisKelamin(ctx context.Context, page, pageSize int, filter *dto.FilterMasterJenisKelaminRequest) ([]models.MasterJenisKelamin, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.MasterJenisKelamin), args.Get(1).(int64), args.Error(2)
}

func (m *MasterRepositoryMock) UpdateJenisKelamin(ctx context.Context, item *models.MasterJenisKelamin) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterRepositoryMock) DeleteJenisKelamin(ctx context.Context, id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
