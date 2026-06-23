package dto

import (
	"time"

	"neosim_go/internal/modules/master/departemen/models"
)

// MasterDepartemenResponse response untuk single MasterDepartemen
type MasterDepartemenResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedBy   *int64    `json:"created_by"`
	UpdatedBy   *int64    `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToMasterDepartemenResponse mengubah model menjadi response
func ToMasterDepartemenResponse(m *models.MasterDepartemen) *MasterDepartemenResponse {
	return &MasterDepartemenResponse{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedBy:   m.CreatedBy,
		UpdatedBy:   m.UpdatedBy,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// ToMasterDepartemenListResponse mengubah slice model menjadi slice response
func ToMasterDepartemenListResponse(items []models.MasterDepartemen) []MasterDepartemenResponse {
	var responses []MasterDepartemenResponse
	for _, m := range items {
		responses = append(responses, *ToMasterDepartemenResponse(&m))
	}
	return responses
}
