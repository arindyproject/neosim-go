package seeders

import (
	"log"

	"neosim_go/internal/modules/kepegawaian/pegawai/models"
	"neosim_go/internal/modules/kepegawaian/pegawai/tests/factories"

	"gorm.io/gorm"
)

// JenisSeeder mengelola seeding data Jenis.
// Seeder tetap punya struct sendiri per entitas (bukan digabung ke seeder
// entitas utama), karena tidak ada interface/contract yang perlu di-embed —
// seeder cuma dipanggil manual dari cmd/seed, tidak lewat DI seperti
// repository/service/handler.
type JenisSeeder struct {
	db *gorm.DB
}

func NewJenisSeeder(db *gorm.DB) *JenisSeeder {
	return &JenisSeeder{db: db}
}

// GetDefaultData mengembalikan daftar master data preset
func GetDefaultJenisData() []models.Jenis {
	creatorID := int64(1)

	return []models.Jenis{
		{
			Code:      "PNS",
			Label:     "PNS",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "PPPK",
			Label:     "PPPK",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "BLUD",
			Label:     "BLUD",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "PTT",
			Label:     "PTT",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "KONTRAK",
			Label:     "KONTRAK",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
	}
}

func (s *JenisSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_pegawai_jeniss...")

	defaults := GetDefaultJenisData()
	for _, item := range defaults {
		// Menggunakan FirstOrCreate berdasarkan `code` agar idempotent
		var existing models.Jenis
		err := s.db.Where("code = ?", item.Code).FirstOrCreate(&existing, item).Error
		if err != nil {
			log.Printf("   ⚠️ Gagal membuat/memeriksa Jenis [%s]: %v", item.Code, err)
			continue
		}
		log.Printf("   ✅ Jenis '%s' (%s) siap.", item.Label, item.Code)
	}

	log.Println("✅ kepegawaian_pegawai_jeniss seeding selesai!")
	return nil
}

func (s *JenisSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data kepegawaian_pegawai_jeniss...")
	if err := s.db.Exec("DELETE FROM kepegawaian_pegawai_jeniss").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE kepegawaian_pegawai_jeniss_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *JenisSeeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.Jenis{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}
	item := factories.NewJenisFactory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
