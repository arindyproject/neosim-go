package ketegori

import (
	"neosim_go/internal/modules/artikel/ketegori/contracts"
	"neosim_go/internal/modules/artikel/ketegori/handlers"
	"neosim_go/internal/modules/artikel/ketegori/repositories"
	"neosim_go/internal/modules/artikel/ketegori/services"
	"neosim_go/internal/shared/utils"

	authContracts "neosim_go/internal/modules/auth/contracts" //auth
	rbacContracts "neosim_go/internal/modules/rbac/contracts" //rbac

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// Module mewakili artikel/ketegori module
type Module struct {
	db         *gorm.DB
	handler    *handlers.ArtikelKetegoriHandler
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
	repo := repositories.NewArtikelKetegoriRepository(db)
	svc := services.NewArtikelKetegoriService(
		repo,
		rbacRepo,
		authRepo,
	)
	handler := handlers.NewArtikelKetegoriHandler(svc)

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
