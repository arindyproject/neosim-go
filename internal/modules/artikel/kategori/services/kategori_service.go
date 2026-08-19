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

// ── Create ────────────────────────────────────────────────────────────────────
func (s *service) CreateKategori(ctx context.Context,req *dto.CreateArtikelKategoriRequest, actor he.AuthContext) (*dto.ArtikelKategoriResponse, error) {
	can, err := s.canCreateArtikelKategori(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat ArtikelKategori baru.", nil)
	}

	m := &models.ArtikelKategori{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   &actor.UserID,
		UpdatedBy:   &actor.UserID,
	}
	if err := s.repo.CreateKategori(ctx,m); err != nil {
		return nil, err
	}
	
	creator := s.buildCreator(m.CreatedBy)

	return dto.ToArtikelKategoriResponse(dto.ArtikelKategoriResponseParams{
		ArtikelKategori: m,
		Creator:    creator,
		Updater:    creator,
	}), nil
}


// ── GetByID ───────────────────────────────────────────────────────────────────
func (s *service) GetKategoriByID(ctx context.Context,id int64, actor he.AuthContext) (*dto.ArtikelKategoriResponse, error) {
	can, err := s.canReadArtikelKategori(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk Melihat ArtikelKategori.", nil)
	}

	m, err := s.repo.GetKategoriByID(ctx,id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("ArtikelKategori tidak ditemukan")
	}
	
	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.ToArtikelKategoriResponse(dto.ArtikelKategoriResponseParams{
		ArtikelKategori: m,
		Creator:    creator,
		Updater:    updater,
	}), nil
}


// ── List ──────────────────────────────────────────────────────────────────────
func (s *service) ListKategori(ctx context.Context,page, pageSize int, filter *dto.FilterArtikelKategoriRequest, actor he.AuthContext) ([]dto.ArtikelKategoriResponse, int64, error) {
	can, err := s.canReadArtikelKategori(actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar ArtikelKategori.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}
	items, total, err := s.repo.ListKategori(ctx,page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMaps(items)
	return dto.ToArtikelKategoriListResponse(items, creatorsMap, updatersMap), total, nil
}


// ── Update ────────────────────────────────────────────────────────────────────
func (s *service) UpdateKategori(ctx context.Context,id int64, req *dto.UpdateArtikelKategoriRequest, actor he.AuthContext) (*dto.ArtikelKategoriResponse, error) {
	can, err := s.canUpdateArtikelKategori(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah ArtikelKategori.", nil)
	}

	m, err := s.repo.GetKategoriByID(ctx,id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("ArtikelKategori tidak ditemukan")
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateKategori(ctx,m); err != nil {
		return nil, err
	}
	
	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.ToArtikelKategoriResponse(dto.ArtikelKategoriResponseParams{
		ArtikelKategori: m,
		Creator:    creator,
		Updater:    updater,
	}), nil
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (s *service) DeleteKategori(ctx context.Context,id int64, actor he.AuthContext) error {
	can, err := s.canDeleteArtikelKategori(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus ArtikelKategori.", nil)
	}

	m, err := s.repo.GetKategoriByID(ctx,id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("ArtikelKategori tidak ditemukan")
	}
	return s.repo.DeleteKategori(ctx,id)
}
