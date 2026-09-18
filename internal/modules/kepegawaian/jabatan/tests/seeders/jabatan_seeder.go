package seeders

import (
	"log"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	"neosim_go/internal/modules/kepegawaian/jabatan/tests/factories"

	"gorm.io/gorm"
)

// KepegawaianJabatanSeeder mengelola seeding data KepegawaianJabatan
type KepegawaianJabatanSeeder struct {
	db *gorm.DB
}

func NewKepegawaianJabatanSeeder(db *gorm.DB) *KepegawaianJabatanSeeder {
	return &KepegawaianJabatanSeeder{db: db}
}

func (s *KepegawaianJabatanSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_jabatans...")

	items := factories.NewKepegawaianJabatanFactory().MakeMany(10)
	for _, item := range items {
		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat KepegawaianJabatan: %v", err)
			continue
		}
		log.Printf("   ✅ KepegawaianJabatan '%s' dibuat.", item.Name)
	}

	log.Println("✅ kepegawaian_jabatans seeding selesai!")
	return nil
}

func (s *KepegawaianJabatanSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data kepegawaian_jabatans...")
	if err := s.db.Exec("DELETE FROM kepegawaian_jabatans").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE kepegawaian_jabatans_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *KepegawaianJabatanSeeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.KepegawaianJabatan{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}
	item := factories.NewKepegawaianJabatanFactory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
