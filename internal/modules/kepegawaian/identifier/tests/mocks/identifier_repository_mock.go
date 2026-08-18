package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"neosim_go/internal/modules/kepegawaian/identifier/dto"
	"neosim_go/internal/modules/kepegawaian/identifier/models"
)

// KepegawaianIdentifierRepositoryMock is a mock implementation of contracts.KepegawaianIdentifierRepository.
type KepegawaianIdentifierRepositoryMock struct {
	mock.Mock
}

func (m *KepegawaianIdentifierRepositoryMock) CreateIdentifier(ctx context.Context, item *models.KepegawaianIdentifier) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *KepegawaianIdentifierRepositoryMock) GetIdentifierByID(ctx context.Context, id int64) (*models.KepegawaianIdentifier, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KepegawaianIdentifier), args.Error(1)
}

func (m *KepegawaianIdentifierRepositoryMock) ListIdentifier(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianIdentifierRequest) ([]models.KepegawaianIdentifier, int64, error) {
	args := m.Called(ctx, page, pageSize, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.KepegawaianIdentifier), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianIdentifierRepositoryMock) FindByPegawaiID(ctx context.Context, pegawaiID int64, page, pageSize int) ([]models.KepegawaianIdentifier, int64, error) {
	args := m.Called(ctx, pegawaiID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.KepegawaianIdentifier), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianIdentifierRepositoryMock) FindByPegawaiIDAndTipe(ctx context.Context, pegawaiID, tipeID int64) ([]models.KepegawaianIdentifier, error) {
	args := m.Called(ctx, pegawaiID, tipeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.KepegawaianIdentifier), args.Error(1)
}

func (m *KepegawaianIdentifierRepositoryMock) FindPrimaryByTipe(ctx context.Context, pegawaiID, tipeID int64) (*models.KepegawaianIdentifier, error) {
	args := m.Called(ctx, pegawaiID, tipeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KepegawaianIdentifier), args.Error(1)
}

func (m *KepegawaianIdentifierRepositoryMock) UpdateIdentifier(ctx context.Context, item *models.KepegawaianIdentifier) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *KepegawaianIdentifierRepositoryMock) DeleteIdentifier(ctx context.Context, id int64, deletedBy int64) error {
	args := m.Called(ctx, id, deletedBy)
	return args.Error(0)
}

func (m *KepegawaianIdentifierRepositoryMock) FindExpiringSoonIdentifier(ctx context.Context, days int) ([]models.KepegawaianIdentifier, error) {
	args := m.Called(ctx, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.KepegawaianIdentifier), args.Error(1)
}

func (m *KepegawaianIdentifierRepositoryMock) FindExpiredIdentifier(ctx context.Context) ([]models.KepegawaianIdentifier, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.KepegawaianIdentifier), args.Error(1)
}

func (m *KepegawaianIdentifierRepositoryMock) ExistsIdentifierByID(ctx context.Context, id int64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *KepegawaianIdentifierRepositoryMock) ExistsByNilaiAndTipe(ctx context.Context, tipeID int64, nilai string, excludeID int64) (bool, error) {
	args := m.Called(ctx, tipeID, nilai, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *KepegawaianIdentifierRepositoryMock) UnsetPrimaryByPegawaiIDAndTipe(ctx context.Context, pegawaiID, tipeID int64, updatedBy int64) error {
	args := m.Called(ctx, pegawaiID, tipeID, updatedBy)
	return args.Error(0)
}
