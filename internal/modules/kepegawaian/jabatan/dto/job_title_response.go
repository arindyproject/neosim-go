package dto

import (
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// JobTitleResponse response untuk single JobTitle
type JobTitleResponse struct {
	ID          int64            `json:"id"`
	Name        string           `json:"name"`
	Description *string          `json:"description"`
	CreatedBy   *he.UserData     `json:"created_by"`
	UpdatedBy   *he.UserData     `json:"updated_by"`
	CreatedAt   types.CustomTime `json:"created_at"`
	UpdatedAt   types.CustomTime `json:"updated_at"`
}

type JobTitleResponseParams struct {
	JobTitle *models.JobTitle
	Creator       *he.UserData
	Updater       *he.UserData
}

// ToJobTitleResponse mengubah model menjadi response
func ToJobTitleResponse(params JobTitleResponseParams) *JobTitleResponse {
	return &JobTitleResponse{
		ID:          params.JobTitle.ID,
		Name:        params.JobTitle.Name,
		Description: params.JobTitle.Description,
		CreatedBy:   params.Creator,
		UpdatedBy:   params.Updater,
		CreatedAt:   types.CustomTime(params.JobTitle.CreatedAt),
		UpdatedAt:   types.CustomTime(params.JobTitle.UpdatedAt),
	}
}

// ToJobTitleListResponse mengubah slice model menjadi slice response
func ToJobTitleListResponse(
	items []models.JobTitle,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []JobTitleResponse {
	responses := make([]JobTitleResponse, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		responses = append(responses, *ToJobTitleResponse(JobTitleResponseParams{
			JobTitle: &m,
			Creator:         creator,
			Updater:         updater,
		}))
	}

	return responses
}
