package dto

// CreateJobTitleKategoriRequest request body untuk membuat JobTitleKategori baru
type CreateJobTitleKategoriRequest struct {
	Code     string  `json:"code" validate:"required,min=1,max=255"`
	Label    string  `json:"label" validate:"required,min=1,max=255"`
	FHIRCode *string `json:"fhir_code" validate:"omitempty,max=500"`
}

// UpdateJobTitleKategoriRequest request body untuk update JobTitleKategori
type UpdateJobTitleKategoriRequest struct {
	Code     *string `json:"code" validate:"required,min=1,max=255"`
	Label    *string `json:"label" validate:"required,min=1,max=255"`
	FHIRCode *string `json:"fhir_code" validate:"omitempty,max=500"`
}

// FilterJobTitleKategoriRequest request body untuk filter JobTitleKategori
type FilterJobTitleKategoriRequest struct {
	Search string `query:"search"`
	Code   string `query:"code"`
	Label  string `query:"label"`
}
