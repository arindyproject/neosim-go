package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"neosim_go/internal/modules/kepegawaian/pendidikan/dto"
	"neosim_go/internal/modules/kepegawaian/pendidikan/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

// ── Create ────────────────────────────────────────────────────────────────────
func (s *service) CreatePendidikan(ctx context.Context, req *dto.CreateKepegawaianPendidikanRequest, actor he.AuthContext) (*dto.KepegawaianPendidikanResponse, error) {
	can, err := s.canCreateKepegawaianPendidikan(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat KepegawaianPendidikan baru.", nil)
	}

	// validasi keberadaan master Tipe
	jenjangMaster, err := s.repo.GetJenjangByID(ctx, req.JenjangID)
	if err != nil {
		return nil, appErrors.Internal("gagal mengambil master jenjang pendidikan")
	}
	if jenjangMaster == nil {
		return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "Jenjang Pendidikan tidak ditemukan.", nil)
	}

	m := &models.KepegawaianPendidikan{
		PegawaiID:       req.PegawaiID,
		JenjangID:       req.JenjangID,
		NamaInstitusi:   req.NamaInstitusi,
		NomorIjazah:     req.NomorIjazah,
		BidangStudi:     req.BidangStudi,
		AlamatInstitusi: req.AlamatInstitusi,
		NilaiAkhir:      req.NilaiAkhir,
		TanggalMasuk:    req.TanggalMasuk.ToTimePtr(),
		TanggalLulus:    req.TanggalLulus.ToTimePtr(),
		FHIRCode:        req.FHIRCode,
		FHIRSystem:      req.FHIRSystem,
		CreatedBy:       &actor.UserID,
		UpdatedBy:       &actor.UserID,
	}
	if err := s.repo.CreatePendidikan(ctx, m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, m.CreatedBy)

	return dto.ToKepegawaianPendidikanResponse(dto.KepegawaianPendidikanResponseParams{
		KepegawaianPendidikan: m,
		Creator:               creator,
		Updater:               creator,
	}), nil
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (s *service) GetPendidikanByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.KepegawaianPendidikanResponse, error) {
	can, err := s.canReadKepegawaianPendidikan(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk Melihat KepegawaianPendidikan.", nil)
	}

	m, err := s.repo.GetPendidikanByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("KepegawaianPendidikan tidak ditemukan")
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToKepegawaianPendidikanResponse(dto.KepegawaianPendidikanResponseParams{
		KepegawaianPendidikan: m,
		Creator:               creator,
		Updater:               updater,
	}), nil
}

// ── List ──────────────────────────────────────────────────────────────────────
func (s *service) ListPendidikan(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianPendidikanRequest, actor he.AuthContext) ([]dto.KepegawaianPendidikanResponse, int64, error) {
	can, err := s.canReadKepegawaianPendidikan(ctx, actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar KepegawaianPendidikan.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}
	items, total, err := s.repo.ListPendidikan(ctx, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMaps(ctx, items)
	return dto.ToKepegawaianPendidikanListResponse(items, creatorsMap, updatersMap), total, nil
}

// ── Update ────────────────────────────────────────────────────────────────────
func (s *service) UpdatePendidikan(ctx context.Context, id int64, req *dto.UpdateKepegawaianPendidikanRequest, actor he.AuthContext) (*dto.KepegawaianPendidikanResponse, error) {
	can, err := s.canUpdateKepegawaianPendidikan(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah KepegawaianPendidikan.", nil)
	}

	m, err := s.repo.GetPendidikanByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("KepegawaianPendidikan tidak ditemukan")
	}

	if req.JenjangID != nil {
		m.JenjangID = *req.JenjangID
	}
	if req.NamaInstitusi != nil {
		m.NamaInstitusi = *req.NamaInstitusi
	}
	m.NomorIjazah = req.NomorIjazah
	m.BidangStudi = req.BidangStudi
	m.AlamatInstitusi = req.AlamatInstitusi
	m.NilaiAkhir = req.NilaiAkhir
	m.TanggalMasuk = req.TanggalMasuk.ToTimePtr()
	m.TanggalLulus = req.TanggalLulus.ToTimePtr()
	m.FHIRCode = req.FHIRCode
	m.FHIRSystem = req.FHIRSystem
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdatePendidikan(ctx, m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToKepegawaianPendidikanResponse(dto.KepegawaianPendidikanResponseParams{
		KepegawaianPendidikan: m,
		Creator:               creator,
		Updater:               updater,
	}), nil
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (s *service) DeletePendidikan(ctx context.Context, id int64, actor he.AuthContext) error {
	can, err := s.canDeleteKepegawaianPendidikan(ctx, actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus KepegawaianPendidikan.", nil)
	}

	m, err := s.repo.GetPendidikanByID(ctx, id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("KepegawaianPendidikan tidak ditemukan")
	}
	return s.repo.DeletePendidikan(ctx, id)
}
