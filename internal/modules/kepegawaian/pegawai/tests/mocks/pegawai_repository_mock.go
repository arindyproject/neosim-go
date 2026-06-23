package mocks

import (
	"neosim_go/internal/modules/kepegawaian/pegawai/dto"
	"neosim_go/internal/modules/kepegawaian/pegawai/models"
	"github.com/stretchr/testify/mock"
)

// KepegawaianPegawaiRepositoryMock is a mock implementation of contracts.Repository
type KepegawaianPegawaiRepositoryMock struct {
	mock.Mock
}

func (m *KepegawaianPegawaiRepositoryMock) Create(item *models.KepegawaianPegawai) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianPegawaiRepositoryMock) GetByID(id int64) (*models.KepegawaianPegawai, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KepegawaianPegawai), args.Error(1)
}

func (m *KepegawaianPegawaiRepositoryMock) List(page, pageSize int, filter *dto.FilterKepegawaianPegawaiRequest) ([]models.KepegawaianPegawai, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.KepegawaianPegawai), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianPegawaiRepositoryMock) Update(item *models.KepegawaianPegawai) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianPegawaiRepositoryMock) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
