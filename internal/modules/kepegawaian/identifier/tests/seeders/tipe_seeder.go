package seeders

import (
	"log"

	"neosim_go/internal/modules/kepegawaian/identifier/models"

	"gorm.io/gorm"
)

// Const Identifier Type Code
const (
	IdentifierNIK     = "NIK"
	IdentifierIHS     = "IHS_NUMBER"
	IdentifierSTR     = "STR"
	IdentifierSIP     = "SIP"
	IdentifierNPAIDI  = "NPA_IDI"
	IdentifierNPAIAI  = "NPA_IAI"
	IdentifierNPAIBI  = "NPA_IBI"
	IdentifierNPAPPNI = "NPA_PPNI"
	IdentifierNPWP    = "NPWP"
	IdentifierBPJSKes = "BPJS_KES"
	IdentifierBPJSTK  = "BPJS_TK"
)

type TipeSeeder struct {
	db *gorm.DB
}

func NewTipeSeeder(db *gorm.DB) *TipeSeeder {
	return &TipeSeeder{db: db}
}

// GetDefaultData mengembalikan daftar master data preset
func GetDefaultData() []models.Tipe {
	strPtr := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}

	creatorID := int64(1)

	return []models.Tipe{
		{
			Code:       IdentifierNIK,
			Label:      "NIK",
			Penerbit:   strPtr("Dukcapil / Kemendagri"),
			FHIRSystem: strPtr("https://fhir.kemkes.go.id/id/nik"),
			HasExpiry:  false,
			IsNakes:    false,
			IsRequired: true,
			CreatedBy:  &creatorID,
			UpdatedBy:  &creatorID,
		},
		{
			Code:       IdentifierIHS,
			Label:      "IHS Number (SATUSEHAT)",
			Penerbit:   strPtr("Kemenkes RI"),
			FHIRSystem: strPtr("https://fhir.kemkes.go.id/id/ihs-number"),
			HasExpiry:  false,
			IsNakes:    true,
			IsRequired: false,
			CreatedBy:  &creatorID,
			UpdatedBy:  &creatorID,
		},
		{
			Code:       IdentifierSTR,
			Label:      "Surat Tanda Registrasi",
			Penerbit:   strPtr("Konsil Profesi"),
			FHIRSystem: strPtr("https://fhir.kemkes.go.id/id/str"),
			HasExpiry:  true,
			IsNakes:    true,
			IsRequired: false,
			CreatedBy:  &creatorID,
			UpdatedBy:  &creatorID,
		},
		{
			Code:       IdentifierSIP,
			Label:      "Surat Izin Praktik",
			Penerbit:   strPtr("Dinkes Kabupaten/Kota"),
			FHIRSystem: strPtr("https://fhir.kemkes.go.id/id/sip"),
			HasExpiry:  true,
			IsNakes:    true,
			IsRequired: false,
			CreatedBy:  &creatorID,
			UpdatedBy:  &creatorID,
		},
		{
			Code:       IdentifierNPAIDI,
			Label:      "NPA IDI",
			Penerbit:   strPtr("Ikatan Dokter Indonesia"),
			FHIRSystem: nil,
			HasExpiry:  false,
			IsNakes:    true,
			IsRequired: false,
			CreatedBy:  &creatorID,
			UpdatedBy:  &creatorID,
		},
		{
			Code:       IdentifierNPAIAI,
			Label:      "NPA IAI",
			Penerbit:   strPtr("Ikatan Apoteker Indonesia"),
			FHIRSystem: nil,
			HasExpiry:  false,
			IsNakes:    true,
			IsRequired: false,
			CreatedBy:  &creatorID,
			UpdatedBy:  &creatorID,
		},
		{
			Code:       IdentifierNPAIBI,
			Label:      "NPA IBI",
			Penerbit:   strPtr("Ikatan Bidan Indonesia"),
			FHIRSystem: nil,
			HasExpiry:  false,
			IsNakes:    true,
			IsRequired: false,
			CreatedBy:  &creatorID,
			UpdatedBy:  &creatorID,
		},
		{
			Code:       IdentifierNPAPPNI,
			Label:      "NPA PPNI",
			Penerbit:   strPtr("Persatuan Perawat Nasional Indonesia"),
			FHIRSystem: nil,
			HasExpiry:  false,
			IsNakes:    true,
			IsRequired: false,
			CreatedBy:  &creatorID,
			UpdatedBy:  &creatorID,
		},
		{
			Code:       IdentifierNPWP,
			Label:      "NPWP",
			Penerbit:   strPtr("Dirjen Pajak"),
			FHIRSystem: nil,
			HasExpiry:  false,
			IsNakes:    false,
			IsRequired: false,
			CreatedBy:  &creatorID,
			UpdatedBy:  &creatorID,
		},
		{
			Code:       IdentifierBPJSKes,
			Label:      "BPJS Kesehatan",
			Penerbit:   strPtr("BPJS Kesehatan"),
			FHIRSystem: nil,
			HasExpiry:  false,
			IsNakes:    false,
			IsRequired: false,
			CreatedBy:  &creatorID,
			UpdatedBy:  &creatorID,
		},
		{
			Code:       IdentifierBPJSTK,
			Label:      "BPJS Ketenagakerjaan",
			Penerbit:   strPtr("BPJS Ketenagakerjaan"),
			FHIRSystem: nil,
			HasExpiry:  false,
			IsNakes:    false,
			IsRequired: false,
			CreatedBy:  &creatorID,
			UpdatedBy:  &creatorID,
		},
	}
}

func (s *TipeSeeder) Run() error {
	log.Println("🌱 Seeding kepegawaian_identifier_tipes...")

	defaults := GetDefaultData()
	for _, item := range defaults {
		// Menggunakan FirstOrCreate berdasarkan `code` agar idempotent
		var existing models.Tipe
		err := s.db.Where("code = ?", item.Code).FirstOrCreate(&existing, item).Error
		if err != nil {
			log.Printf("   ⚠️ Gagal membuat/memeriksa Tipe [%s]: %v", item.Code, err)
			continue
		}
		log.Printf("   ✅ Tipe '%s' (%s) siap.", item.Label, item.Code)
	}

	log.Println("✅ kepegawaian_identifier_tipes seeding selesai!")
	return nil
}

func (s *TipeSeeder) Fresh() error {
	log.Println("🗑️  Menghapus semua data kepegawaian_identifier_tipes...")
	if err := s.db.Exec("DELETE FROM kepegawaian_identifier_tipes").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE kepegawaian_identifier_tipes_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Gagal reset sequence: %v", err)
	}
	return s.Run()
}
