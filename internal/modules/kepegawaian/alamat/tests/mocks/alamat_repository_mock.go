package mocks

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/alamat/dto"
	"neosim_go/internal/modules/kepegawaian/alamat/models"

	"github.com/stretchr/testify/mock"
)

// KepegawaianAlamatRepositoryMock is a mock implementation of contracts.Repository.
// Ketika item ditambahkan (mode add-item), method mock untuk item tersebut
// ditempelkan ke struct INI JUGA (mis. tests/mocks/tag_repository_mock.go),
// bukan membuat mock struct baru.
type KepegawaianAlamatRepositoryMock struct {
	mock.Mock
}

func (m *KepegawaianAlamatRepositoryMock) CheckDuplicateAlamat(ctx context.Context, model *models.KepegawaianAlamat, excludeID *int64) (bool, error) {
	args := m.Called(model, excludeID)
	if args.Get(0) == nil {
		return false, args.Error(1)
	}
	return args.Bool(0), args.Error(1)
}

func (m *KepegawaianAlamatRepositoryMock) CreateAlamat(ctx context.Context, item *models.KepegawaianAlamat) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianAlamatRepositoryMock) FindAlamatByPegawaiID(ctx context.Context, pegawaiID int64) ([]models.KepegawaianAlamat, error) {
	args := m.Called(pegawaiID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.KepegawaianAlamat), args.Error(1)
}

func (m *KepegawaianAlamatRepositoryMock) GetAlamatByPegawaiID(ctx context.Context, pegawaiID int64, page, pageSize int) ([]models.KepegawaianAlamat, int64, error) {
	args := m.Called(pegawaiID, page, pageSize)
	if args.Get(0) == nil {
		return nil, int64(args.Int(1)), args.Error(2)
	}
	return args.Get(0).([]models.KepegawaianAlamat), int64(args.Int(1)), args.Error(2)
}

func (m *KepegawaianAlamatRepositoryMock) GetAlamatByID(ctx context.Context, id int64) (*models.KepegawaianAlamat, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KepegawaianAlamat), args.Error(1)
}

func (m *KepegawaianAlamatRepositoryMock) GetByIDs(ctx context.Context, ids []int64) ([]models.KepegawaianAlamat, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.KepegawaianAlamat), args.Error(1)
}

func (m *KepegawaianAlamatRepositoryMock) ListAlamat(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianAlamatRequest) ([]models.KepegawaianAlamat, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.KepegawaianAlamat), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianAlamatRepositoryMock) UpdateAlamat(ctx context.Context, item *models.KepegawaianAlamat) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianAlamatRepositoryMock) DeleteAlamat(ctx context.Context, id int64, deletedBy int64) error {
	args := m.Called(id, deletedBy)
	return args.Error(0)
}

func (m *KepegawaianAlamatRepositoryMock) UnsetPrimaryAlamatByPegawaiID(ctx context.Context, pegawaiID int64, updatedBy int64) error {
	args := m.Called(pegawaiID, updatedBy)
	return args.Error(0)
}
