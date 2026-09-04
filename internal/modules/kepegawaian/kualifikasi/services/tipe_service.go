package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"neosim_go/internal/modules/kepegawaian/kualifikasi/dto"
	"neosim_go/internal/modules/kepegawaian/kualifikasi/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

// Semua method di bawah ini ditempelkan ke struct 'service' yang sama dengan
// entitas utama (lihat services/service.go). s.repo, s.buildCreator, dan
// s.buildAuditMaps dipakai ulang langsung — tidak perlu field/param baru.

// ── Create ────────────────────────────────────────────────────────────────────
func (s *service) CreateTipe(ctx context.Context, req *dto.CreateTipeRequest, actor he.AuthContext) (*dto.TipeResponse, error) {
	can, err := s.canCreateTipe(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Tipe baru.", nil)
	}

	//ceck duplicate code
	data, err := s.repo.GetTipeByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return nil, appErrors.Wrap(http.StatusConflict, "Tipe dengan kode ini sudah ada", nil)
	}

	//ceck duplicate label
	data, err = s.repo.GetTipeByLabel(ctx, req.Label)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return nil, appErrors.Wrap(http.StatusConflict, "Tipe dengan label ini sudah ada", nil)
	}

	// Buat instance model Tipe baru dengan bidang-bidang yang disesuaikan
	m := &models.Tipe{
		Code:      req.Code,
		Label:     req.Label,
		CreatedBy: &actor.UserID,
		UpdatedBy: &actor.UserID,
	}
	if err := s.repo.CreateTipe(ctx, m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, m.CreatedBy)

	return dto.ToTipeResponse(dto.TipeResponseParams{
		Tipe:    m,
		Creator: creator,
		Updater: creator,
	}), nil
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (s *service) GetTipeByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.TipeResponse, error) {
	can, err := s.canReadTipe(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat Tipe.", nil)
	}

	m, err := s.repo.GetTipeByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("Tipe tidak ditemukan")
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToTipeResponse(dto.TipeResponseParams{
		Tipe:    m,
		Creator: creator,
		Updater: updater,
	}), nil
}

// ── GetByCode ─────────────────────────────────────────────────────────────────
func (s *service) GetTipeByCode(ctx context.Context, code string, actor he.AuthContext) (*dto.TipeResponse, error) {
	can, err := s.canReadTipe(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat Tipe.", nil)
	}

	m, err := s.repo.GetTipeByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("Tipe tidak ditemukan")
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToTipeResponse(dto.TipeResponseParams{
		Tipe:    m,
		Creator: creator,
		Updater: updater,
	}), nil
}

// ── GetByLabel ────────────────────────────────────────────────────────────────
func (s *service) GetTipeByLabel(ctx context.Context, label string, actor he.AuthContext) (*dto.TipeResponse, error) {
	can, err := s.canReadTipe(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat Tipe.", nil)
	}

	m, err := s.repo.GetTipeByLabel(ctx, label)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("Tipe tidak ditemukan")
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToTipeResponse(dto.TipeResponseParams{
		Tipe:    m,
		Creator: creator,
		Updater: updater,
	}), nil
}

// ── List ──────────────────────────────────────────────────────────────────────
func (s *service) ListTipe(ctx context.Context, page, pageSize int, filter *dto.FilterTipeRequest, actor he.AuthContext) ([]dto.TipeResponse, int64, error) {
	can, err := s.canReadTipe(ctx, actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar Tipe.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}
	items, total, err := s.repo.ListTipe(ctx, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMapsForTipe(ctx, items)
	return dto.ToTipeListResponse(items, creatorsMap, updatersMap), total, nil
}

// ── Update ────────────────────────────────────────────────────────────────────
func (s *service) UpdateTipe(ctx context.Context, id int64, req *dto.UpdateTipeRequest, actor he.AuthContext) (*dto.TipeResponse, error) {
	can, err := s.canUpdateTipe(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Tipe.", nil)
	}

	m, err := s.repo.GetTipeByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("Tipe tidak ditemukan")
	}

	//ceck duplicate code
	if req.Code != nil && *req.Code != m.Code {
		data, err := s.repo.GetTipeByCode(ctx, *req.Code)
		if err != nil {
			return nil, err
		}
		if data != nil {
			return nil, appErrors.Wrap(http.StatusConflict, "Tipe dengan kode ini sudah digunakan", nil)
		}
	}

	//ceck duplicate label
	if req.Label != nil && *req.Label != m.Label {
		data, err := s.repo.GetTipeByLabel(ctx, *req.Label)
		if err != nil {
			return nil, err
		}
		if data != nil {
			return nil, appErrors.Wrap(http.StatusConflict, "Tipe dengan label ini sudah digunakan", nil)
		}
	}

	// update fields yang diubah
	if req.Code != nil {
		m.Code = *req.Code
	}
	if req.Label != nil {
		m.Label = *req.Label
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateTipe(ctx, m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToTipeResponse(dto.TipeResponseParams{
		Tipe:    m,
		Creator: creator,
		Updater: updater,
	}), nil
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (s *service) DeleteTipe(ctx context.Context, id int64, actor he.AuthContext) error {
	can, err := s.canDeleteTipe(ctx, actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Tipe.", nil)
	}

	m, err := s.repo.GetTipeByID(ctx, id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("Tipe tidak ditemukan")
	}
	return s.repo.DeleteTipe(ctx, id, actor.UserID)
}

// ── helper khusus Tipe (nama fungsi unik agar tidak bentrok) ───────

func (s *service) buildAuditMapsForTipe(ctx context.Context, items []models.Tipe) (map[int64]*he.UserData, map[int64]*he.UserData) {
	fetchUser := func(id int64) (*he.UserData, error) {
		user, err := s.userRepo.GetByID(ctx, id)
		if err != nil || user == nil {
			return nil, err
		}
		return &he.UserData{ID: user.ID, Username: user.Username, Name: user.Name}, nil
	}

	creatorIDs := make(map[int64]struct{})
	updaterIDs := make(map[int64]struct{})
	for _, item := range items {
		if item.CreatedBy != nil {
			creatorIDs[*item.CreatedBy] = struct{}{}
		}
		if item.UpdatedBy != nil {
			updaterIDs[*item.UpdatedBy] = struct{}{}
		}
	}

	creatorsMap := make(map[int64]*he.UserData)
	for id := range creatorIDs {
		if data, err := fetchUser(id); err == nil && data != nil {
			creatorsMap[id] = data
		}
	}

	updatersMap := make(map[int64]*he.UserData)
	for id := range updaterIDs {
		if data, ok := creatorsMap[id]; ok {
			updatersMap[id] = data
		} else if data, err := fetchUser(id); err == nil && data != nil {
			updatersMap[id] = data
		}
	}

	return creatorsMap, updatersMap
}
