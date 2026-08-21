package seeders

import (
	"log"

	"neosim_go/internal/modules/kepegawaian/pendidikan/models"
	"neosim_go/internal/modules/kepegawaian/pendidikan/tests/factories"

	"gorm.io/gorm"
)

// KepegawaianPendidikanSeeder mengelola seeding data KepegawaianPendidikan
type KepegawaianPendidikanSeeder struct {
	db *gorm.DB
}

func NewKepegawaianPendidikanSeeder(db *gorm.DB) *KepegawaianPendidikanSeeder {
	return &KepegawaianPendidikanSeeder{db: db}
}

func (s *KepegawaianPendidikanSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_pendidikans...")

	items := factories.NewKepegawaianPendidikanFactory().MakeMany(10)
	for _, item := range items {
		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat KepegawaianPendidikan: %v", err)
			continue
		}
		log.Printf("   ✅ KepegawaianPendidikan '%s' dibuat.", item.NamaInstitusi)
	}

	log.Println("✅ kepegawaian_pendidikans seeding selesai!")
	return nil
}

func (s *KepegawaianPendidikanSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data kepegawaian_pendidikans...")
	if err := s.db.Exec("DELETE FROM kepegawaian_pendidikans").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE kepegawaian_pendidikans_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *KepegawaianPendidikanSeeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.KepegawaianPendidikan{}).Where("nama_institusi = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}
	item := factories.NewKepegawaianPendidikanFactory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
