package factories

import (
	"fmt"
	"math/rand"
	"time"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// KepegawaianJabatanFactory membuat data KepegawaianJabatan untuk testing/seeding.
//
// Field yang bisa di-override lewat With():
// pegawai_id, department_id, position_id, job_title_id, specialization_id,
// is_primary, tanggal_mulai, tanggal_selesai, nomor_sk, tanggal_sk, is_aktif
type KepegawaianJabatanFactory struct {
	overrides map[string]interface{}
}

func NewKepegawaianJabatanFactory() *KepegawaianJabatanFactory {
	return &KepegawaianJabatanFactory{overrides: make(map[string]interface{})}
}

func (f *KepegawaianJabatanFactory) With(field string, value interface{}) *KepegawaianJabatanFactory {
	f.overrides[field] = value
	return f
}

func (f *KepegawaianJabatanFactory) Make() *models.KepegawaianJabatan {
	idx := rng.Intn(999999)

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	// tanggal mulai acak dalam 1-5 tahun terakhir
	mulai := today.AddDate(0, 0, -(365 + rng.Intn(365*4)))
	tanggalSK := mulai
	nomorSK := fmt.Sprintf("SK/%d/%06d", mulai.Year(), idx)

	createdBy := int64(rng.Intn(99) + 1)
	updatedBy := int64(rng.Intn(99) + 1)

	m := &models.KepegawaianJabatan{
		PegawaiID:    int64(rng.Intn(5) + 1), // 1-5
		DepartmentID: int64(rng.Intn(2) + 1), // 1 atau 2
		PositionID:   int64(rng.Intn(2) + 1), // 1 atau 2
		JobTitleID:   int64(rng.Intn(2) + 1), // 1 atau 2
		IsPrimary:    false,                  // false agar tidak melanggar aturan 1 primer aktif per pegawai
		TanggalMulai: mulai,
		NomorSK:      &nomorSK,
		TanggalSK:    &tanggalSK,
		IsAktif:      true,
		CreatedBy:    &createdBy,
		UpdatedBy:    &updatedBy,
	}
	specID := int64(rng.Intn(2) + 1) // 1 atau 2
	m.SpecializationID = &specID

	if v, ok := f.overrides["pegawai_id"]; ok {
		m.PegawaiID = toInt64(v)
	}
	if v, ok := f.overrides["department_id"]; ok {
		m.DepartmentID = toInt64(v)
	}
	if v, ok := f.overrides["position_id"]; ok {
		m.PositionID = toInt64(v)
	}
	if v, ok := f.overrides["job_title_id"]; ok {
		m.JobTitleID = toInt64(v)
	}
	if v, ok := f.overrides["specialization_id"]; ok {
		m.SpecializationID = toInt64Ptr(v)
	}
	if v, ok := f.overrides["is_primary"]; ok {
		m.IsPrimary = v.(bool)
	}
	if v, ok := f.overrides["tanggal_mulai"]; ok {
		m.TanggalMulai = v.(time.Time)
	}
	if v, ok := f.overrides["tanggal_selesai"]; ok {
		m.TanggalSelesai = toTimePtr(v)
	}
	if v, ok := f.overrides["nomor_sk"]; ok {
		m.NomorSK = toStringPtr(v)
	}
	if v, ok := f.overrides["tanggal_sk"]; ok {
		m.TanggalSK = toTimePtr(v)
	}
	if v, ok := f.overrides["is_aktif"]; ok {
		m.IsAktif = v.(bool)
	}

	return m
}

func (f *KepegawaianJabatanFactory) MakeMany(count int) []*models.KepegawaianJabatan {
	items := make([]*models.KepegawaianJabatan, count)
	for i := 0; i < count; i++ {
		items[i] = NewKepegawaianJabatanFactory().Make()
	}
	return items
}

// ── helper konversi override ─────────────────────────────────────────────────
// toInt64Ptr dan toStringPtr sudah ada di factory Specialization (package yang sama).

func toInt64(v interface{}) int64 {
	switch t := v.(type) {
	case int64:
		return t
	case int:
		return int64(t)
	case *int64:
		if t != nil {
			return *t
		}
	}
	return 0
}

func toTimePtr(v interface{}) *time.Time {
	switch t := v.(type) {
	case time.Time:
		return &t
	case *time.Time:
		return t
	}
	return nil // termasuk With("tanggal_selesai", nil)
}
