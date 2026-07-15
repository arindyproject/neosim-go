package mocks

import (
	"neosim_go/internal/modules/artikel/kategori/dto"
	"neosim_go/internal/modules/artikel/kategori/models"
	"github.com/stretchr/testify/mock"
)

// ArtikelKategoriRepositoryMock is a mock implementation of contracts.Repository.
// Ketika item ditambahkan (mode add-item), method mock untuk item tersebut
// ditempelkan ke struct INI JUGA (mis. tests/mocks/tag_repository_mock.go),
// bukan membuat mock struct baru.
type ArtikelKategoriRepositoryMock struct {
	mock.Mock
}

func (m *ArtikelKategoriRepositoryMock) Create(item *models.ArtikelKategori) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *ArtikelKategoriRepositoryMock) GetByID(id int64) (*models.ArtikelKategori, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ArtikelKategori), args.Error(1)
}

func (m *ArtikelKategoriRepositoryMock) List(page, pageSize int, filter *dto.FilterArtikelKategoriRequest) ([]models.ArtikelKategori, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.ArtikelKategori), args.Get(1).(int64), args.Error(2)
}

func (m *ArtikelKategoriRepositoryMock) Update(item *models.ArtikelKategori) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *ArtikelKategoriRepositoryMock) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
