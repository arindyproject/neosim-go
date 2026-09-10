package factories

import (
	"fmt"
	"math/rand"
	"time"

	"neosim_go/internal/modules/kepegawaian/alamat/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// KepegawaianAlamatFactory membuat data KepegawaianAlamat untuk testing/seeding
type KepegawaianAlamatFactory struct {
	overrides map[string]interface{}
}

func NewKepegawaianAlamatFactory() *KepegawaianAlamatFactory {
	return &KepegawaianAlamatFactory{overrides: make(map[string]interface{})}
}

func (f *KepegawaianAlamatFactory) With(field string, value interface{}) *KepegawaianAlamatFactory {
	f.overrides[field] = value
	return f
}

func (f *KepegawaianAlamatFactory) Make() *models.KepegawaianAlamat {
	idx := rng.Intn(999999)
	name := fmt.Sprintf("KepegawaianAlamat %d", idx)
	desc := fmt.Sprintf("Deskripsi KepegawaianAlamat %d", idx)

	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	createdBy := int64(rng.Intn(99) + 1)
	updatedBy := int64(rng.Intn(99) + 1)

	return &models.KepegawaianAlamat{
		Name:        name,
		Description: &desc,
		CreatedBy:      &createdBy,
		UpdatedBy:      &updatedBy,
	}
}

func (f *KepegawaianAlamatFactory) MakeMany(count int) []*models.KepegawaianAlamat {
	items := make([]*models.KepegawaianAlamat, count)
	for i := 0; i < count; i++ {
		items[i] = NewKepegawaianAlamatFactory().Make()
	}
	return items
}
