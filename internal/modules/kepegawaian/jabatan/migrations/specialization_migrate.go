package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"

	"gorm.io/gorm"
)

//go:embed 20260918095224_create_kepegawaian_jabatan_specializations_table.sql
var specializationSQL string

// MigrateSpecialization menjalankan GORM auto-migration
func MigrateSpecialization(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.Specialization{})
}

// MigrateSpecializationWithSQL menjalankan migrasi via raw SQL
func MigrateSpecializationWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(specializationSQL)
	if err != nil {
		log.Printf("Error creating kepegawaian_jabatan_specializations table: %v", err)
		return err
	}
	log.Println("kepegawaian_jabatan_specializations table migrated successfully")
	return nil
}

// DropSpecializationTable menghapus tabel (gunakan dengan hati-hati!)
func DropSpecializationTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.Specialization{})
}
