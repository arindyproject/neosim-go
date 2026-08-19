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
func (s *service) GetByIDProvinsi(ctx context.Context, id int64) (*dto.ProvinsiDetailResponse, error) {
	ctxs := context.Background()
	cacheKey := cacheKeyProvinsiDetail(id)

	// 1. Cek Cache
	var cachedRes dto.ProvinsiDetailResponse
	if s.cache.Get(ctxs, cacheKey, &cachedRes) {
		return &cachedRes, nil
	}

	// 2. Ambil dari DB
	m, err := s.repo.GetByIDProvinsi(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Provinsi tidak ditemukan")
	}

	totalKota, err := s.repo.CountKotaByProvinsiID(ctx, m.ID)
	if err != nil {
		return nil, err
	}
	totalKecamatan, err := s.repo.CountKecamatanByProvinsiID(ctx, m.ID)
	if err != nil {
		return nil, err
	}
	totalDesa, err := s.repo.CountDesaByProvinsiID(ctx, m.ID)
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
	s.cache.SetDefault(ctxs, cacheKey, res)
	return res, nil
}

// ─────────────── List ────────────────────────────────────────────────────────────
func (s *service) ListProvinsi(ctx context.Context, page, pageSize int, negaraID *int64, filter *dto.FilterProvinsiRequest) ([]dto.ProvinsiResponse, int64, error) {
	ctxs := context.Background()
	cacheKey := cacheKeyProvinsiList(page, pageSize, negaraID, filter)

	// 1. Cek Cache
	var cachedRes struct {
		Items []dto.ProvinsiResponse `json:"items"`
		Total int64                  `json:"total"`
	}
	if s.cache.Get(ctxs, cacheKey, &cachedRes) {
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

	items, total, err := s.repo.ListProvinsi(ctx, page, pageSize, negaraID, filter)
	if err != nil {
		return nil, 0, err
	}
	if len(items) == 0 {
		return nil, 0, appErrors.NotFound("Provinsi tidak ditemukan")
	}

	res := dto.ToProvinsiListResponse(items)

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctxs, cacheKey, struct {
		Items []dto.ProvinsiResponse `json:"items"`
		Total int64                  `json:"total"`
	}{Items: res, Total: total})

	return res, total, nil
}

// ─────────────── Create ──────────────────────────────────────────────────────────
func (s *service) CreateProvinsi(ctx context.Context, req *dto.CreateProvinsiRequest, actor he.AuthContext) (*dto.ProvinsiResponse, error) {
	// Permission
	can, err := s.canCreateMasterAlamat(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Provinsi baru.", nil)
	}

	// Cek duplikat code
	exists, err := s.repo.ExistsProvinsiByCode(ctx, req.Code, nil)
	if err != nil {
		return nil, appErrors.Internal("gagal cek duplikat code")
	}
	if exists {
		return nil, appErrors.Wrap(http.StatusConflict,
			"Code provinsi '"+req.Code+"' sudah digunakan.", nil)
	}

	// Logic
	negara, err := s.repo.GetByIDNegara(ctx, req.NegaraID)
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
	if err := s.repo.CreateProvinsi(ctx, m); err != nil {
		return nil, err
	}

	res := dto.ToProvinsiResponse(m)

	// Invalidate Cache
	// Provinsi baru mempengaruhi list child (kota, kecamatan, desa) jika mereka embed nama provinsi
	ctxs := context.Background()
	s.cache.InvalidateList(ctxs, cachePrefixProvinsiList)
	s.cache.InvalidateList(ctxs, cachePrefixKotaList)
	s.cache.InvalidateList(ctxs, cachePrefixKecamatanList)
	s.cache.InvalidateList(ctxs, cachePrefixDesaList)

	return res, nil
}

// ─────────────── Update ──────────────────────────────────────────────────────────
func (s *service) UpdateProvinsi(ctx context.Context, id int64, req *dto.UpdateProvinsiRequest, actor he.AuthContext) (*dto.ProvinsiResponse, error) {
	// Permission
	can, err := s.canUpdateMasterAlamat(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Provinsi.", nil)
	}

	// Logic
	m, err := s.repo.GetByIDProvinsi(ctx, id)
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
		exists, err := s.repo.ExistsProvinsiByCode(ctx, *req.Code, &id)
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

	if err := s.repo.UpdateProvinsi(ctx, m); err != nil {
		return nil, err
	}

	res := dto.ToProvinsiResponse(m)

	// Invalidate Cache
	ctxs := context.Background()
	s.cache.InvalidateDetail(ctxs, cacheKeyProvinsiDetail(id))
	s.cache.InvalidateDetail(ctxs, cacheKeyProvinsiGetDetail(id))
	s.cache.InvalidateList(ctxs, cachePrefixProvinsiList)
	s.cache.InvalidateList(ctxs, cachePrefixKotaList)
	s.cache.InvalidateList(ctxs, cachePrefixKecamatanList)
	s.cache.InvalidateList(ctxs, cachePrefixDesaList)

	return res, nil
}

// ─────────────── Delete ──────────────────────────────────────────────────────────
func (s *service) DeleteProvinsi(ctx context.Context, id int64, actor he.AuthContext) error {
	// Permission
	can, err := s.canDeleteMasterAlamat(ctx, actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Provinsi.", nil)
	}

	// Logic
	m, err := s.repo.GetByIDProvinsi(ctx, id)
	if err != nil {
		return err
	}
	if m == nil {
		return appErrors.NotFound("Provinsi tidak ditemukan")
	}

	err = s.repo.DeleteProvinsi(ctx, id)
	if err == nil {
		// Invalidate Cache
		ctxs := context.Background()
		s.cache.InvalidateDetail(ctxs, cacheKeyProvinsiDetail(id))
		s.cache.InvalidateDetail(ctxs, cacheKeyProvinsiGetDetail(id))
		s.cache.InvalidateList(ctxs, cachePrefixProvinsiList)
		s.cache.InvalidateList(ctxs, cachePrefixKotaList)
		s.cache.InvalidateList(ctxs, cachePrefixKecamatanList)
		s.cache.InvalidateList(ctxs, cachePrefixDesaList)
	}
	return err
}
