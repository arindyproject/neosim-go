package factories

import (
	"fmt"
	"math/rand"
	"time"

	"neosim_go/internal/modules/kepegawaian/kontak/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// KepegawaianKontakFactory membuat data KepegawaianKontak untuk testing/seeding
type KepegawaianKontakFactory struct {
	overrides map[string]interface{}
}

func NewKepegawaianKontakFactory() *KepegawaianKontakFactory {
	return &KepegawaianKontakFactory{overrides: make(map[string]interface{})}
}

func (f *KepegawaianKontakFactory) With(field string, value interface{}) *KepegawaianKontakFactory {
	f.overrides[field] = value
	return f
}

func (f *KepegawaianKontakFactory) Make() *models.KepegawaianKontak {
	idx := rng.Intn(999999)
	pegawaiID := int64(rng.Intn(10) + 1)
	tipeID := int64(rng.Intn(4) + 1)
	nilai := fmt.Sprintf("3512345678%06d", idx)
	desc := fmt.Sprintf("Deskripsi KepegawaianKontak %d", idx)
	createdBy := int64(rng.Intn(99) + 1)
	updatedBy := int64(rng.Intn(99) + 1)

	return &models.KepegawaianKontak{
		PegawaiID:   pegawaiID,
		TipeID:      tipeID,
		Nilai:       nilai,
		IsPrimary:   true,
		IsAktif:     true,
		Description: &desc,
		CreatedBy:   &createdBy,
		UpdatedBy:   &updatedBy,
	}
}

func (f *KepegawaianKontakFactory) MakeMany(count int) []*models.KepegawaianKontak {
	items := make([]*models.KepegawaianKontak, count)
	for i := 0; i < count; i++ {
		items[i] = NewKepegawaianKontakFactory().Make()
	}
	return items
}
