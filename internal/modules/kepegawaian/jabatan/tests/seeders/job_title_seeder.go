package seeders

import (
	"log"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	"neosim_go/internal/modules/kepegawaian/jabatan/tests/factories"

	"gorm.io/gorm"
)

// JobTitleSeeder mengelola seeding data JobTitle.
// Seeder tetap punya struct sendiri per entitas (bukan digabung ke seeder
// entitas utama), karena tidak ada interface/contract yang perlu di-embed —
// seeder cuma dipanggil manual dari cmd/seed, tidak lewat DI seperti
// repository/service/handler.
type JobTitleSeeder struct {
	db *gorm.DB
}

func NewJobTitleSeeder(db *gorm.DB) *JobTitleSeeder {
	return &JobTitleSeeder{db: db}
}

func (s *JobTitleSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_jabatan_job_titles...")

	items := factories.NewJobTitleFactory().MakeMany(10)
	for _, item := range items {
		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat JobTitle: %v", err)
			continue
		}
		log.Printf("   ✅ JobTitle '%s' dibuat.", item.Name)
	}

	log.Println("✅ kepegawaian_jabatan_job_titles seeding selesai!")
	return nil
}

func (s *JobTitleSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data kepegawaian_jabatan_job_titles...")
	if err := s.db.Exec("DELETE FROM kepegawaian_jabatan_job_titles").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE kepegawaian_jabatan_job_titles_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *JobTitleSeeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.JobTitle{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}
	item := factories.NewJobTitleFactory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
