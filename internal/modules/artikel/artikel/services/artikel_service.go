package services

import (
	"errors"
	"net/http"
	"time"

	"neosim_go/internal/modules/artikel/artikel/dto"
	"neosim_go/internal/modules/artikel/artikel/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

func (s *service) CreateArtikel(req *dto.CreateArtikelRequest, actor he.AuthContext) (*dto.ArtikelResponse, error) {
	can, err := s.canCreateArtikel(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Artikel baru.", nil)
	}

	m := &models.Artikel{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   &actor.UserID,
		UpdatedBy:   &actor.UserID,
	}
	if err := s.repo.CreateArtikel(m); err != nil {
		return nil, err
	}
	
	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.ToArtikelResponse(dto.ArtikelResponseParams{
		Artikel: m,
		Creator:    creator,
		Updater:    updater,
	}), nil
}

func (s *service) GetArtikelByID(id int64, actor he.AuthContext) (*dto.ArtikelResponse, error) {
	can, err := s.canReadArtikel(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk Melihat Artikel.", nil)
	}

	m, err := s.repo.GetArtikelByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("Artikel tidak ditemukan")
	}
	
	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.ToArtikelResponse(dto.ArtikelResponseParams{
		Artikel: m,
		Creator:    creator,
		Updater:    updater,
	}), nil
}

func (s *service) ListArtikel(page, pageSize int, filter *dto.FilterArtikelRequest, actor he.AuthContext) ([]dto.ArtikelResponse, int64, error) {
	can, err := s.canReadArtikel(actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar Artikel.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	items, total, err := s.repo.ListArtikel(page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMaps(items)
	return toArtikelResponses(items, creatorsMap, updatersMap), total, nil
}

func (s *service) UpdateArtikel(id int64, req *dto.UpdateArtikelRequest, actor he.AuthContext) (*dto.ArtikelResponse, error) {
	can, err := s.canUpdateArtikel(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Artikel.", nil)
	}

	m, err := s.repo.GetArtikelByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("Artikel tidak ditemukan")
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateArtikel(m); err != nil {
		return nil, err
	}
	
	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.ToArtikelResponse(dto.ArtikelResponseParams{
		Artikel: m,
		Creator:    creator,
		Updater:    updater,
	}), nil
}

func (s *service) DeleteArtikel(id int64, actor he.AuthContext) error {
	can, err := s.canDeleteArtikel(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Artikel.", nil)
	}

	m, err := s.repo.GetArtikelByID(id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("Artikel tidak ditemukan")
	}
	return s.repo.DeleteArtikel(id)
}
