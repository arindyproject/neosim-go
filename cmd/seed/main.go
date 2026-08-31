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
	masterKepegawaianIdentifier "neosim_go/internal/modules/kepegawaian/identifier/tests/seeders"
	masterKepegawaianKontak "neosim_go/internal/modules/kepegawaian/kontak/tests/seeders"
	masterKepegawaianKualifikasi "neosim_go/internal/modules/kepegawaian/kualifikasi/tests/seeders"
	masterKepegawaianPegawai "neosim_go/internal/modules/kepegawaian/pegawai/tests/seeders"
	masterKepegawaianPendidikan "neosim_go/internal/modules/kepegawaian/pendidikan/tests/seeders"

	// Artikel---------------------------------------------------------------
	masterArtikel "neosim_go/internal/modules/artikel/artikel/tests/seeders"
	masterArtikelKategori "neosim_go/internal/modules/artikel/kategori/tests/seeders"
	// =====================================================================
)

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
	// TODO: Tambahkan seeder lain di sini
	// =====================================================================
	// Jalankan semua seeder di sini
	// Contoh: user seeder
	// =====================================================================
	userSeeder := userSeed.NewUserSeeder(db)

	if *fresh {
		if err := userSeeder.Fresh(); err != nil {
			log.Fatal("Gagal fresh seed users:", err)
		}
	} else {
		if err := userSeeder.Run(); err != nil {
			log.Fatal("Gagal seed users:", err)
		}
	}
	// =====================================================================
	// Contoh: RBAC seeder
	rbacSeeder := rbacSeed.NewRBACSeeder(db)

	if *fresh {
		if err := rbacSeeder.Fresh(); err != nil {
			log.Fatal("Gagal fresh seed RBAC:", err)
		}
	} else {
		if err := rbacSeeder.Run(); err != nil {
			log.Fatal("Gagal seed RBAC:", err)
		}
	}
	// =====================================================================
	//Master----------------------------------------------------------------
	//------Alamat----------------------------------------------------------

	masterAlamatSeeder := masterAlamat.NewMasterAlamatSeeder(db)
	if *fresh {
		if err := masterAlamatSeeder.Fresh(); err != nil {
			log.Fatal("Gagal fresh seed RBAC:", err)
		}
	} else {
		if err := masterAlamatSeeder.Run(); err != nil {
			log.Fatal("Gagal seed RBAC:", err)
		}
	} //------Alamat--------------------------------------------------------

	//------Master--------------------------------------------------------
	masterMasterSeeder := masterMaster.NewMasterSeeder(db)
	if *fresh {
		if err := masterMasterSeeder.Fresh(); err != nil {
			log.Fatal("Gagal fresh seed RBAC:", err)
		}
	} else {
		if err := masterMasterSeeder.Run(); err != nil {
			log.Fatal("Gagal seed master:", err)
		}
	} //------Master--------------------------------------------------------

	//------Departemen----------------------------------------------------
	masterDepartemenSeeder := masterDepartemen.NewMasterDepartemenSeeder(db)
	if *fresh {
		if err := masterDepartemenSeeder.Fresh(); err != nil {
			log.Fatal("Gagal fresh seed Departemen:", err)
		}
	} else {
		if err := masterDepartemenSeeder.Run(); err != nil {
			log.Fatal("Gagal seed Departemen:", err)
		}
	} //------Departemen----------------------------------------------------

	//Kepegawaian-----------------------------------------------------------
	//------Kepegawaian-----------------------------------------------------
	masterKepegawaianPegawaiSeeder := masterKepegawaianPegawai.NewKepegawaianPegawaiSeeder(db)
	if *fresh {
		if err := masterKepegawaianPegawaiSeeder.Fresh(); err != nil {
			log.Fatal("Gagal fresh seed Pegawai:", err)
		}
	} else {
		if err := masterKepegawaianPegawaiSeeder.Run(); err != nil {
			log.Fatal("Gagal seed Pegawai:", err)
		}
	} //------Kepegawaian----------------------------------------------------

	//------Kepegawaian Identifier Type--------------------------------------
	masterKepegawaianIdentifierTypeSeeder := masterKepegawaianIdentifier.NewTipeSeeder(db)
	if *fresh {
		if err := masterKepegawaianIdentifierTypeSeeder.Fresh(); err != nil {
			log.Fatal("Gagal fresh seed Identifier Type:", err)
		}
	} else {
		if err := masterKepegawaianIdentifierTypeSeeder.Run(); err != nil {
			log.Fatal("Gagal seed Identifier Type:", err)
		}
	} //------Kepegawaian Identifier Type------------------------------------

	//------Kepegawaian Identifier-------------------------------------------
	masterKepegawaianIdentifierSeeder := masterKepegawaianIdentifier.NewKepegawaianIdentifierSeeder(db)
	if *fresh {
		if err := masterKepegawaianIdentifierSeeder.Fresh(); err != nil {
			log.Fatal("Gagal fresh seed Identifier:", err)
		}
	} else {
		if err := masterKepegawaianIdentifierSeeder.Run(); err != nil {
			log.Fatal("Gagal seed Identifier:", err)
		}
	} //------Kepegawaian Identifier-----------------------------------------

	//------Kepegawaian Kontak Tipe------------------------------------------
	masterKepegawaianKontakTipeSeeder := masterKepegawaianKontak.NewTipeSeeder(db)
	if *fresh {
		if err := masterKepegawaianKontakTipeSeeder.Fresh(); err != nil {
			log.Fatal("Gagal fresh seed Kontak Tipe:", err)
		}
	} else {
		if err := masterKepegawaianKontakTipeSeeder.Run(); err != nil {
			log.Fatal("Gagal seed Kontak Tipe:", err)
		}
	} //------Kepegawaian Kontak Tipe----------------------------------------

	//------Kepegawaian Kontak-----------------------------------------------
	masterKepegawaianKontakSeeder := masterKepegawaianKontak.NewKepegawaianKontakSeeder(db)
	if *fresh {
		if err := masterKepegawaianKontakSeeder.Fresh(); err != nil {
			log.Fatal("Gagal fresh seed Kontak:", err)
		}
	} else {
		if err := masterKepegawaianKontakSeeder.Run(); err != nil {
			log.Fatal("Gagal seed Kontak:", err)
		}
	} //------Kepegawaian Kontak---------------------------------------------

	//------Kepegawaian Pendidikan Jenjang-----------------------------------
	masterKepegawaianPendidikanJenjangSeeder := masterKepegawaianPendidikan.NewJenjangSeeder(db)
	if *fresh {
		if err := masterKepegawaianPendidikanJenjangSeeder.Fresh(); err != nil {
			log.Fatal("Gagal fresh seed jenjang pendidikan:", err)
		}
	} else {
		if err := masterKepegawaianPendidikanJenjangSeeder.Run(); err != nil {
			log.Fatal("Gagal seed jenjang pendidikan:", err)
		}
	} //------Kepegawaian Pendidikan Jenjang---------------------------------

	//------Kepegawaian Pendidikan-------------------------------------------
	masterKepegawaianPendidikanSeeder := masterKepegawaianPendidikan.NewKepegawaianPendidikanSeeder(db)
	if *fresh {
		if err := masterKepegawaianPendidikanSeeder.Fresh(); err != nil {
			log.Fatal("Gagal fresh seed  pendidikan:", err)
		}
	} else {
		if err := masterKepegawaianPendidikanSeeder.Run(); err != nil {
			log.Fatal("Gagal seed  pendidikan:", err)
		}
	} //------Kepegawaian Pendidikan-----------------------------------------

	//------Kepegawaian Kualifikasi Tipe-------------------------------------
	masterKepegawaianKualifikasiTipeSeeder := masterKepegawaianKualifikasi.NewTipeSeeder(db)
	if *fresh {
		if err := masterKepegawaianKualifikasiTipeSeeder.Fresh(); err != nil {
			log.Fatal("Gagal fresh seed  tipe kualifikasi:", err)
		}
	} else {
		if err := masterKepegawaianKualifikasiTipeSeeder.Run(); err != nil {
			log.Fatal("Gagal seed   tipe kualifikasi:", err)
		}
	} //------Kepegawaian Kualifikasi Tipe-----------------------------------

	//------Kepegawaian Kualifikasi -----------------------------------------
	masterKepegawaianKualifikasiSeeder := masterKepegawaianKualifikasi.NewKepegawaianKualifikasiSeeder(db)
	if *fresh {
		if err := masterKepegawaianKualifikasiSeeder.Fresh(); err != nil {
			log.Fatal("Gagal fresh seed   kualifikasi:", err)
		}
	} else {
		if err := masterKepegawaianKualifikasiSeeder.Run(); err != nil {
			log.Fatal("Gagal seed  kualifikasi:", err)
		}
	} //------Kepegawaian Kualifikasi ---------------------------------------

	//Kepegawaian------------------------------------------------------------

	// =====================================================================

	//------Artikel----------------------------------------------------
	masterArtikels := masterArtikel.NewArtikelSeeder(db)
	if *fresh {
		if err := masterArtikels.Fresh(); err != nil {
			log.Fatal("Gagal fresh seed Artikel:", err)
		}
	} else {
		if err := masterArtikels.Run(); err != nil {
			log.Fatal("Gagal seed Artikel:", err)
		}
	}
	masterArtikelKategoris := masterArtikelKategori.NewArtikelKategoriSeeder(db)
	if *fresh {
		if err := masterArtikelKategoris.Fresh(); err != nil {
			log.Fatal("Gagal fresh seed Artikel Kategori:", err)
		}
	} else {
		if err := masterArtikelKategoris.Run(); err != nil {
			log.Fatal("Gagal seed Artikel Kategori:", err)
		}
	}
	masterArtikelTags := masterArtikelKategori.NewTagSeeder(db)
	if *fresh {
		if err := masterArtikelTags.Fresh(); err != nil {
			log.Fatal("Gagal fresh seed Artikel Tag:", err)
		}
	} else {
		if err := masterArtikelTags.Run(); err != nil {
			log.Fatal("Gagal seed Artikel Tag:", err)
		}
	}
	//------Artikel----------------------------------------------------

	log.Println("✅ Seeding selesai!")
}
