package departemen

import (
	"neosim_go/internal/modules/master/departemen/contracts"
	"neosim_go/internal/modules/master/departemen/handlers"
	"neosim_go/internal/modules/master/departemen/repositories"
	"neosim_go/internal/modules/master/departemen/services"
	"neosim_go/internal/shared/utils"

	authContracts "neosim_go/internal/modules/auth/contracts"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type Module struct {
	db         *gorm.DB
	handler    *handlers.MasterDepartemenHandler
	jwtManager *utils.JWTManager
	repo       contracts.Repository
	rbacRepo   rbacContracts.RBACRepository
}

func NewModule(
	db *gorm.DB,
	jwtManager *utils.JWTManager,
	rbacRepo rbacContracts.RBACRepository,
	authRepo authContracts.AuthRepository,
) *Module {
	repo := repositories.NewMasterDepartemenRepository(db)
	svc := services.NewMasterDepartemenService(repo, rbacRepo, authRepo)
	handler := handlers.NewMasterDepartemenHandler(svc)

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
