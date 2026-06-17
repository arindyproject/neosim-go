package master

import (
	"neosim_go/internal/modules/master/master/contracts"
	"neosim_go/internal/modules/master/master/handlers"
	"neosim_go/internal/modules/master/master/repositories"
	"neosim_go/internal/modules/master/master/services"
	"neosim_go/internal/shared/utils"

	authContracts "neosim_go/internal/modules/auth/contracts" //auth
	rbacContracts "neosim_go/internal/modules/rbac/contracts" //rbac

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// Module mewakili master/master module
type Module struct {
	db         *gorm.DB
	handler    *handlers.MasterHandler
	jwtManager *utils.JWTManager
	repo       contracts.Repository
	rbacRepo   rbacContracts.RBACRepository //RBAC
}

// NewModule membuat instance module baru dan wire semua layer
func NewModule(
	db *gorm.DB, 
	jwtManager *utils.JWTManager,
	rbacRepo rbacContracts.RBACRepository, //RBAC
	authRepo authContracts.AuthRepository, //AUTH
) *Module {
	repo := repositories.NewMasterRepository(db)
	svc := services.NewMasterService(
		repo,
		rbacRepo,
		authRepo,
	)
	handler := handlers.NewMasterHandler(svc)

	return &Module{
		db:         db,
		handler:    handler,
		jwtManager: jwtManager,
		repo:       repo,
		rbacRepo:   rbacRepo, //RBAC
	}
}

// InitRoutes mendaftarkan routes ke echo instance
func (m *Module) InitRoutes(e *echo.Echo) {
	RegisterRoutes(e, m.handler, m.jwtManager, m.db)
}
