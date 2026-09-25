package dto

import (
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// SpecializationKategoriResponse response untuk single SpecializationKategori
type SpecializationKategoriResponse struct {
	ID        int64            `json:"id"`
	Code      string           `json:"code"`
	Label     string           `json:"label"`
	FHIRCode  *string          `json:"fhir_code"`
	CreatedBy *he.UserData     `json:"created_by"`
	UpdatedBy *he.UserData     `json:"updated_by"`
	CreatedAt types.CustomTime `json:"created_at"`
	UpdatedAt types.CustomTime `json:"updated_at"`
}

type SpecializationKategoriSimpelResponse struct {
	ID    int64  `json:"id"`
	Code  string `json:"code"`
	Label string `json:"label"`
}

type SpecializationKategoriResponseParams struct {
	SpecializationKategori *models.SpecializationKategori
	Creator                *he.UserData
	Updater                *he.UserData
}

func ToSpecializationKategoriSimpelResponse(items []models.SpecializationKategori) []SpecializationKategoriSimpelResponse {
	responses := make([]SpecializationKategoriSimpelResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, SpecializationKategoriSimpelResponse{
			ID:    item.ID,
			Code:  item.Code,
			Label: item.Label,
		})
	}
	return responses
}

// ToSpecializationKategoriResponse mengubah model menjadi response
func ToSpecializationKategoriResponse(params SpecializationKategoriResponseParams) *SpecializationKategoriResponse {
	return &SpecializationKategoriResponse{
		ID:        params.SpecializationKategori.ID,
		Code:      params.SpecializationKategori.Code,
		Label:     params.SpecializationKategori.Label,
		FHIRCode:  params.SpecializationKategori.FHIRCode,
		CreatedBy: params.Creator,
		UpdatedBy: params.Updater,
		CreatedAt: types.CustomTime(params.SpecializationKategori.CreatedAt),
		UpdatedAt: types.CustomTime(params.SpecializationKategori.UpdatedAt),
	}
}

// ToSpecializationKategoriListResponse mengubah slice model menjadi slice response
func ToSpecializationKategoriListResponse(
	items []models.SpecializationKategori,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []SpecializationKategoriResponse {
	responses := make([]SpecializationKategoriResponse, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		responses = append(responses, *ToSpecializationKategoriResponse(SpecializationKategoriResponseParams{
			SpecializationKategori: &m,
			Creator:                creator,
			Updater:                updater,
		}))
	}

	return responses
}
