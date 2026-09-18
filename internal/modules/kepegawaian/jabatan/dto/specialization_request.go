package dto

// CreateSpecializationRequest request body untuk membuat Specialization baru
type CreateSpecializationRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateSpecializationRequest request body untuk update Specialization
type UpdateSpecializationRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// FilterSpecializationRequest request body untuk filter Specialization
type FilterSpecializationRequest struct {
	Name string `query:"name"`
}
