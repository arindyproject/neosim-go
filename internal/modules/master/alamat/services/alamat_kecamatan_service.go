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

// Kecamatan =======================================================================

// ─────────────── GetByID ─────────────────────────────────────────────────────────
func (s *service) GetByIDKecamatan(ctx context.Context, id int64) (*dto.KecamatanDetailResponse, error) {
	ctxs := context.Background()
	cacheKey := cacheKeyKecamatanDetail(id)

	// 1. Cek Cache
	var cachedRes dto.KecamatanDetailResponse
	if s.cache.Get(ctxs, cacheKey, &cachedRes) {
		return &cachedRes, nil
	}

	// 2. Ambil dari DB
	m, err := s.repo.GetByIDKecamatan(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Kecamatan tidak ditemukan")
	}

	totalDesa, err := s.repo.CountDesaByKecamatanID(ctx, m.ID)
	if err != nil {
		return nil, err
	}

	res := &dto.KecamatanDetailResponse{
		ID:                m.ID,
		Code:              m.Code,
		Name:              m.Name,
		KotaKabupatenID:   m.KotaKabupatenID,
		KotaKabupatenName: m.KotaKabupaten.Name,
		ProvinsiID:        m.KotaKabupaten.ProvinsiID,
		ProvinsiName:      m.KotaKabupaten.Provinsi.Name,
		NegaraID:          m.KotaKabupaten.Provinsi.NegaraID,
		NegaraName:        m.KotaKabupaten.Provinsi.Negara.Name,
		TotalDesa:         totalDesa,
	}

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctxs, cacheKey, res)
	return res, nil
}

// ─────────────── List ────────────────────────────────────────────────────────────
func (s *service) ListKecamatan(ctx context.Context, page, pageSize int, kotaKabupatenID *int64, filter *dto.FilterKecamatanRequest) ([]dto.KecamatanResponse, int64, error) {
	ctxs := context.Background()
	cacheKey := cacheKeyKecamatanList(page, pageSize, kotaKabupatenID, filter)

	// 1. Cek Cache
	var cachedRes struct {
		Items []dto.KecamatanResponse `json:"items"`
		Total int64                   `json:"total"`
	}
	if s.cache.Get(ctxs, cacheKey, &cachedRes) {
		if len(cachedRes.Items) == 0 {
			return nil, 0, appErrors.NotFound("Kecamatan tidak ditemukan")
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

	items, total, err := s.repo.ListKecamatan(ctx, page, pageSize, kotaKabupatenID, filter)
	if err != nil {
		return nil, 0, err
	}
	if len(items) == 0 {
		return nil, 0, appErrors.NotFound("Kecamatan tidak ditemukan")
	}

	res := dto.ToKecamatanListResponse(items)

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctxs, cacheKey, struct {
		Items []dto.KecamatanResponse `json:"items"`
		Total int64                   `json:"total"`
	}{Items: res, Total: total})

	return res, total, nil
}

// ─────────────── Create ──────────────────────────────────────────────────────────
func (s *service) CreateKecamatan(ctx context.Context, req *dto.CreateKecamatanRequest, actor he.AuthContext) (*dto.KecamatanResponse, error) {
	// Permission
	can, err := s.canCreateMasterAlamat(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Kecamatan baru.", nil)
	}

	// Cek duplikat code
	exists, err := s.repo.ExistsKecamatanByCode(ctx, req.Code, nil)
	if err != nil {
		return nil, appErrors.Internal("gagal cek duplikat code")
	}
	if exists {
		return nil, appErrors.Wrap(http.StatusConflict,
			"Code kecamatan '"+req.Code+"' sudah digunakan.", nil)
	}

	// Logic
	kota, err := s.repo.GetByIDKotaKabupaten(ctx, req.KotaKabupatenID)
	if err != nil {
		return nil, err
	}
	if kota == nil {
		return nil, appErrors.NotFound("Kota/Kabupaten tidak ditemukan")
	}

	m := &models.MasterAlamatKecamatan{
		KotaKabupatenID: req.KotaKabupatenID,
		Code:            req.Code,
		Name:            req.Name,
		CreatedBy:       &actor.UserID,
		UpdatedBy:       &actor.UserID,
	}
	if err := s.repo.CreateKecamatan(ctx, m); err != nil {
		return nil, err
	}

	res := dto.ToKecamatanResponse(m)

	// Invalidate Cache
	ctxs := context.Background()
	s.cache.InvalidateList(ctxs, cachePrefixKecamatanList)
	s.cache.InvalidateList(ctxs, cachePrefixDesaList)

	return res, nil
}

// ─────────────── Update ──────────────────────────────────────────────────────────
func (s *service) UpdateKecamatan(ctx context.Context, id int64, req *dto.UpdateKecamatanRequest, actor he.AuthContext) (*dto.KecamatanResponse, error) {
	// Permission
	can, err := s.canUpdateMasterAlamat(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Kecamatan.", nil)
	}

	// Logic
	m, err := s.repo.GetByIDKecamatan(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Kecamatan tidak ditemukan")
	}

	if req.KotaKabupatenID != nil {
		m.KotaKabupatenID = *req.KotaKabupatenID
	}
	if req.Code != nil {
		// Cek duplikat code, exclude diri sendiri
		exists, err := s.repo.ExistsKecamatanByCode(ctx, *req.Code, &id)
		if err != nil {
			return nil, appErrors.Internal("gagal cek duplikat code")
		}
		if exists {
			return nil, appErrors.Wrap(http.StatusConflict,
				"Code kecamatan '"+*req.Code+"' sudah digunakan.", nil)
		}
		m.Code = *req.Code
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateKecamatan(ctx, m); err != nil {
		return nil, err
	}

	res := dto.ToKecamatanResponse(m)

	// Invalidate Cache
	ctxs := context.Background()
	s.cache.InvalidateDetail(ctxs, cacheKeyKecamatanDetail(id))
	s.cache.InvalidateDetail(ctxs, cacheKeyKecamatanGetDetail(id))
	s.cache.InvalidateList(ctxs, cachePrefixKecamatanList)
	s.cache.InvalidateList(ctxs, cachePrefixDesaList)

	return res, nil
}

// ─────────────── Delete ──────────────────────────────────────────────────────────
func (s *service) DeleteKecamatan(ctx context.Context, id int64, actor he.AuthContext) error {
	// Permission
	can, err := s.canDeleteMasterAlamat(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Kecamatan.", nil)
	}

	// Logic
	m, err := s.repo.GetByIDKecamatan(ctx, id)
	if err != nil {
		return err
	}
	if m == nil {
		return appErrors.NotFound("Kecamatan tidak ditemukan")
	}

	err = s.repo.DeleteKecamatan(ctx, id)
	if err == nil {
		// Invalidate Cache
		ctxs := context.Background()
		s.cache.InvalidateDetail(ctxs, cacheKeyKecamatanDetail(id))
		s.cache.InvalidateDetail(ctxs, cacheKeyKecamatanGetDetail(id))
		s.cache.InvalidateList(ctxs, cachePrefixKecamatanList)
		s.cache.InvalidateList(ctxs, cachePrefixDesaList)
	}
	return err
}
