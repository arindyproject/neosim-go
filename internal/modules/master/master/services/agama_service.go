package services

import (
	"context"
	"neosim_go/internal/modules/master/master/dto"
	"neosim_go/internal/modules/master/master/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
	"net/http"
)

// ─────────────── GetByID ─────────────────────────────────────────────────────────
func (s *service) GetByIDAgama(ctx context.Context, id int64) (*dto.MasterAgamaResponse, error) {
	ctxs := context.Background()
	cacheKey := cacheKeyAgamaDetail(id)

	// 1. Cek Cache
	var cachedRes dto.MasterAgamaResponse
	if s.cache.Get(ctxs, cacheKey, &cachedRes) {
		return &cachedRes, nil
	}

	// 2. Hit Database
	m, err := s.repo.GetByIDAgama(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Agama tidak ditemukan")
	}

	res := dto.ToMasterAgamaResponse(m)

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctxs, cacheKey, res)
	return res, nil
}

// ─────────────── List ────────────────────────────────────────────────────────────
func (s *service) ListAgama(ctx context.Context, page, pageSize int, filter *dto.FilterMasterAgamaRequest) ([]dto.MasterAgamaResponse, int64, error) {
	ctxs := context.Background()
	cacheKey := cacheKeyAgamaList(page, pageSize, filter)

	// 1. Cek Cache
	var cachedRes struct {
		Items []dto.MasterAgamaResponse `json:"items"`
		Total int64                     `json:"total"`
	}
	if s.cache.Get(ctxs, cacheKey, &cachedRes) {
		if len(cachedRes.Items) == 0 {
			return nil, 0, appErrors.NotFound("Agama tidak ditemukan")
		}
		return cachedRes.Items, cachedRes.Total, nil
	}

	// 2. Hit Database
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSize
	}

	items, total, err := s.repo.ListAgama(ctx, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}
	if len(items) == 0 {
		return nil, 0, appErrors.NotFound("Agama tidak ditemukan")
	}

	res := dto.ToMasterAgamaListResponse(items)

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctxs, cacheKey, struct {
		Items []dto.MasterAgamaResponse `json:"items"`
		Total int64                     `json:"total"`
	}{Items: res, Total: total})

	return res, total, nil
}

// ─────────────── Create ──────────────────────────────────────────────────────────
func (s *service) CreateAgama(ctx context.Context, req *dto.CreateMasterAgamaRequest, actor he.AuthContext) (*dto.MasterAgamaResponse, error) {
	// Permission Check
	can, err := s.canCreateMaster(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Agama baru.", nil)
	}

	// Check Duplicate Name
	data, err := s.repo.GetByNameAgama(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return nil, appErrors.Wrap(http.StatusConflict, "Agama dengan nama ini sudah ada", nil)
	}

	// Logic
	m := &models.MasterAgama{
		KodeKemenkes: req.KodeKemenkes,
		Name:         req.Name,
		Description:  req.Description,
		CreatedBy:    &actor.UserID,
		UpdatedBy:    &actor.UserID,
	}
	if err := s.repo.CreateAgama(ctx, m); err != nil {
		return nil, err
	}

	res := dto.ToMasterAgamaResponse(m)

	// Invalidate Cache
	s.cache.InvalidateList(context.Background(), cachePrefixAgamaList)

	return res, nil

}

// ─────────────── Update ──────────────────────────────────────────────────────────
func (s *service) UpdateAgama(ctx context.Context, id int64, req *dto.UpdateMasterAgamaRequest, actor he.AuthContext) (*dto.MasterAgamaResponse, error) {
	// Permission Check
	can, err := s.canUpdateMaster(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Agama.", nil)
	}

	// Cek Keberadaan Data
	existing, err := s.repo.GetByIDAgama(ctx, id)
	if err != nil || existing == nil {
		return nil, appErrors.NotFound("Agama tidak ditemukan")
	}

	// Check Duplicate Name (jika ada perubahan)
	if req.Name != nil && *req.Name != existing.Name {
		data, err := s.repo.GetByNameAgama(ctx, *req.Name)
		if err != nil {
			return nil, err
		}
		if data != nil {
			return nil, appErrors.Wrap(http.StatusConflict, "Agama dengan nama ini sudah ada", nil)
		}
	}

	// Logic
	if req.KodeKemenkes != nil {
		existing.KodeKemenkes = req.KodeKemenkes
	}
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Description != nil {
		existing.Description = req.Description
	}
	existing.UpdatedBy = &actor.UserID

	if err := s.repo.UpdateAgama(ctx, existing); err != nil {
		return nil, err
	}

	res := dto.ToMasterAgamaResponse(existing)

	// Invalidate Cache
	s.cache.InvalidateDetail(context.Background(), cacheKeyAgamaDetail(id))
	s.cache.InvalidateList(context.Background(), cachePrefixAgamaList)

	return res, nil

}

// ─────────────── Delete ──────────────────────────────────────────────────────────
func (s *service) DeleteAgama(ctx context.Context, id int64, actor he.AuthContext) error {
	// Permission Check
	can, err := s.canDeleteMaster(ctx, actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Agama.", nil)
	}

	// Cek Keberadaan Data
	existing, err := s.repo.GetByIDAgama(ctx, id)
	if err != nil || existing == nil {
		return appErrors.NotFound("Agama tidak ditemukan")
	}

	// delete
	err = s.repo.DeleteAgama(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate Cache
	s.cache.InvalidateList(context.Background(), cachePrefixAgamaList)

	return nil
}
