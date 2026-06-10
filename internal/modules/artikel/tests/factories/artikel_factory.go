package factories

import (
	"fmt"
	"math/rand"
	"time"

	"neosim_go/internal/modules/artikel/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// ArtikelFactory membuat data artikel untuk testing/seeding
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
	desc := fmt.Sprintf("Deskripsi artikel %d", idx)

	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.Artikel{
		Name:        name,
		Description: &desc,
	}
}

func (f *ArtikelFactory) MakeMany(count int) []*models.Artikel {
	items := make([]*models.Artikel, count)
	for i := 0; i < count; i++ {
		items[i] = NewArtikelFactory().Make()
	}
	return items
}
