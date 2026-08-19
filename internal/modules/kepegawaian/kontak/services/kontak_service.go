package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"neosim_go/internal/modules/kepegawaian/kontak/dto"
	"neosim_go/internal/modules/kepegawaian/kontak/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

// ─── Create ───────────────────────────────────────────────────────────────────────────────
func (s *service) CreateKontak(ctx context.Context, req *dto.CreateKepegawaianKontakRequest, actor he.AuthContext) (*dto.KepegawaianKontakResponse, error) {
	can, err := s.canCreateKepegawaianKontak(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat KepegawaianKontak baru.", nil)
	}

	m := &models.KepegawaianKontak{
		PegawaiID:   req.PegawaiID,
		TipeID:      req.TipeID,
		Nilai:       req.Nilai,
		IsPrimary:   req.IsPrimary,
		IsAktif:     req.IsAktif,
		Description: req.Description,
		CreatedBy:   &actor.UserID,
		UpdatedBy:   &actor.UserID,
	}
	if err := s.repo.CreateKontak(ctx, m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, m.CreatedBy)

	return dto.ToKepegawaianKontakResponse(dto.KepegawaianKontakResponseParams{
		KepegawaianKontak: m,
		Creator:           creator,
		Updater:           creator,
	}), nil
}

// ─── GetByID ───────────────────────────────────────────────────────────────────────────────
func (s *service) GetKontakByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.KepegawaianKontakResponse, error) {
	can, err := s.canReadKepegawaianKontak(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk Melihat KepegawaianKontak.", nil)
	}

	m, err := s.repo.GetKontakByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("KepegawaianKontak tidak ditemukan")
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToKepegawaianKontakResponse(dto.KepegawaianKontakResponseParams{
		KepegawaianKontak: m,
		Creator:           creator,
		Updater:           updater,
	}), nil
}

// ─── GetByPegawaiID ───────────────────────────────────────────────────────────────────────────────
func (s *service) GetKontakByPegawaiID(ctx context.Context, pegawaiID int64, actor he.AuthContext) ([]dto.KepegawaianKontakResponse, error) {
	can, err := s.canReadKepegawaianKontak(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk Melihat KepegawaianKontak.", nil)
	}

	items, err := s.repo.GetKontakByPegawaiID(ctx, pegawaiID)
	if err != nil {
		return nil, err
	}
	if items == nil {
		return nil, errors.New("KepegawaianKontak tidak ditemukan")
	}

	creatorsMap, updatersMap := s.buildAuditMaps(ctx, items)
	return dto.ToKepegawaianKontakListResponse(items, creatorsMap, updatersMap), nil
}

// ─── List ───────────────────────────────────────────────────────────────────────────────
func (s *service) ListKontak(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianKontakRequest, actor he.AuthContext) ([]dto.KepegawaianKontakResponse, int64, error) {
	can, err := s.canReadKepegawaianKontak(ctx, actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar KepegawaianKontak.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}
	items, total, err := s.repo.ListKontak(ctx, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMaps(ctx, items)
	return dto.ToKepegawaianKontakListResponse(items, creatorsMap, updatersMap), total, nil
}

// ─── Update ───────────────────────────────────────────────────────────────────────────────
func (s *service) UpdateKontak(ctx context.Context, id int64, req *dto.UpdateKepegawaianKontakRequest, actor he.AuthContext) (*dto.KepegawaianKontakResponse, error) {
	can, err := s.canUpdateKepegawaianKontak(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah KepegawaianKontak.", nil)
	}

	m, err := s.repo.GetKontakByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("KepegawaianKontak tidak ditemukan")
	}
	if req.TipeID != nil {
		m.TipeID = *req.TipeID
	}
	if req.Nilai != nil {
		m.Nilai = *req.Nilai
	}
	if req.IsPrimary != nil {
		m.IsPrimary = *req.IsPrimary
	}
	if req.IsAktif != nil {
		m.IsAktif = *req.IsAktif
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateKontak(ctx, m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToKepegawaianKontakResponse(dto.KepegawaianKontakResponseParams{
		KepegawaianKontak: m,
		Creator:           creator,
		Updater:           updater,
	}), nil
}

// ─── Delete ───────────────────────────────────────────────────────────────────────────────
func (s *service) DeleteKontak(ctx context.Context, id int64, actor he.AuthContext) error {
	can, err := s.canDeleteKepegawaianKontak(ctx, actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus KepegawaianKontak.", nil)
	}

	m, err := s.repo.GetKontakByID(ctx, id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("KepegawaianKontak tidak ditemukan")
	}
	return s.repo.DeleteKontak(ctx, id)
}
