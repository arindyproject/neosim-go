package users

import (
	"database/sql"

	"neosim_go/internal/apps"
	"neosim_go/internal/modules/users/migrations"
	"neosim_go/internal/modules/users/models"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// registryModule adalah implementasi apps.Module untuk users
type registryModule struct {
	db *gorm.DB
}

// init() dipanggil otomatis saat package di-import (via blank import di apps.go)
func init() {
	apps.Register(&registryModule{})
}

// ─── DBInjectable ──────────────────────────────────────────────────────────────

// SetDB menerima db dari registry sebelum InitRoutes dipanggil
func (r *registryModule) SetDB(db *gorm.DB) {
	r.db = db
}

// ─── Routes ────────────────────────────────────────────────────────────────────

func (r *registryModule) InitRoutes(e *echo.Echo) {
	NewModule(r.db).InitRoutes(e)
}

// ─── Migration ─────────────────────────────────────────────────────────────────

func (r *registryModule) Models() []interface{} {
	return []interface{}{
		&models.User{},
		// Tambah model baru di sini jika ada, misal:
		// &models.LoginHistory{},
	}
}

func (r *registryModule) SeedData(db *gorm.DB) error {
	return migrations.SeedDefaultSettings(db)
}

func (r *registryModule) MigrateSQL(sqlDB *sql.DB) error {
	return migrations.MigrateUsersWithSQL(sqlDB)
}
