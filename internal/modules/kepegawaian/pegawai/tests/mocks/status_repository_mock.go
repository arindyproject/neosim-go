package mocks

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/pegawai/dto"
	"neosim_go/internal/modules/kepegawaian/pegawai/models"
)

// Method di bawah ini ditempelkan ke KepegawaianPegawaiRepositoryMock yang
// sama dengan mock entitas utama (lihat tests/mocks/pegawai_repository_mock.go).

func (m *KepegawaianPegawaiRepositoryMock) CreateStatus(ctx context.Context, item *models.Status) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianPegawaiRepositoryMock) GetStatusByID(ctx context.Context, id int64) (*models.Status, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Status), args.Error(1)
}

func (m *KepegawaianPegawaiRepositoryMock) GetStatusByCode(ctx context.Context, code string) (*models.Status, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Status), args.Error(1)
}

func (m *KepegawaianPegawaiRepositoryMock) GetStatusByLabel(ctx context.Context, label string) (*models.Status, error) {
	args := m.Called(label)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Status), args.Error(1)
}

func (m *KepegawaianPegawaiRepositoryMock) ListSelectStatus(ctx context.Context, search string) ([]models.Status, error) {
	args := m.Called(search)
	return args.Get(0).([]models.Status), args.Error(1)
}

func (m *KepegawaianPegawaiRepositoryMock) ListStatus(ctx context.Context, page, pageSize int, filter *dto.FilterStatusRequest) ([]models.Status, int64, error) {
	args := m.Called(page, pageSize, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Status), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianPegawaiRepositoryMock) UpdateStatus(ctx context.Context, item *models.Status) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianPegawaiRepositoryMock) DeleteStatus(ctx context.Context, id int64, deletedBy int64) error {
	args := m.Called(id, deletedBy)
	return args.Error(0)
}
