package dto

// CreateJenisRequest request body untuk membuat Jenis baru
type CreateJenisRequest struct {
	Code     string  `json:"code" validate:"required,min=1,max=255"`
	Label    string  `json:"label" validate:"required,min=1,max=255"`
	FHIRCode *string `json:"fhir_code" validate:"omitempty,max=500"`
}

// UpdateJenisRequest request body untuk update Jenis
type UpdateJenisRequest struct {
	Code     *string `json:"code" validate:"required,min=1,max=255"`
	Label    *string `json:"label" validate:"required,min=1,max=255"`
	FHIRCode *string `json:"fhir_code" validate:"omitempty,max=500"`
}

// FilterJenisRequest request body untuk filter Jenis
type FilterJenisRequest struct {
	Search string `query:"search"`
	Code   string `query:"code"`
	Label  string `query:"label"`
}
