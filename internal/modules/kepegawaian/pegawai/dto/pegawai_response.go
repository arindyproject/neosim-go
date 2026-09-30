package dto

import (
	"neosim_go/internal/modules/kepegawaian/pegawai/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

// MasterRefResponse ringkasan data master (jenis kelamin, agama, dst.)
type MasterRefResponse struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	FhirCode *string `json:"fhir_code"`
}

func newMasterRef(id int64, name string, fhirCode *string) *MasterRefResponse {
	return &MasterRefResponse{ID: id, Name: name, FhirCode: fhirCode}
}

// KepegawaianPegawaiResponse response untuk single KepegawaianPegawai
type KepegawaianPegawaiResponse struct {
	ID     int64  `json:"id"`
	UserID *int64 `json:"user_id"`

	NIK          string  `json:"nik"`
	IHSNumber    *string `json:"ihs_number"`
	NomorPegawai string  `json:"nomor_pegawai"`
	NamaLengkap  string  `json:"nama_lengkap"`

	JenisKelamin     *MasterRefResponse `json:"jenis_kelamin"`
	TanggalLahir     *types.DateOnly    `json:"tanggal_lahir"`
	TempatLahir      string             `json:"tempat_lahir"`
	GolonganDarah    *MasterRefResponse `json:"golongan_darah"`
	Agama            *MasterRefResponse `json:"agama"`
	StatusPernikahan *MasterRefResponse `json:"status_pernikahan"`
	Kewarganegaraan  *string            `json:"kewarganegaraan"`

	TanggalMasuk  *types.DateOnly       `json:"tanggal_masuk"`
	TanggalKeluar *types.DateOnly       `json:"tanggal_keluar"`
	Jenis         *JenisSimpelResponse  `json:"jenis"`
	Status        *StatusSimpelResponse `json:"status"`
	FotoURL       *string               `json:"foto_url"`
	IsAktif       bool                  `json:"is_aktif"`

	CreatedBy *he.UserData     `json:"created_by"`
	UpdatedBy *he.UserData     `json:"updated_by"`
	CreatedAt types.CustomTime `json:"created_at"`
	UpdatedAt types.CustomTime `json:"updated_at"`
}

type KepegawaianPegawaiResponseParams struct {
	KepegawaianPegawai *models.KepegawaianPegawai
	Creator            *he.UserData
	Updater            *he.UserData
}

// ToKepegawaianPegawaiResponse mengubah model menjadi response
func ToKepegawaianPegawaiResponse(params KepegawaianPegawaiResponseParams) *KepegawaianPegawaiResponse {
	if params.KepegawaianPegawai == nil {
		return nil
	}
	m := params.KepegawaianPegawai

	var jenis *JenisSimpelResponse
	if m.Jenis != nil {
		jenis = &JenisSimpelResponse{ID: m.Jenis.ID, Code: m.Jenis.Code, Label: m.Jenis.Label, FHIRCode: m.Jenis.FHIRCode}
	}

	var status *StatusSimpelResponse
	if m.Status != nil {
		status = &StatusSimpelResponse{ID: m.Status.ID, Code: m.Status.Code, Label: m.Status.Label, FHIRCode: m.Status.FHIRCode}
	}

	var jenisKelamin, golDarah, agama, pernikahan *MasterRefResponse
	if m.JenisKelamin != nil {
		jenisKelamin = newMasterRef(m.JenisKelamin.ID, m.JenisKelamin.Name, m.JenisKelamin.FhirCode)
	}
	if m.GolonganDarah != nil {
		golDarah = newMasterRef(m.GolonganDarah.ID, m.GolonganDarah.Name, m.GolonganDarah.FhirCode)
	}
	if m.Agama != nil {
		agama = newMasterRef(m.Agama.ID, m.Agama.Name, m.Agama.FhirCode)
	}
	if m.StatusPernikahan != nil {
		pernikahan = newMasterRef(m.StatusPernikahan.ID, m.StatusPernikahan.Name, m.StatusPernikahan.FhirCode)
	}

	return &KepegawaianPegawaiResponse{
		ID:               m.ID,
		UserID:           m.UserID,
		NIK:              m.NIK,
		IHSNumber:        m.IHSNumber,
		NomorPegawai:     m.NomorPegawai,
		NamaLengkap:      m.NamaLengkap,
		JenisKelamin:     jenisKelamin,
		TanggalLahir:     types.NewDateOnlyPtr(&m.TanggalLahir),
		TempatLahir:      m.TempatLahir,
		GolonganDarah:    golDarah,
		Agama:            agama,
		StatusPernikahan: pernikahan,
		Kewarganegaraan:  m.Kewarganegaraan,
		TanggalMasuk:     types.NewDateOnlyPtr(&m.TanggalMasuk),
		TanggalKeluar:    types.NewDateOnlyPtr(m.TanggalKeluar),
		Jenis:            jenis,
		Status:           status,
		FotoURL:          m.FotoURL,
		IsAktif:          m.IsAktif,
		CreatedBy:        params.Creator,
		UpdatedBy:        params.Updater,
		CreatedAt:        types.CustomTime(m.CreatedAt),
		UpdatedAt:        types.CustomTime(m.UpdatedAt),
	}
}

// ToKepegawaianPegawaiListResponse mengubah slice model menjadi slice response
func ToKepegawaianPegawaiListResponse(
	items []models.KepegawaianPegawai,
	creatorsMap map[int64]*he.UserData,
	updatersMap map[int64]*he.UserData,
) []KepegawaianPegawaiResponse {
	responses := make([]KepegawaianPegawaiResponse, 0, len(items))

	for i := range items {
		m := &items[i]
		var creator, updater *he.UserData

		if creatorsMap != nil && m.CreatedBy != nil {
			creator = creatorsMap[*m.CreatedBy]
		}
		if updatersMap != nil && m.UpdatedBy != nil {
			updater = updatersMap[*m.UpdatedBy]
		}

		responses = append(responses, *ToKepegawaianPegawaiResponse(KepegawaianPegawaiResponseParams{
			KepegawaianPegawai: m,
			Creator:            creator,
			Updater:            updater,
		}))
	}

	return responses
}
