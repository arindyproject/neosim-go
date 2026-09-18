package seeders

import (
	"log"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	"neosim_go/internal/modules/kepegawaian/jabatan/tests/factories"

	"gorm.io/gorm"
)

// SpecializationSeeder mengelola seeding data Specialization.
// Seeder tetap punya struct sendiri per entitas (bukan digabung ke seeder
// entitas utama), karena tidak ada interface/contract yang perlu di-embed —
// seeder cuma dipanggil manual dari cmd/seed, tidak lewat DI seperti
// repository/service/handler.
type SpecializationSeeder struct {
	db *gorm.DB
}

func NewSpecializationSeeder(db *gorm.DB) *SpecializationSeeder {
	return &SpecializationSeeder{db: db}
}

func (s *SpecializationSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_jabatan_specializations...")

	items := factories.NewSpecializationFactory().MakeMany(10)
	for _, item := range items {
		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat Specialization: %v", err)
			continue
		}
		log.Printf("   ✅ Specialization '%s' dibuat.", item.Name)
	}

	log.Println("✅ kepegawaian_jabatan_specializations seeding selesai!")
	return nil
}

func (s *SpecializationSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data kepegawaian_jabatan_specializations...")
	if err := s.db.Exec("DELETE FROM kepegawaian_jabatan_specializations").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE kepegawaian_jabatan_specializations_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *SpecializationSeeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.Specialization{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}
	item := factories.NewSpecializationFactory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
