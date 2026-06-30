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
func (s *service) GetByIDGolonganDarah(id int64) (*dto.MasterGolonganDarahResponse, error) {
	ctx := context.Background()
	cacheKey := cacheKeyGolonganDarahDetail(id)

	// 1. Cek Cache
	var cachedRes dto.MasterGolonganDarahResponse
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		return &cachedRes, nil
	}

	// 2. Hit Database
	m, err := s.repo.GetByIDGolonganDarah(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("GolonganDarah tidak ditemukan")
	}

	res := dto.ToMasterGolonganDarahResponse(m)

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctx, cacheKey, res)
	return res, nil
}

// ─────────────── List ────────────────────────────────────────────────────────────
func (s *service) ListGolonganDarah(page, pageSize int, filter *dto.FilterMasterGolonganDarahRequest) ([]dto.MasterGolonganDarahResponse, int64, error) {
	ctx := context.Background()
	cacheKey := cacheKeyGolonganDarahList(page, pageSize, filter)

	// 1. Cek Cache
	var cachedRes struct {
		Items []dto.MasterGolonganDarahResponse `json:"items"`
		Total int64                             `json:"total"`
	}
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		if len(cachedRes.Items) == 0 {
			return nil, 0, appErrors.NotFound("GolonganDarah tidak ditemukan")
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

	items, total, err := s.repo.ListGolonganDarah(page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}
	if len(items) == 0 {
		return nil, 0, appErrors.NotFound("GolonganDarah tidak ditemukan")
	}

	res := dto.ToMasterGolonganDarahListResponse(items)

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctx, cacheKey, struct {
		Items []dto.MasterGolonganDarahResponse `json:"items"`
		Total int64                             `json:"total"`
	}{Items: res, Total: total})

	return res, total, nil
}

// ─────────────── Create ──────────────────────────────────────────────────────────
func (s *service) CreateGolonganDarah(req *dto.CreateMasterGolonganDarahRequest, actor he.AuthContext) (*dto.MasterGolonganDarahResponse, error) {
	// Permission Check
	can, err := s.canCreateMaster(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat GolonganDarah baru.", nil)
	}

	// Check Duplicate Name
	data, err := s.repo.GetByNameGolonganDarah(req.Name)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return nil, appErrors.Wrap(http.StatusConflict, "GolonganDarah dengan nama ini sudah ada", nil)
	}

	// Logic
	m := &models.MasterGolonganDarah{
		KodeKemenkes: req.KodeKemenkes,
		Name:         req.Name,
		Description:  req.Description,
		CreatedBy:    &actor.UserID,
		UpdatedBy:    &actor.UserID,
	}
	if err := s.repo.CreateGolonganDarah(m); err != nil {
		return nil, err
	}

	res := dto.ToMasterGolonganDarahResponse(m)

	// Invalidate Cache
	s.cache.InvalidateList(context.Background(), cachePrefixGolonganDarahList)

	return res, nil

}

// ─────────────── Update ──────────────────────────────────────────────────────────
func (s *service) UpdateGolonganDarah(id int64, req *dto.UpdateMasterGolonganDarahRequest, actor he.AuthContext) (*dto.MasterGolonganDarahResponse, error) {
	// Permission Check
	can, err := s.canUpdateMaster(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah GolonganDarah.", nil)
	}

	// Cek Keberadaan Data
	existing, err := s.repo.GetByIDGolonganDarah(id)
	if err != nil || existing == nil {
		return nil, appErrors.NotFound("GolonganDarah tidak ditemukan")
	}

	// Check Duplicate Name (jika ada perubahan)
	if req.Name != nil && *req.Name != existing.Name {
		data, err := s.repo.GetByNameGolonganDarah(*req.Name)
		if err != nil {
			return nil, err
		}
		if data != nil {
			return nil, appErrors.Wrap(http.StatusConflict, "GolonganDarah dengan nama ini sudah ada", nil)
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

	if err := s.repo.UpdateGolonganDarah(existing); err != nil {
		return nil, err
	}

	res := dto.ToMasterGolonganDarahResponse(existing)

	// Invalidate Cache
	s.cache.InvalidateDetail(context.Background(), cacheKeyGolonganDarahDetail(id))
	s.cache.InvalidateList(context.Background(), cachePrefixGolonganDarahList)

	return res, nil

}

// ─────────────── Delete ──────────────────────────────────────────────────────────
func (s *service) DeleteGolonganDarah(id int64, actor he.AuthContext) error {
	// Permission Check
	can, err := s.canDeleteMaster(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus GolonganDarah.", nil)
	}

	// Cek Keberadaan Data
	existing, err := s.repo.GetByIDGolonganDarah(id)
	if err != nil || existing == nil {
		return appErrors.NotFound("GolonganDarah tidak ditemukan")
	}

	// delete
	err = s.repo.DeleteGolonganDarah(id)
	if err != nil {
		return err
	}

	// Invalidate Cache
	s.cache.InvalidateList(context.Background(), cachePrefixGolonganDarahList)

	return nil
}
