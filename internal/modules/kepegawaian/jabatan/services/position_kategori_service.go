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
func (s *service) CreatePositionKategori(ctx context.Context, req *dto.CreatePositionKategoriRequest, actor he.AuthContext) (*dto.PositionKategoriResponse, error) {
	can, err := s.canCreatePositionKategori(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses: " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat PositionKategori baru.", nil)
	}

	// Check Duplicate code
	data, err := s.repo.GetPositionKategoriByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return nil, appErrors.Wrap(http.StatusConflict, "Kategori Posisi dengan kode ini sudah ada", nil)
	}

	// Check Duplicate label
	data, err = s.repo.GetPositionKategoriByLabel(ctx, req.Label)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return nil, appErrors.Wrap(http.StatusConflict, "Kategori Posisi dengan label ini sudah ada", nil)
	}

	// Buat instance model Kategori Posisi baru dengan bidang-bidang yang disesuaikan
	m := &models.PositionKategori{
		Code:      req.Code,
		Label:     req.Label,
		FHIRCode:  req.FHIRCode,
		Point:     req.Point,
		CreatedBy: &actor.UserID,
		UpdatedBy: &actor.UserID,
	}

	if err := s.repo.CreatePositionKategori(ctx, m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, m.CreatedBy)

	// Invalidate Cache
	s.cache.InvalidateList(context.Background(), cachePrefixPositionsKategoriSelectList)

	return dto.ToPositionKategoriResponse(dto.PositionKategoriResponseParams{
		PositionKategori: m,
		Creator:          creator,
		Updater:          creator,
	}), nil
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (s *service) GetPositionKategoriByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.PositionKategoriResponse, error) {
	can, err := s.canReadPositionKategori(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses: " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat PositionKategori.", nil)
	}

	m, err := s.repo.GetPositionKategoriByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("PositionKategori tidak ditemukan")
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToPositionKategoriResponse(dto.PositionKategoriResponseParams{
		PositionKategori: m,
		Creator:          creator,
		Updater:          updater,
	}), nil
}

// ── GetByCode ─────────────────────────────────────────────────────────────────
func (s *service) GetPositionKategoriByCode(ctx context.Context, code string, actor he.AuthContext) (*dto.PositionKategoriResponse, error) {
	can, err := s.canReadPositionKategori(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses: " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat PositionKategori.", nil)
	}

	m, err := s.repo.GetPositionKategoriByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("PositionKategori tidak ditemukan")
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToPositionKategoriResponse(dto.PositionKategoriResponseParams{
		PositionKategori: m,
		Creator:          creator,
		Updater:          updater,
	}), nil
}

// ── GetByLabel ────────────────────────────────────────────────────────────────
func (s *service) GetPositionKategoriByLabel(ctx context.Context, label string, actor he.AuthContext) (*dto.PositionKategoriResponse, error) {
	can, err := s.canReadPositionKategori(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses: " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat PositionKategori.", nil)
	}

	m, err := s.repo.GetPositionKategoriByLabel(ctx, label)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("PositionKategori tidak ditemukan")
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToPositionKategoriResponse(dto.PositionKategoriResponseParams{
		PositionKategori: m,
		Creator:          creator,
		Updater:          updater,
	}), nil
}

// ── ListSelect ────────────────────────────────────────────────────────────────
func (s *service) ListSelectPositionKategori(ctx context.Context, search string, actor he.AuthContext) ([]dto.PositionKategoriSelectResponse, error) {
	can, err := s.canReadPositionKategori(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar Tipe.", nil)
	}

	ctxs := context.Background()
	cacheKey := cacheKeyPositionsKategoriSelectList(search)

	// 1. Cek Cache
	var cachedRes []dto.PositionKategoriSelectResponse
	if s.cache.Get(ctxs, cacheKey, &cachedRes) {
		return cachedRes, nil
	}

	items, err := s.repo.ListSelectPositionKategori(ctx, search)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, appErrors.Wrap(http.StatusNotFound, "data tipe tidak ditemukan", nil)
	}

	res := dto.ToPositionKategoriSelectResponse(items)
	s.cache.SetDefault(ctxs, cacheKey, res)
	return res, nil
}

// ── List ──────────────────────────────────────────────────────────────────────
func (s *service) ListPositionKategori(ctx context.Context, page, pageSize int, filter *dto.FilterPositionKategoriRequest, actor he.AuthContext) ([]dto.PositionKategoriResponse, int64, error) {
	can, err := s.canReadPositionKategori(ctx, actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar PositionKategori.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}
	items, total, err := s.repo.ListPositionKategori(ctx, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMapsForPositionKategori(ctx, items)
	return dto.ToPositionKategoriListResponse(items, creatorsMap, updatersMap), total, nil
}

// ── Update ────────────────────────────────────────────────────────────────────
func (s *service) UpdatePositionKategori(ctx context.Context, id int64, req *dto.UpdatePositionKategoriRequest, actor he.AuthContext) (*dto.PositionKategoriResponse, error) {
	can, err := s.canUpdatePositionKategori(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah PositionKategori.", nil)
	}

	m, err := s.repo.GetPositionKategoriByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("PositionKategori tidak ditemukan")
	}

	// Check Duplicate code jika ada perubahan
	if req.Code != nil && *req.Code != m.Code {
		data, err := s.repo.GetPositionKategoriByCode(ctx, *req.Code)
		if err != nil {
			return nil, err
		}
		if data != nil {
			return nil, appErrors.Wrap(http.StatusConflict, "Kategori Posisi dengan kode ini sudah digunakan", nil)
		}
	}

	// Check Duplicate label jika ada perubahan
	if req.Label != nil && *req.Label != m.Label {
		data, err := s.repo.GetPositionKategoriByLabel(ctx, *req.Label)
		if err != nil {
			return nil, err
		}
		if data != nil {
			return nil, appErrors.Wrap(http.StatusConflict, "Kategori Posisi dengan label ini sudah digunakan", nil)
		}
	}

	// Update parsial sesuai pointer/field pada DTO Update
	if req.Code != nil {
		m.Code = *req.Code
	}
	if req.Label != nil {
		m.Label = *req.Label
	}
	if req.FHIRCode != nil {
		m.FHIRCode = req.FHIRCode
	}
	if req.Point != nil {
		m.Point = req.Point
	}

	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdatePositionKategori(ctx, m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToPositionKategoriResponse(dto.PositionKategoriResponseParams{
		PositionKategori: m,
		Creator:          creator,
		Updater:          updater,
	}), nil
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (s *service) DeletePositionKategori(ctx context.Context, id int64, actor he.AuthContext) error {
	can, err := s.canDeletePositionKategori(ctx, actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus PositionKategori.", nil)
	}

	m, err := s.repo.GetPositionKategoriByID(ctx, id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("PositionKategori tidak ditemukan")
	}
	return s.repo.DeletePositionKategori(ctx, id, actor.UserID)
}

// ── helper khusus PositionKategori (nama fungsi unik agar tidak bentrok) ───────

func (s *service) buildAuditMapsForPositionKategori(ctx context.Context, items []models.PositionKategori) (map[int64]*he.UserData, map[int64]*he.UserData) {
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
