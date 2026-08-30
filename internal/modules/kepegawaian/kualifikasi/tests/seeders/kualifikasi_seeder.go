package seeders

import (
	"log"

	"neosim_go/internal/modules/kepegawaian/kualifikasi/models"
	"neosim_go/internal/modules/kepegawaian/kualifikasi/tests/factories"

	"gorm.io/gorm"
)

// KepegawaianKualifikasiSeeder mengelola seeding data KepegawaianKualifikasi
type KepegawaianKualifikasiSeeder struct {
	db *gorm.DB
}

func NewKepegawaianKualifikasiSeeder(db *gorm.DB) *KepegawaianKualifikasiSeeder {
	return &KepegawaianKualifikasiSeeder{db: db}
}

func (s *KepegawaianKualifikasiSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_kualifikasis...")

	items := factories.NewKepegawaianKualifikasiFactory().MakeMany(10)
	for _, item := range items {
		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat KepegawaianKualifikasi: %v", err)
			continue
		}
		log.Printf("   ✅ KepegawaianKualifikasi '%s' dibuat.", item.Name)
	}

	log.Println("✅ kepegawaian_kualifikasis seeding selesai!")
	return nil
}

func (s *KepegawaianKualifikasiSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data kepegawaian_kualifikasis...")
	if err := s.db.Exec("DELETE FROM kepegawaian_kualifikasis").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE kepegawaian_kualifikasis_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *KepegawaianKualifikasiSeeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.KepegawaianKualifikasi{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}
	item := factories.NewKepegawaianKualifikasiFactory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
