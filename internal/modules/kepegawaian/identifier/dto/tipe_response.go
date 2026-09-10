package dto

import (
	"neosim_go/internal/modules/kepegawaian/identifier/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// TipeResponse response untuk single Tipe
type TipeResponse struct {
	ID          int64            `json:"id"`
	Code        string           `json:"code"`
	Label       string           `json:"label"`
	Penerbit    *string          `json:"penerbit"`
	FHIRSystem  *string          `json:"fhir_system"`
	HasExpiry   bool             `json:"has_expiry"`
	IsNakes     bool             `json:"is_nakes"`
	IsRequired  bool             `json:"is_required"`
	Point       *float64         `json:"point"`
	Description *string          `json:"description"`
	CreatedBy   *he.UserData     `json:"created_by"`
	UpdatedBy   *he.UserData     `json:"updated_by"`
	CreatedAt   types.CustomTime `json:"created_at"`
	UpdatedAt   types.CustomTime `json:"updated_at"`
}

type TipeSimpelResponse struct {
	ID          int64    `json:"id"`
	Code        string   `json:"code"`
	Label       string   `json:"label"`
	Penerbit    *string  `json:"penerbit"`
	FHIRSystem  *string  `json:"fhir_system"`
	HasExpiry   bool     `json:"has_expiry"`
	IsNakes     bool     `json:"is_nakes"`
	IsRequired  bool     `json:"is_required"`
	Point       *float64 `json:"point"`
	Description *string  `json:"description"`
}

type TipeSelectResponse struct {
	ID    int64  `json:"id"`
	Code  string `json:"code"`
	Label string `json:"label"`
}

type TipeResponseParams struct {
	Tipe    *models.Tipe
	Creator *he.UserData
	Updater *he.UserData
}

func ToTipeSelectResponse(items []models.Tipe) []TipeSelectResponse {
	responses := make([]TipeSelectResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, TipeSelectResponse{
			ID:    item.ID,
			Code:  item.Code,
			Label: item.Label,
		})
	}
	return responses
}

// ToTipeResponse mengubah model menjadi response
func ToTipeResponse(params TipeResponseParams) *TipeResponse {
	if params.Tipe == nil {
		return nil
	}

	return &TipeResponse{
		ID:          params.Tipe.ID,
		Code:        params.Tipe.Code,
		Label:       params.Tipe.Label,
		Penerbit:    params.Tipe.Penerbit,
		FHIRSystem:  params.Tipe.FHIRSystem,
		HasExpiry:   params.Tipe.HasExpiry,
		IsNakes:     params.Tipe.IsNakes,
		IsRequired:  params.Tipe.IsRequired,
		Point:       params.Tipe.Point,
		Description: params.Tipe.Description,
		CreatedBy:   params.Creator,
		UpdatedBy:   params.Updater,
		CreatedAt:   types.CustomTime(params.Tipe.CreatedAt),
		UpdatedAt:   types.CustomTime(params.Tipe.UpdatedAt),
	}
}

func ToTipeListResponses(
	items []models.Tipe,
	creatorsMap, updatersMap map[int64]*he.UserData,
) []TipeResponse {
	responses := make([]TipeResponse, 0, len(items))
	for _, item := range items {
		var creator, updater *he.UserData
		if creatorsMap != nil && item.CreatedBy != nil {
			creator = creatorsMap[*item.CreatedBy]
		}
		if updatersMap != nil && item.UpdatedBy != nil {
			updater = updatersMap[*item.UpdatedBy]
		}
		responses = append(responses, *ToTipeResponse(TipeResponseParams{
			Tipe:    &item,
			Creator: creator,
			Updater: updater,
		}))
	}
	return responses
}
