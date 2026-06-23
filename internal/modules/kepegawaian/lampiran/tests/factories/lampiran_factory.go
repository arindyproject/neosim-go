package factories

import (
	"fmt"
	"math/rand"
	"time"

	"neosim_go/internal/modules/kepegawaian/lampiran/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// KepegawaianLampiranFactory membuat data KepegawaianLampiran untuk testing/seeding
type KepegawaianLampiranFactory struct {
	overrides map[string]interface{}
}

func NewKepegawaianLampiranFactory() *KepegawaianLampiranFactory {
	return &KepegawaianLampiranFactory{overrides: make(map[string]interface{})}
}

func (f *KepegawaianLampiranFactory) With(field string, value interface{}) *KepegawaianLampiranFactory {
	f.overrides[field] = value
	return f
}

func (f *KepegawaianLampiranFactory) Make() *models.KepegawaianLampiran {
	idx := rng.Intn(999999)
	name := fmt.Sprintf("KepegawaianLampiran %d", idx)
	desc := fmt.Sprintf("Deskripsi KepegawaianLampiran %d", idx)

	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.KepegawaianLampiran{
		Name:        name,
		Description: &desc,
	}
}

func (f *KepegawaianLampiranFactory) MakeMany(count int) []*models.KepegawaianLampiran {
	items := make([]*models.KepegawaianLampiran, count)
	for i := 0; i < count; i++ {
		items[i] = NewKepegawaianLampiranFactory().Make()
	}
	return items
}
