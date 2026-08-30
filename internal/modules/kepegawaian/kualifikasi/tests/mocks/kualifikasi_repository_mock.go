package mocks

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/kualifikasi/dto"
	"neosim_go/internal/modules/kepegawaian/kualifikasi/models"
	"github.com/stretchr/testify/mock"
)

// KepegawaianKualifikasiRepositoryMock is a mock implementation of contracts.Repository.
// Ketika item ditambahkan (mode add-item), method mock untuk item tersebut
// ditempelkan ke struct INI JUGA (mis. tests/mocks/tag_repository_mock.go),
// bukan membuat mock struct baru.
type KepegawaianKualifikasiRepositoryMock struct {
	mock.Mock
}

func (m *KepegawaianKualifikasiRepositoryMock) CreateKualifikasi(ctx context.Context,item *models.KepegawaianKualifikasi) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianKualifikasiRepositoryMock) GetKualifikasiByID(ctx context.Context,id int64) (*models.KepegawaianKualifikasi, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KepegawaianKualifikasi), args.Error(1)
}

func (m *KepegawaianKualifikasiRepositoryMock) GetByIDs(ctx context.Context,ids []int64) ([]models.KepegawaianKualifikasi, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.KepegawaianKualifikasi), args.Error(1)
}

func (m *KepegawaianKualifikasiRepositoryMock) ListKualifikasi(ctx context.Context,page, pageSize int, filter *dto.FilterKepegawaianKualifikasiRequest) ([]models.KepegawaianKualifikasi, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.KepegawaianKualifikasi), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianKualifikasiRepositoryMock) UpdateKualifikasi(ctx context.Context,item *models.KepegawaianKualifikasi) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianKualifikasiRepositoryMock) DeleteKualifikasi(ctx context.Context,id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
