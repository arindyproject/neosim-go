package factories

import (
	"fmt"
	"math/rand"
	"time"

	"neosim_go/internal/modules/artikel/kategori/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// ArtikelKategoriFactory membuat data ArtikelKategori untuk testing/seeding
type ArtikelKategoriFactory struct {
	overrides map[string]interface{}
}

func NewArtikelKategoriFactory() *ArtikelKategoriFactory {
	return &ArtikelKategoriFactory{overrides: make(map[string]interface{})}
}

func (f *ArtikelKategoriFactory) With(field string, value interface{}) *ArtikelKategoriFactory {
	f.overrides[field] = value
	return f
}

func (f *ArtikelKategoriFactory) Make() *models.ArtikelKategori {
	idx := rng.Intn(999999)
	name := fmt.Sprintf("ArtikelKategori %d", idx)
	desc := fmt.Sprintf("Deskripsi ArtikelKategori %d", idx)
	user_id := rand.Int63n(100) + 1 // ID user acak antara 1 dan 300
	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.ArtikelKategori{
		Name:        name,
		Description: &desc,
		CreatedBy:   &user_id,
		UpdatedBy:   &user_id,
	}
}

func (f *ArtikelKategoriFactory) MakeMany(count int) []*models.ArtikelKategori {
	items := make([]*models.ArtikelKategori, count)
	for i := 0; i < count; i++ {
		items[i] = NewArtikelKategoriFactory().Make()
	}
	return items
}
