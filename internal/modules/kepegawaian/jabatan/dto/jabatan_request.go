package dto

import "neosim_go/internal/shared/types"

// CreateKepegawaianJabatanRequest request body untuk membuat KepegawaianJabatan baru
type CreateKepegawaianJabatanRequest struct {
	PegawaiID        int64  `json:"pegawai_id" validate:"required,gt=0"`
	DepartmentID     int64  `json:"department_id" validate:"required,gt=0"`
	PositionID       int64  `json:"position_id" validate:"required,gt=0"`
	JobTitleID       int64  `json:"job_title_id" validate:"required,gt=0"`
	SpecializationID *int64 `json:"specialization_id" validate:"omitempty,gt=0"`

	IsPrimary *bool `json:"is_primary"`

	// Format tanggal: YYYY-MM-DD
	TanggalMulai   types.DateOnly  `json:"tanggal_mulai" validate:"required" swaggertype:"string" format:"date" example:"1990-01-01"`
	TanggalSelesai *types.DateOnly `json:"tanggal_selesai" validate:"omitempty" swaggertype:"string" format:"date" example:"1990-01-01"`

	NomorSK   *string         `json:"nomor_sk" validate:"omitempty,max=100"`
	TanggalSK *types.DateOnly `json:"tanggal_sk" validate:"omitempty" swaggertype:"string" format:"date" example:"1990-01-01"`

	IsAktif *bool `json:"is_aktif"`
}

// UpdateKepegawaianJabatanRequest request body untuk update KepegawaianJabatan
type UpdateKepegawaianJabatanRequest struct {
	DepartmentID     *int64 `json:"department_id" validate:"omitempty,gt=0"`
	PositionID       *int64 `json:"position_id" validate:"omitempty,gt=0"`
	JobTitleID       *int64 `json:"job_title_id" validate:"omitempty,gt=0"`
	SpecializationID *int64 `json:"specialization_id" validate:"omitempty,gt=0"`

	IsPrimary *bool `json:"is_primary"`

	TanggalMulai   *types.DateOnly `json:"tanggal_mulai" validate:"required" swaggertype:"string" format:"date" example:"1990-01-01"`
	TanggalSelesai *types.DateOnly `json:"tanggal_selesai" validate:"omitempty" swaggertype:"string" format:"date" example:"1990-01-01"`

	NomorSK   *string         `json:"nomor_sk" validate:"omitempty,max=100"`
	TanggalSK *types.DateOnly `json:"tanggal_sk" validate:"omitempty" swaggertype:"string" format:"date" example:"1990-01-01"`

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

	// Sorting
	SortBy    string `query:"sort_by"`    // pegawai_id | department_id | position_id | job_title_id | specializations_id | is_primary | is_aktif | created_at | updated_at
	SortOrder string `query:"sort_order"` // asc | desc
}
