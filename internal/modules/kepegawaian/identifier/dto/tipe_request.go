package dto

// CreateTipeRequest request body untuk membuat Tipe baru
type CreateTipeRequest struct {
	Code        string  `json:"code" validate:"required,min=1,max=100"`
	Label       string  `json:"label" validate:"required,min=1,max=255"`
	Penerbit    *string `json:"penerbit" validate:"omitempty,max=255"`
	FHIRSystem  *string `json:"fhir_system" validate:"omitempty,max=255"`
	HasExpiry   bool    `json:"has_expiry"`
	IsNakes     bool    `json:"is_nakes"`
	IsRequired  bool    `json:"is_required"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateTipeRequest request body untuk update Tipe (parsial/optional)
type UpdateTipeRequest struct {
	Code        *string `json:"code" validate:"omitempty,min=1,max=100"`
	Label       *string `json:"label" validate:"omitempty,min=1,max=255"`
	Penerbit    *string `json:"penerbit" validate:"omitempty,max=255"`
	FHIRSystem  *string `json:"fhir_system" validate:"omitempty,max=255"`
	HasExpiry   *bool   `json:"has_expiry"`
	IsNakes     *bool   `json:"is_nakes"`
	IsRequired  *bool   `json:"is_required"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// FilterTipeRequest request query/body untuk filter Tipe
type FilterTipeRequest struct {
	Search     string `query:"search"`
	Code       string `query:"code"`
	Label      string `query:"label"`
	IsNakes    *bool  `query:"is_nakes"`
	IsRequired *bool  `query:"is_required"`
}
