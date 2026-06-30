package alamat

import (
	"neosim_go/config"
	"neosim_go/internal/modules/master/alamat/contracts"
	"neosim_go/internal/modules/master/alamat/handlers"
	"neosim_go/internal/modules/master/alamat/repositories"
	"neosim_go/internal/modules/master/alamat/services"
	"neosim_go/internal/shared/cache"
	"neosim_go/internal/shared/utils"

	authContracts "neosim_go/internal/modules/auth/contracts" //auth
	rbacContracts "neosim_go/internal/modules/rbac/contracts" //rbac

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// Module mewakili master/alamat module
type Module struct {
	db         *gorm.DB
	handler    *handlers.MasterAlamatHandler
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
	cacheManager *cache.Manager, // <--- Cache Manager
	cfg *config.Config,
) *Module {
	repo := repositories.NewMasterAlamatRepository(db)
	svc := services.NewMasterAlamatService(
		repo,
		rbacRepo,
		authRepo,
		cacheManager, //cache
	)
	handler := handlers.NewMasterAlamatHandler(svc, cfg)

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
