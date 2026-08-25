package dto

import "neosim_go/internal/shared/types"

// CreateKepegawaianPendidikanRequest request body untuk membuat KepegawaianPendidikan baru
type CreateKepegawaianPendidikanRequest struct {
	PegawaiID       int64           `json:"pegawai_id" validate:"required"`
	JenjangID       int64           `json:"jenjang_id" validate:"required"`
	NamaInstitusi   string          `json:"nama_institusi" validate:"required,min=1,max=225"`
	NomorIjazah     *string         `json:"nomor_ijazah" validate:"omitempty,max=255"`
	BidangStudi     *string         `json:"bidang_studi" validate:"omitempty,max=225"`
	AlamatInstitusi *string         `json:"alamat_institusi" validate:"omitempty"`
	NilaiAkhir      *string         `json:"nilai_akhir" validate:"omitempty,max=225"`
	TanggalMasuk    *types.DateOnly `json:"tanggal_masuk" swaggertype:"string" format:"date" example:"2026-01-01"`
	TanggalLulus    *types.DateOnly `json:"tanggal_lulus" swaggertype:"string" format:"date" example:"2026-01-01"`
	FHIRCode        *string         `json:"fhir_code" validate:"omitempty,max=255"`
	FHIRSystem      *string         `json:"fhir_system" validate:"omitempty,max=255"`
}

// UpdateKepegawaianPendidikanRequest request body untuk update KepegawaianPendidikan
type UpdateKepegawaianPendidikanRequest struct {
	JenjangID       *int64          `json:"jenjang_id" validate:"omitempty"`
	NamaInstitusi   *string         `json:"nama_institusi" validate:"omitempty,min=1,max=225"`
	NomorIjazah     *string         `json:"nomor_ijazah" validate:"omitempty,max=255"`
	BidangStudi     *string         `json:"bidang_studi" validate:"omitempty,max=225"`
	AlamatInstitusi *string         `json:"alamat_institusi" validate:"omitempty"`
	NilaiAkhir      *string         `json:"nilai_akhir" validate:"omitempty,max=225"`
	TanggalMasuk    *types.DateOnly `json:"tanggal_masuk" swaggertype:"string" format:"date" example:"2026-01-01"`
	TanggalLulus    *types.DateOnly `json:"tanggal_lulus" swaggertype:"string" format:"date" example:"2026-01-01"`
	FHIRCode        *string         `json:"fhir_code" validate:"omitempty,max=255"`
	FHIRSystem      *string         `json:"fhir_system" validate:"omitempty,max=255"`
}

// FilterKepegawaianPendidikanRequest request body untuk filter KepegawaianPendidikan
type FilterKepegawaianPendidikanRequest struct {
	PegawaiID       *int64 `query:"pegawai_id"`
	JenjangID       *int64 `query:"jenjang_id"`
	NamaInstitusi   string `query:"nama_institusi"`
	NomorIjazah     string `query:"nomor_ijazah"`
	BidangStudi     string `query:"bidang_studi"`
	AlamatInstitusi string `query:"alamat_institusi"`
}
