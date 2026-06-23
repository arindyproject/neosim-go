package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/kepegawaian/kontak/models"

	"gorm.io/gorm"
)

//go:embed 001_create_kepegawaian_kontaks_table.sql
var kontakSQL string

// MigrateKepegawaianKontak menjalankan GORM auto-migration
func MigrateKepegawaianKontak(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.KepegawaianKontak{})
}

// MigrateKepegawaianKontakWithSQL menjalankan migrasi via raw SQL
func MigrateKepegawaianKontakWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(kontakSQL)
	if err != nil {
		log.Printf("Error creating kepegawaian_kontaks table: %v", err)
		return err
	}
	log.Println("kepegawaian_kontaks table migrated successfully")
	return nil
}

// DropKepegawaianKontakTable menghapus tabel (gunakan dengan hati-hati!)
func DropKepegawaianKontakTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.KepegawaianKontak{})
}
