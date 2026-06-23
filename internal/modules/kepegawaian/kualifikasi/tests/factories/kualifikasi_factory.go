package factories

import (
	"fmt"
	"math/rand"
	"time"

	"neosim_go/internal/modules/kepegawaian/kualifikasi/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// KepegawaianKualifikasiFactory membuat data KepegawaianKualifikasi untuk testing/seeding
type KepegawaianKualifikasiFactory struct {
	overrides map[string]interface{}
}

func NewKepegawaianKualifikasiFactory() *KepegawaianKualifikasiFactory {
	return &KepegawaianKualifikasiFactory{overrides: make(map[string]interface{})}
}

func (f *KepegawaianKualifikasiFactory) With(field string, value interface{}) *KepegawaianKualifikasiFactory {
	f.overrides[field] = value
	return f
}

func (f *KepegawaianKualifikasiFactory) Make() *models.KepegawaianKualifikasi {
	idx := rng.Intn(999999)
	name := fmt.Sprintf("KepegawaianKualifikasi %d", idx)
	desc := fmt.Sprintf("Deskripsi KepegawaianKualifikasi %d", idx)

	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.KepegawaianKualifikasi{
		Name:        name,
		Description: &desc,
	}
}

func (f *KepegawaianKualifikasiFactory) MakeMany(count int) []*models.KepegawaianKualifikasi {
	items := make([]*models.KepegawaianKualifikasi, count)
	for i := 0; i < count; i++ {
		items[i] = NewKepegawaianKualifikasiFactory().Make()
	}
	return items
}
