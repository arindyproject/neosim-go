package mocks

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/pendidikan/dto"
	"neosim_go/internal/modules/kepegawaian/pendidikan/models"
)

// Method di bawah ini ditempelkan ke KepegawaianPendidikanRepositoryMock yang
// sama dengan mock entitas utama (lihat tests/mocks/pendidikan_repository_mock.go).

func (m *KepegawaianPendidikanRepositoryMock) CreateJenjang(ctx context.Context,item *models.Jenjang) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianPendidikanRepositoryMock) GetJenjangByID(ctx context.Context,id int64) (*models.Jenjang, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Jenjang), args.Error(1)
}

func (m *KepegawaianPendidikanRepositoryMock) ListJenjang(ctx context.Context,page, pageSize int, filter *dto.FilterJenjangRequest) ([]models.Jenjang, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.Jenjang), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianPendidikanRepositoryMock) UpdateJenjang(ctx context.Context,item *models.Jenjang) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianPendidikanRepositoryMock) DeleteJenjang(ctx context.Context,id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
