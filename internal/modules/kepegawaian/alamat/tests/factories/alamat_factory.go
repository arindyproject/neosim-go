package factories

import (
	"fmt"
	"math/rand"
	"time"

	"neosim_go/internal/modules/kepegawaian/alamat/models"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// KepegawaianAlamatFactory membuat data KepegawaianAlamat untuk testing/seeding
type KepegawaianAlamatFactory struct {
	overrides map[string]interface{}
}

func NewKepegawaianAlamatFactory() *KepegawaianAlamatFactory {
	return &KepegawaianAlamatFactory{overrides: make(map[string]interface{})}
}

func (f *KepegawaianAlamatFactory) With(field string, value interface{}) *KepegawaianAlamatFactory {
	f.overrides[field] = value
	return f
}

func (f *KepegawaianAlamatFactory) Make() *models.KepegawaianAlamat {
	idx := rng.Intn(999999)
	pegawaiID := int64(rng.Intn(10) + 1)
	tipeID := int64(rng.Intn(2) + 1)
	jalan := fmt.Sprintf("Jalan %d", idx)
	rt := fmt.Sprintf("RT %d", idx)
	rw := fmt.Sprintf("RW %d", idx)
	kode_pos := fmt.Sprintf("POS %d", idx)
	desc := fmt.Sprintf("Deskripsi KepegawaianAlamat %d", idx)
	negara_id := int64(1)
	provinsi_id := int64(1)
	kota_kabupaten_id := int64(1)
	kecamatan_id := int64(1)
	kelurahan_desa_id := int64(1)

	if v, ok := f.overrides["jalan"]; ok {
		jalan = v.(string)
	}

	if v, ok := f.overrides["rt"]; ok {
		rt = v.(string)
	}

	if v, ok := f.overrides["rw"]; ok {
		rw = v.(string)
	}

	if v, ok := f.overrides["kode_pos"]; ok {
		kode_pos = v.(string)
	}

	createdBy := int64(rng.Intn(99) + 1)
	updatedBy := int64(rng.Intn(99) + 1)

	return &models.KepegawaianAlamat{
		TipeID:    tipeID,
		PegawaiID: pegawaiID,
		Jalan:     jalan,
		RT:        &rt,
		RW:        &rw,
		KodePos:   &kode_pos,

		NegaraID:        &negara_id,
		ProvinsiID:      &provinsi_id,
		KotaKabupatenID: &kota_kabupaten_id,
		KecamatanID:     &kecamatan_id,
		KelurahanDesaID: &kelurahan_desa_id,

		Description: &desc,
		CreatedBy:   &createdBy,
		UpdatedBy:   &updatedBy,
	}
}

func (f *KepegawaianAlamatFactory) MakeMany(count int) []*models.KepegawaianAlamat {
	items := make([]*models.KepegawaianAlamat, count)
	for i := 0; i < count; i++ {
		items[i] = NewKepegawaianAlamatFactory().Make()
	}
	return items
}
