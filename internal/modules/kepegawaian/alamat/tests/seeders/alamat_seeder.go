package seeders

import (
	"log"

	"neosim_go/internal/modules/kepegawaian/alamat/models"
	"neosim_go/internal/modules/kepegawaian/alamat/tests/factories"

	"gorm.io/gorm"
)

// KepegawaianAlamatSeeder mengelola seeding data KepegawaianAlamat
type KepegawaianAlamatSeeder struct {
	db *gorm.DB
}

func NewKepegawaianAlamatSeeder(db *gorm.DB) *KepegawaianAlamatSeeder {
	return &KepegawaianAlamatSeeder{db: db}
}

func (s *KepegawaianAlamatSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_alamats...")

	items := factories.NewKepegawaianAlamatFactory().MakeMany(10)
	for _, item := range items {
		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat KepegawaianAlamat: %v", err)
			continue
		}
		log.Printf("   ✅ KepegawaianAlamat '%s' dibuat.", item.Name)
	}

	log.Println("✅ kepegawaian_alamats seeding selesai!")
	return nil
}

func (s *KepegawaianAlamatSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data kepegawaian_alamats...")
	if err := s.db.Exec("DELETE FROM kepegawaian_alamats").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE kepegawaian_alamats_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *KepegawaianAlamatSeeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.KepegawaianAlamat{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}
	item := factories.NewKepegawaianAlamatFactory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
