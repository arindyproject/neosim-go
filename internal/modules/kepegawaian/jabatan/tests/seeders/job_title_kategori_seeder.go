package seeders

import (
	"log"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	"neosim_go/internal/modules/kepegawaian/jabatan/tests/factories"

	"gorm.io/gorm"
)

// JobTitleKategoriSeeder mengelola seeding data JobTitleKategori.
// Seeder tetap punya struct sendiri per entitas (bukan digabung ke seeder
// entitas utama), karena tidak ada interface/contract yang perlu di-embed —
// seeder cuma dipanggil manual dari cmd/seed, tidak lewat DI seperti
// repository/service/handler.
type JobTitleKategoriSeeder struct {
	db *gorm.DB
}

func NewJobTitleKategoriSeeder(db *gorm.DB) *JobTitleKategoriSeeder {
	return &JobTitleKategoriSeeder{db: db}
}

// GetDefaultData mengembalikan daftar master data preset
func GetDefaultDataJobTitleKategori() []models.JobTitleKategori {
	creatorID := int64(1)

	return []models.JobTitleKategori{
		{
			Code:      "medis",
			Label:     "Medis",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "paramedis",
			Label:     "Paramedis",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "penunjang_medis",
			Label:     "Penunjang Medis",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "non_medis",
			Label:     "Non Medis",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
	}
}

func (s *JobTitleKategoriSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_jabatan_job_title_kategoris...")

	defaults := GetDefaultDataJobTitleKategori()
	for _, item := range defaults {
		// Menggunakan FirstOrCreate berdasarkan `code` agar idempotent
		var existing models.JobTitleKategori
		err := s.db.Where("code = ?", item.Code).FirstOrCreate(&existing, item).Error
		if err != nil {
			log.Printf("   ⚠️ Gagal membuat/memeriksa JobTitleKategori [%s]: %v", item.Code, err)
			continue
		}
		log.Printf("   ✅ Tipe '%s' (%s) siap.", item.Label, item.Code)
	}

	log.Println("✅ JobTitleKategori seeding selesai!")
	return nil
}

func (s *JobTitleKategoriSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data kepegawaian_jabatan_job_title_kategoris...")
	if err := s.db.Exec("DELETE FROM kepegawaian_jabatan_job_title_kategoris").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE kepegawaian_jabatan_job_title_kategoris_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *JobTitleKategoriSeeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.JobTitleKategori{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}
	item := factories.NewJobTitleKategoriFactory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
