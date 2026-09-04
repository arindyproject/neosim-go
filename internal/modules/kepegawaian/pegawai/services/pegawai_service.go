package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"neosim_go/internal/modules/kepegawaian/pegawai/dto"
	"neosim_go/internal/modules/kepegawaian/pegawai/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

// ── Create ────────────────────────────────────────────────────────────────────
func (s *service) CreatePegawai(ctx context.Context, req *dto.CreateKepegawaianPegawaiRequest, actor he.AuthContext) (*dto.KepegawaianPegawaiResponse, error) {
	can, err := s.canCreateKepegawaianPegawai(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat KepegawaianPegawai baru.", nil)
	}

	m := &models.KepegawaianPegawai{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   &actor.UserID,
		UpdatedBy:   &actor.UserID,
	}
	if err := s.repo.CreatePegawai(ctx, m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, m.CreatedBy)

	return dto.ToKepegawaianPegawaiResponse(dto.KepegawaianPegawaiResponseParams{
		KepegawaianPegawai: m,
		Creator:            creator,
		Updater:            creator,
	}), nil
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (s *service) GetPegawaiByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.KepegawaianPegawaiResponse, error) {
	can, err := s.canReadKepegawaianPegawai(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk Melihat KepegawaianPegawai.", nil)
	}

	m, err := s.repo.GetPegawaiByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("KepegawaianPegawai tidak ditemukan")
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToKepegawaianPegawaiResponse(dto.KepegawaianPegawaiResponseParams{
		KepegawaianPegawai: m,
		Creator:            creator,
		Updater:            updater,
	}), nil
}

// ── List ──────────────────────────────────────────────────────────────────────
func (s *service) ListPegawai(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianPegawaiRequest, actor he.AuthContext) ([]dto.KepegawaianPegawaiResponse, int64, error) {
	can, err := s.canReadKepegawaianPegawai(ctx, actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar KepegawaianPegawai.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}
	items, total, err := s.repo.ListPegawai(ctx, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMaps(ctx, items)
	return dto.ToKepegawaianPegawaiListResponse(items, creatorsMap, updatersMap), total, nil
}

// ── Update ────────────────────────────────────────────────────────────────────
func (s *service) UpdatePegawai(ctx context.Context, id int64, req *dto.UpdateKepegawaianPegawaiRequest, actor he.AuthContext) (*dto.KepegawaianPegawaiResponse, error) {
	can, err := s.canUpdateKepegawaianPegawai(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah KepegawaianPegawai.", nil)
	}

	m, err := s.repo.GetPegawaiByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("KepegawaianPegawai tidak ditemukan")
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdatePegawai(ctx, m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToKepegawaianPegawaiResponse(dto.KepegawaianPegawaiResponseParams{
		KepegawaianPegawai: m,
		Creator:            creator,
		Updater:            updater,
	}), nil
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (s *service) DeletePegawai(ctx context.Context, id int64, actor he.AuthContext) error {
	can, err := s.canDeleteKepegawaianPegawai(ctx, actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus KepegawaianPegawai.", nil)
	}

	m, err := s.repo.GetPegawaiByID(ctx, id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("KepegawaianPegawai tidak ditemukan")
	}
	return s.repo.DeletePegawai(ctx, id, actor.UserID)
}
