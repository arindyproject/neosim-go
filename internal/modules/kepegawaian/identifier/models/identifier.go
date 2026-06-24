package models

import (
	"time"

	pegawai "neosim_go/internal/modules/kepegawaian/pegawai/models"

	"gorm.io/gorm"
)

// IdentifierType mendefinisikan tipe-tipe identifier pegawai yang valid
type IdentifierType string

const (
	IdentifierNIK     IdentifierType = "NIK"
	IdentifierIHS     IdentifierType = "IHS_NUMBER"
	IdentifierSTR     IdentifierType = "STR"
	IdentifierSIP     IdentifierType = "SIP"
	IdentifierNPAIDI  IdentifierType = "NPA_IDI"
	IdentifierNPAIAI  IdentifierType = "NPA_IAI"
	IdentifierNPAIBI  IdentifierType = "NPA_IBI"
	IdentifierNPAPPNI IdentifierType = "NPA_PPNI"
	IdentifierNPWP    IdentifierType = "NPWP"
	IdentifierBPJSKes IdentifierType = "BPJS_KES"
	IdentifierBPJSTK  IdentifierType = "BPJS_TK"
)

// IdentifierMeta menyimpan metadata statis untuk setiap tipe identifier
type IdentifierMeta struct {
	Code       IdentifierType
	Label      string
	Penerbit   string
	FHIRSystem string
	HasExpiry  bool
	IsNakes    bool // hanya untuk tenaga kesehatan
	IsRequired bool // wajib dimiliki semua pegawai
}

// identifierRegistry adalah registry metadata — tidak diekspor langsung
var identifierRegistry = map[IdentifierType]IdentifierMeta{
	IdentifierNIK: {
		Code:       IdentifierNIK,
		Label:      "NIK",
		Penerbit:   "Dukcapil / Kemendagri",
		FHIRSystem: "https://fhir.kemkes.go.id/id/nik",
		HasExpiry:  false,
		IsNakes:    false,
		IsRequired: true,
	},
	IdentifierIHS: {
		Code:       IdentifierIHS,
		Label:      "IHS Number (SATUSEHAT)",
		Penerbit:   "Kemenkes RI",
		FHIRSystem: "https://fhir.kemkes.go.id/id/ihs-number",
		HasExpiry:  false,
		IsNakes:    true,
		IsRequired: false,
	},
	IdentifierSTR: {
		Code:       IdentifierSTR,
		Label:      "Surat Tanda Registrasi",
		Penerbit:   "Konsil Profesi",
		FHIRSystem: "https://fhir.kemkes.go.id/id/str",
		HasExpiry:  true,
		IsNakes:    true,
		IsRequired: false,
	},
	IdentifierSIP: {
		Code:       IdentifierSIP,
		Label:      "Surat Izin Praktik",
		Penerbit:   "Dinkes Kabupaten/Kota",
		FHIRSystem: "https://fhir.kemkes.go.id/id/sip",
		HasExpiry:  true,
		IsNakes:    true,
		IsRequired: false,
	},
	IdentifierNPAIDI: {
		Code:       IdentifierNPAIDI,
		Label:      "NPA IDI",
		Penerbit:   "Ikatan Dokter Indonesia",
		FHIRSystem: "",
		HasExpiry:  false,
		IsNakes:    true,
		IsRequired: false,
	},
	IdentifierNPAIAI: {
		Code:       IdentifierNPAIAI,
		Label:      "NPA IAI",
		Penerbit:   "Ikatan Apoteker Indonesia",
		FHIRSystem: "",
		HasExpiry:  false,
		IsNakes:    true,
		IsRequired: false,
	},
	IdentifierNPAIBI: {
		Code:       IdentifierNPAIBI,
		Label:      "NPA IBI",
		Penerbit:   "Ikatan Bidan Indonesia",
		FHIRSystem: "",
		HasExpiry:  false,
		IsNakes:    true,
		IsRequired: false,
	},
	IdentifierNPAPPNI: {
		Code:       IdentifierNPAPPNI,
		Label:      "NPA PPNI",
		Penerbit:   "Persatuan Perawat Nasional Indonesia",
		FHIRSystem: "",
		HasExpiry:  false,
		IsNakes:    true,
		IsRequired: false,
	},
	IdentifierNPWP: {
		Code:       IdentifierNPWP,
		Label:      "NPWP",
		Penerbit:   "Dirjen Pajak",
		FHIRSystem: "",
		HasExpiry:  false,
		IsNakes:    false,
		IsRequired: false,
	},
	IdentifierBPJSKes: {
		Code:       IdentifierBPJSKes,
		Label:      "BPJS Kesehatan",
		Penerbit:   "BPJS Kesehatan",
		FHIRSystem: "",
		HasExpiry:  false,
		IsNakes:    false,
		IsRequired: false,
	},
	IdentifierBPJSTK: {
		Code:       IdentifierBPJSTK,
		Label:      "BPJS Ketenagakerjaan",
		Penerbit:   "BPJS Ketenagakerjaan",
		FHIRSystem: "",
		HasExpiry:  false,
		IsNakes:    false,
		IsRequired: false,
	},
}

// IsValid mengecek apakah tipe identifier dikenal oleh sistem
func (t IdentifierType) IsValid() bool {
	_, ok := identifierRegistry[t]
	return ok
}

// Meta mengembalikan metadata untuk tipe identifier ini
func (t IdentifierType) Meta() (IdentifierMeta, bool) {
	meta, ok := identifierRegistry[t]
	return meta, ok
}

// HasExpiry mengecek apakah tipe ini memiliki tanggal expired
func (t IdentifierType) HasExpiry() bool {
	meta, ok := identifierRegistry[t]
	return ok && meta.HasExpiry
}

// IsNakesOnly mengecek apakah tipe ini khusus tenaga kesehatan
func (t IdentifierType) IsNakesOnly() bool {
	meta, ok := identifierRegistry[t]
	return ok && meta.IsNakes
}

// AllIdentifierTypes mengembalikan semua tipe untuk keperluan dropdown UI
func AllIdentifierTypes() []IdentifierMeta {
	result := make([]IdentifierMeta, 0, len(identifierRegistry))
	for _, meta := range identifierRegistry {
		result = append(result, meta)
	}
	return result
}

// NakesIdentifierTypes mengembalikan hanya tipe khusus tenaga kesehatan
func NakesIdentifierTypes() []IdentifierMeta {
	result := make([]IdentifierMeta, 0)
	for _, meta := range identifierRegistry {
		if meta.IsNakes {
			result = append(result, meta)
		}
	}
	return result
}

// KepegawaianIdentifier represents the kepegawaian_identifiers table in database
// FHIR R4: Practitioner.identifier[]
type KepegawaianIdentifier struct {
	ID        int64 `gorm:"primaryKey;autoIncrement;column:id"                          json:"id"`
	PegawaiID int64 `gorm:"column:pegawai_id;not null;index"                        json:"pegawai_id"`

	// Tipe identifier — divalidasi via IdentifierType.IsValid()
	Tipe IdentifierType `gorm:"column:tipe;type:varchar(30);not null"                       json:"tipe"`

	// Nilai / nomor identifier
	Nilai string `gorm:"column:nilai;type:varchar(100);not null"                     json:"nilai"`

	// Penerbit (KKI, IDI, Dinkes, dll) — opsional, bisa diisi manual
	Penerbit *string `gorm:"column:penerbit;type:varchar(100)"                           json:"penerbit"`

	// Untuk STR dan SIP wajib diisi jika HasExpiry = true
	TanggalTerbit  *time.Time `gorm:"column:tanggal_terbit;type:date"                             json:"tanggal_terbit"`
	TanggalExpired *time.Time `gorm:"column:tanggal_expired;type:date"                            json:"tanggal_expired"`

	// Apakah identifier ini yang utama untuk tipe tersebut
	// Contoh: dokter punya 2 SIP, salah satunya is_primary = true
	IsPrimary bool `gorm:"column:is_primary;default:false"                             json:"is_primary"`
	IsAktif   bool `gorm:"column:is_aktif;default:true"                                json:"is_aktif"`

	// Audit fields — mengikuti konvensi model existing
	CreatedBy *int64         `gorm:"column:created_by"                                           json:"created_by"`
	UpdatedBy *int64         `gorm:"column:updated_by"                                           json:"updated_by"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:NOW()"   json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:NOW()"   json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz"                          json:"deleted_at"`

	// Relasi — preload saat dibutuhkan
	Pegawai *pegawai.KepegawaianPegawai `gorm:"foreignKey:PegawaiID"                                    json:"pegawai,omitempty"`
}

func (KepegawaianIdentifier) TableName() string {
	return "kepegawaian_identifiers"
}

// IsFHIRMappable mengecek apakah identifier ini bisa disync ke SATUSEHAT
func (ki *KepegawaianIdentifier) IsFHIRMappable() bool {
	meta, ok := ki.Tipe.Meta()
	return ok && meta.FHIRSystem != ""
}

// IsExpired mengecek apakah identifier ini sudah expired
func (ki *KepegawaianIdentifier) IsExpired() bool {
	if ki.TanggalExpired == nil {
		return false
	}
	return time.Now().After(*ki.TanggalExpired)
}

// DaysUntilExpired mengembalikan sisa hari sebelum expired, -1 jika tidak ada expired date
func (ki *KepegawaianIdentifier) DaysUntilExpired() int {
	if ki.TanggalExpired == nil {
		return -1
	}
	return int(time.Until(*ki.TanggalExpired).Hours() / 24)
}
