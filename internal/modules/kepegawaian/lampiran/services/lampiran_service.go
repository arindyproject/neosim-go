package services

import (
	"errors"
	"net/http"
	"time"

	"neosim_go/internal/modules/kepegawaian/lampiran/dto"
	"neosim_go/internal/modules/kepegawaian/lampiran/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

func (s *service) Create(req *dto.CreateKepegawaianLampiranRequest, createdBy *int64, actor he.AuthContext) (*dto.KepegawaianLampiranResponse, error) {
	can, err := s.canCreateKepegawaianLampiran(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat KepegawaianLampiran baru.", nil)
	}

	m := &models.KepegawaianLampiran{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
	}
	if err := s.repo.Create(m); err != nil {
		return nil, err
	}
	return dto.ToKepegawaianLampiranResponse(m), nil
}

func (s *service) GetByID(id int64, actor he.AuthContext) (*dto.KepegawaianLampiranResponse, error) {
	can, err := s.canReadKepegawaianLampiran(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk Melihat KepegawaianLampiran.", nil)
	}

	m, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("KepegawaianLampiran tidak ditemukan")
	}
	return dto.ToKepegawaianLampiranResponse(m), nil
}

func (s *service) List(page, pageSize int, filter *dto.FilterKepegawaianLampiranRequest, actor he.AuthContext) ([]dto.KepegawaianLampiranResponse, int64, error) {
	can, err := s.canReadKepegawaianLampiran(actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar KepegawaianLampiran.", nil)
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
	return dto.ToKepegawaianLampiranListResponse(items), total, nil
}

func (s *service) Update(id int64, req *dto.UpdateKepegawaianLampiranRequest, updatedBy *int64, actor he.AuthContext) (*dto.KepegawaianLampiranResponse, error) {
	can, err := s.canUpdateKepegawaianLampiran(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah KepegawaianLampiran.", nil)
	}

	m, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("KepegawaianLampiran tidak ditemukan")
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
	return dto.ToKepegawaianLampiranResponse(m), nil
}

func (s *service) Delete(id int64, actor he.AuthContext) error {
	can, err := s.canDeleteKepegawaianLampiran(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus KepegawaianLampiran.", nil)
	}

	m, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("KepegawaianLampiran tidak ditemukan")
	}
	return s.repo.Delete(id)
}
