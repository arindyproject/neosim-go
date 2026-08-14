package mocks

import (
	userDto "neosim_go/internal/modules/users/dto"
	userModels "neosim_go/internal/modules/users/models"

	"github.com/stretchr/testify/mock"
)

// UserRepositoryMock adalah mock dari userContracts.Repository (modul users).
type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) Create(user *userModels.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepositoryMock) GetByID(id int64) (*userModels.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userModels.User), args.Error(1)
}

func (m *UserRepositoryMock) GetByUsername(username string) (*userModels.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userModels.User), args.Error(1)
}

func (m *UserRepositoryMock) GetByEmail(email string) (*userModels.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userModels.User), args.Error(1)
}

func (m *UserRepositoryMock) List(page, pageSize int, filter *userDto.UserFilter) ([]userModels.User, int64, error) {
	args := m.Called(page, pageSize, filter)
	var users []userModels.User
	if args.Get(0) != nil {
		users = args.Get(0).([]userModels.User)
	}
	return users, args.Get(1).(int64), args.Error(2)
}

func (m *UserRepositoryMock) Update(user *userModels.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepositoryMock) Delete(id int64, deletedBy int64, reason string) error {
	args := m.Called(id, deletedBy, reason)
	return args.Error(0)
}

func (m *UserRepositoryMock) DeletedList(page, pageSize int, filter *userDto.UserDeletedFilter) ([]userModels.User, int64, error) {
	args := m.Called(page, pageSize, filter)
	var users []userModels.User
	if args.Get(0) != nil {
		users = args.Get(0).([]userModels.User)
	}
	return users, args.Get(1).(int64), args.Error(2)
}

func (m *UserRepositoryMock) GetSettings(id int64) ([]userModels.UserSetting, error) {
	args := m.Called(id)
	var settings []userModels.UserSetting
	if args.Get(0) != nil {
		settings = args.Get(0).([]userModels.UserSetting)
	}
	return settings, args.Error(1)
}

func (m *UserRepositoryMock) UpdateSettings(id int64, settings []userModels.UserSetting) error {
	args := m.Called(id, settings)
	return args.Error(0)
}
