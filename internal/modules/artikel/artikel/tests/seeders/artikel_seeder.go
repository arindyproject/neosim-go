package seeders

import (
	"log"

	"neosim_go/internal/modules/artikel/artikel/models"
	"neosim_go/internal/modules/artikel/artikel/tests/factories"

	"gorm.io/gorm"
)

// ArtikelSeeder mengelola seeding data Artikel
type ArtikelSeeder struct {
	db *gorm.DB
}

func NewArtikelSeeder(db *gorm.DB) *ArtikelSeeder {
	return &ArtikelSeeder{db: db}
}

func (s *ArtikelSeeder) Run() error {
	log.Println("🌱 Seeding artikels...")

	items := factories.NewArtikelFactory().MakeMany(10)
	for _, item := range items {
		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat Artikel: %v", err)
			continue
		}
		log.Printf("   ✅ Artikel '%s' dibuat.", item.Name)
	}

	log.Println("✅ artikels seeding selesai!")
	return nil
}

func (s *ArtikelSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data artikels...")
	if err := s.db.Exec("DELETE FROM artikels").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE artikels_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *ArtikelSeeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.Artikel{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}
	item := factories.NewArtikelFactory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
