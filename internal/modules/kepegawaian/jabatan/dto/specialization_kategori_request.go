package dto

// CreateSpecializationKategoriRequest request body untuk membuat SpecializationKategori baru
type CreateSpecializationKategoriRequest struct {
	Code     string  `json:"code" validate:"required,min=1,max=255"`
	Label    string  `json:"label" validate:"required,min=1,max=255"`
	FHIRCode *string `json:"fhir_code" validate:"omitempty,max=500"`
}

// UpdateSpecializationKategoriRequest request body untuk update SpecializationKategori
type UpdateSpecializationKategoriRequest struct {
	Code     *string `json:"code" validate:"required,min=1,max=255"`
	Label    *string `json:"label" validate:"required,min=1,max=255"`
	FHIRCode *string `json:"fhir_code" validate:"omitempty,max=500"`
}

// FilterSpecializationKategoriRequest request body untuk filter SpecializationKategori
type FilterSpecializationKategoriRequest struct {
	Search string `query:"search"`
	Code   string `query:"code"`
	Label  string `query:"label"`
}
