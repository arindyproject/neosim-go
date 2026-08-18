package factories

import (
	"fmt"
	"math/rand"

	"neosim_go/internal/modules/artikel/kategori/models"
)

// TagFactory membuat data Tag untuk testing/seeding.
// Memakai 'rng' package-level yang sudah dideklarasikan di factory entitas
// utama sub-module ini.
type TagFactory struct {
	overrides map[string]interface{}
}

func NewTagFactory() *TagFactory {
	return &TagFactory{overrides: make(map[string]interface{})}
}

func (f *TagFactory) With(field string, value interface{}) *TagFactory {
	f.overrides[field] = value
	return f
}

func (f *TagFactory) Make() *models.Tag {
	idx := rng.Intn(999999)
	name := fmt.Sprintf("Tag %d", idx)
	desc := fmt.Sprintf("Deskripsi Tag %d", idx)
	user_id := rand.Int63n(100) + 1 // ID user acak antara 1 dan 300
	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.Tag{
		Name:        name,
		Description: &desc,
		CreatedBy:   &user_id,
		UpdatedBy:   &user_id,
	}
}

func (f *TagFactory) MakeMany(count int) []*models.Tag {
	items := make([]*models.Tag, count)
	for i := 0; i < count; i++ {
		items[i] = NewTagFactory().Make()
	}
	return items
}
