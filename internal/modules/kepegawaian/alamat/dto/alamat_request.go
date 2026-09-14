package dto

// CreateKepegawaianAlamatRequest request body untuk membuat KepegawaianAlamat baru
type CreateKepegawaianAlamatRequest struct {
	PegawaiID int64 `json:"pegawai_id" validate:"required"`
	TipeID    int64 `json:"tipe_id" validate:"required"`

	Jalan   string  `json:"jalan" validate:"required,min=1,max=500"`
	RT      *string `json:"rt" validate:"omitempty,min=1,max=20"`
	RW      *string `json:"rw" validate:"omitempty,min=1,max=20"`
	KodePos *string `json:"kode_pos" validate:"omitempty,min=1,max=20"`

	NegaraID        *int64 `json:"negara_id" validate:"omitempty"`
	ProvinsiID      *int64 `json:"provinsi_id" validate:"omitempty"`
	KotaKabupatenID *int64 `json:"kota_kabupaten_id" validate:"omitempty"`
	KecamatanID     *int64 `json:"kecamatan_id" validate:"omitempty"`
	KelurahanDesaID *int64 `json:"kelurahan_desa_id" validate:"omitempty"`

	IsPrimary bool `json:"is_primary"`

	Description *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateKepegawaianAlamatRequest request body untuk update KepegawaianAlamat
type UpdateKepegawaianAlamatRequest struct {
	TipeID *int64 `json:"tipe_id" validate:"required"`

	Jalan   *string `json:"jalan" validate:"required,min=1,max=500"`
	RT      *string `json:"rt" validate:"omitempty,min=1,max=20"`
	RW      *string `json:"rw" validate:"omitempty,min=1,max=20"`
	KodePos *string `json:"kode_pos" validate:"omitempty,min=1,max=20"`

	NegaraID        *int64 `json:"negara_id" validate:"omitempty"`
	ProvinsiID      *int64 `json:"provinsi_id" validate:"omitempty"`
	KotaKabupatenID *int64 `json:"kota_kabupaten_id" validate:"omitempty"`
	KecamatanID     *int64 `json:"kecamatan_id" validate:"omitempty"`
	KelurahanDesaID *int64 `json:"kelurahan_desa_id" validate:"omitempty"`

	IsPrimary *bool `json:"is_primary"`

	Description *string `json:"description" validate:"omitempty,max=500"`
}

// FilterKepegawaianAlamatRequest request body untuk filter KepegawaianAlamat
type FilterKepegawaianAlamatRequest struct {
	TipeID *int64  `query:"tipe_id"`
	Jalan  *string `query:"jalan"`

	NegaraID        *int64 `query:"negara_id"`
	ProvinsiID      *int64 `query:"provinsi_id"`
	KotaKabupatenID *int64 `query:"kota_kabupaten_id"`
	KecamatanID     *int64 `query:"kecamatan_id"`
	KelurahanDesaID *int64 `query:"kelurahan_desa_id"`
}
