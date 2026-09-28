package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

// Semua method di bawah ini ditempelkan ke struct 'service' yang sama dengan
// entitas utama (lihat services/service.go). s.repo, s.buildCreator, dan
// s.buildAuditMaps dipakai ulang langsung — tidak perlu field/param baru.

// ── Create ────────────────────────────────────────────────────────────────────
func (s *service) CreateSpecializationKategori(ctx context.Context, req *dto.CreateSpecializationKategoriRequest, actor he.AuthContext) (*dto.SpecializationKategoriResponse, error) {
	can, err := s.canCreateSpecializationKategori(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses: " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat SpecializationKategori baru.", nil)
	}

	// Check Duplicate code
	data, err := s.repo.GetSpecializationKategoriByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return nil, appErrors.Wrap(http.StatusConflict, "Specialization Kategori dengan kode ini sudah ada", nil)
	}

	// Check Duplicate label
	data, err = s.repo.GetSpecializationKategoriByLabel(ctx, req.Label)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return nil, appErrors.Wrap(http.StatusConflict, "Specialization Kategori dengan label ini sudah ada", nil)
	}

	m := &models.SpecializationKategori{
		Code:      req.Code,
		Label:     req.Label,
		FHIRCode:  req.FHIRCode,
		CreatedBy: &actor.UserID,
		UpdatedBy: &actor.UserID,
	}
	if err := s.repo.CreateSpecializationKategori(ctx, m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, m.CreatedBy)

	// Invalidate Cache
	s.cache.InvalidateList(context.Background(), cachePrefixSpecializationsKategoriSelectList)

	return dto.ToSpecializationKategoriResponse(dto.SpecializationKategoriResponseParams{
		SpecializationKategori: m,
		Creator:                creator,
		Updater:                creator, // saat create, creator dan updater sama
	}), nil
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (s *service) GetSpecializationKategoriByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.SpecializationKategoriResponse, error) {
	can, err := s.canReadSpecializationKategori(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses: " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat SpecializationKategori.", nil)
	}

	m, err := s.repo.GetSpecializationKategoriByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("SpecializationKategori tidak ditemukan")
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToSpecializationKategoriResponse(dto.SpecializationKategoriResponseParams{
		SpecializationKategori: m,
		Creator:                creator,
		Updater:                updater,
	}), nil
}

// ── ListSelect ────────────────────────────────────────────────────────────────
func (s *service) ListSelectSpecializationKategori(ctx context.Context, search string, actor he.AuthContext) ([]dto.SpecializationKategoriSelectResponse, error) {
	can, err := s.canReadSpecializationKategori(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar Tipe.", nil)
	}

	ctxs := context.Background()
	cacheKey := cacheKeySpecializationsKategoriSelectList(search)

	// 1. Cek Cache
	var cachedRes []dto.SpecializationKategoriSelectResponse
	if s.cache.Get(ctxs, cacheKey, &cachedRes) {
		return cachedRes, nil
	}

	items, err := s.repo.ListSelectSpecializationKategori(ctx, search)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, appErrors.Wrap(http.StatusNotFound, "data tipe tidak ditemukan", nil)
	}

	res := dto.ToSpecializationKategoriSelectResponse(items)
	s.cache.SetDefault(ctxs, cacheKey, res)
	return res, nil
}

// ── List ──────────────────────────────────────────────────────────────────────
func (s *service) ListSpecializationKategori(ctx context.Context, page, pageSize int, filter *dto.FilterSpecializationKategoriRequest, actor he.AuthContext) ([]dto.SpecializationKategoriResponse, int64, error) {
	can, err := s.canReadSpecializationKategori(ctx, actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar SpecializationKategori.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}
	items, total, err := s.repo.ListSpecializationKategori(ctx, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMapsForSpecializationKategori(ctx, items)
	return dto.ToSpecializationKategoriListResponse(items, creatorsMap, updatersMap), total, nil
}

// ── Update ────────────────────────────────────────────────────────────────────
func (s *service) UpdateSpecializationKategori(ctx context.Context, id int64, req *dto.UpdateSpecializationKategoriRequest, actor he.AuthContext) (*dto.SpecializationKategoriResponse, error) {
	can, err := s.canUpdateSpecializationKategori(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah SpecializationKategori.", nil)
	}

	//------------------------------------------------
	m, err := s.repo.GetSpecializationKategoriByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("SpecializationKategori tidak ditemukan")
	}

	//ceck duplicate code
	if req.Code != nil && *req.Code != m.Code {
		data, err := s.repo.GetSpecializationKategoriByCode(ctx, *req.Code)
		if err != nil {
			return nil, err
		}
		if data != nil {
			return nil, appErrors.Wrap(http.StatusConflict, "SpecializationKategori dengan kode ini sudah digunakan", nil)
		}
	}

	//ceck duplicate label
	if req.Label != nil && *req.Label != m.Label {
		data, err := s.repo.GetSpecializationKategoriByLabel(ctx, *req.Label)
		if err != nil {
			return nil, err
		}
		if data != nil {
			return nil, appErrors.Wrap(http.StatusConflict, "SpecializationKategori dengan label ini sudah digunakan", nil)
		}
	}

	// update fields yang diubah
	if req.Code != nil {
		m.Code = *req.Code
	}
	if req.Label != nil {
		m.Label = *req.Label
	}
	if req.FHIRCode != nil {
		m.FHIRCode = req.FHIRCode
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateSpecializationKategori(ctx, m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	res := dto.ToSpecializationKategoriResponse(dto.SpecializationKategoriResponseParams{
		SpecializationKategori: m,
		Creator:                creator,
		Updater:                updater,
	})

	ctxs := context.Background()
	s.cache.InvalidateList(ctxs, cachePrefixSpecializationsKategoriSelectList)
	return res, nil
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (s *service) DeleteSpecializationKategori(ctx context.Context, id int64, actor he.AuthContext) error {
	can, err := s.canDeleteSpecializationKategori(ctx, actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus SpecializationKategori.", nil)
	}

	m, err := s.repo.GetSpecializationKategoriByID(ctx, id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("SpecializationKategori tidak ditemukan")
	}
	err = s.repo.DeleteSpecializationKategori(ctx, id, actor.UserID)
	if err == nil {
		// Invalidate Cache
		ctxs := context.Background()
		s.cache.InvalidateList(ctxs, cachePrefixSpecializationsKategoriSelectList)
	}
	return err
}

// ── helper khusus SpecializationKategori (nama fungsi unik agar tidak bentrok) ───────

func (s *service) buildAuditMapsForSpecializationKategori(ctx context.Context, items []models.SpecializationKategori) (map[int64]*he.UserData, map[int64]*he.UserData) {
	idSet := make(map[int64]struct{})
	for _, item := range items {
		if item.CreatedBy != nil {
			idSet[*item.CreatedBy] = struct{}{}
		}
		if item.UpdatedBy != nil {
			idSet[*item.UpdatedBy] = struct{}{}
		}
	}
	ids := make([]int64, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}

	users, err := s.userRepo.GetByIDs(ctx, ids) // ← 1 query total, bukan 40
	if err != nil {
		return map[int64]*he.UserData{}, map[int64]*he.UserData{}
	}

	userMap := make(map[int64]*he.UserData, len(users))
	for _, u := range users {
		userMap[u.ID] = &he.UserData{ID: u.ID, Username: u.Username, Name: u.Name}
	}
	// creator dan updater sekarang share map yang sama — reuse otomatis, kode lebih pendek juga
	return userMap, userMap
}
