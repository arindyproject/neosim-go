package jabatan

import (
	"neosim_go/internal/modules/kepegawaian/jabatan/contracts"
	"neosim_go/internal/modules/kepegawaian/jabatan/handlers"
	"neosim_go/internal/modules/kepegawaian/jabatan/repositories"
	"neosim_go/internal/modules/kepegawaian/jabatan/services"
	"neosim_go/internal/shared/utils"

	authContracts "neosim_go/internal/modules/auth/contracts"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type Module struct {
	db         *gorm.DB
	handler    *handlers.KepegawaianJabatanHandler
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
	repo := repositories.NewKepegawaianJabatanRepository(db)
	svc := services.NewKepegawaianJabatanService(repo, rbacRepo, authRepo)
	handler := handlers.NewKepegawaianJabatanHandler(svc)

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
