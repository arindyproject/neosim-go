package helpers

import (
	"log"

	"neosim_go/config"
	"neosim_go/internal/modules/kepegawaian/lampiran/models"

	"gorm.io/gorm"
)

// SetupTestDB membuat koneksi DB untuk keperluan test
func SetupTestDB() *gorm.DB {
	cfg := config.LoadConfig()
	db, err := cfg.ConnectDB()
	if err != nil {
		log.Fatal("Gagal koneksi DB untuk test:", err)
	}
	return db
}

func MigrateTestDB(db *gorm.DB) {
	if err := db.AutoMigrate(&models.KepegawaianLampiran{}); err != nil {
		log.Fatal("Gagal migrasi test DB:", err)
	}
}

func TruncateTable(db *gorm.DB, tables ...string) {
	for _, table := range tables {
		if err := db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE").Error; err != nil {
			log.Printf("Warning: Gagal truncate table %s: %v", table, err)
		}
	}
}
