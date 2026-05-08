package apps

import (
	"database/sql"
	"log"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// Module interface yang harus diimplementasi tiap module
type Module interface {
	Models() []interface{}
	SeedData(db *gorm.DB) error
	MigrateSQL(sqlDB *sql.DB) error
	InitRoutes(e *echo.Echo)
}

// DBInjectable optional interface untuk module yang butuh DB
type DBInjectable interface {
	SetDB(db *gorm.DB)
}

var registeredModules []Module

// Register mendaftarkan module ke registry
func Register(m Module) {
	registeredModules = append(registeredModules, m)
}

// ─── Routes ────────────────────────────────────────────────────────────────────

// InitAllRoutes menginisialisasi semua routes dari semua module
func InitAllRoutes(e *echo.Echo) {
	for _, m := range registeredModules {
		m.InitRoutes(e)
	}
}

// ─── DB Injection ──────────────────────────────────────────────────────────────

// InjectDB menyuntikkan db ke semua module yang mengimplementasi DBInjectable
func InjectDB(db *gorm.DB) {
	for _, m := range registeredModules {
		if injectable, ok := m.(DBInjectable); ok {
			injectable.SetDB(db)
		}
	}
}

// ─── Migration ─────────────────────────────────────────────────────────────────

// AllModels mengumpulkan semua model dari semua module
func AllModels() []interface{} {
	var all []interface{}
	for _, m := range registeredModules {
		all = append(all, m.Models()...)
	}
	return all
}

// DropAll menghapus semua tabel dalam urutan reverse (FK constraint)
func DropAll(db *gorm.DB) error {
	models := AllModels()
	reversed := make([]interface{}, len(models))
	for i, m := range models {
		reversed[len(models)-1-i] = m
	}
	return db.Migrator().DropTable(reversed...)
}

// MigrateAll menjalankan GORM auto-migration semua module
func MigrateAll(db *gorm.DB) error {
	for _, m := range registeredModules {
		if err := db.AutoMigrate(m.Models()...); err != nil {
			return err
		}
	}
	return nil
}

// SeedAll menjalankan seed data semua module
func SeedAll(db *gorm.DB) {
	for _, m := range registeredModules {
		if err := m.SeedData(db); err != nil {
			log.Printf("Warning: seed failed for %T: %v", m, err)
		}
	}
}

// MigrateAllSQL menjalankan SQL migration semua module
func MigrateAllSQL(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	for _, m := range registeredModules {
		if err := m.MigrateSQL(sqlDB); err != nil {
			return err
		}
	}
	return nil
}
