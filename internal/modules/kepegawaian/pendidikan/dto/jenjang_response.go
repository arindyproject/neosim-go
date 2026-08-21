package dto

import (
	"neosim_go/internal/modules/kepegawaian/pendidikan/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// JenjangResponse response untuk single Jenjang
type JenjangResponse struct {
	ID         int64            `json:"id"`
	Code       string           `json:"code"`
	Label      string           `json:"label"`
	FHIRSystem *string          `json:"fhir_system"`
	CreatedBy  *he.UserData     `json:"created_by"`
	UpdatedBy  *he.UserData     `json:"updated_by"`
	CreatedAt  types.CustomTime `json:"created_at"`
	UpdatedAt  types.CustomTime `json:"updated_at"`
}

type JenjangSimpelResponse struct {
	ID         int64   `json:"id"`
	Code       string  `json:"code"`
	Label      string  `json:"label"`
	FHIRSystem *string `json:"fhir_system"`
}

type JenjangResponseParams struct {
	Jenjang *models.Jenjang
	Creator *he.UserData
	Updater *he.UserData
}

// ToJenjangResponse mengubah model menjadi response
func ToJenjangResponse(params JenjangResponseParams) *JenjangResponse {
	return &JenjangResponse{
		ID:         params.Jenjang.ID,
		Code:       params.Jenjang.Code,
		Label:      params.Jenjang.Label,
		FHIRSystem: params.Jenjang.FHIRSystem,
		CreatedBy:  params.Creator,
		UpdatedBy:  params.Updater,
		CreatedAt:  types.CustomTime(params.Jenjang.CreatedAt),
		UpdatedAt:  types.CustomTime(params.Jenjang.UpdatedAt),
	}
}

// ToJenjangListResponse mengubah slice model menjadi slice response
func ToJenjangListResponse(
	items []models.Jenjang,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []JenjangResponse {
	responses := make([]JenjangResponse, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		responses = append(responses, *ToJenjangResponse(JenjangResponseParams{
			Jenjang: &m,
			Creator: creator,
			Updater: updater,
		}))
	}

	return responses
}
