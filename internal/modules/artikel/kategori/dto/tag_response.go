package dto

import (
	"neosim_go/internal/modules/artikel/kategori/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// TagResponse response untuk single Tag
type TagResponse struct {
	ID          int64            `json:"id"`
	Name        string           `json:"name"`
	Description *string          `json:"description"`
	CreatedBy   *he.UserData     `json:"created_by"`
	UpdatedBy   *he.UserData     `json:"updated_by"`
	CreatedAt   types.CustomTime `json:"created_at"`
	UpdatedAt   types.CustomTime `json:"updated_at"`
}

type TagResponseParams struct {
	Tag *models.Tag
	Creator       *he.UserData
	Updater       *he.UserData
}

// ToTagResponse mengubah model menjadi response
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
