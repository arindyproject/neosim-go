package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"neosim_go/internal/modules/master/departemen/dto"
	"neosim_go/internal/modules/master/departemen/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

// ── Create ────────────────────────────────────────────────────────────────────
func (s *service) CreateDepartemen(ctx context.Context,req *dto.CreateMasterDepartemenRequest, actor he.AuthContext) (*dto.MasterDepartemenResponse, error) {
	can, err := s.canCreateMasterDepartemen(ctx,actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat MasterDepartemen baru.", nil)
	}

	m := &models.MasterDepartemen{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   &actor.UserID,
		UpdatedBy:   &actor.UserID,
	}
	if err := s.repo.CreateDepartemen(ctx,m); err != nil {
		return nil, err
	}
	
	creator := s.buildCreator(ctx,m.CreatedBy)

	return dto.ToMasterDepartemenResponse(dto.MasterDepartemenResponseParams{
		MasterDepartemen: m,
		Creator:    creator,
		Updater:    creator,
	}), nil
}


// ── GetByID ───────────────────────────────────────────────────────────────────
func (s *service) GetDepartemenByID(ctx context.Context,id int64, actor he.AuthContext) (*dto.MasterDepartemenResponse, error) {
	can, err := s.canReadMasterDepartemen(ctx,actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk Melihat MasterDepartemen.", nil)
	}

	m, err := s.repo.GetDepartemenByID(ctx,id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("MasterDepartemen tidak ditemukan")
	}
	
	creator := s.buildCreator(ctx,m.CreatedBy)
	updater := s.buildCreator(ctx,m.UpdatedBy)

	return dto.ToMasterDepartemenResponse(dto.MasterDepartemenResponseParams{
		MasterDepartemen: m,
		Creator:    creator,
		Updater:    updater,
	}), nil
}


// ── List ──────────────────────────────────────────────────────────────────────
func (s *service) ListDepartemen(ctx context.Context,page, pageSize int, filter *dto.FilterMasterDepartemenRequest, actor he.AuthContext) ([]dto.MasterDepartemenResponse, int64, error) {
	can, err := s.canReadMasterDepartemen(ctx,actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar MasterDepartemen.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}
	items, total, err := s.repo.ListDepartemen(ctx,page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMaps(ctx, items)
	return dto.ToMasterDepartemenListResponse(items, creatorsMap, updatersMap), total, nil
}


// ── Update ────────────────────────────────────────────────────────────────────
func (s *service) UpdateDepartemen(ctx context.Context,id int64, req *dto.UpdateMasterDepartemenRequest, actor he.AuthContext) (*dto.MasterDepartemenResponse, error) {
	can, err := s.canUpdateMasterDepartemen(ctx,actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah MasterDepartemen.", nil)
	}

	m, err := s.repo.GetDepartemenByID(ctx,id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("MasterDepartemen tidak ditemukan")
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateDepartemen(ctx,m); err != nil {
		return nil, err
	}
	
	creator := s.buildCreator(ctx,m.CreatedBy)
	updater := s.buildCreator(ctx,m.UpdatedBy)

	return dto.ToMasterDepartemenResponse(dto.MasterDepartemenResponseParams{
		MasterDepartemen: m,
		Creator:    creator,
		Updater:    updater,
	}), nil
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (s *service) DeleteDepartemen(ctx context.Context,id int64, actor he.AuthContext) error {
	can, err := s.canDeleteMasterDepartemen(ctx,actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus MasterDepartemen.", nil)
	}

	m, err := s.repo.GetDepartemenByID(ctx,id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("MasterDepartemen tidak ditemukan")
	}
	return s.repo.DeleteDepartemen(ctx,id)
}
