package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"neosim_go/internal/modules/artikel/kategori/dto"
	"neosim_go/internal/modules/artikel/kategori/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

// Semua method di bawah ini ditempelkan ke struct 'service' yang sama dengan
// entitas utama (lihat services/service.go). s.repo, s.buildCreator, dan
// s.buildAuditMaps dipakai ulang langsung — tidak perlu field/param baru.

// ── Create ────────────────────────────────────────────────────────────────────
func (s *service) CreateTag(ctx context.Context,req *dto.CreateTagRequest, actor he.AuthContext) (*dto.TagResponse, error) {
	can, err := s.canCreateTag(ctx,actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Tag baru.", nil)
	}

	m := &models.Tag{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   &actor.UserID,
		UpdatedBy:   &actor.UserID,
	}
	if err := s.repo.CreateTag(ctx,m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx,m.CreatedBy)

	return dto.ToTagResponse(dto.TagResponseParams{
		Tag: m,
		Creator:       creator,
		Updater:       creator, // saat create, creator dan updater sama
	}), nil
}


// ── GetByID ───────────────────────────────────────────────────────────────────
func (s *service) GetTagByID(ctx context.Context,id int64, actor he.AuthContext) (*dto.TagResponse, error) {
	can, err := s.canReadTag(ctx,actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat Tag.", nil)
	}

	m, err := s.repo.GetTagByID(ctx,id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("Tag tidak ditemukan")
	}

	creator := s.buildCreator(ctx,m.CreatedBy)
	updater := s.buildCreator(ctx,m.UpdatedBy)

	return dto.ToTagResponse(dto.TagResponseParams{
		Tag: m,
		Creator:       creator,
		Updater:       updater,
	}), nil
}


// ── List ──────────────────────────────────────────────────────────────────────	
func (s *service) ListTag(ctx context.Context,page, pageSize int, filter *dto.FilterTagRequest, actor he.AuthContext) ([]dto.TagResponse, int64, error) {
	can, err := s.canReadTag(ctx, actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar Tag.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}
	items, total, err := s.repo.ListTag(ctx,page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMapsForTag(ctx,items)
	return dto.ToTagListResponse(items, creatorsMap, updatersMap), total, nil
}


// ── Update ────────────────────────────────────────────────────────────────────
func (s *service) UpdateTag(ctx context.Context,id int64, req *dto.UpdateTagRequest, actor he.AuthContext) (*dto.TagResponse, error) {
	can, err := s.canUpdateTag(ctx,actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Tag.", nil)
	}

	m, err := s.repo.GetTagByID(ctx,id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("Tag tidak ditemukan")
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateTag(ctx,m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx,m.CreatedBy)
	updater := s.buildCreator(ctx,m.UpdatedBy)

	return dto.ToTagResponse(dto.TagResponseParams{
		Tag: m,
		Creator:       creator,
		Updater:       updater,
	}), nil
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (s *service) DeleteTag(ctx context.Context,id int64, actor he.AuthContext) error {
	can, err := s.canDeleteTag(ctx,actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Tag.", nil)
	}

	m, err := s.repo.GetTagByID(ctx,id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("Tag tidak ditemukan")
	}
	return s.repo.DeleteTag(ctx,id)
}

// ── helper khusus Tag (nama fungsi unik agar tidak bentrok) ───────

func (s *service) buildAuditMapsForTag(ctx context.Context,items []models.Tag) (map[int64]*he.UserData, map[int64]*he.UserData) {
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

	users, err := s.userRepo.GetByIDs(ctx,ids) // ← 1 query total, bukan 40
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


