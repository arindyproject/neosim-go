package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"neosim_go/internal/modules/kepegawaian/alamat/dto"
	"neosim_go/internal/modules/kepegawaian/alamat/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

// ── Create ────────────────────────────────────────────────────────────────────
func (s *service) CreateAlamat(ctx context.Context,req *dto.CreateKepegawaianAlamatRequest, actor he.AuthContext) (*dto.KepegawaianAlamatResponse, error) {
	can, err := s.canCreateKepegawaianAlamat(ctx,actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat KepegawaianAlamat baru.", nil)
	}

	m := &models.KepegawaianAlamat{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   &actor.UserID,
		UpdatedBy:   &actor.UserID,
	}
	if err := s.repo.CreateAlamat(ctx,m); err != nil {
		return nil, err
	}
	
	creator := s.buildCreator(ctx,m.CreatedBy)

	return dto.ToKepegawaianAlamatResponse(dto.KepegawaianAlamatResponseParams{
		KepegawaianAlamat: m,
		Creator:    creator,
		Updater:    creator,
	}), nil
}


// ── GetByID ───────────────────────────────────────────────────────────────────
func (s *service) GetAlamatByID(ctx context.Context,id int64, actor he.AuthContext) (*dto.KepegawaianAlamatResponse, error) {
	can, err := s.canReadKepegawaianAlamat(ctx,actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk Melihat KepegawaianAlamat.", nil)
	}

	m, err := s.repo.GetAlamatByID(ctx,id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("KepegawaianAlamat tidak ditemukan")
	}
	
	creator := s.buildCreator(ctx,m.CreatedBy)
	updater := s.buildCreator(ctx,m.UpdatedBy)

	return dto.ToKepegawaianAlamatResponse(dto.KepegawaianAlamatResponseParams{
		KepegawaianAlamat: m,
		Creator:    creator,
		Updater:    updater,
	}), nil
}


// ── List ──────────────────────────────────────────────────────────────────────
func (s *service) ListAlamat(ctx context.Context,page, pageSize int, filter *dto.FilterKepegawaianAlamatRequest, actor he.AuthContext) ([]dto.KepegawaianAlamatResponse, int64, error) {
	can, err := s.canReadKepegawaianAlamat(ctx,actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar KepegawaianAlamat.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}
	items, total, err := s.repo.ListAlamat(ctx,page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMaps(ctx, items)
	return dto.ToKepegawaianAlamatListResponse(items, creatorsMap, updatersMap), total, nil
}


// ── Update ────────────────────────────────────────────────────────────────────
func (s *service) UpdateAlamat(ctx context.Context,id int64, req *dto.UpdateKepegawaianAlamatRequest, actor he.AuthContext) (*dto.KepegawaianAlamatResponse, error) {
	can, err := s.canUpdateKepegawaianAlamat(ctx,actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah KepegawaianAlamat.", nil)
	}

	m, err := s.repo.GetAlamatByID(ctx,id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("KepegawaianAlamat tidak ditemukan")
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateAlamat(ctx,m); err != nil {
		return nil, err
	}
	
	creator := s.buildCreator(ctx,m.CreatedBy)
	updater := s.buildCreator(ctx,m.UpdatedBy)

	return dto.ToKepegawaianAlamatResponse(dto.KepegawaianAlamatResponseParams{
		KepegawaianAlamat: m,
		Creator:    creator,
		Updater:    updater,
	}), nil
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (s *service) DeleteAlamat(ctx context.Context,id int64, actor he.AuthContext) error {
	can, err := s.canDeleteKepegawaianAlamat(ctx,actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus KepegawaianAlamat.", nil)
	}

	m, err := s.repo.GetAlamatByID(ctx,id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("KepegawaianAlamat tidak ditemukan")
	}
	return s.repo.DeleteAlamat(ctx,id, actor.UserID)
}
