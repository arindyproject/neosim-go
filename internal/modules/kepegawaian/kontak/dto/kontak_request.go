package dto

// CreateKepegawaianKontakRequest request body untuk membuat KepegawaianKontak baru
type CreateKepegawaianKontakRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateKepegawaianKontakRequest request body untuk update KepegawaianKontak
type UpdateKepegawaianKontakRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// FilterKepegawaianKontakRequest request body untuk filter KepegawaianKontak
type FilterKepegawaianKontakRequest struct {
	Name string `query:"name"`
}
