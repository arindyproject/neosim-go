package dto

// CreateTipeRequest request body untuk membuat Tipe baru
type CreateTipeRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateTipeRequest request body untuk update Tipe
type UpdateTipeRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// FilterTipeRequest request body untuk filter Tipe
type FilterTipeRequest struct {
	Name string `query:"name"`
}
