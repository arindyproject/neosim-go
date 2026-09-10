package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/kepegawaian/alamat/models"

	"gorm.io/gorm"
)

//go:embed 001_create_kepegawaian_alamats_table.sql
var alamatSQL string

// MigrateKepegawaianAlamat menjalankan GORM auto-migration
func MigrateKepegawaianAlamat(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.KepegawaianAlamat{})
}

// MigrateKepegawaianAlamatWithSQL menjalankan migrasi via raw SQL
func MigrateKepegawaianAlamatWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(alamatSQL)
	if err != nil {
		log.Printf("Error creating kepegawaian_alamats table: %v", err)
		return err
	}
	log.Println("kepegawaian_alamats table migrated successfully")
	return nil
}

// DropKepegawaianAlamatTable menghapus tabel (gunakan dengan hati-hati!)
func DropKepegawaianAlamatTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.KepegawaianAlamat{})
}
