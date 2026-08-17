package dto

import (
	"neosim_go/internal/modules/artikel/artikel/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// ArtikelResponse response untuk single Artikel
type ArtikelResponse struct {
	ID          int64            `json:"id"`
	Name        string           `json:"name"`
	Description *string          `json:"description"`
	CreatedBy   *he.UserData     `json:"created_by"`
	UpdatedBy   *he.UserData     `json:"updated_by"`
	CreatedAt   types.CustomTime `json:"created_at"`
	UpdatedAt   types.CustomTime `json:"updated_at"`
}

type ArtikelResponseParams struct {
	Artikel *models.Artikel
	Creator *he.UserData
	Updater *he.UserData
}

// ToArtikelResponse mengubah model menjadi response
func ToArtikelResponse(params ArtikelResponseParams) *ArtikelResponse {
	return &ArtikelResponse{
		ID:          params.Artikel.ID,
		Name:        params.Artikel.Name,
		Description: params.Artikel.Description,
		CreatedBy:   params.Creator,
		UpdatedBy:   params.Updater,
		CreatedAt:   types.CustomTime(params.Artikel.CreatedAt),
		UpdatedAt:   types.CustomTime(params.Artikel.UpdatedAt),
	}
}

// ToArtikelListResponse mengubah slice model menjadi slice response
func ToArtikelListResponse(
	items []models.Artikel,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []ArtikelResponse {
	responses := make([]ArtikelResponse, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		responses = append(responses, *ToArtikelResponse(ArtikelResponseParams{
			Artikel: &m,
			Creator: creator,
			Updater: updater,
		}))
	}

	return responses
}
