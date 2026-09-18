package dto

// CreateJobTitleRequest request body untuk membuat JobTitle baru
type CreateJobTitleRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateJobTitleRequest request body untuk update JobTitle
type UpdateJobTitleRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// FilterJobTitleRequest request body untuk filter JobTitle
type FilterJobTitleRequest struct {
	Name string `query:"name"`
}
