package dto

import (

	"neosim_go/internal/modules/artikel/artikel/models"
	"neosim_go/internal/shared/types"
)

// ArtikelResponse response untuk single Artikel
type ArtikelResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedBy   *int64    `json:"created_by"`
	UpdatedBy   *int64    `json:"updated_by"`
	CreatedAt   types.CustomTime `json:"created_at"`
	UpdatedAt   types.CustomTime `json:"updated_at"`
}

// ToArtikelResponse mengubah model menjadi response
func ToArtikelResponse(m *models.Artikel) *ArtikelResponse {
	return &ArtikelResponse{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedBy:   m.CreatedBy,
		UpdatedBy:   m.UpdatedBy,
		CreatedAt:   types.CustomTime(m.CreatedAt),
		UpdatedAt:   types.CustomTime(m.UpdatedAt),
	}
}

// ToArtikelListResponse mengubah slice model menjadi slice response
func ToArtikelListResponse(items []models.Artikel) []ArtikelResponse {
	var responses []ArtikelResponse
	for _, m := range items {
		responses = append(responses, *ToArtikelResponse(&m))
	}
	return responses
}
