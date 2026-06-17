package seeders

import (
	_ "embed"
	"encoding/json"
	"log"

	"neosim_go/internal/modules/master/alamat/models"

	"gorm.io/gorm"
)

// Ambil file JSON menggunakan go:embed (pastikan letak file json sesuai path relatifnya)
//
//go:embed negara.json
var negaraJSON []byte

//go:embed provinsi.json
var provinsiJSON []byte

//go:embed kota_kabupaten.json
var kotaKabupatenJSON []byte

//go:embed kecamatan.json
var kecamatanJSON []byte

//go:embed kelurahan_desa.json
var kelurahanDesaJSON []byte

// Struct pembantu untuk parsing JSON relasional
type jsonProvinsi struct {
	Code       string `json:"code"`
	NegaraCode string `json:"negara_code"`
	Name       string `json:"name"`
}

type jsonKota struct {
	Code         string `json:"code"`
	ProvinsiCode string `json:"provinsi_code"`
	Name         string `json:"name"`
}

type jsonKecamatan struct {
	Code              string `json:"code"`
	KotaKabupatenCode string `json:"kota_kabupaten_code"`
	Name              string `json:"name"`
}

type jsonDesa struct {
	Code          string `json:"code"`
	KecamatanCode string `json:"kecamatan_code"`
	Name          string `json:"name"`
	PostalCode    string `json:"postal_code"`
}

type MasterAlamatSeeder struct {
	db *gorm.DB
}

func NewMasterAlamatSeeder(db *gorm.DB) *MasterAlamatSeeder {
	return &MasterAlamatSeeder{db: db}
}

// Run menjalankan seeder alamat berjenjang
func (s *MasterAlamatSeeder) Run() error {
	log.Println("🌱 Seeding master alamat bertingkat...")

	// 1. SEED NEGARA
	var negaraItems []models.MasterAlamatNegara
	if err := json.Unmarshal(negaraJSON, &negaraItems); err != nil {
		return err
	}
	// Peta/Map untuk menyimpan referensi ID berdasarkan Code (guna efisiensi pencarian relasi)
	mapNegara := make(map[string]int64)

	for _, item := range negaraItems {
		// Cek eksistensi agar tidak duplikat jika dijalankan ulang
		var existing models.MasterAlamatNegara
		err := s.db.Where("code = ?", item.Code).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := s.db.Create(&item).Error; err != nil {
				log.Printf("   ⚠️  Gagal membuat Negara: %v", err)
				continue
			}
			mapNegara[item.Code] = item.ID
			log.Printf("   ✅ Negara '%s' dibuat.", item.Name)
		} else {
			mapNegara[existing.Code] = existing.ID
		}
	}

	// 2. SEED PROVINSI
	var provItems []jsonProvinsi
	if err := json.Unmarshal(provinsiJSON, &provItems); err != nil {
		return err
	}
	mapProvinsi := make(map[string]int64)

	for _, item := range provItems {
		negaraID, exists := mapNegara[item.NegaraCode]
		if !exists {
			continue
		}

		var existing models.MasterAlamatProvinsi
		err := s.db.Where("code = ?", item.Code).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			newProv := models.MasterAlamatProvinsi{
				NegaraID: negaraID,
				Code:     item.Code,
				Name:     item.Name,
			}
			if err := s.db.Create(&newProv).Error; err != nil {
				log.Printf("   ⚠️  Gagal membuat Provinsi: %v", err)
				continue
			}
			mapProvinsi[item.Code] = newProv.ID
		} else {
			mapProvinsi[existing.Code] = existing.ID
		}
	}
	log.Println("   ✅ Seeding Provinsi selesai.")

	// 3. SEED KOTA / KABUPATEN
	var kotaItems []jsonKota
	if err := json.Unmarshal(kotaKabupatenJSON, &kotaItems); err != nil {
		return err
	}
	mapKota := make(map[string]int64)

	for _, item := range kotaItems {
		provID, exists := mapProvinsi[item.ProvinsiCode]
		if !exists {
			continue
		}

		var existing models.MasterAlamatKotaKabupaten
		err := s.db.Where("code = ?", item.Code).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			newKota := models.MasterAlamatKotaKabupaten{
				ProvinsiID: provID,
				Code:       item.Code,
				Name:       item.Name,
			}
			if err := s.db.Create(&newKota).Error; err != nil {
				continue
			}
			mapKota[item.Code] = newKota.ID
		} else {
			mapKota[existing.Code] = existing.ID
		}
	}
	log.Println("   ✅ Seeding Kota/Kabupaten selesai.")

	// 4. SEED KECAMATAN
	var kecItems []jsonKecamatan
	if err := json.Unmarshal(kecamatanJSON, &kecItems); err != nil {
		return err
	}
	mapKecamatan := make(map[string]int64)

	for _, item := range kecItems {
		kotaID, exists := mapKota[item.KotaKabupatenCode]
		if !exists {
			continue
		}

		var existing models.MasterAlamatKecamatan
		err := s.db.Where("code = ?", item.Code).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			newKec := models.MasterAlamatKecamatan{
				KotaKabupatenID: kotaID,
				Code:            item.Code,
				Name:            item.Name,
			}
			if err := s.db.Create(&newKec).Error; err != nil {
				continue
			}
			mapKecamatan[item.Code] = newKec.ID
		} else {
			mapKecamatan[existing.Code] = existing.ID
		}
	}
	log.Println("   ✅ Seeding Kecamatan selesai.")

	// 5. SEED KELURAHAN / DESA
	var desaItems []jsonDesa
	if err := json.Unmarshal(kelurahanDesaJSON, &desaItems); err != nil {
		return err
	}

	// Karena jumlah data desa cukup banyak (1000+), gunakan GORM Batch Insert agar proses sangat cepat
	var batchDesa []models.MasterAlamatKelurahanDesa
	for _, item := range desaItems {
		kecID, exists := mapKecamatan[item.KecamatanCode]
		if !exists {
			continue
		}

		// Validasi eksistensi lokal/sederhana (Opsional, asalkan kode desa unik)
		var count int64
		s.db.Model(&models.MasterAlamatKelurahanDesa{}).Where("code = ?", item.Code).Count(&count)
		if count == 0 {
			postal := item.PostalCode
			batchDesa = append(batchDesa, models.MasterAlamatKelurahanDesa{
				KecamatanID: kecID,
				Code:        item.Code,
				Name:        item.Name,
				PostalCode:  &postal,
			})
		}
	}

	if len(batchDesa) > 0 {
		// Melakukan insert berkelompok (batching) per 200 data sekaligus agar hemat query I/O
		if err := s.db.CreateInBatches(&batchDesa, 200).Error; err != nil {
			log.Printf("   ⚠️  Gagal batch insert Kelurahan/Desa: %v", err)
			return err
		}
	}

	log.Println("✅ Semua tingkatan master alamat seeding selesai!")
	return nil
}

// Fresh menghapus semua data alamat dari level terbawah hingga teratas, lalu seed ulang
func (s *MasterAlamatSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data master alamat (Cascading style)...")

	// Urutan hapus dibalik agar tidak melanggar Foreign Key Constraints
	tables := []string{
		"master_alamat_kelurahan_desa",
		"master_alamat_kecamatan",
		"master_alamat_kota_kabupaten",
		"master_alamat_provinsi",
		"master_alamat_negara",
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
