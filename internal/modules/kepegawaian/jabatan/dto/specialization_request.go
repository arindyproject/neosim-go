package dto

// CreateSpecializationRequest request body untuk membuat Specialization baru
type CreateSpecializationRequest struct {
	Code                string  `json:"code" validate:"required,min=1,max=20"`
	Label               string  `json:"label" validate:"required,min=1,max=150"`
	JobTitleID          *int64  `json:"job_title_id" validate:"omitempty,gt=0"`
	KategoriID          *int64  `json:"kategori_id" validate:"omitempty,gt=0"`
	Gelar               *string `json:"gelar" validate:"omitempty,max=20"`
	LamaPendidikanTahun *int16  `json:"lama_pendidikan_tahun" validate:"omitempty,gte=0,lte=20"`
	FHIRCode            *string `json:"fhir_code" validate:"omitempty,max=50"`
	FHIRSystem          *string `json:"fhir_system" validate:"omitempty,max=200"`
	IsAktif             *bool   `json:"is_aktif"`
}

// UpdateSpecializationRequest request body untuk update Specialization
type UpdateSpecializationRequest struct {
	Code                *string `json:"code" validate:"omitempty,min=1,max=20"`
	Label               *string `json:"label" validate:"omitempty,min=1,max=150"`
	JobTitleID          *int64  `json:"job_title_id" validate:"omitempty,gt=0"`
	KategoriID          *int64  `json:"kategori_id" validate:"omitempty,gt=0"`
	Gelar               *string `json:"gelar" validate:"omitempty,max=20"`
	LamaPendidikanTahun *int16  `json:"lama_pendidikan_tahun" validate:"omitempty,gte=0,lte=20"`
	FHIRCode            *string `json:"fhir_code" validate:"omitempty,max=50"`
	FHIRSystem          *string `json:"fhir_system" validate:"omitempty,max=200"`
	IsAktif             *bool   `json:"is_aktif"`
}

// FilterSpecializationRequest request body untuk filter Specialization
type FilterSpecializationRequest struct {
	Code       string `query:"code"`
	Label      string `query:"label"`
	JobTitleID *int64 `query:"job_title_id"`
	KategoriID *int64 `query:"kategori_id"`
	IsAktif    *bool  `query:"is_aktif"`

	// Sorting
	SortBy    string `query:"sort_by"`    // code | label | job_title_id | kategori_id | is_aktif | created_at | updated_at
	SortOrder string `query:"sort_order"` // asc | desc
}
