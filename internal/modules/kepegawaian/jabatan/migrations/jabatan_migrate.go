package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"

	"gorm.io/gorm"
)

//go:embed 001_create_kepegawaian_jabatans_table.sql
var jabatanSQL string

// MigrateKepegawaianJabatan menjalankan GORM auto-migration
func MigrateKepegawaianJabatan(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.KepegawaianJabatan{})
}

// MigrateKepegawaianJabatanWithSQL menjalankan migrasi via raw SQL
func MigrateKepegawaianJabatanWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(jabatanSQL)
	if err != nil {
		log.Printf("Error creating kepegawaian_jabatans table: %v", err)
		return err
	}
	log.Println("kepegawaian_jabatans table migrated successfully")
	return nil
}

// DropKepegawaianJabatanTable menghapus tabel (gunakan dengan hati-hati!)
func DropKepegawaianJabatanTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.KepegawaianJabatan{})
}
