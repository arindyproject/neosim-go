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
func (s *service) CreateJobTitle(ctx context.Context,req *dto.CreateJobTitleRequest, actor he.AuthContext) (*dto.JobTitleResponse, error) {
	can, err := s.canCreateJobTitle(ctx,actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses: " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat JobTitle baru.", nil)
	}

	m := &models.JobTitle{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   &actor.UserID,
		UpdatedBy:   &actor.UserID,
	}
	if err := s.repo.CreateJobTitle(ctx,m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx,m.CreatedBy)

	return dto.ToJobTitleResponse(dto.JobTitleResponseParams{
		JobTitle: m,
		Creator:       creator,
		Updater:       creator, // saat create, creator dan updater sama
	}), nil
}


// ── GetByID ───────────────────────────────────────────────────────────────────
func (s *service) GetJobTitleByID(ctx context.Context,id int64, actor he.AuthContext) (*dto.JobTitleResponse, error) {
	can, err := s.canReadJobTitle(ctx,actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses: " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat JobTitle.", nil)
	}

	m, err := s.repo.GetJobTitleByID(ctx,id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("JobTitle tidak ditemukan")
	}

	creator := s.buildCreator(ctx,m.CreatedBy)
	updater := s.buildCreator(ctx,m.UpdatedBy)

	return dto.ToJobTitleResponse(dto.JobTitleResponseParams{
		JobTitle: m,
		Creator:       creator,
		Updater:       updater,
	}), nil
}


// ── List ──────────────────────────────────────────────────────────────────────	
func (s *service) ListJobTitle(ctx context.Context,page, pageSize int, filter *dto.FilterJobTitleRequest, actor he.AuthContext) ([]dto.JobTitleResponse, int64, error) {
	can, err := s.canReadJobTitle(ctx, actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar JobTitle.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}
	items, total, err := s.repo.ListJobTitle(ctx,page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMapsForJobTitle(ctx,items)
	return dto.ToJobTitleListResponse(items, creatorsMap, updatersMap), total, nil
}


// ── Update ────────────────────────────────────────────────────────────────────
func (s *service) UpdateJobTitle(ctx context.Context,id int64, req *dto.UpdateJobTitleRequest, actor he.AuthContext) (*dto.JobTitleResponse, error) {
	can, err := s.canUpdateJobTitle(ctx,actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah JobTitle.", nil)
	}

	m, err := s.repo.GetJobTitleByID(ctx,id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("JobTitle tidak ditemukan")
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateJobTitle(ctx,m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx,m.CreatedBy)
	updater := s.buildCreator(ctx,m.UpdatedBy)

	return dto.ToJobTitleResponse(dto.JobTitleResponseParams{
		JobTitle: m,
		Creator:       creator,
		Updater:       updater,
	}), nil
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (s *service) DeleteJobTitle(ctx context.Context,id int64, actor he.AuthContext) error {
	can, err := s.canDeleteJobTitle(ctx,actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus JobTitle.", nil)
	}

	m, err := s.repo.GetJobTitleByID(ctx,id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("JobTitle tidak ditemukan")
	}
	return s.repo.DeleteJobTitle(ctx,id, actor.UserID)
}

// ── helper khusus JobTitle (nama fungsi unik agar tidak bentrok) ───────

func (s *service) buildAuditMapsForJobTitle(ctx context.Context,items []models.JobTitle) (map[int64]*he.UserData, map[int64]*he.UserData) {
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


