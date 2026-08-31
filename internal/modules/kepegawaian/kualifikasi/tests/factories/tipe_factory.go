package factories

import (
	"fmt"

	"neosim_go/internal/modules/kepegawaian/kualifikasi/models"
)

// TipeFactory membuat data Tipe untuk testing/seeding.
// Memakai 'rng' package-level yang sudah dideklarasikan di factory entitas
// utama sub-module ini.
type TipeFactory struct {
	overrides map[string]interface{}
}

func NewTipeFactory() *TipeFactory {
	return &TipeFactory{overrides: make(map[string]interface{})}
}

func (f *TipeFactory) With(field string, value interface{}) *TipeFactory {
	f.overrides[field] = value
	return f
}

func (f *TipeFactory) Make() *models.Tipe {
	idx := rng.Intn(999999)
	name := fmt.Sprintf("Tipe %d", idx)

	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.Tipe{
		Code:  name,
		Label: name,
	}
}

func (f *TipeFactory) MakeMany(count int) []*models.Tipe {
	items := make([]*models.Tipe, count)
	for i := 0; i < count; i++ {
		items[i] = NewTipeFactory().Make()
	}
	return items
}
