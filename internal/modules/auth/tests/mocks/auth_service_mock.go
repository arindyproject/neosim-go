package mocks

import (
	"context"
	"neosim_go/internal/modules/auth/dto"

	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(
	ctx context.Context,
	req *dto.LoginRequest,
	ip,
	userAgent string,
) (*dto.TokenResponse, error) {

	args := m.Called(req, ip, userAgent)

	var resp *dto.TokenResponse
	if v := args.Get(0); v != nil {
		resp = v.(*dto.TokenResponse)
	}

	return resp, args.Error(1)
}

func (m *MockAuthService) Register(
	ctx context.Context,
	req *dto.RegisterRequest,
) (*dto.RegisterResponse, error) {

	args := m.Called(req)

	var resp *dto.RegisterResponse
	if v := args.Get(0); v != nil {
		resp = v.(*dto.RegisterResponse)
	}

	return resp, args.Error(1)
}

func (m *MockAuthService) RefreshToken(
	ctx context.Context,
	req *dto.RefreshTokenRequest,
) (*dto.TokenResponse, error) {

	args := m.Called(req)

	var resp *dto.TokenResponse
	if v := args.Get(0); v != nil {
		resp = v.(*dto.TokenResponse)
	}

	return resp, args.Error(1)
}

func (m *MockAuthService) Logout(
	ctx context.Context,
	req *dto.LogoutRequest,
) error {

	args := m.Called(req)
	return args.Error(0)
}

func (m *MockAuthService) LogoutAll(
	ctx context.Context,
	userID int64,
) error {

	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockAuthService) ForgotPassword(
	ctx context.Context,
	req *dto.ForgotPasswordRequest,
) error {

	args := m.Called(req)
	return args.Error(0)
}

func (m *MockAuthService) ResetPassword(
	ctx context.Context,
	req *dto.ResetPasswordRequest,
) error {

	args := m.Called(req)
	return args.Error(0)
}
