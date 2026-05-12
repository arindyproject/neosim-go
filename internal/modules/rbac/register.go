package rbac

import (
	"database/sql"

	"neosim_go/config"
	"neosim_go/internal/apps"
	"neosim_go/internal/modules/rbac/models"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type registryModule struct {
	db  *gorm.DB
	cfg *config.Config
}

func init() {
	apps.Register(&registryModule{})
}

func (r *registryModule) SetDB(db *gorm.DB)            { r.db = db }
func (r *registryModule) SetConfig(cfg *config.Config) { r.cfg = cfg }

func (r *registryModule) InitRoutes(e *echo.Echo) {
	NewModule(r.db, r.cfg).InitRoutes(e)
}

func (r *registryModule) Models() []interface{} {
	return []interface{}{
		&models.Permission{},
		&models.Role{},
		&models.RolePermission{},
		&models.UserRole{},
		&models.UserPermission{},
	}
}

func (r *registryModule) SeedData(db *gorm.DB) error {
	return nil // seed dihandle oleh cmd/seed
}

func (r *registryModule) MigrateSQL(sqlDB *sql.DB) error {
	return nil // gunakan GORM auto-migrate
}
