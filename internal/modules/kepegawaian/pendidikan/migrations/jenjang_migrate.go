package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/kepegawaian/pendidikan/models"

	"gorm.io/gorm"
)

//go:embed 20260821085543_create_kepegawaian_pendidikan_jenjangs_table.sql
var jenjangSQL string

// MigrateJenjang menjalankan GORM auto-migration
func MigrateJenjang(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.Jenjang{})
}

// MigrateJenjangWithSQL menjalankan migrasi via raw SQL
func MigrateJenjangWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(jenjangSQL)
	if err != nil {
		log.Printf("Error creating kepegawaian_pendidikan_jenjangs table: %v", err)
		return err
	}
	log.Println("kepegawaian_pendidikan_jenjangs table migrated successfully")
	return nil
}

// DropJenjangTable menghapus tabel (gunakan dengan hati-hati!)
func DropJenjangTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.Jenjang{})
}
