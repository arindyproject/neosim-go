package dto

// CreateStatusRequest request body untuk membuat Status baru
type CreateStatusRequest struct {
	Code     string  `json:"code" validate:"required,min=1,max=255"`
	Label    string  `json:"label" validate:"required,min=1,max=255"`
	FHIRCode *string `json:"fhir_code" validate:"omitempty,max=500"`
}

// UpdateStatusRequest request body untuk update Status
type UpdateStatusRequest struct {
	Code     *string `json:"code" validate:"required,min=1,max=255"`
	Label    *string `json:"label" validate:"required,min=1,max=255"`
	FHIRCode *string `json:"fhir_code" validate:"omitempty,max=500"`
}

// FilterStatusRequest request body untuk filter Status
type FilterStatusRequest struct {
	Search string `query:"search"`
	Code   string `query:"code"`
	Label  string `query:"label"`

	// Sorting
	SortBy    string `query:"sort_by"`    // code | label | fhir_code | created_at | updated_at
	SortOrder string `query:"sort_order"` // asc | desc
}
