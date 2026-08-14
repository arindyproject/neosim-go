package services

import (
	"errors"
	"net/http"
	"time"

	"neosim_go/internal/modules/kepegawaian/identifier/dto"
	"neosim_go/internal/modules/kepegawaian/identifier/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

func (s *service) CreateIdentifier(req *dto.CreateKepegawaianIdentifierRequest, actor he.AuthContext) (*dto.KepegawaianIdentifierResponse, error) {
	can, err := s.canCreateKepegawaianIdentifier(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat KepegawaianIdentifier baru.", nil)
	}

	// 1. Validasi keberadaan master Tipe
	tipeMaster, err := s.repo.GetTipeByID(req.TipeID)
	if err != nil {
		return nil, err
	}
	if tipeMaster == nil {
		return nil, errors.New("Tipe identifier tidak ditemukan")
	}

	m := &models.KepegawaianIdentifier{
		PegawaiID:      req.PegawaiID,
		TipeID:         req.TipeID,
		Nilai:          req.Nilai,
		TanggalTerbit:  req.TanggalTerbit.ToTimePtr(),
		TanggalExpired: req.TanggalExpired.ToTimePtr(),
		IsPrimary:      req.IsPrimary,
		IsAktif:        req.IsAktif,
		CreatedBy:      &actor.UserID,
		UpdatedBy:      &actor.UserID,
	}

	if err := s.repo.CreateIdentifier(m); err != nil {
		return nil, err
	}

	// Attach Tipe master untuk kebutuhan response mapping
	m.Tipe = tipeMaster

	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.ToKepegawaianIdentifierResponse(dto.KepegawaianIdentifierResponseParams{
		KepegawaianIdentifier: m,
		Creator:               creator,
		Updater:               updater,
	}), nil
}

func (s *service) GetIdentifierByID(id int64, actor he.AuthContext) (*dto.KepegawaianIdentifierResponse, error) {
	can, err := s.canReadKepegawaianIdentifier(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk Melihat KepegawaianIdentifier.", nil)
	}

	m, err := s.repo.GetIdentifierByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("KepegawaianIdentifier tidak ditemukan")
	}

	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.ToKepegawaianIdentifierResponse(dto.KepegawaianIdentifierResponseParams{
		KepegawaianIdentifier: m,
		Creator:               creator,
		Updater:               updater,
	}), nil
}

func (s *service) ListIdentifier(page, pageSize int, filter *dto.FilterKepegawaianIdentifierRequest, actor he.AuthContext) ([]dto.KepegawaianIdentifierResponse, int64, error) {
	can, err := s.canReadKepegawaianIdentifier(actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar KepegawaianIdentifier.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	items, total, err := s.repo.ListIdentifier(page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMaps(items)
	return dto.ToKepegawaianIdentifierListResponse(items, creatorsMap, updatersMap), total, nil
}

func (s *service) UpdateIdentifier(id int64, req *dto.UpdateKepegawaianIdentifierRequest, actor he.AuthContext) (*dto.KepegawaianIdentifierResponse, error) {
	can, err := s.canUpdateKepegawaianIdentifier(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah KepegawaianIdentifier.", nil)
	}

	m, err := s.repo.GetIdentifierByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("KepegawaianIdentifier tidak ditemukan")
	}

	// Update parsial jika pointer dikirimkan (not nil)
	if req.TipeID != nil {
		tipeMaster, err := s.repo.GetTipeByID(*req.TipeID)
		if err != nil {
			return nil, err
		}
		if tipeMaster == nil {
			return nil, errors.New("Tipe identifier tidak ditemukan")
		}
		m.TipeID = *req.TipeID
		m.Tipe = tipeMaster
	}
	if req.Nilai != nil {
		m.Nilai = *req.Nilai
	}

	if req.TanggalTerbit != nil {
		m.TanggalTerbit = req.TanggalTerbit.ToTimePtr()
	}
	if req.TanggalExpired != nil {
		m.TanggalExpired = req.TanggalExpired.ToTimePtr()
	}
	if req.IsPrimary != nil {
		m.IsPrimary = *req.IsPrimary
	}
	if req.IsAktif != nil {
		m.IsAktif = *req.IsAktif
	}

	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateIdentifier(m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.ToKepegawaianIdentifierResponse(dto.KepegawaianIdentifierResponseParams{
		KepegawaianIdentifier: m,
		Creator:               creator,
		Updater:               updater,
	}), nil
}

func (s *service) DeleteIdentifier(id int64, actor he.AuthContext) error {
	can, err := s.canDeleteKepegawaianIdentifier(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus KepegawaianIdentifier.", nil)
	}

	m, err := s.repo.GetIdentifierByID(id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("KepegawaianIdentifier tidak ditemukan")
	}
	return s.repo.DeleteIdentifier(id)
}
