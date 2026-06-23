package mocks

import (
	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	"github.com/stretchr/testify/mock"
)

// KepegawaianJabatanRepositoryMock is a mock implementation of contracts.Repository
type KepegawaianJabatanRepositoryMock struct {
	mock.Mock
}

func (m *KepegawaianJabatanRepositoryMock) Create(item *models.KepegawaianJabatan) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianJabatanRepositoryMock) GetByID(id int64) (*models.KepegawaianJabatan, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KepegawaianJabatan), args.Error(1)
}

func (m *KepegawaianJabatanRepositoryMock) List(page, pageSize int, filter *dto.FilterKepegawaianJabatanRequest) ([]models.KepegawaianJabatan, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.KepegawaianJabatan), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianJabatanRepositoryMock) Update(item *models.KepegawaianJabatan) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianJabatanRepositoryMock) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
