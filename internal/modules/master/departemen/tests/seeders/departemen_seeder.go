package seeders

import (
	_ "embed"
	"encoding/json"
	"log"

	"neosim_go/internal/modules/master/departemen/models"
	"neosim_go/internal/modules/master/departemen/tests/factories"

	"gorm.io/gorm"
)

//go:embed departemen.json
var departemenJSON []byte

type jsonDepartemen struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	SystemModule string `json:"system_module"`
}

// MasterDepartemenSeeder mengelola seeding data MasterDepartemen
type MasterDepartemenSeeder struct {
	db *gorm.DB
}

func NewMasterDepartemenSeeder(db *gorm.DB) *MasterDepartemenSeeder {
	return &MasterDepartemenSeeder{db: db}
}

func (s *MasterDepartemenSeeder) Run() error {
	log.Println("🌱 Seeding master_departemens...")

	// 1. Pekerjaan
	//----------------------------------------------------------------------------------
	var departemenItems []models.MasterDepartemen
	if err := json.Unmarshal(departemenJSON, &departemenItems); err != nil {
		return err
	}
	// Peta/Map untuk menyimpan referensi ID berdasarkan Code (guna efisiensi pencarian relasi)
	mapDepartemen := make(map[string]int64)

	for _, item := range departemenItems {
		// Cek eksistensi agar tidak duplikat jika dijalankan ulang
		var existing models.MasterDepartemen
		err := s.db.Where("name = ?", item.Name).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := s.db.Create(&item).Error; err != nil {
				log.Printf("   ⚠️  Gagal membuat Departemen: %v", err)
				continue
			}
			mapDepartemen[item.Name] = item.ID
			log.Printf("   ✅ Departemen '%s' dibuat.", item.Name)
		} else {
			mapDepartemen[existing.Name] = existing.ID
		}
	}
	log.Println("   ✅ Seeding Departemen selesai.")
	//----------------------------------------------------------------------------------

	log.Println("✅ master_departemens seeding selesai!")
	return nil
}

func (s *MasterDepartemenSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data master_departemens...")
	if err := s.db.Exec("DELETE FROM master_departemens").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE master_departemens_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}

func (s *MasterDepartemenSeeder) seedDefault(name string) error {
	var count int64
	s.db.Model(&models.MasterDepartemen{}).Where("name = ?", name).Count(&count)
	if count > 0 {
		log.Printf("   ⏭️  '%s' sudah ada, skip.", name)
		return nil
	}
	item := factories.NewMasterDepartemenFactory().With("name", name).Make()
	if err := s.db.Create(item).Error; err != nil {
		return err
	}
	log.Printf("   ✅ '%s' dibuat.", name)
	return nil
}
