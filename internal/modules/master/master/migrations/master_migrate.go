package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/master/master/models"

	"gorm.io/gorm"
)

//go:embed 001_create_masters_table.sql
var masterSQL string

// MigrateMaster menjalankan GORM auto-migration
func MigrateMaster(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.Master{})
}

// MigrateMasterWithSQL menjalankan migrasi via raw SQL
func MigrateMasterWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(masterSQL)
	if err != nil {
		log.Printf("Error creating masters table: %v", err)
		return err
	}
	log.Println("masters table migrated successfully")
	return nil
}

// DropMasterTable menghapus tabel (gunakan dengan hati-hati!)
func DropMasterTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.Master{})
}
