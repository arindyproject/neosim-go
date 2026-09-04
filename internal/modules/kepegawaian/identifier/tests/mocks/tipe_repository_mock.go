package mocks

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/identifier/dto"
	"neosim_go/internal/modules/kepegawaian/identifier/models"
)

// Method di bawah ini ditempelkan ke KepegawaianIdentifierRepositoryMock yang
// sama dengan mock entitas utama (lihat tests/mocks/identifier_repository_mock.go).

func (m *KepegawaianIdentifierRepositoryMock) CreateTipe(ctx context.Context, item *models.Tipe) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianIdentifierRepositoryMock) GetTipeByID(ctx context.Context, id int64) (*models.Tipe, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tipe), args.Error(1)
}

func (m *KepegawaianIdentifierRepositoryMock) GetTipeByLabel(ctx context.Context, label string) (*models.Tipe, error) {
	args := m.Called(label)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tipe), args.Error(1)
}

func (m *KepegawaianIdentifierRepositoryMock) GetTipeByCode(ctx context.Context, code string) (*models.Tipe, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tipe), args.Error(1)
}

func (m *KepegawaianIdentifierRepositoryMock) ListTipe(ctx context.Context, page, pageSize int, filter *dto.FilterTipeRequest) ([]models.Tipe, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.Tipe), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianIdentifierRepositoryMock) UpdateTipe(ctx context.Context, item *models.Tipe) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianIdentifierRepositoryMock) DeleteTipe(ctx context.Context, id int64, deletedBy int64) error {
	args := m.Called(id, deletedBy)
	return args.Error(0)
}
