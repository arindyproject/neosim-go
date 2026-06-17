package dto

import (
	"time"

	"neosim_go/internal/modules/master/master/models"
)

// MasterResponse response untuk single Master
type MasterResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedBy   *int64    `json:"created_by"`
	UpdatedBy   *int64    `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToMasterResponse mengubah model menjadi response
func ToMasterResponse(m *models.Master) *MasterResponse {
	return &MasterResponse{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedBy:   m.CreatedBy,
		UpdatedBy:   m.UpdatedBy,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// ToMasterListResponse mengubah slice model menjadi slice response
func ToMasterListResponse(items []models.Master) []MasterResponse {
	var responses []MasterResponse
	for _, m := range items {
		responses = append(responses, *ToMasterResponse(&m))
	}
	return responses
}
