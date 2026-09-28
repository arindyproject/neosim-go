package seeders

import (
	"log"
	"time"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	"neosim_go/internal/modules/kepegawaian/jabatan/tests/factories"

	"gorm.io/gorm"
)

// KepegawaianJabatanSeeder mengelola seeding data KepegawaianJabatan.
//
// PENTING: jalankan setelah seeder pegawai, department, positions, job_titles,
// dan specializations. pegawai_id 1-5, department_id/position_id/job_title_id/
// specialization_id 1-2 harus sudah ada, kalau tidak insert gagal karena FK.
type KepegawaianJabatanSeeder struct {
	db *gorm.DB
}

func NewKepegawaianJabatanSeeder(db *gorm.DB) *KepegawaianJabatanSeeder {
	return &KepegawaianJabatanSeeder{db: db}
}

const seedPegawaiCount = 5

// Run membuat riwayat penugasan untuk pegawai 1..seedPegawaiCount:
// satu jabatan lama yang sudah ditutup + satu jabatan primer yang aktif.
func (s *KepegawaianJabatanSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_jabatans...")

	for pegawaiID := int64(1); pegawaiID <= seedPegawaiCount; pegawaiID++ {
		if err := s.seedDefault(pegawaiID); err != nil {
			log.Printf("   ⚠️  Gagal membuat jabatan pegawai %d: %v", pegawaiID, err)
		}
	}

	log.Println("✅ kepegawaian_jabatans seeding selesai!")
	return nil
}

// RunDummy menambah data acak dari factory (untuk testing/load).
// Semua baris dibuat non-primer sehingga aturan 1 primer aktif per pegawai tetap aman.
func (s *KepegawaianJabatanSeeder) RunDummy(count int) error {
	log.Printf("🌱 Seeding %d dummy KepegawaianJabatan...", count)

	items := factories.NewKepegawaianJabatanFactory().MakeMany(count)
	for _, item := range items {
		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat KepegawaianJabatan: %v", err)
			continue
		}
		log.Printf("   ✅ KepegawaianJabatan pegawai=%d position=%d job_title=%d dibuat.",
			item.PegawaiID, item.PositionID, item.JobTitleID)
	}
	return nil
}

func (s *KepegawaianJabatanSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data kepegawaian_jabatans...")
	if err := s.db.Exec("DELETE FROM kepegawaian_jabatans").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE kepegawaian_jabatans_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

// seedDefault idempoten per pegawai: dilewati kalau pegawai sudah punya jabatan.
func (s *KepegawaianJabatanSeeder) seedDefault(pegawaiID int64) error {
	var count int64
	if err := s.db.Model(&models.KepegawaianJabatan{}).
		Where("pegawai_id = ?", pegawaiID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		log.Printf("   ⏭️  pegawai %d sudah punya jabatan, skip.", pegawaiID)
		return nil
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// 1) Jabatan lama: sudah ditutup (tanggal_selesai terisi, is_aktif = false)
	lamaMulai := today.AddDate(-3, 0, 0)
	lamaSelesai := today.AddDate(-1, 0, 0)
	lama := factories.NewKepegawaianJabatanFactory().
		With("pegawai_id", pegawaiID).
		With("is_primary", false).
		With("tanggal_mulai", lamaMulai).
		With("tanggal_selesai", lamaSelesai).
		With("is_aktif", false).
		Make()
	if err := s.db.Create(lama).Error; err != nil {
		return err
	}

	// 2) Jabatan sekarang: primer dan aktif (tanggal_selesai NULL)
	sekarang := factories.NewKepegawaianJabatanFactory().
		With("pegawai_id", pegawaiID).
		With("is_primary", true).
		With("tanggal_mulai", lamaSelesai).
		With("tanggal_selesai", nil).
		With("is_aktif", true).
		Make()
	if err := s.db.Create(sekarang).Error; err != nil {
		return err
	}

	log.Printf("   ✅ pegawai %d: 1 jabatan lama + 1 jabatan primer aktif dibuat.", pegawaiID)
	return nil
}
