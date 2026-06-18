package alamat

import (
	"database/sql"
	"log"

	"neosim_go/config"
	"neosim_go/internal/apps"
	"neosim_go/internal/modules/master/alamat/models"
	"neosim_go/internal/shared/cache"
	"neosim_go/internal/shared/utils"

	authContracts "neosim_go/internal/modules/auth/contracts"
	authRepositories "neosim_go/internal/modules/auth/repositories"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
	rbacRepositories "neosim_go/internal/modules/rbac/repositories"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type registryModule struct {
	db           *gorm.DB
	cfg          *config.Config
	rbacRepo     rbacContracts.RBACRepository //RBAC
	authRepo     authContracts.AuthRepository //AUTH
	cacheManager *cache.Manager               //CACHE
}

// init dipanggil otomatis saat package di-import (blank import)
func init() {
	apps.Register(&registryModule{})
}

func (r *registryModule) SetDB(db *gorm.DB) {
	r.db = db
	r.rbacRepo = rbacRepositories.NewRBACRepository(db) //RBAC
	r.authRepo = authRepositories.NewAuthRepository(db) //AUTH
}
func (r *registryModule) SetConfig(cfg *config.Config) {
	r.cfg = cfg

	// ─── Inisialisasi Cache Manager ─────────────────────────────────────────
	if cfg.CacheMasterAlamat {
		client, err := cfg.ConnectRedis()
		if err != nil {
			log.Printf("⚠️ Warning: Gagal koneksi ke Redis untuk cache Master Alamat: %v", err)
			r.cacheManager = cache.NewManager(nil, false, 0)
		} else {
			r.cacheManager = cache.NewManager(client, true, cfg.CacheMasterAlamatTTLDay)
			log.Println("✅ Cache Manager untuk Master Alamat berhasil diinisialisasi")
		}
	} else {
		r.cacheManager = cache.NewManager(nil, false, 0)
	}
}

func (r *registryModule) InitRoutes(e *echo.Echo) {
	jwtManager := utils.NewJWTManager(
		r.cfg.JWTSecret,
		r.cfg.JWTIssuer,
		r.cfg.JWTAccessTokenExpMinutes,
		r.cfg.JWTRefreshTokenExpDays,
	)
	NewModule(
		r.db,
		jwtManager,
		r.rbacRepo,
		r.authRepo,
		r.cacheManager, // <---  Cache Redis
	).InitRoutes(e)
}

func (r *registryModule) Models() []interface{} {
	return []interface{}{
		&models.MasterAlamatNegara{},
		&models.MasterAlamatProvinsi{},
		&models.MasterAlamatKotaKabupaten{},
		&models.MasterAlamatKecamatan{},
		&models.MasterAlamatKelurahanDesa{},
	}
}

func (r *registryModule) SeedData(db *gorm.DB) error {
	return nil
}

func (r *registryModule) MigrateSQL(sqlDB *sql.DB) error {
	return nil
}
