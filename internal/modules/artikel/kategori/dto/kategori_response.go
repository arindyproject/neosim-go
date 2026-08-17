package dto

import (

	"neosim_go/internal/modules/artikel/kategori/models"
	"neosim_go/internal/shared/types"
	he "neosim_go/internal/shared/httputil"
)

// ArtikelKategoriResponse response untuk single ArtikelKategori
type ArtikelKategoriResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedBy   *he.UserData `json:"created_by"`
	UpdatedBy   *he.UserData `json:"updated_by"`
	CreatedAt   types.CustomTime `json:"created_at"`
	UpdatedAt   types.CustomTime `json:"updated_at"`
}

type ArtikelKategoriResponseParams struct {
	ArtikelKategori *models.ArtikelKategori
	Creator         *he.UserData
	Updater         *he.UserData
}

// ToArtikelKategoriResponse mengubah model menjadi response
func ToArtikelKategoriResponse(params ArtikelKategoriResponseParams) *ArtikelKategoriResponse {
	return &ArtikelKategoriResponse{
		ID:          params.ArtikelKategori.ID,
		Name:        params.ArtikelKategori.Name,
		Description: params.ArtikelKategori.Description,
		CreatedBy:   params.Creator,
		UpdatedBy:   params.Updater,
		CreatedAt:   types.CustomTime(params.ArtikelKategori.CreatedAt),
		UpdatedAt:   types.CustomTime(params.ArtikelKategori.UpdatedAt),
	}
}

// ToArtikelKategoriListResponse mengubah slice model menjadi slice response
func ToArtikelKategoriListResponse(
	items []models.ArtikelKategori,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []ArtikelKategoriResponse {
	responses := make([]ArtikelKategoriResponse, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		responses = append(responses, *ToArtikelKategoriResponse(ArtikelKategoriResponseParams{
			ArtikelKategori: &m,
			Creator:    creator,
			Updater:    updater,
		}))
	}

	return responses
}
