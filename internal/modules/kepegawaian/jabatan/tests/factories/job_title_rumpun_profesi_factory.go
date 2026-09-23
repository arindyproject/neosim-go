package factories

import (
	"fmt"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"
)

// JobTitleRumpunProfesiFactory membuat data JobTitleRumpunProfesi untuk testing/seeding.
// Memakai 'rng' package-level yang sudah dideklarasikan di factory entitas
// utama sub-module ini.
type JobTitleRumpunProfesiFactory struct {
	overrides map[string]interface{}
}

func NewJobTitleRumpunProfesiFactory() *JobTitleRumpunProfesiFactory {
	return &JobTitleRumpunProfesiFactory{overrides: make(map[string]interface{})}
}

func (f *JobTitleRumpunProfesiFactory) With(field string, value interface{}) *JobTitleRumpunProfesiFactory {
	f.overrides[field] = value
	return f
}

func (f *JobTitleRumpunProfesiFactory) Make() *models.JobTitleRumpunProfesi {
	idx := rng.Intn(999999)
	code := fmt.Sprintf("Code %d", idx)
	label := fmt.Sprintf("Label %d", idx)

	if v, ok := f.overrides["code"]; ok {
		code = v.(string)
	}
	if v, ok := f.overrides["label"]; ok {
		label = v.(string)
	}

	return &models.JobTitleRumpunProfesi{
		Code:  code,
		Label: label,
	}
}

func (f *JobTitleRumpunProfesiFactory) MakeMany(count int) []*models.JobTitleRumpunProfesi {
	items := make([]*models.JobTitleRumpunProfesi, count)
	for i := 0; i < count; i++ {
		items[i] = NewJobTitleRumpunProfesiFactory().Make()
	}
	return items
}
