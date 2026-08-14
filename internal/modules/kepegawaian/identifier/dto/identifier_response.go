package dto

import (
	"neosim_go/internal/modules/kepegawaian/identifier/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// KepegawaianIdentifierResponse response untuk single KepegawaianIdentifier
type KepegawaianIdentifierResponse struct {
	ID        int64 `json:"id"`
	PegawaiID int64 `json:"pegawai_id"`
	//TipeID         int64            `json:"tipe_id"`
	Tipe           *TipeSimpelResponse `json:"tipe,omitempty"`
	Nilai          string              `json:"nilai"`
	Penerbit       *string             `json:"penerbit"`
	TanggalTerbit  *types.DateOnly     `json:"tanggal_terbit"`
	TanggalExpired *types.DateOnly     `json:"tanggal_expired"`
	IsPrimary      bool                `json:"is_primary"`
	IsAktif        bool                `json:"is_aktif"`
	IsExpired      bool                `json:"is_expired"`
	DaysUntilExp   int                 `json:"days_until_expired"`
	IsFHIRMappable bool                `json:"is_fhir_mappable"`
	CreatedBy      *he.UserData        `json:"created_by"`
	UpdatedBy      *he.UserData        `json:"updated_by"`
	CreatedAt      types.CustomTime    `json:"created_at"`
	UpdatedAt      types.CustomTime    `json:"updated_at"`
}

type KepegawaianIdentifierResponseParams struct {
	KepegawaianIdentifier *models.KepegawaianIdentifier
	Creator               *he.UserData
	Updater               *he.UserData
}

// ToKepegawaianIdentifierResponse mengubah model menjadi response
func ToKepegawaianIdentifierResponse(params KepegawaianIdentifierResponseParams) *KepegawaianIdentifierResponse {
	if params.KepegawaianIdentifier == nil {
		return nil
	}

	m := params.KepegawaianIdentifier

	var tipeResponse *TipeSimpelResponse
	if m.Tipe != nil {
		tipeResponse = &TipeSimpelResponse{
			ID:          m.Tipe.ID,
			Code:        m.Tipe.Code,
			Label:       m.Tipe.Label,
			FHIRSystem:  m.Tipe.FHIRSystem,
			HasExpiry:   m.Tipe.HasExpiry,
			IsNakes:     m.Tipe.IsNakes,
			IsRequired:  m.Tipe.IsRequired,
			Description: m.Tipe.Description,
		}
	}

	return &KepegawaianIdentifierResponse{
		ID:        m.ID,
		PegawaiID: m.PegawaiID,
		//TipeID:         m.TipeID,
		Tipe:           tipeResponse,
		Nilai:          m.Nilai,
		TanggalTerbit:  types.NewDateOnlyPtr(m.TanggalTerbit),
		TanggalExpired: types.NewDateOnlyPtr(m.TanggalExpired),
		IsPrimary:      m.IsPrimary,
		IsAktif:        m.IsAktif,
		IsExpired:      m.IsExpired(),
		DaysUntilExp:   m.DaysUntilExpired(),
		IsFHIRMappable: m.IsFHIRMappable(),
		CreatedBy:      params.Creator,
		UpdatedBy:      params.Updater,
		CreatedAt:      types.CustomTime(m.CreatedAt),
		UpdatedAt:      types.CustomTime(m.UpdatedAt),
	}
}

// ToKepegawaianIdentifierListResponse mengubah slice model menjadi slice response
func ToKepegawaianIdentifierListResponse(
	items []models.KepegawaianIdentifier,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []KepegawaianIdentifierResponse {
	responses := make([]KepegawaianIdentifierResponse, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		responses = append(responses, *ToKepegawaianIdentifierResponse(KepegawaianIdentifierResponseParams{
			KepegawaianIdentifier: &m,
			Creator:               creator,
			Updater:               updater,
		}))
	}

	return responses
}
