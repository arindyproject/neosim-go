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
	name := fmt.Sprintf("SpecializationKategori %d", idx)
	desc := fmt.Sprintf("Deskripsi SpecializationKategori %d", idx)

	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.SpecializationKategori{
		Name:        name,
		Description: &desc,
	}
}

func (f *SpecializationKategoriFactory) MakeMany(count int) []*models.SpecializationKategori {
	items := make([]*models.SpecializationKategori, count)
	for i := 0; i < count; i++ {
		items[i] = NewSpecializationKategoriFactory().Make()
	}
	return items
}
