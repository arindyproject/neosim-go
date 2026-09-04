package mocks

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/pendidikan/dto"
	"neosim_go/internal/modules/kepegawaian/pendidikan/models"

	"github.com/stretchr/testify/mock"
)

// KepegawaianPendidikanRepositoryMock is a mock implementation of contracts.Repository.
// Ketika item ditambahkan (mode add-item), method mock untuk item tersebut
// ditempelkan ke struct INI JUGA (mis. tests/mocks/tag_repository_mock.go),
// bukan membuat mock struct baru.
type KepegawaianPendidikanRepositoryMock struct {
	mock.Mock
}

func (m *KepegawaianPendidikanRepositoryMock) CreatePendidikan(ctx context.Context, item *models.KepegawaianPendidikan) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianPendidikanRepositoryMock) GetPendidikanByID(ctx context.Context, id int64) (*models.KepegawaianPendidikan, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KepegawaianPendidikan), args.Error(1)
}

func (m *KepegawaianPendidikanRepositoryMock) GetByIDs(ctx context.Context, ids []int64) ([]models.KepegawaianPendidikan, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.KepegawaianPendidikan), args.Error(1)
}

func (m *KepegawaianPendidikanRepositoryMock) GetPendidikanByPegawaiID(ctx context.Context, pegawaiID int64, page, pageSize int) ([]models.KepegawaianPendidikan, int64, error) {
	args := m.Called(pegawaiID)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.KepegawaianPendidikan), args.Get(1).(int64), args.Error(2)
}
func (m *KepegawaianPendidikanRepositoryMock) GetByPegawaiIDAndTipe(ctx context.Context, pegawaiID, jenjangID int64) ([]models.KepegawaianPendidikan, error) {
	args := m.Called(ctx, pegawaiID, jenjangID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.KepegawaianPendidikan), args.Error(1)
}

func (m *KepegawaianPendidikanRepositoryMock) ListPendidikan(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianPendidikanRequest) ([]models.KepegawaianPendidikan, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.KepegawaianPendidikan), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianPendidikanRepositoryMock) UpdatePendidikan(ctx context.Context, item *models.KepegawaianPendidikan) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianPendidikanRepositoryMock) DeletePendidikan(ctx context.Context, id int64, deletedBy int64) error {
	args := m.Called(id, deletedBy)
	return args.Error(0)
}

func (m *KepegawaianPendidikanRepositoryMock) ExistsPendidikanByID(ctx context.Context, id int64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}
func (m *KepegawaianPendidikanRepositoryMock) ExistsByNomorIjazah(ctx context.Context, jenjangID int64, nomorIjazah string, excludeID int64) (bool, error) {
	args := m.Called(ctx, jenjangID, nomorIjazah, excludeID)
	return args.Bool(0), args.Error(1)
}
func (m *KepegawaianPendidikanRepositoryMock) ExistsByNomorIjazahOnly(ctx context.Context, nomorIjazah string, excludeID int64) (bool, error) {
	args := m.Called(ctx, nomorIjazah, excludeID)
	return args.Bool(0), args.Error(1)
}
