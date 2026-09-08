package dto

// CreateMasterDepartemenRequest request body untuk membuat MasterDepartemen baru
type CreateMasterDepartemenRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	FhirCode    *string `json:"fhir_code" validate:"omitempty,max=255"`
	FhirSystem  *string `json:"fhir_system" validate:"omitempty,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateMasterDepartemenRequest request body untuk update MasterDepartemen
type UpdateMasterDepartemenRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	FhirCode    *string `json:"fhir_code" validate:"omitempty,max=255"`
	FhirSystem  *string `json:"fhir_system" validate:"omitempty,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// FilterMasterDepartemenRequest request body untuk filter MasterDepartemen
type FilterMasterDepartemenRequest struct {
	Name       string  `query:"name"`
	FhirCode   *string `query:"fhir_code"`
	FhirSystem *string `query:"fhir_system"`
}
