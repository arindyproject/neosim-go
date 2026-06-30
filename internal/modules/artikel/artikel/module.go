package artikel

import (
	"neosim_go/config"
	"neosim_go/internal/modules/artikel/artikel/contracts"
	"neosim_go/internal/modules/artikel/artikel/handlers"
	"neosim_go/internal/modules/artikel/artikel/repositories"
	"neosim_go/internal/modules/artikel/artikel/services"
	"neosim_go/internal/shared/utils"

	authContracts "neosim_go/internal/modules/auth/contracts"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type Module struct {
	db         *gorm.DB
	handler    *handlers.ArtikelHandler
	jwtManager *utils.JWTManager
	repo       contracts.Repository
	rbacRepo   rbacContracts.RBACRepository
}

func NewModule(
	db *gorm.DB,
	jwtManager *utils.JWTManager,
	rbacRepo rbacContracts.RBACRepository,
	authRepo authContracts.AuthRepository,
	cfg *config.Config,
) *Module {
	repo := repositories.NewArtikelRepository(db)
	svc := services.NewArtikelService(repo, rbacRepo, authRepo, cfg)
	handler := handlers.NewArtikelHandler(svc, cfg)

	return &Module{
		db:         db,
		handler:    handler,
		jwtManager: jwtManager,
		repo:       repo,
		rbacRepo:   rbacRepo,
	}
}

func (m *Module) InitRoutes(e *echo.Echo) {
	RegisterRoutes(e, m.handler, m.jwtManager, m.db)
}
