package dto

import (
	"neosim_go/internal/modules/kepegawaian/pendidikan/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// KepegawaianPendidikanResponse response untuk single KepegawaianPendidikan
type KepegawaianPendidikanResponse struct {
	ID        int64 `json:"id"`
	PegawaiID int64 `json:"pegawai_id"`
	//JenjangID       int64                  `json:"jenjang_id"`
	Jenjang         *JenjangSimpelResponse `json:"jenjang,omitempty"`
	NamaInstitusi   string                 `json:"nama_institusi"`
	NomorIjazah     *string                `json:"nomor_ijazah"`
	BidangStudi     *string                `json:"bidang_studi"`
	AlamatInstitusi *string                `json:"alamat_institusi"`
	NilaiAkhir      *string                `json:"nilai_akhir"`
	TanggalMasuk    *types.DateOnly        `json:"tanggal_masuk"`
	TanggalLulus    *types.DateOnly        `json:"tanggal_lulus"`
	FHIRCode        *string                `json:"fhir_code"`
	FHIRSystem      *string                `json:"fhir_system"`

	CreatedBy *he.UserData     `json:"created_by"`
	UpdatedBy *he.UserData     `json:"updated_by"`
	CreatedAt types.CustomTime `json:"created_at"`
	UpdatedAt types.CustomTime `json:"updated_at"`
}

type KepegawaianPendidikanResponseParams struct {
	KepegawaianPendidikan *models.KepegawaianPendidikan
	Creator               *he.UserData
	Updater               *he.UserData
}

// ToKepegawaianPendidikanResponse mengubah model menjadi response
func ToKepegawaianPendidikanResponse(params KepegawaianPendidikanResponseParams) *KepegawaianPendidikanResponse {
	if params.KepegawaianPendidikan == nil {
		return nil
	}

	m := params.KepegawaianPendidikan

	var jenjangResponse *JenjangSimpelResponse
	if m.Jenjang != nil {
		jenjangResponse = &JenjangSimpelResponse{
			ID:         m.Jenjang.ID,
			Code:       m.Jenjang.Code,
			Label:      m.Jenjang.Label,
			FHIRSystem: m.Jenjang.FHIRSystem,
		}
	}

	return &KepegawaianPendidikanResponse{
		ID:        params.KepegawaianPendidikan.ID,
		PegawaiID: params.KepegawaianPendidikan.PegawaiID,
		//JenjangID:       params.KepegawaianPendidikan.JenjangID,
		Jenjang:         jenjangResponse,
		NamaInstitusi:   params.KepegawaianPendidikan.NamaInstitusi,
		NomorIjazah:     params.KepegawaianPendidikan.NomorIjazah,
		BidangStudi:     params.KepegawaianPendidikan.BidangStudi,
		AlamatInstitusi: params.KepegawaianPendidikan.AlamatInstitusi,
		NilaiAkhir:      params.KepegawaianPendidikan.NilaiAkhir,
		TanggalMasuk:    types.NewDateOnlyPtr(params.KepegawaianPendidikan.TanggalMasuk),
		TanggalLulus:    types.NewDateOnlyPtr(params.KepegawaianPendidikan.TanggalLulus),
		FHIRCode:        params.KepegawaianPendidikan.FHIRCode,
		FHIRSystem:      params.KepegawaianPendidikan.FHIRSystem,

		CreatedBy: params.Creator,
		UpdatedBy: params.Updater,
		CreatedAt: types.CustomTime(params.KepegawaianPendidikan.CreatedAt),
		UpdatedAt: types.CustomTime(params.KepegawaianPendidikan.UpdatedAt),
	}
}

// ToKepegawaianPendidikanListResponse mengubah slice model menjadi slice response
func ToKepegawaianPendidikanListResponse(
	items []models.KepegawaianPendidikan,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []KepegawaianPendidikanResponse {
	responses := make([]KepegawaianPendidikanResponse, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		responses = append(responses, *ToKepegawaianPendidikanResponse(KepegawaianPendidikanResponseParams{
			KepegawaianPendidikan: &m,
			Creator:               creator,
			Updater:               updater,
		}))
	}

	return responses
}
