
package services

import (
	"errors"
	"time"
	"net/http"

	"neosim_go/internal/modules/artikel/ketegori/dto"
	"neosim_go/internal/modules/artikel/ketegori/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)






// ─── Service ────────────────────────────────────────────────────────────────────

// ------------ Create ------------------------------------------------------------
func (s *service) Create(req *dto.CreateArtikelKetegoriRequest, createdBy *int64, actor  he.AuthContext) (*dto.ArtikelKetegoriResponse, error) {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canCreateArtikelKetegori(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat ArtikelKetegori baru.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	m := &models.ArtikelKetegori{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
	}
	if err := s.repo.Create(m); err != nil {
		return nil, err
	}
	return dto.ToArtikelKetegoriResponse(m), nil
}

// ------------ GetByID -----------------------------------------------------------
func (s *service) GetByID(id int64, actor  he.AuthContext) (*dto.ArtikelKetegoriResponse, error) {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canReadArtikelKetegori(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk Melihat ArtikelKetegori.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	m, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("ArtikelKetegori tidak ditemukan")
	}
	return dto.ToArtikelKetegoriResponse(m), nil
}

// ------------ List --------------------------------------------------------------
func (s *service) List(page, pageSize int,filter *dto.FilterArtikelKetegoriRequest, actor  he.AuthContext) ([]dto.ArtikelKetegoriResponse, int64, error) {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canReadArtikelKetegori(actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar ArtikelKetegori.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
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
	return dto.ToArtikelKetegoriListResponse(items), total, nil
}

// ------------ Update ------------------------------------------------------------
func (s *service) Update(id int64, req *dto.UpdateArtikelKetegoriRequest, updatedBy *int64, actor  he.AuthContext) (*dto.ArtikelKetegoriResponse, error) {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canUpdateArtikelKetegori(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah ArtikelKetegori.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	m, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("ArtikelKetegori tidak ditemukan")
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	m.UpdatedBy = updatedBy
	m.UpdatedAt = time.Now()

	if err := s.repo.Update(m); err != nil {
		return nil, err
	}
	return dto.ToArtikelKetegoriResponse(m), nil
}

// ------------ Delete ------------------------------------------------------------
func (s *service) Delete(id int64, actor  he.AuthContext) error {
	// ─── Permission ─────────────────────────────────────────────────────────────
	can, err := s.canDeleteArtikelKetegori(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus ArtikelKetegori.", nil)
	}

	// ─── Logic ──────────────────────────────────────────────────────────────────
	m, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("ArtikelKetegori tidak ditemukan")
	}
	return s.repo.Delete(id)
}
