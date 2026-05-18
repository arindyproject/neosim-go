package mocks

import (
	"neosim_go/internal/modules/auth/dto"

	"github.com/stretchr/testify/mock"
)

// MockAuthService adalah mock untuk contracts.AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(req *dto.LoginRequest, ip, userAgent string) (*dto.TokenResponse, error) {
	args := m.Called(req, ip, userAgent)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TokenResponse), args.Error(1)
}

func (m *MockAuthService) Register(req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RegisterResponse), args.Error(1)
}

func (m *MockAuthService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.TokenResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TokenResponse), args.Error(1)
}

func (m *MockAuthService) Logout(req *dto.LogoutRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockAuthService) LogoutAll(userID int64) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockAuthService) ForgotPassword(req *dto.ForgotPasswordRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockAuthService) ResetPassword(req *dto.ResetPasswordRequest) error {
	args := m.Called(req)
	return args.Error(0)
}
