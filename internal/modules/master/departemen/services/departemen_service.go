package services

import (
	"errors"
	"net/http"
	"time"

	"neosim_go/internal/modules/master/departemen/dto"
	"neosim_go/internal/modules/master/departemen/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

func (s *service) Create(req *dto.CreateMasterDepartemenRequest, createdBy *int64, actor he.AuthContext) (*dto.MasterDepartemenResponse, error) {
	can, err := s.canCreateMasterDepartemen(actor)
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
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
	}
	if err := s.repo.Create(m); err != nil {
		return nil, err
	}
	return dto.ToMasterDepartemenResponse(m), nil
}

func (s *service) GetByID(id int64, actor he.AuthContext) (*dto.MasterDepartemenResponse, error) {
	can, err := s.canReadMasterDepartemen(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk Melihat MasterDepartemen.", nil)
	}

	m, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("MasterDepartemen tidak ditemukan")
	}
	return dto.ToMasterDepartemenResponse(m), nil
}

func (s *service) List(page, pageSize int, filter *dto.FilterMasterDepartemenRequest, actor he.AuthContext) ([]dto.MasterDepartemenResponse, int64, error) {
	can, err := s.canReadMasterDepartemen(actor)
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
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	items, total, err := s.repo.List(page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}
	return dto.ToMasterDepartemenListResponse(items), total, nil
}

func (s *service) Update(id int64, req *dto.UpdateMasterDepartemenRequest, updatedBy *int64, actor he.AuthContext) (*dto.MasterDepartemenResponse, error) {
	can, err := s.canUpdateMasterDepartemen(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah MasterDepartemen.", nil)
	}

	m, err := s.repo.GetByID(id)
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
	m.UpdatedBy = updatedBy
	m.UpdatedAt = time.Now()

	if err := s.repo.Update(m); err != nil {
		return nil, err
	}
	return dto.ToMasterDepartemenResponse(m), nil
}

func (s *service) Delete(id int64, actor he.AuthContext) error {
	can, err := s.canDeleteMasterDepartemen(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus MasterDepartemen.", nil)
	}

	m, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("MasterDepartemen tidak ditemukan")
	}
	return s.repo.Delete(id)
}
