package factories

import (
	"fmt"
	"math/rand"
	"time"

	"neosim_go/internal/modules/kepegawaian/pegawai/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// KepegawaianPegawaiFactory membuat data KepegawaianPegawai untuk testing/seeding
type KepegawaianPegawaiFactory struct {
	overrides map[string]interface{}
}

func NewKepegawaianPegawaiFactory() *KepegawaianPegawaiFactory {
	return &KepegawaianPegawaiFactory{overrides: make(map[string]interface{})}
}

func (f *KepegawaianPegawaiFactory) With(field string, value interface{}) *KepegawaianPegawaiFactory {
	f.overrides[field] = value
	return f
}

func (f *KepegawaianPegawaiFactory) Make() *models.KepegawaianPegawai {
	idx := rng.Intn(999999)
	name := fmt.Sprintf("KepegawaianPegawai %d", idx)
	desc := fmt.Sprintf("Deskripsi KepegawaianPegawai %d", idx)

	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.KepegawaianPegawai{
		Name:        name,
		Description: &desc,
	}
}

func (f *KepegawaianPegawaiFactory) MakeMany(count int) []*models.KepegawaianPegawai {
	items := make([]*models.KepegawaianPegawai, count)
	for i := 0; i < count; i++ {
		items[i] = NewKepegawaianPegawaiFactory().Make()
	}
	return items
}
