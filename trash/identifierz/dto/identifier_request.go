// request.go
package dto

import (
	"fmt"

	"neosim_go/internal/modules/kepegawaian/identifier/models"
)

// CreateKepegawaianIdentifierRequest request body untuk membuat KepegawaianIdentifier baru
type CreateKepegawaianIdentifierRequest struct {
	PegawaiID      int64                 `json:"pegawai_id"   validate:"required"`
	Tipe           models.IdentifierType `json:"tipe"             validate:"required"`
	Nilai          string                `json:"nilai"            validate:"required,min=1,max=100"`
	Penerbit       *string               `json:"penerbit"         validate:"omitempty,max=100"`
	TanggalTerbit  *string               `json:"tanggal_terbit"   validate:"omitempty"`
	TanggalExpired *string               `json:"tanggal_expired"  validate:"omitempty"`
	IsPrimary      bool                  `json:"is_primary"`
	IsAktif        bool                  `json:"is_aktif"`
}

// Validate validasi bisnis di luar tag validate
func (r *CreateKepegawaianIdentifierRequest) Validate() error {
	if !r.Tipe.IsValid() {
		return fmt.Errorf("tipe identifier '%s' tidak dikenal", r.Tipe)
	}
	if r.Tipe.HasExpiry() && (r.TanggalExpired == nil || *r.TanggalExpired == "") {
		return fmt.Errorf("tanggal_expired wajib diisi untuk tipe %s", r.Tipe)
	}
	return nil
}

// UpdateKepegawaianIdentifierRequest request body untuk update KepegawaianIdentifier
type UpdateKepegawaianIdentifierRequest struct {
	Tipe           *models.IdentifierType `json:"tipe"             validate:"omitempty"`
	Nilai          *string                `json:"nilai"            validate:"omitempty,min=1,max=100"`
	Penerbit       *string                `json:"penerbit"         validate:"omitempty,max=100"`
	TanggalTerbit  *string                `json:"tanggal_terbit"   validate:"omitempty"`
	TanggalExpired *string                `json:"tanggal_expired"  validate:"omitempty"`
	IsPrimary      *bool                  `json:"is_primary"       validate:"omitempty"`
	IsAktif        *bool                  `json:"is_aktif"         validate:"omitempty"`
}

// Validate validasi bisnis di luar tag validate
func (r *UpdateKepegawaianIdentifierRequest) Validate() error {
	if r.Tipe != nil {
		if !r.Tipe.IsValid() {
			return fmt.Errorf("tipe identifier '%s' tidak dikenal", *r.Tipe)
		}
		if r.Tipe.HasExpiry() && (r.TanggalExpired == nil || *r.TanggalExpired == "") {
			return fmt.Errorf("tanggal_expired wajib diisi untuk tipe %s", *r.Tipe)
		}
	}
	return nil
}

// FilterKepegawaianIdentifierRequest request body untuk filter KepegawaianIdentifier
type FilterKepegawaianIdentifierRequest struct {
	PegawaiID *int64                 `query:"pegawai_id"`
	Tipe      *models.IdentifierType `query:"tipe"`
	IsAktif   *bool                  `query:"is_aktif"`
	IsExpired *bool                  `query:"is_expired"` // filter identifier yg sudah expired
	IsPrimary *bool                  `query:"is_primary"`
}
