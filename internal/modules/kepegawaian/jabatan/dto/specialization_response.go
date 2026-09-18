package dto

import (
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// SpecializationResponse response untuk single Specialization
type SpecializationResponse struct {
	ID          int64            `json:"id"`
	Name        string           `json:"name"`
	Description *string          `json:"description"`
	CreatedBy   *he.UserData     `json:"created_by"`
	UpdatedBy   *he.UserData     `json:"updated_by"`
	CreatedAt   types.CustomTime `json:"created_at"`
	UpdatedAt   types.CustomTime `json:"updated_at"`
}

type SpecializationResponseParams struct {
	Specialization *models.Specialization
	Creator       *he.UserData
	Updater       *he.UserData
}

// ToSpecializationResponse mengubah model menjadi response
func ToSpecializationResponse(params SpecializationResponseParams) *SpecializationResponse {
	return &SpecializationResponse{
		ID:          params.Specialization.ID,
		Name:        params.Specialization.Name,
		Description: params.Specialization.Description,
		CreatedBy:   params.Creator,
		UpdatedBy:   params.Updater,
		CreatedAt:   types.CustomTime(params.Specialization.CreatedAt),
		UpdatedAt:   types.CustomTime(params.Specialization.UpdatedAt),
	}
}

// ToSpecializationListResponse mengubah slice model menjadi slice response
func ToSpecializationListResponse(
	items []models.Specialization,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []SpecializationResponse {
	responses := make([]SpecializationResponse, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		responses = append(responses, *ToSpecializationResponse(SpecializationResponseParams{
			Specialization: &m,
			Creator:         creator,
			Updater:         updater,
		}))
	}

	return responses
}
