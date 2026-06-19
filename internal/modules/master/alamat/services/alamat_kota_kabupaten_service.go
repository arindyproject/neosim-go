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

// Kota/Kabupaten ==================================================================

// ─────────────── GetByID ─────────────────────────────────────────────────────────
func (s *service) GetByIDKotaKabupaten(id int64) (*dto.KotaKabupatenDetailResponse, error) {
	ctx := context.Background()
	cacheKey := cacheKeyKotaDetail(id)

	// 1. Cek Cache
	var cachedRes dto.KotaKabupatenDetailResponse
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		return &cachedRes, nil
	}

	// 2. Ambil dari DB
	m, err := s.repo.GetByIDKotaKabupaten(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Kota/Kabupaten tidak ditemukan")
	}

	totalKecamatan, err := s.repo.CountKecamatanByKotaID(m.ID)
	if err != nil {
		return nil, err
	}
	totalDesa, err := s.repo.CountDesaByKotaID(m.ID)
	if err != nil {
		return nil, err
	}

	res := &dto.KotaKabupatenDetailResponse{
		ID:             m.ID,
		Code:           m.Code,
		Name:           m.Name,
		ProvinsiID:     m.ProvinsiID,
		ProvinsiName:   m.Provinsi.Name,
		NegaraID:       m.Provinsi.NegaraID,
		NegaraName:     m.Provinsi.Negara.Name,
		TotalKecamatan: totalKecamatan,
		TotalDesa:      totalDesa,
	}

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctx, cacheKey, res)
	return res, nil
}

// ─────────────── List ────────────────────────────────────────────────────────────
func (s *service) ListKotaKabupaten(page, pageSize int, provinsiID *int64, filter *dto.FilterKotaKabupatenRequest) ([]dto.KotaKabupatenResponse, int64, error) {
	ctx := context.Background()
	cacheKey := cacheKeyKotaList(page, pageSize, provinsiID, filter)

	// 1. Cek Cache
	var cachedRes struct {
		Items []dto.KotaKabupatenResponse `json:"items"`
		Total int64                       `json:"total"`
	}
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		if len(cachedRes.Items) == 0 {
			return nil, 0, appErrors.NotFound("Kota/Kabupaten tidak ditemukan")
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

	items, total, err := s.repo.ListKotaKabupaten(page, pageSize, provinsiID, filter)
	if err != nil {
		return nil, 0, err
	}
	if len(items) == 0 {
		return nil, 0, appErrors.NotFound("Kota/Kabupaten tidak ditemukan")
	}

	res := dto.ToKotaKabupatenListResponse(items)

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctx, cacheKey, struct {
		Items []dto.KotaKabupatenResponse `json:"items"`
		Total int64                       `json:"total"`
	}{Items: res, Total: total})

	return res, total, nil
}

// ─────────────── Create ──────────────────────────────────────────────────────────
func (s *service) CreateKotaKabupaten(req *dto.CreateKotaKabupatenRequest, actor he.AuthContext) (*dto.KotaKabupatenResponse, error) {
	// Permission
	can, err := s.canCreateMasterAlamat(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Kota/Kabupaten baru.", nil)
	}

	// Cek duplikat code
	exists, err := s.repo.ExistsKotaKabupatenByCode(req.Code, nil)
	if err != nil {
		return nil, appErrors.Internal("gagal cek duplikat code")
	}
	if exists {
		return nil, appErrors.Wrap(http.StatusConflict,
			"Code kota/kabupaten '"+req.Code+"' sudah digunakan.", nil)
	}

	// Logic
	provinsi, err := s.repo.GetByIDProvinsi(req.ProvinsiID)
	if err != nil {
		return nil, err
	}
	if provinsi == nil {
		return nil, appErrors.NotFound("Provinsi tidak ditemukan")
	}

	m := &models.MasterAlamatKotaKabupaten{
		ProvinsiID: req.ProvinsiID,
		Code:       req.Code,
		Name:       req.Name,
		CreatedBy:  &actor.UserID,
		UpdatedBy:  &actor.UserID,
	}
	if err := s.repo.CreateKotaKabupaten(m); err != nil {
		return nil, err
	}

	res := dto.ToKotaKabupatenResponse(m)

	// Invalidate Cache
	ctx := context.Background()
	s.cache.InvalidateList(ctx, cachePrefixKotaList)
	s.cache.InvalidateList(ctx, cachePrefixKecamatanList)
	s.cache.InvalidateList(ctx, cachePrefixDesaList)

	return res, nil
}

// ─────────────── Update ──────────────────────────────────────────────────────────
func (s *service) UpdateKotaKabupaten(id int64, req *dto.UpdateKotaKabupatenRequest, actor he.AuthContext) (*dto.KotaKabupatenResponse, error) {
	// Permission
	can, err := s.canUpdateMasterAlamat(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Kota/Kabupaten.", nil)
	}

	// Logic
	m, err := s.repo.GetByIDKotaKabupaten(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.NotFound("Kota/Kabupaten tidak ditemukan")
	}

	if req.ProvinsiID != nil {
		m.ProvinsiID = *req.ProvinsiID
	}
	if req.Code != nil {
		// Cek duplikat code, exclude diri sendiri
		exists, err := s.repo.ExistsKotaKabupatenByCode(*req.Code, &id)
		if err != nil {
			return nil, appErrors.Internal("gagal cek duplikat code")
		}
		if exists {
			return nil, appErrors.Wrap(http.StatusConflict,
				"Code kota/kabupaten '"+*req.Code+"' sudah digunakan.", nil)
		}
		m.Code = *req.Code
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateKotaKabupaten(m); err != nil {
		return nil, err
	}

	res := dto.ToKotaKabupatenResponse(m)

	// Invalidate Cache
	ctx := context.Background()
	s.cache.InvalidateDetail(ctx, cacheKeyKotaDetail(id))
	s.cache.InvalidateDetail(ctx, cacheKeyKotaGetDetail(id))
	s.cache.InvalidateList(ctx, cachePrefixKotaList)
	s.cache.InvalidateList(ctx, cachePrefixKecamatanList)
	s.cache.InvalidateList(ctx, cachePrefixDesaList)

	return res, nil
}

// ─────────────── Delete ──────────────────────────────────────────────────────────
func (s *service) DeleteKotaKabupaten(id int64, actor he.AuthContext) error {
	// Permission
	can, err := s.canDeleteMasterAlamat(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Kota/Kabupaten.", nil)
	}

	// Logic
	m, err := s.repo.GetByIDKotaKabupaten(id)
	if err != nil {
		return err
	}
	if m == nil {
		return appErrors.NotFound("Kota/Kabupaten tidak ditemukan")
	}

	err = s.repo.DeleteKotaKabupaten(id)
	if err == nil {
		// Invalidate Cache
		ctx := context.Background()
		s.cache.InvalidateDetail(ctx, cacheKeyKotaDetail(id))
		s.cache.InvalidateDetail(ctx, cacheKeyKotaGetDetail(id))
		s.cache.InvalidateList(ctx, cachePrefixKotaList)
		s.cache.InvalidateList(ctx, cachePrefixKecamatanList)
		s.cache.InvalidateList(ctx, cachePrefixDesaList)
	}
	return err
}
