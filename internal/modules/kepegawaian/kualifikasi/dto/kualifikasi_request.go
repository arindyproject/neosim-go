package dto

// CreateKepegawaianKualifikasiRequest request body untuk membuat KepegawaianKualifikasi baru
type CreateKepegawaianKualifikasiRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateKepegawaianKualifikasiRequest request body untuk update KepegawaianKualifikasi
type UpdateKepegawaianKualifikasiRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// FilterKepegawaianKualifikasiRequest request body untuk filter KepegawaianKualifikasi
type FilterKepegawaianKualifikasiRequest struct {
	Name string `query:"name"`
}
