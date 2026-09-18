package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"

	"gorm.io/gorm"
)

//go:embed 20260918095201_create_kepegawaian_jabatan_job_titles_table.sql
var job_titleSQL string

// MigrateJobTitle menjalankan GORM auto-migration
func MigrateJobTitle(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.JobTitle{})
}

// MigrateJobTitleWithSQL menjalankan migrasi via raw SQL
func MigrateJobTitleWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(job_titleSQL)
	if err != nil {
		log.Printf("Error creating kepegawaian_jabatan_job_titles table: %v", err)
		return err
	}
	log.Println("kepegawaian_jabatan_job_titles table migrated successfully")
	return nil
}

// DropJobTitleTable menghapus tabel (gunakan dengan hati-hati!)
func DropJobTitleTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.JobTitle{})
}
