package dto

import (

	"neosim_go/internal/modules/kepegawaian/alamat/models"
	"neosim_go/internal/shared/types"
	he "neosim_go/internal/shared/httputil"
)

// KepegawaianAlamatResponse response untuk single KepegawaianAlamat
type KepegawaianAlamatResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedBy   *he.UserData `json:"created_by"`
	UpdatedBy   *he.UserData `json:"updated_by"`
	CreatedAt   types.CustomTime `json:"created_at"`
	UpdatedAt   types.CustomTime `json:"updated_at"`
}

type KepegawaianAlamatResponseParams struct {
	KepegawaianAlamat *models.KepegawaianAlamat
	Creator         *he.UserData
	Updater         *he.UserData
}

// ToKepegawaianAlamatResponse mengubah model menjadi response
func ToKepegawaianAlamatResponse(params KepegawaianAlamatResponseParams) *KepegawaianAlamatResponse {
	return &KepegawaianAlamatResponse{
		ID:          params.KepegawaianAlamat.ID,
		Name:        params.KepegawaianAlamat.Name,
		Description: params.KepegawaianAlamat.Description,
		CreatedBy:   params.Creator,
		UpdatedBy:   params.Updater,
		CreatedAt:   types.CustomTime(params.KepegawaianAlamat.CreatedAt),
		UpdatedAt:   types.CustomTime(params.KepegawaianAlamat.UpdatedAt),
	}
}

// ToKepegawaianAlamatListResponse mengubah slice model menjadi slice response
func ToKepegawaianAlamatListResponse(
	items []models.KepegawaianAlamat,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
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

		responses = append(responses, *ToKepegawaianAlamatResponse(KepegawaianAlamatResponseParams{
			KepegawaianAlamat: &m,
			Creator:    creator,
			Updater:    updater,
		}))
	}

	return responses
}
