package dto

import (

	"neosim_go/internal/modules/kepegawaian/pegawai/models"
	"neosim_go/internal/shared/types"
	he "neosim_go/internal/shared/httputil"
)

// KepegawaianPegawaiResponse response untuk single KepegawaianPegawai
type KepegawaianPegawaiResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedBy   *he.UserData `json:"created_by"`
	UpdatedBy   *he.UserData `json:"updated_by"`
	CreatedAt   types.CustomTime `json:"created_at"`
	UpdatedAt   types.CustomTime `json:"updated_at"`
}

type KepegawaianPegawaiResponseParams struct {
	KepegawaianPegawai *models.KepegawaianPegawai
	Creator         *he.UserData
	Updater         *he.UserData
}

// ToKepegawaianPegawaiResponse mengubah model menjadi response
func ToKepegawaianPegawaiResponse(params KepegawaianPegawaiResponseParams) *KepegawaianPegawaiResponse {
	return &KepegawaianPegawaiResponse{
		ID:          params.KepegawaianPegawai.ID,
		Name:        params.KepegawaianPegawai.Name,
		Description: params.KepegawaianPegawai.Description,
		CreatedBy:   params.Creator,
		UpdatedBy:   params.Updater,
		CreatedAt:   types.CustomTime(params.KepegawaianPegawai.CreatedAt),
		UpdatedAt:   types.CustomTime(params.KepegawaianPegawai.UpdatedAt),
	}
}

// ToKepegawaianPegawaiListResponse mengubah slice model menjadi slice response
func ToKepegawaianPegawaiListResponse(
	items []models.KepegawaianPegawai,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []KepegawaianPegawaiResponse {
	responses := make([]KepegawaianPegawaiResponse, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		responses = append(responses, *ToKepegawaianPegawaiResponse(KepegawaianPegawaiResponseParams{
			KepegawaianPegawai: &m,
			Creator:    creator,
			Updater:    updater,
		}))
	}

	return responses
}
