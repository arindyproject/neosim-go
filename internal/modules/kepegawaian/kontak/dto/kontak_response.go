package dto

import (
	"neosim_go/internal/modules/kepegawaian/kontak/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// KepegawaianKontakResponse response untuk single KepegawaianKontak
type KepegawaianKontakResponse struct {
	ID        int64 `json:"id"`
	PegawaiID int64 `json:"pegawai_id"`
	//TipeID      int64            `json:"tipe_id"`
	Tipe        *TipeSimpelResponse `json:"tipe,omitempty"`
	Nilai       string              `json:"nilai"`
	IsPrimary   bool                `json:"is_primary"`
	IsAktif     bool                `json:"is_aktif"`
	Description *string             `json:"description"`
	CreatedBy   *he.UserData        `json:"created_by"`
	UpdatedBy   *he.UserData        `json:"updated_by"`
	CreatedAt   types.CustomTime    `json:"created_at"`
	UpdatedAt   types.CustomTime    `json:"updated_at"`
}

type KepegawaianKontakResponseParams struct {
	KepegawaianKontak *models.KepegawaianKontak
	Creator           *he.UserData
	Updater           *he.UserData
}

// ToKepegawaianKontakResponse mengubah model menjadi response
func ToKepegawaianKontakResponse(params KepegawaianKontakResponseParams) *KepegawaianKontakResponse {
	if params.KepegawaianKontak == nil {
		return nil
	}

	m := params.KepegawaianKontak

	var tipeResponse *TipeSimpelResponse

	if m.Tipe != nil {
		tipeResponse = &TipeSimpelResponse{
			ID:    m.Tipe.ID,
			Code:  m.Tipe.Code,
			Label: m.Tipe.Label,
		}
	}

	return &KepegawaianKontakResponse{
		ID:        m.ID,
		PegawaiID: m.PegawaiID,
		//TipeID:      params.KepegawaianKontak.TipeID,
		Tipe:        tipeResponse,
		Nilai:       m.Nilai,
		IsPrimary:   m.IsPrimary,
		IsAktif:     m.IsAktif,
		Description: m.Description,
		CreatedBy:   params.Creator,
		UpdatedBy:   params.Updater,
		CreatedAt:   types.CustomTime(params.KepegawaianKontak.CreatedAt),
		UpdatedAt:   types.CustomTime(params.KepegawaianKontak.UpdatedAt),
	}
}

// ToKepegawaianKontakListResponse mengubah slice model menjadi slice response
func ToKepegawaianKontakListResponse(
	items []models.KepegawaianKontak,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []KepegawaianKontakResponse {
	responses := make([]KepegawaianKontakResponse, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		responses = append(responses, *ToKepegawaianKontakResponse(KepegawaianKontakResponseParams{
			KepegawaianKontak: &m,
			Creator:           creator,
			Updater:           updater,
		}))
	}

	return responses
}
