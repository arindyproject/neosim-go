package dto

// CreateJenjangRequest request body untuk membuat Jenjang baru
type CreateJenjangRequest struct {
	Code       string  `json:"code" validate:"required,min=1,max=100"`
	Label      string  `json:"label" validate:"required,min=1,max=255"`
	FHIRSystem *string `json:"fhir_system" validate:"omitempty,max=500"`
}

// UpdateJenjangRequest request body untuk update Jenjang
type UpdateJenjangRequest struct {
	Code       *string `json:"code" validate:"omitempty,min=1,max=100"`
	Label      *string `json:"label" validate:"omitempty,min=1,max=255"`
	FHIRSystem *string `json:"fhir_system" validate:"omitempty,max=500"`
}

// FilterJenjangRequest request body untuk filter Jenjang
type FilterJenjangRequest struct {
	Code  string `query:"code"`
	Label string `query:"label"`
}
