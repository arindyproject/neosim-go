package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/artikel/models"

	"gorm.io/gorm"
)

//go:embed 001_create_artikel_table.sql
var artikelSQL string

// MigrateArtikel menjalankan GORM auto-migration
func MigrateArtikel(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.Artikel{})
}

// MigrateArtikelWithSQL menjalankan migrasi via raw SQL
func MigrateArtikelWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(artikelSQL)
	if err != nil {
		log.Printf("Error creating artikels table: %v", err)
		return err
	}
	log.Println("Artikels table migrated successfully")
	return nil
}

// DropArtikelTable menghapus tabel (gunakan dengan hati-hati!)
func DropArtikelTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.Artikel{})
}
