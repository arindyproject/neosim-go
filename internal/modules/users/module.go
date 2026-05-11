package users

import (
	"neosim_go/internal/modules/users/contracts"
	"neosim_go/internal/modules/users/handlers"
	"neosim_go/internal/modules/users/repositories"
	"neosim_go/internal/modules/users/services"

	"neosim_go/internal/modules/auth/utils"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// Module represents the users module
type Module struct {
	db         *gorm.DB
	handler    *handlers.Handler
	service    contracts.Service
	repo       contracts.Repository
	jwtManager *utils.JWTManager
}

// NewModule creates a new users module instance
func NewModule(db *gorm.DB, jwtManager *utils.JWTManager) *Module {
	// Create repository
	repo := repositories.NewRepository(db)

	// Create service
	service := services.NewService(repo)

	// Create handler
	handler := handlers.NewHandler(service)

	return &Module{
		db:         db,
		handler:    handler,
		service:    service,
		repo:       repo,
		jwtManager: jwtManager,
	}
}

// InitRoutes registers the module routes to the echo instance
func (m *Module) InitRoutes(e *echo.Echo) {
	RegisterRoutes(e, m.handler, m.jwtManager)
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
