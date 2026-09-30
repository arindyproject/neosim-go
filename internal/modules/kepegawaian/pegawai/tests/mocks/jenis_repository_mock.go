package mocks

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/pegawai/dto"
	"neosim_go/internal/modules/kepegawaian/pegawai/models"
)

// Method di bawah ini ditempelkan ke KepegawaianPegawaiRepositoryMock yang
// sama dengan mock entitas utama (lihat tests/mocks/pegawai_repository_mock.go).

func (m *KepegawaianPegawaiRepositoryMock) CreateJenis(ctx context.Context, item *models.Jenis) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianPegawaiRepositoryMock) GetJenisByID(ctx context.Context, id int64) (*models.Jenis, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Jenis), args.Error(1)
}

func (m *KepegawaianPegawaiRepositoryMock) GetJenisByCode(ctx context.Context, code string) (*models.Jenis, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Jenis), args.Error(1)
}

func (m *KepegawaianPegawaiRepositoryMock) GetJenisByLabel(ctx context.Context, label string) (*models.Jenis, error) {
	args := m.Called(label)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Jenis), args.Error(1)
}

func (m *KepegawaianPegawaiRepositoryMock) ListSelectJenis(ctx context.Context, search string) ([]models.Jenis, error) {
	args := m.Called(search)
	return args.Get(0).([]models.Jenis), args.Error(1)
}

func (m *KepegawaianPegawaiRepositoryMock) ListJenis(ctx context.Context, page, pageSize int, filter *dto.FilterJenisRequest) ([]models.Jenis, int64, error) {
	args := m.Called(page, pageSize, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Jenis), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianPegawaiRepositoryMock) UpdateJenis(ctx context.Context, item *models.Jenis) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianPegawaiRepositoryMock) DeleteJenis(ctx context.Context, id int64, deletedBy int64) error {
	args := m.Called(id, deletedBy)
	return args.Error(0)
}
