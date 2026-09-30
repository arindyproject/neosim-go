package seeders

import (
	"log"

	"neosim_go/internal/modules/kepegawaian/pegawai/models"
	"neosim_go/internal/modules/kepegawaian/pegawai/tests/factories"

	"gorm.io/gorm"
)

// StatusSeeder mengelola seeding data Status.
// Seeder tetap punya struct sendiri per entitas (bukan digabung ke seeder
// entitas utama), karena tidak ada interface/contract yang perlu di-embed —
// seeder cuma dipanggil manual dari cmd/seed, tidak lewat DI seperti
// repository/service/handler.
type StatusSeeder struct {
	db *gorm.DB
}

func NewStatusSeeder(db *gorm.DB) *StatusSeeder {
	return &StatusSeeder{db: db}
}

// GetDefaultData mengembalikan daftar master data preset
func GetDefaultStatusData() []models.Status {
	creatorID := int64(1)

	return []models.Status{
		{
			Code:      "AKTIF",
			Label:     "AKTIF",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "CUTI",
			Label:     "CUTI",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "TUGAS_BELAJAR",
			Label:     "TUGAS_BELAJAR",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "PENSIUN",
			Label:     "PENSIUN",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
	}
}

func (s *StatusSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_pegawai_statuss...")

	defaults := GetDefaultStatusData()
	for _, item := range defaults {
		// Menggunakan FirstOrCreate berdasarkan `code` agar idempotent
		var existing models.Status
		err := s.db.Where("code = ?", item.Code).FirstOrCreate(&existing, item).Error
		if err != nil {
			log.Printf("   ⚠️ Gagal membuat/memeriksa Status [%s]: %v", item.Code, err)
			continue
		}
		log.Printf("   ✅ Status '%s' (%s) siap.", item.Label, item.Code)
	}

	log.Println("✅ kepegawaian_pegawai_statuss seeding selesai!")
	return nil
}

func (s *StatusSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data kepegawaian_pegawai_statuss...")
	if err := s.db.Exec("DELETE FROM kepegawaian_pegawai_statuss").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE kepegawaian_pegawai_statuss_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *StatusSeeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.Status{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}
	item := factories.NewStatusFactory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
