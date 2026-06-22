package mocks

import (
	"neosim_go/internal/modules/master/master/dto"
	"neosim_go/internal/modules/master/master/models"

	"github.com/stretchr/testify/mock"
)

type MasterRepositoryMock struct {
	mock.Mock
}

// Pekerjaan ============================================================

func (m *MasterRepositoryMock) CreatePekerjaan(item *models.MasterPekerjaan) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterRepositoryMock) GetByIDPekerjaan(id int64) (*models.MasterPekerjaan, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterPekerjaan), args.Error(1)
}

func (m *MasterRepositoryMock) GetByNamePekerjaan(name string) (*models.MasterPekerjaan, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterPekerjaan), args.Error(1)
}

func (m *MasterRepositoryMock) ListPekerjaan(page, pageSize int, filter *dto.FilterMasterPekerjaanRequest) ([]models.MasterPekerjaan, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.MasterPekerjaan), args.Get(1).(int64), args.Error(2)
}

func (m *MasterRepositoryMock) UpdatePekerjaan(item *models.MasterPekerjaan) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterRepositoryMock) DeletePekerjaan(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

// Pendidikan =========================================================

func (m *MasterRepositoryMock) CreatePendidikan(item *models.MasterPendidikan) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterRepositoryMock) GetByIDPendidikan(id int64) (*models.MasterPendidikan, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterPendidikan), args.Error(1)
}

func (m *MasterRepositoryMock) GetByNamePendidikan(name string) (*models.MasterPendidikan, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterPendidikan), args.Error(1)
}

func (m *MasterRepositoryMock) ListPendidikan(page, pageSize int, filter *dto.FilterMasterPendidikanRequest) ([]models.MasterPendidikan, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.MasterPendidikan), args.Get(1).(int64), args.Error(2)
}

func (m *MasterRepositoryMock) UpdatePendidikan(item *models.MasterPendidikan) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterRepositoryMock) DeletePendidikan(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

// Agama ================================================================

func (m *MasterRepositoryMock) CreateAgama(item *models.MasterAgama) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterRepositoryMock) GetByIDAgama(id int64) (*models.MasterAgama, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterAgama), args.Error(1)
}

func (m *MasterRepositoryMock) GetByNameAgama(name string) (*models.MasterAgama, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterAgama), args.Error(1)
}

func (m *MasterRepositoryMock) ListAgama(page, pageSize int, filter *dto.FilterMasterAgamaRequest) ([]models.MasterAgama, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.MasterAgama), args.Get(1).(int64), args.Error(2)
}

func (m *MasterRepositoryMock) UpdateAgama(item *models.MasterAgama) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterRepositoryMock) DeleteAgama(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

// Status Pernikahan =================================================

func (m *MasterRepositoryMock) CreateStatusPernikahan(item *models.MasterStatusPernikahan) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterRepositoryMock) GetByIDStatusPernikahan(id int64) (*models.MasterStatusPernikahan, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterStatusPernikahan), args.Error(1)
}

func (m *MasterRepositoryMock) GetByNameStatusPernikahan(name string) (*models.MasterStatusPernikahan, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterStatusPernikahan), args.Error(1)
}

func (m *MasterRepositoryMock) ListStatusPernikahan(page, pageSize int, filter *dto.FilterMasterStatusPernikahanRequest) ([]models.MasterStatusPernikahan, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.MasterStatusPernikahan), args.Get(1).(int64), args.Error(2)
}

func (m *MasterRepositoryMock) UpdateStatusPernikahan(item *models.MasterStatusPernikahan) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterRepositoryMock) DeleteStatusPernikahan(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
