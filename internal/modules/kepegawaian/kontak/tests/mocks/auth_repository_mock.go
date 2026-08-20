package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	authModels "neosim_go/internal/modules/auth/models"
)

// AuthRepositoryMock is a mock implementation of authContracts.AuthRepository
type AuthRepositoryMock struct {
	mock.Mock
}

func (m *AuthRepositoryMock) SaveToken(ctx context.Context, token *authModels.AuthToken) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *AuthRepositoryMock) GetTokenByJTI(ctx context.Context, jti string) (*authModels.AuthToken, error) {
	args := m.Called(jti)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authModels.AuthToken), args.Error(1)
}

func (m *AuthRepositoryMock) BlacklistToken(ctx context.Context, jti string) error {
	args := m.Called(jti)
	return args.Error(0)
}

func (m *AuthRepositoryMock) BlacklistAllUserTokens(ctx context.Context, userID int64) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *AuthRepositoryMock) CountActiveTokens(ctx context.Context, userID int64) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *AuthRepositoryMock) SaveLoginHistory(ctx context.Context, history *authModels.LoginHistory) error {
	args := m.Called(history)
	return args.Error(0)
}

func (m *AuthRepositoryMock) GetUserLoginHistories(ctx context.Context, userID int64, limit int) ([]authModels.LoginHistory, error) {
	args := m.Called(userID, limit)
	return args.Get(0).([]authModels.LoginHistory), args.Error(1)
}

func (m *AuthRepositoryMock) SavePasswordHistory(ctx context.Context, history *authModels.PasswordHistory) error {
	args := m.Called(history)
	return args.Error(0)
}

func (m *AuthRepositoryMock) GetPasswordHistories(ctx context.Context, userID int64, limit int) ([]authModels.PasswordHistory, error) {
	args := m.Called(userID, limit)
	return args.Get(0).([]authModels.PasswordHistory), args.Error(1)
}
