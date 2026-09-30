package services

import (
	"context"
	"net/http"
	"time"

	"neosim_go/internal/modules/kepegawaian/pegawai/dto"
	"neosim_go/internal/modules/kepegawaian/pegawai/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

const layoutTanggal = "2006-01-02"

func parseTanggal(label, value string) (time.Time, error) {
	t, err := time.Parse(layoutTanggal, value)
	if err != nil {
		return time.Time{}, appErrors.Wrap(http.StatusBadRequest,
			label+" tidak valid, gunakan format YYYY-MM-DD", nil)
	}
	return t, nil
}

// cekJenisStatus memastikan jenis_id & status_id ada di tabel master
func (s *service) cekJenisStatus(ctx context.Context, jenisID, statusID int64) error {
	jenis, err := s.repo.GetJenisByID(ctx, jenisID)
	if err != nil {
		return err
	}
	if jenis == nil {
		return appErrors.Wrap(http.StatusBadRequest, "Jenis pegawai tidak ditemukan", nil)
	}

	status, err := s.repo.GetStatusByID(ctx, statusID)
	if err != nil {
		return err
	}
	if status == nil {
		return appErrors.Wrap(http.StatusBadRequest, "Status pegawai tidak ditemukan", nil)
	}
	return nil
}

// ── Create ────────────────────────────────────────────────────────────────────
func (s *service) CreatePegawai(ctx context.Context, req *dto.CreateKepegawaianPegawaiRequest, actor he.AuthContext) (*dto.KepegawaianPegawaiResponse, error) {
	can, err := s.canCreateKepegawaianPegawai(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat Pegawai baru.", nil)
	}

	// parse tanggal
	tglLahir, err := parseTanggal("Tanggal lahir", req.TanggalLahir)
	if err != nil {
		return nil, err
	}
	tglMasuk, err := parseTanggal("Tanggal masuk", req.TanggalMasuk)
	if err != nil {
		return nil, err
	}
	var tglKeluar *time.Time
	if req.TanggalKeluar != nil {
		t, err := parseTanggal("Tanggal keluar", *req.TanggalKeluar)
		if err != nil {
			return nil, err
		}
		if t.Before(tglMasuk) {
			return nil, appErrors.Wrap(http.StatusBadRequest,
				"Tanggal keluar tidak boleh lebih awal dari tanggal masuk", nil)
		}
		tglKeluar = &t
	}

	// is_aktif: default true, dipaksa false jika sudah ada tanggal keluar
	isAktif := true
	if req.IsAktif != nil {
		isAktif = *req.IsAktif
	}
	if tglKeluar != nil {
		isAktif = false
	}

	// cek duplikat NIK
	dup, err := s.repo.GetPegawaiByNIK(ctx, req.NIK)
	if err != nil {
		return nil, err
	}
	if dup != nil {
		return nil, appErrors.Wrap(http.StatusConflict, "Pegawai dengan NIK ini sudah ada", nil)
	}

	// cek duplikat nomor pegawai
	dup, err = s.repo.GetPegawaiByNomorPegawai(ctx, req.NomorPegawai)
	if err != nil {
		return nil, err
	}
	if dup != nil {
		return nil, appErrors.Wrap(http.StatusConflict, "Pegawai dengan nomor pegawai ini sudah ada", nil)
	}

	// cek user sudah tertaut ke pegawai lain
	if req.UserID != nil {
		dup, err = s.repo.GetPegawaiByUserID(ctx, *req.UserID)
		if err != nil {
			return nil, err
		}
		if dup != nil {
			return nil, appErrors.Wrap(http.StatusConflict, "User ini sudah tertaut ke pegawai lain", nil)
		}
	}

	// cek master jenis & status
	if err := s.cekJenisStatus(ctx, req.JenisID, req.StatusID); err != nil {
		return nil, err
	}

	m := &models.KepegawaianPegawai{
		UserID:           req.UserID,
		NIK:              req.NIK,
		IHSNumber:        req.IHSNumber,
		NomorPegawai:     req.NomorPegawai,
		NamaLengkap:      req.NamaLengkap,
		JenisKelamin:     req.JenisKelamin,
		TanggalLahir:     tglLahir,
		TempatLahir:      req.TempatLahir,
		GolonganDarah:    req.GolonganDarah,
		Agama:            req.Agama,
		StatusPerkawinan: req.StatusPerkawinan,
		Kewarganegaraan:  req.Kewarganegaraan,
		TanggalMasuk:     tglMasuk,
		TanggalKeluar:    tglKeluar,
		JenisID:          req.JenisID,
		StatusID:         req.StatusID,
		FotoURL:          req.FotoURL,
		IsAktif:          isAktif,
		CreatedBy:        &actor.UserID,
		UpdatedBy:        &actor.UserID,
	}
	if err := s.repo.CreatePegawai(ctx, m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, m.CreatedBy)

	return dto.ToKepegawaianPegawaiResponse(dto.KepegawaianPegawaiResponseParams{
		KepegawaianPegawai: m,
		Creator:            creator,
		Updater:            creator,
	}), nil
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (s *service) GetPegawaiByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.KepegawaianPegawaiResponse, error) {
	can, err := s.canReadKepegawaianPegawai(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat Pegawai.", nil)
	}

	m, err := s.repo.GetPegawaiByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.Wrap(http.StatusNotFound, "Pegawai tidak ditemukan", nil)
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToKepegawaianPegawaiResponse(dto.KepegawaianPegawaiResponseParams{
		KepegawaianPegawai: m,
		Creator:            creator,
		Updater:            updater,
	}), nil
}

// ── List ──────────────────────────────────────────────────────────────────────
func (s *service) ListPegawai(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianPegawaiRequest, actor he.AuthContext) ([]dto.KepegawaianPegawaiResponse, int64, error) {
	can, err := s.canReadKepegawaianPegawai(ctx, actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar Pegawai.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}
	items, total, err := s.repo.ListPegawai(ctx, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMaps(ctx, items)
	return dto.ToKepegawaianPegawaiListResponse(items, creatorsMap, updatersMap), total, nil
}

// ── Update ────────────────────────────────────────────────────────────────────
func (s *service) UpdatePegawai(ctx context.Context, id int64, req *dto.UpdateKepegawaianPegawaiRequest, actor he.AuthContext) (*dto.KepegawaianPegawaiResponse, error) {
	can, err := s.canUpdateKepegawaianPegawai(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah Pegawai.", nil)
	}

	m, err := s.repo.GetPegawaiByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.Wrap(http.StatusNotFound, "Pegawai tidak ditemukan", nil)
	}

	// cek duplikat hanya jika nilainya berubah
	if req.NIK != nil && *req.NIK != m.NIK {
		dup, err := s.repo.GetPegawaiByNIK(ctx, *req.NIK)
		if err != nil {
			return nil, err
		}
		if dup != nil {
			return nil, appErrors.Wrap(http.StatusConflict, "Pegawai dengan NIK ini sudah ada", nil)
		}
		m.NIK = *req.NIK
	}
	if req.NomorPegawai != nil && *req.NomorPegawai != m.NomorPegawai {
		dup, err := s.repo.GetPegawaiByNomorPegawai(ctx, *req.NomorPegawai)
		if err != nil {
			return nil, err
		}
		if dup != nil {
			return nil, appErrors.Wrap(http.StatusConflict, "Pegawai dengan nomor pegawai ini sudah ada", nil)
		}
		m.NomorPegawai = *req.NomorPegawai
	}
	if req.UserID != nil && (m.UserID == nil || *req.UserID != *m.UserID) {
		dup, err := s.repo.GetPegawaiByUserID(ctx, *req.UserID)
		if err != nil {
			return nil, err
		}
		if dup != nil && dup.ID != m.ID {
			return nil, appErrors.Wrap(http.StatusConflict, "User ini sudah tertaut ke pegawai lain", nil)
		}
		m.UserID = req.UserID
	}

	// master jenis & status
	if req.JenisID != nil || req.StatusID != nil {
		jenisID, statusID := m.JenisID, m.StatusID
		if req.JenisID != nil {
			jenisID = *req.JenisID
		}
		if req.StatusID != nil {
			statusID = *req.StatusID
		}
		if err := s.cekJenisStatus(ctx, jenisID, statusID); err != nil {
			return nil, err
		}
		m.JenisID = jenisID
		m.StatusID = statusID
	}

	// tanggal
	if req.TanggalLahir != nil {
		t, err := parseTanggal("Tanggal lahir", *req.TanggalLahir)
		if err != nil {
			return nil, err
		}
		m.TanggalLahir = t
	}
	if req.TanggalMasuk != nil {
		t, err := parseTanggal("Tanggal masuk", *req.TanggalMasuk)
		if err != nil {
			return nil, err
		}
		m.TanggalMasuk = t
	}
	if req.TanggalKeluar != nil {
		t, err := parseTanggal("Tanggal keluar", *req.TanggalKeluar)
		if err != nil {
			return nil, err
		}
		m.TanggalKeluar = &t
		if req.IsAktif == nil {
			m.IsAktif = false
		}
	}

	// field biasa
	if req.IHSNumber != nil {
		m.IHSNumber = req.IHSNumber
	}
	if req.NamaLengkap != nil {
		m.NamaLengkap = *req.NamaLengkap
	}
	if req.JenisKelamin != nil {
		m.JenisKelamin = *req.JenisKelamin
	}
	if req.TempatLahir != nil {
		m.TempatLahir = *req.TempatLahir
	}
	if req.GolonganDarah != nil {
		m.GolonganDarah = req.GolonganDarah
	}
	if req.Agama != nil {
		m.Agama = *req.Agama
	}
	if req.StatusPerkawinan != nil {
		m.StatusPerkawinan = *req.StatusPerkawinan
	}
	if req.Kewarganegaraan != nil {
		m.Kewarganegaraan = *req.Kewarganegaraan
	}
	if req.FotoURL != nil {
		m.FotoURL = req.FotoURL
	}
	if req.IsAktif != nil {
		m.IsAktif = *req.IsAktif
	}

	// validasi akhir setelah semua perubahan diterapkan
	if m.TanggalKeluar != nil {
		if m.TanggalKeluar.Before(m.TanggalMasuk) {
			return nil, appErrors.Wrap(http.StatusBadRequest,
				"Tanggal keluar tidak boleh lebih awal dari tanggal masuk", nil)
		}
		if m.IsAktif {
			return nil, appErrors.Wrap(http.StatusBadRequest,
				"Pegawai yang sudah memiliki tanggal keluar tidak boleh berstatus aktif", nil)
		}
	}

	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	if err := s.repo.UpdatePegawai(ctx, m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToKepegawaianPegawaiResponse(dto.KepegawaianPegawaiResponseParams{
		KepegawaianPegawai: m,
		Creator:            creator,
		Updater:            updater,
	}), nil
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (s *service) DeletePegawai(ctx context.Context, id int64, actor he.AuthContext) error {
	can, err := s.canDeleteKepegawaianPegawai(ctx, actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus Pegawai.", nil)
	}

	m, err := s.repo.GetPegawaiByID(ctx, id)
	if err != nil {
		return err
	}
	if m == nil {
		return appErrors.Wrap(http.StatusNotFound, "Pegawai tidak ditemukan", nil)
	}
	return s.repo.DeletePegawai(ctx, id, actor.UserID)
}
