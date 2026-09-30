package dto

import (
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// KepegawaianJabatanResponse response untuk single KepegawaianJabatan
type KepegawaianJabatanResponse struct {
	ID int64 `json:"id"`

	PegawaiID    int64 `json:"pegawai_id"`
	DepartmentID int64 `json:"department_id"`

	PositionID       int64                         `json:"position_id"`
	Position         *PositionSimpelResponse       `json:"position"`
	JobTitleID       int64                         `json:"job_title_id"`
	JobTitle         *JobTitleSimpelResponse       `json:"job_title"`
	SpecializationID *int64                        `json:"specialization_id"`
	Specialization   *SpecializationSimpelResponse `json:"specialization"`

	IsPrimary bool `json:"is_primary"`

	TanggalMulai   *types.DateOnly `json:"tanggal_mulai"`
	TanggalSelesai *types.DateOnly `json:"tanggal_selesai"`

	NomorSK   *string         `json:"nomor_sk"`
	TanggalSK *types.DateOnly `json:"tanggal_sk"`

	IsAktif bool `json:"is_aktif"`

	CreatedBy *he.UserData     `json:"created_by"`
	UpdatedBy *he.UserData     `json:"updated_by"`
	CreatedAt types.CustomTime `json:"created_at"`
	UpdatedAt types.CustomTime `json:"updated_at"`
}

type KepegawaianJabatanResponseParams struct {
	KepegawaianJabatan *models.KepegawaianJabatan
	Creator            *he.UserData
	Updater            *he.UserData
}

// ToKepegawaianJabatanResponse mengubah model menjadi response
func ToKepegawaianJabatanResponse(params KepegawaianJabatanResponseParams) *KepegawaianJabatanResponse {
	m := params.KepegawaianJabatan

	var position *PositionSimpelResponse
	if m.Position != nil {
		position = &PositionSimpelResponse{
			ID:   m.Position.ID,
			Name: m.Position.Name,
		}
	}

	var jobTitle *JobTitleSimpelResponse
	if m.JobTitle != nil {
		jobTitle = &JobTitleSimpelResponse{
			ID:    m.JobTitle.ID,
			Code:  m.JobTitle.Code,
			Label: m.JobTitle.Label,
		}
	}

	var specialization *SpecializationSimpelResponse
	if m.Specialization != nil {
		specialization = &SpecializationSimpelResponse{
			ID:         m.Specialization.ID,
			Code:       m.Specialization.Code,
			Label:      m.Specialization.Label,
			Gelar:      m.Specialization.Gelar,
			FHIRCode:   m.Specialization.FHIRCode,
			FHIRSystem: m.Specialization.FHIRSystem,
		}
	}

	return &KepegawaianJabatanResponse{
		ID:               m.ID,
		PegawaiID:        m.PegawaiID,
		DepartmentID:     m.DepartmentID,
		PositionID:       m.PositionID,
		Position:         position,
		JobTitleID:       m.JobTitleID,
		JobTitle:         jobTitle,
		SpecializationID: m.SpecializationID,
		Specialization:   specialization,
		IsPrimary:        m.IsPrimary,
		TanggalMulai:     types.NewDateOnlyPtr(&m.TanggalMulai),
		TanggalSelesai:   types.NewDateOnlyPtr(m.TanggalSelesai),
		NomorSK:          m.NomorSK,
		TanggalSK:        types.NewDateOnlyPtr(m.TanggalSK),
		IsAktif:          m.IsAktif,
		CreatedBy:        params.Creator,
		UpdatedBy:        params.Updater,
		CreatedAt:        types.CustomTime(m.CreatedAt),
		UpdatedAt:        types.CustomTime(m.UpdatedAt),
	}
}

// ToKepegawaianJabatanListResponse mengubah slice model menjadi slice response
func ToKepegawaianJabatanListResponse(
	items []models.KepegawaianJabatan,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []KepegawaianJabatanResponse {
	responses := make([]KepegawaianJabatanResponse, 0, len(items))

	for i := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && items[i].CreatedBy != nil {
			creator = creatorsMap[*items[i].CreatedBy]
		}
		if updatersMap != nil && items[i].UpdatedBy != nil {
			updater = updatersMap[*items[i].UpdatedBy]
		}

		responses = append(responses, *ToKepegawaianJabatanResponse(KepegawaianJabatanResponseParams{
			KepegawaianJabatan: &items[i],
			Creator:            creator,
			Updater:            updater,
		}))
	}

	return responses
}
