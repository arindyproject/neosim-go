package factories

import (
	"fmt"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"
)

// JobTitleKategoriFactory membuat data JobTitleKategori untuk testing/seeding.
// Memakai 'rng' package-level yang sudah dideklarasikan di factory entitas
// utama sub-module ini.
type JobTitleKategoriFactory struct {
	overrides map[string]interface{}
}

func NewJobTitleKategoriFactory() *JobTitleKategoriFactory {
	return &JobTitleKategoriFactory{overrides: make(map[string]interface{})}
}

func (f *JobTitleKategoriFactory) With(field string, value interface{}) *JobTitleKategoriFactory {
	f.overrides[field] = value
	return f
}

func (f *JobTitleKategoriFactory) Make() *models.JobTitleKategori {
	idx := rng.Intn(999999)
	code := fmt.Sprintf("Code %d", idx)
	label := fmt.Sprintf("Label %d", idx)

	if v, ok := f.overrides["code"]; ok {
		code = v.(string)
	}
	if v, ok := f.overrides["label"]; ok {
		label = v.(string)
	}

	return &models.JobTitleKategori{
		Code:  code,
		Label: label,
	}
}

func (f *JobTitleKategoriFactory) MakeMany(count int) []*models.JobTitleKategori {
	items := make([]*models.JobTitleKategori, count)
	for i := 0; i < count; i++ {
		items[i] = NewJobTitleKategoriFactory().Make()
	}
	return items
}
