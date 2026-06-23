package dto

// CreateKepegawaianPegawaiRequest request body untuk membuat KepegawaianPegawai baru
type CreateKepegawaianPegawaiRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateKepegawaianPegawaiRequest request body untuk update KepegawaianPegawai
type UpdateKepegawaianPegawaiRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// FilterKepegawaianPegawaiRequest request body untuk filter KepegawaianPegawai
type FilterKepegawaianPegawaiRequest struct {
	Name string `query:"name"`
}
