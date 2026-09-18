package seeders

import (
	"log"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	"neosim_go/internal/modules/kepegawaian/jabatan/tests/factories"

	"gorm.io/gorm"
)

// PositionSeeder mengelola seeding data Position.
// Seeder tetap punya struct sendiri per entitas (bukan digabung ke seeder
// entitas utama), karena tidak ada interface/contract yang perlu di-embed —
// seeder cuma dipanggil manual dari cmd/seed, tidak lewat DI seperti
// repository/service/handler.
type PositionSeeder struct {
	db *gorm.DB
}

func NewPositionSeeder(db *gorm.DB) *PositionSeeder {
	return &PositionSeeder{db: db}
}

func (s *PositionSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_jabatan_positions...")

	items := factories.NewPositionFactory().MakeMany(10)
	for _, item := range items {
		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat Position: %v", err)
			continue
		}
		log.Printf("   ✅ Position '%s' dibuat.", item.Name)
	}

	log.Println("✅ kepegawaian_jabatan_positions seeding selesai!")
	return nil
}

func (s *PositionSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data kepegawaian_jabatan_positions...")
	if err := s.db.Exec("DELETE FROM kepegawaian_jabatan_positions").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE kepegawaian_jabatan_positions_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *PositionSeeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.Position{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}
	item := factories.NewPositionFactory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
