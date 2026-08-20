package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"neosim_go/internal/modules/kepegawaian/pegawai/dto"
	"neosim_go/internal/modules/kepegawaian/pegawai/models"
)

// KepegawaianPegawaiRepositoryMock is a mock implementation of contracts.Repository.
// Ketika item ditambahkan (mode add-item), method mock untuk item tersebut
// ditempelkan ke struct INI JUGA (mis. tests/mocks/tag_repository_mock.go),
// bukan membuat mock struct baru.
type KepegawaianPegawaiRepositoryMock struct {
	mock.Mock
}

func (m *KepegawaianPegawaiRepositoryMock) CreatePegawai(ctx context.Context, item *models.KepegawaianPegawai) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianPegawaiRepositoryMock) GetPegawaiByID(ctx context.Context, id int64) (*models.KepegawaianPegawai, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KepegawaianPegawai), args.Error(1)
}

func (m *KepegawaianPegawaiRepositoryMock) GetByIDs(ctx context.Context, ids []int64) ([]models.KepegawaianPegawai, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.KepegawaianPegawai), args.Error(1)
}

func (m *KepegawaianPegawaiRepositoryMock) ListPegawai(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianPegawaiRequest) ([]models.KepegawaianPegawai, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.KepegawaianPegawai), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianPegawaiRepositoryMock) UpdatePegawai(ctx context.Context, item *models.KepegawaianPegawai) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianPegawaiRepositoryMock) DeletePegawai(ctx context.Context, id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
