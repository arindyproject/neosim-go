package seeders

import (
	"log"

	"neosim_go/internal/modules/kepegawaian/kualifikasi/models"
	"neosim_go/internal/modules/kepegawaian/kualifikasi/tests/factories"

	pegawaiModel "neosim_go/internal/modules/kepegawaian/pegawai/models"

	"gorm.io/gorm"
)

// KepegawaianKualifikasiSeeder mengelola seeding data KepegawaianKualifikasi
type KepegawaianKualifikasiSeeder struct {
	db *gorm.DB
}

func NewKepegawaianKualifikasiSeeder(db *gorm.DB) *KepegawaianKualifikasiSeeder {
	return &KepegawaianKualifikasiSeeder{db: db}
}

func (s *KepegawaianKualifikasiSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_kualifikasis...")

	// 1. Ambil daftar Tipe ID yang ada untuk memastikan relasi valid
	var tipes []models.Tipe
	if err := s.db.Find(&tipes).Error; err != nil || len(tipes) == 0 {
		log.Println("   ⚠️  Master data 'Tipe' tidak ditemukan. Harap pastikan TipeSeeder dijalankan terlebih dahulu.")
		return err
	}

	// 2. Ambil daftar Pegawai yang ada untuk memastikan PegawaiID valid
	var pegawais []pegawaiModel.KepegawaianPegawai
	if err := s.db.Find(&pegawais).Error; err != nil || len(pegawais) == 0 {
		log.Println("   ⚠️  Data 'Pegawai' tidak ditemukan. Harap pastikan PegawaiSeeder dijalankan terlebih dahulu.")
		return err
	}

	// 3. Buat sample data menggunakan factory dengan TipeID dan PegawaiID yang valid dari DB
	for i := 0; i < 10; i++ {
		validTipe := tipes[i%len(tipes)]          // Rotasi pilihan tipe
		validPegawai := pegawais[i%len(pegawais)] // Rotasi pilihan pegawai

		item := factories.NewKepegawaianKualifikasiFactory().
			With("TipeID", validTipe.ID).
			With("PegawaiID", validPegawai.ID).
			Make()

		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat KepegawaianKualifikasi: %v", err)
			continue
		}
		log.Printf("   ✅ KepegawaianKualifikasi ID %d (PegawaiID: %d, Tipe: %s, Nama: %s) dibuat.", item.ID, validPegawai.ID, validTipe.Code, item.Nama)
	}

	log.Println("✅ kepegawaian_kualifikasis seeding selesai!")
	return nil
}

func (s *KepegawaianKualifikasiSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data kepegawaian_kualifikasis...")
	if err := s.db.Exec("DELETE FROM kepegawaian_kualifikasis").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE kepegawaian_kualifikasis_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

// seedDefault membuat satu KepegawaianKualifikasi default untuk pegawai tertentu,
// menghindari duplikasi berdasarkan kombinasi pegawai_id + nomor_sertifikat.
func (s *KepegawaianKualifikasiSeeder) seedDefault(pegawaiID int64, nama string, nomorSertifikat string) error {
	var count int64
	s.db.Model(&models.KepegawaianKualifikasi{}).
		Where("pegawai_id = ? AND nomor_sertifikat = ?", pegawaiID, nomorSertifikat).
		Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' (pegawai_id=%d) sudah ada, skip.", nama, pegawaiID)
		return nil
	}
	item := factories.NewKepegawaianKualifikasiFactory().
		With("pegawai_id", pegawaiID).
		With("nama", nama).
		With("nomor_sertifikat", nomorSertifikat).
		Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' (pegawai_id=%d) dibuat.", nama, pegawaiID)
	return nil
}
