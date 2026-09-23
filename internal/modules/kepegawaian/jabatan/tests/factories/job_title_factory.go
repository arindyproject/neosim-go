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

// Make membangun satu *models.JobTitle dari nilai random + override.
//
// KategoriID WAJIB diisi lewat .With("kategori_id", id) karena kolomnya
// NOT NULL FK ke kepegawaian_jabatan_job_title_kategoris — factory ini
// tidak query DB sendiri, jadi tidak bisa menebak ID kategori yang valid.
// Kalau tidak di-override, KategoriID akan 0 dan Create akan gagal karena
// FK constraint — disengaja, bukan bug, supaya caller (biasanya seeder)
// sadar harus sediakan kategori dulu.
func (f *JobTitleFactory) Make() *models.JobTitle {
	idx := rng.Intn(999999)
	code := fmt.Sprintf("JT-%06d", idx)
	label := fmt.Sprintf("JobTitle %d", idx)
	desc := fmt.Sprintf("Deskripsi JobTitle %d", idx)
	isAktif := true

	var kategoriID int64
	var rumpunProfesiID *int64
	var point *float64
	var memerlukanSTR, memerlukanSIP bool
	var jenjangMin *string
	var fhirCode, fhirSystem *string

	if v, ok := f.overrides["code"]; ok {
		if s, ok2 := v.(string); ok2 {
			code = s
		}
	}
	if v, ok := f.overrides["label"]; ok {
		if s, ok2 := v.(string); ok2 {
			label = s
		}
	}
	if v, ok := f.overrides["description"]; ok {
		if s, ok2 := v.(string); ok2 {
			desc = s
		}
	}
	if v, ok := f.overrides["kategori_id"]; ok {
		if id, ok2 := v.(int64); ok2 {
			kategoriID = id
		}
	}
	if v, ok := f.overrides["rumpun_profesi_id"]; ok {
		if id, ok2 := v.(*int64); ok2 {
			rumpunProfesiID = id
		} else if id, ok2 := v.(int64); ok2 {
			rumpunProfesiID = &id
		}
	}
	if v, ok := f.overrides["point"]; ok {
		if p, ok2 := v.(*float64); ok2 {
			point = p
		} else if p, ok2 := v.(float64); ok2 {
			point = &p
		}
	}
	if v, ok := f.overrides["memerlukan_str"]; ok {
		if b, ok2 := v.(bool); ok2 {
			memerlukanSTR = b
		}
	}
	if v, ok := f.overrides["memerlukan_sip"]; ok {
		if b, ok2 := v.(bool); ok2 {
			memerlukanSIP = b
		}
	}
	if v, ok := f.overrides["jenjang_min"]; ok {
		if s, ok2 := v.(*string); ok2 {
			jenjangMin = s
		} else if s, ok2 := v.(string); ok2 {
			jenjangMin = &s
		}
	}
	if v, ok := f.overrides["fhir_code"]; ok {
		if s, ok2 := v.(*string); ok2 {
			fhirCode = s
		} else if s, ok2 := v.(string); ok2 {
			fhirCode = &s
		}
	}
	if v, ok := f.overrides["fhir_system"]; ok {
		if s, ok2 := v.(*string); ok2 {
			fhirSystem = s
		} else if s, ok2 := v.(string); ok2 {
			fhirSystem = &s
		}
	}
	if v, ok := f.overrides["is_aktif"]; ok {
		if b, ok2 := v.(bool); ok2 {
			isAktif = b
		}
	}

	return &models.JobTitle{
		Code:            code,
		Label:           label,
		Description:     &desc,
		KategoriID:      kategoriID,
		RumpunProfesiID: rumpunProfesiID,
		Point:           point,
		MemerlukanSTR:   memerlukanSTR,
		MemerlukanSIP:   memerlukanSIP,
		JenjangMin:      jenjangMin,
		FHIRCode:        fhirCode,
		FHIRSystem:      fhirSystem,
		IsAktif:         isAktif,
	}
}

// MakeMany membuat 'count' JobTitle dari factory 'f' YANG SAMA, supaya
// override yang sudah di-set (mis. kategori_id) ikut terpakai di semua
// item — bukan bikin NewJobTitleFactory() baru tiap iterasi (itu bug yang
// sama seperti di PositionFactory sebelumnya, sudah diperbaiki di sini).
func (f *JobTitleFactory) MakeMany(count int) []*models.JobTitle {
	items := make([]*models.JobTitle, count)
	for i := 0; i < count; i++ {
		items[i] = f.Make()
	}
	return items
}
