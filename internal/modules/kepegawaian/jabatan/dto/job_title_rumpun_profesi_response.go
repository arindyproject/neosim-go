package dto

import (
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// JobTitleRumpunProfesiResponse response untuk single JobTitleRumpunProfesi
type JobTitleRumpunProfesiResponse struct {
	ID        int64            `json:"id"`
	Code      string           `json:"code"`
	Label     string           `json:"label"`
	FHIRCode  *string          `json:"fhir_code"`
	CreatedBy *he.UserData     `json:"created_by"`
	UpdatedBy *he.UserData     `json:"updated_by"`
	CreatedAt types.CustomTime `json:"created_at"`
	UpdatedAt types.CustomTime `json:"updated_at"`
}

type JobTitleRumpunProfesiResponseParams struct {
	JobTitleRumpunProfesi *models.JobTitleRumpunProfesi
	Creator               *he.UserData
	Updater               *he.UserData
}

type JobTitleRumpunProfesiSimpelResponse struct {
	ID    int64  `json:"id"`
	Code  string `json:"code"`
	Label string `json:"label"`
}

func ToJobTitleRumpunProfesiSimpelResponse(items []models.JobTitleRumpunProfesi) []JobTitleRumpunProfesiSimpelResponse {
	responses := make([]JobTitleRumpunProfesiSimpelResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, JobTitleRumpunProfesiSimpelResponse{
			ID:    item.ID,
			Code:  item.Code,
			Label: item.Label,
		})
	}
	return responses
}

// ToJobTitleRumpunProfesiResponse mengubah model menjadi response
func ToJobTitleRumpunProfesiResponse(params JobTitleRumpunProfesiResponseParams) *JobTitleRumpunProfesiResponse {
	return &JobTitleRumpunProfesiResponse{
		ID:        params.JobTitleRumpunProfesi.ID,
		Code:      params.JobTitleRumpunProfesi.Code,
		Label:     params.JobTitleRumpunProfesi.Label,
		FHIRCode:  params.JobTitleRumpunProfesi.FHIRCode,
		CreatedBy: params.Creator,
		UpdatedBy: params.Updater,
		CreatedAt: types.CustomTime(params.JobTitleRumpunProfesi.CreatedAt),
		UpdatedAt: types.CustomTime(params.JobTitleRumpunProfesi.UpdatedAt),
	}
}

// ToJobTitleRumpunProfesiListResponse mengubah slice model menjadi slice response
func ToJobTitleRumpunProfesiListResponse(
	items []models.JobTitleRumpunProfesi,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []JobTitleRumpunProfesiResponse {
	responses := make([]JobTitleRumpunProfesiResponse, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		responses = append(responses, *ToJobTitleRumpunProfesiResponse(JobTitleRumpunProfesiResponseParams{
			JobTitleRumpunProfesi: &m,
			Creator:               creator,
			Updater:               updater,
		}))
	}

	return responses
}
