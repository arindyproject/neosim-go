package seeders

import (
	"errors"
	"fmt"
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
//
// PENTING: JobTitle.KategoriID adalah FK NOT NULL, jadi JobTitleKategoriSeeder
// WAJIB dijalankan sebelum JobTitleSeeder di cmd/seed. Kalau belum ada
// kategori sama sekali, Run() akan gagal dengan pesan jelas alih-alih
// diam-diam insert dengan kategori_id = 0.
type JobTitleSeeder struct {
	db *gorm.DB
}

func NewJobTitleSeeder(db *gorm.DB) *JobTitleSeeder {
	return &JobTitleSeeder{db: db}
}

func (s *JobTitleSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_jabatan_job_titles...")

	defaultKategoriID, err := s.firstKategoriID()
	if err != nil {
		return fmt.Errorf("gagal ambil default kategori_id: %w", err)
	}
	if defaultKategoriID == 0 {
		return errors.New("kepegawaian_jabatan_job_title_kategoris masih kosong — jalankan JobTitleKategoriSeeder dulu sebelum JobTitleSeeder")
	}

	items := factories.NewJobTitleFactory().
		With("kategori_id", defaultKategoriID).
		MakeMany(10)

	for _, item := range items {
		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat JobTitle: %v", err)
			continue
		}
		log.Printf("   ✅ JobTitle '%s' dibuat.", item.Label)
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

// SeedDefaultJobTitleParams input untuk seed satu JobTitle tetap/default
// (mis. profesi baku RS: Dokter Umum, Perawat, Apoteker) — beda dari Run()
// yang isinya data random buat testing.
type SeedDefaultJobTitleParams struct {
	Code            string
	Label           string
	KategoriID      int64
	RumpunProfesiID *int64 // nil = tidak punya rumpun profesi spesifik
	Point           *float64
	MemerlukanSTR   bool
	MemerlukanSIP   bool
	JenjangMin      *string
	FHIRCode        *string
	FHIRSystem      *string
	IsAktif         bool
}

func (s *JobTitleSeeder) seedDefault(p SeedDefaultJobTitleParams) error {
	var count int64
	s.db.Model(&models.JobTitle{}).Where("code = ?", p.Code).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", p.Label)
		return nil
	}

	item := factories.NewJobTitleFactory().
		With("code", p.Code).
		With("label", p.Label).
		With("kategori_id", p.KategoriID).
		With("rumpun_profesi_id", p.RumpunProfesiID).
		With("point", p.Point).
		With("memerlukan_str", p.MemerlukanSTR).
		With("memerlukan_sip", p.MemerlukanSIP).
		With("jenjang_min", p.JenjangMin).
		With("fhir_code", p.FHIRCode).
		With("fhir_system", p.FHIRSystem).
		With("is_aktif", p.IsAktif).
		Make()

	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", p.Label)
	return nil
}

// firstKategoriID mengambil ID JobTitleKategori pertama (order by id) untuk
// dipakai sebagai default saat seeding data random di Run(). Return 0 kalau
// tabel kategori masih kosong (bukan error — caller yang memutuskan).
func (s *JobTitleSeeder) firstKategoriID() (int64, error) {
	var m models.JobTitleKategori
	err := s.db.Order("id ASC").First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return m.ID, err
}
