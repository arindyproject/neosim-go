package seeders

import (
	_ "embed"
	"encoding/json"
	"log"

	"neosim_go/internal/modules/master/master/models"

	"gorm.io/gorm"
)

//go:embed pekerjaan.json
var pekerjaanJSON []byte

//go:embed pendidikan.json
var pendidikanJSON []byte

//go:embed agama.json
var agamaJSON []byte

//go:embed status_pernikahan.json
var statusPernikahanJSON []byte

type jsonPekerjaan struct {
	KodeKemenkes string `json:"kode_kemenkes"`
	Name         string `json:"name"`
}

type jsonPendidikan struct {
	KodeKemenkes string `json:"kode_kemenkes"`
	Name         string `json:"name"`
}

type jsonAgama struct {
	KodeKemenkes string `json:"kode_kemenkes"`
	Name         string `json:"name"`
}

type jsonStatusPernikahan struct {
	KodeKemenkes string `json:"kode_kemenkes"`
	Name         string `json:"name"`
}

// MasterSeeder mengelola seeding data Master
type MasterSeeder struct {
	db *gorm.DB
}

func NewMasterSeeder(db *gorm.DB) *MasterSeeder {
	return &MasterSeeder{db: db}
}

// Run menjalankan seeder
func (s *MasterSeeder) Run() error {
	log.Println("🌱 Seeding masters...")

	// 1. Pekerjaan
	var pekerjaanItems []models.MasterPekerjaan
	if err := json.Unmarshal(pekerjaanJSON, &pekerjaanItems); err != nil {
		return err
	}
	// Peta/Map untuk menyimpan referensi ID berdasarkan Code (guna efisiensi pencarian relasi)
	mapPekerjaan := make(map[string]int64)

	for _, item := range pekerjaanItems {
		// Cek eksistensi agar tidak duplikat jika dijalankan ulang
		var existing models.MasterPekerjaan
		err := s.db.Where("name = ?", item.Name).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := s.db.Create(&item).Error; err != nil {
				log.Printf("   ⚠️  Gagal membuat Pekerjaan: %v", err)
				continue
			}
			mapPekerjaan[item.Name] = item.ID
			log.Printf("   ✅ Pekerjaan '%s' dibuat.", item.Name)
		} else {
			mapPekerjaan[existing.Name] = existing.ID
		}
	}
	log.Println("   ✅ Seeding Pekerjaan selesai.")

	// 2. Pendidikan
	var pendidikanItems []models.MasterPendidikan
	if err := json.Unmarshal(pendidikanJSON, &pendidikanItems); err != nil {
		return err
	}
	// Peta/Map untuk menyimpan referensi ID berdasarkan Code (guna efisiensi pencarian relasi)
	mapPendidikan := make(map[string]int64)

	for _, item := range pendidikanItems {
		// Cek eksistensi agar tidak duplikat jika dijalankan ulang
		var existing models.MasterPendidikan
		err := s.db.Where("name = ?", item.Name).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := s.db.Create(&item).Error; err != nil {
				log.Printf("   ⚠️  Gagal membuat Pendidikan: %v", err)
				continue
			}
			mapPendidikan[item.Name] = item.ID
			log.Printf("   ✅ Pendidikan '%s' dibuat.", item.Name)
		} else {
			mapPendidikan[existing.Name] = existing.ID
		}
	}
	log.Println("   ✅ Seeding Pendidikan selesai.")

	// 3. Agama
	var agamaItems []models.MasterAgama
	if err := json.Unmarshal(agamaJSON, &agamaItems); err != nil {
		return err
	}
	// Peta/Map untuk menyimpan referensi ID berdasarkan Code (guna efisiensi pencarian relasi)
	mapAgama := make(map[string]int64)

	for _, item := range agamaItems {
		// Cek eksistensi agar tidak duplikat jika dijalankan ulang
		var existing models.MasterAgama
		err := s.db.Where("name = ?", item.Name).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := s.db.Create(&item).Error; err != nil {
				log.Printf("   ⚠️  Gagal membuat Agama: %v", err)
				continue
			}
			mapAgama[item.Name] = item.ID
			log.Printf("   ✅ Agama '%s' dibuat.", item.Name)
		} else {
			mapAgama[existing.Name] = existing.ID
		}
	}
	log.Println("   ✅ Seeding Agama selesai.")

	// 4. Status Pernikahan
	var statusPernikahanItems []models.MasterStatusPernikahan
	if err := json.Unmarshal(statusPernikahanJSON, &statusPernikahanItems); err != nil {
		return err
	}
	// Peta/Map untuk menyimpan referensi ID berdasarkan Code (guna efisiensi pencarian relasi)
	mapStatusPernikahan := make(map[string]int64)

	for _, item := range statusPernikahanItems {
		// Cek eksistensi agar tidak duplikat jika dijalankan ulang
		var existing models.MasterStatusPernikahan
		err := s.db.Where("name = ?", item.Name).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := s.db.Create(&item).Error; err != nil {
				log.Printf("   ⚠️  Gagal membuat Status Pernikahan: %v", err)
				continue
			}
			mapStatusPernikahan[item.Name] = item.ID
			log.Printf("   ✅ Status Pernikahan '%s' dibuat.", item.Name)
		} else {
			mapStatusPernikahan[existing.Name] = existing.ID
		}
	}
	log.Println("   ✅ Seeding Status Pernikahan selesai.")

	log.Println("   ✅ Seeding masters selesai.")

	return nil
}

// Fresh menghapus semua data lalu seed ulang
func (s *MasterSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data masters...")

	tables := []string{
		"master_pekerjaan",
		"master_pendidikan",
		"master_agama",
		"master_status_pernikahan",
	}

	for _, table := range tables {
		if err := s.db.Exec("DELETE FROM " + table).Error; err != nil {
			return err
		}
		// Reset serial id sequence di PostgreSQL
		if err := s.db.Exec("ALTER SEQUENCE " + table + "_id_seq RESTART WITH 1").Error; err != nil {
			log.Printf("Warning: Gagal reset sequence untuk %s: %v", table, err)
		}
	}

	return s.Run()
}
