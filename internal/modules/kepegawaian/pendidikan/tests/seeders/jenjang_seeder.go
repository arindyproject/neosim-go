package seeders

import (
	"log"

	"neosim_go/internal/modules/kepegawaian/pendidikan/models"
	"neosim_go/internal/modules/kepegawaian/pendidikan/tests/factories"

	"gorm.io/gorm"
)

// JenjangSeeder mengelola seeding data Jenjang.
// Seeder tetap punya struct sendiri per entitas (bukan digabung ke seeder
// entitas utama), karena tidak ada interface/contract yang perlu di-embed —
// seeder cuma dipanggil manual dari cmd/seed, tidak lewat DI seperti
// repository/service/handler.
type JenjangSeeder struct {
	db *gorm.DB
}

func NewJenjangSeeder(db *gorm.DB) *JenjangSeeder {
	return &JenjangSeeder{db: db}
}

// GetDefaultData mengembalikan daftar master data preset
func GetDefaultData() []models.Jenjang {
	creatorID := int64(1)

	return []models.Jenjang{
		{
			Code:      "SD",
			Label:     "SD",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "SMP",
			Label:     "SMP",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "SMA/SMK",
			Label:     "SMA/SMK",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "D1",
			Label:     "D1",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "D2",
			Label:     "D2",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "D3",
			Label:     "D3",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "D4",
			Label:     "D4",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "S1",
			Label:     "S1",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "Profesi",
			Label:     "Profesi",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "S2",
			Label:     "S2",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "Sp1",
			Label:     "Sp1",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "S3",
			Label:     "S3",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "Sp2",
			Label:     "Sp2",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
	}
}

func (s *JenjangSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_pendidikan_jenjangs...")

	defaults := GetDefaultData()
	for _, item := range defaults {
		// Menggunakan FirstOrCreate berdasarkan `code` agar idempotent
		var existing models.Jenjang
		err := s.db.Where("code = ?", item.Code).FirstOrCreate(&existing, item).Error
		if err != nil {
			log.Printf("   ⚠️ Gagal membuat/memeriksa Tipe [%s]: %v", item.Code, err)
			continue
		}
		log.Printf("   ✅ Jenjang '%s' (%s) siap.", item.Label, item.Code)
	}

	log.Println("✅ kepegawaian_pendidikan_jenjangs seeding selesai!")
	return nil
}

func (s *JenjangSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data kepegawaian_pendidikan_jenjangs...")
	if err := s.db.Exec("DELETE FROM kepegawaian_pendidikan_jenjangs").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE kepegawaian_pendidikan_jenjangs_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *JenjangSeeder) seedDefault(label string) error {
	var count int64
	s.db.Model(&models.Jenjang{}).Where("label = ?", label).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", label)
		return nil
	}
	item := factories.NewJenjangFactory().With("label", label).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", label)
	return nil
}
