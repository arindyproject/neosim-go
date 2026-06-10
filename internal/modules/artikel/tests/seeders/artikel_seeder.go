package seeders

import (
	"log"

	"neosim_go/internal/modules/artikel/models"
	"neosim_go/internal/modules/artikel/tests/factories"

	"gorm.io/gorm"
)

// ArtikelSeeder mengelola seeding data artikel
type ArtikelSeeder struct {
	db *gorm.DB
}

func NewArtikelSeeder(db *gorm.DB) *ArtikelSeeder {
	return &ArtikelSeeder{db: db}
}

// Run menjalankan seeder
func (s *ArtikelSeeder) Run() error {
	log.Println("🌱 Seeding artikels...")

	items := factories.NewArtikelFactory().MakeMany(10)
	for _, item := range items {
		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat artikel: %v", err)
			continue
		}
		log.Printf("   ✅ Artikel '%s' dibuat.", item.Name)
	}

	log.Println("✅ Artikels seeding selesai!")
	return nil
}

// Fresh menghapus semua data lalu seed ulang
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

// seedDefault menyimpan satu item jika belum ada
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
