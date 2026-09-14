package dto

// CreateTipeRequest request body untuk membuat Tipe baru
type CreateTipeRequest struct {
	Code     string  `json:"code" validate:"required,min=1,max=255"`
	Label    string  `json:"label" validate:"required,min=1,max=255"`
	FHIRCode *string `json:"fhir_code" validate:"omitempty,max=500"`
}

// UpdateTipeRequest request body untuk update Tipe
type UpdateTipeRequest struct {
	Code     *string `json:"code" validate:"omitempty,min=1,max=255"`
	Label    *string `json:"label" validate:"omitempty,min=1,max=255"`
	FHIRCode *string `json:"fhir_code" validate:"omitempty,max=500"`
}

// FilterTipeRequest request body untuk filter Tipe
type FilterTipeRequest struct {
	Code  string `query:"code"`
	Label string `query:"label"`
}
