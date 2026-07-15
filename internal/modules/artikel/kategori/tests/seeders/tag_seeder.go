package seeders

import (
	"log"

	"neosim_go/internal/modules/artikel/kategori/models"
	"neosim_go/internal/modules/artikel/kategori/tests/factories"

	"gorm.io/gorm"
)

// TagSeeder mengelola seeding data Tag.
// Seeder tetap punya struct sendiri per entitas (bukan digabung ke seeder
// entitas utama), karena tidak ada interface/contract yang perlu di-embed —
// seeder cuma dipanggil manual dari cmd/seed, tidak lewat DI seperti
// repository/service/handler.
type TagSeeder struct {
	db *gorm.DB
}

func NewTagSeeder(db *gorm.DB) *TagSeeder {
	return &TagSeeder{db: db}
}

func (s *TagSeeder) Run() error {
	log.Println("🌱 Seeding artikel_kategori_tags...")

	items := factories.NewTagFactory().MakeMany(10)
	for _, item := range items {
		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat Tag: %v", err)
			continue
		}
		log.Printf("   ✅ Tag '%s' dibuat.", item.Name)
	}

	log.Println("✅ artikel_kategori_tags seeding selesai!")
	return nil
}

func (s *TagSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data artikel_kategori_tags...")
	if err := s.db.Exec("DELETE FROM artikel_kategori_tags").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE artikel_kategori_tags_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *TagSeeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.Tag{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}
	item := factories.NewTagFactory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
