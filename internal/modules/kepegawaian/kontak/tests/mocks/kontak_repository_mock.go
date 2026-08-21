package mocks

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/kontak/dto"
	"neosim_go/internal/modules/kepegawaian/kontak/models"

	"github.com/stretchr/testify/mock"
)

// KepegawaianKontakRepositoryMock is a mock implementation of contracts.Repository.
// Ketika item ditambahkan (mode add-item), method mock untuk item tersebut
// ditempelkan ke struct INI JUGA (mis. tests/mocks/tag_repository_mock.go),
// bukan membuat mock struct baru.
type KepegawaianKontakRepositoryMock struct {
	mock.Mock
}

func (m *KepegawaianKontakRepositoryMock) CreateKontak(ctx context.Context, item *models.KepegawaianKontak) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianKontakRepositoryMock) GetKontakByID(ctx context.Context, id int64) (*models.KepegawaianKontak, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KepegawaianKontak), args.Error(1)
}

func (m *KepegawaianKontakRepositoryMock) GetKontakByPegawaiID(ctx context.Context, pegawaiID int64, page, pageSize int) ([]models.KepegawaianKontak, int64, error) {
	args := m.Called(pegawaiID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.KepegawaianKontak), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianKontakRepositoryMock) GetByPegawaiIDAndTipe(ctx context.Context, pegawaiID, tipeID int64) ([]models.KepegawaianKontak, error) {
	args := m.Called(pegawaiID, tipeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.KepegawaianKontak), args.Error(1)
}

func (m *KepegawaianKontakRepositoryMock) GetPrimaryByTipe(ctx context.Context, pegawaiID, tipeID int64) (*models.KepegawaianKontak, error) {
	args := m.Called(pegawaiID, tipeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KepegawaianKontak), args.Error(1)
}

func (m *KepegawaianKontakRepositoryMock) ListKontak(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianKontakRequest) ([]models.KepegawaianKontak, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.KepegawaianKontak), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianKontakRepositoryMock) UpdateKontak(ctx context.Context, item *models.KepegawaianKontak) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianKontakRepositoryMock) DeleteKontak(ctx context.Context, id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *KepegawaianKontakRepositoryMock) ExistsByNilaiAndTipe(ctx context.Context, tipeID int64,
	nilai string,
	excludeID int64) (bool, error) {
	args := m.Called(tipeID, nilai, excludeID)

	if args.Get(0) == nil {
		return false, args.Error(1)
	}
	return args.Bool(0), args.Error(1)

}

func (m *KepegawaianKontakRepositoryMock) UnsetPrimaryByPegawaiIDAndTipe(ctx context.Context, pegawaiID, tipeID int64, updatedBy int64) error {
	args := m.Called(pegawaiID, tipeID, updatedBy)
	return args.Error(0)
}
