package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"

	"gorm.io/gorm"
)

//go:embed 20260923083332_create_kepegawaian_jabatan_job_title_kategoris_table.sql
var job_title_kategoriSQL string

// MigrateJobTitleKategori menjalankan GORM auto-migration
func MigrateJobTitleKategori(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.JobTitleKategori{})
}

// MigrateJobTitleKategoriWithSQL menjalankan migrasi via raw SQL
func MigrateJobTitleKategoriWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(job_title_kategoriSQL)
	if err != nil {
		log.Printf("Error creating kepegawaian_jabatan_job_title_kategoris table: %v", err)
		return err
	}
	log.Println("kepegawaian_jabatan_job_title_kategoris table migrated successfully")
	return nil
}

// DropJobTitleKategoriTable menghapus tabel (gunakan dengan hati-hati!)
func DropJobTitleKategoriTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.JobTitleKategori{})
}
