package users

import (
	"database/sql"

	"neosim_go/config"
	"neosim_go/internal/apps"
	"neosim_go/internal/modules/auth/utils"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
	rbacRepositories "neosim_go/internal/modules/rbac/repositories"
	"neosim_go/internal/modules/users/migrations"
	"neosim_go/internal/modules/users/models"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type registryModule struct {
	db       *gorm.DB
	cfg      *config.Config
	rbacRepo rbacContracts.RBACRepository
}

func init() {
	apps.Register(&registryModule{})
}

// ─── Injections ────────────────────────────────────────────────────────────────

func (r *registryModule) SetDB(db *gorm.DB) {
	r.db = db
	// Build rbacRepo saat DB tersedia
	// rbacRepo tidak punya table sendiri di module ini, reuse dari rbac module
	r.rbacRepo = rbacRepositories.NewRBACRepository(db)
}

func (r *registryModule) SetConfig(cfg *config.Config) {
	r.cfg = cfg
}

// ─── Routes ────────────────────────────────────────────────────────────────────

func (r *registryModule) InitRoutes(e *echo.Echo) {
	jwtManager := utils.NewJWTManager(
		r.cfg.JWTSecret,
		r.cfg.JWTIssuer,
		r.cfg.JWTAccessTokenExpMinutes,
		r.cfg.JWTRefreshTokenExpDays,
	)
	NewModule(r.db, jwtManager, r.rbacRepo).InitRoutes(e)
}

// ─── Migration ─────────────────────────────────────────────────────────────────

func (r *registryModule) Models() []interface{} {
	return []interface{}{
		&models.User{},
	}
}

func (r *registryModule) SeedData(db *gorm.DB) error {
	return migrations.SeedDefaultSettings(db)
}

func (r *registryModule) MigrateSQL(sqlDB *sql.DB) error {
	return migrations.MigrateUsersWithSQL(sqlDB)
}
