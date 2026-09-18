package factories

import (
	"fmt"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"
)

// PositionFactory membuat data Position untuk testing/seeding.
// Memakai 'rng' package-level yang sudah dideklarasikan di factory entitas
// utama sub-module ini.
type PositionFactory struct {
	overrides map[string]interface{}
}

func NewPositionFactory() *PositionFactory {
	return &PositionFactory{overrides: make(map[string]interface{})}
}

func (f *PositionFactory) With(field string, value interface{}) *PositionFactory {
	f.overrides[field] = value
	return f
}

func (f *PositionFactory) Make() *models.Position {
	idx := rng.Intn(999999)
	name := fmt.Sprintf("Position %d", idx)
	desc := fmt.Sprintf("Deskripsi Position %d", idx)

	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.Position{
		Name:        name,
		Description: &desc,
	}
}

func (f *PositionFactory) MakeMany(count int) []*models.Position {
	items := make([]*models.Position, count)
	for i := 0; i < count; i++ {
		items[i] = NewPositionFactory().Make()
	}
	return items
}
