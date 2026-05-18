package auth

import (
	"neosim_go/config"
	"neosim_go/internal/modules/auth/handlers"
	"neosim_go/internal/modules/auth/repositories"
	"neosim_go/internal/modules/auth/services"
	"neosim_go/internal/modules/auth/utils"
	userContracts "neosim_go/internal/modules/users/contracts"
	userRepositories "neosim_go/internal/modules/users/repositories"

	"github.com/labstack/echo/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Module mewakili auth module
type Module struct {
	db         *gorm.DB
	redis      *redis.Client
	handler    *handlers.AuthHandler
	jwtManager *utils.JWTManager
}

// NewModule membuat instance module dan wire semua layer
func NewModule(db *gorm.DB, redisClient *redis.Client, cfg *config.Config) *Module {
	// Auth repository
	authRepo := repositories.NewAuthRepository(db)

	// Reuse users repository yang sudah ada
	var userRepo userContracts.Repository = userRepositories.NewRepository(db)

	// Build mailer jika SMTP dikonfigurasi
	var mailer *utils.Mailer
	if cfg.SMTPHost != "" {
		mailer = utils.NewMailer(utils.MailConfig{
			Host:     cfg.SMTPHost,
			Port:     cfg.SMTPPort,
			Username: cfg.SMTPUsername,
			Password: cfg.SMTPPassword,
			From:     cfg.SMTPFrom,
			FromName: cfg.SMTPFromName,
		})
	}

	// Build service config dari Config struct
	svcCfg := services.AuthServiceConfig{
		JWTManager: utils.NewJWTManager(
			cfg.JWTSecret,
			cfg.JWTIssuer,
			cfg.JWTAccessTokenExpMinutes,
			cfg.JWTRefreshTokenExpDays,
		),
		LoginMaxAttempts:             cfg.LoginMaxAttempts,
		LoginLockDurationMinutes:     cfg.LoginLockDurationMinutes,
		MaxConcurrentSessions:        cfg.MaxConcurrentSessions,
		RateLimitLoginPerIPPerMinute: cfg.RateLimitLoginPerIPPerMinute,
		PasswordPolicy: &utils.PasswordPolicy{
			MinLength:        cfg.PasswordMinLength,
			RequireUppercase: cfg.PasswordRequireUppercase,
			RequireNumber:    cfg.PasswordRequireNumber,
			RequireSymbol:    cfg.PasswordRequireSymbol,
		},
		PasswordHistoryCount:     cfg.PasswordHistoryCount,
		IsRegistrationActive:     cfg.IsRegistrationActive,
		AutoActiveUser:           cfg.AutoActiveUser,
		MailResetTokenExpMinutes: cfg.MailResetTokenExpMinutes,
		AppFrontendURL:           cfg.AppFrontendURL,
		Mailer:                   mailer,
	}

	svc := services.NewAuthService(authRepo, userRepo, redisClient, svcCfg)
	handler := handlers.NewAuthHandler(svc)
	jwtManager := utils.NewJWTManager(
		cfg.JWTSecret,
		cfg.JWTIssuer,
		cfg.JWTAccessTokenExpMinutes,
		cfg.JWTRefreshTokenExpDays,
	)

	return &Module{
		db:         db,
		redis:      redisClient,
		handler:    handler,
		jwtManager: jwtManager,
	}
}

// InitRoutes mendaftarkan routes ke echo instance
func (m *Module) InitRoutes(e *echo.Echo) {
	RegisterRoutes(e, m.handler, m.jwtManager)
}
