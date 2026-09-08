package factories

import (
	"fmt"
	"math/rand"
	"time"

	"neosim_go/internal/modules/master/departemen/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// MasterDepartemenFactory membuat data MasterDepartemen untuk testing/seeding
type MasterDepartemenFactory struct {
	overrides map[string]interface{}
}

func NewMasterDepartemenFactory() *MasterDepartemenFactory {
	return &MasterDepartemenFactory{overrides: make(map[string]interface{})}
}

func (f *MasterDepartemenFactory) With(field string, value interface{}) *MasterDepartemenFactory {
	f.overrides[field] = value
	return f
}

func (f *MasterDepartemenFactory) Make() *models.MasterDepartemen {
	idx := rng.Intn(999999)
	name := fmt.Sprintf("MasterDepartemen %d", idx)
	fhirCode := fmt.Sprintf("FHIR_CODE_%d", idx)
	fhirSystem := fmt.Sprintf("FHIR_SYSTEM_%d", idx)
	desc := fmt.Sprintf("Deskripsi MasterDepartemen %d", idx)

	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	if v, ok := f.overrides["fhir_code"]; ok {
		fhirCode = v.(string)
	}

	if v, ok := f.overrides["fhir_system"]; ok {
		fhirSystem = v.(string)
	}

	if v, ok := f.overrides["description"]; ok {
		desc = v.(string)
	}

	return &models.MasterDepartemen{
		Name:        name,
		FhirCode:    &fhirCode,
		FhirSystem:  &fhirSystem,
		Description: &desc,
	}
}

func (f *MasterDepartemenFactory) MakeMany(count int) []*models.MasterDepartemen {
	items := make([]*models.MasterDepartemen, count)
	for i := 0; i < count; i++ {
		items[i] = NewMasterDepartemenFactory().Make()
	}
	return items
}
