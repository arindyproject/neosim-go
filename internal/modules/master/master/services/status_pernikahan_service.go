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
func (s *service) GetByIDStatusPernikahan(id int64) (*dto.MasterStatusPernikahanResponse, error) {
	ctx := context.Background()
	cacheKey := cacheKeyStatusPernikahanDetail(id)

	// 1. Cek Cache
	var cachedRes dto.MasterStatusPernikahanResponse
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		return &cachedRes, nil
	}

	// 2. Hit Database
	m, err := s.repo.GetByIDStatusPernikahan(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("StatusPernikahan tidak ditemukan")
	}

	res := dto.ToMasterStatusPernikahanResponse(m)

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctx, cacheKey, res)
	return res, nil
}

// ─────────────── List ────────────────────────────────────────────────────────────
func (s *service) ListStatusPernikahan(page, pageSize int, filter *dto.FilterMasterStatusPernikahanRequest) ([]dto.MasterStatusPernikahanResponse, int64, error) {
	ctx := context.Background()
	cacheKey := cacheKeyStatusPernikahanList(page, pageSize, filter)

	// 1. Cek Cache
	var cachedRes struct {
		Items []dto.MasterStatusPernikahanResponse `json:"items"`
		Total int64                                `json:"total"`
	}
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		if len(cachedRes.Items) == 0 {
			return nil, 0, appErrors.NotFound("StatusPernikahan tidak ditemukan")
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

	items, total, err := s.repo.ListStatusPernikahan(page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}
	if len(items) == 0 {
		return nil, 0, appErrors.NotFound("StatusPernikahan tidak ditemukan")
	}

	res := dto.ToMasterStatusPernikahanListResponse(items)

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctx, cacheKey, struct {
		Items []dto.MasterStatusPernikahanResponse `json:"items"`
		Total int64                                `json:"total"`
	}{Items: res, Total: total})

	return res, total, nil
}

// ─────────────── Create ──────────────────────────────────────────────────────────
func (s *service) CreateStatusPernikahan(req *dto.CreateMasterStatusPernikahanRequest, actor he.AuthContext) (*dto.MasterStatusPernikahanResponse, error) {
	// Permission Check
	can, err := s.canCreateMaster(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat StatusPernikahan baru.", nil)
	}

	// Check Duplicate Name
	data, err := s.repo.GetByNameStatusPernikahan(req.Name)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return nil, appErrors.Wrap(http.StatusConflict, "StatusPernikahan dengan nama ini sudah ada", nil)
	}

	// Logic
	m := &models.MasterStatusPernikahan{
		KodeKemenkes: req.KodeKemenkes,
		Name:         req.Name,
		Description:  req.Description,
		CreatedBy:    &actor.UserID,
		UpdatedBy:    &actor.UserID,
	}
	if err := s.repo.CreateStatusPernikahan(m); err != nil {
		return nil, err
	}

	res := dto.ToMasterStatusPernikahanResponse(m)

	// Invalidate Cache
	s.cache.InvalidateList(context.Background(), cachePrefixStatusPernikahanList)

	return res, nil

}

// ─────────────── Update ──────────────────────────────────────────────────────────
func (s *service) UpdateStatusPernikahan(id int64, req *dto.UpdateMasterStatusPernikahanRequest, actor he.AuthContext) (*dto.MasterStatusPernikahanResponse, error) {
	// Permission Check
	can, err := s.canUpdateMaster(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah StatusPernikahan.", nil)
	}

	// Cek Keberadaan Data
	existing, err := s.repo.GetByIDStatusPernikahan(id)
	if err != nil || existing == nil {
		return nil, appErrors.NotFound("StatusPernikahan tidak ditemukan")
	}

	// Check Duplicate Name (jika ada perubahan)
	if req.Name != nil && *req.Name != existing.Name {
		data, err := s.repo.GetByNameStatusPernikahan(*req.Name)
		if err != nil {
			return nil, err
		}
		if data != nil {
			return nil, appErrors.Wrap(http.StatusConflict, "StatusPernikahan dengan nama ini sudah ada", nil)
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

	if err := s.repo.UpdateStatusPernikahan(existing); err != nil {
		return nil, err
	}

	res := dto.ToMasterStatusPernikahanResponse(existing)

	// Invalidate Cache
	s.cache.InvalidateDetail(context.Background(), cacheKeyStatusPernikahanDetail(id))
	s.cache.InvalidateList(context.Background(), cachePrefixStatusPernikahanList)

	return res, nil

}

// ─────────────── Delete ──────────────────────────────────────────────────────────
func (s *service) DeleteStatusPernikahan(id int64, actor he.AuthContext) error {
	// Permission Check
	can, err := s.canDeleteMaster(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus StatusPernikahan.", nil)
	}

	// Cek Keberadaan Data
	existing, err := s.repo.GetByIDStatusPernikahan(id)
	if err != nil || existing == nil {
		return appErrors.NotFound("StatusPernikahan tidak ditemukan")
	}

	// delete
	err = s.repo.DeleteStatusPernikahan(id)
	if err != nil {
		return err
	}

	// Invalidate Cache
	s.cache.InvalidateList(context.Background(), cachePrefixStatusPernikahanList)

	return nil
}
