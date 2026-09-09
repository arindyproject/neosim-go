package dto

import (
	"neosim_go/internal/modules/master/alamat/models"
	"neosim_go/internal/shared/types"
)

// ProvinsiDetailResponse response detail provinsi dengan statistik turunan
// -------------------------------------------------------------------------
type ProvinsiDetailResponse struct {
	ID             int64   `json:"id"`
	Code           string  `json:"code"`
	Name           string  `json:"name"`
	FhirCode       *string `json:"fhir_code"`
	NegaraID       int64   `json:"negara_id"`
	NegaraName     string  `json:"negara_name"`
	TotalKota      int64   `json:"total_kota"`
	TotalKecamatan int64   `json:"total_kecamatan"`
	TotalDesa      int64   `json:"total_desa"`
} //------------------------------------------------------------------------

// KotaKabupatenDetailResponse response detail kota/kabupaten dengan statistik turunan
// -------------------------------------------------------------------------
type KotaKabupatenDetailResponse struct {
	ID             int64   `json:"id"`
	Code           string  `json:"code"`
	Name           string  `json:"name"`
	FhirCode       *string `json:"fhir_code"`
	ProvinsiID     int64   `json:"provinsi_id"`
	ProvinsiName   string  `json:"provinsi_name"`
	NegaraID       int64   `json:"negara_id"`
	NegaraName     string  `json:"negara_name"`
	TotalKecamatan int64   `json:"total_kecamatan"`
	TotalDesa      int64   `json:"total_desa"`
} //------------------------------------------------------------------------

// KecamatanDetailResponse response detail kecamatan dengan statistik turunan
// -------------------------------------------------------------------------
type KecamatanDetailResponse struct {
	ID                int64   `json:"id"`
	Code              string  `json:"code"`
	Name              string  `json:"name"`
	FhirCode          *string `json:"fhir_code"`
	KotaKabupatenID   int64   `json:"kota_kabupaten_id"`
	KotaKabupatenName string  `json:"kota_kabupaten_name"`
	ProvinsiID        int64   `json:"provinsi_id"`
	ProvinsiName      string  `json:"provinsi_name"`
	NegaraID          int64   `json:"negara_id"`
	NegaraName        string  `json:"negara_name"`
	TotalDesa         int64   `json:"total_desa"`
} //------------------------------------------------------------------------

// KelurahanDesaDetailResponse response detail desa/kelurahan dengan jalur hierarki lengkap
// -------------------------------------------------------------------------
type KelurahanDesaDetailResponse struct {
	ID                int64   `json:"id"`
	Code              string  `json:"code"`
	Name              string  `json:"name"`
	FhirCode          *string `json:"fhir_code"`
	PostalCode        *string `json:"postal_code"`
	KecamatanID       int64   `json:"kecamatan_id"`
	KecamatanName     string  `json:"kecamatan_name"`
	KotaKabupatenID   int64   `json:"kota_kabupaten_id"`
	KotaKabupatenName string  `json:"kota_kabupaten_name"`
	ProvinsiID        int64   `json:"provinsi_id"`
	ProvinsiName      string  `json:"provinsi_name"`
	NegaraID          int64   `json:"negara_id"`
	NegaraName        string  `json:"negara_name"`
} //------------------------------------------------------------------------

// =====================================================================
// RESPONSE DASAR (untuk list & get biasa, tanpa statistik turunan)
// =====================================================================

// NegaraResponse response untuk single Negara
type NegaraResponse struct {
	ID          int64            `json:"id"`
	Code        string           `json:"code"`
	Name        string           `json:"name"`
	FhirCode    *string          `json:"fhir_code"`
	Description *string          `json:"description"`
	CreatedBy   *int64           `json:"created_by"`
	UpdatedBy   *int64           `json:"updated_by"`
	CreatedAt   types.CustomTime `json:"created_at"`
	UpdatedAt   types.CustomTime `json:"updated_at"`
}

type NegaraSimpelResponse struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// ProvinsiResponse response untuk single Provinsi
type ProvinsiResponse struct {
	ID        int64            `json:"id"`
	NegaraID  int64            `json:"negara_id"`
	Code      string           `json:"code"`
	Name      string           `json:"name"`
	FhirCode  *string          `json:"fhir_code"`
	CreatedBy *int64           `json:"created_by"`
	UpdatedBy *int64           `json:"updated_by"`
	CreatedAt types.CustomTime `json:"created_at"`
	UpdatedAt types.CustomTime `json:"updated_at"`
}

type ProvinsiSimpelResponse struct {
	ID       int64  `json:"id"`
	NegaraID int64  `json:"negara_id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
}

// KotaKabupatenResponse response untuk single Kota/Kabupaten
type KotaKabupatenResponse struct {
	ID         int64            `json:"id"`
	ProvinsiID int64            `json:"provinsi_id"`
	Code       string           `json:"code"`
	Name       string           `json:"name"`
	FhirCode   *string          `json:"fhir_code"`
	CreatedBy  *int64           `json:"created_by"`
	UpdatedBy  *int64           `json:"updated_by"`
	CreatedAt  types.CustomTime `json:"created_at"`
	UpdatedAt  types.CustomTime `json:"updated_at"`
}

type KotaKabupatenSimpelResponse struct {
	ID         int64  `json:"id"`
	ProvinsiID int64  `json:"provinsi_id"`
	Code       string `json:"code"`
	Name       string `json:"name"`
}

// KecamatanResponse response untuk single Kecamatan
type KecamatanResponse struct {
	ID              int64            `json:"id"`
	KotaKabupatenID int64            `json:"kota_kabupaten_id"`
	Code            string           `json:"code"`
	Name            string           `json:"name"`
	FhirCode        *string          `json:"fhir_code"`
	CreatedBy       *int64           `json:"created_by"`
	UpdatedBy       *int64           `json:"updated_by"`
	CreatedAt       types.CustomTime `json:"created_at"`
	UpdatedAt       types.CustomTime `json:"updated_at"`
}

type KecamatanSimpelResponse struct {
	ID              int64  `json:"id"`
	KotaKabupatenID int64  `json:"kota_kabupaten_id"`
	Code            string `json:"code"`
	Name            string `json:"name"`
}

// KelurahanDesaResponse response untuk single Kelurahan/Desa
type KelurahanDesaResponse struct {
	ID          int64            `json:"id"`
	KecamatanID int64            `json:"kecamatan_id"`
	Code        string           `json:"code"`
	Name        string           `json:"name"`
	FhirCode    *string          `json:"fhir_code"`
	PostalCode  *string          `json:"postal_code"`
	CreatedBy   *int64           `json:"created_by"`
	UpdatedBy   *int64           `json:"updated_by"`
	CreatedAt   types.CustomTime `json:"created_at"`
	UpdatedAt   types.CustomTime `json:"updated_at"`
}

type KelurahanDesaSimpelResponse struct {
	ID          int64   `json:"id"`
	KecamatanID int64   `json:"kecamatan_id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	PostalCode  *string `json:"postal_code"`
}

// =====================================================================
// NEGARA
// =====================================================================

// ToNegaraResponse mengubah model Negara menjadi response
func ToNegaraResponse(m *models.MasterAlamatNegara) *NegaraResponse {
	if m == nil {
		return nil
	}
	return &NegaraResponse{
		ID:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		FhirCode:    m.FhirCode,
		Description: m.Description,
		CreatedBy:   m.CreatedBy,
		UpdatedBy:   m.UpdatedBy,
		CreatedAt:   types.CustomTime(m.CreatedAt),
		UpdatedAt:   types.CustomTime(m.UpdatedAt),
	}
}

func ToNegaraListSimpelResponse(items []models.MasterAlamatNegara) []NegaraSimpelResponse {
	result := make([]NegaraSimpelResponse, 0, len(items))
	for _, item := range items {
		result = append(result, NegaraSimpelResponse{
			ID:   item.ID,
			Code: item.Code,
			Name: item.Name,
		})
	}
	return result
}

// ToNegaraListResponse mengubah slice model Negara menjadi slice response
func ToNegaraListResponse(items []models.MasterAlamatNegara) []NegaraResponse {
	result := make([]NegaraResponse, 0, len(items))
	for _, item := range items {
		result = append(result, *ToNegaraResponse(&item))
	}
	return result
}

// =====================================================================
// PROVINSI
// =====================================================================

// ToProvinsiResponse mengubah model Provinsi menjadi response
func ToProvinsiResponse(m *models.MasterAlamatProvinsi) *ProvinsiResponse {
	if m == nil {
		return nil
	}
	return &ProvinsiResponse{
		ID:        m.ID,
		NegaraID:  m.NegaraID,
		Code:      m.Code,
		Name:      m.Name,
		FhirCode:  m.FhirCode,
		CreatedBy: m.CreatedBy,
		UpdatedBy: m.UpdatedBy,
		CreatedAt: types.CustomTime(m.CreatedAt),
		UpdatedAt: types.CustomTime(m.UpdatedAt),
	}
}

func ToProvinsiListSimpelResponse(items []models.MasterAlamatProvinsi) []ProvinsiSimpelResponse {
	result := make([]ProvinsiSimpelResponse, 0, len(items))
	for _, item := range items {
		result = append(result, ProvinsiSimpelResponse{
			ID:       item.ID,
			NegaraID: item.NegaraID,
			Code:     item.Code,
			Name:     item.Name,
		})
	}
	return result
}

// ToProvinsiListResponse mengubah slice model Provinsi menjadi slice response
func ToProvinsiListResponse(items []models.MasterAlamatProvinsi) []ProvinsiResponse {
	result := make([]ProvinsiResponse, 0, len(items))
	for _, item := range items {
		result = append(result, *ToProvinsiResponse(&item))
	}
	return result
}

// =====================================================================
// KOTA / KABUPATEN
// =====================================================================

// ToKotaKabupatenResponse mengubah model Kota/Kabupaten menjadi response
func ToKotaKabupatenResponse(m *models.MasterAlamatKotaKabupaten) *KotaKabupatenResponse {
	if m == nil {
		return nil
	}
	return &KotaKabupatenResponse{
		ID:         m.ID,
		ProvinsiID: m.ProvinsiID,
		Code:       m.Code,
		Name:       m.Name,
		FhirCode:   m.FhirCode,
		CreatedBy:  m.CreatedBy,
		UpdatedBy:  m.UpdatedBy,
		CreatedAt:  types.CustomTime(m.CreatedAt),
		UpdatedAt:  types.CustomTime(m.UpdatedAt),
	}
}

func ToKotaKabupatenListSimpelResponse(items []models.MasterAlamatKotaKabupaten) []KotaKabupatenSimpelResponse {
	result := make([]KotaKabupatenSimpelResponse, 0, len(items))
	for _, item := range items {
		result = append(result, KotaKabupatenSimpelResponse{
			ID:         item.ID,
			ProvinsiID: item.ProvinsiID,
			Code:       item.Code,
			Name:       item.Name,
		})
	}
	return result
}

// ToKotaKabupatenListResponse mengubah slice model Kota/Kabupaten menjadi slice response
func ToKotaKabupatenListResponse(items []models.MasterAlamatKotaKabupaten) []KotaKabupatenResponse {
	result := make([]KotaKabupatenResponse, 0, len(items))
	for _, item := range items {
		result = append(result, *ToKotaKabupatenResponse(&item))
	}
	return result
}

// =====================================================================
// KECAMATAN
// =====================================================================

// ToKecamatanResponse mengubah model Kecamatan menjadi response
func ToKecamatanResponse(m *models.MasterAlamatKecamatan) *KecamatanResponse {
	if m == nil {
		return nil
	}
	return &KecamatanResponse{
		ID:              m.ID,
		KotaKabupatenID: m.KotaKabupatenID,
		Code:            m.Code,
		Name:            m.Name,
		FhirCode:        m.FhirCode,
		CreatedBy:       m.CreatedBy,
		UpdatedBy:       m.UpdatedBy,
		CreatedAt:       types.CustomTime(m.CreatedAt),
		UpdatedAt:       types.CustomTime(m.UpdatedAt),
	}
}

func ToKecamatanListSimpelResponse(items []models.MasterAlamatKecamatan) []KecamatanSimpelResponse {
	result := make([]KecamatanSimpelResponse, 0, len(items))
	for _, item := range items {
		result = append(result, KecamatanSimpelResponse{
			ID:              item.ID,
			KotaKabupatenID: item.KotaKabupatenID,
			Code:            item.Code,
			Name:            item.Name,
		})
	}
	return result
}

// ToKecamatanListResponse mengubah slice model Kecamatan menjadi slice response
func ToKecamatanListResponse(items []models.MasterAlamatKecamatan) []KecamatanResponse {
	result := make([]KecamatanResponse, 0, len(items))
	for _, item := range items {
		result = append(result, *ToKecamatanResponse(&item))
	}
	return result
}

// =====================================================================
// KELURAHAN / DESA
// =====================================================================

// ToKelurahanDesaResponse mengubah model Kelurahan/Desa menjadi response
func ToKelurahanDesaResponse(m *models.MasterAlamatKelurahanDesa) *KelurahanDesaResponse {
	if m == nil {
		return nil
	}
	return &KelurahanDesaResponse{
		ID:          m.ID,
		KecamatanID: m.KecamatanID,
		Code:        m.Code,
		Name:        m.Name,
		FhirCode:    m.FhirCode,
		PostalCode:  m.PostalCode,
		CreatedBy:   m.CreatedBy,
		UpdatedBy:   m.UpdatedBy,
		CreatedAt:   types.CustomTime(m.CreatedAt),
		UpdatedAt:   types.CustomTime(m.UpdatedAt),
	}
}

func ToKelurahanDesaListSimpelResponse(items []models.MasterAlamatKelurahanDesa) []KelurahanDesaSimpelResponse {
	result := make([]KelurahanDesaSimpelResponse, 0, len(items))
	for _, item := range items {
		result = append(result, KelurahanDesaSimpelResponse{
			ID:          item.ID,
			KecamatanID: item.KecamatanID,
			Code:        item.Code,
			Name:        item.Name,
			PostalCode:  item.PostalCode,
		})
	}
	return result
}

// ToKelurahanDesaListResponse mengubah slice model Kelurahan/Desa menjadi slice response
func ToKelurahanDesaListResponse(items []models.MasterAlamatKelurahanDesa) []KelurahanDesaResponse {
	result := make([]KelurahanDesaResponse, 0, len(items))
	for _, item := range items {
		result = append(result, *ToKelurahanDesaResponse(&item))
	}
	return result
}
