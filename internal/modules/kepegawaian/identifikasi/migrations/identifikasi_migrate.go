package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/kepegawaian/identifikasi/models"

	"gorm.io/gorm"
)

//go:embed 001_create_kepegawaian_identifikasis_table.sql
var identifikasiSQL string

// MigrateKepegawaianIdentifikasi menjalankan GORM auto-migration
func MigrateKepegawaianIdentifikasi(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.KepegawaianIdentifikasi{})
}

// MigrateKepegawaianIdentifikasiWithSQL menjalankan migrasi via raw SQL
func MigrateKepegawaianIdentifikasiWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(identifikasiSQL)
	if err != nil {
		log.Printf("Error creating kepegawaian_identifikasis table: %v", err)
		return err
	}
	log.Println("kepegawaian_identifikasis table migrated successfully")
	return nil
}

// DropKepegawaianIdentifikasiTable menghapus tabel (gunakan dengan hati-hati!)
func DropKepegawaianIdentifikasiTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.KepegawaianIdentifikasi{})
}
