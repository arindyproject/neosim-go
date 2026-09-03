package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"neosim_go/internal/modules/kepegawaian/kualifikasi/dto"
	"neosim_go/internal/modules/kepegawaian/kualifikasi/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

// ── Create ────────────────────────────────────────────────────────────────────
func (s *service) CreateKualifikasi(ctx context.Context, req *dto.CreateKepegawaianKualifikasiRequest, actor he.AuthContext) (*dto.KepegawaianKualifikasiResponse, error) {
	can, err := s.canCreateKepegawaianKualifikasi(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat KepegawaianKualifikasi baru.", nil)
	}

	// validasi keberadaan master Tipe
	tipeMaster, err := s.repo.GetTipeByID(ctx, req.TipeID)
	if err != nil {
		return nil, appErrors.Internal("gagal mengambil master tipe identifier")
	}
	if tipeMaster == nil {
		return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "Tipe identifier tidak ditemukan.", nil)
	}

	// cek duplikasi nomor sertifikat + tipe
	//exists, err := s.repo.ExistsByNomorSertifikatAndTipe(ctx, req.TipeID, *req.NomorSertifikat, 0)
	//if err != nil {
	//	return nil, appErrors.Internal("gagal memeriksa duplikasi nomor sertifikat")
	//}
	//if exists {
	//	return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "Nomor sertifikat sudah digunakan.", nil)
	//}

	m := &models.KepegawaianKualifikasi{
		PegawaiID:       req.PegawaiID,
		TipeID:          req.TipeID,
		Nama:            req.Nama,
		Penyelenggara:   req.Penyelenggara,
		NomorSertifikat: req.NomorSertifikat,
		TanggalTerbit:   req.TanggalTerbit.ToTimePtr(),
		TanggalExpired:  req.TanggalExpired.ToTimePtr(),
		IsAktif:         req.IsAktif,
		FhirCode:        req.FhirCode,
		FhirSystem:      req.FhirSystem,
		CreatedBy:       &actor.UserID,
		UpdatedBy:       &actor.UserID,
	}
	if err := s.repo.CreateKualifikasi(ctx, m); err != nil {
		return nil, appErrors.Internal("gagal menyimpan kualifikasi pegawai")
	}

	m.Tipe = tipeMaster
	creator := s.buildCreator(ctx, m.CreatedBy)

	return dto.ToKepegawaianKualifikasiResponse(dto.KepegawaianKualifikasiResponseParams{
		KepegawaianKualifikasi: m,
		Creator:                creator,
		Updater:                creator,
	}), nil
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (s *service) GetKualifikasiByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.KepegawaianKualifikasiResponse, error) {
	can, err := s.canReadKepegawaianKualifikasi(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk Melihat KepegawaianKualifikasi.", nil)
	}

	m, err := s.repo.GetKualifikasiByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("KepegawaianKualifikasi tidak ditemukan")
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToKepegawaianKualifikasiResponse(dto.KepegawaianKualifikasiResponseParams{
		KepegawaianKualifikasi: m,
		Creator:                creator,
		Updater:                updater,
	}), nil
}

// ── List ──────────────────────────────────────────────────────────────────────
func (s *service) ListKualifikasi(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianKualifikasiRequest, actor he.AuthContext) ([]dto.KepegawaianKualifikasiResponse, int64, error) {
	can, err := s.canReadKepegawaianKualifikasi(ctx, actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar KepegawaianKualifikasi.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}
	items, total, err := s.repo.ListKualifikasi(ctx, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMaps(ctx, items)
	return dto.ToKepegawaianKualifikasiListResponse(items, creatorsMap, updatersMap), total, nil
}

// ── ListByPegawai ─────────────────────────────────────────────────────────────
func (s *service) ListByPegawai(
	ctx context.Context,
	pegawaiID int64,
	page, pageSize int,
	actor he.AuthContext,
) ([]dto.KepegawaianKualifikasiResponse, int64, error) {
	can, err := s.canReadKepegawaianKualifikasi(ctx, actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat identifier pegawai.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSize
	}

	items, total, err := s.repo.GetKualifikasiByPegawaiID(ctx, pegawaiID, page, pageSize)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal mengambil daftar kualifikasi pegawai")
	}

	creatorsMap, updatersMap := s.buildAuditMaps(ctx, items)
	return dto.ToKepegawaianKualifikasiListResponse(items, creatorsMap, updatersMap), total, nil
}

// ── Update ────────────────────────────────────────────────────────────────────
func (s *service) UpdateKualifikasi(ctx context.Context, id int64, req *dto.UpdateKepegawaianKualifikasiRequest, actor he.AuthContext) (*dto.KepegawaianKualifikasiResponse, error) {
	can, err := s.canUpdateKepegawaianKualifikasi(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah KepegawaianKualifikasi.", nil)
	}

	m, err := s.repo.GetKualifikasiByID(ctx, id)
	if err != nil {
		return nil, appErrors.Internal("gagal mengambil data kualifikasi")
	}
	if m == nil {
		return nil, appErrors.Wrap(http.StatusNotFound, "Kepegawaian Kualifikasi tidak ditemukan.", nil)
	}

	//----------------------
	tipeCheck := m.TipeID
	if req.TipeID != nil {
		tipeCheck = *req.TipeID
	}

	nomorSertifikatCheck := m.NomorSertifikat
	if req.NomorSertifikat != nil {
		nomorSertifikatCheck = req.NomorSertifikat
	}

	if nomorSertifikatCheck != nil {
		duplicate, err := s.repo.ExistsByNomorSertifikatAndTipe(ctx, tipeCheck, *nomorSertifikatCheck, id)
		if err != nil {
			return nil, appErrors.Internal("gagal cek duplikasi identifier")
		}
		if duplicate {
			return nil, appErrors.Wrap(http.StatusConflict,
				"Nilai identifier sudah pernah digunakan oleh Anda.", nil)
		}
	}

	// update parsial jika pointer dikirimkan (not nil)
	if req.TipeID != nil {
		tipeMaster, err := s.repo.GetTipeByID(ctx, *req.TipeID)
		if err != nil {
			return nil, appErrors.Internal("gagal mengambil master tipe kualifikasi")
		}
		if tipeMaster == nil {
			return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "Tipe kualifikasi tidak ditemukan.", nil)
		}
		m.TipeID = *req.TipeID
		m.Tipe = tipeMaster
	}

	if req.Nama != nil {
		m.Nama = *req.Nama
	}

	if req.Penyelenggara != nil {
		m.Penyelenggara = *req.Penyelenggara
	}

	if req.NomorSertifikat != nil {
		m.NomorSertifikat = req.NomorSertifikat
	}

	if req.IsAktif != nil {
		m.IsAktif = *req.IsAktif
	}

	if req.TanggalTerbit != nil {
		m.TanggalTerbit = req.TanggalTerbit.ToTimePtr()
	}
	if req.TanggalExpired != nil {
		m.TanggalExpired = req.TanggalExpired.ToTimePtr()
	}

	if req.FhirCode != nil {
		m.FhirCode = req.FhirCode
	}

	if req.FhirSystem != nil {
		m.FhirSystem = req.FhirSystem
	}

	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateKualifikasi(ctx, m); err != nil {
		return nil, appErrors.Internal("gagal mengupdate kualifikasi pegawai")
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToKepegawaianKualifikasiResponse(dto.KepegawaianKualifikasiResponseParams{
		KepegawaianKualifikasi: m,
		Creator:                creator,
		Updater:                updater,
	}), nil
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (s *service) DeleteKualifikasi(ctx context.Context, id int64, actor he.AuthContext) error {
	can, err := s.canDeleteKepegawaianKualifikasi(ctx, actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus KepegawaianKualifikasi.", nil)
	}

	m, err := s.repo.GetKualifikasiByID(ctx, id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("KepegawaianKualifikasi tidak ditemukan")
	}
	return s.repo.DeleteKualifikasi(ctx, id, actor.UserID)
}

// ── GetExpiringSoon ───────────────────────────────────────────────────────────
func (s *service) GetExpiringSoonKualifikasi(
	ctx context.Context,
	days int,
	page, pageSize int,
	actor he.AuthContext,
) ([]dto.KepegawaianKualifikasiResponse, int64, error) {
	can, err := s.canReadKepegawaianKualifikasi(ctx, actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak untuk melihat kualifikasi pegawai.", nil)
	}

	if days < 1 {
		days = 30 // default 30 hari
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSize
	}

	items, total, err := s.repo.GetExpiringSoonKualifikasi(ctx, days, page, pageSize)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal mengambil kualifikasi pegawai")
	}
	creatorsMap, updatersMap := s.buildAuditMaps(ctx, items)
	return dto.ToKepegawaianKualifikasiListResponse(items, creatorsMap, updatersMap), total, nil
}

// ── GetExpired ────────────────────────────────────────────────────────────────
func (s *service) GetExpiredKualifikasi(
	ctx context.Context,
	page, pageSize int,
	actor he.AuthContext,
) ([]dto.KepegawaianKualifikasiResponse, int64, error) {
	can, err := s.canReadKepegawaianKualifikasi(ctx, actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak untuk melihat kualifikasi pegawai.", nil)
	}

	items, total, err := s.repo.GetExpiredKualifikasi(ctx, page, pageSize)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal mengambil data kualifikasi yang sudah expired")
	}
	creatorsMap, updatersMap := s.buildAuditMaps(ctx, items)
	return dto.ToKepegawaianKualifikasiListResponse(items, creatorsMap, updatersMap), total, nil
}
