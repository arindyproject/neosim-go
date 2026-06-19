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
func (s *service) GetByIDKecamatan(id int64) (*dto.KecamatanDetailResponse, error) {
	ctx := context.Background()
	cacheKey := cacheKeyKecamatanDetail(id)

	// 1. Cek Cache
	var cachedRes dto.KecamatanDetailResponse
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		return &cachedRes, nil
	}

	// 2. Ambil dari DB
	m, err := s.repo.GetByIDKecamatan(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Kecamatan tidak ditemukan")
	}

	totalDesa, err := s.repo.CountDesaByKecamatanID(m.ID)
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
	s.cache.SetDefault(ctx, cacheKey, res)
	return res, nil
}

// ─────────────── List ────────────────────────────────────────────────────────────
func (s *service) ListKecamatan(page, pageSize int, kotaKabupatenID *int64, filter *dto.FilterKecamatanRequest) ([]dto.KecamatanResponse, int64, error) {
	ctx := context.Background()
	cacheKey := cacheKeyKecamatanList(page, pageSize, kotaKabupatenID, filter)

	// 1. Cek Cache
	var cachedRes struct {
		Items []dto.KecamatanResponse `json:"items"`
		Total int64                   `json:"total"`
	}
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
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

	items, total, err := s.repo.ListKecamatan(page, pageSize, kotaKabupatenID, filter)
	if err != nil {
		return nil, 0, err
	}
	if len(items) == 0 {
		return nil, 0, appErrors.NotFound("Kecamatan tidak ditemukan")
	}

	res := dto.ToKecamatanListResponse(items)

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctx, cacheKey, struct {
		Items []dto.KecamatanResponse `json:"items"`
		Total int64                   `json:"total"`
	}{Items: res, Total: total})

	return res, total, nil
}

// ─────────────── Create ──────────────────────────────────────────────────────────
func (s *service) CreateKecamatan(req *dto.CreateKecamatanRequest, actor he.AuthContext) (*dto.KecamatanResponse, error) {
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
	exists, err := s.repo.ExistsKecamatanByCode(req.Code, nil)
	if err != nil {
		return nil, appErrors.Internal("gagal cek duplikat code")
	}
	if exists {
		return nil, appErrors.Wrap(http.StatusConflict,
			"Code kecamatan '"+req.Code+"' sudah digunakan.", nil)
	}

	// Logic
	kota, err := s.repo.GetByIDKotaKabupaten(req.KotaKabupatenID)
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
	if err := s.repo.CreateKecamatan(m); err != nil {
		return nil, err
	}

	res := dto.ToKecamatanResponse(m)

	// Invalidate Cache
	ctx := context.Background()
	s.cache.InvalidateList(ctx, cachePrefixKecamatanList)
	s.cache.InvalidateList(ctx, cachePrefixDesaList)

	return res, nil
}

// ─────────────── Update ──────────────────────────────────────────────────────────
func (s *service) UpdateKecamatan(id int64, req *dto.UpdateKecamatanRequest, actor he.AuthContext) (*dto.KecamatanResponse, error) {
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
	m, err := s.repo.GetByIDKecamatan(id)
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
		exists, err := s.repo.ExistsKecamatanByCode(*req.Code, &id)
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

	if err := s.repo.UpdateKecamatan(m); err != nil {
		return nil, err
	}

	res := dto.ToKecamatanResponse(m)

	// Invalidate Cache
	ctx := context.Background()
	s.cache.InvalidateDetail(ctx, cacheKeyKecamatanDetail(id))
	s.cache.InvalidateDetail(ctx, cacheKeyKecamatanGetDetail(id))
	s.cache.InvalidateList(ctx, cachePrefixKecamatanList)
	s.cache.InvalidateList(ctx, cachePrefixDesaList)

	return res, nil
}

// ─────────────── Delete ──────────────────────────────────────────────────────────
func (s *service) DeleteKecamatan(id int64, actor he.AuthContext) error {
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
	m, err := s.repo.GetByIDKecamatan(id)
	if err != nil {
		return err
	}
	if m == nil {
		return appErrors.NotFound("Kecamatan tidak ditemukan")
	}

	err = s.repo.DeleteKecamatan(id)
	if err == nil {
		// Invalidate Cache
		ctx := context.Background()
		s.cache.InvalidateDetail(ctx, cacheKeyKecamatanDetail(id))
		s.cache.InvalidateDetail(ctx, cacheKeyKecamatanGetDetail(id))
		s.cache.InvalidateList(ctx, cachePrefixKecamatanList)
		s.cache.InvalidateList(ctx, cachePrefixDesaList)
	}
	return err
}
