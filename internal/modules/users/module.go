package users

import (
	authContracts "neosim_go/internal/modules/auth/contracts"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
	"neosim_go/internal/modules/users/contracts"
	"neosim_go/internal/modules/users/handlers"
	"neosim_go/internal/modules/users/repositories"
	"neosim_go/internal/modules/users/services"
	"neosim_go/internal/shared/utils"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// Module represents the users module
type Module struct {
	db         *gorm.DB
	handler    *handlers.Handler
	service    contracts.Service
	repo       contracts.Repository
	rbacRepo   rbacContracts.RBACRepository
	jwtManager *utils.JWTManager
}

// NewModule membuat instance module dan wire semua layer
func NewModule(db *gorm.DB, jwtManager *utils.JWTManager, rbacRepo rbacContracts.RBACRepository, authRepo authContracts.AuthRepository) *Module {
	repo := repositories.NewRepository(db)
	svc := services.NewUserService(repo, rbacRepo, authRepo) // ← hanya repo + rbacRepo
	handler := handlers.NewHandler(svc)

	return &Module{
		db:         db,
		handler:    handler,
		service:    svc,
		repo:       repo,
		rbacRepo:   rbacRepo,
		jwtManager: jwtManager,
	}
}

// InitRoutes mendaftarkan routes — db diteruskan ke JWTMiddleware untuk realtime check
func (m *Module) InitRoutes(e *echo.Echo) {
	RegisterRoutes(e, m.handler, m.rbacRepo, m.jwtManager, m.db)
}

func (m *Module) GetRepository() contracts.Repository { return m.repo }
func (m *Module) GetService() contracts.Service       { return m.service }
func (m *Module) GetHandler() *handlers.Handler       { return m.handler }
