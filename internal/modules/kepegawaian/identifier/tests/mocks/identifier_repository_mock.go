package mocks

import (
	"neosim_go/internal/modules/kepegawaian/identifier/dto"
	"neosim_go/internal/modules/kepegawaian/identifier/models"
	"github.com/stretchr/testify/mock"
)

// KepegawaianIdentifierRepositoryMock is a mock implementation of contracts.Repository
type KepegawaianIdentifierRepositoryMock struct {
	mock.Mock
}

func (m *KepegawaianIdentifierRepositoryMock) Create(item *models.KepegawaianIdentifier) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianIdentifierRepositoryMock) GetByID(id int64) (*models.KepegawaianIdentifier, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KepegawaianIdentifier), args.Error(1)
}

func (m *KepegawaianIdentifierRepositoryMock) List(page, pageSize int, filter *dto.FilterKepegawaianIdentifierRequest) ([]models.KepegawaianIdentifier, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.KepegawaianIdentifier), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianIdentifierRepositoryMock) Update(item *models.KepegawaianIdentifier) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianIdentifierRepositoryMock) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
