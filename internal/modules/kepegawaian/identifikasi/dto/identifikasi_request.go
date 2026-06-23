package dto

// CreateKepegawaianIdentifikasiRequest request body untuk membuat KepegawaianIdentifikasi baru
type CreateKepegawaianIdentifikasiRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateKepegawaianIdentifikasiRequest request body untuk update KepegawaianIdentifikasi
type UpdateKepegawaianIdentifikasiRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// FilterKepegawaianIdentifikasiRequest request body untuk filter KepegawaianIdentifikasi
type FilterKepegawaianIdentifikasiRequest struct {
	Name string `query:"name"`
}
