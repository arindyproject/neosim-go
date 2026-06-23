package dto

// CreateMasterDepartemenRequest request body untuk membuat MasterDepartemen baru
type CreateMasterDepartemenRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateMasterDepartemenRequest request body untuk update MasterDepartemen
type UpdateMasterDepartemenRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// FilterMasterDepartemenRequest request body untuk filter MasterDepartemen
type FilterMasterDepartemenRequest struct {
	Name string `query:"name"`
}
