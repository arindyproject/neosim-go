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

func (s *service) CreateKontak(req *dto.CreateKepegawaianKontakRequest, actor he.AuthContext) (*dto.KepegawaianKontakResponse, error) {
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
		CreatedBy:   &actor.UserID,
		UpdatedBy:   &actor.UserID,
	}
	if err := s.repo.CreateKontak(m); err != nil {
		return nil, err
	}
	
	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.ToKepegawaianKontakResponse(dto.KepegawaianKontakResponseParams{
		KepegawaianKontak: m,
		Creator:    creator,
		Updater:    updater,
	}), nil
}

func (s *service) GetKontakByID(id int64, actor he.AuthContext) (*dto.KepegawaianKontakResponse, error) {
	can, err := s.canReadKepegawaianKontak(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk Melihat KepegawaianKontak.", nil)
	}

	m, err := s.repo.GetKontakByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("KepegawaianKontak tidak ditemukan")
	}
	
	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.ToKepegawaianKontakResponse(dto.KepegawaianKontakResponseParams{
		KepegawaianKontak: m,
		Creator:    creator,
		Updater:    updater,
	}), nil
}

func (s *service) ListKontak(page, pageSize int, filter *dto.FilterKepegawaianKontakRequest, actor he.AuthContext) ([]dto.KepegawaianKontakResponse, int64, error) {
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
	items, total, err := s.repo.ListKontak(page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMaps(items)
	return toKepegawaianKontakResponses(items, creatorsMap, updatersMap), total, nil
}

func (s *service) UpdateKontak(id int64, req *dto.UpdateKepegawaianKontakRequest, actor he.AuthContext) (*dto.KepegawaianKontakResponse, error) {
	can, err := s.canUpdateKepegawaianKontak(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah KepegawaianKontak.", nil)
	}

	m, err := s.repo.GetKontakByID(id)
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
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateKontak(m); err != nil {
		return nil, err
	}
	
	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.ToKepegawaianKontakResponse(dto.KepegawaianKontakResponseParams{
		KepegawaianKontak: m,
		Creator:    creator,
		Updater:    updater,
	}), nil
}

func (s *service) DeleteKontak(id int64, actor he.AuthContext) error {
	can, err := s.canDeleteKepegawaianKontak(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus KepegawaianKontak.", nil)
	}

	m, err := s.repo.GetKontakByID(id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("KepegawaianKontak tidak ditemukan")
	}
	return s.repo.DeleteKontak(id)
}
