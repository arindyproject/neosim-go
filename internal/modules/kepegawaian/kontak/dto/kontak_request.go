package dto

// CreateKepegawaianKontakRequest request body untuk membuat KepegawaianKontak baru
type CreateKepegawaianKontakRequest struct {
	PegawaiID   int64   `json:"pegawai_id" validate:"required"`
	TipeID      int64   `json:"tipe_id" validate:"required"`
	Nilai       string  `json:"nilai" validate:"required,max=225"`
	IsPrimary   bool    `json:"is_primary"`
	IsAktif     bool    `json:"is_aktif"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateKepegawaianKontakRequest request body untuk update KepegawaianKontak
type UpdateKepegawaianKontakRequest struct {
	TipeID      *int64  `json:"tipe_id" validate:"omitempty"`
	Nilai       *string `json:"nilai" validate:"omitempty,max=225"`
	IsPrimary   *bool   `json:"is_primary" validate:"omitempty"`
	IsAktif     *bool   `json:"is_aktif" validate:"omitempty"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// FilterKepegawaianKontakRequest request body untuk filter KepegawaianKontak
type FilterKepegawaianKontakRequest struct {
	PegawaiID *int64  `query:"pegawai_id"`
	TipeID    *int64  `query:"tipe_id"`
	Nilai     *string `query:"nilai"`
	IsPrimary *bool   `query:"is_primary"`
	IsAktif   *bool   `query:"is_aktif"`
}
