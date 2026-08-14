package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/kepegawaian/identifier/models"

	"gorm.io/gorm"
)

//go:embed 001_create_kepegawaian_identifiers_table.sql
var identifierSQL string

// MigrateKepegawaianIdentifier menjalankan GORM auto-migration
func MigrateKepegawaianIdentifier(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.KepegawaianIdentifier{})
}

// MigrateKepegawaianIdentifierWithSQL menjalankan migrasi via raw SQL
func MigrateKepegawaianIdentifierWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(identifierSQL)
	if err != nil {
		log.Printf("Error creating kepegawaian_identifiers table: %v", err)
		return err
	}
	log.Println("kepegawaian_identifiers table migrated successfully")
	return nil
}

// DropKepegawaianIdentifierTable menghapus tabel (gunakan dengan hati-hati!)
func DropKepegawaianIdentifierTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.KepegawaianIdentifier{})
}
