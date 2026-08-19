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

// Negara ==========================================================================

// ─────────────── GetByID ─────────────────────────────────────────────────────────
func (s *service) GetByIDNegara(ctx context.Context, id int64) (*dto.NegaraResponse, error) {
	ctxs := context.Background()
	cacheKey := cacheKeyNegaraDetail(id)

	// 1. Cek Cache
	var cachedRes dto.NegaraResponse
	if s.cache.Get(ctxs, cacheKey, &cachedRes) {
		return &cachedRes, nil
	}

	// 2. Ambil dari DB
	m, err := s.repo.GetByIDNegara(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Negara tidak ditemukan")
	}

	res := dto.ToNegaraResponse(m)

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctxs, cacheKey, res)
	return res, nil
}

// ─────────────── List ────────────────────────────────────────────────────────────
func (s *service) ListNegara(ctx context.Context, page, pageSize int, filter *dto.FilterNegaraRequest) ([]dto.NegaraResponse, int64, error) {
	ctxs := context.Background()
	cacheKey := cacheKeyNegaraList(page, pageSize, filter)

	// 1. Cek Cache
	var cachedRes struct {
		Items []dto.NegaraResponse `json:"items"`
		Total int64                `json:"total"`
	}
	if s.cache.Get(ctxs, cacheKey, &cachedRes) {
		if len(cachedRes.Items) == 0 {
			return nil, 0, appErrors.NotFound("Negara tidak ditemukan")
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

	items, total, err := s.repo.ListNegara(ctx, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}
	if len(items) == 0 {
		return nil, 0, appErrors.NotFound("Negara tidak ditemukan")
	}

	res := dto.ToNegaraListResponse(items)

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctxs, cacheKey, struct {
		Items []dto.NegaraResponse `json:"items"`
		Total int64                `json:"total"`
	}{Items: res, Total: total})

	return res, total, nil
}

// ─────────────── Create ──────────────────────────────────────────────────────────
func (s *service) CreateNegara(ctx context.Context, req *dto.CreateNegaraRequest, actor he.AuthContext) (*dto.NegaraResponse, error) {
	// Permission
	can, err := s.canCreateMasterAlamat(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat MasterAlamat baru.", nil)
	}

	// Cek duplikat code
	exists, err := s.repo.ExistsNegaraByCode(ctx, req.Code, nil)
	if err != nil {
		return nil, appErrors.Internal("gagal cek duplikat code")
	}
	if exists {
		return nil, appErrors.Wrap(http.StatusConflict,
			"Code negara '"+req.Code+"' sudah digunakan.", nil)
	}

	// Logic
	m := &models.MasterAlamatNegara{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   &actor.UserID,
		UpdatedBy:   &actor.UserID,
	}
	if err := s.repo.CreateNegara(ctx, m); err != nil {
		return nil, err
	}

	res := dto.ToNegaraResponse(m)

	// Invalidate Cache
	s.cache.InvalidateList(context.Background(), cachePrefixNegaraList)

	return res, nil
}

// ─────────────── Update ──────────────────────────────────────────────────────────
func (s *service) UpdateNegara(ctx context.Context, id int64, req *dto.UpdateNegaraRequest, actor he.AuthContext) (*dto.NegaraResponse, error) {
	// Permission
	can, err := s.canUpdateMasterAlamat(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah MasterAlamat.", nil)
	}

	// Logic
	m, err := s.repo.GetByIDNegara(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Negara tidak ditemukan")
	}

	if req.Code != nil {
		// Cek duplikat code, exclude diri sendiri
		exists, err := s.repo.ExistsNegaraByCode(ctx, *req.Code, &id)
		if err != nil {
			return nil, appErrors.Internal("gagal cek duplikat code")
		}
		if exists {
			return nil, appErrors.Wrap(http.StatusConflict,
				"Code negara '"+*req.Code+"' sudah digunakan.", nil)
		}
		m.Code = *req.Code
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateNegara(ctx, m); err != nil {
		return nil, err
	}

	res := dto.ToNegaraResponse(m)

	// Invalidate Cache
	ctxs := context.Background()
	s.cache.InvalidateDetail(ctxs, cacheKeyNegaraDetail(id))
	s.cache.InvalidateList(ctxs, cachePrefixNegaraList)

	return res, nil
}

// ─────────────── Delete ──────────────────────────────────────────────────────────
func (s *service) DeleteNegara(ctx context.Context, id int64, actor he.AuthContext) error {
	// Permission
	can, err := s.canDeleteMasterAlamat(ctx, actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus MasterAlamat.", nil)
	}

	// Logic
	m, err := s.repo.GetByIDNegara(ctx, id)
	if err != nil {
		return err
	}
	if m == nil {
		return appErrors.NotFound("Negara tidak ditemukan")
	}

	err = s.repo.DeleteNegara(ctx, id)
	if err == nil {
		// Invalidate Cache
		ctxs := context.Background()
		s.cache.InvalidateDetail(ctxs, cacheKeyNegaraDetail(id))
		s.cache.InvalidateList(ctxs, cachePrefixNegaraList)
	}
	return err
}
