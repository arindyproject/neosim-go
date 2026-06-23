package factories

import (
	"fmt"
	"math/rand"
	"time"

	"neosim_go/internal/modules/kepegawaian/identifikasi/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// KepegawaianIdentifikasiFactory membuat data KepegawaianIdentifikasi untuk testing/seeding
type KepegawaianIdentifikasiFactory struct {
	overrides map[string]interface{}
}

func NewKepegawaianIdentifikasiFactory() *KepegawaianIdentifikasiFactory {
	return &KepegawaianIdentifikasiFactory{overrides: make(map[string]interface{})}
}

func (f *KepegawaianIdentifikasiFactory) With(field string, value interface{}) *KepegawaianIdentifikasiFactory {
	f.overrides[field] = value
	return f
}

func (f *KepegawaianIdentifikasiFactory) Make() *models.KepegawaianIdentifikasi {
	idx := rng.Intn(999999)
	name := fmt.Sprintf("KepegawaianIdentifikasi %d", idx)
	desc := fmt.Sprintf("Deskripsi KepegawaianIdentifikasi %d", idx)

	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.KepegawaianIdentifikasi{
		Name:        name,
		Description: &desc,
	}
}

func (f *KepegawaianIdentifikasiFactory) MakeMany(count int) []*models.KepegawaianIdentifikasi {
	items := make([]*models.KepegawaianIdentifikasi, count)
	for i := 0; i < count; i++ {
		items[i] = NewKepegawaianIdentifikasiFactory().Make()
	}
	return items
}
