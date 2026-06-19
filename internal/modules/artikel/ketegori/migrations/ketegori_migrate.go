package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/artikel/ketegori/models"

	"gorm.io/gorm"
)

//go:embed 001_create_artikel_ketegoris_table.sql
var ketegoriSQL string

// MigrateArtikelKetegori menjalankan GORM auto-migration
func MigrateArtikelKetegori(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.ArtikelKetegori{})
}

// MigrateArtikelKetegoriWithSQL menjalankan migrasi via raw SQL
func MigrateArtikelKetegoriWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(ketegoriSQL)
	if err != nil {
		log.Printf("Error creating artikel_ketegoris table: %v", err)
		return err
	}
	log.Println("artikel_ketegoris table migrated successfully")
	return nil
}

// DropArtikelKetegoriTable menghapus tabel (gunakan dengan hati-hati!)
func DropArtikelKetegoriTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.ArtikelKetegori{})
}
