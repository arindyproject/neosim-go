package seeders

import (
	"log"

	"neosim_go/internal/modules/kepegawaian/pegawai/models"
	"neosim_go/internal/modules/kepegawaian/pegawai/tests/factories"

	"gorm.io/gorm"
)

// KepegawaianPegawaiSeeder mengelola seeding data KepegawaianPegawai
type KepegawaianPegawaiSeeder struct {
	db *gorm.DB
}

func NewKepegawaianPegawaiSeeder(db *gorm.DB) *KepegawaianPegawaiSeeder {
	return &KepegawaianPegawaiSeeder{db: db}
}

func (s *KepegawaianPegawaiSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_pegawais...")

	items := factories.NewKepegawaianPegawaiFactory().MakeMany(10)
	for _, item := range items {
		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat KepegawaianPegawai: %v", err)
			continue
		}
		log.Printf("   ✅ KepegawaianPegawai '%s' dibuat.", item.Name)
	}

	log.Println("✅ kepegawaian_pegawais seeding selesai!")
	return nil
}

func (s *KepegawaianPegawaiSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data kepegawaian_pegawais...")
	if err := s.db.Exec("DELETE FROM kepegawaian_pegawais").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE kepegawaian_pegawais_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *KepegawaianPegawaiSeeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.KepegawaianPegawai{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}
	item := factories.NewKepegawaianPegawaiFactory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
