package factories

import (
	"fmt"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"
)

// SpecializationFactory membuat data Specialization untuk testing/seeding.
// Memakai 'rng' package-level yang sudah dideklarasikan di factory entitas
// utama sub-module ini.
//
// Field yang bisa di-override lewat With():
// code, label, job_title_id, kategori_id, gelar, lama_pendidikan_tahun,
// fhir_code, fhir_system, is_aktif
type SpecializationFactory struct {
	overrides map[string]interface{}
}

func NewSpecializationFactory() *SpecializationFactory {
	return &SpecializationFactory{overrides: make(map[string]interface{})}
}

func (f *SpecializationFactory) With(field string, value interface{}) *SpecializationFactory {
	f.overrides[field] = value
	return f
}

func (f *SpecializationFactory) Make() *models.Specialization {
	idx := rng.Intn(999999)

	code := fmt.Sprintf("SP-%06d", idx)
	label := fmt.Sprintf("Spesialis %d", idx)
	gelar := fmt.Sprintf("Sp.%d", idx%100)
	lama := int16(4 + rng.Intn(3))       // 4-6 tahun
	jobTitleID := int64(rng.Intn(2) + 1) // 1 atau 2
	kategoriID := int64(rng.Intn(2) + 1) // 1 atau 2

	m := &models.Specialization{
		Code:                code,
		Label:               label,
		JobTitleID:          &jobTitleID,
		KategoriID:          &kategoriID,
		Gelar:               &gelar,
		LamaPendidikanTahun: &lama,
		IsAktif:             true,
	}

	if v, ok := f.overrides["code"]; ok {
		m.Code = v.(string)
	}
	if v, ok := f.overrides["label"]; ok {
		m.Label = v.(string)
	}
	if v, ok := f.overrides["job_title_id"]; ok {
		m.JobTitleID = toInt64Ptr(v)
	}
	if v, ok := f.overrides["kategori_id"]; ok {
		m.KategoriID = toInt64Ptr(v)
	}
	if v, ok := f.overrides["gelar"]; ok {
		m.Gelar = toStringPtr(v)
	}
	if v, ok := f.overrides["lama_pendidikan_tahun"]; ok {
		m.LamaPendidikanTahun = toInt16Ptr(v)
	}
	if v, ok := f.overrides["fhir_code"]; ok {
		m.FHIRCode = toStringPtr(v)
	}
	if v, ok := f.overrides["fhir_system"]; ok {
		m.FHIRSystem = toStringPtr(v)
	}
	if v, ok := f.overrides["is_aktif"]; ok {
		m.IsAktif = v.(bool)
	}

	return m
}

func (f *SpecializationFactory) MakeMany(count int) []*models.Specialization {
	items := make([]*models.Specialization, count)
	for i := 0; i < count; i++ {
		items[i] = NewSpecializationFactory().Make()
	}
	return items
}

// ── helper konversi override (menerima nilai langsung maupun pointer/nil) ────

func toInt64Ptr(v interface{}) *int64 {
	switch t := v.(type) {
	case int64:
		return &t
	case int:
		n := int64(t)
		return &n
	case *int64:
		return t
	}
	return nil
}

func toInt16Ptr(v interface{}) *int16 {
	switch t := v.(type) {
	case int16:
		return &t
	case int:
		n := int16(t)
		return &n
	case *int16:
		return t
	}
	return nil
}

func toStringPtr(v interface{}) *string {
	switch t := v.(type) {
	case string:
		return &t
	case *string:
		return t
	}
	return nil
}
