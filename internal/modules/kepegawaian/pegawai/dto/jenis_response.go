package dto

import (
	"neosim_go/internal/modules/kepegawaian/pegawai/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// JenisResponse response untuk single Jenis
type JenisResponse struct {
	ID        int64            `json:"id"`
	Code      string           `json:"code"`
	Label     string           `json:"label"`
	FHIRCode  *string          `json:"fhir_code"`
	CreatedBy *he.UserData     `json:"created_by"`
	UpdatedBy *he.UserData     `json:"updated_by"`
	CreatedAt types.CustomTime `json:"created_at"`
	UpdatedAt types.CustomTime `json:"updated_at"`
}

type JenisSimpelResponse struct {
	ID       int64   `json:"id"`
	Code     string  `json:"code"`
	Label    string  `json:"label"`
	FHIRCode *string `json:"fhir_code"`
}

type JenisSelectResponse struct {
	ID    int64  `json:"id"`
	Code  string `json:"code"`
	Label string `json:"label"`
}

type JenisResponseParams struct {
	Jenis   *models.Jenis
	Creator *he.UserData
	Updater *he.UserData
}

func ToJenisSelectResponse(items []models.Jenis) []JenisSelectResponse {
	responses := make([]JenisSelectResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, JenisSelectResponse{
			ID:    item.ID,
			Code:  item.Code,
			Label: item.Label,
		})
	}
	return responses
}

// ToJenisResponse mengubah model menjadi response
func ToJenisResponse(params JenisResponseParams) *JenisResponse {
	return &JenisResponse{
		ID:        params.Jenis.ID,
		Code:      params.Jenis.Code,
		Label:     params.Jenis.Label,
		FHIRCode:  params.Jenis.FHIRCode,
		CreatedBy: params.Creator,
		UpdatedBy: params.Updater,
		CreatedAt: types.CustomTime(params.Jenis.CreatedAt),
		UpdatedAt: types.CustomTime(params.Jenis.UpdatedAt),
	}
}

// ToJenisListResponse mengubah slice model menjadi slice response
func ToJenisListResponse(
	items []models.Jenis,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []JenisResponse {
	responses := make([]JenisResponse, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		responses = append(responses, *ToJenisResponse(JenisResponseParams{
			Jenis:   &m,
			Creator: creator,
			Updater: updater,
		}))
	}

	return responses
}
