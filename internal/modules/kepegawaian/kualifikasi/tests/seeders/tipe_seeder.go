package seeders

import (
	"log"

	"neosim_go/internal/modules/kepegawaian/kualifikasi/models"
	"neosim_go/internal/modules/kepegawaian/kualifikasi/tests/factories"

	"gorm.io/gorm"
)

// TipeSeeder mengelola seeding data Tipe.
// Seeder tetap punya struct sendiri per entitas (bukan digabung ke seeder
// entitas utama), karena tidak ada interface/contract yang perlu di-embed —
// seeder cuma dipanggil manual dari cmd/seed, tidak lewat DI seperti
// repository/service/handler.
type TipeSeeder struct {
	db *gorm.DB
}

func NewTipeSeeder(db *gorm.DB) *TipeSeeder {
	return &TipeSeeder{db: db}
}

// GetDefaultData mengembalikan daftar master data preset
func GetDefaultData() []models.Tipe {
	creatorID := int64(1)

	return []models.Tipe{
		{
			Code:      "sertifikasi",
			Label:     "Sertifikasi",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "pelatihan",
			Label:     "Pelatihan",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "kompetensi",
			Label:     "Kompetensi",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
	}
}

func (s *TipeSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_kualifikasi_tipes...")

	defaults := GetDefaultData()
	for _, item := range defaults {
		// Menggunakan FirstOrCreate berdasarkan `code` agar idempotent
		var existing models.Tipe
		err := s.db.Where("code = ?", item.Code).FirstOrCreate(&existing, item).Error
		if err != nil {
			log.Printf("   ⚠️ Gagal membuat/memeriksa Tipe [%s]: %v", item.Code, err)
			continue
		}
		log.Printf("   ✅ Tipe '%s' (%s) siap.", item.Label, item.Code)
	}

	log.Println("✅ kepegawaian_kualifikasi_tipes seeding selesai!")
	return nil
}

func (s *TipeSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data kepegawaian_kualifikasi_tipes...")
	if err := s.db.Exec("DELETE FROM kepegawaian_kualifikasi_tipes").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE kepegawaian_kualifikasi_tipes_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *TipeSeeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.Tipe{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}
	item := factories.NewTipeFactory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
