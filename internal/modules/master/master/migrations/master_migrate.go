package migrations

import (
	_ "embed"

	"neosim_go/internal/modules/master/master/models"

	"gorm.io/gorm"
)

// MigrateMaster menjalankan GORM auto-migration
func MigrateMaster(db *gorm.DB) error {
	return db.Migrator().CreateTable(
		&models.MasterPekerjaan{},
		&models.MasterPendidikan{},
		&models.MasterStatusPernikahan{},
		&models.MasterAgama{},
	)
}

// DropMasterTable menghapus tabel (gunakan dengan hati-hati!)
func DropMasterTable(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&models.MasterPekerjaan{},
		&models.MasterPendidikan{},
		&models.MasterStatusPernikahan{},
		&models.MasterAgama{},
	)
}
