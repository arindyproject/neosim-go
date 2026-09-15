package alamat

import (
	"neosim_go/config"
	"neosim_go/internal/modules/kepegawaian/alamat/contracts"
	"neosim_go/internal/modules/kepegawaian/alamat/handlers"
	"neosim_go/internal/modules/kepegawaian/alamat/repositories"
	"neosim_go/internal/modules/kepegawaian/alamat/services"
	"neosim_go/internal/shared/cache"
	"neosim_go/internal/shared/utils"

	authContracts "neosim_go/internal/modules/auth/contracts"
	pegawaiContracts "neosim_go/internal/modules/kepegawaian/pegawai/contracts"
	masterAlamatContracts "neosim_go/internal/modules/master/alamat/contracts"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
	userContracts "neosim_go/internal/modules/users/contracts"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type Module struct {
	db         *gorm.DB
	handler    *handlers.KepegawaianAlamatHandler
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
	pegawaiRepo pegawaiContracts.Repository,
	masterAlamatRepo masterAlamatContracts.Repository,
	cfg *config.Config,
	cacheManager *cache.Manager, // <--- Cache Manager
) *Module {
	repo := repositories.NewKepegawaianAlamatRepository(db)
	svc := services.NewKepegawaianAlamatService(repo, rbacRepo, authRepo, userRepo, pegawaiRepo, masterAlamatRepo, cfg, cacheManager)
	handler := handlers.NewKepegawaianAlamatHandler(svc, cfg)

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
