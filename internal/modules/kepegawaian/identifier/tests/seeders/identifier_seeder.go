package seeders

import (
	"log"

	"neosim_go/internal/modules/kepegawaian/identifier/models"
	"neosim_go/internal/modules/kepegawaian/identifier/tests/factories"
	pegawaiModel "neosim_go/internal/modules/kepegawaian/pegawai/models"

	"gorm.io/gorm"
)

// KepegawaianIdentifierSeeder mengelola seeding data KepegawaianIdentifier
type KepegawaianIdentifierSeeder struct {
	db *gorm.DB
}

func NewKepegawaianIdentifierSeeder(db *gorm.DB) *KepegawaianIdentifierSeeder {
	return &KepegawaianIdentifierSeeder{db: db}
}

func (s *KepegawaianIdentifierSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_identifiers...")

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

		item := factories.NewKepegawaianIdentifierFactory().
			With("TipeID", validTipe.ID).
			With("PegawaiID", validPegawai.ID).
			Make()

		if err := s.db.Create(item).Error; err != nil {
			log.Printf("   ⚠️  Gagal membuat KepegawaianIdentifier: %v", err)
			continue
		}
		log.Printf("   ✅ KepegawaianIdentifier ID %d (PegawaiID: %d, Tipe: %s, Nilai: %s) dibuat.", item.ID, validPegawai.ID, validTipe.Code, item.Nilai)
	}

	log.Println("✅ kepegawaian_identifiers seeding selesai!")
	return nil
}

func (s *KepegawaianIdentifierSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data kepegawaian_identifiers...")
	if err := s.db.Exec("DELETE FROM kepegawaian_identifiers").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE kepegawaian_identifiers_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *KepegawaianIdentifierSeeder) seedDefault(pegawaiID int64, tipeID int64, nilai string) error {
	var count int64
	s.db.Model(&models.KepegawaianIdentifier{}).
		Where("pegawai_id = ? AND tipe_id = ?", pegawaiID, tipeID).
		Count(&count)

	if count > 0 {
		log.Printf("   ⏭️  Identifier untuk PegawaiID %d dengan TipeID %d sudah ada, skip.", pegawaiID, tipeID)
		return nil
	}

	item := factories.NewKepegawaianIdentifierFactory().
		With("PegawaiID", pegawaiID).
		With("TipeID", tipeID).
		With("Nilai", nilai).
		Make()

	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ Identifier '%s' untuk PegawaiID %d dibuat.", nilai, pegawaiID)
	return nil
}
