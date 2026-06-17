package migrations

import (
	_ "embed"

	"neosim_go/internal/modules/master/alamat/models"

	"gorm.io/gorm"
)

// MigrateMasterAlamat menjalankan GORM auto-migration
func MigrateMasterAlamat(db *gorm.DB) error {
	return db.Migrator().CreateTable(
		&models.MasterAlamatNegara{},
		&models.MasterAlamatProvinsi{},
		&models.MasterAlamatKotaKabupaten{},
		&models.MasterAlamatKecamatan{},
		&models.MasterAlamatKelurahanDesa{},
	)
}

// DropMasterAlamatTable menghapus tabel (gunakan dengan hati-hati!)
func DropMasterAlamatTable(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&models.MasterAlamatNegara{},
		&models.MasterAlamatProvinsi{},
		&models.MasterAlamatKotaKabupaten{},
		&models.MasterAlamatKecamatan{},
		&models.MasterAlamatKelurahanDesa{},
	)
}
