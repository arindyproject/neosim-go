package dto

import (
	"time"

	"neosim_go/internal/modules/artikel/models"
)

// ArtikelResponse response untuk single artikel
type ArtikelResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedBy   *int64    `json:"created_by"`
	UpdatedBy   *int64    `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToArtikelResponse mengubah model menjadi response
func ToArtikelResponse(m *models.Artikel) *ArtikelResponse {
	return &ArtikelResponse{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedBy:   m.CreatedBy,
		UpdatedBy:   m.UpdatedBy,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
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
