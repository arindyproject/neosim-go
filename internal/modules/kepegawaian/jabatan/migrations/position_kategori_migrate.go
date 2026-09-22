package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"

	"gorm.io/gorm"
)

//go:embed 20260922090240_create_kepegawaian_jabatan_position_kategoris_table.sql
var position_kategoriSQL string

// MigratePositionKategori menjalankan GORM auto-migration
func MigratePositionKategori(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.PositionKategori{})
}

// MigratePositionKategoriWithSQL menjalankan migrasi via raw SQL
func MigratePositionKategoriWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(position_kategoriSQL)
	if err != nil {
		log.Printf("Error creating kepegawaian_jabatan_position_kategoris table: %v", err)
		return err
	}
	log.Println("kepegawaian_jabatan_position_kategoris table migrated successfully")
	return nil
}

// DropPositionKategoriTable menghapus tabel (gunakan dengan hati-hati!)
func DropPositionKategoriTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.PositionKategori{})
}
