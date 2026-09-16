package dto

import (
	"neosim_go/internal/modules/kepegawaian/alamat/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// Helper functions untuk aman dari nil pointer dereference
func stringVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func int64Val(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

// KepegawaianAlamatResponse response untuk single KepegawaianAlamat
type KepegawaianAlamatResponse struct {
	ID        int64               `json:"id"`
	PegawaiID int64               `json:"pegawai_id"`
	Tipe      *TipeSimpelResponse `json:"tipe,omitempty"`

	Jalan   string `json:"jalan"`
	RT      string `json:"rt"`
	RW      string `json:"rw"`
	KodePos string `json:"kode_pos"`

	// Diubah ke pointer atau tetap int64 (jika int64, nil akan jadi 0)
	//NegaraID        *int64 `json:"negara_id,omitempty"`
	//ProvinsiID      *int64 `json:"provinsi_id,omitempty"`
	//KotaKabupatenID *int64 `json:"kota_kabupaten_id,omitempty"`
	//KecamatanID     *int64 `json:"kecamatan_id,omitempty"`
	//KelurahanDesaID *int64 `json:"kelurahan_desa_id,omitempty"`

	Negara        *WilayahSimpelResponse `json:"negara,omitempty"`
	Provinsi      *WilayahSimpelResponse `json:"provinsi,omitempty"`
	KotaKabupaten *WilayahSimpelResponse `json:"kota_kabupaten,omitempty"`
	Kecamatan     *WilayahSimpelResponse `json:"kecamatan,omitempty"`
	KelurahanDesa *WilayahSimpelResponse `json:"kelurahan_desa,omitempty"`

	IsPrimary bool `json:"is_primary"`

	Description *string          `json:"description"`
	CreatedBy   *he.UserData     `json:"created_by"`
	UpdatedBy   *he.UserData     `json:"updated_by"`
	CreatedAt   types.CustomTime `json:"created_at"`
	UpdatedAt   types.CustomTime `json:"updated_at"`
}

type KepegawaianAlamatResponseParams struct {
	KepegawaianAlamat *models.KepegawaianAlamat
	Creator           *he.UserData
	Updater           *he.UserData

	Negara        *WilayahSimpelResponse
	Provinsi      *WilayahSimpelResponse
	KotaKabupaten *WilayahSimpelResponse
	Kecamatan     *WilayahSimpelResponse
	KelurahanDesa *WilayahSimpelResponse
}

// ToKepegawaianAlamatResponse mengubah model menjadi response
func ToKepegawaianAlamatResponse(params KepegawaianAlamatResponseParams) *KepegawaianAlamatResponse {
	if params.KepegawaianAlamat == nil {
		return nil
	}

	m := params.KepegawaianAlamat

	var tipeResponse *TipeSimpelResponse
	if m.Tipe != nil {
		tipeResponse = &TipeSimpelResponse{
			ID:    m.Tipe.ID,
			Code:  m.Tipe.Code,
			Label: m.Tipe.Label,
		}
	}

	return &KepegawaianAlamatResponse{
		ID:        m.ID,
		PegawaiID: m.PegawaiID,
		Tipe:      tipeResponse,

		Jalan:   m.Jalan,
		RT:      stringVal(m.RT),
		RW:      stringVal(m.RW),
		KodePos: stringVal(m.KodePos),

		// Tetapkan langsung pointer ID wilayah (tanpa dereference *)
		// Agar jika nil di DB, JSON yang keluar bernilai null
		//NegaraID:        m.NegaraID,
		//ProvinsiID:      m.ProvinsiID,
		//KotaKabupatenID: m.KotaKabupatenID,
		//KecamatanID:     m.KecamatanID,
		//KelurahanDesaID: m.KelurahanDesaID,
		Negara:        params.Negara,
		Provinsi:      params.Provinsi,
		KotaKabupaten: params.KotaKabupaten,
		Kecamatan:     params.Kecamatan,
		KelurahanDesa: params.KelurahanDesa,

		IsPrimary:   m.IsPrimary,
		Description: m.Description,
		CreatedBy:   params.Creator,
		UpdatedBy:   params.Updater,
		CreatedAt:   types.CustomTime(m.CreatedAt),
		UpdatedAt:   types.CustomTime(m.UpdatedAt),
	}
}

// ToKepegawaianAlamatListResponse mengubah slice model menjadi slice response
func ToKepegawaianAlamatListResponse(
	items []models.KepegawaianAlamat,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
	negaraMap, provinsiMap, kotaMap, kecamatanMap, kelurahanMap map[int64]*WilayahSimpelResponse,
) []KepegawaianAlamatResponse {
	responses := make([]KepegawaianAlamatResponse, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData
		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		var negara, provinsi, kota, kecamatan, kelurahan *WilayahSimpelResponse
		if m.NegaraID != nil {
			negara = negaraMap[*m.NegaraID]
		}
		if m.ProvinsiID != nil {
			provinsi = provinsiMap[*m.ProvinsiID]
		}
		if m.KotaKabupatenID != nil {
			kota = kotaMap[*m.KotaKabupatenID]
		}
		if m.KecamatanID != nil {
			kecamatan = kecamatanMap[*m.KecamatanID]
		}
		if m.KelurahanDesaID != nil {
			kelurahan = kelurahanMap[*m.KelurahanDesaID]
		}

		res := ToKepegawaianAlamatResponse(KepegawaianAlamatResponseParams{
			KepegawaianAlamat: &m,
			Creator:           creator,
			Updater:           updater,
			Negara:            negara,
			Provinsi:          provinsi,
			KotaKabupaten:     kota,
			Kecamatan:         kecamatan,
			KelurahanDesa:     kelurahan,
		})
		if res != nil {
			responses = append(responses, *res)
		}
	}

	return responses
}

// WilayahSimpelResponse representasi ringkas satu level wilayah
// (Negara, Provinsi, Kota/Kabupaten, Kecamatan, Kelurahan/Desa).
type WilayahSimpelResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
