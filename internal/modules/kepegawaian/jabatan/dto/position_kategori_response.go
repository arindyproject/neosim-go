package dto

import (
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// PositionKategoriResponse response untuk single PositionKategori
type PositionKategoriResponse struct {
	ID       int64    `json:"id"`
	Code     string   `json:"code"`
	Label    string   `json:"label"`
	FHIRCode *string  `json:"fhir_code"`
	Point    *float64 `json:"point"`

	CreatedBy *he.UserData     `json:"created_by"`
	UpdatedBy *he.UserData     `json:"updated_by"`
	CreatedAt types.CustomTime `json:"created_at"`
	UpdatedAt types.CustomTime `json:"updated_at"`
}

type PositionKategoriSimpelResponse struct {
	ID       int64    `json:"id"`
	Code     string   `json:"code"`
	Label    string   `json:"label"`
	FHIRCode *string  `json:"fhir_code"`
	Point    *float64 `json:"point"`
}

type PositionKategoriSelectResponse struct {
	ID    int64  `json:"id"`
	Code  string `json:"code"`
	Label string `json:"label"`
}

type PositionKategoriResponseParams struct {
	PositionKategori *models.PositionKategori
	Creator          *he.UserData
	Updater          *he.UserData
}

func ToPositionKategoriSelectResponse(items []models.PositionKategori) []PositionKategoriSelectResponse {
	responses := make([]PositionKategoriSelectResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, PositionKategoriSelectResponse{
			ID:    item.ID,
			Code:  item.Code,
			Label: item.Label,
		})
	}
	return responses
}

// ToPositionKategoriResponse mengubah model menjadi response
func ToPositionKategoriResponse(params PositionKategoriResponseParams) *PositionKategoriResponse {
	return &PositionKategoriResponse{
		ID:        params.PositionKategori.ID,
		Code:      params.PositionKategori.Code,
		Label:     params.PositionKategori.Label,
		FHIRCode:  params.PositionKategori.FHIRCode,
		Point:     params.PositionKategori.Point,
		CreatedBy: params.Creator,
		UpdatedBy: params.Updater,
		CreatedAt: types.CustomTime(params.PositionKategori.CreatedAt),
		UpdatedAt: types.CustomTime(params.PositionKategori.UpdatedAt),
	}
}

// ToPositionKategoriListResponse mengubah slice model menjadi slice response
func ToPositionKategoriListResponse(
	items []models.PositionKategori,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []PositionKategoriResponse {
	responses := make([]PositionKategoriResponse, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		responses = append(responses, *ToPositionKategoriResponse(PositionKategoriResponseParams{
			PositionKategori: &m,
			Creator:          creator,
			Updater:          updater,
		}))
	}

	return responses
}
