package dto

// CreatePositionRequest request body untuk membuat Position baru
type CreatePositionRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// UpdatePositionRequest request body untuk update Position
type UpdatePositionRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// FilterPositionRequest request body untuk filter Position
type FilterPositionRequest struct {
	Name string `query:"name"`
}
