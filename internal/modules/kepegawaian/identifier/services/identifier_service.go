package services

import (
	"context"
	"net/http"
	"time"

	"neosim_go/internal/modules/kepegawaian/identifier/dto"
	"neosim_go/internal/modules/kepegawaian/identifier/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

// ── Create ────────────────────────────────────────────────────────────────────

func (s *service) CreateIdentifier(
	ctx context.Context,
	req *dto.CreateKepegawaianIdentifierRequest,
	actor he.AuthContext,
) (*dto.KepegawaianIdentifierResponse, error) {
	can, err := s.canCreateKepegawaianIdentifier(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat KepegawaianIdentifier baru.", nil)
	}

	// validasi keberadaan master Tipe
	tipeMaster, err := s.repo.GetTipeByID(ctx, req.TipeID)
	if err != nil {
		return nil, appErrors.Internal("gagal mengambil master tipe identifier")
	}
	if tipeMaster == nil {
		return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "Tipe identifier tidak ditemukan.", nil)
	}

	// cek duplikasi nilai+tipe (mis. NIK tidak boleh sama antar pegawai)
	duplicate, err := s.repo.ExistsByNilaiAndTipe(ctx, req.TipeID, req.Nilai, 0)
	if err != nil {
		return nil, appErrors.Internal("gagal cek duplikasi identifier")
	}
	if duplicate {
		return nil, appErrors.Wrap(http.StatusConflict,
			"Nilai identifier sudah digunakan oleh pegawai lain.", nil)
	}

	// jika is_primary = true, unset primary lama untuk tipe yang sama
	if req.IsPrimary {
		if err := s.repo.UnsetPrimaryByPegawaiIDAndTipe(ctx, req.PegawaiID, req.TipeID, actor.UserID); err != nil {
			return nil, appErrors.Internal("gagal mereset primary identifier sebelumnya")
		}
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

	if err := s.repo.CreateIdentifier(ctx, m); err != nil {
		return nil, appErrors.Internal("gagal menyimpan identifier pegawai")
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

// ── GetByID ───────────────────────────────────────────────────────────────────

func (s *service) GetIdentifierByID(
	ctx context.Context,
	id int64,
	actor he.AuthContext,
) (*dto.KepegawaianIdentifierResponse, error) {
	can, err := s.canReadKepegawaianIdentifier(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat KepegawaianIdentifier.", nil)
	}

	m, err := s.repo.GetIdentifierByID(ctx, id)
	if err != nil {
		return nil, appErrors.Internal("gagal mengambil data identifier")
	}
	if m == nil {
		return nil, appErrors.Wrap(http.StatusNotFound, "KepegawaianIdentifier tidak ditemukan.", nil)
	}

	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.ToKepegawaianIdentifierResponse(dto.KepegawaianIdentifierResponseParams{
		KepegawaianIdentifier: m,
		Creator:               creator,
		Updater:               updater,
	}), nil
}

// ── List ──────────────────────────────────────────────────────────────────────

func (s *service) ListIdentifier(
	ctx context.Context,
	page, pageSize int,
	filter *dto.FilterKepegawaianIdentifierRequest,
	actor he.AuthContext,
) ([]dto.KepegawaianIdentifierResponse, int64, error) {
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
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSize
	}

	items, total, err := s.repo.ListIdentifier(ctx, page, pageSize, filter)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal mengambil daftar identifier")
	}

	creatorsMap, updatersMap := s.buildAuditMaps(items)
	return dto.ToKepegawaianIdentifierListResponse(items, creatorsMap, updatersMap), total, nil
}

// ── ListByPegawai ─────────────────────────────────────────────────────────────

func (s *service) ListByPegawai(
	ctx context.Context,
	pegawaiID int64,
	page, pageSize int,
	actor he.AuthContext,
) ([]dto.KepegawaianIdentifierResponse, int64, error) {
	can, err := s.canReadKepegawaianIdentifier(actor)
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

	items, total, err := s.repo.FindByPegawaiID(ctx, pegawaiID, page, pageSize)

	if err != nil {
		return nil, 0, appErrors.Internal("gagal mengambil identifier pegawai")
	}

	creatorsMap, updatersMap := s.buildAuditMaps(items)
	return dto.ToKepegawaianIdentifierListResponse(items, creatorsMap, updatersMap), total, nil
}

// ── Update ────────────────────────────────────────────────────────────────────

func (s *service) UpdateIdentifier(
	ctx context.Context,
	id int64,
	req *dto.UpdateKepegawaianIdentifierRequest,
	actor he.AuthContext,
) (*dto.KepegawaianIdentifierResponse, error) {
	can, err := s.canUpdateKepegawaianIdentifier(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah KepegawaianIdentifier.", nil)
	}

	m, err := s.repo.GetIdentifierByID(ctx, id)
	if err != nil {
		return nil, appErrors.Internal("gagal mengambil data identifier")
	}
	if m == nil {
		return nil, appErrors.Wrap(http.StatusNotFound, "KepegawaianIdentifier tidak ditemukan.", nil)
	}

	// tentukan tipe & nilai final untuk cek duplikasi
	tipeCheck := m.TipeID
	nilaiCheck := m.Nilai
	if req.TipeID != nil {
		tipeCheck = *req.TipeID
	}
	if req.Nilai != nil {
		nilaiCheck = *req.Nilai
	}

	duplicate, err := s.repo.ExistsByNilaiAndTipe(ctx, tipeCheck, nilaiCheck, id)
	if err != nil {
		return nil, appErrors.Internal("gagal cek duplikasi identifier")
	}
	if duplicate {
		return nil, appErrors.Wrap(http.StatusConflict,
			"Nilai identifier sudah digunakan oleh pegawai lain.", nil)
	}

	// jika diubah menjadi primary, unset primary lama terlebih dahulu
	if req.IsPrimary != nil && *req.IsPrimary && !m.IsPrimary {
		if err := s.repo.UnsetPrimaryByPegawaiIDAndTipe(ctx, m.PegawaiID, tipeCheck, actor.UserID); err != nil {
			return nil, appErrors.Internal("gagal mereset primary identifier sebelumnya")
		}
	}

	// update parsial jika pointer dikirimkan (not nil)
	if req.TipeID != nil {
		tipeMaster, err := s.repo.GetTipeByID(ctx, *req.TipeID)
		if err != nil {
			return nil, appErrors.Internal("gagal mengambil master tipe identifier")
		}
		if tipeMaster == nil {
			return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "Tipe identifier tidak ditemukan.", nil)
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

	if err := s.repo.UpdateIdentifier(ctx, m); err != nil {
		return nil, appErrors.Internal("gagal menyimpan perubahan identifier")
	}

	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.ToKepegawaianIdentifierResponse(dto.KepegawaianIdentifierResponseParams{
		KepegawaianIdentifier: m,
		Creator:               creator,
		Updater:               updater,
	}), nil
}

// ── Delete ────────────────────────────────────────────────────────────────────

func (s *service) DeleteIdentifier(
	ctx context.Context,
	id int64,
	actor he.AuthContext,
) error {
	can, err := s.canDeleteKepegawaianIdentifier(actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus KepegawaianIdentifier.", nil)
	}

	m, err := s.repo.GetIdentifierByID(ctx, id)
	if err != nil {
		return appErrors.Internal("gagal mengambil data identifier")
	}
	if m == nil {
		return appErrors.Wrap(http.StatusNotFound, "KepegawaianIdentifier tidak ditemukan.", nil)
	}

	// identifier primary tidak boleh dihapus langsung —
	// harus ada identifier lain (aktif) dengan tipe sama sebagai primary dulu
	if m.IsPrimary {
		others, err := s.repo.FindByPegawaiIDAndTipe(ctx, m.PegawaiID, m.TipeID)
		if err != nil {
			return appErrors.Internal("gagal cek identifier lain")
		}
		activeCount := 0
		for _, o := range others {
			if o.ID != id && o.IsAktif {
				activeCount++
			}
		}
		if activeCount > 0 {
			return appErrors.Wrap(http.StatusUnprocessableEntity,
				"Identifier ini adalah primary. Tetapkan identifier lain sebagai primary terlebih dahulu sebelum menghapus.", nil)
		}
	}

	return s.repo.DeleteIdentifier(ctx, id, actor.UserID)
}

// ── GetExpiringSoon ───────────────────────────────────────────────────────────

func (s *service) GetExpiringSoonIdentifier(
	ctx context.Context,
	days int,
	actor he.AuthContext,
) ([]dto.KepegawaianIdentifierResponse, error) {
	can, err := s.canReadKepegawaianIdentifier(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak untuk melihat identifier pegawai.", nil)
	}

	if days < 1 {
		days = 30 // default 30 hari
	}

	items, err := s.repo.FindExpiringSoonIdentifier(ctx, days)
	if err != nil {
		return nil, appErrors.Internal("gagal mengambil data identifier yang akan expired")
	}
	creatorsMap, updatersMap := s.buildAuditMaps(items)
	return dto.ToKepegawaianIdentifierListResponse(items, creatorsMap, updatersMap), nil
}

// ── GetExpired ────────────────────────────────────────────────────────────────

func (s *service) GetExpiredIdentifier(
	ctx context.Context,
	actor he.AuthContext,
) ([]dto.KepegawaianIdentifierResponse, error) {
	can, err := s.canReadKepegawaianIdentifier(actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak untuk melihat identifier pegawai.", nil)
	}

	items, err := s.repo.FindExpiredIdentifier(ctx)
	if err != nil {
		return nil, appErrors.Internal("gagal mengambil data identifier yang sudah expired")
	}
	creatorsMap, updatersMap := s.buildAuditMaps(items)
	return dto.ToKepegawaianIdentifierListResponse(items, creatorsMap, updatersMap), nil
}
