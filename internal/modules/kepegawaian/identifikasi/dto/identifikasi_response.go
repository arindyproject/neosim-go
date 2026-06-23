package dto

import (
	"time"

	"neosim_go/internal/modules/kepegawaian/identifikasi/models"
)

// KepegawaianIdentifikasiResponse response untuk single KepegawaianIdentifikasi
type KepegawaianIdentifikasiResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedBy   *int64    `json:"created_by"`
	UpdatedBy   *int64    `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToKepegawaianIdentifikasiResponse mengubah model menjadi response
func ToKepegawaianIdentifikasiResponse(m *models.KepegawaianIdentifikasi) *KepegawaianIdentifikasiResponse {
	return &KepegawaianIdentifikasiResponse{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedBy:   m.CreatedBy,
		UpdatedBy:   m.UpdatedBy,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// ToKepegawaianIdentifikasiListResponse mengubah slice model menjadi slice response
func ToKepegawaianIdentifikasiListResponse(items []models.KepegawaianIdentifikasi) []KepegawaianIdentifikasiResponse {
	var responses []KepegawaianIdentifikasiResponse
	for _, m := range items {
		responses = append(responses, *ToKepegawaianIdentifikasiResponse(&m))
	}
	return responses
}
