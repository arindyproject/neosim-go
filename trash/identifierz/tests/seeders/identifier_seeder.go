package seeders

import (
	"fmt"
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

	// 1. Ambil list ID Pegawai yang ada di database agar Foreign Key valid
	var pegawaiIDs []int64
	if err := s.db.Model(&pegawaiModel.KepegawaianPegawai{}).
		Pluck("id", &pegawaiIDs).Error; err != nil || len(pegawaiIDs) == 0 {
		log.Println("   ⚠️  Tidak ada data pegawai ditemukan. Menggunakan fallback PegawaiID 1-5.")
		pegawaiIDs = []int64{1, 2, 3, 4, 5}
	}

	// 2. Tipe-tipe identifier penting yang ingin kita generate secara terstruktur
	defaultTypes := []models.IdentifierType{
		models.IdentifierNIK,
		models.IdentifierSTR,
		models.IdentifierSIP,
		models.IdentifierIHS,
		models.IdentifierBPJSKes,
	}

	successCount := 0

	// 3. Loop untuk setiap PegawaiID yang ada dan buatkan beberapa identifier
	for _, pegawaiID := range pegawaiIDs {
		for _, tipe := range defaultTypes {
			// Cek apakah identifier tipe ini sudah ada untuk pegawai tersebut (mencegah duplikat seeding)
			var count int64
			s.db.Model(&models.KepegawaianIdentifier{}).
				Where("pegawai_id = ? AND tipe = ?", pegawaiID, tipe).
				Count(&count)

			if count > 0 {
				continue
			}

			// Generate dari factory dengan meng-override PegawaiID dan Tipe
			item := factories.NewKepegawaianIdentifierFactory().
				With("PegawaiID", pegawaiID).
				With("Tipe", tipe).
				Make()

			if err := s.db.Create(item).Error; err != nil {
				log.Printf("   ⚠️  Gagal membuat Identifier [%s] untuk Pegawai ID %d: %v", tipe, pegawaiID, err)
				continue
			}

			successCount++
			log.Printf("   ✅ KepegawaianIdentifier [%s] - Nilai: %s (Pegawai ID: %d) dibuat.", item.Tipe, item.Nilai, item.PegawaiID)
		}
	}

	log.Printf("✅ kepegawaian_identifiers seeding selesai! Total %d data ditambahkan.", successCount)
	return nil
}

func (s *KepegawaianIdentifierSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data kepegawaian_identifiers...")
	if err := s.db.Exec("TRUNCATE TABLE kepegawaian_identifiers RESTART IDENTITY CASCADE").Error; err != nil {
		// Fallback jika DB bukan PostgreSQL atau syntax TRUNCATE gagal
		if errDel := s.db.Exec("DELETE FROM kepegawaian_identifiers").Error; errDel != nil {
			return errDel
		}
		_ = s.db.Exec("ALTER SEQUENCE kepegawaian_identifiers_id_seq RESTART WITH 1")
	}

	return s.Run()
}

// SeedSpecificIdentifier membantu seeding spesifik jika dibutuhkan (misal untuk testing/setup awal)
func (s *KepegawaianIdentifierSeeder) SeedSpecificIdentifier(pegawaiID int64, tipe models.IdentifierType, nilai string) error {
	var count int64
	s.db.Model(&models.KepegawaianIdentifier{}).
		Where("pegawai_id = ? AND tipe = ? AND nilai = ?", pegawaiID, tipe, nilai).
		Count(&count)

	if count > 0 {
		log.Printf("   ⏭️  Identifier [%s] '%s' untuk Pegawai ID %d sudah ada, skip.", tipe, nilai, pegawaiID)
		return nil
	}

	item := factories.NewKepegawaianIdentifierFactory().
		With("PegawaiID", pegawaiID).
		With("Tipe", tipe).
		With("Nilai", nilai).
		Make()

	if err := s.db.Create(item).Error; err != nil {
		return fmt.Errorf("gagal membuat identifier: %w", err)
	}

	log.Printf("   ✅ Identifier [%s] '%s' dibuat untuk Pegawai ID %d.", tipe, nilai, pegawaiID)
	return nil
}
