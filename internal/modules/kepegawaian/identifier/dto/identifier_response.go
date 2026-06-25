// response.go
package dto

import (
	"neosim_go/internal/modules/kepegawaian/identifier/models"
	he "neosim_go/internal/shared/httputil"
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
	CreatedBy        *he.UserData      `json:"created_by"`
	UpdatedBy        *he.UserData      `json:"updated_by"`
	CreatedAt        types.CustomTime  `json:"created_at"`
	UpdatedAt        types.CustomTime  `json:"updated_at"`
}

type KepegawaianIdentifierResponseParams struct {
	Identifier *models.KepegawaianIdentifier
	Creator    *he.UserData
	Updater    *he.UserData
}

func ToKepegawaianIdentifierResponse(params KepegawaianIdentifierResponseParams) *KepegawaianIdentifierResponse {
	tipeLabel := string(params.Identifier.Tipe)
	if meta, ok := params.Identifier.Tipe.Meta(); ok {
		tipeLabel = meta.Label
	}

	return &KepegawaianIdentifierResponse{
		ID:               params.Identifier.ID,
		PegawaiID:        params.Identifier.PegawaiID,
		Tipe:             string(params.Identifier.Tipe),
		TipeLabel:        tipeLabel,
		Nilai:            params.Identifier.Nilai,
		Penerbit:         params.Identifier.Penerbit,
		TanggalTerbit:    types.ToCustomTimePtr(params.Identifier.TanggalTerbit),
		TanggalExpired:   types.ToCustomTimePtr(params.Identifier.TanggalExpired),
		IsPrimary:        params.Identifier.IsPrimary,
		IsAktif:          params.Identifier.IsAktif,
		IsFHIRMappable:   params.Identifier.IsFHIRMappable(),
		IsExpired:        params.Identifier.IsExpired(),
		DaysUntilExpired: params.Identifier.DaysUntilExpired(),
		CreatedBy:        params.Creator,
		UpdatedBy:        params.Updater,
		CreatedAt:        types.CustomTime(params.Identifier.CreatedAt),
		UpdatedAt:        types.CustomTime(params.Identifier.UpdatedAt),
	}
}

// ToKepegawaianIdentifierListResponse mengubah slice model menjadi slice response
// ToKepegawaianIdentifierListResponse mengubah slice model menjadi slice response
func ToKepegawaianIdentifierListResponse(
	items []models.KepegawaianIdentifier,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []KepegawaianIdentifierResponse {
	responses := make([]KepegawaianIdentifierResponse, 0, len(items))

	for _, m := range items {
		// Ambil data creator dan updater dari map berdasarkan ID.
		// CATATAN: Sesuaikan 'm.CreatedByID' dan 'm.UpdatedByID' dengan nama field foreign key di model Anda.
		// Jika di model Anda field ID-nya bernama 'CreatedBy' dan 'UpdatedBy', ubah menjadi m.CreatedBy dan m.UpdatedBy.
		var creator, updater *he.UserData

		if creatorsMap != nil {
			creator = creatorsMap[m.ID]
		}
		if updatersMap != nil {
			updater = updatersMap[m.ID]
		}

		responses = append(responses, *ToKepegawaianIdentifierResponse(KepegawaianIdentifierResponseParams{
			Identifier: &m,
			Creator:    creator,
			Updater:    updater,
		}))
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
