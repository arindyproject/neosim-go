package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/master/departemen/models"

	"gorm.io/gorm"
)

//go:embed 001_create_master_departemens_table.sql
var departemenSQL string

// MigrateMasterDepartemen menjalankan GORM auto-migration
func MigrateMasterDepartemen(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.MasterDepartemen{})
}

// MigrateMasterDepartemenWithSQL menjalankan migrasi via raw SQL
func MigrateMasterDepartemenWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(departemenSQL)
	if err != nil {
		log.Printf("Error creating master_departemens table: %v", err)
		return err
	}
	log.Println("master_departemens table migrated successfully")
	return nil
}

// DropMasterDepartemenTable menghapus tabel (gunakan dengan hati-hati!)
func DropMasterDepartemenTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.MasterDepartemen{})
}
