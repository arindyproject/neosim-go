package mocks

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
)

// Method di bawah ini ditempelkan ke KepegawaianJabatanRepositoryMock yang
// sama dengan mock entitas utama (lihat tests/mocks/jabatan_repository_mock.go).

func (m *KepegawaianJabatanRepositoryMock) CreatePositionKategori(ctx context.Context,item *models.PositionKategori) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianJabatanRepositoryMock) GetPositionKategoriByID(ctx context.Context,id int64) (*models.PositionKategori, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PositionKategori), args.Error(1)
}

func (m *KepegawaianJabatanRepositoryMock) ListPositionKategori(ctx context.Context,page, pageSize int, filter *dto.FilterPositionKategoriRequest) ([]models.PositionKategori, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.PositionKategori), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianJabatanRepositoryMock) UpdatePositionKategori(ctx context.Context,item *models.PositionKategori) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianJabatanRepositoryMock) DeletePositionKategori(ctx context.Context,id int64, deletedBy int64) error {
	args := m.Called(id, deletedBy)
	return args.Error(0)
}
