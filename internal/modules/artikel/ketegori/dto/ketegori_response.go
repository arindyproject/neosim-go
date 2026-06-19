package dto

import (
	"time"

	"neosim_go/internal/modules/artikel/ketegori/models"
)

// ArtikelKetegoriResponse response untuk single ArtikelKetegori
type ArtikelKetegoriResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedBy   *int64    `json:"created_by"`
	UpdatedBy   *int64    `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToArtikelKetegoriResponse mengubah model menjadi response
func ToArtikelKetegoriResponse(m *models.ArtikelKetegori) *ArtikelKetegoriResponse {
	return &ArtikelKetegoriResponse{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedBy:   m.CreatedBy,
		UpdatedBy:   m.UpdatedBy,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// ToArtikelKetegoriListResponse mengubah slice model menjadi slice response
func ToArtikelKetegoriListResponse(items []models.ArtikelKetegori) []ArtikelKetegoriResponse {
	var responses []ArtikelKetegoriResponse
	for _, m := range items {
		responses = append(responses, *ToArtikelKetegoriResponse(&m))
	}
	return responses
}
