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
func (s *service) GetByIDSuku(id int64) (*dto.MasterSukuResponse, error) {
	ctx := context.Background()
	cacheKey := cacheKeySukuDetail(id)

	// 1. Cek Cache
	var cachedRes dto.MasterSukuResponse
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		return &cachedRes, nil
	}

	// 2. Hit Database
	m, err := s.repo.GetByIDSuku(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Suku tidak ditemukan")
	}

	res := dto.ToMasterSukuResponse(m)

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctx, cacheKey, res)
	return res, nil
}

// ─────────────── List ────────────────────────────────────────────────────────────
func (s *service) ListSuku(page, pageSize int, filter *dto.FilterMasterSukuRequest) ([]dto.MasterSukuResponse, int64, error) {
	ctx := context.Background()
	cacheKey := cacheKeySukuList(page, pageSize, filter)

	// 1. Cek Cache
	var cachedRes struct {
		Items []dto.MasterSukuResponse `json:"items"`
		Total int64                    `json:"total"`
	}
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		if len(cachedRes.Items) == 0 {
			return nil, 0, appErrors.NotFound("Suku tidak ditemukan")
		}
		return cachedRes.Items, cachedRes.Total, nil
	}

	// 2. Hit Database
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	items, total, err := s.repo.ListSuku(page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}
	if len(items) == 0 {
		return nil, 0, appErrors.NotFound("Suku tidak ditemukan")
	}

	res := dto.ToMasterSukuListResponse(items)

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctx, cacheKey, struct {
		Items []dto.MasterSukuResponse `json:"items"`
		Total int64                    `json:"total"`
	}{Items: res, Total: total})

	return res, total, nil
}

// ─────────────── Create ──────────────────────────────────────────────────────────
func (s *service) CreateSuku(req *dto.CreateMasterSukuRequest, actor he.AuthContext) (*dto.MasterSukuResponse, error) {
	// Permission Check
	can, err := s.canCreateMaster(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Suku baru.", nil)
	}

	// Check Duplicate Name
	data, err := s.repo.GetByNameSuku(req.Name)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return nil, appErrors.Wrap(http.StatusConflict, "Suku dengan nama ini sudah ada", nil)
	}

	// Logic
	m := &models.MasterSuku{
		KodeKemenkes: req.KodeKemenkes,
		Name:         req.Name,
		Description:  req.Description,
		CreatedBy:    &actor.UserID,
		UpdatedBy:    &actor.UserID,
	}
	if err := s.repo.CreateSuku(m); err != nil {
		return nil, err
	}

	res := dto.ToMasterSukuResponse(m)

	// Invalidate Cache
	s.cache.InvalidateList(context.Background(), cachePrefixSukuList)

	return res, nil

}

// ─────────────── Update ──────────────────────────────────────────────────────────
func (s *service) UpdateSuku(id int64, req *dto.UpdateMasterSukuRequest, actor he.AuthContext) (*dto.MasterSukuResponse, error) {
	// Permission Check
	can, err := s.canUpdateMaster(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Suku.", nil)
	}

	// Cek Keberadaan Data
	existing, err := s.repo.GetByIDSuku(id)
	if err != nil || existing == nil {
		return nil, appErrors.NotFound("Suku tidak ditemukan")
	}

	// Check Duplicate Name (jika ada perubahan)
	if req.Name != nil && *req.Name != existing.Name {
		data, err := s.repo.GetByNameSuku(*req.Name)
		if err != nil {
			return nil, err
		}
		if data != nil {
			return nil, appErrors.Wrap(http.StatusConflict, "Suku dengan nama ini sudah ada", nil)
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

	if err := s.repo.UpdateSuku(existing); err != nil {
		return nil, err
	}

	res := dto.ToMasterSukuResponse(existing)

	// Invalidate Cache
	s.cache.InvalidateDetail(context.Background(), cacheKeySukuDetail(id))
	s.cache.InvalidateList(context.Background(), cachePrefixSukuList)

	return res, nil

}

// ─────────────── Delete ──────────────────────────────────────────────────────────
func (s *service) DeleteSuku(id int64, actor he.AuthContext) error {
	// Permission Check
	can, err := s.canDeleteMaster(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Suku.", nil)
	}

	// Cek Keberadaan Data
	existing, err := s.repo.GetByIDSuku(id)
	if err != nil || existing == nil {
		return appErrors.NotFound("Suku tidak ditemukan")
	}

	// delete
	err = s.repo.DeleteSuku(id)
	if err != nil {
		return err
	}

	// Invalidate Cache
	s.cache.InvalidateList(context.Background(), cachePrefixSukuList)

	return nil
}
