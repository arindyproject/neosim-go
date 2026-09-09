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

// Kelurahan/Desa ==================================================================

// ─────────────── GetByID ─────────────────────────────────────────────────────────
func (s *service) GetByIDKelurahanDesa(ctx context.Context, id int64) (*dto.KelurahanDesaDetailResponse, error) {
	ctxs := context.Background()
	cacheKey := cacheKeyDesaDetail(id)

	// 1. Cek Cache
	var cachedRes dto.KelurahanDesaDetailResponse
	if s.cache.Get(ctxs, cacheKey, &cachedRes) {
		return &cachedRes, nil
	}

	// 2. Ambil dari DB
	m, err := s.repo.GetByIDKelurahanDesa(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Kelurahan/Desa tidak ditemukan")
	}

	res := &dto.KelurahanDesaDetailResponse{
		ID:                m.ID,
		Code:              m.Code,
		Name:              m.Name,
		FhirCode:          m.FhirCode,
		PostalCode:        m.PostalCode,
		KecamatanID:       m.KecamatanID,
		KecamatanName:     m.Kecamatan.Name,
		KotaKabupatenID:   m.Kecamatan.KotaKabupatenID,
		KotaKabupatenName: m.Kecamatan.KotaKabupaten.Name,
		ProvinsiID:        m.Kecamatan.KotaKabupaten.ProvinsiID,
		ProvinsiName:      m.Kecamatan.KotaKabupaten.Provinsi.Name,
		NegaraID:          m.Kecamatan.KotaKabupaten.Provinsi.NegaraID,
		NegaraName:        m.Kecamatan.KotaKabupaten.Provinsi.Negara.Name,
	}

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctxs, cacheKey, res)
	return res, nil
}

// ─────────────── ListSelect ──────────────────────────────────────────────────────
func (s *service) ListSelectKelurahanDesa(ctx context.Context, kecamatanID int64, search string) ([]dto.KelurahanDesaSimpelResponse, error) {
	ctxs := context.Background()

	if kecamatanID <= 0 {
		return nil, appErrors.BadRequest("ID kecamatan wajib diisi dan harus lebih dari 0")
	}

	cacheKey := cacheKeyDesaSelectList(kecamatanID, search)

	// 1. Cek Cache
	var cachedRes []dto.KelurahanDesaSimpelResponse
	if s.cache.Get(ctxs, cacheKey, &cachedRes) {
		return cachedRes, nil
	}

	// 2. Ambil dari DB
	items, err := s.repo.ListSelectKelurahanDesa(ctx, kecamatanID, search)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, appErrors.NotFound("Kelurahan/Desa tidak ditemukan")
	}

	res := dto.ToKelurahanDesaListSimpelResponse(items)

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctxs, cacheKey, res)
	return res, nil
}

// ─────────────── List ────────────────────────────────────────────────────────────
func (s *service) ListKelurahanDesa(ctx context.Context, page, pageSize int, kecamatanID *int64, filter *dto.FilterKelurahanDesaRequest) ([]dto.KelurahanDesaResponse, int64, error) {
	ctxs := context.Background()
	cacheKey := cacheKeyDesaList(page, pageSize, kecamatanID, filter)

	// 1. Cek Cache
	var cachedRes struct {
		Items []dto.KelurahanDesaResponse `json:"items"`
		Total int64                       `json:"total"`
	}
	if s.cache.Get(ctxs, cacheKey, &cachedRes) {
		if len(cachedRes.Items) == 0 {
			return nil, 0, appErrors.NotFound("Kelurahan/Desa tidak ditemukan")
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

	items, total, err := s.repo.ListKelurahanDesa(ctx, page, pageSize, kecamatanID, filter)
	if err != nil {
		return nil, 0, err
	}
	if len(items) == 0 {
		return nil, 0, appErrors.NotFound("Kelurahan/Desa tidak ditemukan")
	}

	res := dto.ToKelurahanDesaListResponse(items)

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctxs, cacheKey, struct {
		Items []dto.KelurahanDesaResponse `json:"items"`
		Total int64                       `json:"total"`
	}{Items: res, Total: total})

	return res, total, nil
}

// ─────────────── Create ──────────────────────────────────────────────────────────
func (s *service) CreateKelurahanDesa(ctx context.Context, req *dto.CreateKelurahanDesaRequest, actor he.AuthContext) (*dto.KelurahanDesaResponse, error) {
	// Permission
	can, err := s.canCreateMasterAlamat(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Kelurahan/Desa baru.", nil)
	}

	// Cek duplikat code
	exists, err := s.repo.ExistsKelurahanDesaByCode(ctx, req.Code, nil)
	if err != nil {
		return nil, appErrors.Internal("gagal cek duplikat code")
	}
	if exists {
		return nil, appErrors.Wrap(http.StatusConflict,
			"Code kelurahan/desa '"+req.Code+"' sudah digunakan.", nil)
	}

	// Logic
	kecamatan, err := s.repo.GetByIDKecamatan(ctx, req.KecamatanID)
	if err != nil {
		return nil, err
	}
	if kecamatan == nil {
		return nil, appErrors.NotFound("Kecamatan tidak ditemukan")
	}

	m := &models.MasterAlamatKelurahanDesa{
		KecamatanID: req.KecamatanID,
		Code:        req.Code,
		Name:        req.Name,
		FhirCode:    req.FhirCode,
		PostalCode:  req.PostalCode,
		CreatedBy:   &actor.UserID,
		UpdatedBy:   &actor.UserID,
	}
	if err := s.repo.CreateKelurahanDesa(ctx, m); err != nil {
		return nil, err
	}

	res := dto.ToKelurahanDesaResponse(m)

	// Invalidate Cache
	s.cache.InvalidateList(context.Background(), cachePrefixDesaList)
	s.cache.InvalidateList(context.Background(), cachePrefixDesaSelectList)

	return res, nil
}

// ─────────────── Update ──────────────────────────────────────────────────────────
func (s *service) UpdateKelurahanDesa(ctx context.Context, id int64, req *dto.UpdateKelurahanDesaRequest, actor he.AuthContext) (*dto.KelurahanDesaResponse, error) {
	// Permission
	can, err := s.canUpdateMasterAlamat(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Kelurahan/Desa.", nil)
	}

	// Logic
	m, err := s.repo.GetByIDKelurahanDesa(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Kelurahan/Desa tidak ditemukan")
	}

	if req.KecamatanID != nil {
		m.KecamatanID = *req.KecamatanID
	}
	if req.Code != nil {
		// Cek duplikat code, exclude diri sendiri
		exists, err := s.repo.ExistsKelurahanDesaByCode(ctx, *req.Code, &id)
		if err != nil {
			return nil, appErrors.Internal("gagal cek duplikat code")
		}
		if exists {
			return nil, appErrors.Wrap(http.StatusConflict,
				"Code kelurahan/desa '"+*req.Code+"' sudah digunakan.", nil)
		}
		m.Code = *req.Code
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.PostalCode != nil {
		m.PostalCode = req.PostalCode
	}
	if req.FhirCode != nil {
		if *req.FhirCode == "" {
			m.FhirCode = nil
		} else {
			m.FhirCode = req.FhirCode
		}
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateKelurahanDesa(ctx, m); err != nil {
		return nil, err
	}

	res := dto.ToKelurahanDesaResponse(m)

	// Invalidate Cache
	ctxs := context.Background()
	s.cache.InvalidateDetail(ctxs, cacheKeyDesaDetail(id))
	s.cache.InvalidateDetail(ctxs, cacheKeyDesaGetDetail(id))
	s.cache.InvalidateList(ctxs, cachePrefixDesaSelectList)
	s.cache.InvalidateList(ctxs, cachePrefixDesaList)

	return res, nil
}

// ─────────────── Delete ──────────────────────────────────────────────────────────
func (s *service) DeleteKelurahanDesa(ctx context.Context, id int64, actor he.AuthContext) error {
	// Permission
	can, err := s.canDeleteMasterAlamat(ctx, actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Kelurahan/Desa.", nil)
	}

	// Logic
	m, err := s.repo.GetByIDKelurahanDesa(ctx, id)
	if err != nil {
		return err
	}
	if m == nil {
		return appErrors.NotFound("Kelurahan/Desa tidak ditemukan")
	}

	err = s.repo.DeleteKelurahanDesa(ctx, id, actor.UserID)
	if err == nil {
		// Invalidate Cache
		ctxs := context.Background()
		s.cache.InvalidateDetail(ctxs, cacheKeyDesaDetail(id))
		s.cache.InvalidateDetail(ctxs, cacheKeyDesaGetDetail(id))
		s.cache.InvalidateList(ctxs, cachePrefixDesaList)
		s.cache.InvalidateList(ctxs, cachePrefixDesaSelectList)
	}
	return err
}
