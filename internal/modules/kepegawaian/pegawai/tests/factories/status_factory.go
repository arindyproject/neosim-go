package factories

import (
	"fmt"

	"neosim_go/internal/modules/kepegawaian/pegawai/models"
)

// StatusFactory membuat data Status untuk testing/seeding.
// Memakai 'rng' package-level yang sudah dideklarasikan di factory entitas
// utama sub-module ini.
type StatusFactory struct {
	overrides map[string]interface{}
}

func NewStatusFactory() *StatusFactory {
	return &StatusFactory{overrides: make(map[string]interface{})}
}

func (f *StatusFactory) With(field string, value interface{}) *StatusFactory {
	f.overrides[field] = value
	return f
}

func (f *StatusFactory) Make() *models.Status {
	idx := rng.Intn(999999)
	code := fmt.Sprintf("Status %d", idx)
	label := fmt.Sprintf("Label Tipe %d", idx)

	if v, ok := f.overrides["code"]; ok {
		code = v.(string)
	}
	if v, ok := f.overrides["label"]; ok {
		label = v.(string)
	}

	return &models.Status{
		Code:  code,
		Label: label,
	}
}

func (f *StatusFactory) MakeMany(count int) []*models.Status {
	items := make([]*models.Status, count)
	for i := 0; i < count; i++ {
		items[i] = NewStatusFactory().Make()
	}
	return items
}
