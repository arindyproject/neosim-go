package dto

// CreateMasterRequest request body untuk membuat Master baru
type CreateMasterRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateMasterRequest request body untuk update Master
type UpdateMasterRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// FilterMasterRequest request body untuk filter Master
type FilterMasterRequest struct {
	Name        string `query:"name"`
}

