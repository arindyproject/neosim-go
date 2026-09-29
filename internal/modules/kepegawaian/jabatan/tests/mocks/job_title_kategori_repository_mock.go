package mocks

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
)

// Method di bawah ini ditempelkan ke KepegawaianJabatanRepositoryMock yang
// sama dengan mock entitas utama (lihat tests/mocks/jabatan_repository_mock.go).

func (m *KepegawaianJabatanRepositoryMock) CreateJobTitleKategori(ctx context.Context, item *models.JobTitleKategori) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianJabatanRepositoryMock) GetJobTitleKategoriByID(ctx context.Context, id int64) (*models.JobTitleKategori, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.JobTitleKategori), args.Error(1)
}

func (m *KepegawaianJabatanRepositoryMock) ListSelectJobTitleKategori(ctx context.Context, search string) ([]models.JobTitleKategori, error) {
	args := m.Called(search)
	return args.Get(0).([]models.JobTitleKategori), args.Error(1)
}

func (m *KepegawaianJabatanRepositoryMock) ListJobTitleKategori(ctx context.Context, page, pageSize int, filter *dto.FilterJobTitleKategoriRequest) ([]models.JobTitleKategori, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.JobTitleKategori), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianJabatanRepositoryMock) UpdateJobTitleKategori(ctx context.Context, item *models.JobTitleKategori) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianJabatanRepositoryMock) DeleteJobTitleKategori(ctx context.Context, id int64, deletedBy int64) error {
	args := m.Called(id, deletedBy)
	return args.Error(0)
}

func (m *KepegawaianJabatanRepositoryMock) GetJobTitleKategoriByCode(ctx context.Context, code string) (*models.JobTitleKategori, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.JobTitleKategori), args.Error(1)
}

func (m *KepegawaianJabatanRepositoryMock) GetJobTitleKategoriByLabel(ctx context.Context, label string) (*models.JobTitleKategori, error) {
	args := m.Called(label)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.JobTitleKategori), args.Error(1)
}

// PERBAIKAN: Pastikan ctx DAN id dimasukkan ke dalam m.Called
func (m *KepegawaianJabatanRepositoryMock) CheckJobTitleKategori(ctx context.Context, id int64) (bool, error) {
	args := m.Called(ctx, id) // <-- Tambahkan ctx di sini
	return args.Bool(0), args.Error(1)
}
