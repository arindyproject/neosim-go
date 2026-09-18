package factories

import (
	"fmt"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"
)

// JobTitleFactory membuat data JobTitle untuk testing/seeding.
// Memakai 'rng' package-level yang sudah dideklarasikan di factory entitas
// utama sub-module ini.
type JobTitleFactory struct {
	overrides map[string]interface{}
}

func NewJobTitleFactory() *JobTitleFactory {
	return &JobTitleFactory{overrides: make(map[string]interface{})}
}

func (f *JobTitleFactory) With(field string, value interface{}) *JobTitleFactory {
	f.overrides[field] = value
	return f
}

func (f *JobTitleFactory) Make() *models.JobTitle {
	idx := rng.Intn(999999)
	name := fmt.Sprintf("JobTitle %d", idx)
	desc := fmt.Sprintf("Deskripsi JobTitle %d", idx)

	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.JobTitle{
		Name:        name,
		Description: &desc,
	}
}

func (f *JobTitleFactory) MakeMany(count int) []*models.JobTitle {
	items := make([]*models.JobTitle, count)
	for i := 0; i < count; i++ {
		items[i] = NewJobTitleFactory().Make()
	}
	return items
}
