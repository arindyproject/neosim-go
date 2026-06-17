package factories

import (
	"fmt"
	"math/rand"
	"time"

	"neosim_go/internal/modules/master/alamat/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// MasterAlamatFactory membuat data MasterAlamat untuk testing/seeding
type MasterAlamatFactory struct {
	overrides map[string]interface{}
}

func NewMasterAlamatFactory() *MasterAlamatFactory {
	return &MasterAlamatFactory{overrides: make(map[string]interface{})}
}

func (f *MasterAlamatFactory) With(field string, value interface{}) *MasterAlamatFactory {
	f.overrides[field] = value
	return f
}

func (f *MasterAlamatFactory) Make() *models.MasterAlamat {
	idx := rng.Intn(999999)
	name := fmt.Sprintf("MasterAlamat %d", idx)
	desc := fmt.Sprintf("Deskripsi MasterAlamat %d", idx)

	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.MasterAlamat{
		Name:        name,
		Description: &desc,
	}
}

func (f *MasterAlamatFactory) MakeMany(count int) []*models.MasterAlamat {
	items := make([]*models.MasterAlamat, count)
	for i := 0; i < count; i++ {
		items[i] = NewMasterAlamatFactory().Make()
	}
	return items
}
