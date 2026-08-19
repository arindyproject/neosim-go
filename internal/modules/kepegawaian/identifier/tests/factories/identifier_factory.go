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
	idx := rng.Intn(999999)
	pegawaiID := int64(rng.Intn(10) + 1)
	tipeID := int64(rng.Intn(10) + 1)
	nilai := fmt.Sprintf("3512345678%06d", idx)

	now := time.Now()
	tanggalTerbit := now.AddDate(-1, 0, 0)
	tanggalExpired := now.AddDate(4, 0, 0)
	createdBy := int64(rng.Intn(99) + 1)
	updatedBy := int64(rng.Intn(99) + 1)

	item := &models.KepegawaianIdentifier{
		PegawaiID:      pegawaiID,
		TipeID:         tipeID,
		Nilai:          nilai,
		TanggalTerbit:  &tanggalTerbit,
		TanggalExpired: &tanggalExpired,
		IsPrimary:      true,
		IsAktif:        true,
		CreatedBy:      &createdBy,
		UpdatedBy:      &updatedBy,
	}

	// Apply overrides jika ada
	if v, ok := f.overrides["PegawaiID"]; ok {
		item.PegawaiID = v.(int64)
	}
	if v, ok := f.overrides["TipeID"]; ok {
		item.TipeID = v.(int64)
	}
	if v, ok := f.overrides["Nilai"]; ok {
		item.Nilai = v.(string)
	}

	if v, ok := f.overrides["TanggalTerbit"]; ok {
		item.TanggalTerbit = v.(*time.Time)
	}
	if v, ok := f.overrides["TanggalExpired"]; ok {
		item.TanggalExpired = v.(*time.Time)
	}
	if v, ok := f.overrides["IsPrimary"]; ok {
		item.IsPrimary = v.(bool)
	}
	if v, ok := f.overrides["IsAktif"]; ok {
		item.IsAktif = v.(bool)
	}
	if v, ok := f.overrides["Tipe"]; ok {
		item.Tipe = v.(*models.Tipe)
	}

	return item
}

func (f *KepegawaianIdentifierFactory) MakeMany(count int) []*models.KepegawaianIdentifier {
	items := make([]*models.KepegawaianIdentifier, count)
	for i := 0; i < count; i++ {
		items[i] = NewKepegawaianIdentifierFactory().Make()
	}
	return items
}
