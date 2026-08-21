package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/kepegawaian/pendidikan/models"

	"gorm.io/gorm"
)

//go:embed 001_create_kepegawaian_pendidikans_table.sql
var pendidikanSQL string

//go:embed 20260821090000_add_nomor_ijazah_to_pendidikans.sql
var nomorIjazahSQL string

// MigrateKepegawaianPendidikan menjalankan GORM auto-migration
func MigrateKepegawaianPendidikan(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.KepegawaianPendidikan{})
}

// MigrateKepegawaianPendidikanWithSQL menjalankan migrasi via raw SQL
func MigrateKepegawaianPendidikanWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(pendidikanSQL)
	if err != nil {
		log.Printf("Error creating kepegawaian_pendidikans table: %v", err)
		return err
	}
	if _, err := sqlDB.Exec(nomorIjazahSQL); err != nil {
		log.Printf("Error adding nomor_ijazah to kepegawaian_pendidikans: %v", err)
		return err
	}
	log.Println("kepegawaian_pendidikans table migrated successfully")
	return nil
}

// DropKepegawaianPendidikanTable menghapus tabel (gunakan dengan hati-hati!)
func DropKepegawaianPendidikanTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.KepegawaianPendidikan{})
}
