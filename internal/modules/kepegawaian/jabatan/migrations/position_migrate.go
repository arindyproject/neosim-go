package migrations

import (
	"database/sql"
	_ "embed"
	"log"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"

	"gorm.io/gorm"
)

//go:embed 20260918095114_create_kepegawaian_jabatan_positions_table.sql
var positionSQL string

// MigratePosition menjalankan GORM auto-migration
func MigratePosition(db *gorm.DB) error {
	return db.Migrator().CreateTable(&models.Position{})
}

// MigratePositionWithSQL menjalankan migrasi via raw SQL
func MigratePositionWithSQL(sqlDB *sql.DB) error {
	_, err := sqlDB.Exec(positionSQL)
	if err != nil {
		log.Printf("Error creating kepegawaian_jabatan_positions table: %v", err)
		return err
	}
	log.Println("kepegawaian_jabatan_positions table migrated successfully")
	return nil
}

// DropPositionTable menghapus tabel (gunakan dengan hati-hati!)
func DropPositionTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.Position{})
}
