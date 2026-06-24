package factories

import (
	"fmt"
	"math/rand"
	"time"

	"neosim_go/internal/modules/kepegawaian/identifier/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// KepegawaianIdentifierFactory membuat data KepegawaianIdentifier untuk testing/seeding
type KepegawaianIdentifierFactory struct {
	overrides map[string]interface{}
}

func NewKepegawaianIdentifierFactory() *KepegawaianIdentifierFactory {
	return &KepegawaianIdentifierFactory{overrides: make(map[string]interface{})}
}

func (f *KepegawaianIdentifierFactory) With(field string, value interface{}) *KepegawaianIdentifierFactory {
	f.overrides[field] = value
	return f
}

func (f *KepegawaianIdentifierFactory) Make() *models.KepegawaianIdentifier {
	idx := rng.Intn(999999)
	name := fmt.Sprintf("KepegawaianIdentifier %d", idx)
	desc := fmt.Sprintf("Deskripsi KepegawaianIdentifier %d", idx)

	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.KepegawaianIdentifier{
		Name:        name,
		Description: &desc,
	}
}

func (f *KepegawaianIdentifierFactory) MakeMany(count int) []*models.KepegawaianIdentifier {
	items := make([]*models.KepegawaianIdentifier, count)
	for i := 0; i < count; i++ {
		items[i] = NewKepegawaianIdentifierFactory().Make()
	}
	return items
}
