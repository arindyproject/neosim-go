package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/kepegawaian/kualifikasi/models"

	"gorm.io/gorm"
)

//go:embed 001_create_kepegawaian_kualifikasis_table.sql
var kualifikasiSQL string

// MigrateKepegawaianKualifikasi menjalankan GORM auto-migration
func MigrateKepegawaianKualifikasi(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.KepegawaianKualifikasi{})
}

// MigrateKepegawaianKualifikasiWithSQL menjalankan migrasi via raw SQL
func MigrateKepegawaianKualifikasiWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(kualifikasiSQL)
	if err != nil {
		log.Printf("Error creating kepegawaian_kualifikasis table: %v", err)
		return err
	}
	log.Println("kepegawaian_kualifikasis table migrated successfully")
	return nil
}

// DropKepegawaianKualifikasiTable menghapus tabel (gunakan dengan hati-hati!)
func DropKepegawaianKualifikasiTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.KepegawaianKualifikasi{})
}
