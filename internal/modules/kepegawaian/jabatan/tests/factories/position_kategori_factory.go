package factories

import (
	"fmt"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"
)

// PositionKategoriFactory membuat data PositionKategori untuk testing/seeding.
// Memakai 'rng' package-level yang sudah dideklarasikan di factory entitas
// utama sub-module ini.
type PositionKategoriFactory struct {
	overrides map[string]interface{}
}

func NewPositionKategoriFactory() *PositionKategoriFactory {
	return &PositionKategoriFactory{overrides: make(map[string]interface{})}
}

func (f *PositionKategoriFactory) With(field string, value interface{}) *PositionKategoriFactory {
	f.overrides[field] = value
	return f
}

func (f *PositionKategoriFactory) Make() *models.PositionKategori {
	idx := rng.Intn(999999)
	code := fmt.Sprintf("Kategori %d", idx)
	label := fmt.Sprintf("Label Kategori %d", idx)

	if v, ok := f.overrides["code"]; ok {
		code = v.(string)
	}
	if v, ok := f.overrides["label"]; ok {
		label = v.(string)
	}

	return &models.PositionKategori{
		Code:  code,
		Label: label,
	}
}

func (f *PositionKategoriFactory) MakeMany(count int) []*models.PositionKategori {
	items := make([]*models.PositionKategori, count)
	for i := 0; i < count; i++ {
		items[i] = NewPositionKategoriFactory().Make()
	}
	return items
}
