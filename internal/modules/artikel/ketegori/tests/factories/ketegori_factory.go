package factories

import (
	"fmt"
	"math/rand"
	"time"

	"neosim_go/internal/modules/artikel/ketegori/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// ArtikelKetegoriFactory membuat data ArtikelKetegori untuk testing/seeding
type ArtikelKetegoriFactory struct {
	overrides map[string]interface{}
}

func NewArtikelKetegoriFactory() *ArtikelKetegoriFactory {
	return &ArtikelKetegoriFactory{overrides: make(map[string]interface{})}
}

func (f *ArtikelKetegoriFactory) With(field string, value interface{}) *ArtikelKetegoriFactory {
	f.overrides[field] = value
	return f
}

func (f *ArtikelKetegoriFactory) Make() *models.ArtikelKetegori {
	idx := rng.Intn(999999)
	name := fmt.Sprintf("ArtikelKetegori %d", idx)
	desc := fmt.Sprintf("Deskripsi ArtikelKetegori %d", idx)

	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.ArtikelKetegori{
		Name:        name,
		Description: &desc,
	}
}

func (f *ArtikelKetegoriFactory) MakeMany(count int) []*models.ArtikelKetegori {
	items := make([]*models.ArtikelKetegori, count)
	for i := 0; i < count; i++ {
		items[i] = NewArtikelKetegoriFactory().Make()
	}
	return items
}
