package contracts

import (
	"neosim_go/internal/modules/auth/dto"
	"neosim_go/internal/modules/auth/models"
)

// ─── Repository ────────────────────────────────────────────────────────────────

// AuthRepository mendefinisikan operasi database untuk auth
type AuthRepository interface {
	// Auth Token
	SaveToken(token *models.AuthToken) error
	GetTokenByJTI(jti string) (*models.AuthToken, error)
	BlacklistToken(jti string) error
	BlacklistAllUserTokens(userID int64) error
	CountActiveTokens(userID int64) (int64, error)

	// Login History
	SaveLoginHistory(history *models.LoginHistory) error

	// Password History
	SavePasswordHistory(history *models.PasswordHistory) error
	GetPasswordHistories(userID int64, limit int) ([]models.PasswordHistory, error)
}

// ─── Service ───────────────────────────────────────────────────────────────────

// AuthService mendefinisikan business logic untuk auth
type AuthService interface {
	Login(req *dto.LoginRequest, ip, userAgent string) (*dto.TokenResponse, error)
	Register(req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	RefreshToken(req *dto.RefreshTokenRequest) (*dto.TokenResponse, error)
	ForgotPassword(req *dto.ForgotPasswordRequest) error
	ResetPassword(req *dto.ResetPasswordRequest) error
}
