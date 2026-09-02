package dto

import (
	"neosim_go/internal/modules/kepegawaian/kualifikasi/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// KepegawaianKualifikasiResponse response untuk single KepegawaianKualifikasi
type KepegawaianKualifikasiResponse struct {
	ID        int64 `json:"id"`
	PegawaiID int64 `json:"pegawai_id"`
	//TipeID      int64            `json:"tipe_id"`
	Tipe *TipeSimpelResponse `json:"tipe,omitempty"`

	Nama          string `json:"nama"`
	Penyelenggara string `json:"penyelenggara"`

	NomorSertifikat *string `json:"nomor_sertifikat"`

	TanggalTerbit  *types.DateOnly `json:"tanggal_terbit"`
	TanggalExpired *types.DateOnly `json:"tanggal_expired"`

	IsAktif      bool `json:"is_aktif"`
	IsExpired    bool `json:"is_expired"`
	DaysUntilExp int  `json:"days_until_expired"`

	FhirCode   *string `json:"fhir_code"`
	FhirSystem *string `json:"fhir_system"`

	CreatedBy *he.UserData     `json:"created_by"`
	UpdatedBy *he.UserData     `json:"updated_by"`
	CreatedAt types.CustomTime `json:"created_at"`
	UpdatedAt types.CustomTime `json:"updated_at"`
}

type KepegawaianKualifikasiResponseParams struct {
	KepegawaianKualifikasi *models.KepegawaianKualifikasi
	Creator                *he.UserData
	Updater                *he.UserData
}

// ToKepegawaianKualifikasiResponse mengubah model menjadi response
func ToKepegawaianKualifikasiResponse(params KepegawaianKualifikasiResponseParams) *KepegawaianKualifikasiResponse {
	if params.KepegawaianKualifikasi == nil {
		return nil
	}

	m := params.KepegawaianKualifikasi
	var tipeResponse *TipeSimpelResponse
	if m.Tipe != nil {
		tipeResponse = &TipeSimpelResponse{
			ID:    m.Tipe.ID,
			Code:  m.Tipe.Code,
			Label: m.Tipe.Label,
		}
	}

	return &KepegawaianKualifikasiResponse{
		ID:        m.ID,
		PegawaiID: m.PegawaiID,
		Tipe:      tipeResponse,

		Nama:            m.Nama,
		Penyelenggara:   m.Penyelenggara,
		NomorSertifikat: m.NomorSertifikat,

		TanggalTerbit:  types.NewDateOnlyPtr(m.TanggalTerbit),
		TanggalExpired: types.NewDateOnlyPtr(m.TanggalExpired),

		IsAktif:      m.IsAktif,
		IsExpired:    m.IsExpired(),
		DaysUntilExp: m.DaysUntilExpired(),

		FhirCode:   m.FhirCode,
		FhirSystem: m.FhirSystem,

		CreatedBy: params.Creator,
		UpdatedBy: params.Updater,
		CreatedAt: types.CustomTime(params.KepegawaianKualifikasi.CreatedAt),
		UpdatedAt: types.CustomTime(params.KepegawaianKualifikasi.UpdatedAt),
	}
}

// ToKepegawaianKualifikasiListResponse mengubah slice model menjadi slice response
func ToKepegawaianKualifikasiListResponse(
	items []models.KepegawaianKualifikasi,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []KepegawaianKualifikasiResponse {
	responses := make([]KepegawaianKualifikasiResponse, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		responses = append(responses, *ToKepegawaianKualifikasiResponse(KepegawaianKualifikasiResponseParams{
			KepegawaianKualifikasi: &m,
			Creator:                creator,
			Updater:                updater,
		}))
	}

	return responses
}
