package dto

import (
	"time"

	"neosim_go/internal/modules/master/alamat/models"
)

// ProvinsiDetailResponse response detail provinsi dengan statistik turunan
// -------------------------------------------------------------------------
type ProvinsiDetailResponse struct {
	ID             int64  `json:"id"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	NegaraID       int64  `json:"negara_id"`
	NegaraName     string `json:"negara_name"`
	TotalKota      int64  `json:"total_kota"`
	TotalKecamatan int64  `json:"total_kecamatan"`
	TotalDesa      int64  `json:"total_desa"`
} //------------------------------------------------------------------------

// KotaKabupatenDetailResponse response detail kota/kabupaten dengan statistik turunan
// -------------------------------------------------------------------------
type KotaKabupatenDetailResponse struct {
	ID             int64  `json:"id"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	ProvinsiID     int64  `json:"provinsi_id"`
	ProvinsiName   string `json:"provinsi_name"`
	NegaraID       int64  `json:"negara_id"`
	NegaraName     string `json:"negara_name"`
	TotalKecamatan int64  `json:"total_kecamatan"`
	TotalDesa      int64  `json:"total_desa"`
} //------------------------------------------------------------------------

// KecamatanDetailResponse response detail kecamatan dengan statistik turunan
// -------------------------------------------------------------------------
type KecamatanDetailResponse struct {
	ID                int64  `json:"id"`
	Code              string `json:"code"`
	Name              string `json:"name"`
	KotaKabupatenID   int64  `json:"kota_kabupaten_id"`
	KotaKabupatenName string `json:"kota_kabupaten_name"`
	ProvinsiID        int64  `json:"provinsi_id"`
	ProvinsiName      string `json:"provinsi_name"`
	NegaraID          int64  `json:"negara_id"`
	NegaraName        string `json:"negara_name"`
	TotalDesa         int64  `json:"total_desa"`
} //------------------------------------------------------------------------

// KelurahanDesaDetailResponse response detail desa/kelurahan dengan jalur hierarki lengkap
// -------------------------------------------------------------------------
type KelurahanDesaDetailResponse struct {
	ID                int64   `json:"id"`
	Code              string  `json:"code"`
	Name              string  `json:"name"`
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
	ID          int64     `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedBy   *int64    `json:"created_by"`
	UpdatedBy   *int64    `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProvinsiResponse response untuk single Provinsi
type ProvinsiResponse struct {
	ID        int64     `json:"id"`
	NegaraID  int64     `json:"negara_id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedBy *int64    `json:"created_by"`
	UpdatedBy *int64    `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// KotaKabupatenResponse response untuk single Kota/Kabupaten
type KotaKabupatenResponse struct {
	ID         int64     `json:"id"`
	ProvinsiID int64     `json:"provinsi_id"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	CreatedBy  *int64    `json:"created_by"`
	UpdatedBy  *int64    `json:"updated_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// KecamatanResponse response untuk single Kecamatan
type KecamatanResponse struct {
	ID              int64     `json:"id"`
	KotaKabupatenID int64     `json:"kota_kabupaten_id"`
	Code            string    `json:"code"`
	Name            string    `json:"name"`
	CreatedBy       *int64    `json:"created_by"`
	UpdatedBy       *int64    `json:"updated_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// KelurahanDesaResponse response untuk single Kelurahan/Desa
type KelurahanDesaResponse struct {
	ID          int64     `json:"id"`
	KecamatanID int64     `json:"kecamatan_id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	PostalCode  *string   `json:"postal_code"`
	CreatedBy   *int64    `json:"created_by"`
	UpdatedBy   *int64    `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
		Description: m.Description,
		CreatedBy:   m.CreatedBy,
		UpdatedBy:   m.UpdatedBy,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
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
		CreatedBy: m.CreatedBy,
		UpdatedBy: m.UpdatedBy,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
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
		CreatedBy:  m.CreatedBy,
		UpdatedBy:  m.UpdatedBy,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
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
		CreatedBy:       m.CreatedBy,
		UpdatedBy:       m.UpdatedBy,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
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
		PostalCode:  m.PostalCode,
		CreatedBy:   m.CreatedBy,
		UpdatedBy:   m.UpdatedBy,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// ToKelurahanDesaListResponse mengubah slice model Kelurahan/Desa menjadi slice response
func ToKelurahanDesaListResponse(items []models.MasterAlamatKelurahanDesa) []KelurahanDesaResponse {
	result := make([]KelurahanDesaResponse, 0, len(items))
	for _, item := range items {
		result = append(result, *ToKelurahanDesaResponse(&item))
	}
	return result
}
