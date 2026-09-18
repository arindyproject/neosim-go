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
func (s *service) CreateSpecialization(ctx context.Context,req *dto.CreateSpecializationRequest, actor he.AuthContext) (*dto.SpecializationResponse, error) {
	can, err := s.canCreateSpecialization(ctx,actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses: " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Specialization baru.", nil)
	}

	m := &models.Specialization{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   &actor.UserID,
		UpdatedBy:   &actor.UserID,
	}
	if err := s.repo.CreateSpecialization(ctx,m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx,m.CreatedBy)

	return dto.ToSpecializationResponse(dto.SpecializationResponseParams{
		Specialization: m,
		Creator:       creator,
		Updater:       creator, // saat create, creator dan updater sama
	}), nil
}


// ── GetByID ───────────────────────────────────────────────────────────────────
func (s *service) GetSpecializationByID(ctx context.Context,id int64, actor he.AuthContext) (*dto.SpecializationResponse, error) {
	can, err := s.canReadSpecialization(ctx,actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses: " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat Specialization.", nil)
	}

	m, err := s.repo.GetSpecializationByID(ctx,id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("Specialization tidak ditemukan")
	}

	creator := s.buildCreator(ctx,m.CreatedBy)
	updater := s.buildCreator(ctx,m.UpdatedBy)

	return dto.ToSpecializationResponse(dto.SpecializationResponseParams{
		Specialization: m,
		Creator:       creator,
		Updater:       updater,
	}), nil
}


// ── List ──────────────────────────────────────────────────────────────────────	
func (s *service) ListSpecialization(ctx context.Context,page, pageSize int, filter *dto.FilterSpecializationRequest, actor he.AuthContext) ([]dto.SpecializationResponse, int64, error) {
	can, err := s.canReadSpecialization(ctx, actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar Specialization.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}
	items, total, err := s.repo.ListSpecialization(ctx,page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMapsForSpecialization(ctx,items)
	return dto.ToSpecializationListResponse(items, creatorsMap, updatersMap), total, nil
}


// ── Update ────────────────────────────────────────────────────────────────────
func (s *service) UpdateSpecialization(ctx context.Context,id int64, req *dto.UpdateSpecializationRequest, actor he.AuthContext) (*dto.SpecializationResponse, error) {
	can, err := s.canUpdateSpecialization(ctx,actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Specialization.", nil)
	}

	m, err := s.repo.GetSpecializationByID(ctx,id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("Specialization tidak ditemukan")
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateSpecialization(ctx,m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx,m.CreatedBy)
	updater := s.buildCreator(ctx,m.UpdatedBy)

	return dto.ToSpecializationResponse(dto.SpecializationResponseParams{
		Specialization: m,
		Creator:       creator,
		Updater:       updater,
	}), nil
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (s *service) DeleteSpecialization(ctx context.Context,id int64, actor he.AuthContext) error {
	can, err := s.canDeleteSpecialization(ctx,actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Specialization.", nil)
	}

	m, err := s.repo.GetSpecializationByID(ctx,id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("Specialization tidak ditemukan")
	}
	return s.repo.DeleteSpecialization(ctx,id, actor.UserID)
}

// ── helper khusus Specialization (nama fungsi unik agar tidak bentrok) ───────

func (s *service) buildAuditMapsForSpecialization(ctx context.Context,items []models.Specialization) (map[int64]*he.UserData, map[int64]*he.UserData) {
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


