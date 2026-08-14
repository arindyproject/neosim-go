package factories

import (
	"fmt"

	"neosim_go/internal/modules/kepegawaian/identifier/models"
)

// TipeFactory membuat data Tipe untuk testing/seeding.
type TipeFactory struct {
	overrides map[string]interface{}
}

func NewTipeFactory() *TipeFactory {
	return &TipeFactory{overrides: make(map[string]interface{})}
}

func (f *TipeFactory) With(field string, value interface{}) *TipeFactory {
	f.overrides[field] = value
	return f
}

func (f *TipeFactory) Make() *models.Tipe {
	idx := rng.Intn(999999)
	code := fmt.Sprintf("CODE_%d", idx)
	label := fmt.Sprintf("Label Tipe %d", idx)
	penerbit := fmt.Sprintf("Penerbit %d", idx)
	fhirSystem := fmt.Sprintf("https://fhir.kemkes.go.id/id/system-%d", idx)
	desc := fmt.Sprintf("Deskripsi Tipe %d", idx)

	item := &models.Tipe{
		Code:        code,
		Label:       label,
		Penerbit:    &penerbit,
		FHIRSystem:  &fhirSystem,
		HasExpiry:   false,
		IsNakes:     false,
		IsRequired:  false,
		Description: &desc,
	}

	// Apply overrides
	if v, ok := f.overrides["code"]; ok {
		item.Code = v.(string)
	}
	if v, ok := f.overrides["label"]; ok {
		item.Label = v.(string)
	}
	if v, ok := f.overrides["penerbit"]; ok {
		if valStr, isStr := v.(string); isStr {
			item.Penerbit = &valStr
		} else if valPtr, isPtr := v.(*string); isPtr {
			item.Penerbit = valPtr
		}
	}
	if v, ok := f.overrides["fhir_system"]; ok {
		if valStr, isStr := v.(string); isStr {
			item.FHIRSystem = &valStr
		} else if valPtr, isPtr := v.(*string); isPtr {
			item.FHIRSystem = valPtr
		}
	}
	if v, ok := f.overrides["has_expiry"]; ok {
		item.HasExpiry = v.(bool)
	}
	if v, ok := f.overrides["is_nakes"]; ok {
		item.IsNakes = v.(bool)
	}
	if v, ok := f.overrides["is_required"]; ok {
		item.IsRequired = v.(bool)
	}
	if v, ok := f.overrides["description"]; ok {
		if valStr, isStr := v.(string); isStr {
			item.Description = &valStr
		} else if valPtr, isPtr := v.(*string); isPtr {
			item.Description = valPtr
		}
	}

	return item
}

func (f *TipeFactory) MakeMany(count int) []*models.Tipe {
	items := make([]*models.Tipe, count)
	for i := 0; i < count; i++ {
		items[i] = NewTipeFactory().Make()
	}
	return items
}
