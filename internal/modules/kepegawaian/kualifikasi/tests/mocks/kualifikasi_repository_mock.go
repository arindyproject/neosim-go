package mocks

import (
	"neosim_go/internal/modules/kepegawaian/kualifikasi/dto"
	"neosim_go/internal/modules/kepegawaian/kualifikasi/models"
	"github.com/stretchr/testify/mock"
)

// KepegawaianKualifikasiRepositoryMock is a mock implementation of contracts.Repository
type KepegawaianKualifikasiRepositoryMock struct {
	mock.Mock
}

func (m *KepegawaianKualifikasiRepositoryMock) Create(item *models.KepegawaianKualifikasi) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianKualifikasiRepositoryMock) GetByID(id int64) (*models.KepegawaianKualifikasi, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KepegawaianKualifikasi), args.Error(1)
}

func (m *KepegawaianKualifikasiRepositoryMock) List(page, pageSize int, filter *dto.FilterKepegawaianKualifikasiRequest) ([]models.KepegawaianKualifikasi, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.KepegawaianKualifikasi), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianKualifikasiRepositoryMock) Update(item *models.KepegawaianKualifikasi) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianKualifikasiRepositoryMock) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
