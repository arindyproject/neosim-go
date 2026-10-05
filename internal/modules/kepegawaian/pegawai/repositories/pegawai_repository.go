package repositories

import (
	"context"
	"errors"
	"time"

	"neosim_go/internal/modules/kepegawaian/pegawai/dto"
	"neosim_go/internal/modules/kepegawaian/pegawai/models"
	"neosim_go/internal/shared/sorting"

	"gorm.io/gorm"
)

// relasi yang di-preload untuk response lengkap
var pegawaiPreloads = []string{
	"Jenis",
	"Status",
	"JenisKelamin",
	"GolonganDarah",
	"Agama",
	"StatusPernikahan",
}

// withPegawaiPreloads menambahkan semua Preload relasi ke query
func withPegawaiPreloads(db *gorm.DB) *gorm.DB {
	for _, rel := range pegawaiPreloads {
		db = db.Preload(rel)
	}
	return db
}

// allowedSortColumns memetakan nilai sort_by dari client ke ekspresi kolom DB.
// Hanya key di map ini yang boleh dipakai (mencegah SQL injection lewat Order()).
var KepegawaianPegawaiSort = sorting.Config{
	Allowed: map[string]string{
		"user_id":              "user_id",
		"nik":                  "nik",
		"ihs_number":           "ihs_number",
		"nomor_pegawai":        "nomor_pegawai",
		"nama_lengkap":         "nama_lengkap",
		"tanggal_lahir":        "tanggal_lahir",
		"tempat_lahir":         "tempat_lahir",
		"golongan_darah_id":    "golongan_darah_id",
		"agama_id":             "agama_id",
		"status_pernikahan_id": "status_pernikahan_id",
		"kewarganegaraan":      "kewarganegaraan",
		"tanggal_masuk":        "tanggal_masuk",
		"tanggal_keluar":       "tanggal_keluar",
		"jenis_id":             "jenis_id",
		"status_id":            "status_id",
		"foto_url":             "foto_url",
		"is_aktif":             "is_aktif",
		"created_at":           "created_at",
		"updated_at":           "updated_at",
	},
	DefaultColumn: "created_at",
	DefaultDesc:   true,
	TieBreaker:    "id",
}

// ── Create ────────────────────────────────────────────────────────────────────
func (r *repository) CreatePegawai(ctx context.Context, m *models.KepegawaianPegawai) error {
	if err := r.db.WithContext(ctx).Omit(pegawaiPreloads...).Create(m).Error; err != nil {
		return err
	}
	// muat relasi agar response langsung lengkap
	return withPegawaiPreloads(r.db.WithContext(ctx)).First(m, m.ID).Error
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (r *repository) GetPegawaiByID(ctx context.Context, id int64) (*models.KepegawaianPegawai, error) {
	var m models.KepegawaianPegawai
	err := withPegawaiPreloads(r.db.WithContext(ctx)).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// ── GetByNIK ──────────────────────────────────────────────────────────────────
func (r *repository) GetPegawaiByNIK(ctx context.Context, nik string) (*models.KepegawaianPegawai, error) {
	var m models.KepegawaianPegawai
	err := r.db.WithContext(ctx).
		Where("nik = ? AND deleted_at IS NULL", nik).
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// ── GetByNomorPegawai ─────────────────────────────────────────────────────────
func (r *repository) GetPegawaiByNomorPegawai(ctx context.Context, nomor string) (*models.KepegawaianPegawai, error) {
	var m models.KepegawaianPegawai
	err := r.db.WithContext(ctx).
		Where("nomor_pegawai = ? AND deleted_at IS NULL", nomor).
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// ── GetByUserID ───────────────────────────────────────────────────────────────
func (r *repository) GetPegawaiByUserID(ctx context.Context, userID int64) (*models.KepegawaianPegawai, error) {
	var m models.KepegawaianPegawai
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// ── GetByIHSNumber ────────────────────────────────────────────────────────────
func (r *repository) GetPegawaiByIHSNumber(ctx context.Context, ihs string) (*models.KepegawaianPegawai, error) {
	var m models.KepegawaianPegawai
	err := r.db.WithContext(ctx).
		Where("ihs_number = ? AND deleted_at IS NULL", ihs).
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// ── List ──────────────────────────────────────────────────────────────────────
func (r *repository) ListPegawai(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianPegawaiRequest) ([]models.KepegawaianPegawai, int64, error) {
	var items []models.KepegawaianPegawai
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.KepegawaianPegawai{}).
		Where("deleted_at IS NULL")

	if filter != nil {
		if filter.Name != "" {
			query = query.Where("nama_lengkap ILIKE ?", "%"+filter.Name+"%")
		}
		if filter.NIK != "" {
			query = query.Where("nik = ?", filter.NIK)
		}
		if filter.NomorPegawai != "" {
			query = query.Where("nomor_pegawai ILIKE ?", "%"+filter.NomorPegawai+"%")
		}
		if filter.JenisKelaminID != nil {
			query = query.Where("jenis_kelamin_id = ?", *filter.JenisKelaminID)
		}
		if filter.AgamaID != nil {
			query = query.Where("agama_id = ?", *filter.AgamaID)
		}
		if filter.StatusPernikahanID != nil {
			query = query.Where("status_pernikahan_id = ?", *filter.StatusPernikahanID)
		}
		if filter.JenisID != nil {
			query = query.Where("jenis_id = ?", *filter.JenisID)
		}
		if filter.StatusID != nil {
			query = query.Where("status_id = ?", *filter.StatusID)
		}
		if filter.IsAktif != nil {
			query = query.Where("is_aktif = ?", *filter.IsAktif)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := withPegawaiPreloads(query).
		Scopes(KepegawaianPegawaiSort.Scope(filter.SortBy, filter.SortOrder)).
		Offset(offset).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// ── Update ────────────────────────────────────────────────────────────────────
func (r *repository) UpdatePegawai(ctx context.Context, m *models.KepegawaianPegawai) error {
	// Omit semua relasi agar Save tidak menulis ulang data master
	if err := r.db.WithContext(ctx).
		Omit(pegawaiPreloads...).
		Save(m).Error; err != nil {
		return err
	}
	// muat ulang relasi karena ID master mungkin berubah
	return withPegawaiPreloads(r.db.WithContext(ctx)).First(m, m.ID).Error
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (r *repository) DeletePegawai(ctx context.Context, id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.KepegawaianPegawai{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
		}).Error
}
