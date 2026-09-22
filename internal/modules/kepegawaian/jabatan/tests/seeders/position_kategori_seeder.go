package seeders

import (
	"log"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	"neosim_go/internal/modules/kepegawaian/jabatan/tests/factories"

	"gorm.io/gorm"
)

// PositionKategoriSeeder mengelola seeding data PositionKategori.
// Seeder tetap punya struct sendiri per entitas (bukan digabung ke seeder
// entitas utama), karena tidak ada interface/contract yang perlu di-embed —
// seeder cuma dipanggil manual dari cmd/seed, tidak lewat DI seperti
// repository/service/handler.
type PositionKategoriSeeder struct {
	db *gorm.DB
}

func NewPositionKategoriSeeder(db *gorm.DB) *PositionKategoriSeeder {
	return &PositionKategoriSeeder{db: db}
}

// GetDefaultData mengembalikan daftar master data preset
func GetDefaultData() []models.PositionKategori {
	creatorID := int64(1)

	return []models.PositionKategori{
		{
			Code:      "direksi",
			Label:     "direksi",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "manajerial",
			Label:     "manajerial",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "koordinator",
			Label:     "koordinator",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
		{
			Code:      "pelaksana",
			Label:     "pelaksana",
			CreatedBy: &creatorID,
			UpdatedBy: &creatorID,
		},
	}
}

func (s *PositionKategoriSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_jabatan_position_kategoris...")

	defaults := GetDefaultData()
	for _, item := range defaults {
		// Menggunakan FirstOrCreate berdasarkan `code` agar idempotent
		var existing models.PositionKategori
		err := s.db.Where("code = ?", item.Code).FirstOrCreate(&existing, item).Error
		if err != nil {
			log.Printf("   ⚠️ Gagal membuat/memeriksa PositionKategori [%s]: %v", item.Code, err)
			continue
		}
		log.Printf("   ✅ PositionKategori '%s' (%s) siap.", item.Label, item.Code)
	}

	log.Println("✅ kepegawaian_jabatan_position_kategoris seeding selesai!")
	return nil
}

func (s *PositionKategoriSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data kepegawaian_jabatan_position_kategoris...")
	if err := s.db.Exec("DELETE FROM kepegawaian_jabatan_position_kategoris").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE kepegawaian_jabatan_position_kategoris_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *PositionKategoriSeeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.PositionKategori{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}
	item := factories.NewPositionKategoriFactory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
