package factories

import (
	"fmt"
	"math/rand"
	"time"

	"neosim_go/internal/modules/kepegawaian/identifier/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// KepegawaianIdentifierFactory membuat data KepegawaianIdentifier untuk testing/seeding
type KepegawaianIdentifierFactory struct {
	overrides map[string]interface{}
}

func NewKepegawaianIdentifierFactory() *KepegawaianIdentifierFactory {
	return &KepegawaianIdentifierFactory{overrides: make(map[string]interface{})}
}

func (f *KepegawaianIdentifierFactory) With(field string, value interface{}) *KepegawaianIdentifierFactory {
	f.overrides[field] = value
	return f
}

func (f *KepegawaianIdentifierFactory) Make() *models.KepegawaianIdentifier {
	// 1. Setup Default Values
	pegawaiID := int64(rng.Intn(1000) + 1)
	tipe := models.IdentifierNIK
	nilai := fmt.Sprintf("31710123%08d", rng.Intn(100000000))

	// Get metadata default
	meta, _ := tipe.Meta()
	penerbitStr := meta.Penerbit
	penerbit := &penerbitStr

	var tglTerbit *time.Time
	var tglExpired *time.Time

	// Jika tipe default/pilihan butuh tanggal kadaluarsa (misal STR/SIP)
	if meta.HasExpiry {
		terbit := time.Now().AddDate(-2, 0, 0) // 2 tahun lalu
		expired := time.Now().AddDate(3, 0, 0) // 3 tahun ke depan
		tglTerbit = &terbit
		tglExpired = &expired
	}

	isPrimary := true
	isAktif := true

	// 2. Apply Overrides
	if v, ok := f.overrides["PegawaiID"]; ok {
		pegawaiID = v.(int64)
	}
	if v, ok := f.overrides["Tipe"]; ok {
		tipe = v.(models.IdentifierType)
		// Jika Tipe di-override, sesuaikan juga metadata default-nya jika tidak di-override manual
		if metaNew, okMeta := tipe.Meta(); okMeta {
			if _, hasPenerbit := f.overrides["Penerbit"]; !hasPenerbit {
				penerbitStr = metaNew.Penerbit
				penerbit = &penerbitStr
			}
			if metaNew.HasExpiry && f.overrides["TanggalExpired"] == nil {
				terbit := time.Now().AddDate(-2, 0, 0)
				expired := time.Now().AddDate(3, 0, 0)
				tglTerbit = &terbit
				tglExpired = &expired
			}
		}
	}
	if v, ok := f.overrides["Nilai"]; ok {
		nilai = v.(string)
	}
	if v, ok := f.overrides["Penerbit"]; ok {
		if val, isStrPtr := v.(*string); isStrPtr {
			penerbit = val
		} else if valStr, isStr := v.(string); isStr {
			penerbit = &valStr
		}
	}
	if v, ok := f.overrides["TanggalTerbit"]; ok {
		tglTerbit = v.(*time.Time)
	}
	if v, ok := f.overrides["TanggalExpired"]; ok {
		tglExpired = v.(*time.Time)
	}
	if v, ok := f.overrides["IsPrimary"]; ok {
		isPrimary = v.(bool)
	}
	if v, ok := f.overrides["IsAktif"]; ok {
		isAktif = v.(bool)
	}

	// 3. Return Model Instance
	return &models.KepegawaianIdentifier{
		PegawaiID:      pegawaiID,
		Tipe:           tipe,
		Nilai:          nilai,
		Penerbit:       penerbit,
		TanggalTerbit:  tglTerbit,
		TanggalExpired: tglExpired,
		IsPrimary:      isPrimary,
		IsAktif:        isAktif,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func (f *KepegawaianIdentifierFactory) MakeMany(count int) []*models.KepegawaianIdentifier {
	items := make([]*models.KepegawaianIdentifier, count)
	for i := 0; i < count; i++ {
		// Menggunakan chain instance f agar override bawaan tetap diteruskan jika ada
		items[i] = f.Make()
	}
	return items
}
