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

// ── Create ────────────────────────────────────────────────────────────────────
func (s *service) CreateJabatan(ctx context.Context,req *dto.CreateKepegawaianJabatanRequest, actor he.AuthContext) (*dto.KepegawaianJabatanResponse, error) {
	can, err := s.canCreateKepegawaianJabatan(ctx,actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses: " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat KepegawaianJabatan baru.", nil)
	}

	m := &models.KepegawaianJabatan{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   &actor.UserID,
		UpdatedBy:   &actor.UserID,
	}
	if err := s.repo.CreateJabatan(ctx,m); err != nil {
		return nil, err
	}
	
	creator := s.buildCreator(ctx,m.CreatedBy)

	return dto.ToKepegawaianJabatanResponse(dto.KepegawaianJabatanResponseParams{
		KepegawaianJabatan: m,
		Creator:    creator,
		Updater:    creator,
	}), nil
}


// ── GetByID ───────────────────────────────────────────────────────────────────
func (s *service) GetJabatanByID(ctx context.Context,id int64, actor he.AuthContext) (*dto.KepegawaianJabatanResponse, error) {
	can, err := s.canReadKepegawaianJabatan(ctx,actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses: " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk Melihat KepegawaianJabatan.", nil)
	}

	m, err := s.repo.GetJabatanByID(ctx,id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("KepegawaianJabatan tidak ditemukan")
	}
	
	creator := s.buildCreator(ctx,m.CreatedBy)
	updater := s.buildCreator(ctx,m.UpdatedBy)

	return dto.ToKepegawaianJabatanResponse(dto.KepegawaianJabatanResponseParams{
		KepegawaianJabatan: m,
		Creator:    creator,
		Updater:    updater,
	}), nil
}


// ── List ──────────────────────────────────────────────────────────────────────
func (s *service) ListJabatan(ctx context.Context,page, pageSize int, filter *dto.FilterKepegawaianJabatanRequest, actor he.AuthContext) ([]dto.KepegawaianJabatanResponse, int64, error) {
	can, err := s.canReadKepegawaianJabatan(ctx,actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar KepegawaianJabatan.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}
	items, total, err := s.repo.ListJabatan(ctx,page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMaps(ctx, items)
	return dto.ToKepegawaianJabatanListResponse(items, creatorsMap, updatersMap), total, nil
}


// ── Update ────────────────────────────────────────────────────────────────────
func (s *service) UpdateJabatan(ctx context.Context,id int64, req *dto.UpdateKepegawaianJabatanRequest, actor he.AuthContext) (*dto.KepegawaianJabatanResponse, error) {
	can, err := s.canUpdateKepegawaianJabatan(ctx,actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses: " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah KepegawaianJabatan.", nil)
	}

	m, err := s.repo.GetJabatanByID(ctx,id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("KepegawaianJabatan tidak ditemukan")
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateJabatan(ctx,m); err != nil {
		return nil, err
	}
	
	creator := s.buildCreator(ctx,m.CreatedBy)
	updater := s.buildCreator(ctx,m.UpdatedBy)

	return dto.ToKepegawaianJabatanResponse(dto.KepegawaianJabatanResponseParams{
		KepegawaianJabatan: m,
		Creator:    creator,
		Updater:    updater,
	}), nil
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (s *service) DeleteJabatan(ctx context.Context,id int64, actor he.AuthContext) error {
	can, err := s.canDeleteKepegawaianJabatan(ctx,actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus KepegawaianJabatan.", nil)
	}

	m, err := s.repo.GetJabatanByID(ctx,id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("KepegawaianJabatan tidak ditemukan")
	}
	return s.repo.DeleteJabatan(ctx,id, actor.UserID)
}
