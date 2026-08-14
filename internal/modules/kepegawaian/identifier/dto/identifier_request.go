package dto

import "neosim_go/internal/shared/types"

// CreateKepegawaianIdentifierRequest request body untuk membuat KepegawaianIdentifier baru
type CreateKepegawaianIdentifierRequest struct {
	PegawaiID      int64           `json:"pegawai_id" validate:"required,gt=0"`
	TipeID         int64           `json:"tipe_id" validate:"required,gt=0"`
	Nilai          string          `json:"nilai" validate:"required,min=1,max=255"`
	TanggalTerbit  *types.DateOnly `json:"tanggal_terbit" swaggertype:"string" format:"date" example:"2026-01-01"`
	TanggalExpired *types.DateOnly `json:"tanggal_expired" swaggertype:"string" format:"date" example:"2026-01-01"`
	IsPrimary      bool            `json:"is_primary"`
	IsAktif        bool            `json:"is_aktif"`
}

// UpdateKepegawaianIdentifierRequest request body untuk update KepegawaianIdentifier
type UpdateKepegawaianIdentifierRequest struct {
	TipeID         *int64          `json:"tipe_id" validate:"omitempty,gt=0"`
	Nilai          *string         `json:"nilai" validate:"omitempty,min=1,max=255"`
	TanggalTerbit  *types.DateOnly `json:"tanggal_terbit" swaggertype:"string" format:"date" example:"2026-01-01"`
	TanggalExpired *types.DateOnly `json:"tanggal_expired" swaggertype:"string" format:"date" example:"2026-01-01"`
	IsPrimary      *bool           `json:"is_primary"`
	IsAktif        *bool           `json:"is_aktif"`
}

// FilterKepegawaianIdentifierRequest request query untuk filter KepegawaianIdentifier
type FilterKepegawaianIdentifierRequest struct {
	PegawaiID *int64 `query:"pegawai_id"`
	TipeID    *int64 `query:"tipe_id"`
	Nilai     string `query:"nilai"`
	IsPrimary *bool  `query:"is_primary"`
	IsAktif   *bool  `query:"is_aktif"`
}
