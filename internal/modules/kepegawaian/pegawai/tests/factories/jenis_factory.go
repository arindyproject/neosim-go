package factories

import (
	"fmt"

	"neosim_go/internal/modules/kepegawaian/pegawai/models"
)

// JenisFactory membuat data Jenis untuk testing/seeding.
// Memakai 'rng' package-level yang sudah dideklarasikan di factory entitas
// utama sub-module ini.
type JenisFactory struct {
	overrides map[string]interface{}
}

func NewJenisFactory() *JenisFactory {
	return &JenisFactory{overrides: make(map[string]interface{})}
}

func (f *JenisFactory) With(field string, value interface{}) *JenisFactory {
	f.overrides[field] = value
	return f
}

func (f *JenisFactory) Make() *models.Jenis {
	idx := rng.Intn(999999)
	code := fmt.Sprintf("Jenis %d", idx)
	label := fmt.Sprintf("Label Tipe %d", idx)

	if v, ok := f.overrides["code"]; ok {
		code = v.(string)
	}
	if v, ok := f.overrides["label"]; ok {
		label = v.(string)
	}

	return &models.Jenis{
		Code:  code,
		Label: label,
	}
}

func (f *JenisFactory) MakeMany(count int) []*models.Jenis {
	items := make([]*models.Jenis, count)
	for i := 0; i < count; i++ {
		items[i] = NewJenisFactory().Make()
	}
	return items
}
