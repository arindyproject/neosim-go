package dto

import (
	"neosim_go/internal/modules/kepegawaian/pegawai/models"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/types"
)

const dateLayout = "2006-01-02"

// KepegawaianPegawaiResponse response untuk single KepegawaianPegawai
type KepegawaianPegawaiResponse struct {
	ID     int64  `json:"id"`
	UserID *int64 `json:"user_id"`

	NIK          string  `json:"nik"`
	IHSNumber    *string `json:"ihs_number"`
	NomorPegawai string  `json:"nomor_pegawai"`
	NamaLengkap  string  `json:"nama_lengkap"`

	JenisKelamin     string  `json:"jenis_kelamin"`
	TanggalLahir     string  `json:"tanggal_lahir"`
	TempatLahir      string  `json:"tempat_lahir"`
	GolonganDarah    *string `json:"golongan_darah"`
	Agama            string  `json:"agama"`
	StatusPerkawinan string  `json:"status_perkawinan"`
	Kewarganegaraan  string  `json:"kewarganegaraan"`

	TanggalMasuk  string                `json:"tanggal_masuk"`
	TanggalKeluar *string               `json:"tanggal_keluar"`
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

func toFHIRGender(jenisKelamin string) string {
	switch jenisKelamin {
	case "Laki-laki":
		return "male"
	case "Perempuan":
		return "female"
	default:
		return "unknown"
	}
}

// ToKepegawaianPegawaiResponse mengubah model menjadi response
func ToKepegawaianPegawaiResponse(params KepegawaianPegawaiResponseParams) *KepegawaianPegawaiResponse {
	m := params.KepegawaianPegawai

	var tglKeluar *string
	if m.TanggalKeluar != nil {
		s := m.TanggalKeluar.Format(dateLayout)
		tglKeluar = &s
	}

	var jenis *JenisSimpelResponse
	if m.Jenis != nil {
		jenis = &JenisSimpelResponse{ID: m.Jenis.ID, Code: m.Jenis.Code, Label: m.Jenis.Label, FHIRCode: m.Jenis.FHIRCode}
	}

	var status *StatusSimpelResponse
	if m.Status != nil {
		status = &StatusSimpelResponse{ID: m.Status.ID, Code: m.Status.Code, Label: m.Status.Label, FHIRCode: m.Status.FHIRCode}
	}

	return &KepegawaianPegawaiResponse{
		ID:               m.ID,
		UserID:           m.UserID,
		NIK:              m.NIK,
		IHSNumber:        m.IHSNumber,
		NomorPegawai:     m.NomorPegawai,
		NamaLengkap:      m.NamaLengkap,
		JenisKelamin:     m.JenisKelamin,
		TanggalLahir:     m.TanggalLahir.Format(dateLayout),
		TempatLahir:      m.TempatLahir,
		GolonganDarah:    m.GolonganDarah,
		Agama:            m.Agama,
		StatusPerkawinan: m.StatusPerkawinan,
		Kewarganegaraan:  m.Kewarganegaraan,
		TanggalMasuk:     m.TanggalMasuk.Format(dateLayout),
		TanggalKeluar:    tglKeluar,
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
