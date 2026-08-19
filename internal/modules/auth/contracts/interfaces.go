package contracts

import (
	"context"
	"neosim_go/internal/modules/auth/dto"
	"neosim_go/internal/modules/auth/models"
)

// ─── Repository ────────────────────────────────────────────────────────────────

// AuthRepository mendefinisikan operasi database untuk auth
type AuthRepository interface {
	// Auth Token
	SaveToken(ctx context.Context, token *models.AuthToken) error
	GetTokenByJTI(ctx context.Context, jti string) (*models.AuthToken, error)
	BlacklistToken(ctx context.Context, jti string) error
	BlacklistAllUserTokens(ctx context.Context, userID int64) error
	CountActiveTokens(uctx context.Context, serID int64) (int64, error)

	// Login History
	SaveLoginHistory(ctx context.Context, history *models.LoginHistory) error
	GetUserLoginHistories(ctx context.Context, userID int64, limit int) ([]models.LoginHistory, error)

	// Password History
	SavePasswordHistory(ctx context.Context, history *models.PasswordHistory) error
	GetPasswordHistories(ctx context.Context, userID int64, limit int) ([]models.PasswordHistory, error)
}

// ─── Service ───────────────────────────────────────────────────────────────────

// AuthService mendefinisikan business logic untuk auth
type AuthService interface {
	// Auth
	Login(ctx context.Context, req *dto.LoginRequest, ip, userAgent string) (*dto.TokenResponse, error)
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.TokenResponse, error)

	// Logout
	Logout(ctx context.Context, req *dto.LogoutRequest) error // logout device saat ini
	LogoutAll(ctx context.Context, userID int64) error        // logout semua device

	// Password
	ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error
}
