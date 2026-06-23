package dto

// CreateKepegawaianJabatanRequest request body untuk membuat KepegawaianJabatan baru
type CreateKepegawaianJabatanRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateKepegawaianJabatanRequest request body untuk update KepegawaianJabatan
type UpdateKepegawaianJabatanRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// FilterKepegawaianJabatanRequest request body untuk filter KepegawaianJabatan
type FilterKepegawaianJabatanRequest struct {
	Name string `query:"name"`
}
