package seeders

import (
	"log"

	"neosim_go/internal/modules/kepegawaian/lampiran/models"
	"neosim_go/internal/modules/kepegawaian/lampiran/tests/factories"

	"gorm.io/gorm"
)

// KepegawaianLampiranSeeder mengelola seeding data KepegawaianLampiran
type KepegawaianLampiranSeeder struct {
	db *gorm.DB
}

func NewKepegawaianLampiranSeeder(db *gorm.DB) *KepegawaianLampiranSeeder {
	return &KepegawaianLampiranSeeder{db: db}
}

func (s *KepegawaianLampiranSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_lampirans...")

	items := factories.NewKepegawaianLampiranFactory().MakeMany(10)
	for _, item := range items {
		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat KepegawaianLampiran: %v", err)
			continue
		}
		log.Printf("   ✅ KepegawaianLampiran '%s' dibuat.", item.Name)
	}

	log.Println("✅ kepegawaian_lampirans seeding selesai!")
	return nil
}

func (s *KepegawaianLampiranSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data kepegawaian_lampirans...")
	if err := s.db.Exec("DELETE FROM kepegawaian_lampirans").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE kepegawaian_lampirans_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *KepegawaianLampiranSeeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.KepegawaianLampiran{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}
	item := factories.NewKepegawaianLampiranFactory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
