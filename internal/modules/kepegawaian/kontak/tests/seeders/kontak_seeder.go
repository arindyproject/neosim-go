package seeders

import (
	"log"

	"neosim_go/internal/modules/kepegawaian/kontak/models"
	"neosim_go/internal/modules/kepegawaian/kontak/tests/factories"

	"gorm.io/gorm"
)

// KepegawaianKontakSeeder mengelola seeding data KepegawaianKontak
type KepegawaianKontakSeeder struct {
	db *gorm.DB
}

func NewKepegawaianKontakSeeder(db *gorm.DB) *KepegawaianKontakSeeder {
	return &KepegawaianKontakSeeder{db: db}
}

func (s *KepegawaianKontakSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_kontaks...")

	items := factories.NewKepegawaianKontakFactory().MakeMany(10)
	for _, item := range items {
		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat KepegawaianKontak: %v", err)
			continue
		}
		log.Printf("   ✅ KepegawaianKontak '%s' dibuat.", item.Nilai)
	}

	log.Println("✅ kepegawaian_kontaks seeding selesai!")
	return nil
}

func (s *KepegawaianKontakSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data kepegawaian_kontaks...")
	if err := s.db.Exec("DELETE FROM kepegawaian_kontaks").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE kepegawaian_kontaks_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *KepegawaianKontakSeeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.KepegawaianKontak{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}
	item := factories.NewKepegawaianKontakFactory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
