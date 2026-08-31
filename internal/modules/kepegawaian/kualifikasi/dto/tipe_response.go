package dto

import (
	"neosim_go/internal/modules/kepegawaian/kualifikasi/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// TipeResponse response untuk single Tipe
type TipeResponse struct {
	ID        int64            `json:"id"`
	Code      string           `json:"code"`
	Label     string           `json:"label"`
	CreatedBy *he.UserData     `json:"created_by"`
	UpdatedBy *he.UserData     `json:"updated_by"`
	CreatedAt types.CustomTime `json:"created_at"`
	UpdatedAt types.CustomTime `json:"updated_at"`
}

type TipeSimpelResponse struct {
	ID    int64  `json:"id"`
	Code  string `json:"code"`
	Label string `json:"label"`
}

type TipeResponseParams struct {
	Tipe    *models.Tipe
	Creator *he.UserData
	Updater *he.UserData
}

// ToTipeResponse mengubah model menjadi response
func ToTipeResponse(params TipeResponseParams) *TipeResponse {
	return &TipeResponse{
		ID:        params.Tipe.ID,
		Code:      params.Tipe.Code,
		Label:     params.Tipe.Label,
		CreatedBy: params.Creator,
		UpdatedBy: params.Updater,
		CreatedAt: types.CustomTime(params.Tipe.CreatedAt),
		UpdatedAt: types.CustomTime(params.Tipe.UpdatedAt),
	}
}

// ToTipeListResponse mengubah slice model menjadi slice response
func ToTipeListResponse(
	items []models.Tipe,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []TipeResponse {
	responses := make([]TipeResponse, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		responses = append(responses, *ToTipeResponse(TipeResponseParams{
			Tipe:    &m,
			Creator: creator,
			Updater: updater,
		}))
	}

	return responses
}
