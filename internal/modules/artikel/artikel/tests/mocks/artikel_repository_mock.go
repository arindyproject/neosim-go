package mocks

import (
	"context"
	"neosim_go/internal/modules/artikel/artikel/dto"
	"neosim_go/internal/modules/artikel/artikel/models"
	"github.com/stretchr/testify/mock"
)

// ArtikelRepositoryMock is a mock implementation of contracts.Repository.
// Ketika item ditambahkan (mode add-item), method mock untuk item tersebut
// ditempelkan ke struct INI JUGA (mis. tests/mocks/tag_repository_mock.go),
// bukan membuat mock struct baru.
type ArtikelRepositoryMock struct {
	mock.Mock
}

func (m *ArtikelRepositoryMock) CreateArtikel(ctx context.Context,item *models.Artikel) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *ArtikelRepositoryMock) GetArtikelByID(ctx context.Context,id int64) (*models.Artikel, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Artikel), args.Error(1)
}

func (m *ArtikelRepositoryMock) GetByIDs(ctx context.Context,ids []int64) ([]models.Artikel, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Artikel), args.Error(1)
}

func (m *ArtikelRepositoryMock) ListArtikel(ctx context.Context,page, pageSize int, filter *dto.FilterArtikelRequest) ([]models.Artikel, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.Artikel), args.Get(1).(int64), args.Error(2)
}

func (m *ArtikelRepositoryMock) UpdateArtikel(ctx context.Context,item *models.Artikel) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *ArtikelRepositoryMock) DeleteArtikel(ctx context.Context,id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
