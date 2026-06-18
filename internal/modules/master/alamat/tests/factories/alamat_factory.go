package factories

import (
	"fmt"
	"math/rand"
	"time"

	"neosim_go/internal/modules/master/alamat/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// =====================================================================
// NEGARA
// =====================================================================

type NegaraFactory struct {
	overrides map[string]interface{}
}

func NewNegaraFactory() *NegaraFactory {
	return &NegaraFactory{overrides: make(map[string]interface{})}
}

func (f *NegaraFactory) With(field string, value interface{}) *NegaraFactory {
	f.overrides[field] = value
	return f
}

func (f *NegaraFactory) Make() *models.MasterAlamatNegara {
	idx := rng.Intn(999999)
	code := fmt.Sprintf("N%d", idx%100)
	name := fmt.Sprintf("Negara %d", idx)
	desc := fmt.Sprintf("Deskripsi Negara %d", idx)

	if v, ok := f.overrides["code"]; ok {
		code = v.(string)
	}
	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}

	return &models.MasterAlamatNegara{
		Code:        code,
		Name:        name,
		Description: &desc,
	}
}

func (f *NegaraFactory) MakeMany(count int) []*models.MasterAlamatNegara {
	items := make([]*models.MasterAlamatNegara, count)
	for i := 0; i < count; i++ {
		items[i] = NewNegaraFactory().Make()
	}
	return items
}

// =====================================================================
// PROVINSI
// =====================================================================

type ProvinsiFactory struct {
	overrides map[string]interface{}
}

func NewProvinsiFactory() *ProvinsiFactory {
	return &ProvinsiFactory{overrides: make(map[string]interface{})}
}

func (f *ProvinsiFactory) With(field string, value interface{}) *ProvinsiFactory {
	f.overrides[field] = value
	return f
}

func (f *ProvinsiFactory) Make() *models.MasterAlamatProvinsi {
	idx := rng.Intn(999999)
	code := fmt.Sprintf("%d", idx%99)
	name := fmt.Sprintf("Provinsi %d", idx)
	var negaraID int64 = 1

	if v, ok := f.overrides["code"]; ok {
		code = v.(string)
	}
	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}
	if v, ok := f.overrides["negara_id"]; ok {
		negaraID = v.(int64)
	}

	return &models.MasterAlamatProvinsi{
		NegaraID: negaraID,
		Code:     code,
		Name:     name,
	}
}

func (f *ProvinsiFactory) MakeMany(count int) []*models.MasterAlamatProvinsi {
	items := make([]*models.MasterAlamatProvinsi, count)
	for i := 0; i < count; i++ {
		items[i] = NewProvinsiFactory().Make()
	}
	return items
}

// =====================================================================
// KOTA / KABUPATEN
// =====================================================================

type KotaKabupatenFactory struct {
	overrides map[string]interface{}
}

func NewKotaKabupatenFactory() *KotaKabupatenFactory {
	return &KotaKabupatenFactory{overrides: make(map[string]interface{})}
}

func (f *KotaKabupatenFactory) With(field string, value interface{}) *KotaKabupatenFactory {
	f.overrides[field] = value
	return f
}

func (f *KotaKabupatenFactory) Make() *models.MasterAlamatKotaKabupaten {
	idx := rng.Intn(999999)
	code := fmt.Sprintf("35.%d", idx%99)
	name := fmt.Sprintf("Kota %d", idx)
	var provinsiID int64 = 1

	if v, ok := f.overrides["code"]; ok {
		code = v.(string)
	}
	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}
	if v, ok := f.overrides["provinsi_id"]; ok {
		provinsiID = v.(int64)
	}

	return &models.MasterAlamatKotaKabupaten{
		ProvinsiID: provinsiID,
		Code:       code,
		Name:       name,
	}
}

func (f *KotaKabupatenFactory) MakeMany(count int) []*models.MasterAlamatKotaKabupaten {
	items := make([]*models.MasterAlamatKotaKabupaten, count)
	for i := 0; i < count; i++ {
		items[i] = NewKotaKabupatenFactory().Make()
	}
	return items
}

// =====================================================================
// KECAMATAN
// =====================================================================

type KecamatanFactory struct {
	overrides map[string]interface{}
}

func NewKecamatanFactory() *KecamatanFactory {
	return &KecamatanFactory{overrides: make(map[string]interface{})}
}

func (f *KecamatanFactory) With(field string, value interface{}) *KecamatanFactory {
	f.overrides[field] = value
	return f
}

func (f *KecamatanFactory) Make() *models.MasterAlamatKecamatan {
	idx := rng.Intn(999999)
	code := fmt.Sprintf("35.21.%d", idx%99)
	name := fmt.Sprintf("Kecamatan %d", idx)
	var kotaKabupatenID int64 = 1

	if v, ok := f.overrides["code"]; ok {
		code = v.(string)
	}
	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}
	if v, ok := f.overrides["kota_kabupaten_id"]; ok {
		kotaKabupatenID = v.(int64)
	}

	return &models.MasterAlamatKecamatan{
		KotaKabupatenID: kotaKabupatenID,
		Code:            code,
		Name:            name,
	}
}

func (f *KecamatanFactory) MakeMany(count int) []*models.MasterAlamatKecamatan {
	items := make([]*models.MasterAlamatKecamatan, count)
	for i := 0; i < count; i++ {
		items[i] = NewKecamatanFactory().Make()
	}
	return items
}

// =====================================================================
// KELURAHAN / DESA
// =====================================================================

type KelurahanDesaFactory struct {
	overrides map[string]interface{}
}

func NewKelurahanDesaFactory() *KelurahanDesaFactory {
	return &KelurahanDesaFactory{overrides: make(map[string]interface{})}
}

func (f *KelurahanDesaFactory) With(field string, value interface{}) *KelurahanDesaFactory {
	f.overrides[field] = value
	return f
}

func (f *KelurahanDesaFactory) Make() *models.MasterAlamatKelurahanDesa {
	idx := rng.Intn(999999)
	code := fmt.Sprintf("35.21.01.%d", idx%9999)
	name := fmt.Sprintf("Desa %d", idx)
	postal := fmt.Sprintf("%05d", idx%99999)
	var kecamatanID int64 = 1

	if v, ok := f.overrides["code"]; ok {
		code = v.(string)
	}
	if v, ok := f.overrides["name"]; ok {
		name = v.(string)
	}
	if v, ok := f.overrides["kecamatan_id"]; ok {
		kecamatanID = v.(int64)
	}

	return &models.MasterAlamatKelurahanDesa{
		KecamatanID: kecamatanID,
		Code:        code,
		Name:        name,
		PostalCode:  &postal,
	}
}

func (f *KelurahanDesaFactory) MakeMany(count int) []*models.MasterAlamatKelurahanDesa {
	items := make([]*models.MasterAlamatKelurahanDesa, count)
	for i := 0; i < count; i++ {
		items[i] = NewKelurahanDesaFactory().Make()
	}
	return items
}
