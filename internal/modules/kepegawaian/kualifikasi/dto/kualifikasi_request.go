package dto

import "neosim_go/internal/shared/types"

// CreateKepegawaianKualifikasiRequest request body untuk membuat KepegawaianKualifikasi baru
type CreateKepegawaianKualifikasiRequest struct {
	PegawaiID int64 `json:"pegawai_id" validate:"required,gt=0"`
	TipeID    int64 `json:"tipe_id" validate:"required,gt=0"`

	Nama            string  `json:"nama" validate:"required,min=1,max=255"`
	Penyelenggara   string  `json:"penyelenggara" validate:"required,min=1,max=255"`
	NomorSertifikat *string `json:"nomor_sertifikat" validate:"omitempty,min=1,max=255"`

	TanggalTerbit  *types.DateOnly `json:"tanggal_terbit" swaggertype:"string" format:"date" example:"2026-01-01"`
	TanggalExpired *types.DateOnly `json:"tanggal_expired" swaggertype:"string" format:"date" example:"2026-01-01"`

	IsAktif bool `json:"is_aktif"`

	FhirCode   *string `json:"fhir_code"`
	FhirSystem *string `json:"fhir_system"`
}

// UpdateKepegawaianKualifikasiRequest request body untuk update KepegawaianKualifikasi
type UpdateKepegawaianKualifikasiRequest struct {
	TipeID *int64 `json:"tipe_id" validate:"omitempty,gt=0"`

	Nama            *string `json:"nama" validate:"omitempty,min=1,max=255"`
	Penyelenggara   *string `json:"penyelenggara" validate:"omitempty,min=1,max=255"`
	NomorSertifikat *string `json:"nomor_sertifikat" validate:"omitempty,min=1,max=255"`

	TanggalTerbit  *types.DateOnly `json:"tanggal_terbit" swaggertype:"string" format:"date" example:"2026-01-01"`
	TanggalExpired *types.DateOnly `json:"tanggal_expired" swaggertype:"string" format:"date" example:"2026-01-01"`

	IsAktif *bool `json:"is_aktif"`

	FhirCode   *string `json:"fhir_code"`
	FhirSystem *string `json:"fhir_system"`
}

// FilterKepegawaianKualifikasiRequest request body untuk filter KepegawaianKualifikasi
type FilterKepegawaianKualifikasiRequest struct {
	TipeID        int64  `query:"tipe_id"`
	Nama          string `query:"nama"`
	Penyelenggara string `query:"penyelenggara"`
	IsAktif       *bool  `query:"is_aktif"`
	IsExpired     *bool  `query:"is_expired"`
}
