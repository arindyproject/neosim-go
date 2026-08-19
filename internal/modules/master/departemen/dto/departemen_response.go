package dto

import (

	"neosim_go/internal/modules/master/departemen/models"
	"neosim_go/internal/shared/types"
	he "neosim_go/internal/shared/httputil"
)

// MasterDepartemenResponse response untuk single MasterDepartemen
type MasterDepartemenResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedBy   *he.UserData `json:"created_by"`
	UpdatedBy   *he.UserData `json:"updated_by"`
	CreatedAt   types.CustomTime `json:"created_at"`
	UpdatedAt   types.CustomTime `json:"updated_at"`
}

type MasterDepartemenResponseParams struct {
	MasterDepartemen *models.MasterDepartemen
	Creator         *he.UserData
	Updater         *he.UserData
}

// ToMasterDepartemenResponse mengubah model menjadi response
func ToMasterDepartemenResponse(params MasterDepartemenResponseParams) *MasterDepartemenResponse {
	return &MasterDepartemenResponse{
		ID:          params.MasterDepartemen.ID,
		Name:        params.MasterDepartemen.Name,
		Description: params.MasterDepartemen.Description,
		CreatedBy:   params.Creator,
		UpdatedBy:   params.Updater,
		CreatedAt:   types.CustomTime(params.MasterDepartemen.CreatedAt),
		UpdatedAt:   types.CustomTime(params.MasterDepartemen.UpdatedAt),
	}
}

// ToMasterDepartemenListResponse mengubah slice model menjadi slice response
func ToMasterDepartemenListResponse(
	items []models.MasterDepartemen,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []MasterDepartemenResponse {
	responses := make([]MasterDepartemenResponse, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		responses = append(responses, *ToMasterDepartemenResponse(MasterDepartemenResponseParams{
			MasterDepartemen: &m,
			Creator:    creator,
			Updater:    updater,
		}))
	}

	return responses
}
