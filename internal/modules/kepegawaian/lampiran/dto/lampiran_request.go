package dto

// CreateKepegawaianLampiranRequest request body untuk membuat KepegawaianLampiran baru
type CreateKepegawaianLampiranRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateKepegawaianLampiranRequest request body untuk update KepegawaianLampiran
type UpdateKepegawaianLampiranRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// FilterKepegawaianLampiranRequest request body untuk filter KepegawaianLampiran
type FilterKepegawaianLampiranRequest struct {
	Name string `query:"name"`
}
