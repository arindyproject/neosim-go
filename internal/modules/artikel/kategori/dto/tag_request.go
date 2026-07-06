package dto

// CreateTagRequest request body untuk membuat Tag baru
type CreateTagRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateTagRequest request body untuk update Tag
type UpdateTagRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// FilterTagRequest request body untuk filter Tag
type FilterTagRequest struct {
	Name string `query:"name"`
}
