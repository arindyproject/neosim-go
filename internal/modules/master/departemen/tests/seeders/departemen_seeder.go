package seeders

import (
	"log"

	"neosim_go/internal/modules/master/departemen/models"
	"neosim_go/internal/modules/master/departemen/tests/factories"

	"gorm.io/gorm"
)

// MasterDepartemenSeeder mengelola seeding data MasterDepartemen
type MasterDepartemenSeeder struct {
	db *gorm.DB
}

func NewMasterDepartemenSeeder(db *gorm.DB) *MasterDepartemenSeeder {
	return &MasterDepartemenSeeder{db: db}
}

func (s *MasterDepartemenSeeder) Run() error {
	log.Println("🌱 Seeding master_departemens...")

	items := factories.NewMasterDepartemenFactory().MakeMany(10)
	for _, item := range items {
		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat MasterDepartemen: %v", err)
			continue
		}
		log.Printf("   ✅ MasterDepartemen '%s' dibuat.", item.Name)
	}

	log.Println("✅ master_departemens seeding selesai!")
	return nil
}

func (s *MasterDepartemenSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data master_departemens...")
	if err := s.db.Exec("DELETE FROM master_departemens").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE master_departemens_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *MasterDepartemenSeeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.MasterDepartemen{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}
	item := factories.NewMasterDepartemenFactory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
