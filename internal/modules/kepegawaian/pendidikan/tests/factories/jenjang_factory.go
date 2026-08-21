package factories

import (
	"fmt"

	"neosim_go/internal/modules/kepegawaian/pendidikan/models"
)

// JenjangFactory membuat data Jenjang untuk testing/seeding.
// Memakai 'rng' package-level yang sudah dideklarasikan di factory entitas
// utama sub-module ini.
type JenjangFactory struct {
	overrides map[string]interface{}
}

func NewJenjangFactory() *JenjangFactory {
	return &JenjangFactory{overrides: make(map[string]interface{})}
}

func (f *JenjangFactory) With(field string, value interface{}) *JenjangFactory {
	f.overrides[field] = value
	return f
}

func (f *JenjangFactory) Make() *models.Jenjang {
	idx := rng.Intn(999999)
	code := fmt.Sprintf("JENJANG_%d", idx)
	label := fmt.Sprintf("Jenjang %d", idx)

	if v, ok := f.overrides["code"]; ok {
		code = v.(string)
	}
	if v, ok := f.overrides["label"]; ok {
		label = v.(string)
	}

	return &models.Jenjang{
		Code:  code,
		Label: label,
	}
}

func (f *JenjangFactory) MakeMany(count int) []*models.Jenjang {
	items := make([]*models.Jenjang, count)
	for i := 0; i < count; i++ {
		items[i] = NewJenjangFactory().Make()
	}
	return items
}
