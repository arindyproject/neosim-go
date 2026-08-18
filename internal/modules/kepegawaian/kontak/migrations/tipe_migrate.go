package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/kepegawaian/kontak/models"

	"gorm.io/gorm"
)

//go:embed 20260818111041_create_kepegawaian_kontak_tipes_table.sql
var tipeSQL string

// MigrateTipe menjalankan GORM auto-migration
func MigrateTipe(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.Tipe{})
}

// MigrateTipeWithSQL menjalankan migrasi via raw SQL
func MigrateTipeWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(tipeSQL)
	if err != nil {
		log.Printf("Error creating kepegawaian_kontak_tipes table: %v", err)
		return err
	}
	log.Println("kepegawaian_kontak_tipes table migrated successfully")
	return nil
}

// DropTipeTable menghapus tabel (gunakan dengan hati-hati!)
func DropTipeTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.Tipe{})
}
