package dto

import (
	"time"

	"neosim_go/internal/modules/master/master/models"
)

// =====================================================================
// Pekerjaan
// =====================================================================
// MasterPekerjaanResponse response untuk single MasterPekerjaan
type MasterPekerjaanResponse struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	KodeKemenkes *string   `json:"kode_kemenkes"`
	FhirCode     *string   `json:"fhir_code"`
	Description  *string   `json:"description"`
	CreatedBy    *int64    `json:"created_by"`
	UpdatedBy    *int64    `json:"updated_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type MasterPekerjaanListSimpelResponse struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	KodeKemenkes *string `json:"kode_kemenkes"`
}

// ToMasterPekerjaanListSimpelResponse mengubah slice model menjadi slice response simpel
func ToMasterPekerjaanListSimpelResponse(items []models.MasterPekerjaan) []MasterPekerjaanListSimpelResponse {
	var responses []MasterPekerjaanListSimpelResponse
	for _, m := range items {
		responses = append(responses, MasterPekerjaanListSimpelResponse{
			ID:           m.ID,
			Name:         m.Name,
			KodeKemenkes: m.KodeKemenkes,
		})
	}
	return responses
} // -------------------------------------------------------------------

// ToMasterPekerjaanResponse mengubah model menjadi response
func ToMasterPekerjaanResponse(m *models.MasterPekerjaan) *MasterPekerjaanResponse {
	return &MasterPekerjaanResponse{
		ID:           m.ID,
		Name:         m.Name,
		KodeKemenkes: m.KodeKemenkes,
		FhirCode:     m.FhirCode,
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
	FhirCode     *string   `json:"fhir_code"`
	Description  *string   `json:"description"`
	CreatedBy    *int64    `json:"created_by"`
	UpdatedBy    *int64    `json:"updated_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type MasterPendidikanListSimpelResponse struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	KodeKemenkes *string `json:"kode_kemenkes"`
}

// ToMasterPendidikanListSimpelResponse mengubah slice model menjadi slice response simpel
func ToMasterPendidikanListSimpelResponse(items []models.MasterPendidikan) []MasterPendidikanListSimpelResponse {
	var responses []MasterPendidikanListSimpelResponse
	for _, m := range items {
		responses = append(responses, MasterPendidikanListSimpelResponse{
			ID:           m.ID,
			Name:         m.Name,
			KodeKemenkes: m.KodeKemenkes,
		})
	}
	return responses
} // -------------------------------------------------------------------

// ToMasterPendidikanResponse mengubah model menjadi response
func ToMasterPendidikanResponse(m *models.MasterPendidikan) *MasterPendidikanResponse {
	return &MasterPendidikanResponse{
		ID:           m.ID,
		Name:         m.Name,
		KodeKemenkes: m.KodeKemenkes,
		FhirCode:     m.FhirCode,
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
	FhirCode     *string   `json:"fhir_code"`
	Description  *string   `json:"description"`
	CreatedBy    *int64    `json:"created_by"`
	UpdatedBy    *int64    `json:"updated_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type MasterStatusPernikahanListSimpelResponse struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	KodeKemenkes *string `json:"kode_kemenkes"`
}

// ToMasterStatusPernikahanListSimpelResponse mengubah slice model menjadi slice response simpel
func ToMasterStatusPernikahanListSimpelResponse(items []models.MasterStatusPernikahan) []MasterStatusPernikahanListSimpelResponse {
	var responses []MasterStatusPernikahanListSimpelResponse
	for _, m := range items {
		responses = append(responses, MasterStatusPernikahanListSimpelResponse{
			ID:           m.ID,
			Name:         m.Name,
			KodeKemenkes: m.KodeKemenkes,
		})
	}
	return responses
} // -------------------------------------------------------------------

// ToMasterStatusPernikahanResponse mengubah model menjadi response
func ToMasterStatusPernikahanResponse(m *models.MasterStatusPernikahan) *MasterStatusPernikahanResponse {
	return &MasterStatusPernikahanResponse{
		ID:           m.ID,
		Name:         m.Name,
		KodeKemenkes: m.KodeKemenkes,
		FhirCode:     m.FhirCode,
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
	FhirCode     *string   `json:"fhir_code"`
	Description  *string   `json:"description"`
	CreatedBy    *int64    `json:"created_by"`
	UpdatedBy    *int64    `json:"updated_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type MasterAgamaListSimpelResponse struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	KodeKemenkes *string `json:"kode_kemenkes"`
}

func ToMasterAgamaListSimpelResponse(items []models.MasterAgama) []MasterAgamaListSimpelResponse {
	var responses []MasterAgamaListSimpelResponse
	for _, m := range items {
		responses = append(responses, MasterAgamaListSimpelResponse{
			ID:           m.ID,
			Name:         m.Name,
			KodeKemenkes: m.KodeKemenkes,
		})
	}
	return responses
} // -------------------------------------------------------------------

// ToMasterAgamaResponse mengubah model menjadi response
func ToMasterAgamaResponse(m *models.MasterAgama) *MasterAgamaResponse {
	return &MasterAgamaResponse{
		ID:           m.ID,
		Name:         m.Name,
		KodeKemenkes: m.KodeKemenkes,
		FhirCode:     m.FhirCode,
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

// =====================================================================
// JenisKelamin
// =====================================================================
// MasterJenisKelaminResponse response untuk single MasterJenisKelamin
type MasterJenisKelaminResponse struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	KodeKemenkes *string   `json:"kode_kemenkes"`
	FhirCode     *string   `json:"fhir_code"`
	Description  *string   `json:"description"`
	CreatedBy    *int64    `json:"created_by"`
	UpdatedBy    *int64    `json:"updated_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type MasterJenisKelaminListSimpelResponse struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	KodeKemenkes *string `json:"kode_kemenkes"`
}

func ToMasterJenisKelaminListSimpelResponse(items []models.MasterJenisKelamin) []MasterJenisKelaminListSimpelResponse {
	var responses []MasterJenisKelaminListSimpelResponse
	for _, m := range items {
		responses = append(responses, MasterJenisKelaminListSimpelResponse{
			ID:           m.ID,
			Name:         m.Name,
			KodeKemenkes: m.KodeKemenkes,
		})
	}
	return responses
} // -------------------------------------------------------------------

// ToMasterJenisKelaminResponse mengubah model menjadi response
func ToMasterJenisKelaminResponse(m *models.MasterJenisKelamin) *MasterJenisKelaminResponse {
	return &MasterJenisKelaminResponse{
		ID:           m.ID,
		Name:         m.Name,
		KodeKemenkes: m.KodeKemenkes,
		FhirCode:     m.FhirCode,
		Description:  m.Description,
		CreatedBy:    m.CreatedBy,
		UpdatedBy:    m.UpdatedBy,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

// ToMasterJenisKelaminListResponse mengubah slice model menjadi slice response
func ToMasterJenisKelaminListResponse(items []models.MasterJenisKelamin) []MasterJenisKelaminResponse {
	var responses []MasterJenisKelaminResponse
	for _, m := range items {
		responses = append(responses, *ToMasterJenisKelaminResponse(&m))
	}
	return responses
} // -------------------------------------------------------------------

// =====================================================================
// GolonganDarah
// =====================================================================
// MasterGolonganDarahResponse response untuk single MasterGolonganDarah
type MasterGolonganDarahResponse struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	KodeKemenkes *string   `json:"kode_kemenkes"`
	FhirCode     *string   `json:"fhir_code"`
	Description  *string   `json:"description"`
	CreatedBy    *int64    `json:"created_by"`
	UpdatedBy    *int64    `json:"updated_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type MasterGolonganDarahListSimpelResponse struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	KodeKemenkes *string `json:"kode_kemenkes"`
}

func ToMasterGolonganDarahListSimpelResponse(items []models.MasterGolonganDarah) []MasterGolonganDarahListSimpelResponse {
	var responses []MasterGolonganDarahListSimpelResponse
	for _, m := range items {
		responses = append(responses, MasterGolonganDarahListSimpelResponse{
			ID:           m.ID,
			Name:         m.Name,
			KodeKemenkes: m.KodeKemenkes,
		})
	}
	return responses
} // -------------------------------------------------------------------

// ToMasterGolonganDarahResponse mengubah model menjadi response
func ToMasterGolonganDarahResponse(m *models.MasterGolonganDarah) *MasterGolonganDarahResponse {
	return &MasterGolonganDarahResponse{
		ID:           m.ID,
		Name:         m.Name,
		KodeKemenkes: m.KodeKemenkes,
		FhirCode:     m.FhirCode,
		Description:  m.Description,
		CreatedBy:    m.CreatedBy,
		UpdatedBy:    m.UpdatedBy,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

// ToMasterGolonganDarahListResponse mengubah slice model menjadi slice response
func ToMasterGolonganDarahListResponse(items []models.MasterGolonganDarah) []MasterGolonganDarahResponse {
	var responses []MasterGolonganDarahResponse
	for _, m := range items {
		responses = append(responses, *ToMasterGolonganDarahResponse(&m))
	}
	return responses
} // -------------------------------------------------------------------

// =====================================================================
// Suku
// =====================================================================
// MasterSukuResponse response untuk single MasterSuku
type MasterSukuResponse struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	KodeKemenkes *string   `json:"kode_kemenkes"`
	FhirCode     *string   `json:"fhir_code"`
	Description  *string   `json:"description"`
	CreatedBy    *int64    `json:"created_by"`
	UpdatedBy    *int64    `json:"updated_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type MasterSukuListSimpelResponse struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	KodeKemenkes *string `json:"kode_kemenkes"`
}

func ToMasterSukuListSimpelResponse(items []models.MasterSuku) []MasterSukuListSimpelResponse {
	var responses []MasterSukuListSimpelResponse
	for _, m := range items {
		responses = append(responses, MasterSukuListSimpelResponse{
			ID:           m.ID,
			Name:         m.Name,
			KodeKemenkes: m.KodeKemenkes,
		})
	}
	return responses
} // -------------------------------------------------------------------

// ToMasterSukuResponse mengubah model menjadi response
func ToMasterSukuResponse(m *models.MasterSuku) *MasterSukuResponse {
	return &MasterSukuResponse{
		ID:           m.ID,
		Name:         m.Name,
		KodeKemenkes: m.KodeKemenkes,
		FhirCode:     m.FhirCode,
		Description:  m.Description,
		CreatedBy:    m.CreatedBy,
		UpdatedBy:    m.UpdatedBy,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

// ToMasterSukuListResponse mengubah slice model menjadi slice response
func ToMasterSukuListResponse(items []models.MasterSuku) []MasterSukuResponse {
	var responses []MasterSukuResponse
	for _, m := range items {
		responses = append(responses, *ToMasterSukuResponse(&m))
	}
	return responses
} // -------------------------------------------------------------------
