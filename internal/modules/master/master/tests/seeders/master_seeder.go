package seeders

import (
	"log"

	"neosim_go/internal/modules/master/master/models"
	"neosim_go/internal/modules/master/master/tests/factories"

	"gorm.io/gorm"
)

// MasterSeeder mengelola seeding data Master
type MasterSeeder struct {
	db *gorm.DB
}

func NewMasterSeeder(db *gorm.DB) *MasterSeeder {
	return &MasterSeeder{db: db}
}

// Run menjalankan seeder
func (s *MasterSeeder) Run() error {
	log.Println("🌱 Seeding masters...")

	items := factories.NewMasterFactory().MakeMany(10)
	for _, item := range items {
		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat Master: %v", err)
			continue
		}
		log.Printf("   ✅ Master '%s' dibuat.", item.Name)
	}

	log.Println("✅ masters seeding selesai!")
	return nil
}

// Fresh menghapus semua data lalu seed ulang
func (s *MasterSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data masters...")

	if err := s.db.Exec("DELETE FROM masters").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE masters_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

// seedDefault menyimpan satu item jika belum ada
func (s *MasterSeeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.Master{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}

	item := factories.NewMasterFactory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}

	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
