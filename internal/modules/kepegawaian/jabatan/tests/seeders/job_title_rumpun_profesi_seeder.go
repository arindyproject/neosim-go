package seeders

import (
	"log"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	"neosim_go/internal/modules/kepegawaian/jabatan/tests/factories"

	"gorm.io/gorm"
)

// JobTitleRumpunProfesiSeeder mengelola seeding data JobTitleRumpunProfesi.
// Seeder tetap punya struct sendiri per entitas (bukan digabung ke seeder
// entitas utama), karena tidak ada interface/contract yang perlu di-embed —
// seeder cuma dipanggil manual dari cmd/seed, tidak lewat DI seperti
// repository/service/handler.
type JobTitleRumpunProfesiSeeder struct {
	db *gorm.DB
}

func NewJobTitleRumpunProfesiSeeder(db *gorm.DB) *JobTitleRumpunProfesiSeeder {
	return &JobTitleRumpunProfesiSeeder{db: db}
}

// GetDefaultData mengembalikan daftar master data preset
func GetDefaultDataJobTitleRumpunProfesi() []models.JobTitleRumpunProfesi {
	creatorID := int64(1)

	return []models.JobTitleRumpunProfesi{
		{
			Code:      "dokter",
			Label:     "Dokter",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "perawat",
			Label:     "Perawat",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "bidan",
			Label:     "Bidan",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "apoteker",
			Label:     "Apoteker",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "tenaga_kesehatan_lain",
			Label:     "Tenaga Kesehatan Lain",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "non_kesehatan",
			Label:     "Non Kesehatan",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
	}
}

func (s *JobTitleRumpunProfesiSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_jabatan_job_title_rumpun_profesis...")

	defaults := GetDefaultDataJobTitleRumpunProfesi()
	for _, item := range defaults {
		// Menggunakan FirstOrCreate berdasarkan `code` agar idempotent
		var existing models.JobTitleRumpunProfesi
		err := s.db.Where("code = ?", item.Code).FirstOrCreate(&existing, item).Error
		if err != nil {
			log.Printf("   ⚠️ Gagal membuat/memeriksa JobTitleRumpunProfesi [%s]: %v", item.Code, err)
			continue
		}
		log.Printf("   ✅ Tipe '%s' (%s) siap.", item.Label, item.Code)
	}

	log.Println("✅ JobTitleRumpunProfesi seeding selesai!")
	return nil
}

func (s *JobTitleRumpunProfesiSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data kepegawaian_jabatan_job_title_rumpun_profesis...")
	if err := s.db.Exec("DELETE FROM kepegawaian_jabatan_job_title_rumpun_profesis").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE kepegawaian_jabatan_job_title_rumpun_profesis_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *JobTitleRumpunProfesiSeeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.JobTitleRumpunProfesi{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}
	item := factories.NewJobTitleRumpunProfesiFactory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
