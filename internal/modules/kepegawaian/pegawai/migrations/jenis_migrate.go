package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/kepegawaian/pegawai/models"

	"gorm.io/gorm"
)

//go:embed 20260929135921_create_kepegawaian_pegawai_jeniss_table.sql
var jenisSQL string

// MigrateJenis menjalankan GORM auto-migration
func MigrateJenis(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.Jenis{})
}

// MigrateJenisWithSQL menjalankan migrasi via raw SQL
func MigrateJenisWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(jenisSQL)
	if err != nil {
		log.Printf("Error creating kepegawaian_pegawai_jeniss table: %v", err)
		return err
	}
	log.Println("kepegawaian_pegawai_jeniss table migrated successfully")
	return nil
}

// DropJenisTable menghapus tabel (gunakan dengan hati-hati!)
func DropJenisTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.Jenis{})
}
