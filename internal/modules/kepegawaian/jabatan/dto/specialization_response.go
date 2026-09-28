package dto

import (
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// SpecializationResponse response untuk single Specialization
type SpecializationResponse struct {
	ID                  int64                                 `json:"id"`
	Code                string                                `json:"code"`
	Label               string                                `json:"label"`
	JobTitle            *JobTitleSimpelResponse               `json:"job_title"`
	Kategori            *SpecializationKategoriSimpelResponse `json:"kategori"`
	Gelar               *string                               `json:"gelar"`
	LamaPendidikanTahun *int16                                `json:"lama_pendidikan_tahun"`
	FHIRCode            *string                               `json:"fhir_code"`
	FHIRSystem          *string                               `json:"fhir_system"`
	IsAktif             bool                                  `json:"is_aktif"`
	CreatedBy           *he.UserData                          `json:"created_by"`
	UpdatedBy           *he.UserData                          `json:"updated_by"`
	CreatedAt           types.CustomTime                      `json:"created_at"`
	UpdatedAt           types.CustomTime                      `json:"updated_at"`
}

// SpecializationSimpelResponse response ringkas, untuk dipakai sebagai relasi di DTO lain
type SpecializationSimpelResponse struct {
	ID         int64   `json:"id"`
	Code       string  `json:"code"`
	Label      string  `json:"label"`
	Gelar      *string `json:"gelar"`
	FHIRCode   *string `json:"fhir_code"`
	FHIRSystem *string `json:"fhir_system"`
}

// SpecializationSelectResponse response untuk dropdown/select
type SpecializationSelectResponse struct {
	ID    int64   `json:"id"`
	Code  string  `json:"code"`
	Label string  `json:"label"`
	Gelar *string `json:"gelar"`
}

type SpecializationResponseParams struct {
	Specialization *models.Specialization
	Creator        *he.UserData
	Updater        *he.UserData
}

// ToSpecializationResponse mengubah model menjadi response
func ToSpecializationResponse(params SpecializationResponseParams) *SpecializationResponse {
	m := params.Specialization

	var jobTitle *JobTitleSimpelResponse
	if m.JobTitle != nil {
		jobTitle = &JobTitleSimpelResponse{
			ID:    m.JobTitle.ID,
			Code:  m.JobTitle.Code,
			Label: m.JobTitle.Label,
		}
	}

	var kategori *SpecializationKategoriSimpelResponse
	if m.Kategori != nil {
		kategori = &SpecializationKategoriSimpelResponse{
			ID:       m.Kategori.ID,
			Code:     m.Kategori.Code,
			Label:    m.Kategori.Label,
			FHIRCode: m.Kategori.FHIRCode,
		}
	}

	return &SpecializationResponse{
		ID:                  m.ID,
		Code:                m.Code,
		Label:               m.Label,
		JobTitle:            jobTitle,
		Kategori:            kategori,
		Gelar:               m.Gelar,
		LamaPendidikanTahun: m.LamaPendidikanTahun,
		FHIRCode:            m.FHIRCode,
		FHIRSystem:          m.FHIRSystem,
		IsAktif:             m.IsAktif,
		CreatedBy:           params.Creator,
		UpdatedBy:           params.Updater,
		CreatedAt:           types.CustomTime(m.CreatedAt),
		UpdatedAt:           types.CustomTime(m.UpdatedAt),
	}
}

// ToSpecializationListResponse mengubah slice model menjadi slice response
func ToSpecializationListResponse(
	items []models.Specialization,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []SpecializationResponse {
	responses := make([]SpecializationResponse, 0, len(items))

	for i := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && items[i].CreatedBy != nil {
			creator = creatorsMap[*items[i].CreatedBy]
		}
		if updatersMap != nil && items[i].UpdatedBy != nil {
			updater = updatersMap[*items[i].UpdatedBy]
		}

		responses = append(responses, *ToSpecializationResponse(SpecializationResponseParams{
			Specialization: &items[i],
			Creator:        creator,
			Updater:        updater,
		}))
	}

	return responses
}

// ToSpecializationSelectResponse mengubah slice model menjadi response untuk select
func ToSpecializationSelectResponse(items []models.Specialization) []SpecializationSelectResponse {
	responses := make([]SpecializationSelectResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, SpecializationSelectResponse{
			ID:    item.ID,
			Code:  item.Code,
			Label: item.Label,
			Gelar: item.Gelar,
		})
	}
	return responses
}
