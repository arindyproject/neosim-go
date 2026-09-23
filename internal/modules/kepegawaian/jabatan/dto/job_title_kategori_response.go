package dto

import (
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// JobTitleKategoriResponse response untuk single JobTitleKategori
type JobTitleKategoriResponse struct {
	ID        int64            `json:"id"`
	Code      string           `json:"code"`
	Label     string           `json:"label"`
	FHIRCode  *string          `json:"fhir_code"`
	CreatedBy *he.UserData     `json:"created_by"`
	UpdatedBy *he.UserData     `json:"updated_by"`
	CreatedAt types.CustomTime `json:"created_at"`
	UpdatedAt types.CustomTime `json:"updated_at"`
}

type JobTitleKategoriSimpelResponse struct {
	ID    int64  `json:"id"`
	Code  string `json:"code"`
	Label string `json:"label"`
}

type JobTitleKategoriResponseParams struct {
	JobTitleKategori *models.JobTitleKategori
	Creator          *he.UserData
	Updater          *he.UserData
}

func ToJobTitleKategoriSimpelResponse(items []models.JobTitleKategori) []JobTitleKategoriSimpelResponse {
	responses := make([]JobTitleKategoriSimpelResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, JobTitleKategoriSimpelResponse{
			ID:    item.ID,
			Code:  item.Code,
			Label: item.Label,
		})
	}
	return responses
}

// ToJobTitleKategoriResponse mengubah model menjadi response
func ToJobTitleKategoriResponse(params JobTitleKategoriResponseParams) *JobTitleKategoriResponse {
	return &JobTitleKategoriResponse{
		ID:        params.JobTitleKategori.ID,
		Code:      params.JobTitleKategori.Code,
		Label:     params.JobTitleKategori.Label,
		FHIRCode:  params.JobTitleKategori.FHIRCode,
		CreatedBy: params.Creator,
		UpdatedBy: params.Updater,
		CreatedAt: types.CustomTime(params.JobTitleKategori.CreatedAt),
		UpdatedAt: types.CustomTime(params.JobTitleKategori.UpdatedAt),
	}
}

// ToJobTitleKategoriListResponse mengubah slice model menjadi slice response
func ToJobTitleKategoriListResponse(
	items []models.JobTitleKategori,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []JobTitleKategoriResponse {
	responses := make([]JobTitleKategoriResponse, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		responses = append(responses, *ToJobTitleKategoriResponse(JobTitleKategoriResponseParams{
			JobTitleKategori: &m,
			Creator:          creator,
			Updater:          updater,
		}))
	}

	return responses
}
