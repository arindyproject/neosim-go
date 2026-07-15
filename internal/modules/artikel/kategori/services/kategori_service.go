package services

import (
	"errors"
	"net/http"
	"time"

	"neosim_go/internal/modules/artikel/kategori/dto"
	"neosim_go/internal/modules/artikel/kategori/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

func (s *service) Create(req *dto.CreateArtikelKategoriRequest, actor he.AuthContext) (*dto.ArtikelKategoriResponse, error) {
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
	if err := s.repo.Create(m); err != nil {
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

func (s *service) GetByID(id int64, actor he.AuthContext) (*dto.ArtikelKategoriResponse, error) {
	can, err := s.canReadArtikelKategori(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk Melihat ArtikelKategori.", nil)
	}

	m, err := s.repo.GetByID(id)
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

func (s *service) List(page, pageSize int, filter *dto.FilterArtikelKategoriRequest, actor he.AuthContext) ([]dto.ArtikelKategoriResponse, int64, error) {
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
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	items, total, err := s.repo.List(page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMaps(items)
	return toArtikelKategoriResponses(items, creatorsMap, updatersMap), total, nil
}

func (s *service) Update(id int64, req *dto.UpdateArtikelKategoriRequest, actor he.AuthContext) (*dto.ArtikelKategoriResponse, error) {
	can, err := s.canUpdateArtikelKategori(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah ArtikelKategori.", nil)
	}

	m, err := s.repo.GetByID(id)
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

	if err := s.repo.Update(m); err != nil {
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

func (s *service) Delete(id int64, actor he.AuthContext) error {
	can, err := s.canDeleteArtikelKategori(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus ArtikelKategori.", nil)
	}

	m, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("ArtikelKategori tidak ditemukan")
	}
	return s.repo.Delete(id)
}
