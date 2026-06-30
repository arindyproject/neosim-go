package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"neosim_go/internal/modules/kepegawaian/identifier/dto"
	"neosim_go/internal/modules/kepegawaian/identifier/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

// ── Create ────────────────────────────────────────────────────────────────────

func (s *service) Create(
	ctx context.Context,
	req *dto.CreateKepegawaianIdentifierRequest,
	actor he.AuthContext,
) (*dto.KepegawaianIdentifierResponse, error) {
	can, err := s.canCreateKepegawaianIdentifier(actor, req.PegawaiID)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak untuk membuat identifier pegawai.", nil)
	}

	// validasi bisnis: tipe valid + tanggal_expired jika diperlukan
	if err := req.Validate(); err != nil {
		return nil, appErrors.Wrap(http.StatusUnprocessableEntity, err.Error(), nil)
	}

	// cek duplikasi nilai+tipe lintas pegawai (misal NIK tidak boleh sama)
	duplicate, err := s.repo.ExistsByNilaiAndTipe(ctx, req.Tipe, req.Nilai, 0)
	if err != nil {
		return nil, appErrors.Internal("gagal cek duplikasi identifier")
	}
	if duplicate {
		return nil, appErrors.Wrap(http.StatusConflict,
			"Nilai identifier sudah digunakan oleh pegawai lain.", nil)
	}

	// jika is_primary = true, unset primary lama untuk tipe yang sama
	if req.IsPrimary {
		if err := s.repo.UnsetPrimaryByPegawaiDAndTipe(ctx, req.PegawaiID, req.Tipe, actor.UserID); err != nil {
			return nil, appErrors.Internal("gagal mereset primary identifier sebelumnya")
		}
	}

	// parse tanggal dari string
	tanggalTerbit, tanggalExpired, err := parseTanggalIdentifier(req.TanggalTerbit, req.TanggalExpired)
	if err != nil {
		return nil, appErrors.Wrap(http.StatusUnprocessableEntity, err.Error(), nil)
	}

	m := &models.KepegawaianIdentifier{
		PegawaiID:      req.PegawaiID,
		Tipe:           req.Tipe,
		Nilai:          req.Nilai,
		Penerbit:       req.Penerbit,
		TanggalTerbit:  tanggalTerbit,
		TanggalExpired: tanggalExpired,
		IsPrimary:      req.IsPrimary,
		IsAktif:        req.IsAktif,
		CreatedBy:      &actor.UserID,
		UpdatedBy:      &actor.UserID,
	}

	if err := s.repo.Create(ctx, m); err != nil {
		return nil, appErrors.Internal("gagal menyimpan identifier pegawai")
	}

	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)

	return dto.ToKepegawaianIdentifierResponse(dto.KepegawaianIdentifierResponseParams{
		Identifier: m,
		Creator:    creator,
		Updater:    updater,
	}), nil
}

// ── GetByID ───────────────────────────────────────────────────────────────────

func (s *service) GetByID(
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
			"Akses ditolak. Anda tidak memiliki hak untuk melihat identifier pegawai.", nil)
	}

	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, appErrors.Internal("gagal mengambil data identifier")
	}
	if m == nil {
		return nil, appErrors.Wrap(http.StatusNotFound, "Identifier pegawai tidak ditemukan.", nil)
	}

	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)
	return dto.ToKepegawaianIdentifierResponse(dto.KepegawaianIdentifierResponseParams{
		Identifier: m,
		Creator:    creator,
		Updater:    updater,
	}), nil
}

// ── List ──────────────────────────────────────────────────────────────────────

func (s *service) List(
	ctx context.Context,
	page, pageSize int,
	filter dto.FilterKepegawaianIdentifierRequest,
	actor he.AuthContext,
) ([]dto.KepegawaianIdentifierResponse, int64, error) {
	can, err := s.canReadKepegawaianIdentifier(actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak untuk melihat daftar identifier pegawai.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSize
	}

	items, total, err := s.repo.FindAll(ctx, filter, page, pageSize)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal mengambil daftar identifier")
	}

	creatorsMap, updatersMap := s.buildAuditMaps(items)
	return toIdentifierResponses(items, creatorsMap, updatersMap), total, nil

}

// ── ListByPegawai ─────────────────────────────────────────────────────────────

func (s *service) ListByPegawai(
	ctx context.Context,
	kepegawaianID int64,
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

	items, err := s.repo.FindByPegawaiD(ctx, kepegawaianID)
	if err != nil {
		return nil, appErrors.Internal("gagal mengambil identifier pegawai")
	}

	creatorsMap, updatersMap := s.buildAuditMaps(items)
	return toIdentifierResponses(items, creatorsMap, updatersMap), nil
}

// ── Update ────────────────────────────────────────────────────────────────────

func (s *service) Update(
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
			"Akses ditolak. Anda tidak memiliki hak untuk mengubah identifier pegawai.", nil)
	}

	// validasi bisnis request
	if err := req.Validate(); err != nil {
		return nil, appErrors.Wrap(http.StatusUnprocessableEntity, err.Error(), nil)
	}

	// ambil data existing
	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, appErrors.Internal("gagal mengambil data identifier")
	}
	if m == nil {
		return nil, appErrors.Wrap(http.StatusNotFound, "Identifier pegawai tidak ditemukan.", nil)
	}

	// cek duplikasi jika nilai atau tipe diubah
	tipeCheck := m.Tipe
	nilaiCheck := m.Nilai
	if req.Tipe != nil {
		tipeCheck = *req.Tipe
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
		if err := s.repo.UnsetPrimaryByPegawaiDAndTipe(ctx, m.PegawaiID, tipeCheck, actor.UserID); err != nil {
			return nil, appErrors.Internal("gagal mereset primary identifier sebelumnya")
		}
	}

	// terapkan perubahan ke model
	if req.Tipe != nil {
		m.Tipe = *req.Tipe
	}
	if req.Nilai != nil {
		m.Nilai = *req.Nilai
	}
	if req.Penerbit != nil {
		m.Penerbit = req.Penerbit
	}
	if req.IsPrimary != nil {
		m.IsPrimary = *req.IsPrimary
	}
	if req.IsAktif != nil {
		m.IsAktif = *req.IsAktif
	}

	// parse ulang tanggal jika ada perubahan
	if req.TanggalTerbit != nil || req.TanggalExpired != nil {
		tanggalTerbit, tanggalExpired, err := parseTanggalIdentifier(req.TanggalTerbit, req.TanggalExpired)
		if err != nil {
			return nil, appErrors.Wrap(http.StatusUnprocessableEntity, err.Error(), nil)
		}
		if req.TanggalTerbit != nil {
			m.TanggalTerbit = tanggalTerbit
		}
		if req.TanggalExpired != nil {
			m.TanggalExpired = tanggalExpired
		}
	}

	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, m); err != nil {
		return nil, appErrors.Internal("gagal menyimpan perubahan identifier")
	}

	creator := s.buildCreator(m.CreatedBy)
	updater := s.buildCreator(m.UpdatedBy)
	return dto.ToKepegawaianIdentifierResponse(dto.KepegawaianIdentifierResponseParams{
		Identifier: m,
		Creator:    creator,
		Updater:    updater,
	}), nil
}

// ── Delete ────────────────────────────────────────────────────────────────────

func (s *service) Delete(
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
			"Akses ditolak. Anda tidak memiliki hak untuk menghapus identifier pegawai.", nil)
	}

	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return appErrors.Internal("gagal mengambil data identifier")
	}
	if m == nil {
		return appErrors.Wrap(http.StatusNotFound, "Identifier pegawai tidak ditemukan.", nil)
	}

	// identifier primary tidak boleh dihapus langsung
	// harus set identifier lain sebagai primary dulu
	if m.IsPrimary {
		others, err := s.repo.FindByPegawaiDAndTipe(ctx, m.PegawaiID, m.Tipe)
		if err != nil {
			return appErrors.Internal("gagal cek identifier lain")
		}
		// hitung identifier aktif selain yang akan dihapus
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

	return s.repo.Delete(ctx, id, actor.UserID)
}

// ── GetExpiringSoon ───────────────────────────────────────────────────────────

func (s *service) GetExpiringSoon(
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

	items, err := s.repo.FindExpiringSoon(ctx, days)
	if err != nil {
		return nil, appErrors.Internal("gagal mengambil data identifier yang akan expired")
	}
	creatorsMap, updatersMap := s.buildAuditMaps(items)
	return toIdentifierResponses(items, creatorsMap, updatersMap), nil
}

// ── GetIdentifierTypes ────────────────────────────────────────────────────────

// GetIdentifierTypes mengembalikan semua tipe identifier untuk dropdown UI
// tidak perlu cek RBAC — data ini bersifat publik dalam sistem
func (s *service) GetIdentifierTypes() []dto.IdentifierMetaResponse {
	return dto.ToIdentifierMetaListResponse(models.AllIdentifierTypes())
}

// ── Helper ────────────────────────────────────────────────────────────────────

const tanggalLayout = "2006-01-02"

// parseTanggalIdentifier mem-parse string tanggal menjadi *time.Time
func parseTanggalIdentifier(terbit, expired *string) (tanggalTerbit, tanggalExpired *time.Time, err error) {
	if terbit != nil && *terbit != "" {
		t, e := time.Parse(tanggalLayout, *terbit)
		if e != nil {
			return nil, nil, errors.New("format tanggal_terbit tidak valid, gunakan YYYY-MM-DD")
		}
		tanggalTerbit = &t
	}
	if expired != nil && *expired != "" {
		t, e := time.Parse(tanggalLayout, *expired)
		if e != nil {
			return nil, nil, errors.New("format tanggal_expired tidak valid, gunakan YYYY-MM-DD")
		}
		tanggalExpired = &t
	}
	return tanggalTerbit, tanggalExpired, nil
}
