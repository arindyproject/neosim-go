package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/artikel/kategori/models"

	"gorm.io/gorm"
)

//go:embed 001_create_artikel_kategoris_table.sql
var kategoriSQL string

// MigrateArtikelKategori menjalankan GORM auto-migration
func MigrateArtikelKategori(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.ArtikelKategori{})
}

// MigrateArtikelKategoriWithSQL menjalankan migrasi via raw SQL
func MigrateArtikelKategoriWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(kategoriSQL)
	if err != nil {
		log.Printf("Error creating artikel_kategoris table: %v", err)
		return err
	}
	log.Println("artikel_kategoris table migrated successfully")
	return nil
}

// DropArtikelKategoriTable menghapus tabel (gunakan dengan hati-hati!)
func DropArtikelKategoriTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.ArtikelKategori{})
}
