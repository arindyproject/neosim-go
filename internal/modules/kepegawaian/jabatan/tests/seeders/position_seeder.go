package seeders

import (
	"errors"
	"fmt"
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
//
// PENTING: Position.PositionKategoriID adalah FK NOT NULL, jadi
// PositionKategoriSeeder WAJIB dijalankan sebelum PositionSeeder di
// cmd/seed. Kalau belum ada kategori sama sekali, Run() akan gagal dengan
// pesan jelas alih-alih diam-diam insert dengan position_kategori_id = 0.
type PositionSeeder struct {
	db *gorm.DB
}

func NewPositionSeeder(db *gorm.DB) *PositionSeeder {
	return &PositionSeeder{db: db}
}

func (s *PositionSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_jabatan_positions...")

	defaultKategoriID, err := s.firstKategoriID()
	if err != nil {
		return fmt.Errorf("gagal ambil default position_kategori_id: %w", err)
	}
	if defaultKategoriID == 0 {
		return errors.New("kepegawaian_jabatan_position_kategoris masih kosong — jalankan PositionKategoriSeeder dulu sebelum PositionSeeder")
	}

	items := factories.NewPositionFactory().
		With("position_kategori_id", defaultKategoriID).
		MakeMany(10)

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

// SeedDefaultPositionParams input untuk seed satu Position tetap/default
// (mis. jabatan struktural baku RS: Direktur, Wakil Direktur, dst) —
// beda dari Run() yang isinya data random buat testing.
type SeedDefaultPositionParams struct {
	Name          string
	KategoriID    int64
	ParentID      *int64 // nil = puncak hierarki
	DepartmentID  *int64
	LevelHierarki int16
	IsAktif       bool
}

func (s *PositionSeeder) seedDefault(p SeedDefaultPositionParams) error {
	var count int64
	s.db.Model(&models.Position{}).Where("name = ?", p.Name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", p.Name)
		return nil
	}

	item := factories.NewPositionFactory().
		With("name", p.Name).
		With("position_kategori_id", p.KategoriID).
		With("parent_id", p.ParentID).
		With("department_id", p.DepartmentID).
		With("level_hierarki", p.LevelHierarki).
		With("is_aktif", p.IsAktif).
		Make()

	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", p.Name)
	return nil
}

// firstKategoriID mengambil ID PositionKategori pertama (order by id) untuk
// dipakai sebagai default saat seeding data random di Run(). Return 0 kalau
// tabel kategori masih kosong (bukan error — caller yang memutuskan).
func (s *PositionSeeder) firstKategoriID() (int64, error) {
	var m models.PositionKategori
	err := s.db.Order("id ASC").First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return m.ID, err
}
