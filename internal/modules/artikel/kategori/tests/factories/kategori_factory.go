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

	//random 1 - 10 for createdBy and updatedBy
	var createdBy int64 = rng.Int63n(10) + 1
	var updatedBy int64 = rng.Int63n(10) + 1

	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}
	if v, ok := f.overrides["description"]; ok {
		desc = v.(string)
	}
	if v, ok := f.overrides["created_by"]; ok {
		createdBy = v.(int64)
	}
	if v, ok := f.overrides["updated_by"]; ok {
		updatedBy = v.(int64)
	}

	return &models.ArtikelKategori{
		Name:        name,
		Description: &desc,
		CreatedBy:   &createdBy,
		UpdatedBy:   &updatedBy,
	}
}

func (f *ArtikelKategoriFactory) MakeMany(count int) []*models.ArtikelKategori {
	items := make([]*models.ArtikelKategori, count)
	for i := 0; i < count; i++ {
		items[i] = NewArtikelKategoriFactory().Make()
	}
	return items
}
