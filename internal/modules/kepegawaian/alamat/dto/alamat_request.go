package dto

// CreateKepegawaianAlamatRequest request body untuk membuat KepegawaianAlamat baru
type CreateKepegawaianAlamatRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateKepegawaianAlamatRequest request body untuk update KepegawaianAlamat
type UpdateKepegawaianAlamatRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// FilterKepegawaianAlamatRequest request body untuk filter KepegawaianAlamat
type FilterKepegawaianAlamatRequest struct {
	Name string `query:"name"`
}
