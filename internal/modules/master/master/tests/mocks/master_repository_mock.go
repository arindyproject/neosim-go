package mocks

import (
	"neosim_go/internal/modules/master/master/dto"
	"neosim_go/internal/modules/master/master/models"
	"github.com/stretchr/testify/mock"
)

// MasterRepositoryMock is a mock implementation of contracts.Repository
type MasterRepositoryMock struct {
	mock.Mock
}

func (m *MasterRepositoryMock) Create(item *models.Master) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterRepositoryMock) GetByID(id int64) (*models.Master, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Master), args.Error(1)
}

func (m *MasterRepositoryMock) List(page, pageSize int, filter *dto.FilterMasterRequest) ([]models.Master, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.Master), args.Get(1).(int64), args.Error(2)
}

func (m *MasterRepositoryMock) Update(item *models.Master) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterRepositoryMock) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
