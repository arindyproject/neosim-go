package factories

import (
	"fmt"
	"math/rand"
	"time"

	"neosim_go/internal/modules/master/master/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// MasterFactory membuat data Master untuk testing/seeding
type MasterFactory struct {
	overrides map[string]interface{}
}

func NewMasterFactory() *MasterFactory {
	return &MasterFactory{overrides: make(map[string]interface{})}
}

func (f *MasterFactory) With(field string, value interface{}) *MasterFactory {
	f.overrides[field] = value
	return f
}

func (f *MasterFactory) Make() *models.Master {
	idx := rng.Intn(999999)
	name := fmt.Sprintf("Master %d", idx)
	desc := fmt.Sprintf("Deskripsi Master %d", idx)

	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.Master{
		Name:        name,
		Description: &desc,
	}
}

func (f *MasterFactory) MakeMany(count int) []*models.Master {
	items := make([]*models.Master, count)
	for i := 0; i < count; i++ {
		items[i] = NewMasterFactory().Make()
	}
	return items
}
