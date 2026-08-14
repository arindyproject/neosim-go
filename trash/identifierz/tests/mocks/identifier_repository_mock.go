package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"neosim_go/internal/modules/kepegawaian/identifier/dto"
	"neosim_go/internal/modules/kepegawaian/identifier/models"
)

// KepegawaianIdentifierRepositoryMock adalah mock implementation dari contracts.Repository
type KepegawaianIdentifierRepositoryMock struct {
	mock.Mock
}

func (m *KepegawaianIdentifierRepositoryMock) Create(ctx context.Context, item *models.KepegawaianIdentifier) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *KepegawaianIdentifierRepositoryMock) FindByID(ctx context.Context, id int64) (*models.KepegawaianIdentifier, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KepegawaianIdentifier), args.Error(1)
}

func (m *KepegawaianIdentifierRepositoryMock) ExistsByID(ctx context.Context, id int64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *KepegawaianIdentifierRepositoryMock) FindAll(ctx context.Context, filter dto.FilterKepegawaianIdentifierRequest, page, pageSize int) ([]models.KepegawaianIdentifier, int64, error) {
	args := m.Called(ctx, filter, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.KepegawaianIdentifier), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianIdentifierRepositoryMock) FindByPegawaiD(ctx context.Context, pegawaiID int64) ([]models.KepegawaianIdentifier, error) {
	args := m.Called(ctx, pegawaiID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.KepegawaianIdentifier), args.Error(1)
}

func (m *KepegawaianIdentifierRepositoryMock) FindByPegawaiDAndTipe(ctx context.Context, pegawaiID int64, tipe models.IdentifierType) ([]models.KepegawaianIdentifier, error) {
	args := m.Called(ctx, pegawaiID, tipe)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.KepegawaianIdentifier), args.Error(1)
}

// ── Method Tambahan: FindPrimaryByTipe ───────────────────────────────────────

func (m *KepegawaianIdentifierRepositoryMock) FindPrimaryByTipe(ctx context.Context, pegawaiID int64, tipe models.IdentifierType) (*models.KepegawaianIdentifier, error) {
	args := m.Called(ctx, pegawaiID, tipe)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KepegawaianIdentifier), args.Error(1)
}

// ──────────────────────────────────────────────────────────────────────────────

func (m *KepegawaianIdentifierRepositoryMock) Update(ctx context.Context, item *models.KepegawaianIdentifier) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *KepegawaianIdentifierRepositoryMock) Delete(ctx context.Context, id int64, deletedBy int64) error {
	args := m.Called(ctx, id, deletedBy)
	return args.Error(0)
}

func (m *KepegawaianIdentifierRepositoryMock) ExistsByNilaiAndTipe(ctx context.Context, tipe models.IdentifierType, nilai string, excludeID int64) (bool, error) {
	args := m.Called(ctx, tipe, nilai, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *KepegawaianIdentifierRepositoryMock) UnsetPrimaryByPegawaiDAndTipe(ctx context.Context, pegawaiID int64, tipe models.IdentifierType, updatedBy int64) error {
	args := m.Called(ctx, pegawaiID, tipe, updatedBy)
	return args.Error(0)
}

func (m *KepegawaianIdentifierRepositoryMock) FindExpiringSoon(ctx context.Context, days int) ([]models.KepegawaianIdentifier, error) {
	args := m.Called(ctx, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.KepegawaianIdentifier), args.Error(1)
}

func (m *KepegawaianIdentifierRepositoryMock) FindExpired(ctx context.Context) ([]models.KepegawaianIdentifier, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.KepegawaianIdentifier), args.Error(1)
}
