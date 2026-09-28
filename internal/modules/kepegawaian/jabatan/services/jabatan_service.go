package services

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

const jabatanDateLayout = "2006-01-02"

// ── Create ────────────────────────────────────────────────────────────────────
func (s *service) CreateJabatan(ctx context.Context, req *dto.CreateKepegawaianJabatanRequest, actor he.AuthContext) (*dto.KepegawaianJabatanResponse, error) {
	can, err := s.canCreateKepegawaianJabatan(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses: " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat KepegawaianJabatan baru.", nil)
	}

	tanggalMulai, err := parseJabatanDate(req.TanggalMulai)
	if err != nil {
		return nil, unprocessable("tanggal_mulai tidak valid (format YYYY-MM-DD)")
	}
	tanggalSelesai, err := parseJabatanDatePtr(req.TanggalSelesai)
	if err != nil {
		return nil, unprocessable("tanggal_selesai tidak valid (format YYYY-MM-DD)")
	}
	tanggalSK, err := parseJabatanDatePtr(req.TanggalSK)
	if err != nil {
		return nil, unprocessable("tanggal_sk tidak valid (format YYYY-MM-DD)")
	}

	isPrimary := false
	if req.IsPrimary != nil {
		isPrimary = *req.IsPrimary
	}
	// Default: aktif selama belum ada tanggal_selesai
	isAktif := tanggalSelesai == nil
	if req.IsAktif != nil {
		isAktif = *req.IsAktif
	}

	m := &models.KepegawaianJabatan{
		PegawaiID:        req.PegawaiID,
		DepartmentID:     req.DepartmentID,
		PositionID:       req.PositionID,
		JobTitleID:       req.JobTitleID,
		SpecializationID: req.SpecializationID,
		IsPrimary:        isPrimary,
		TanggalMulai:     tanggalMulai,
		TanggalSelesai:   tanggalSelesai,
		NomorSK:          trimPtr(req.NomorSK),
		TanggalSK:        tanggalSK,
		IsAktif:          isAktif,
		CreatedBy:        &actor.UserID,
		UpdatedBy:        &actor.UserID,
	}

	if err := s.validateJabatan(ctx, m, 0); err != nil {
		return nil, err
	}

	if err := s.repo.CreateJabatan(ctx, m); err != nil {
		return nil, err
	}

	// Reload agar relasi Position, JobTitle, Specialization ter-preload
	created, err := s.repo.GetJabatanByID(ctx, m.ID)
	if err != nil {
		return nil, err
	}
	if created == nil {
		created = m
	}

	creator := s.buildCreator(ctx, created.CreatedBy)

	return dto.ToKepegawaianJabatanResponse(dto.KepegawaianJabatanResponseParams{
		KepegawaianJabatan: created,
		Creator:            creator,
		Updater:            creator,
	}), nil
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (s *service) GetJabatanByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.KepegawaianJabatanResponse, error) {
	can, err := s.canReadKepegawaianJabatan(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses: " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat KepegawaianJabatan.", nil)
	}

	m, err := s.repo.GetJabatanByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("KepegawaianJabatan tidak ditemukan")
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToKepegawaianJabatanResponse(dto.KepegawaianJabatanResponseParams{
		KepegawaianJabatan: m,
		Creator:            creator,
		Updater:            updater,
	}), nil
}

// ── GetByPegawaiID ────────────────────────────────────────────────────────────
func (s *service) GetJabatanByPegawaiID(ctx context.Context, pegawaiID int64, page, pageSize int, actor he.AuthContext) ([]dto.KepegawaianJabatanResponse, int64, error) {
	can, err := s.canReadKepegawaianJabatan(ctx, actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat KepegawaianJabatan.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}

	items, total, err := s.repo.GetJabatanByPegawaiID(ctx, pegawaiID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMaps(ctx, items)
	return dto.ToKepegawaianJabatanListResponse(items, creatorsMap, updatersMap), total, nil
}

// ── List ──────────────────────────────────────────────────────────────────────
func (s *service) ListJabatan(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianJabatanRequest, actor he.AuthContext) ([]dto.KepegawaianJabatanResponse, int64, error) {
	can, err := s.canReadKepegawaianJabatan(ctx, actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar KepegawaianJabatan.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}
	items, total, err := s.repo.ListJabatan(ctx, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMaps(ctx, items)
	return dto.ToKepegawaianJabatanListResponse(items, creatorsMap, updatersMap), total, nil
}

// ── Update ────────────────────────────────────────────────────────────────────
func (s *service) UpdateJabatan(ctx context.Context, id int64, req *dto.UpdateKepegawaianJabatanRequest, actor he.AuthContext) (*dto.KepegawaianJabatanResponse, error) {
	can, err := s.canUpdateKepegawaianJabatan(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah KepegawaianJabatan.", nil)
	}

	m, err := s.repo.GetJabatanByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("KepegawaianJabatan tidak ditemukan")
	}

	if req.DepartmentID != nil {
		m.DepartmentID = *req.DepartmentID
	}
	if req.PositionID != nil {
		m.PositionID = *req.PositionID
	}
	if req.JobTitleID != nil {
		m.JobTitleID = *req.JobTitleID
	}
	if req.SpecializationID != nil {
		m.SpecializationID = req.SpecializationID
	}
	if req.IsPrimary != nil {
		m.IsPrimary = *req.IsPrimary
	}
	if req.TanggalMulai != nil {
		t, err := parseJabatanDate(*req.TanggalMulai)
		if err != nil {
			return nil, unprocessable("tanggal_mulai tidak valid (format YYYY-MM-DD)")
		}
		m.TanggalMulai = t
	}
	if req.TanggalSelesai != nil {
		t, err := parseJabatanDate(*req.TanggalSelesai)
		if err != nil {
			return nil, unprocessable("tanggal_selesai tidak valid (format YYYY-MM-DD)")
		}
		m.TanggalSelesai = &t
		// Menutup jabatan: nonaktifkan otomatis kecuali is_aktif dikirim eksplisit
		if req.IsAktif == nil {
			m.IsAktif = false
		}
	}
	if req.NomorSK != nil {
		m.NomorSK = trimPtr(req.NomorSK)
	}
	if req.TanggalSK != nil {
		t, err := parseJabatanDate(*req.TanggalSK)
		if err != nil {
			return nil, unprocessable("tanggal_sk tidak valid (format YYYY-MM-DD)")
		}
		m.TanggalSK = &t
	}
	if req.IsAktif != nil {
		m.IsAktif = *req.IsAktif
	}

	if err := s.validateJabatan(ctx, m, m.ID); err != nil {
		return nil, err
	}

	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdateJabatan(ctx, m); err != nil {
		return nil, err
	}

	// Reload agar relasi sesuai ID terbaru (bukan data lama hasil preload)
	updated, err := s.repo.GetJabatanByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		updated = m
	}

	creator := s.buildCreator(ctx, updated.CreatedBy)
	updater := s.buildCreator(ctx, updated.UpdatedBy)

	return dto.ToKepegawaianJabatanResponse(dto.KepegawaianJabatanResponseParams{
		KepegawaianJabatan: updated,
		Creator:            creator,
		Updater:            updater,
	}), nil
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (s *service) DeleteJabatan(ctx context.Context, id int64, actor he.AuthContext) error {
	can, err := s.canDeleteKepegawaianJabatan(ctx, actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus KepegawaianJabatan.", nil)
	}

	m, err := s.repo.GetJabatanByID(ctx, id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("KepegawaianJabatan tidak ditemukan")
	}
	return s.repo.DeleteJabatan(ctx, id, actor.UserID)
}

// ── helper khusus KepegawaianJabatan ─────────────────────────────────────────

// validateJabatan memeriksa aturan lintas field pada state akhir model.
// excludeID = 0 saat create, id yang sedang diedit saat update.
func (s *service) validateJabatan(ctx context.Context, m *models.KepegawaianJabatan, excludeID int64) error {
	// 1. tanggal_selesai tidak boleh sebelum tanggal_mulai
	if m.TanggalSelesai != nil && m.TanggalSelesai.Before(m.TanggalMulai) {
		return unprocessable("tanggal_selesai tidak boleh lebih awal dari tanggal_mulai")
	}

	// 2. specialization (jika ada) harus milik job_title yang sama
	if m.SpecializationID != nil {
		spec, err := s.repo.GetSpecializationByID(ctx, *m.SpecializationID)
		if err != nil {
			return err
		}
		if spec == nil {
			return unprocessable("Specialization tidak ditemukan")
		}
		if spec.JobTitleID != nil && *spec.JobTitleID != m.JobTitleID {
			return unprocessable("Specialization tidak sesuai dengan job title yang dipilih")
		}
	}

	// 3. hanya satu jabatan primer aktif per pegawai
	if m.IsPrimary && m.IsAktif {
		exists, err := s.repo.ExistsPrimaryAktifByPegawai(ctx, m.PegawaiID, excludeID)
		if err != nil {
			return err
		}
		if exists {
			return appErrors.Wrap(http.StatusConflict,
				"Pegawai sudah memiliki jabatan primer yang aktif.", nil)
		}
	}

	return nil
}

func unprocessable(msg string) error {
	return appErrors.Wrap(http.StatusUnprocessableEntity, msg, nil)
}

func parseJabatanDate(s string) (time.Time, error) {
	return time.Parse(jabatanDateLayout, strings.TrimSpace(s))
}

func parseJabatanDatePtr(s *string) (*time.Time, error) {
	if s == nil || strings.TrimSpace(*s) == "" {
		return nil, nil
	}
	t, err := parseJabatanDate(*s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func trimPtr(s *string) *string {
	if s == nil {
		return nil
	}
	v := strings.TrimSpace(*s)
	if v == "" {
		return nil
	}
	return &v
}
