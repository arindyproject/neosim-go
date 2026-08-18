package factories

import (
	"fmt"
	"math/rand"
	"time"

	"neosim_go/internal/modules/artikel/artikel/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// ArtikelFactory membuat data Artikel untuk testing/seeding
type ArtikelFactory struct {
	overrides map[string]interface{}
}

func NewArtikelFactory() *ArtikelFactory {
	return &ArtikelFactory{overrides: make(map[string]interface{})}
}

func (f *ArtikelFactory) With(field string, value interface{}) *ArtikelFactory {
	f.overrides[field] = value
	return f
}

func (f *ArtikelFactory) Make() *models.Artikel {
	idx := rng.Intn(999999)
	name := fmt.Sprintf("Artikel %d", idx)
	desc := fmt.Sprintf("Deskripsi Artikel %d", idx)

	user_id := rand.Int63n(100) + 1 // ID user acak antara 1 dan 300

	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.Artikel{
		Name:        name,
		Description: &desc,
		CreatedBy:   &user_id,
		UpdatedBy:   &user_id,
	}
}

func (f *ArtikelFactory) MakeMany(count int) []*models.Artikel {
	items := make([]*models.Artikel, count)
	for i := 0; i < count; i++ {
		items[i] = NewArtikelFactory().Make()
	}
	return items
}
