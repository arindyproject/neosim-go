package seeders

import (
	"fmt"
	"log"

	"neosim_go/internal/modules/kepegawaian/pegawai/models"
	"neosim_go/internal/modules/kepegawaian/pegawai/tests/factories"

	"gorm.io/gorm"
)

type masterSeed struct{ Code, Label string }

var defaultJenis = []masterSeed{
	{"PNS", "Pegawai Negeri Sipil"},
	{"PPPK", "Pegawai Pemerintah dengan Perjanjian Kerja"},
	{"TETAP", "Pegawai Tetap"},
	{"KONTRAK", "Pegawai Kontrak"},
	{"HONORER", "Tenaga Honorer"},
}

var defaultStatus = []masterSeed{
	{"AKTIF", "Aktif"},
	{"CUTI", "Cuti"},
	{"NONAKTIF", "Nonaktif"},
	{"PENSIUN", "Pensiun"},
}

// KepegawaianPegawaiSeeder mengelola seeding data KepegawaianPegawai
type KepegawaianPegawaiSeeder struct {
	db *gorm.DB
}

func NewKepegawaianPegawaiSeeder(db *gorm.DB) *KepegawaianPegawaiSeeder {
	return &KepegawaianPegawaiSeeder{db: db}
}

// ensureMaster memastikan Jenis & Status default ada (idempotent), lalu
// mengembalikan ID Jenis dan ID status AKTIF.
func (s *KepegawaianPegawaiSeeder) ensureMaster() (jenisIDs []int64, statusAktifID int64, err error) {
	for _, j := range defaultJenis {
		m := models.Jenis{Code: j.Code, Label: j.Label}
		if err := s.db.Where("code = ?", j.Code).FirstOrCreate(&m).Error; err != nil {
			return nil, 0, fmt.Errorf("seed jenis %s: %w", j.Code, err)
		}
		jenisIDs = append(jenisIDs, m.ID)
	}
	for _, st := range defaultStatus {
		m := models.Status{Code: st.Code, Label: st.Label}
		if err := s.db.Where("code = ?", st.Code).FirstOrCreate(&m).Error; err != nil {
			return nil, 0, fmt.Errorf("seed status %s: %w", st.Code, err)
		}
		if st.Code == "AKTIF" {
			statusAktifID = m.ID
		}
	}
	return jenisIDs, statusAktifID, nil
}

func (s *KepegawaianPegawaiSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_pegawais...")

	jenisIDs, statusAktifID, err := s.ensureMaster()
	if err != nil {
		return err
	}

	const total = 50 // user_id 1..50

	items := factories.NewKepegawaianPegawaiFactory().MakeMany(total)
	for i, item := range items {
		userID := int64(i + 1)

		// skip jika user ini sudah tertaut ke pegawai lain (aman dijalankan ulang)
		var count int64
		s.db.Model(&models.KepegawaianPegawai{}).
			Where("user_id = ?", userID).
			Count(&count)
		if count > 0 {
			log.Printf("   ⏭️  user_id %d sudah tertaut, skip.", userID)
			continue
		}

		item.UserID = &userID
		item.JenisID = jenisIDs[i%len(jenisIDs)]
		item.StatusID = statusAktifID

		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat pegawai (user_id %d): %v", userID, err)
			continue
		}
		log.Printf("   ✅ Pegawai '%s' (%s) dibuat, user_id %d.", item.NamaLengkap, item.NomorPegawai, userID)
	}

	log.Println("✅ kepegawaian_pegawais seeding selesai!")
	return nil
}

func (s *KepegawaianPegawaiSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data kepegawaian_pegawais...")
	if err := s.db.Exec("DELETE FROM kepegawaian_pegawais").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE kepegawaian_pegawais_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

// seedDefault membuat satu pegawai dengan nomor_pegawai tertentu (skip jika sudah ada)
func (s *KepegawaianPegawaiSeeder) seedDefault(nomorPegawai, namaLengkap string) error {
	var count int64
	s.db.Model(&models.KepegawaianPegawai{}).
		Where("nomor_pegawai = ?", nomorPegawai).
		Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", nomorPegawai)
		return nil
	}

	jenisIDs, statusAktifID, err := s.ensureMaster()
	if err != nil {
		return err
	}

	item := factories.NewKepegawaianPegawaiFactory().
		With("nomor_pegawai", nomorPegawai).
		With("nama_lengkap", namaLengkap).
		With("jenis_id", jenisIDs[0]).
		With("status_id", statusAktifID).
		Make()

	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", nomorPegawai)
	return nil
}
