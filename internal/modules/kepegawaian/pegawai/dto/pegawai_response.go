package dto

import (
	"time"

	"neosim_go/internal/modules/kepegawaian/pegawai/models"
)

// KepegawaianPegawaiResponse response untuk single KepegawaianPegawai
type KepegawaianPegawaiResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedBy   *int64    `json:"created_by"`
	UpdatedBy   *int64    `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToKepegawaianPegawaiResponse mengubah model menjadi response
func ToKepegawaianPegawaiResponse(m *models.KepegawaianPegawai) *KepegawaianPegawaiResponse {
	return &KepegawaianPegawaiResponse{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedBy:   m.CreatedBy,
		UpdatedBy:   m.UpdatedBy,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// ToKepegawaianPegawaiListResponse mengubah slice model menjadi slice response
func ToKepegawaianPegawaiListResponse(items []models.KepegawaianPegawai) []KepegawaianPegawaiResponse {
	var responses []KepegawaianPegawaiResponse
	for _, m := range items {
		responses = append(responses, *ToKepegawaianPegawaiResponse(&m))
	}
	return responses
}
