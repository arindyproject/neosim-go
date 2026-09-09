package dto

import (
	"neosim_go/internal/modules/master/departemen/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// MasterDepartemenResponse response untuk single MasterDepartemen
type MasterDepartemenResponse struct {
	ID          int64            `json:"id"`
	Name        string           `json:"name"`
	FhirCode    *string          `json:"fhir_code"`
	FhirSystem  *string          `json:"fhir_system"`
	Description *string          `json:"description"`
	CreatedBy   *he.UserData     `json:"created_by"`
	UpdatedBy   *he.UserData     `json:"updated_by"`
	CreatedAt   types.CustomTime `json:"created_at"`
	UpdatedAt   types.CustomTime `json:"updated_at"`
}

type MasterDepartemenListSimpelResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type MasterDepartemenResponseParams struct {
	MasterDepartemen *models.MasterDepartemen
	Creator          *he.UserData
	Updater          *he.UserData
}

func ToMasterDepartemenListSimpelResponse(items []models.MasterDepartemen) []MasterDepartemenListSimpelResponse {
	responses := make([]MasterDepartemenListSimpelResponse, 0, len(items))

	for _, m := range items {
		responses = append(responses, MasterDepartemenListSimpelResponse{
			ID:   m.ID,
			Name: m.Name,
		})
	}

	return responses
}

// ToMasterDepartemenResponse mengubah model menjadi response
func ToMasterDepartemenResponse(params MasterDepartemenResponseParams) *MasterDepartemenResponse {
	return &MasterDepartemenResponse{
		ID:          params.MasterDepartemen.ID,
		Name:        params.MasterDepartemen.Name,
		FhirCode:    params.MasterDepartemen.FhirCode,
		FhirSystem:  params.MasterDepartemen.FhirSystem,
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
			Creator:          creator,
			Updater:          updater,
		}))
	}

	return responses
}
