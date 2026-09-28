package dto

// CreateKepegawaianJabatanRequest request body untuk membuat KepegawaianJabatan baru
type CreateKepegawaianJabatanRequest struct {
	PegawaiID        int64  `json:"pegawai_id" validate:"required,gt=0"`
	DepartmentID     int64  `json:"department_id" validate:"required,gt=0"`
	PositionID       int64  `json:"position_id" validate:"required,gt=0"`
	JobTitleID       int64  `json:"job_title_id" validate:"required,gt=0"`
	SpecializationID *int64 `json:"specialization_id" validate:"omitempty,gt=0"`

	IsPrimary *bool `json:"is_primary"`

	// Format tanggal: YYYY-MM-DD
	TanggalMulai   string  `json:"tanggal_mulai" validate:"required,datetime=2006-01-02"`
	TanggalSelesai *string `json:"tanggal_selesai" validate:"omitempty,datetime=2006-01-02"`

	NomorSK   *string `json:"nomor_sk" validate:"omitempty,max=100"`
	TanggalSK *string `json:"tanggal_sk" validate:"omitempty,datetime=2006-01-02"`

	IsAktif *bool `json:"is_aktif"`
}

// UpdateKepegawaianJabatanRequest request body untuk update KepegawaianJabatan
type UpdateKepegawaianJabatanRequest struct {
	DepartmentID     *int64 `json:"department_id" validate:"omitempty,gt=0"`
	PositionID       *int64 `json:"position_id" validate:"omitempty,gt=0"`
	JobTitleID       *int64 `json:"job_title_id" validate:"omitempty,gt=0"`
	SpecializationID *int64 `json:"specialization_id" validate:"omitempty,gt=0"`

	IsPrimary *bool `json:"is_primary"`

	TanggalMulai   *string `json:"tanggal_mulai" validate:"omitempty,datetime=2006-01-02"`
	TanggalSelesai *string `json:"tanggal_selesai" validate:"omitempty,datetime=2006-01-02"`

	NomorSK   *string `json:"nomor_sk" validate:"omitempty,max=100"`
	TanggalSK *string `json:"tanggal_sk" validate:"omitempty,datetime=2006-01-02"`

	IsAktif *bool `json:"is_aktif"`
}

// FilterKepegawaianJabatanRequest request body untuk filter KepegawaianJabatan
type FilterKepegawaianJabatanRequest struct {
	PegawaiID        *int64 `query:"pegawai_id"`
	DepartmentID     *int64 `query:"department_id"`
	PositionID       *int64 `query:"position_id"`
	JobTitleID       *int64 `query:"job_title_id"`
	SpecializationID *int64 `query:"specialization_id"`
	IsPrimary        *bool  `query:"is_primary"`
	IsAktif          *bool  `query:"is_aktif"`
}
