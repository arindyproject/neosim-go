package seeders

import (
	"log"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	"neosim_go/internal/modules/kepegawaian/jabatan/tests/factories"

	"gorm.io/gorm"
)

// SpecializationSeeder mengelola seeding data Specialization.
// Seeder tetap punya struct sendiri per entitas (bukan digabung ke seeder
// entitas utama), karena tidak ada interface/contract yang perlu di-embed —
// seeder cuma dipanggil manual dari cmd/seed, tidak lewat DI seperti
// repository/service/handler.
//
// PENTING: jalankan setelah seeder job_titles dan specialization_kategoris.
// job_title_id dan kategori_id diisi acak 1 atau 2 oleh factory, jadi
// masing-masing tabel induk harus punya baris dengan ID 1 dan 2.
type SpecializationSeeder struct {
	db *gorm.DB
}

func NewSpecializationSeeder(db *gorm.DB) *SpecializationSeeder {
	return &SpecializationSeeder{db: db}
}

type specializationSeed struct {
	Code           string
	Label          string
	Gelar          string
	LamaPendidikan int16
}

var defaultSpecializations = []specializationSeed{
	{"SP-A", "Spesialis Anak", "Sp.A", 4},
	{"SP-PD", "Spesialis Penyakit Dalam", "Sp.PD", 4},
	{"SP-OG", "Spesialis Obstetri dan Ginekologi", "Sp.OG", 4},
	{"SP-B", "Spesialis Bedah", "Sp.B", 4},
	{"SP-JP", "Spesialis Jantung dan Pembuluh Darah", "Sp.JP", 4},
	{"SP-S", "Spesialis Saraf", "Sp.S", 4},
	{"SP-THT", "Spesialis Telinga Hidung Tenggorok", "Sp.THT-KL", 4},
	{"SP-M", "Spesialis Mata", "Sp.M", 4},
	{"SP-RAD", "Spesialis Radiologi", "Sp.Rad", 4},
	{"SP-PK", "Spesialis Patologi Klinik", "Sp.PK", 4},
}

func (s *SpecializationSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_jabatan_specializations...")

	for _, d := range defaultSpecializations {
		if err := s.seedDefault(d); err != nil {
			log.Printf("   ⚠️  Gagal membuat Specialization '%s': %v", d.Code, err)
		}
	}

	log.Println("✅ kepegawaian_jabatan_specializations seeding selesai!")
	return nil
}

// RunDummy menambah data acak dari factory (untuk testing/load, bukan data master).
func (s *SpecializationSeeder) RunDummy(count int) error {
	log.Printf("🌱 Seeding %d dummy Specialization...", count)

	items := factories.NewSpecializationFactory().MakeMany(count)
	for _, item := range items {
		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat Specialization: %v", err)
			continue
		}
		log.Printf("   ✅ Specialization '%s' dibuat.", item.Label)
	}
	return nil
}

func (s *SpecializationSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data kepegawaian_jabatan_specializations...")
	if err := s.db.Exec("DELETE FROM kepegawaian_jabatan_specializations").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE kepegawaian_jabatan_specializations_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *SpecializationSeeder) seedDefault(d specializationSeed) error {
	var count int64
	if err := s.db.Model(&models.Specialization{}).Where("code = ?", d.Code).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", d.Code)
		return nil
	}

	// job_title_id & kategori_id tidak di-override: factory mengisi acak 1 atau 2.
	item := factories.NewSpecializationFactory().
		With("code", d.Code).
		With("label", d.Label).
		With("gelar", d.Gelar).
		With("lama_pendidikan_tahun", d.LamaPendidikan).
		Make()

	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", d.Code)
	return nil
}
