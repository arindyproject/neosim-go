package dto

import (
	"time"

	"neosim_go/internal/modules/kepegawaian/lampiran/models"
)

// KepegawaianLampiranResponse response untuk single KepegawaianLampiran
type KepegawaianLampiranResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedBy   *int64    `json:"created_by"`
	UpdatedBy   *int64    `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToKepegawaianLampiranResponse mengubah model menjadi response
func ToKepegawaianLampiranResponse(m *models.KepegawaianLampiran) *KepegawaianLampiranResponse {
	return &KepegawaianLampiranResponse{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedBy:   m.CreatedBy,
		UpdatedBy:   m.UpdatedBy,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// ToKepegawaianLampiranListResponse mengubah slice model menjadi slice response
func ToKepegawaianLampiranListResponse(items []models.KepegawaianLampiran) []KepegawaianLampiranResponse {
	var responses []KepegawaianLampiranResponse
	for _, m := range items {
		responses = append(responses, *ToKepegawaianLampiranResponse(&m))
	}
	return responses
}
