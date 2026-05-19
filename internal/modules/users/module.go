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

// NewModule creates a new users module instance
func NewModule(db *gorm.DB, jwtManager *utils.JWTManager, rbacRepo rbacContracts.RBACRepository, authRepo authContracts.AuthRepository) *Module {
	// Layer 1: Repository
	repo := repositories.NewRepository(db)

	// Layer 2: Service — inject rbacRepo untuk authorization
	svc := services.NewUserService(repo, rbacRepo, authRepo)

	// Layer 3: Handler
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

// InitRoutes mendaftarkan routes ke echo instance
func (m *Module) InitRoutes(e *echo.Echo) {
	// Inject rbacRepo ke routes untuk RBAC middleware
	RegisterRoutes(e, m.handler, m.rbacRepo, m.jwtManager)
}

// GetRepository returns the repository
func (m *Module) GetRepository() contracts.Repository {
	return m.repo
}

// GetService returns the service
func (m *Module) GetService() contracts.Service {
	return m.service
}

// GetHandler returns the handler
func (m *Module) GetHandler() *handlers.Handler {
	return m.handler
}
