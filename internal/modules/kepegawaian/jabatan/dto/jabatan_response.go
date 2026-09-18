package dto

import (

	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	"neosim_go/internal/shared/types"
	he "neosim_go/internal/shared/httputil"
)

// KepegawaianJabatanResponse response untuk single KepegawaianJabatan
type KepegawaianJabatanResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedBy   *he.UserData `json:"created_by"`
	UpdatedBy   *he.UserData `json:"updated_by"`
	CreatedAt   types.CustomTime `json:"created_at"`
	UpdatedAt   types.CustomTime `json:"updated_at"`
}

type KepegawaianJabatanResponseParams struct {
	KepegawaianJabatan *models.KepegawaianJabatan
	Creator         *he.UserData
	Updater         *he.UserData
}

// ToKepegawaianJabatanResponse mengubah model menjadi response
func ToKepegawaianJabatanResponse(params KepegawaianJabatanResponseParams) *KepegawaianJabatanResponse {
	return &KepegawaianJabatanResponse{
		ID:          params.KepegawaianJabatan.ID,
		Name:        params.KepegawaianJabatan.Name,
		Description: params.KepegawaianJabatan.Description,
		CreatedBy:   params.Creator,
		UpdatedBy:   params.Updater,
		CreatedAt:   types.CustomTime(params.KepegawaianJabatan.CreatedAt),
		UpdatedAt:   types.CustomTime(params.KepegawaianJabatan.UpdatedAt),
	}
}

// ToKepegawaianJabatanListResponse mengubah slice model menjadi slice response
func ToKepegawaianJabatanListResponse(
	items []models.KepegawaianJabatan,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []KepegawaianJabatanResponse {
	responses := make([]KepegawaianJabatanResponse, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		responses = append(responses, *ToKepegawaianJabatanResponse(KepegawaianJabatanResponseParams{
			KepegawaianJabatan: &m,
			Creator:    creator,
			Updater:    updater,
		}))
	}

	return responses
}
