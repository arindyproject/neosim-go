package dto

import (
	"time"

	"neosim_go/internal/modules/kepegawaian/kualifikasi/models"
)

// KepegawaianKualifikasiResponse response untuk single KepegawaianKualifikasi
type KepegawaianKualifikasiResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedBy   *int64    `json:"created_by"`
	UpdatedBy   *int64    `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToKepegawaianKualifikasiResponse mengubah model menjadi response
func ToKepegawaianKualifikasiResponse(m *models.KepegawaianKualifikasi) *KepegawaianKualifikasiResponse {
	return &KepegawaianKualifikasiResponse{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedBy:   m.CreatedBy,
		UpdatedBy:   m.UpdatedBy,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// ToKepegawaianKualifikasiListResponse mengubah slice model menjadi slice response
func ToKepegawaianKualifikasiListResponse(items []models.KepegawaianKualifikasi) []KepegawaianKualifikasiResponse {
	var responses []KepegawaianKualifikasiResponse
	for _, m := range items {
		responses = append(responses, *ToKepegawaianKualifikasiResponse(&m))
	}
	return responses
}
