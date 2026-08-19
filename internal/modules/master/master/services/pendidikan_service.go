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
func (s *service) GetByIDPendidikan(ctx context.Context, id int64) (*dto.MasterPendidikanResponse, error) {
	ctxs := context.Background()
	cacheKey := cacheKeyPendidikanDetail(id)

	// 1. Cek Cache
	var cachedRes dto.MasterPendidikanResponse
	if s.cache.Get(ctxs, cacheKey, &cachedRes) {
		return &cachedRes, nil
	}

	// 2. Hit Database
	m, err := s.repo.GetByIDPendidikan(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Pendidikan tidak ditemukan")
	}

	res := dto.ToMasterPendidikanResponse(m)

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctxs, cacheKey, res)
	return res, nil
}

// ─────────────── List ────────────────────────────────────────────────────────────
func (s *service) ListPendidikan(ctx context.Context, page, pageSize int, filter *dto.FilterMasterPendidikanRequest) ([]dto.MasterPendidikanResponse, int64, error) {
	ctxs := context.Background()
	cacheKey := cacheKeyPendidikanList(page, pageSize, filter)

	// 1. Cek Cache
	var cachedRes struct {
		Items []dto.MasterPendidikanResponse `json:"items"`
		Total int64                          `json:"total"`
	}
	if s.cache.Get(ctxs, cacheKey, &cachedRes) {
		if len(cachedRes.Items) == 0 {
			return nil, 0, appErrors.NotFound("Pendidikan tidak ditemukan")
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

	items, total, err := s.repo.ListPendidikan(ctx, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}
	if len(items) == 0 {
		return nil, 0, appErrors.NotFound("Pendidikan tidak ditemukan")
	}

	res := dto.ToMasterPendidikanListResponse(items)

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctxs, cacheKey, struct {
		Items []dto.MasterPendidikanResponse `json:"items"`
		Total int64                          `json:"total"`
	}{Items: res, Total: total})

	return res, total, nil
}

// ─────────────── Create ──────────────────────────────────────────────────────────
func (s *service) CreatePendidikan(ctx context.Context, req *dto.CreateMasterPendidikanRequest, actor he.AuthContext) (*dto.MasterPendidikanResponse, error) {
	// Permission Check
	can, err := s.canCreateMaster(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Pendidikan baru.", nil)
	}

	// Check Duplicate Name
	data, err := s.repo.GetByNamePendidikan(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return nil, appErrors.Wrap(http.StatusConflict, "Pendidikan dengan nama ini sudah ada", nil)
	}

	// Logic
	m := &models.MasterPendidikan{
		KodeKemenkes: req.KodeKemenkes,
		Name:         req.Name,
		Description:  req.Description,
		CreatedBy:    &actor.UserID,
		UpdatedBy:    &actor.UserID,
	}
	if err := s.repo.CreatePendidikan(ctx, m); err != nil {
		return nil, err
	}

	res := dto.ToMasterPendidikanResponse(m)

	// Invalidate Cache
	s.cache.InvalidateList(context.Background(), cachePrefixPendidikanList)

	return res, nil

}

// ─────────────── Update ──────────────────────────────────────────────────────────
func (s *service) UpdatePendidikan(ctx context.Context, id int64, req *dto.UpdateMasterPendidikanRequest, actor he.AuthContext) (*dto.MasterPendidikanResponse, error) {
	// Permission Check
	can, err := s.canUpdateMaster(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Pendidikan.", nil)
	}

	// Cek Keberadaan Data
	existing, err := s.repo.GetByIDPendidikan(ctx, id)
	if err != nil || existing == nil {
		return nil, appErrors.NotFound("Pendidikan tidak ditemukan")
	}

	// Check Duplicate Name (jika ada perubahan)
	if req.Name != nil && *req.Name != existing.Name {
		data, err := s.repo.GetByNamePendidikan(ctx, *req.Name)
		if err != nil {
			return nil, err
		}
		if data != nil {
			return nil, appErrors.Wrap(http.StatusConflict, "Pendidikan dengan nama ini sudah ada", nil)
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

	if err := s.repo.UpdatePendidikan(ctx, existing); err != nil {
		return nil, err
	}

	res := dto.ToMasterPendidikanResponse(existing)

	// Invalidate Cache
	s.cache.InvalidateDetail(context.Background(), cacheKeyPendidikanDetail(id))
	s.cache.InvalidateList(context.Background(), cachePrefixPendidikanList)

	return res, nil

}

// ─────────────── Delete ──────────────────────────────────────────────────────────
func (s *service) DeletePendidikan(ctx context.Context, id int64, actor he.AuthContext) error {
	// Permission Check
	can, err := s.canDeleteMaster(ctx, actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Pendidikan.", nil)
	}

	// Cek Keberadaan Data
	existing, err := s.repo.GetByIDPendidikan(ctx, id)
	if err != nil || existing == nil {
		return appErrors.NotFound("Pendidikan tidak ditemukan")
	}

	// delete
	err = s.repo.DeletePendidikan(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate Cache
	s.cache.InvalidateList(context.Background(), cachePrefixPendidikanList)

	return nil
}
