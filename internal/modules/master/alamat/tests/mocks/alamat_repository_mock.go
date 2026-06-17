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

func (m *MasterAlamatRepositoryMock) Create(item *models.MasterAlamat) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) GetByID(id int64) (*models.MasterAlamat, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterAlamat), args.Error(1)
}

func (m *MasterAlamatRepositoryMock) List(page, pageSize int, filter *dto.FilterMasterAlamatRequest) ([]models.MasterAlamat, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.MasterAlamat), args.Get(1).(int64), args.Error(2)
}

func (m *MasterAlamatRepositoryMock) Update(item *models.MasterAlamat) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterAlamatRepositoryMock) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
