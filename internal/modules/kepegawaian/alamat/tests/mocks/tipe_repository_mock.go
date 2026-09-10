package mocks

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/alamat/dto"
	"neosim_go/internal/modules/kepegawaian/alamat/models"
)

// Method di bawah ini ditempelkan ke KepegawaianAlamatRepositoryMock yang
// sama dengan mock entitas utama (lihat tests/mocks/alamat_repository_mock.go).

func (m *KepegawaianAlamatRepositoryMock) CreateTipe(ctx context.Context,item *models.Tipe) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianAlamatRepositoryMock) GetTipeByID(ctx context.Context,id int64) (*models.Tipe, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tipe), args.Error(1)
}

func (m *KepegawaianAlamatRepositoryMock) ListTipe(ctx context.Context,page, pageSize int, filter *dto.FilterTipeRequest) ([]models.Tipe, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.Tipe), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianAlamatRepositoryMock) UpdateTipe(ctx context.Context,item *models.Tipe) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianAlamatRepositoryMock) DeleteTipe(ctx context.Context,id int64, deletedBy int64) error {
	args := m.Called(id, deletedBy)
	return args.Error(0)
}
