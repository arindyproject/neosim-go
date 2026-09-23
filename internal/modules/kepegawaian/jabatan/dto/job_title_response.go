package dto

import (
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// JobTitleResponse response untuk single JobTitle
type JobTitleResponse struct {
	ID          int64   `json:"id"`
	Code        string  `json:"code"`
	Label       string  `json:"label"`
	Description *string `json:"description"`

	//KategoriID int64                           `json:"kategori_id"`
	Kategori *JobTitleKategoriSimpelResponse `json:"kategori,omitempty"`

	//RumpunProfesiID *int64                               `json:"rumpun_profesi_id"`
	RumpunProfesi *JobTitleRumpunProfesiSimpelResponse `json:"rumpun_profesi,omitempty"`

	Point *float64 `json:"point"`

	MemerlukanSTR bool `json:"memerlukan_str"`
	MemerlukanSIP bool `json:"memerlukan_sip"`

	JenjangMin *string `json:"jenjang_min"`

	FHIRCode   *string `json:"fhir_code"`
	FHIRSystem *string `json:"fhir_system"`

	IsAktif bool `json:"is_aktif"`

	CreatedBy *he.UserData     `json:"created_by"`
	UpdatedBy *he.UserData     `json:"updated_by"`
	CreatedAt types.CustomTime `json:"created_at"`
	UpdatedAt types.CustomTime `json:"updated_at"`
}

type JobTitleSimpelResponse struct {
	ID    int64  `json:"id"`
	Code  string `json:"code"`
	Label string `json:"label"`
}

type JobTitleResponseParams struct {
	JobTitle *models.JobTitle
	Creator  *he.UserData
	Updater  *he.UserData
}

func ToJobTitleSimpelResponse(items []models.JobTitle) []JobTitleSimpelResponse {
	responses := make([]JobTitleSimpelResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, JobTitleSimpelResponse{
			ID:    item.ID,
			Code:  item.Code,
			Label: item.Label,
		})
	}
	return responses
}

// ToJobTitleResponse mengubah model menjadi response
func ToJobTitleResponse(params JobTitleResponseParams) *JobTitleResponse {
	if params.JobTitle == nil {
		return nil
	}

	m := params.JobTitle

	var kategoriResponse *JobTitleKategoriSimpelResponse
	if m.Kategori != nil {
		kategoriResponse = &JobTitleKategoriSimpelResponse{
			ID:    m.Kategori.ID,
			Code:  m.Kategori.Code,
			Label: m.Kategori.Label,
		}
	}

	var rumpunResponse *JobTitleRumpunProfesiSimpelResponse
	if m.RumpunProfesi != nil {
		rumpunResponse = &JobTitleRumpunProfesiSimpelResponse{
			ID:    m.RumpunProfesi.ID,
			Code:  m.RumpunProfesi.Code,
			Label: m.RumpunProfesi.Label,
		}
	}

	return &JobTitleResponse{
		ID:            m.ID,
		Code:          m.Code,
		Label:         m.Label,
		Description:   m.Description,
		Kategori:      kategoriResponse,
		RumpunProfesi: rumpunResponse,
		Point:         m.Point,
		MemerlukanSTR: m.MemerlukanSTR,
		MemerlukanSIP: m.MemerlukanSIP,
		JenjangMin:    m.JenjangMin,
		FHIRCode:      m.FHIRCode,
		FHIRSystem:    m.FHIRSystem,
		IsAktif:       m.IsAktif,
		CreatedBy:     params.Creator,
		UpdatedBy:     params.Updater,
		CreatedAt:     types.CustomTime(m.CreatedAt),
		UpdatedAt:     types.CustomTime(m.UpdatedAt),
	}
}

// ToJobTitleListResponse mengubah slice model menjadi slice response
func ToJobTitleListResponse(
	items []models.JobTitle,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []JobTitleResponse {
	responses := make([]JobTitleResponse, 0, len(items))

	for _, m := range items {
		var creator, updater *he.UserData

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		responses = append(responses, *ToJobTitleResponse(JobTitleResponseParams{
			JobTitle: &m,
			Creator:  creator,
			Updater:  updater,
		}))
	}

	return responses
}
