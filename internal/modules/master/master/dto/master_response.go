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

// =====================================================================
// Pekerjaan
// =====================================================================
// MasterPekerjaanResponse response untuk single MasterPekerjaan
type MasterPekerjaanResponse struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	KodeKemenkes *string   `json:"kode_kemenkes"`
	Description  *string   `json:"description"`
	CreatedBy    *int64    `json:"created_by"`
	UpdatedBy    *int64    `json:"updated_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ToMasterPekerjaanResponse mengubah model menjadi response
func ToMasterPekerjaanResponse(m *models.MasterPekerjaan) *MasterPekerjaanResponse {
	return &MasterPekerjaanResponse{
		ID:           m.ID,
		Name:         m.Name,
		KodeKemenkes: m.KodeKemenkes,
		Description:  m.Description,
		CreatedBy:    m.CreatedBy,
		UpdatedBy:    m.UpdatedBy,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

// ToMasterPekerjaanListResponse mengubah slice model menjadi slice response
func ToMasterPekerjaanListResponse(items []models.MasterPekerjaan) []MasterPekerjaanResponse {
	var responses []MasterPekerjaanResponse
	for _, m := range items {
		responses = append(responses, *ToMasterPekerjaanResponse(&m))
	}
	return responses
} // -------------------------------------------------------------------

// =====================================================================
// Pendidikan
// =====================================================================
// MasterPendidikanResponse response untuk single MasterPendidikan
type MasterPendidikanResponse struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	KodeKemenkes *string   `json:"kode_kemenkes"`
	Description  *string   `json:"description"`
	CreatedBy    *int64    `json:"created_by"`
	UpdatedBy    *int64    `json:"updated_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ToMasterPendidikanResponse mengubah model menjadi response
func ToMasterPendidikanResponse(m *models.MasterPendidikan) *MasterPendidikanResponse {
	return &MasterPendidikanResponse{
		ID:           m.ID,
		Name:         m.Name,
		KodeKemenkes: m.KodeKemenkes,
		Description:  m.Description,
		CreatedBy:    m.CreatedBy,
		UpdatedBy:    m.UpdatedBy,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

// ToMasterPendidikanListResponse mengubah slice model menjadi slice response
func ToMasterPendidikanListResponse(items []models.MasterPendidikan) []MasterPendidikanResponse {
	var responses []MasterPendidikanResponse
	for _, m := range items {
		responses = append(responses, *ToMasterPendidikanResponse(&m))
	}
	return responses
} // -------------------------------------------------------------------

// =====================================================================
// StatusPernikahan
// =====================================================================
// MasterStatusPernikahanResponse response untuk single MasterStatusPernikahan
type MasterStatusPernikahanResponse struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	KodeKemenkes *string   `json:"kode_kemenkes"`
	Description  *string   `json:"description"`
	CreatedBy    *int64    `json:"created_by"`
	UpdatedBy    *int64    `json:"updated_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ToMasterStatusPernikahanResponse mengubah model menjadi response
func ToMasterStatusPernikahanResponse(m *models.MasterStatusPernikahan) *MasterStatusPernikahanResponse {
	return &MasterStatusPernikahanResponse{
		ID:           m.ID,
		Name:         m.Name,
		KodeKemenkes: m.KodeKemenkes,
		Description:  m.Description,
		CreatedBy:    m.CreatedBy,
		UpdatedBy:    m.UpdatedBy,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

// ToMasterStatusPernikahanListResponse mengubah slice model menjadi slice response
func ToMasterStatusPernikahanListResponse(items []models.MasterStatusPernikahan) []MasterStatusPernikahanResponse {
	var responses []MasterStatusPernikahanResponse
	for _, m := range items {
		responses = append(responses, *ToMasterStatusPernikahanResponse(&m))
	}
	return responses
} // -------------------------------------------------------------------

// =====================================================================
// Agama
// =====================================================================
// MasterAgamaResponse response untuk single MasterAgama
type MasterAgamaResponse struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	KodeKemenkes *string   `json:"kode_kemenkes"`
	Description  *string   `json:"description"`
	CreatedBy    *int64    `json:"created_by"`
	UpdatedBy    *int64    `json:"updated_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ToMasterAgamaResponse mengubah model menjadi response
func ToMasterAgamaResponse(m *models.MasterAgama) *MasterAgamaResponse {
	return &MasterAgamaResponse{
		ID:           m.ID,
		Name:         m.Name,
		KodeKemenkes: m.KodeKemenkes,
		Description:  m.Description,
		CreatedBy:    m.CreatedBy,
		UpdatedBy:    m.UpdatedBy,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

// ToMasterAgamaListResponse mengubah slice model menjadi slice response
func ToMasterAgamaListResponse(items []models.MasterAgama) []MasterAgamaResponse {
	var responses []MasterAgamaResponse
	for _, m := range items {
		responses = append(responses, *ToMasterAgamaResponse(&m))
	}
	return responses
} // -------------------------------------------------------------------
