// response.go
package dto

import (
	"neosim_go/internal/modules/kepegawaian/identifier/models"
	"neosim_go/internal/shared/types"
)

// IdentifierMetaResponse response metadata tipe identifier untuk dropdown UI
type IdentifierMetaResponse struct {
	Code       string `json:"code"`
	Label      string `json:"label"`
	Penerbit   string `json:"penerbit"`
	FHIRSystem string `json:"fhir_system"`
	HasExpiry  bool   `json:"has_expiry"`
	IsNakes    bool   `json:"is_nakes"`
	IsRequired bool   `json:"is_required"`
}

// KepegawaianIdentifierResponse response untuk single KepegawaianIdentifier
type KepegawaianIdentifierResponse struct {
	ID               int64             `json:"id"`
	PegawaiID        int64             `json:"pegawai_id"`
	Tipe             string            `json:"tipe"`
	TipeLabel        string            `json:"tipe_label"` // label human-readable, ex: "Surat Tanda Registrasi"
	Nilai            string            `json:"nilai"`
	Penerbit         *string           `json:"penerbit"`
	TanggalTerbit    *types.CustomTime `json:"tanggal_terbit"`
	TanggalExpired   *types.CustomTime `json:"tanggal_expired"`
	IsPrimary        bool              `json:"is_primary"`
	IsAktif          bool              `json:"is_aktif"`
	IsFHIRMappable   bool              `json:"is_fhir_mappable"`   // bisa disync ke SATUSEHAT?
	IsExpired        bool              `json:"is_expired"`         // sudah lewat tanggal expired?
	DaysUntilExpired int               `json:"days_until_expired"` // -1 jika tidak ada expired date
	CreatedBy        *int64            `json:"created_by"`
	UpdatedBy        *int64            `json:"updated_by"`
	CreatedAt        types.CustomTime  `json:"created_at"`
	UpdatedAt        types.CustomTime  `json:"updated_at"`
}

// ToKepegawaianIdentifierResponse mengubah model menjadi response
func ToKepegawaianIdentifierResponse(m *models.KepegawaianIdentifier) *KepegawaianIdentifierResponse {
	tipeLabel := string(m.Tipe)
	if meta, ok := m.Tipe.Meta(); ok {
		tipeLabel = meta.Label
	}

	return &KepegawaianIdentifierResponse{
		ID:               m.ID,
		PegawaiID:        m.PegawaiID,
		Tipe:             string(m.Tipe),
		TipeLabel:        tipeLabel,
		Nilai:            m.Nilai,
		Penerbit:         m.Penerbit,
		TanggalTerbit:    types.ToCustomTimePtr(m.TanggalTerbit),
		TanggalExpired:   types.ToCustomTimePtr(m.TanggalExpired),
		IsPrimary:        m.IsPrimary,
		IsAktif:          m.IsAktif,
		IsFHIRMappable:   m.IsFHIRMappable(),
		IsExpired:        m.IsExpired(),
		DaysUntilExpired: m.DaysUntilExpired(),
		CreatedBy:        m.CreatedBy,
		UpdatedBy:        m.UpdatedBy,
		CreatedAt:        types.CustomTime(m.CreatedAt),
		UpdatedAt:        types.CustomTime(m.UpdatedAt),
	}
}

// ToKepegawaianIdentifierListResponse mengubah slice model menjadi slice response
func ToKepegawaianIdentifierListResponse(items []models.KepegawaianIdentifier) []KepegawaianIdentifierResponse {
	responses := make([]KepegawaianIdentifierResponse, 0, len(items))
	for _, m := range items {
		responses = append(responses, *ToKepegawaianIdentifierResponse(&m))
	}
	return responses
}

// ToIdentifierMetaResponse mengubah IdentifierMeta menjadi response untuk dropdown
func ToIdentifierMetaResponse(m models.IdentifierMeta) IdentifierMetaResponse {
	return IdentifierMetaResponse{
		Code:       string(m.Code),
		Label:      m.Label,
		Penerbit:   m.Penerbit,
		FHIRSystem: m.FHIRSystem,
		HasExpiry:  m.HasExpiry,
		IsNakes:    m.IsNakes,
		IsRequired: m.IsRequired,
	}
}

// ToIdentifierMetaListResponse mengubah slice IdentifierMeta menjadi slice response
func ToIdentifierMetaListResponse(items []models.IdentifierMeta) []IdentifierMetaResponse {
	responses := make([]IdentifierMetaResponse, 0, len(items))
	for _, m := range items {
		responses = append(responses, ToIdentifierMetaResponse(m))
	}
	return responses
}
