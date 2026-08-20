package mocks

import (
	"context"
	"neosim_go/internal/modules/users/dto"
	"neosim_go/internal/modules/users/models"

	"github.com/stretchr/testify/mock"
)

// MockUserRepository mock untuk userContracts.Repository
// dipakai di rbac service karena rbacService inject userRepo
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	return m.Called(user).Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByIDs(ctx context.Context, ids []int64) ([]models.User, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context, page, pageSize int, filter *dto.UserFilter) ([]models.User, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	return m.Called(user).Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64, deletedBy int64, reason string) error {
	return m.Called(id, deletedBy, reason).Error(0)
}

func (m *MockUserRepository) DeletedList(ctx context.Context, page, pageSize int, filter *dto.UserDeletedFilter) ([]models.User, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) GetSettings(ctx context.Context, id int64) ([]models.UserSetting, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.UserSetting), args.Error(1)
}

func (m *MockUserRepository) UpdateSettings(ctx context.Context, id int64, settings []models.UserSetting) error {
	return m.Called(id, settings).Error(0)
}
