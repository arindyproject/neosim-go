package pendidikan

import (
	"neosim_go/config"
	"neosim_go/internal/modules/kepegawaian/pendidikan/contracts"
	"neosim_go/internal/modules/kepegawaian/pendidikan/handlers"
	"neosim_go/internal/modules/kepegawaian/pendidikan/repositories"
	"neosim_go/internal/modules/kepegawaian/pendidikan/services"
	"neosim_go/internal/shared/utils"

	authContracts "neosim_go/internal/modules/auth/contracts"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
	userContracts "neosim_go/internal/modules/users/contracts"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type Module struct {
	db         *gorm.DB
	handler    *handlers.KepegawaianPendidikanHandler
	jwtManager *utils.JWTManager
	repo       contracts.Repository
	rbacRepo   rbacContracts.RBACRepository
}

func NewModule(
	db *gorm.DB,
	jwtManager *utils.JWTManager,
	rbacRepo rbacContracts.RBACRepository,
	authRepo authContracts.AuthRepository,
	userRepo userContracts.Repository,
	cfg *config.Config,
) *Module {
	repo := repositories.NewKepegawaianPendidikanRepository(db)
	svc := services.NewKepegawaianPendidikanService(repo, rbacRepo, authRepo,userRepo, cfg)
	handler := handlers.NewKepegawaianPendidikanHandler(svc, cfg)

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
