package seeders

import (
	"log"

	"neosim_go/internal/modules/artikel/kategori/models"
	"neosim_go/internal/modules/artikel/kategori/tests/factories"

	"gorm.io/gorm"
)

// ArtikelKategoriSeeder mengelola seeding data ArtikelKategori
type ArtikelKategoriSeeder struct {
	db *gorm.DB
}

func NewArtikelKategoriSeeder(db *gorm.DB) *ArtikelKategoriSeeder {
	return &ArtikelKategoriSeeder{db: db}
}

func (s *ArtikelKategoriSeeder) Run() error {
	log.Println("🌱 Seeding artikel_kategoris...")

	items := factories.NewArtikelKategoriFactory().MakeMany(10)
	for _, item := range items {
		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat ArtikelKategori: %v", err)
			continue
		}
		log.Printf("   ✅ ArtikelKategori '%s' dibuat.", item.Name)
	}

	log.Println("✅ artikel_kategoris seeding selesai!")
	return nil
}

func (s *ArtikelKategoriSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data artikel_kategoris...")
	if err := s.db.Exec("DELETE FROM artikel_kategoris").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE artikel_kategoris_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *ArtikelKategoriSeeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.ArtikelKategori{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}
	item := factories.NewArtikelKategoriFactory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
