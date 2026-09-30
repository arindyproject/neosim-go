package dto

import (
	"neosim_go/internal/modules/kepegawaian/pegawai/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// StatusResponse response untuk single Status
type StatusResponse struct {
	ID        int64            `json:"id"`
	Code      string           `json:"code"`
	Label     string           `json:"label"`
	FHIRCode  *string          `json:"fhir_code"`
	CreatedBy *he.UserData     `json:"created_by"`
	UpdatedBy *he.UserData     `json:"updated_by"`
	CreatedAt types.CustomTime `json:"created_at"`
	UpdatedAt types.CustomTime `json:"updated_at"`
}

type StatusSimpelResponse struct {
	ID       int64   `json:"id"`
	Code     string  `json:"code"`
	Label    string  `json:"label"`
	FHIRCode *string `json:"fhir_code"`
}

type StatusSelectResponse struct {
	ID    int64  `json:"id"`
	Code  string `json:"code"`
	Label string `json:"label"`
}

type StatusResponseParams struct {
	Status  *models.Status
	Creator *he.UserData
	Updater *he.UserData
}

func ToStatusSelectResponse(items []models.Status) []StatusSelectResponse {
	responses := make([]StatusSelectResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, StatusSelectResponse{
			ID:    item.ID,
			Code:  item.Code,
			Label: item.Label,
		})
	}
	return responses
}

// ToStatusResponse mengubah model menjadi response
func ToStatusResponse(params StatusResponseParams) *StatusResponse {
	return &StatusResponse{
		ID:        params.Status.ID,
		Code:      params.Status.Code,
		Label:     params.Status.Label,
		FHIRCode:  params.Status.FHIRCode,
		CreatedBy: params.Creator,
		UpdatedBy: params.Updater,
		CreatedAt: types.CustomTime(params.Status.CreatedAt),
		UpdatedAt: types.CustomTime(params.Status.UpdatedAt),
	}
}

// ToStatusListResponse mengubah slice model menjadi slice response
func ToStatusListResponse(
	items []models.Status,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []StatusResponse {
	responses := make([]StatusResponse, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		responses = append(responses, *ToStatusResponse(StatusResponseParams{
			Status:  &m,
			Creator: creator,
			Updater: updater,
		}))
	}

	return responses
}
