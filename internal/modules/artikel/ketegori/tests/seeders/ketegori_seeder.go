package seeders

import (
	"log"

	"neosim_go/internal/modules/artikel/ketegori/models"
	"neosim_go/internal/modules/artikel/ketegori/tests/factories"

	"gorm.io/gorm"
)

// ArtikelKetegoriSeeder mengelola seeding data ArtikelKetegori
type ArtikelKetegoriSeeder struct {
	db *gorm.DB
}

func NewArtikelKetegoriSeeder(db *gorm.DB) *ArtikelKetegoriSeeder {
	return &ArtikelKetegoriSeeder{db: db}
}

// Run menjalankan seeder
func (s *ArtikelKetegoriSeeder) Run() error {
	log.Println("🌱 Seeding artikel_ketegoris...")

	items := factories.NewArtikelKetegoriFactory().MakeMany(10)
	for _, item := range items {
		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat ArtikelKetegori: %v", err)
			continue
		}
		log.Printf("   ✅ ArtikelKetegori '%s' dibuat.", item.Name)
	}

	log.Println("✅ artikel_ketegoris seeding selesai!")
	return nil
}

// Fresh menghapus semua data lalu seed ulang
func (s *ArtikelKetegoriSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data artikel_ketegoris...")

	if err := s.db.Exec("DELETE FROM artikel_ketegoris").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE artikel_ketegoris_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

// seedDefault menyimpan satu item jika belum ada
func (s *ArtikelKetegoriSeeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.ArtikelKetegori{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}

	item := factories.NewArtikelKetegoriFactory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}

	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
