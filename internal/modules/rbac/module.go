package rbac

import (
	"neosim_go/config"
	"neosim_go/internal/modules/rbac/contracts"
	"neosim_go/internal/modules/rbac/handlers"
	"neosim_go/internal/modules/rbac/repositories"
	"neosim_go/internal/modules/rbac/services"
	"neosim_go/internal/shared/utils"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// Module mewakili rbac module
type Module struct {
	db         *gorm.DB
	repo       contracts.RBACRepository
	handler    *handlers.RBACHandler
	jwtManager *utils.JWTManager
}

func NewModule(db *gorm.DB, cfg *config.Config) *Module {
	repo := repositories.NewRBACRepository(db)
	svc := services.NewRBACService(repo)
	handler := handlers.NewRBACHandler(svc)
	jwtManager := utils.NewJWTManager(
		cfg.JWTSecret,
		cfg.JWTIssuer,
		cfg.JWTAccessTokenExpMinutes,
		cfg.JWTRefreshTokenExpDays,
	)

	return &Module{
		db:         db,
		repo:       repo,
		handler:    handler,
		jwtManager: jwtManager,
	}
}

func (m *Module) InitRoutes(e *echo.Echo) {
	RegisterRoutes(e, m.handler, m.repo, m.jwtManager, m.db) // ← tambah m.db
}

func (m *Module) GetRepository() contracts.RBACRepository {
	return m.repo
}
