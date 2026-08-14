package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/kepegawaian/identifier/models"

	"gorm.io/gorm"
)

//go:embed 20260814105037_create_kepegawaian_identifier_tipes_table.sql
var tipeSQL string

// MigrateTipe menjalankan GORM auto-migration
func MigrateTipe(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.Tipe{})
}

// MigrateTipeWithSQL menjalankan migrasi via raw SQL
func MigrateTipeWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(tipeSQL)
	if err != nil {
		log.Printf("Error creating kepegawaian_identifier_tipes table: %v", err)
		return err
	}
	log.Println("kepegawaian_identifier_tipes table migrated successfully")
	return nil
}

// DropTipeTable menghapus tabel (gunakan dengan hati-hati!)
func DropTipeTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.Tipe{})
}
