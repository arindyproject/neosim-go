package dto

// CreateJobTitleRumpunProfesiRequest request body untuk membuat JobTitleRumpunProfesi baru
type CreateJobTitleRumpunProfesiRequest struct {
	Code     string  `json:"code" validate:"required,min=1,max=255"`
	Label    string  `json:"label" validate:"required,min=1,max=255"`
	FHIRCode *string `json:"fhir_code" validate:"omitempty,max=500"`
}

// UpdateJobTitleRumpunProfesiRequest request body untuk update JobTitleRumpunProfesi
type UpdateJobTitleRumpunProfesiRequest struct {
	Code     *string `json:"code" validate:"required,min=1,max=255"`
	Label    *string `json:"label" validate:"required,min=1,max=255"`
	FHIRCode *string `json:"fhir_code" validate:"omitempty,max=500"`
}

// FilterJobTitleRumpunProfesiRequest request body untuk filter JobTitleRumpunProfesi
type FilterJobTitleRumpunProfesiRequest struct {
	Search string `query:"search"`
	Code   string `query:"code"`
	Label  string `query:"label"`
}
