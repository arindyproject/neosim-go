package factories

import (
	"fmt"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"
)

// SpecializationKategoriFactory membuat data SpecializationKategori untuk testing/seeding.
// Memakai 'rng' package-level yang sudah dideklarasikan di factory entitas
// utama sub-module ini.
type SpecializationKategoriFactory struct {
	overrides map[string]interface{}
}

func NewSpecializationKategoriFactory() *SpecializationKategoriFactory {
	return &SpecializationKategoriFactory{overrides: make(map[string]interface{})}
}

func (f *SpecializationKategoriFactory) With(field string, value interface{}) *SpecializationKategoriFactory {
	f.overrides[field] = value
	return f
}

func (f *SpecializationKategoriFactory) Make() *models.SpecializationKategori {
	idx := rng.Intn(999999)
	code := fmt.Sprintf("Code %d", idx)
	label := fmt.Sprintf("Label %d", idx)

	if v, ok := f.overrides["code"]; ok {
		code = v.(string)
	}
	if v, ok := f.overrides["label"]; ok {
		label = v.(string)
	}

	return &models.SpecializationKategori{
		Code:  code,
		Label: label,
	}
}

func (f *SpecializationKategoriFactory) MakeMany(count int) []*models.SpecializationKategori {
	items := make([]*models.SpecializationKategori, count)
	for i := 0; i < count; i++ {
		items[i] = NewSpecializationKategoriFactory().Make()
	}
	return items
}
