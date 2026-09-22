package dto

// CreatePositionKategoriRequest request body untuk membuat PositionKategori baru
type CreatePositionKategoriRequest struct {
	Code     string   `json:"code" validate:"required,min=1,max=100"`
	Label    string   `json:"label" validate:"required,min=1,max=255"`
	FHIRCode *string  `json:"fhir_code" validate:"omitempty,max=500"`
	Point    *float64 `json:"point" validate:"omitempty,gte=0"`
}

// UpdatePositionKategoriRequest request body untuk update PositionKategori
type UpdatePositionKategoriRequest struct {
	Code     *string  `json:"code" validate:"required,min=1,max=100"`
	Label    *string  `json:"label" validate:"required,min=1,max=255"`
	FHIRCode *string  `json:"fhir_code" validate:"omitempty,max=500"`
	Point    *float64 `json:"point" validate:"omitempty,gte=0"`
}

// FilterPositionKategoriRequest request body untuk filter PositionKategori
type FilterPositionKategoriRequest struct {
	Search string `query:"search"`
	Code   string `query:"code"`
	Label  string `query:"label"`
}
