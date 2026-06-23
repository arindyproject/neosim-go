package dto

import (
	"time"

	"neosim_go/internal/modules/kepegawaian/jabatan/models"
)

// KepegawaianJabatanResponse response untuk single KepegawaianJabatan
type KepegawaianJabatanResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedBy   *int64    `json:"created_by"`
	UpdatedBy   *int64    `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToKepegawaianJabatanResponse mengubah model menjadi response
func ToKepegawaianJabatanResponse(m *models.KepegawaianJabatan) *KepegawaianJabatanResponse {
	return &KepegawaianJabatanResponse{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedBy:   m.CreatedBy,
		UpdatedBy:   m.UpdatedBy,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// ToKepegawaianJabatanListResponse mengubah slice model menjadi slice response
func ToKepegawaianJabatanListResponse(items []models.KepegawaianJabatan) []KepegawaianJabatanResponse {
	var responses []KepegawaianJabatanResponse
	for _, m := range items {
		responses = append(responses, *ToKepegawaianJabatanResponse(&m))
	}
	return responses
}
