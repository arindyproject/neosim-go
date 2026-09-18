package dto

import (
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// PositionResponse response untuk single Position
type PositionResponse struct {
	ID          int64            `json:"id"`
	Name        string           `json:"name"`
	Description *string          `json:"description"`
	CreatedBy   *he.UserData     `json:"created_by"`
	UpdatedBy   *he.UserData     `json:"updated_by"`
	CreatedAt   types.CustomTime `json:"created_at"`
	UpdatedAt   types.CustomTime `json:"updated_at"`
}

type PositionResponseParams struct {
	Position *models.Position
	Creator       *he.UserData
	Updater       *he.UserData
}

// ToPositionResponse mengubah model menjadi response
func ToPositionResponse(params PositionResponseParams) *PositionResponse {
	return &PositionResponse{
		ID:          params.Position.ID,
		Name:        params.Position.Name,
		Description: params.Position.Description,
		CreatedBy:   params.Creator,
		UpdatedBy:   params.Updater,
		CreatedAt:   types.CustomTime(params.Position.CreatedAt),
		UpdatedAt:   types.CustomTime(params.Position.UpdatedAt),
	}
}

// ToPositionListResponse mengubah slice model menjadi slice response
func ToPositionListResponse(
	items []models.Position,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []PositionResponse {
	responses := make([]PositionResponse, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		responses = append(responses, *ToPositionResponse(PositionResponseParams{
			Position: &m,
			Creator:         creator,
			Updater:         updater,
		}))
	}

	return responses
}
