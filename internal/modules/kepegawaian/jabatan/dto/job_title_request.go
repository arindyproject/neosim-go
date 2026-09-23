package dto

// CreateJobTitleRequest request body untuk membuat JobTitle baru
type CreateJobTitleRequest struct {
	Code        string  `json:"code" validate:"required,min=1,max=20"`
	Label       string  `json:"label" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`

	KategoriID      int64  `json:"kategori_id" validate:"required,gt=0"`
	RumpunProfesiID *int64 `json:"rumpun_profesi_id" validate:"omitempty,gt=0"`

	// Point jabatan (mis. untuk perhitungan remunerasi/tunjangan). Nullable —
	// DB default 0.00 kalau tidak diisi.
	Point *float64 `json:"point" validate:"omitempty,gte=0"`

	MemerlukanSTR bool `json:"memerlukan_str"`
	MemerlukanSIP bool `json:"memerlukan_sip"`

	JenjangMin *string `json:"jenjang_min" validate:"omitempty,max=20"`

	FHIRCode   *string `json:"fhir_code" validate:"omitempty,max=50"`
	FHIRSystem *string `json:"fhir_system" validate:"omitempty,max=200"`

	IsAktif bool `json:"is_aktif"`
}

// UpdateJobTitleRequest request body untuk update JobTitle
type UpdateJobTitleRequest struct {
	Code        *string `json:"code" validate:"omitempty,min=1,max=20"`
	Label       *string `json:"label" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=500"`

	KategoriID      *int64 `json:"kategori_id" validate:"omitempty,gt=0"`
	RumpunProfesiID *int64 `json:"rumpun_profesi_id" validate:"omitempty,gt=0"`

	Point *float64 `json:"point" validate:"omitempty,gte=0"`

	MemerlukanSTR *bool `json:"memerlukan_str"`
	MemerlukanSIP *bool `json:"memerlukan_sip"`

	JenjangMin *string `json:"jenjang_min" validate:"omitempty,max=20"`

	FHIRCode   *string `json:"fhir_code" validate:"omitempty,max=50"`
	FHIRSystem *string `json:"fhir_system" validate:"omitempty,max=200"`

	IsAktif *bool `json:"is_aktif"`
}

// FilterJobTitleRequest request query untuk filter JobTitle
type FilterJobTitleRequest struct {
	Label           string `query:"label"`
	Code            string `query:"code"`
	KategoriID      *int64 `query:"kategori_id"`
	RumpunProfesiID *int64 `query:"rumpun_profesi_id"`
	MemerlukanSTR   *bool  `query:"memerlukan_str"`
	MemerlukanSIP   *bool  `query:"memerlukan_sip"`
	IsAktif         *bool  `query:"is_aktif"`
}
