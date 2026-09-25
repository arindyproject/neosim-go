package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"

	"gorm.io/gorm"
)

//go:embed 20260925083852_create_kepegawaian_jabatan_specialization_kategoris_table.sql
var specialization_kategoriSQL string

// MigrateSpecializationKategori menjalankan GORM auto-migration
func MigrateSpecializationKategori(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.SpecializationKategori{})
}

// MigrateSpecializationKategoriWithSQL menjalankan migrasi via raw SQL
func MigrateSpecializationKategoriWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(specialization_kategoriSQL)
	if err != nil {
		log.Printf("Error creating kepegawaian_jabatan_specialization_kategoris table: %v", err)
		return err
	}
	log.Println("kepegawaian_jabatan_specialization_kategoris table migrated successfully")
	return nil
}

// DropSpecializationKategoriTable menghapus tabel (gunakan dengan hati-hati!)
func DropSpecializationKategoriTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.SpecializationKategori{})
}
