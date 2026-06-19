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
func (s *service) GetByIDKelurahanDesa(id int64) (*dto.KelurahanDesaDetailResponse, error) {
	ctx := context.Background()
	cacheKey := cacheKeyDesaDetail(id)

	// 1. Cek Cache
	var cachedRes dto.KelurahanDesaDetailResponse
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
		return &cachedRes, nil
	}

	// 2. Ambil dari DB
	m, err := s.repo.GetByIDKelurahanDesa(id)
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
	s.cache.SetDefault(ctx, cacheKey, res)
	return res, nil
}

// ─────────────── List ────────────────────────────────────────────────────────────
func (s *service) ListKelurahanDesa(page, pageSize int, kecamatanID *int64, filter *dto.FilterKelurahanDesaRequest) ([]dto.KelurahanDesaResponse, int64, error) {
	ctx := context.Background()
	cacheKey := cacheKeyDesaList(page, pageSize, kecamatanID, filter)

	// 1. Cek Cache
	var cachedRes struct {
		Items []dto.KelurahanDesaResponse `json:"items"`
		Total int64                       `json:"total"`
	}
	if s.cache.Get(ctx, cacheKey, &cachedRes) {
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

	items, total, err := s.repo.ListKelurahanDesa(page, pageSize, kecamatanID, filter)
	if err != nil {
		return nil, 0, err
	}
	if len(items) == 0 {
		return nil, 0, appErrors.NotFound("Kelurahan/Desa tidak ditemukan")
	}

	res := dto.ToKelurahanDesaListResponse(items)

	// 3. Simpan ke Cache
	s.cache.SetDefault(ctx, cacheKey, struct {
		Items []dto.KelurahanDesaResponse `json:"items"`
		Total int64                       `json:"total"`
	}{Items: res, Total: total})

	return res, total, nil
}

// ─────────────── Create ──────────────────────────────────────────────────────────
func (s *service) CreateKelurahanDesa(req *dto.CreateKelurahanDesaRequest, actor he.AuthContext) (*dto.KelurahanDesaResponse, error) {
	// Permission
	can, err := s.canCreateMasterAlamat(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Kelurahan/Desa baru.", nil)
	}

	// Cek duplikat code
	exists, err := s.repo.ExistsKelurahanDesaByCode(req.Code, nil)
	if err != nil {
		return nil, appErrors.Internal("gagal cek duplikat code")
	}
	if exists {
		return nil, appErrors.Wrap(http.StatusConflict,
			"Code kelurahan/desa '"+req.Code+"' sudah digunakan.", nil)
	}

	// Logic
	kecamatan, err := s.repo.GetByIDKecamatan(req.KecamatanID)
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
		PostalCode:  req.PostalCode,
		CreatedBy:   &actor.UserID,
		UpdatedBy:   &actor.UserID,
	}
	if err := s.repo.CreateKelurahanDesa(m); err != nil {
		return nil, err
	}

	res := dto.ToKelurahanDesaResponse(m)

	// Invalidate Cache
	s.cache.InvalidateList(context.Background(), cachePrefixDesaList)

	return res, nil
}

// ─────────────── Update ──────────────────────────────────────────────────────────
func (s *service) UpdateKelurahanDesa(id int64, req *dto.UpdateKelurahanDesaRequest, actor he.AuthContext) (*dto.KelurahanDesaResponse, error) {
	// Permission
	can, err := s.canUpdateMasterAlamat(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Kelurahan/Desa.", nil)
	}

	// Logic
	m, err := s.repo.GetByIDKelurahanDesa(id)
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
		exists, err := s.repo.ExistsKelurahanDesaByCode(*req.Code, &id)
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
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateKelurahanDesa(m); err != nil {
		return nil, err
	}

	res := dto.ToKelurahanDesaResponse(m)

	// Invalidate Cache
	ctx := context.Background()
	s.cache.InvalidateDetail(ctx, cacheKeyDesaDetail(id))
	s.cache.InvalidateDetail(ctx, cacheKeyDesaGetDetail(id))
	s.cache.InvalidateList(ctx, cachePrefixDesaList)

	return res, nil
}

// ─────────────── Delete ──────────────────────────────────────────────────────────
func (s *service) DeleteKelurahanDesa(id int64, actor he.AuthContext) error {
	// Permission
	can, err := s.canDeleteMasterAlamat(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Kelurahan/Desa.", nil)
	}

	// Logic
	m, err := s.repo.GetByIDKelurahanDesa(id)
	if err != nil {
		return err
	}
	if m == nil {
		return appErrors.NotFound("Kelurahan/Desa tidak ditemukan")
	}

	err = s.repo.DeleteKelurahanDesa(id)
	if err == nil {
		// Invalidate Cache
		ctx := context.Background()
		s.cache.InvalidateDetail(ctx, cacheKeyDesaDetail(id))
		s.cache.InvalidateDetail(ctx, cacheKeyDesaGetDetail(id))
		s.cache.InvalidateList(ctx, cachePrefixDesaList)
	}
	return err
}
