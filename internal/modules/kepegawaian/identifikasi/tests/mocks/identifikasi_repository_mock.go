package mocks

import (
	"neosim_go/internal/modules/kepegawaian/identifikasi/dto"
	"neosim_go/internal/modules/kepegawaian/identifikasi/models"
	"github.com/stretchr/testify/mock"
)

// KepegawaianIdentifikasiRepositoryMock is a mock implementation of contracts.Repository
type KepegawaianIdentifikasiRepositoryMock struct {
	mock.Mock
}

func (m *KepegawaianIdentifikasiRepositoryMock) Create(item *models.KepegawaianIdentifikasi) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianIdentifikasiRepositoryMock) GetByID(id int64) (*models.KepegawaianIdentifikasi, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KepegawaianIdentifikasi), args.Error(1)
}

func (m *KepegawaianIdentifikasiRepositoryMock) List(page, pageSize int, filter *dto.FilterKepegawaianIdentifikasiRequest) ([]models.KepegawaianIdentifikasi, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.KepegawaianIdentifikasi), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianIdentifikasiRepositoryMock) Update(item *models.KepegawaianIdentifikasi) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianIdentifikasiRepositoryMock) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
