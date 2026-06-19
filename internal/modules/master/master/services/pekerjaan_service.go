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
func (s *service) GetByIDPekerjaan(id int64) (*dto.MasterPekerjaanResponse, error) {
	ctx := context.Background()
	cacheKey := cacheKeyPekerjaanDetail(id)

	// 1. Cek Cache
	var cachedRes dto.MasterPekerjaanResponse
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		return &cachedRes, nil
	}

	// 2. Hit Database
	m, err := s.repo.GetByIDPekerjaan(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Pekerjaan tidak ditemukan")
	}

	res := dto.ToMasterPekerjaanResponse(m)

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctx, cacheKey, res)
	return res, nil
}

// ─────────────── List ────────────────────────────────────────────────────────────
func (s *service) ListPekerjaan(page, pageSize int, filter *dto.FilterMasterPekerjaanRequest) ([]dto.MasterPekerjaanResponse, int64, error) {
	ctx := context.Background()
	cacheKey := cacheKeyPekerjaanList(page, pageSize, filter)

	// 1. Cek Cache
	var cachedRes struct {
		Items []dto.MasterPekerjaanResponse `json:"items"`
		Total int64                         `json:"total"`
	}
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		if len(cachedRes.Items) == 0 {
			return nil, 0, appErrors.NotFound("Pekerjaan tidak ditemukan")
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

	items, total, err := s.repo.ListPekerjaan(page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}
	if len(items) == 0 {
		return nil, 0, appErrors.NotFound("Pekerjaan tidak ditemukan")
	}

	res := dto.ToMasterPekerjaanListResponse(items)

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctx, cacheKey, struct {
		Items []dto.MasterPekerjaanResponse `json:"items"`
		Total int64                         `json:"total"`
	}{Items: res, Total: total})

	return res, total, nil
}

// ─────────────── Create ──────────────────────────────────────────────────────────
func (s *service) CreatePekerjaan(req *dto.CreateMasterPekerjaanRequest, actor he.AuthContext) (*dto.MasterPekerjaanResponse, error) {
	// Permission Check
	can, err := s.canCreateMaster(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Pekerjaan baru.", nil)
	}

	// Check Duplicate Name
	if duplicate, err := s.repo.GetByNamePekerjaan(req.Name); duplicate != nil || err == nil {
		return nil, appErrors.Wrap(http.StatusConflict, "Pekerjaan dengan nama ini sudah ada", nil)
	}

	// Logic
	m := &models.MasterPekerjaan{
		KodeKemenkes: req.KodeKemenkes,
		Name:         req.Name,
		Description:  req.Description,
		CreatedBy:    &actor.UserID,
		UpdatedBy:    &actor.UserID,
	}
	if err := s.repo.CreatePekerjaan(m); err != nil {
		return nil, err
	}

	res := dto.ToMasterPekerjaanResponse(m)

	// Invalidate Cache
	s.cache.InvalidateList(context.Background(), cachePrefixPekerjaanList)

	return res, nil

}

// ─────────────── Update ──────────────────────────────────────────────────────────
func (s *service) UpdatePekerjaan(id int64, req *dto.UpdateMasterPekerjaanRequest, actor he.AuthContext) (*dto.MasterPekerjaanResponse, error) {
	// Permission Check
	can, err := s.canUpdateMaster(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Pekerjaan.", nil)
	}

	// Cek Keberadaan Data
	existing, err := s.repo.GetByIDPekerjaan(id)
	if err != nil || existing == nil {
		return nil, appErrors.NotFound("Pekerjaan tidak ditemukan")
	}

	// Check Duplicate Name (jika ada perubahan)
	if req.Name != nil && *req.Name != existing.Name {
		if duplicate, err := s.repo.GetByNamePekerjaan(*req.Name); duplicate != nil || err == nil {
			return nil, appErrors.Wrap(http.StatusConflict, "Pekerjaan dengan nama ini sudah ada", nil)
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

	if err := s.repo.UpdatePekerjaan(existing); err != nil {
		return nil, err
	}

	res := dto.ToMasterPekerjaanResponse(existing)

	// Invalidate Cache
	s.cache.InvalidateDetail(context.Background(), cacheKeyPekerjaanDetail(id))
	s.cache.InvalidateList(context.Background(), cachePrefixPekerjaanList)

	return res, nil

}

// ─────────────── Delete ──────────────────────────────────────────────────────────
func (s *service) DeletePekerjaan(id int64, actor he.AuthContext) error {
	// Permission Check
	can, err := s.canDeleteMaster(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Pekerjaan.", nil)
	}

	// Cek Keberadaan Data
	existing, err := s.repo.GetByIDPekerjaan(id)
	if err != nil || existing == nil {
		return appErrors.NotFound("Pekerjaan tidak ditemukan")
	}

	// delete
	err = s.repo.DeletePekerjaan(id)
	if err != nil {
		return err
	}

	// Invalidate Cache
	s.cache.InvalidateList(context.Background(), cachePrefixPekerjaanList)

	return nil
}
