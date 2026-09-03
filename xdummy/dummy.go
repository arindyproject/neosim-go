// neosim_go/xdummy/dummy.go
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"neosim_go/config"
	artikelFactories "neosim_go/internal/modules/artikel/artikel/tests/factories"
	kategoriFactories "neosim_go/internal/modules/artikel/kategori/tests/factories"
	kepegawaianIdentifierFactories "neosim_go/internal/modules/kepegawaian/identifier/tests/factories"
	kepegawaianKontakFactories "neosim_go/internal/modules/kepegawaian/kontak/tests/factories"
	kepegawaianKualifikasiFactories "neosim_go/internal/modules/kepegawaian/kualifikasi/tests/factories"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ─── Seeder Registry ───────────────────────────────────────────────
//
// Cara menambah model baru:
//   1. Pastikan model tersebut punya Factory dengan method:
//        MakeMany(count int) []*models.NamaModel
//      (lihat contoh: internal/modules/artikel/artikel/tests/factories)
//   2. Import package factory-nya di atas.
//   3. Tambahkan satu entry baru ke map `registry` di bawah ini,
//      formatnya seragam untuk semua model:
//
//        "pasien": {
//            Run: func(db *gorm.DB, total, batch int) error {
//                factory := pasienFactories.NewPasienFactory()
//                return seedBatch(db, "pasien", total, batch, factory.MakeMany)
//            },
//        },
//
//   Tidak perlu ubah kode lain — logika batching, progress bar, dan
//   flag -model otomatis berlaku untuk model baru tersebut.

type Seeder struct {
	Run func(db *gorm.DB, total, batchSize int) error
}

var registry = map[string]Seeder{
	"artikel": {
		Run: func(db *gorm.DB, total, batchSize int) error {
			factory := artikelFactories.NewArtikelFactory()
			return seedBatch(db, "artikel", total, batchSize, factory.MakeMany)
		},
	},

	"kategori": {
		Run: func(db *gorm.DB, total, batchSize int) error {
			factory := kategoriFactories.NewArtikelKategoriFactory()
			return seedBatch(db, "kategori", total, batchSize, factory.MakeMany)
		},
	},

	"tag": {
		Run: func(db *gorm.DB, total, batchSize int) error {
			factory := kategoriFactories.NewTagFactory()
			return seedBatch(db, "tag", total, batchSize, factory.MakeMany)
		},
	},

	"kepegawaian_identifier": {
		Run: func(db *gorm.DB, total, batchSize int) error {
			factory := kepegawaianIdentifierFactories.NewKepegawaianIdentifierFactory()
			return seedBatch(db, "kepegawaian_identifier", total, batchSize, factory.MakeMany)
		},
	},

	"kepegawaian_kontak": {
		Run: func(db *gorm.DB, total, batchSize int) error {
			factory := kepegawaianKontakFactories.NewKepegawaianKontakFactory()
			return seedBatch(db, "kepegawaian_kontak", total, batchSize, factory.MakeMany)
		},
	},
	"kepegawaian_kualifikasis": {
		Run: func(db *gorm.DB, total, batchSize int) error {
			factory := kepegawaianKualifikasiFactories.NewKepegawaianKualifikasiFactory()
			return seedBatch(db, "kepegawaian_kualifikasis", total, batchSize, factory.MakeMany)
		},
	},
	// Tambahkan seeder model lain di sini, contoh pola di atas.
}

func main() {
	// ─── Flags ──────────────────────────────────────────────────────
	totalRecords := flag.Int("total", 100000, "jumlah total data dummy per model")
	batchSize := flag.Int("batch", 2000, "jumlah record per batch insert")
	modelFlag := flag.String("model", "all", "nama model yang di-seed, pisahkan koma (misal: artikel,pasien), atau 'all'")
	flag.Parse()

	selected, err := resolveSeeders(*modelFlag)
	if err != nil {
		log.Fatalf("❌ %v", err)
	}

	// ─── 1. Load .env tambahan (opsional) ──────────────────────────
	if err := godotenv.Load(".env", "config/.env", "config/.env.dev"); err != nil {
		log.Printf("⚠️  Info: tidak ada file .env tambahan yang ter-load (%v) — lanjut pakai config/.env.dev bawaan LoadConfig()", err)
	}

	// ─── 2. Load Config & Connect DB ───────────────────────────────
	cfg := config.LoadConfig()
	db, err := cfg.ConnectDB()
	if err != nil {
		log.Fatalf("❌ Koneksi database gagal: %v", err)
	}
	defer config.CloseDB(db)

	// ─── 3. Matikan query log GORM selama seeding ──────────────────
	db.Logger = logger.Default.LogMode(logger.Silent)

	fmt.Printf("🚀 Seeding model: %s | total: %d | batch: %d\n", strings.Join(selected, ", "), *totalRecords, *batchSize)
	overallStart := time.Now()

	// ─── 4. Jalankan tiap seeder yang dipilih ──────────────────────
	for _, name := range selected {
		if err := registry[name].Run(db, *totalRecords, *batchSize); err != nil {
			log.Fatalf("❌ Seeder '%s' gagal: %v", name, err)
		}
	}

	fmt.Printf("✅ Semua seeder selesai dalam %v\n", time.Since(overallStart))
}

// resolveSeeders menerjemahkan flag -model ("all" atau "a,b,c") menjadi
// daftar nama seeder yang valid, urut sesuai urutan didaftarkan di map
// tidak dijamin (Go map unordered) — kalau butuh urutan pasti, ganti
// registry jadi slice. Untuk kebutuhan seeding paralel-independen ini
// urutan biasanya tidak masalah.
func resolveSeeders(modelFlag string) ([]string, error) {
	if modelFlag == "all" {
		names := make([]string, 0, len(registry))
		for name := range registry {
			names = append(names, name)
		}
		if len(names) == 0 {
			return nil, fmt.Errorf("registry seeder kosong, belum ada model yang didaftarkan")
		}
		return names, nil
	}

	parts := strings.Split(modelFlag, ",")
	names := make([]string, 0, len(parts))
	for _, p := range parts {
		name := strings.TrimSpace(p)
		if name == "" {
			continue
		}
		if _, ok := registry[name]; !ok {
			available := make([]string, 0, len(registry))
			for n := range registry {
				available = append(available, n)
			}
			return nil, fmt.Errorf("model '%s' tidak dikenal, tersedia: %s", name, strings.Join(available, ", "))
		}
		names = append(names, name)
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("tidak ada model valid pada flag -model=%q", modelFlag)
	}
	return names, nil
}

// seedBatch adalah logika generik: generate + insert per-batch + progress bar.
// Berlaku untuk model apa pun selama factory-nya punya method
// MakeMany(count int) []T.
func seedBatch[T any](db *gorm.DB, label string, total, batchSize int, makeMany func(int) []T) error {
	inserted := 0
	start := time.Now()

	for inserted < total {
		currentBatchSize := batchSize
		if remaining := total - inserted; remaining < batchSize {
			currentBatchSize = remaining
		}

		items := makeMany(currentBatchSize)

		if err := db.CreateInBatches(&items, currentBatchSize).Error; err != nil {
			return fmt.Errorf("insert batch pada offset ke-%d: %w", inserted, err)
		}

		inserted += currentBatchSize
		elapsed := time.Since(start)
		rate := float64(inserted) / elapsed.Seconds()
		percent := float64(inserted) / float64(total) * 100

		eta := "-"
		if rate > 0 {
			remaining := total - inserted
			eta = time.Duration(float64(remaining) / rate * float64(time.Second)).Round(time.Second).String()
		}

		fmt.Printf("   ⏳ [%s] %d / %d (%.1f%%) | %.0f rec/s | ETA: %s\n",
			label, inserted, total, percent, rate, eta)
	}

	fmt.Printf("🎉 [%s] Selesai! %d data dibuat dalam %v\n", label, total, time.Since(start))
	return nil
}
