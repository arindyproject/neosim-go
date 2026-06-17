package dto

// =====================================================================
// NEGARA
// =====================================================================

// CreateNegaraRequest request body untuk membuat Negara baru
type CreateNegaraRequest struct {
	Code        string  `json:"code" validate:"required,min=2,max=5"`
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// UpdateNegaraRequest request body untuk update Negara
type UpdateNegaraRequest struct {
	Code        *string `json:"code" validate:"omitempty,min=2,max=5"`
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// FilterNegaraRequest request query untuk filter Negara
type FilterNegaraRequest struct {
	Code string `query:"code"`
	Name string `query:"name"`
}

// =====================================================================
// PROVINSI
// =====================================================================

// CreateProvinsiRequest request body untuk membuat Provinsi baru
type CreateProvinsiRequest struct {
	NegaraID int64  `json:"negara_id" validate:"required,gt=0"`
	Code     string `json:"code" validate:"required,min=1,max=10"`
	Name     string `json:"name" validate:"required,min=1,max=255"`
}

// UpdateProvinsiRequest request body untuk update Provinsi
type UpdateProvinsiRequest struct {
	NegaraID *int64  `json:"negara_id" validate:"omitempty,gt=0"`
	Code     *string `json:"code" validate:"omitempty,min=1,max=10"`
	Name     *string `json:"name" validate:"omitempty,min=1,max=255"`
}

// FilterProvinsiRequest request query untuk filter Provinsi
type FilterProvinsiRequest struct {
	NegaraID *int64 `query:"negara_id"`
	Code     string `query:"code"`
	Name     string `query:"name"`
}

// =====================================================================
// KOTA / KABUPATEN
// =====================================================================

// CreateKotaKabupatenRequest request body untuk membuat Kota/Kabupaten baru
type CreateKotaKabupatenRequest struct {
	ProvinsiID int64  `json:"provinsi_id" validate:"required,gt=0"`
	Code       string `json:"code" validate:"required,min=1,max=10"`
	Name       string `json:"name" validate:"required,min=1,max=255"`
}

// UpdateKotaKabupatenRequest request body untuk update Kota/Kabupaten
type UpdateKotaKabupatenRequest struct {
	ProvinsiID *int64  `json:"provinsi_id" validate:"omitempty,gt=0"`
	Code       *string `json:"code" validate:"omitempty,min=1,max=10"`
	Name       *string `json:"name" validate:"omitempty,min=1,max=255"`
}

// FilterKotaKabupatenRequest request query untuk filter Kota/Kabupaten
type FilterKotaKabupatenRequest struct {
	ProvinsiID *int64 `query:"provinsi_id"`
	Code       string `query:"code"`
	Name       string `query:"name"`
}

// =====================================================================
// KECAMATAN
// =====================================================================

// CreateKecamatanRequest request body untuk membuat Kecamatan baru
type CreateKecamatanRequest struct {
	KotaKabupatenID int64  `json:"kota_kabupaten_id" validate:"required,gt=0"`
	Code            string `json:"code" validate:"required,min=1,max=10"`
	Name            string `json:"name" validate:"required,min=1,max=255"`
}

// UpdateKecamatanRequest request body untuk update Kecamatan
type UpdateKecamatanRequest struct {
	KotaKabupatenID *int64  `json:"kota_kabupaten_id" validate:"omitempty,gt=0"`
	Code            *string `json:"code" validate:"omitempty,min=1,max=10"`
	Name            *string `json:"name" validate:"omitempty,min=1,max=255"`
}

// FilterKecamatanRequest request query untuk filter Kecamatan
type FilterKecamatanRequest struct {
	KotaKabupatenID *int64 `query:"kota_kabupaten_id"`
	Code            string `query:"code"`
	Name            string `query:"name"`
}

// =====================================================================
// KELURAHAN / DESA
// =====================================================================

// CreateKelurahanDesaRequest request body untuk membuat Kelurahan/Desa baru
type CreateKelurahanDesaRequest struct {
	KecamatanID int64   `json:"kecamatan_id" validate:"required,gt=0"`
	Code        string  `json:"code" validate:"required,min=1,max=15"`
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	PostalCode  *string `json:"postal_code" validate:"omitempty,max=10"`
}

// UpdateKelurahanDesaRequest request body untuk update Kelurahan/Desa
type UpdateKelurahanDesaRequest struct {
	KecamatanID *int64  `json:"kecamatan_id" validate:"omitempty,gt=0"`
	Code        *string `json:"code" validate:"omitempty,min=1,max=15"`
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	PostalCode  *string `json:"postal_code" validate:"omitempty,max=10"`
}

// FilterKelurahanDesaRequest request query untuk filter Kelurahan/Desa
type FilterKelurahanDesaRequest struct {
	KecamatanID *int64 `query:"kecamatan_id"`
	Code        string `query:"code"`
	Name        string `query:"name"`
	PostalCode  string `query:"postal_code"`
}
