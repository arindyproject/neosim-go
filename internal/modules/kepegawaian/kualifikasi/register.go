package kualifikasi

import (
	"database/sql"

	"neosim_go/config"
	"neosim_go/internal/apps"
	"neosim_go/internal/modules/kepegawaian/kualifikasi/migrations"
	"neosim_go/internal/modules/kepegawaian/kualifikasi/models"
	"neosim_go/internal/shared/utils"

	authContracts "neosim_go/internal/modules/auth/contracts"
	authRepositories "neosim_go/internal/modules/auth/repositories"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
	rbacRepositories "neosim_go/internal/modules/rbac/repositories"
	userRepositories "neosim_go/internal/modules/users/repositories"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type registryModule struct {
	db       *gorm.DB
	cfg      *config.Config
	rbacRepo rbacContracts.RBACRepository
	authRepo authContracts.AuthRepository
}

func init() {
	apps.Register(&registryModule{})
}

func (r *registryModule) SetDB(db *gorm.DB) {
	r.db = db
	r.rbacRepo = rbacRepositories.NewRBACRepository(db)
	r.authRepo = authRepositories.NewAuthRepository(db)
}

func (r *registryModule) SetConfig(cfg *config.Config) {
	r.cfg = cfg
}

func (r *registryModule) InitRoutes(e *echo.Echo) {
	jwtManager := utils.NewJWTManager(
		r.cfg.JWTSecret,
		r.cfg.JWTIssuer,
		r.cfg.JWTAccessTokenExpMinutes,
		r.cfg.JWTRefreshTokenExpDays,
	)
	userRepo := userRepositories.NewRepository(r.db)
	NewModule(r.db, jwtManager, r.rbacRepo, r.authRepo, userRepo, r.cfg).InitRoutes(e)
}

func (r *registryModule) Models() []interface{} {
	return []interface{}{
		&models.Tipe{},
		&models.KepegawaianKualifikasi{},

		// GEN:ITEM_MODELS
	}
}

func (r *registryModule) SeedData(db *gorm.DB) error {
	return nil
}

func (r *registryModule) MigrateSQL(sqlDB *sql.DB) error {
	if err := migrations.MigrateKepegawaianKualifikasiWithSQL(sqlDB); err != nil {
		return err
	}
	if err := migrations.MigrateTipeWithSQL(sqlDB); err != nil {
		return err
	}
	// GEN:ITEM_MIGRATIONS
	return nil
}
