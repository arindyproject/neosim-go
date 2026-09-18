package factories

import (
	"fmt"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"
)

// SpecializationFactory membuat data Specialization untuk testing/seeding.
// Memakai 'rng' package-level yang sudah dideklarasikan di factory entitas
// utama sub-module ini.
type SpecializationFactory struct {
	overrides map[string]interface{}
}

func NewSpecializationFactory() *SpecializationFactory {
	return &SpecializationFactory{overrides: make(map[string]interface{})}
}

func (f *SpecializationFactory) With(field string, value interface{}) *SpecializationFactory {
	f.overrides[field] = value
	return f
}

func (f *SpecializationFactory) Make() *models.Specialization {
	idx := rng.Intn(999999)
	name := fmt.Sprintf("Specialization %d", idx)
	desc := fmt.Sprintf("Deskripsi Specialization %d", idx)

	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.Specialization{
		Name:        name,
		Description: &desc,
	}
}

func (f *SpecializationFactory) MakeMany(count int) []*models.Specialization {
	items := make([]*models.Specialization, count)
	for i := 0; i < count; i++ {
		items[i] = NewSpecializationFactory().Make()
	}
	return items
}
