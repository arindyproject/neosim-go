package services

import (
	"context"
	"net/http"
	"time"

	"neosim_go/internal/modules/master/alamat/dto"
	"neosim_go/internal/modules/master/alamat/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

// Provinsi ========================================================================

// ─────────────── GetByID ─────────────────────────────────────────────────────────
func (s *service) GetByIDProvinsi(id int64) (*dto.ProvinsiDetailResponse, error) {
	ctx := context.Background()
	cacheKey := cacheKeyProvinsiDetail(id)

	// 1. Cek Cache
	var cachedRes dto.ProvinsiDetailResponse
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		return &cachedRes, nil
	}

	// 2. Ambil dari DB
	m, err := s.repo.GetByIDProvinsi(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Provinsi tidak ditemukan")
	}

	totalKota, err := s.repo.CountKotaByProvinsiID(m.ID)
	if err != nil {
		return nil, err
	}
	totalKecamatan, err := s.repo.CountKecamatanByProvinsiID(m.ID)
	if err != nil {
		return nil, err
	}
	totalDesa, err := s.repo.CountDesaByProvinsiID(m.ID)
	if err != nil {
		return nil, err
	}

	res := &dto.ProvinsiDetailResponse{
		ID:             m.ID,
		Code:           m.Code,
		Name:           m.Name,
		NegaraID:       m.NegaraID,
		NegaraName:     m.Negara.Name,
		TotalKota:      totalKota,
		TotalKecamatan: totalKecamatan,
		TotalDesa:      totalDesa,
	}

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctx, cacheKey, res)
	return res, nil
}

// ─────────────── List ────────────────────────────────────────────────────────────
func (s *service) ListProvinsi(page, pageSize int, negaraID *int64, filter *dto.FilterProvinsiRequest) ([]dto.ProvinsiResponse, int64, error) {
	ctx := context.Background()
	cacheKey := cacheKeyProvinsiList(page, pageSize, negaraID, filter)

	// 1. Cek Cache
	var cachedRes struct {
		Items []dto.ProvinsiResponse `json:"items"`
		Total int64                  `json:"total"`
	}
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		if len(cachedRes.Items) == 0 {
			return nil, 0, appErrors.NotFound("Provinsi tidak ditemukan")
		}
		return cachedRes.Items, cachedRes.Total, nil
	}

	// 2. Ambil dari DB
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	items, total, err := s.repo.ListProvinsi(page, pageSize, negaraID, filter)
	if err != nil {
		return nil, 0, err
	}
	if len(items) == 0 {
		return nil, 0, appErrors.NotFound("Provinsi tidak ditemukan")
	}

	res := dto.ToProvinsiListResponse(items)

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctx, cacheKey, struct {
		Items []dto.ProvinsiResponse `json:"items"`
		Total int64                  `json:"total"`
	}{Items: res, Total: total})

	return res, total, nil
}

// ─────────────── Create ──────────────────────────────────────────────────────────
func (s *service) CreateProvinsi(req *dto.CreateProvinsiRequest, actor he.AuthContext) (*dto.ProvinsiResponse, error) {
	// Permission
	can, err := s.canCreateMasterAlamat(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Provinsi baru.", nil)
	}

	// Cek duplikat code
	exists, err := s.repo.ExistsProvinsiByCode(req.Code, nil)
	if err != nil {
		return nil, appErrors.Internal("gagal cek duplikat code")
	}
	if exists {
		return nil, appErrors.Wrap(http.StatusConflict,
			"Code provinsi '"+req.Code+"' sudah digunakan.", nil)
	}

	// Logic
	negara, err := s.repo.GetByIDNegara(req.NegaraID)
	if err != nil {
		return nil, err
	}
	if negara == nil {
		return nil, appErrors.NotFound("Negara tidak ditemukan")
	}

	m := &models.MasterAlamatProvinsi{
		NegaraID:  req.NegaraID,
		Code:      req.Code,
		Name:      req.Name,
		CreatedBy: &actor.UserID,
		UpdatedBy: &actor.UserID,
	}
	if err := s.repo.CreateProvinsi(m); err != nil {
		return nil, err
	}

	res := dto.ToProvinsiResponse(m)

	// Invalidate Cache
	// Provinsi baru mempengaruhi list child (kota, kecamatan, desa) jika mereka embed nama provinsi
	ctx := context.Background()
	s.cache.InvalidateList(ctx, cachePrefixProvinsiList)
	s.cache.InvalidateList(ctx, cachePrefixKotaList)
	s.cache.InvalidateList(ctx, cachePrefixKecamatanList)
	s.cache.InvalidateList(ctx, cachePrefixDesaList)

	return res, nil
}

// ─────────────── Update ──────────────────────────────────────────────────────────
func (s *service) UpdateProvinsi(id int64, req *dto.UpdateProvinsiRequest, actor he.AuthContext) (*dto.ProvinsiResponse, error) {
	// Permission
	can, err := s.canUpdateMasterAlamat(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Provinsi.", nil)
	}

	// Logic
	m, err := s.repo.GetByIDProvinsi(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Provinsi tidak ditemukan")
	}

	if req.NegaraID != nil {
		m.NegaraID = *req.NegaraID
	}
	if req.Code != nil {
		// Cek duplikat code, exclude diri sendiri
		exists, err := s.repo.ExistsProvinsiByCode(*req.Code, &id)
		if err != nil {
			return nil, appErrors.Internal("gagal cek duplikat code")
		}
		if exists {
			return nil, appErrors.Wrap(http.StatusConflict,
				"Code provinsi '"+*req.Code+"' sudah digunakan.", nil)
		}
		m.Code = *req.Code
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateProvinsi(m); err != nil {
		return nil, err
	}

	res := dto.ToProvinsiResponse(m)

	// Invalidate Cache
	ctx := context.Background()
	s.cache.InvalidateDetail(ctx, cacheKeyProvinsiDetail(id))
	s.cache.InvalidateDetail(ctx, cacheKeyProvinsiGetDetail(id))
	s.cache.InvalidateList(ctx, cachePrefixProvinsiList)
	s.cache.InvalidateList(ctx, cachePrefixKotaList)
	s.cache.InvalidateList(ctx, cachePrefixKecamatanList)
	s.cache.InvalidateList(ctx, cachePrefixDesaList)

	return res, nil
}

// ─────────────── Delete ──────────────────────────────────────────────────────────
func (s *service) DeleteProvinsi(id int64, actor he.AuthContext) error {
	// Permission
	can, err := s.canDeleteMasterAlamat(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Provinsi.", nil)
	}

	// Logic
	m, err := s.repo.GetByIDProvinsi(id)
	if err != nil {
		return err
	}
	if m == nil {
		return appErrors.NotFound("Provinsi tidak ditemukan")
	}

	err = s.repo.DeleteProvinsi(id)
	if err == nil {
		// Invalidate Cache
		ctx := context.Background()
		s.cache.InvalidateDetail(ctx, cacheKeyProvinsiDetail(id))
		s.cache.InvalidateDetail(ctx, cacheKeyProvinsiGetDetail(id))
		s.cache.InvalidateList(ctx, cachePrefixProvinsiList)
		s.cache.InvalidateList(ctx, cachePrefixKotaList)
		s.cache.InvalidateList(ctx, cachePrefixKecamatanList)
		s.cache.InvalidateList(ctx, cachePrefixDesaList)
	}
	return err
}
