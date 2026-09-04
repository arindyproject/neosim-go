package mocks

import (
	"context"
	"neosim_go/internal/modules/artikel/kategori/dto"
	"neosim_go/internal/modules/artikel/kategori/models"
)

// Method di bawah ini ditempelkan ke ArtikelKategoriRepositoryMock yang
// sama dengan mock entitas utama (lihat tests/mocks/kategori_repository_mock.go).

func (m *ArtikelKategoriRepositoryMock) CreateTag(ctx context.Context,item *models.Tag) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *ArtikelKategoriRepositoryMock) GetTagByID(ctx context.Context,id int64) (*models.Tag, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tag), args.Error(1)
}

func (m *ArtikelKategoriRepositoryMock) ListTag(ctx context.Context,page, pageSize int, filter *dto.FilterTagRequest) ([]models.Tag, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.Tag), args.Get(1).(int64), args.Error(2)
}

func (m *ArtikelKategoriRepositoryMock) UpdateTag(ctx context.Context,item *models.Tag) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *ArtikelKategoriRepositoryMock) DeleteTag(ctx context.Context,id int64, deletedBy int64) error {
	args := m.Called(id, deletedBy)
	return args.Error(0)
}
