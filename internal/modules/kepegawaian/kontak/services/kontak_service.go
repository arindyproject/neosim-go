package services

import (
	"errors"
	"net/http"
	"time"

	"neosim_go/internal/modules/kepegawaian/kontak/dto"
	"neosim_go/internal/modules/kepegawaian/kontak/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

func (s *service) Create(req *dto.CreateKepegawaianKontakRequest, createdBy *int64, actor he.AuthContext) (*dto.KepegawaianKontakResponse, error) {
	can, err := s.canCreateKepegawaianKontak(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat KepegawaianKontak baru.", nil)
	}

	m := &models.KepegawaianKontak{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
	}
	if err := s.repo.Create(m); err != nil {
		return nil, err
	}
	return dto.ToKepegawaianKontakResponse(m), nil
}

func (s *service) GetByID(id int64, actor he.AuthContext) (*dto.KepegawaianKontakResponse, error) {
	can, err := s.canReadKepegawaianKontak(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk Melihat KepegawaianKontak.", nil)
	}

	m, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("KepegawaianKontak tidak ditemukan")
	}
	return dto.ToKepegawaianKontakResponse(m), nil
}

func (s *service) List(page, pageSize int, filter *dto.FilterKepegawaianKontakRequest, actor he.AuthContext) ([]dto.KepegawaianKontakResponse, int64, error) {
	can, err := s.canReadKepegawaianKontak(actor)
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
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	items, total, err := s.repo.List(page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}
	return dto.ToKepegawaianKontakListResponse(items), total, nil
}

func (s *service) Update(id int64, req *dto.UpdateKepegawaianKontakRequest, updatedBy *int64, actor he.AuthContext) (*dto.KepegawaianKontakResponse, error) {
	can, err := s.canUpdateKepegawaianKontak(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah KepegawaianKontak.", nil)
	}

	m, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("KepegawaianKontak tidak ditemukan")
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
	return dto.ToKepegawaianKontakResponse(m), nil
}

func (s *service) Delete(id int64, actor he.AuthContext) error {
	can, err := s.canDeleteKepegawaianKontak(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus KepegawaianKontak.", nil)
	}

	m, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("KepegawaianKontak tidak ditemukan")
	}
	return s.repo.Delete(id)
}
