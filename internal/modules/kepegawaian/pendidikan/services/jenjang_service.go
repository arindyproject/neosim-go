package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"neosim_go/internal/modules/kepegawaian/pendidikan/dto"
	"neosim_go/internal/modules/kepegawaian/pendidikan/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

// Semua method di bawah ini ditempelkan ke struct 'service' yang sama dengan
// entitas utama (lihat services/service.go). s.repo, s.buildCreator, dan
// s.buildAuditMaps dipakai ulang langsung — tidak perlu field/param baru.

// ── Create ────────────────────────────────────────────────────────────────────
func (s *service) CreateJenjang(ctx context.Context, req *dto.CreateJenjangRequest, actor he.AuthContext) (*dto.JenjangResponse, error) {
	can, err := s.canCreateJenjang(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Jenjang baru.", nil)
	}

	m := &models.Jenjang{
		Code:       req.Code,
		Label:      req.Label,
		FHIRSystem: req.FHIRSystem,
		CreatedBy:  &actor.UserID,
		UpdatedBy:  &actor.UserID,
	}
	if err := s.repo.CreateJenjang(ctx, m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, m.CreatedBy)

	return dto.ToJenjangResponse(dto.JenjangResponseParams{
		Jenjang: m,
		Creator: creator,
		Updater: creator, // saat create, creator dan updater sama
	}), nil
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (s *service) GetJenjangByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.JenjangResponse, error) {
	can, err := s.canReadJenjang(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat Jenjang.", nil)
	}

	m, err := s.repo.GetJenjangByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("Jenjang tidak ditemukan")
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToJenjangResponse(dto.JenjangResponseParams{
		Jenjang: m,
		Creator: creator,
		Updater: updater,
	}), nil
}

// ── List ──────────────────────────────────────────────────────────────────────
func (s *service) ListJenjang(ctx context.Context, page, pageSize int, filter *dto.FilterJenjangRequest, actor he.AuthContext) ([]dto.JenjangResponse, int64, error) {
	can, err := s.canReadJenjang(ctx, actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar Jenjang.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}
	items, total, err := s.repo.ListJenjang(ctx, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMapsForJenjang(ctx, items)
	return dto.ToJenjangListResponse(items, creatorsMap, updatersMap), total, nil
}

// ── Update ────────────────────────────────────────────────────────────────────
func (s *service) UpdateJenjang(ctx context.Context, id int64, req *dto.UpdateJenjangRequest, actor he.AuthContext) (*dto.JenjangResponse, error) {
	can, err := s.canUpdateJenjang(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Jenjang.", nil)
	}

	m, err := s.repo.GetJenjangByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("Jenjang tidak ditemukan")
	}
	if req.Code != nil {
		m.Code = *req.Code
	}
	if req.Label != nil {
		m.Label = *req.Label
	}
	if req.FHIRSystem != nil {
		m.FHIRSystem = req.FHIRSystem
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateJenjang(ctx, m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToJenjangResponse(dto.JenjangResponseParams{
		Jenjang: m,
		Creator: creator,
		Updater: updater,
	}), nil
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (s *service) DeleteJenjang(ctx context.Context, id int64, actor he.AuthContext) error {
	can, err := s.canDeleteJenjang(ctx, actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Jenjang.", nil)
	}

	m, err := s.repo.GetJenjangByID(ctx, id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("Jenjang tidak ditemukan")
	}
	return s.repo.DeleteJenjang(ctx, id)
}

// ── helper khusus Jenjang (nama fungsi unik agar tidak bentrok) ───────

func (s *service) buildAuditMapsForJenjang(ctx context.Context, items []models.Jenjang) (map[int64]*he.UserData, map[int64]*he.UserData) {
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
