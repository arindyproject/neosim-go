package dto

import (
	"time"

	"neosim_go/internal/modules/kepegawaian/kontak/models"
)

// KepegawaianKontakResponse response untuk single KepegawaianKontak
type KepegawaianKontakResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedBy   *int64    `json:"created_by"`
	UpdatedBy   *int64    `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToKepegawaianKontakResponse mengubah model menjadi response
func ToKepegawaianKontakResponse(m *models.KepegawaianKontak) *KepegawaianKontakResponse {
	return &KepegawaianKontakResponse{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedBy:   m.CreatedBy,
		UpdatedBy:   m.UpdatedBy,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// ToKepegawaianKontakListResponse mengubah slice model menjadi slice response
func ToKepegawaianKontakListResponse(items []models.KepegawaianKontak) []KepegawaianKontakResponse {
	var responses []KepegawaianKontakResponse
	for _, m := range items {
		responses = append(responses, *ToKepegawaianKontakResponse(&m))
	}
	return responses
}
