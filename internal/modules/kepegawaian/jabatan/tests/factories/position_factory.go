package factories

import (
	"fmt"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"
)

// PositionFactory membuat data Position untuk testing/seeding.
// Memakai 'rng' package-level yang sudah dideklarasikan di factory entitas
// utama sub-module ini.
type PositionFactory struct {
	overrides map[string]interface{}
}

func NewPositionFactory() *PositionFactory {
	return &PositionFactory{overrides: make(map[string]interface{})}
}

func (f *PositionFactory) With(field string, value interface{}) *PositionFactory {
	f.overrides[field] = value
	return f
}

// Make membangun satu *models.Position dari nilai random + override.
//
// PositionKategoriID WAJIB diisi lewat .With("position_kategori_id", id)
// karena kolomnya NOT NULL FK ke kepegawaian_jabatan_position_kategoris —
// factory ini tidak query DB sendiri, jadi tidak bisa menebak ID kategori
// yang valid. Kalau tidak di-override, PositionKategoriID akan 0 dan Create
// akan gagal karena FK constraint — itu disengaja, bukan bug, supaya
// caller (biasanya seeder) sadar harus sediakan kategori dulu.
func (f *PositionFactory) Make() *models.Position {
	idx := rng.Intn(999999)
	name := fmt.Sprintf("Position %d", idx)
	desc := fmt.Sprintf("Deskripsi Position %d", idx)
	level := int16(rng.Intn(5) + 1)
	isAktif := true

	var kategoriID int64
	var parentID *int64
	var departmentID *int64
	var kuota *int16

	if v, ok := f.overrides["name"]; ok {
		if s, ok2 := v.(string); ok2 {
			name = s
		}
	}
	if v, ok := f.overrides["description"]; ok {
		if s, ok2 := v.(string); ok2 {
			desc = s
		}
	}
	if v, ok := f.overrides["position_kategori_id"]; ok {
		if id, ok2 := v.(int64); ok2 {
			kategoriID = id
		}
	}
	if v, ok := f.overrides["parent_id"]; ok {
		if id, ok2 := v.(*int64); ok2 {
			parentID = id
		} else if id, ok2 := v.(int64); ok2 {
			parentID = &id
		}
	}
	if v, ok := f.overrides["department_id"]; ok {
		if id, ok2 := v.(*int64); ok2 {
			departmentID = id
		} else if id, ok2 := v.(int64); ok2 {
			departmentID = &id
		}
	}
	if v, ok := f.overrides["level_hierarki"]; ok {
		if l, ok2 := v.(int16); ok2 {
			level = l
		}
	}
	if v, ok := f.overrides["kuota"]; ok {
		if k, ok2 := v.(*int16); ok2 {
			kuota = k
		} else if k, ok2 := v.(int16); ok2 {
			kuota = &k
		}
	}
	if v, ok := f.overrides["is_aktif"]; ok {
		if b, ok2 := v.(bool); ok2 {
			isAktif = b
		}
	}

	return &models.Position{
		Name:               name,
		Description:        &desc,
		PositionKategoriID: kategoriID,
		ParentID:           parentID,
		DepartmentID:       departmentID,
		LevelHierarki:      level,
		Kuota:              kuota,
		IsAktif:            isAktif,
	}
}

// MakeMany membuat 'count' Position dari factory 'f' YANG SAMA, supaya
// override yang sudah di-set (mis. position_kategori_id) ikut terpakai
// di semua item. Versi sebelumnya membuat NewPositionFactory() baru di
// tiap iterasi sehingga override dari pemanggil hilang — itu bug, sudah
// diperbaiki di sini.
func (f *PositionFactory) MakeMany(count int) []*models.Position {
	items := make([]*models.Position, count)
	for i := 0; i < count; i++ {
		items[i] = f.Make()
	}
	return items
}
