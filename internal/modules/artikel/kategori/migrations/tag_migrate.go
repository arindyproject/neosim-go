package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/artikel/kategori/models"

	"gorm.io/gorm"
)

//go:embed 20260904112409_create_artikel_kategori_tags_table.sql
var tagSQL string

// MigrateTag menjalankan GORM auto-migration
func MigrateTag(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.Tag{})
}

// MigrateTagWithSQL menjalankan migrasi via raw SQL
func MigrateTagWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(tagSQL)
	if err != nil {
		log.Printf("Error creating artikel_kategori_tags table: %v", err)
		return err
	}
	log.Println("artikel_kategori_tags table migrated successfully")
	return nil
}

// DropTagTable menghapus tabel (gunakan dengan hati-hati!)
func DropTagTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.Tag{})
}
