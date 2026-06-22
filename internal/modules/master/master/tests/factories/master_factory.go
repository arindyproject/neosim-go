package factories

import (
	"fmt"
	"math/rand"
	"time"

	"neosim_go/internal/modules/master/master/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// =====================================================================
// Pekerjaan
// =====================================================================

type PekerjaanFactory struct {
	overrides map[string]interface{}
}

func NewPekerjaanFactory() *PekerjaanFactory {
	return &PekerjaanFactory{overrides: make(map[string]interface{})}
}

func (f *PekerjaanFactory) With(field string, value interface{}) *PekerjaanFactory {
	f.overrides[field] = value
	return f
}

func (f *PekerjaanFactory) Make() *models.MasterPekerjaan {
	idx := rng.Intn(999999)
	kodeKemenkes := fmt.Sprintf("N%d", idx%100)
	name := fmt.Sprintf("Negara %d", idx)
	desc := fmt.Sprintf("Deskripsi Negara %d", idx)

	if v, ok := f.overrides["kode_kemenkes"]; ok {
		kodeKemenkes = v.(string)
	}
	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.MasterPekerjaan{
		KodeKemenkes: &kodeKemenkes,
		Name:         name,
		Description:  &desc,
	}
}

func (f *PekerjaanFactory) MakeMany(count int) []*models.MasterPekerjaan {
	items := make([]*models.MasterPekerjaan, count)
	for i := 0; i < count; i++ {
		items[i] = NewPekerjaanFactory().Make()
	}
	return items
}

// =====================================================================
// Pendidikan
// =====================================================================

type PendidikanFactory struct {
	overrides map[string]interface{}
}

func NewPendidikanFactory() *PendidikanFactory {
	return &PendidikanFactory{overrides: make(map[string]interface{})}
}

func (f *PendidikanFactory) With(field string, value interface{}) *PendidikanFactory {
	f.overrides[field] = value
	return f
}

func (f *PendidikanFactory) Make() *models.MasterPendidikan {
	idx := rng.Intn(999999)
	kodeKemenkes := fmt.Sprintf("N%d", idx%100)
	name := fmt.Sprintf("Negara %d", idx)
	desc := fmt.Sprintf("Deskripsi Negara %d", idx)

	if v, ok := f.overrides["kode_kemenkes"]; ok {
		kodeKemenkes = v.(string)
	}
	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.MasterPendidikan{
		KodeKemenkes: &kodeKemenkes,
		Name:         name,
		Description:  &desc,
	}
}

func (f *PendidikanFactory) MakeMany(count int) []*models.MasterPendidikan {
	items := make([]*models.MasterPendidikan, count)
	for i := 0; i < count; i++ {
		items[i] = NewPendidikanFactory().Make()
	}
	return items
}

// =====================================================================
// Agama
// =====================================================================

type AgamaFactory struct {
	overrides map[string]interface{}
}

func NewAgamaFactory() *AgamaFactory {
	return &AgamaFactory{overrides: make(map[string]interface{})}
}

func (f *AgamaFactory) With(field string, value interface{}) *AgamaFactory {
	f.overrides[field] = value
	return f
}

func (f *AgamaFactory) Make() *models.MasterAgama {
	idx := rng.Intn(999999)
	kodeKemenkes := fmt.Sprintf("N%d", idx%100)
	name := fmt.Sprintf("Negara %d", idx)
	desc := fmt.Sprintf("Deskripsi Negara %d", idx)

	if v, ok := f.overrides["kode_kemenkes"]; ok {
		kodeKemenkes = v.(string)
	}
	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.MasterAgama{
		KodeKemenkes: &kodeKemenkes,
		Name:         name,
		Description:  &desc,
	}
}

func (f *AgamaFactory) MakeMany(count int) []*models.MasterAgama {
	items := make([]*models.MasterAgama, count)
	for i := 0; i < count; i++ {
		items[i] = NewAgamaFactory().Make()
	}
	return items
}

// =====================================================================
// Status Pernikahan
// =====================================================================

type StatusPernikahanFactory struct {
	overrides map[string]interface{}
}

func NewStatusPernikahanFactory() *StatusPernikahanFactory {
	return &StatusPernikahanFactory{overrides: make(map[string]interface{})}
}

func (f *StatusPernikahanFactory) With(field string, value interface{}) *StatusPernikahanFactory {
	f.overrides[field] = value
	return f
}

func (f *StatusPernikahanFactory) Make() *models.MasterStatusPernikahan {
	idx := rng.Intn(999999)
	kodeKemenkes := fmt.Sprintf("N%d", idx%100)
	name := fmt.Sprintf("Negara %d", idx)
	desc := fmt.Sprintf("Deskripsi Negara %d", idx)

	if v, ok := f.overrides["kode_kemenkes"]; ok {
		kodeKemenkes = v.(string)
	}
	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.MasterStatusPernikahan{
		KodeKemenkes: &kodeKemenkes,
		Name:         name,
		Description:  &desc,
	}
}

func (f *StatusPernikahanFactory) MakeMany(count int) []*models.MasterStatusPernikahan {
	items := make([]*models.MasterStatusPernikahan, count)
	for i := 0; i < count; i++ {
		items[i] = NewStatusPernikahanFactory().Make()
	}
	return items
}
