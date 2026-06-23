package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/kepegawaian/lampiran/models"

	"gorm.io/gorm"
)

//go:embed 001_create_kepegawaian_lampirans_table.sql
var lampiranSQL string

// MigrateKepegawaianLampiran menjalankan GORM auto-migration
func MigrateKepegawaianLampiran(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.KepegawaianLampiran{})
}

// MigrateKepegawaianLampiranWithSQL menjalankan migrasi via raw SQL
func MigrateKepegawaianLampiranWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(lampiranSQL)
	if err != nil {
		log.Printf("Error creating kepegawaian_lampirans table: %v", err)
		return err
	}
	log.Println("kepegawaian_lampirans table migrated successfully")
	return nil
}

// DropKepegawaianLampiranTable menghapus tabel (gunakan dengan hati-hati!)
func DropKepegawaianLampiranTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.KepegawaianLampiran{})
}
