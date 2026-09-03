package factories

import (
	"fmt"
	"math/rand"
	"time"

	"neosim_go/internal/modules/kepegawaian/kualifikasi/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// KepegawaianKualifikasiFactory membuat data KepegawaianKualifikasi untuk testing/seeding
type KepegawaianKualifikasiFactory struct {
	overrides map[string]interface{}
}

func NewKepegawaianKualifikasiFactory() *KepegawaianKualifikasiFactory {
	return &KepegawaianKualifikasiFactory{overrides: make(map[string]interface{})}
}

func (f *KepegawaianKualifikasiFactory) With(field string, value interface{}) *KepegawaianKualifikasiFactory {
	f.overrides[field] = value
	return f
}

func (f *KepegawaianKualifikasiFactory) Make() *models.KepegawaianKualifikasi {
	idx := rng.Intn(999999)

	nama := fmt.Sprintf("Kualifikasi %d", idx)
	penyelenggara := fmt.Sprintf("Penyelenggara %d", idx)
	nomorSertifikat := fmt.Sprintf("CERT-%06d", idx)

	pegawaiID := int64(rng.Intn(9) + 1)
	tipeID := int64(rng.Intn(3) + 1)

	terbit := time.Now().AddDate(0, -rng.Intn(24), 0) // 0–24 bulan lalu
	expired := terbit.AddDate(2, 0, 0)                // berlaku 2 tahun

	createdBy := int64(rng.Intn(49) + 1)
	updatedBy := int64(rng.Intn(49) + 1)

	m := &models.KepegawaianKualifikasi{
		PegawaiID:       pegawaiID,
		TipeID:          tipeID,
		Nama:            nama,
		Penyelenggara:   penyelenggara,
		NomorSertifikat: &nomorSertifikat,
		TanggalTerbit:   &terbit,
		TanggalExpired:  &expired,
		IsAktif:         true,
		CreatedBy:       &createdBy,
		UpdatedBy:       &updatedBy,
	}

	// Terapkan override, jika ada
	if v, ok := f.overrides["pegawai_id"]; ok {
		m.PegawaiID = v.(int64)
	}
	if v, ok := f.overrides["tipe_id"]; ok {
		m.TipeID = v.(int64)
	}
	if v, ok := f.overrides["nama"]; ok {
		m.Nama = v.(string)
	}
	if v, ok := f.overrides["penyelenggara"]; ok {
		m.Penyelenggara = v.(string)
	}
	if v, ok := f.overrides["nomor_sertifikat"]; ok {
		val := v.(string)
		m.NomorSertifikat = &val
	}
	if v, ok := f.overrides["tanggal_terbit"]; ok {
		val := v.(time.Time)
		m.TanggalTerbit = &val
	}
	if v, ok := f.overrides["tanggal_expired"]; ok {
		val := v.(time.Time)
		m.TanggalExpired = &val
	}
	if v, ok := f.overrides["is_aktif"]; ok {
		m.IsAktif = v.(bool)
	}
	if v, ok := f.overrides["fhir_code"]; ok {
		val := v.(string)
		m.FhirCode = &val
	}
	if v, ok := f.overrides["fhir_system"]; ok {
		val := v.(string)
		m.FhirSystem = &val
	}
	if v, ok := f.overrides["created_by"]; ok {
		val := v.(int64)
		m.CreatedBy = &val
	}
	if v, ok := f.overrides["updated_by"]; ok {
		val := v.(int64)
		m.UpdatedBy = &val
	}

	return m
}

func (f *KepegawaianKualifikasiFactory) MakeMany(count int) []*models.KepegawaianKualifikasi {
	items := make([]*models.KepegawaianKualifikasi, count)
	for i := 0; i < count; i++ {
		items[i] = f.Make()
	}
	return items
}
