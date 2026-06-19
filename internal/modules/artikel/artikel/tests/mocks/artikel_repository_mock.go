package mocks

import (
	"neosim_go/internal/modules/artikel/artikel/dto"
	"neosim_go/internal/modules/artikel/artikel/models"
	"github.com/stretchr/testify/mock"
)

// ArtikelRepositoryMock is a mock implementation of contracts.Repository
type ArtikelRepositoryMock struct {
	mock.Mock
}

func (m *ArtikelRepositoryMock) Create(item *models.Artikel) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *ArtikelRepositoryMock) GetByID(id int64) (*models.Artikel, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Artikel), args.Error(1)
}

func (m *ArtikelRepositoryMock) List(page, pageSize int, filter *dto.FilterArtikelRequest) ([]models.Artikel, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.Artikel), args.Get(1).(int64), args.Error(2)
}

func (m *ArtikelRepositoryMock) Update(item *models.Artikel) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *ArtikelRepositoryMock) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
