package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/kepegawaian/pegawai/models"

	"gorm.io/gorm"
)

//go:embed 20260929135806_create_kepegawaian_pegawai_statuss_table.sql
var statusSQL string

// MigrateStatus menjalankan GORM auto-migration
func MigrateStatus(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.Status{})
}

// MigrateStatusWithSQL menjalankan migrasi via raw SQL
func MigrateStatusWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(statusSQL)
	if err != nil {
		log.Printf("Error creating kepegawaian_pegawai_statuss table: %v", err)
		return err
	}
	log.Println("kepegawaian_pegawai_statuss table migrated successfully")
	return nil
}

// DropStatusTable menghapus tabel (gunakan dengan hati-hati!)
func DropStatusTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.Status{})
}
