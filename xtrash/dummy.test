package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Definisikan tipe IdentifierType agar sesuai dengan model Anda
type IdentifierType string

const (
	TypeNIK IdentifierType = "NIK"
	TypeSTR IdentifierType = "STR"
)

// Struct model yang disalin dari kode Anda
type KepegawaianIdentifier struct {
	ID             int64          `gorm:"primaryKey;autoIncrement;column:id"`
	PegawaiID      int64          `gorm:"column:pegawai_id;not null;index"`
	Tipe           IdentifierType `gorm:"column:tipe;type:varchar(30);not null"`
	Nilai          string         `gorm:"column:nilai;type:varchar(100);not null"`
	Penerbit       *string        `gorm:"column:penerbit;type:varchar(100)"`
	TanggalTerbit  *time.Time     `gorm:"column:tanggal_terbit;type:date"`
	TanggalExpired *time.Time     `gorm:"column:tanggal_expired;type:date"`
	IsPrimary      bool           `gorm:"column:is_primary;default:false"`
	IsAktif        bool           `gorm:"column:is_aktif;default:true"`
	CreatedBy      *int64         `gorm:"column:created_by"`
	UpdatedBy      *int64         `gorm:"column:updated_by"`
	CreatedAt      time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz"`
}

func (KepegawaianIdentifier) TableName() string {
	return "kepegawaian_identifiers"
}

func main() {
	// 1. Koneksi ke Database
	dsn := "host=localhost port=5432 user=users password=rahasia dbname=neosim_go_dev_db sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Ubah ke logger.Info jika ingin melihat query SQL
	})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}
	fmt.Println("✅ Berhasil terhubung ke database!")

	// Konfigurasi Seed
	totalRecords := 100000
	batchSize := 1000 // Insert 1000 baris per query agar cepat dan aman dari limit parameter Postgres
	pegawaiID := int64(1)
	userID := int64(1)

	// Inisialisasi Randomizer
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	fmt.Printf("🚀 Memulai proses generate dan insert %d data dummy...\n", totalRecords)
	startTime := time.Now()

	// 2. Loop untuk membuat data secara batch (menghindari OOM / Out of Memory)
	for i := 0; i < totalRecords; i += batchSize {
		batch := make([]KepegawaianIdentifier, 0, batchSize)

		// Hitung ukuran batch terakhir jika sisa data kurang dari batchSize
		currentBatchSize := batchSize
		if i+batchSize > totalRecords {
			currentBatchSize = totalRecords - i
		}

		// 3. Generate data dummy
		for j := 0; j < currentBatchSize; j++ {
			// Random Tipe: NIK atau STR
			tipe := TypeNIK
			if r.Intn(2) == 0 {
				tipe = TypeSTR
			}

			// Random Nilai (NIK = 16 digit, STR = kombinasi huruf & angka)
			var nilai string
			if tipe == TypeNIK {
				nilai = fmt.Sprintf("%016d", r.Int63n(10000000000000000))
			} else {
				nilai = fmt.Sprintf("STR%d", r.Int63n(10000000000000))
			}

			// Random Penerbit
			var penerbit string
			if tipe == TypeNIK {
				penerbit = "Dukcapil"
			} else {
				penerbitOptions := []string{"Konsil Kedokteran Indonesia", "Dinas Kesehatan", "Ikatan Dokter Indonesia"}
				penerbit = penerbitOptions[r.Intn(len(penerbitOptions))]
			}

			// Random Tanggal Terbit (antara 5 tahun lalu sampai 1 tahun lalu)
			tglTerbit := time.Now().AddDate(-5, 0, 0).Add(time.Duration(r.Intn(365*4)) * 24 * time.Hour)

			// Random Tanggal Expired (Hanya diisi jika STR, NIK dikosongkan)
			var tglExpired *time.Time
			if tipe == TypeSTR {
				exp := tglTerbit.AddDate(5, 0, 0).Add(time.Duration(r.Intn(365)) * 24 * time.Hour)
				tglExpired = &exp
			}

			// Random IsAktif (90% true, 10% false)
			isAktif := true
			if r.Intn(10) == 0 {
				isAktif = false
			}

			// Random Timestamp Audit
			createdAt := time.Now().AddDate(-2, 0, 0).Add(time.Duration(r.Intn(365*2)) * 24 * time.Hour)
			updatedAt := createdAt.Add(time.Duration(r.Intn(30)) * 24 * time.Hour)

			// Susun record
			record := KepegawaianIdentifier{
				PegawaiID:      pegawaiID,
				Tipe:           tipe,
				Nilai:          nilai,
				Penerbit:       &penerbit,
				TanggalTerbit:  &tglTerbit,
				TanggalExpired: tglExpired,
				IsPrimary:      r.Intn(2) == 0,
				IsAktif:        isAktif,
				CreatedBy:      &userID,
				UpdatedBy:      &userID,
				CreatedAt:      createdAt,
				UpdatedAt:      updatedAt,
			}
			batch = append(batch, record)
		}

		// 4. Eksekusi Insert ke Database secara Batch
		result := db.CreateInBatches(batch, len(batch))
		if result.Error != nil {
			log.Printf("❌ Error saat insert batch ke-%d: %v", i/batchSize+1, result.Error)
			break
		}

		// Progress bar sederhana
		fmt.Printf("📊 Progress: %d / %d records inserted...\n", i+currentBatchSize, totalRecords)
	}

	fmt.Printf("🎉 Selesai! Total waktu yang dibutuhkan: %v\n", time.Since(startTime))
}
