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
func (s *service) GetByIDJenisKelamin(id int64) (*dto.MasterJenisKelaminResponse, error) {
	ctx := context.Background()
	cacheKey := cacheKeyJenisKelaminDetail(id)

	// 1. Cek Cache
	var cachedRes dto.MasterJenisKelaminResponse
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		return &cachedRes, nil
	}

	// 2. Hit Database
	m, err := s.repo.GetByIDJenisKelamin(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("JenisKelamin tidak ditemukan")
	}

	res := dto.ToMasterJenisKelaminResponse(m)

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctx, cacheKey, res)
	return res, nil
}

// ─────────────── List ────────────────────────────────────────────────────────────
func (s *service) ListJenisKelamin(page, pageSize int, filter *dto.FilterMasterJenisKelaminRequest) ([]dto.MasterJenisKelaminResponse, int64, error) {
	ctx := context.Background()
	cacheKey := cacheKeyJenisKelaminList(page, pageSize, filter)

	// 1. Cek Cache
	var cachedRes struct {
		Items []dto.MasterJenisKelaminResponse `json:"items"`
		Total int64                            `json:"total"`
	}
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		if len(cachedRes.Items) == 0 {
			return nil, 0, appErrors.NotFound("JenisKelamin tidak ditemukan")
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

	items, total, err := s.repo.ListJenisKelamin(page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}
	if len(items) == 0 {
		return nil, 0, appErrors.NotFound("JenisKelamin tidak ditemukan")
	}

	res := dto.ToMasterJenisKelaminListResponse(items)

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctx, cacheKey, struct {
		Items []dto.MasterJenisKelaminResponse `json:"items"`
		Total int64                            `json:"total"`
	}{Items: res, Total: total})

	return res, total, nil
}

// ─────────────── Create ──────────────────────────────────────────────────────────
func (s *service) CreateJenisKelamin(req *dto.CreateMasterJenisKelaminRequest, actor he.AuthContext) (*dto.MasterJenisKelaminResponse, error) {
	// Permission Check
	can, err := s.canCreateMaster(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat JenisKelamin baru.", nil)
	}

	// Check Duplicate Name
	data, err := s.repo.GetByNameJenisKelamin(req.Name)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return nil, appErrors.Wrap(http.StatusConflict, "JenisKelamin dengan nama ini sudah ada", nil)
	}

	// Logic
	m := &models.MasterJenisKelamin{
		KodeKemenkes: req.KodeKemenkes,
		Name:         req.Name,
		Description:  req.Description,
		CreatedBy:    &actor.UserID,
		UpdatedBy:    &actor.UserID,
	}
	if err := s.repo.CreateJenisKelamin(m); err != nil {
		return nil, err
	}

	res := dto.ToMasterJenisKelaminResponse(m)

	// Invalidate Cache
	s.cache.InvalidateList(context.Background(), cachePrefixJenisKelaminList)

	return res, nil

}

// ─────────────── Update ──────────────────────────────────────────────────────────
func (s *service) UpdateJenisKelamin(id int64, req *dto.UpdateMasterJenisKelaminRequest, actor he.AuthContext) (*dto.MasterJenisKelaminResponse, error) {
	// Permission Check
	can, err := s.canUpdateMaster(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah JenisKelamin.", nil)
	}

	// Cek Keberadaan Data
	existing, err := s.repo.GetByIDJenisKelamin(id)
	if err != nil || existing == nil {
		return nil, appErrors.NotFound("JenisKelamin tidak ditemukan")
	}

	// Check Duplicate Name (jika ada perubahan)
	if req.Name != nil && *req.Name != existing.Name {
		data, err := s.repo.GetByNameJenisKelamin(*req.Name)
		if err != nil {
			return nil, err
		}
		if data != nil {
			return nil, appErrors.Wrap(http.StatusConflict, "JenisKelamin dengan nama ini sudah ada", nil)
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

	if err := s.repo.UpdateJenisKelamin(existing); err != nil {
		return nil, err
	}

	res := dto.ToMasterJenisKelaminResponse(existing)

	// Invalidate Cache
	s.cache.InvalidateDetail(context.Background(), cacheKeyJenisKelaminDetail(id))
	s.cache.InvalidateList(context.Background(), cachePrefixJenisKelaminList)

	return res, nil

}

// ─────────────── Delete ──────────────────────────────────────────────────────────
func (s *service) DeleteJenisKelamin(id int64, actor he.AuthContext) error {
	// Permission Check
	can, err := s.canDeleteMaster(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus JenisKelamin.", nil)
	}

	// Cek Keberadaan Data
	existing, err := s.repo.GetByIDJenisKelamin(id)
	if err != nil || existing == nil {
		return appErrors.NotFound("JenisKelamin tidak ditemukan")
	}

	// delete
	err = s.repo.DeleteJenisKelamin(id)
	if err != nil {
		return err
	}

	// Invalidate Cache
	s.cache.InvalidateList(context.Background(), cachePrefixJenisKelaminList)

	return nil
}
