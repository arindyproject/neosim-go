package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"neosim_go/internal/modules/master/departemen/dto"
	"neosim_go/internal/modules/master/departemen/models"
)

// MasterDepartemenRepositoryMock is a mock implementation of contracts.Repository.
// Ketika item ditambahkan (mode add-item), method mock untuk item tersebut
// ditempelkan ke struct INI JUGA (mis. tests/mocks/tag_repository_mock.go),
// bukan membuat mock struct baru.
type MasterDepartemenRepositoryMock struct {
	mock.Mock
}

func (m *MasterDepartemenRepositoryMock) CreateDepartemen(ctx context.Context, item *models.MasterDepartemen) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterDepartemenRepositoryMock) GetDepartemenByID(ctx context.Context, id int64) (*models.MasterDepartemen, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasterDepartemen), args.Error(1)
}

func (m *MasterDepartemenRepositoryMock) GetByIDs(ctx context.Context, ids []int64) ([]models.MasterDepartemen, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.MasterDepartemen), args.Error(1)
}

func (m *MasterDepartemenRepositoryMock) ListDepartemen(ctx context.Context, page, pageSize int, filter *dto.FilterMasterDepartemenRequest) ([]models.MasterDepartemen, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.MasterDepartemen), args.Get(1).(int64), args.Error(2)
}

func (m *MasterDepartemenRepositoryMock) UpdateDepartemen(ctx context.Context, item *models.MasterDepartemen) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MasterDepartemenRepositoryMock) DeleteDepartemen(ctx context.Context, id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
