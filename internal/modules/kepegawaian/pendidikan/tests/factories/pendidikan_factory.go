package factories

import (
	"fmt"
	"math/rand"
	"time"

	"neosim_go/internal/modules/kepegawaian/pendidikan/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// KepegawaianPendidikanFactory membuat data KepegawaianPendidikan untuk testing/seeding
type KepegawaianPendidikanFactory struct {
	overrides map[string]interface{}
}

func NewKepegawaianPendidikanFactory() *KepegawaianPendidikanFactory {
	return &KepegawaianPendidikanFactory{overrides: make(map[string]interface{})}
}

func (f *KepegawaianPendidikanFactory) With(field string, value interface{}) *KepegawaianPendidikanFactory {
	f.overrides[field] = value
	return f
}

func (f *KepegawaianPendidikanFactory) Make() *models.KepegawaianPendidikan {
	idx := rng.Intn(999999)
	namaInstitusi := fmt.Sprintf("Institusi Pendidikan %d", idx)

	if v, ok := f.overrides["name"]; ok {
		namaInstitusi = v.(string)
	}

	createdBy := int64(rng.Intn(99) + 1)
	updatedBy := int64(rng.Intn(99) + 1)

	return &models.KepegawaianPendidikan{
		PegawaiID:     1,
		JenjangID:     1,
		NamaInstitusi: namaInstitusi,
		CreatedBy:     &createdBy,
		UpdatedBy:     &updatedBy,
	}
}

func (f *KepegawaianPendidikanFactory) MakeMany(count int) []*models.KepegawaianPendidikan {
	items := make([]*models.KepegawaianPendidikan, count)
	for i := 0; i < count; i++ {
		items[i] = NewKepegawaianPendidikanFactory().Make()
	}
	return items
}
