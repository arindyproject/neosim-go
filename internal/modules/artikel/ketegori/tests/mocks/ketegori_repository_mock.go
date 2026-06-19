package mocks

import (
	"neosim_go/internal/modules/artikel/ketegori/dto"
	"neosim_go/internal/modules/artikel/ketegori/models"
	"github.com/stretchr/testify/mock"
)

// ArtikelKetegoriRepositoryMock is a mock implementation of contracts.Repository
type ArtikelKetegoriRepositoryMock struct {
	mock.Mock
}

func (m *ArtikelKetegoriRepositoryMock) Create(item *models.ArtikelKetegori) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *ArtikelKetegoriRepositoryMock) GetByID(id int64) (*models.ArtikelKetegori, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ArtikelKetegori), args.Error(1)
}

func (m *ArtikelKetegoriRepositoryMock) List(page, pageSize int, filter *dto.FilterArtikelKetegoriRequest) ([]models.ArtikelKetegori, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.ArtikelKetegori), args.Get(1).(int64), args.Error(2)
}

func (m *ArtikelKetegoriRepositoryMock) Update(item *models.ArtikelKetegori) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *ArtikelKetegoriRepositoryMock) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
