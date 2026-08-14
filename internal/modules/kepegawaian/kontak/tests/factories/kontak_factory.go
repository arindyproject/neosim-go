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
	name := fmt.Sprintf("KepegawaianKontak %d", idx)
	desc := fmt.Sprintf("Deskripsi KepegawaianKontak %d", idx)

	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.KepegawaianKontak{
		Name:        name,
		Description: &desc,
	}
}

func (f *KepegawaianKontakFactory) MakeMany(count int) []*models.KepegawaianKontak {
	items := make([]*models.KepegawaianKontak, count)
	for i := 0; i < count; i++ {
		items[i] = NewKepegawaianKontakFactory().Make()
	}
	return items
}
