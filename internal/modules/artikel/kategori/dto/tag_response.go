package dto

import (
	"neosim_go/internal/shared/types"
	"neosim_go/internal/modules/artikel/kategori/models"
	he "neosim_go/internal/shared/httputil"
)

// TagResponse response untuk single Tag
type TagResponse struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	Description *string      `json:"description"`
	CreatedBy   *he.UserData `json:"created_by"`
	UpdatedBy   *he.UserData `json:"updated_by"`
	CreatedAt   types.CustomTime    `json:"created_at"`
	UpdatedAt   types.CustomTime    `json:"updated_at"`
}

type TagResponseParams struct {
	Tag *models.Tag
	Creator         *he.UserData
	Updater         *he.UserData
}

func ToTagResponse(params TagResponseParams) *TagResponse {
	return &TagResponse{
		ID:          params.Tag.ID,
		Name:        params.Tag.Name,
		Description: params.Tag.Description,
		CreatedBy:   params.Creator,
		UpdatedBy:   params.Updater,
		CreatedAt:   types.CustomTime(params.Tag.CreatedAt),
		UpdatedAt:   types.CustomTime(params.Tag.UpdatedAt),
	}
}

func ToTagListResponse(
	items []models.Tag,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []TagResponse {
	responses := make([]TagResponse, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil {
			creator = creatorsMap[m.ID]
		}
		if updatersMap != nil {
			updater = updatersMap[m.ID]
		}

		responses = append(responses, *ToTagResponse(TagResponseParams{
			Tag: &m,
			Creator:    creator,
			Updater:    updater,
		}))
	}

	return responses
}
