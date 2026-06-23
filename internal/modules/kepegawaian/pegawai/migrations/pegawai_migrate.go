package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/kepegawaian/pegawai/models"

	"gorm.io/gorm"
)

//go:embed 001_create_kepegawaian_pegawais_table.sql
var pegawaiSQL string

// MigrateKepegawaianPegawai menjalankan GORM auto-migration
func MigrateKepegawaianPegawai(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.KepegawaianPegawai{})
}

// MigrateKepegawaianPegawaiWithSQL menjalankan migrasi via raw SQL
func MigrateKepegawaianPegawaiWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(pegawaiSQL)
	if err != nil {
		log.Printf("Error creating kepegawaian_pegawais table: %v", err)
		return err
	}
	log.Println("kepegawaian_pegawais table migrated successfully")
	return nil
}

// DropKepegawaianPegawaiTable menghapus tabel (gunakan dengan hati-hati!)
func DropKepegawaianPegawaiTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.KepegawaianPegawai{})
}
