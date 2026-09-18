package factories

import (
	"fmt"
	"math/rand"
	"time"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// KepegawaianJabatanFactory membuat data KepegawaianJabatan untuk testing/seeding
type KepegawaianJabatanFactory struct {
	overrides map[string]interface{}
}

func NewKepegawaianJabatanFactory() *KepegawaianJabatanFactory {
	return &KepegawaianJabatanFactory{overrides: make(map[string]interface{})}
}

func (f *KepegawaianJabatanFactory) With(field string, value interface{}) *KepegawaianJabatanFactory {
	f.overrides[field] = value
	return f
}

func (f *KepegawaianJabatanFactory) Make() *models.KepegawaianJabatan {
	idx := rng.Intn(999999)
	name := fmt.Sprintf("KepegawaianJabatan %d", idx)
	desc := fmt.Sprintf("Deskripsi KepegawaianJabatan %d", idx)

	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	createdBy := int64(rng.Intn(99) + 1)
	updatedBy := int64(rng.Intn(99) + 1)

	return &models.KepegawaianJabatan{
		Name:        name,
		Description: &desc,
		CreatedBy:      &createdBy,
		UpdatedBy:      &updatedBy,
	}
}

func (f *KepegawaianJabatanFactory) MakeMany(count int) []*models.KepegawaianJabatan {
	items := make([]*models.KepegawaianJabatan, count)
	for i := 0; i < count; i++ {
		items[i] = NewKepegawaianJabatanFactory().Make()
	}
	return items
}
