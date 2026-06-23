package mocks

import (
	"neosim_go/internal/modules/kepegawaian/lampiran/dto"
	"neosim_go/internal/modules/kepegawaian/lampiran/models"
	"github.com/stretchr/testify/mock"
)

// KepegawaianLampiranRepositoryMock is a mock implementation of contracts.Repository
type KepegawaianLampiranRepositoryMock struct {
	mock.Mock
}

func (m *KepegawaianLampiranRepositoryMock) Create(item *models.KepegawaianLampiran) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianLampiranRepositoryMock) GetByID(id int64) (*models.KepegawaianLampiran, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KepegawaianLampiran), args.Error(1)
}

func (m *KepegawaianLampiranRepositoryMock) List(page, pageSize int, filter *dto.FilterKepegawaianLampiranRequest) ([]models.KepegawaianLampiran, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.KepegawaianLampiran), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianLampiranRepositoryMock) Update(item *models.KepegawaianLampiran) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianLampiranRepositoryMock) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
