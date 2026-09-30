package repositories

import (
	"context"
	"errors"
	"time"

	"neosim_go/internal/modules/kepegawaian/pegawai/dto"
	"neosim_go/internal/modules/kepegawaian/pegawai/models"

	"gorm.io/gorm"
)

// ── Create ────────────────────────────────────────────────────────────────────
func (r *repository) CreatePegawai(ctx context.Context, m *models.KepegawaianPegawai) error {
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	// muat relasi Jenis & Status agar response langsung lengkap
	return r.db.WithContext(ctx).
		Preload("Jenis").
		Preload("Status").
		First(m, m.ID).Error
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (r *repository) GetPegawaiByID(ctx context.Context, id int64) (*models.KepegawaianPegawai, error) {
	var m models.KepegawaianPegawai
	err := r.db.WithContext(ctx).
		Preload("Jenis").
		Preload("Status").
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
		if filter.JenisKelamin != "" {
			query = query.Where("jenis_kelamin = ?", filter.JenisKelamin)
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
	if err := query.
		Preload("Jenis").
		Preload("Status").
		Order("nama_lengkap ASC, id ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// ── Update ────────────────────────────────────────────────────────────────────
func (r *repository) UpdatePegawai(ctx context.Context, m *models.KepegawaianPegawai) error {
	// Omit relasi agar Save tidak ikut menulis ulang Jenis/Status
	if err := r.db.WithContext(ctx).
		Omit("Jenis", "Status").
		Save(m).Error; err != nil {
		return err
	}
	// muat ulang relasi karena JenisID/StatusID mungkin berubah
	return r.db.WithContext(ctx).
		Preload("Jenis").
		Preload("Status").
		First(m, m.ID).Error
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
