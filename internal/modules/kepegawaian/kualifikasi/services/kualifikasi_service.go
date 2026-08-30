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

// ── Create ────────────────────────────────────────────────────────────────────
func (s *service) CreateKualifikasi(ctx context.Context,req *dto.CreateKepegawaianKualifikasiRequest, actor he.AuthContext) (*dto.KepegawaianKualifikasiResponse, error) {
	can, err := s.canCreateKepegawaianKualifikasi(ctx,actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat KepegawaianKualifikasi baru.", nil)
	}

	m := &models.KepegawaianKualifikasi{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   &actor.UserID,
		UpdatedBy:   &actor.UserID,
	}
	if err := s.repo.CreateKualifikasi(ctx,m); err != nil {
		return nil, err
	}
	
	creator := s.buildCreator(ctx,m.CreatedBy)

	return dto.ToKepegawaianKualifikasiResponse(dto.KepegawaianKualifikasiResponseParams{
		KepegawaianKualifikasi: m,
		Creator:    creator,
		Updater:    creator,
	}), nil
}


// ── GetByID ───────────────────────────────────────────────────────────────────
func (s *service) GetKualifikasiByID(ctx context.Context,id int64, actor he.AuthContext) (*dto.KepegawaianKualifikasiResponse, error) {
	can, err := s.canReadKepegawaianKualifikasi(ctx,actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk Melihat KepegawaianKualifikasi.", nil)
	}

	m, err := s.repo.GetKualifikasiByID(ctx,id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("KepegawaianKualifikasi tidak ditemukan")
	}
	
	creator := s.buildCreator(ctx,m.CreatedBy)
	updater := s.buildCreator(ctx,m.UpdatedBy)

	return dto.ToKepegawaianKualifikasiResponse(dto.KepegawaianKualifikasiResponseParams{
		KepegawaianKualifikasi: m,
		Creator:    creator,
		Updater:    updater,
	}), nil
}


// ── List ──────────────────────────────────────────────────────────────────────
func (s *service) ListKualifikasi(ctx context.Context,page, pageSize int, filter *dto.FilterKepegawaianKualifikasiRequest, actor he.AuthContext) ([]dto.KepegawaianKualifikasiResponse, int64, error) {
	can, err := s.canReadKepegawaianKualifikasi(ctx,actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar KepegawaianKualifikasi.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}
	items, total, err := s.repo.ListKualifikasi(ctx,page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMaps(ctx, items)
	return dto.ToKepegawaianKualifikasiListResponse(items, creatorsMap, updatersMap), total, nil
}


// ── Update ────────────────────────────────────────────────────────────────────
func (s *service) UpdateKualifikasi(ctx context.Context,id int64, req *dto.UpdateKepegawaianKualifikasiRequest, actor he.AuthContext) (*dto.KepegawaianKualifikasiResponse, error) {
	can, err := s.canUpdateKepegawaianKualifikasi(ctx,actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah KepegawaianKualifikasi.", nil)
	}

	m, err := s.repo.GetKualifikasiByID(ctx,id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("KepegawaianKualifikasi tidak ditemukan")
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateKualifikasi(ctx,m); err != nil {
		return nil, err
	}
	
	creator := s.buildCreator(ctx,m.CreatedBy)
	updater := s.buildCreator(ctx,m.UpdatedBy)

	return dto.ToKepegawaianKualifikasiResponse(dto.KepegawaianKualifikasiResponseParams{
		KepegawaianKualifikasi: m,
		Creator:    creator,
		Updater:    updater,
	}), nil
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (s *service) DeleteKualifikasi(ctx context.Context,id int64, actor he.AuthContext) error {
	can, err := s.canDeleteKepegawaianKualifikasi(ctx,actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus KepegawaianKualifikasi.", nil)
	}

	m, err := s.repo.GetKualifikasiByID(ctx,id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("KepegawaianKualifikasi tidak ditemukan")
	}
	return s.repo.DeleteKualifikasi(ctx,id)
}
