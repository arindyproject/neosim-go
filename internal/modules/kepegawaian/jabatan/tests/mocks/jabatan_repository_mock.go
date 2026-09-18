package mocks

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	"github.com/stretchr/testify/mock"
)

// KepegawaianJabatanRepositoryMock is a mock implementation of contracts.Repository.
// Ketika item ditambahkan (mode add-item), method mock untuk item tersebut
// ditempelkan ke struct INI JUGA (mis. tests/mocks/tag_repository_mock.go),
// bukan membuat mock struct baru.
type KepegawaianJabatanRepositoryMock struct {
	mock.Mock
}

func (m *KepegawaianJabatanRepositoryMock) CreateJabatan(ctx context.Context,item *models.KepegawaianJabatan) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianJabatanRepositoryMock) GetJabatanByID(ctx context.Context,id int64) (*models.KepegawaianJabatan, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KepegawaianJabatan), args.Error(1)
}

func (m *KepegawaianJabatanRepositoryMock) GetByIDs(ctx context.Context,ids []int64) ([]models.KepegawaianJabatan, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.KepegawaianJabatan), args.Error(1)
}

func (m *KepegawaianJabatanRepositoryMock) ListJabatan(ctx context.Context,page, pageSize int, filter *dto.FilterKepegawaianJabatanRequest) ([]models.KepegawaianJabatan, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.KepegawaianJabatan), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianJabatanRepositoryMock) UpdateJabatan(ctx context.Context,item *models.KepegawaianJabatan) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianJabatanRepositoryMock) DeleteJabatan(ctx context.Context,id int64, deletedBy int64) error {
	args := m.Called(id, deletedBy)
	return args.Error(0)
}
