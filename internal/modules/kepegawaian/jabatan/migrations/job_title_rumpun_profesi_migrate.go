package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"

	"gorm.io/gorm"
)

//go:embed 20260923083416_create_kepegawaian_jabatan_job_title_rumpun_profesis_table.sql
var job_title_rumpun_profesiSQL string

// MigrateJobTitleRumpunProfesi menjalankan GORM auto-migration
func MigrateJobTitleRumpunProfesi(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.JobTitleRumpunProfesi{})
}

// MigrateJobTitleRumpunProfesiWithSQL menjalankan migrasi via raw SQL
func MigrateJobTitleRumpunProfesiWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(job_title_rumpun_profesiSQL)
	if err != nil {
		log.Printf("Error creating kepegawaian_jabatan_job_title_rumpun_profesis table: %v", err)
		return err
	}
	log.Println("kepegawaian_jabatan_job_title_rumpun_profesis table migrated successfully")
	return nil
}

// DropJobTitleRumpunProfesiTable menghapus tabel (gunakan dengan hati-hati!)
func DropJobTitleRumpunProfesiTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.JobTitleRumpunProfesi{})
}
