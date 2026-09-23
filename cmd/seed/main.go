package main

import (
	"flag"
	"log"

	"neosim_go/config"

	// =====================================================================
	// import seeder di sini
	// =====================================================================
	// Menggunakan alias untuk membedakan kedua paket seeders

	rbacSeed "neosim_go/internal/modules/rbac/tests/seeders"
	userSeed "neosim_go/internal/modules/users/tests/seeders"

	// Master---------------------------------------------------------------
	masterAlamat "neosim_go/internal/modules/master/alamat/tests/seeders"
	masterDepartemen "neosim_go/internal/modules/master/departemen/tests/seeders"
	masterMaster "neosim_go/internal/modules/master/master/tests/seeders"

	// Kepegawaian----------------------------------------------------------
	masterKepegawaianAlamat "neosim_go/internal/modules/kepegawaian/alamat/tests/seeders"
	masterKepegawaianIdentifier "neosim_go/internal/modules/kepegawaian/identifier/tests/seeders"
	masterKepegawaianJabatan "neosim_go/internal/modules/kepegawaian/jabatan/tests/seeders"
	masterKepegawaianKontak "neosim_go/internal/modules/kepegawaian/kontak/tests/seeders"
	masterKepegawaianKualifikasi "neosim_go/internal/modules/kepegawaian/kualifikasi/tests/seeders"
	masterKepegawaianPegawai "neosim_go/internal/modules/kepegawaian/pegawai/tests/seeders"
	masterKepegawaianPendidikan "neosim_go/internal/modules/kepegawaian/pendidikan/tests/seeders"

	// Artikel---------------------------------------------------------------
	masterArtikel "neosim_go/internal/modules/artikel/artikel/tests/seeders"
	masterArtikelKategori "neosim_go/internal/modules/artikel/kategori/tests/seeders"
	// =====================================================================
)

// Seeder kontrak minimal yang dipakai semua *Seeder di project ini — semua
// sudah punya Run() dan Fresh(), jadi tidak perlu interface baru per modul.
type Seeder interface {
	Run() error
	Fresh() error
}

// runSeeder menjalankan satu seeder sesuai flag --fresh, dan log.Fatal kalau
// gagal — persis perilaku tiap blok if/else berulang di kode aslinya, cuma
// ditulis sekali. 'label' dipakai buat pesan log, meniru teks error asli
// per bagian (mis. "Gagal seed Pegawai:", "Gagal fresh seed Pegawai:").
func runSeeder(label string, s Seeder, fresh bool) {
	var err error
	if fresh {
		err = s.Fresh()
	} else {
		err = s.Run()
	}
	if err != nil {
		verb := "seed"
		if fresh {
			verb = "fresh seed"
		}
		log.Fatalf("❌ Gagal %s %s: %v", verb, label, err)
	}
}

func main() {
	cfg := config.LoadConfig()
	env := flag.String("env", cfg.EnvCode, "Environment (DEV atau PROD)")
	fresh := flag.Bool("fresh", false, "Hapus semua data lalu seed ulang")
	flag.Parse()

	if *env != "DEV" && *env != "PROD" {
		log.Fatal("❌ Environment tidak valid. Gunakan DEV atau PROD")
	}

	// Safety guard: fresh seed tidak boleh di PROD
	if *fresh && *env == "PROD" {
		log.Fatal("❌ Fresh seed TIDAK diizinkan di environment PROD!")
	}

	log.Printf("🚀 Menjalankan seeder untuk environment: %s", *env)

	db, err := cfg.ConnectDB()
	if err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}
	defer config.CloseDB(db)

	// =====================================================================
	// Urutan di bawah ini SENGAJA dipertahankan sama persis dengan versi
	// sebelumnya — banyak seeder di sini punya dependensi FK ke seeder di
	// atasnya (mis. Position butuh PositionKategori, JobTitle butuh
	// JobTitleKategori), jadi urutan bukan sekadar kosmetik.
	// =====================================================================

	runSeeder("users", userSeed.NewUserSeeder(db), *fresh)
	runSeeder("RBAC", rbacSeed.NewRBACSeeder(db), *fresh)

	// Master-----------------------------------------------------------------
	runSeeder("master alamat", masterAlamat.NewMasterAlamatSeeder(db), *fresh)
	runSeeder("master", masterMaster.NewMasterSeeder(db), *fresh)
	runSeeder("departemen", masterDepartemen.NewMasterDepartemenSeeder(db), *fresh)

	// Kepegawaian==============================================================
	runSeeder("Pegawai", masterKepegawaianPegawai.NewKepegawaianPegawaiSeeder(db), *fresh)

	// Kepegawaian - Identifier
	runSeeder("Identifier Type", masterKepegawaianIdentifier.NewTipeSeeder(db), *fresh)
	runSeeder("Identifier", masterKepegawaianIdentifier.NewKepegawaianIdentifierSeeder(db), *fresh)

	// Kepegawaian - Kontak
	runSeeder("Kontak Tipe", masterKepegawaianKontak.NewTipeSeeder(db), *fresh)
	runSeeder("Kontak", masterKepegawaianKontak.NewKepegawaianKontakSeeder(db), *fresh)

	// Kepegawaian - Pendidikan
	runSeeder("jenjang pendidikan", masterKepegawaianPendidikan.NewJenjangSeeder(db), *fresh)
	runSeeder("pendidikan", masterKepegawaianPendidikan.NewKepegawaianPendidikanSeeder(db), *fresh)

	// Kepegawaian - Kualifikasi
	runSeeder("tipe kualifikasi", masterKepegawaianKualifikasi.NewTipeSeeder(db), *fresh)
	runSeeder("kualifikasi", masterKepegawaianKualifikasi.NewKepegawaianKualifikasiSeeder(db), *fresh)

	// Kepegawaian - Alamat
	runSeeder("tipe alamat", masterKepegawaianAlamat.NewTipeSeeder(db), *fresh)
	runSeeder("alamat", masterKepegawaianAlamat.NewKepegawaianAlamatSeeder(db), *fresh)

	// Kepegawaian - Jabatan (urutan mengikuti dependensi FK)
	runSeeder("jabatan - posisi kategori", masterKepegawaianJabatan.NewPositionKategoriSeeder(db), *fresh)
	runSeeder("jabatan - posisi", masterKepegawaianJabatan.NewPositionSeeder(db), *fresh)
	runSeeder("jabatan - job title - kategori", masterKepegawaianJabatan.NewJobTitleKategoriSeeder(db), *fresh)
	runSeeder("jabatan - job title - rumpun_profesi", masterKepegawaianJabatan.NewJobTitleRumpunProfesiSeeder(db), *fresh)
	runSeeder("jabatan - job title", masterKepegawaianJabatan.NewJobTitleSeeder(db), *fresh)
	runSeeder("jabatan - spesialisasi", masterKepegawaianJabatan.NewSpecializationSeeder(db), *fresh)
	runSeeder("jabatan", masterKepegawaianJabatan.NewKepegawaianJabatanSeeder(db), *fresh)
	// Kepegawaian==============================================================

	// Artikel------------------------------------------------------------------
	runSeeder("Artikel", masterArtikel.NewArtikelSeeder(db), *fresh)
	runSeeder("Artikel Kategori", masterArtikelKategori.NewArtikelKategoriSeeder(db), *fresh)
	runSeeder("Artikel Tag", masterArtikelKategori.NewTagSeeder(db), *fresh)
	// Artikel------------------------------------------------------------------

	log.Println("✅ Seeding selesai!")
}
