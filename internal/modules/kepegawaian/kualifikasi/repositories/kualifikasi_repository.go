package repositories

import (
	"context"
	"errors"
	"time"

	"neosim_go/internal/modules/kepegawaian/kualifikasi/dto"
	"neosim_go/internal/modules/kepegawaian/kualifikasi/models"

	"gorm.io/gorm"
)

// ── Create ────────────────────────────────────────────────────────────────────
func (r *repository) CreateKualifikasi(ctx context.Context, m *models.KepegawaianKualifikasi) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (r *repository) GetKualifikasiByID(ctx context.Context, id int64) (*models.KepegawaianKualifikasi, error) {
	var m models.KepegawaianKualifikasi
	err := r.db.WithContext(ctx).
		Preload("Tipe").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, err
}

// ── GetByPegawaiID ────────────────────────────────────────────────────────────
func (r *repository) GetKualifikasiByPegawaiID(
	ctx context.Context,
	pegawaiID int64,
	page, pageSize int,
) ([]models.KepegawaianKualifikasi, int64, error) {
	var items []models.KepegawaianKualifikasi
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.KepegawaianKualifikasi{}).
		Preload("Tipe").
		Where("pegawai_id = ? AND deleted_at IS NULL", pegawaiID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Order("tipe_id ASC,  created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&items).Error

	return items, total, err
}

// ── GetByPegawaiIDAndTipe ─────────────────────────────────────────────────────
func (r *repository) GetKualifikasiByPegawaiIDAndTipe(ctx context.Context, pegawaiID, tipeID int64) ([]models.KepegawaianKualifikasi, error) {
	var m []models.KepegawaianKualifikasi
	err := r.db.WithContext(ctx).
		Preload("Tipe").
		Where("pegawai_id = ? AND tipe_id = ? AND deleted_at IS NULL", pegawaiID, tipeID).
		Find(&m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return m, err
}

// ── GetByTipe ─────────────────────────────────────────────────────────────────
func (r *repository) GetKualifikasiByTipe(ctx context.Context, tipeID int64, page, pageSize int) ([]models.KepegawaianKualifikasi, int64, error) {
	var m []models.KepegawaianKualifikasi
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.KepegawaianKualifikasi{}).
		Preload("Tipe").
		Where("tipe_id = ? AND deleted_at IS NULL", tipeID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&m).Error

	return m, total, err
}

// ── GetExpiringSoonKualifikasi ───────────────────────────────────────────────
func (r *repository) GetExpiringSoonKualifikasi(ctx context.Context, days int, page, pageSize int) ([]models.KepegawaianKualifikasi, int64, error) {
	var m []models.KepegawaianKualifikasi
	var total int64
	deadline := time.Now().AddDate(0, 0, days)

	query := r.db.WithContext(ctx).
		Model(&models.KepegawaianKualifikasi{}).
		Preload("Tipe").
		Where(
			"tanggal_expired IS NOT NULL AND tanggal_expired <= ? AND tanggal_expired >= ? AND is_aktif = true AND deleted_at IS NULL",
			deadline,
			time.Now(),
		)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Order("tanggal_expired ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&m).Error

	return m, total, err
}

// ── GetExpiredKualifikasi ─────────────────────────────────────────────────────
func (r *repository) GetExpiredKualifikasi(ctx context.Context, page, pageSize int) ([]models.KepegawaianKualifikasi, int64, error) {
	var m []models.KepegawaianKualifikasi
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.KepegawaianKualifikasi{}).
		Preload("Tipe").
		Where("tanggal_expired < ? AND is_aktif = true AND deleted_at IS NULL", time.Now())

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Order("tanggal_expired ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&m).Error

	return m, total, err
}

// ── List ──────────────────────────────────────────────────────────────────────
func (r *repository) ListKualifikasi(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianKualifikasiRequest) ([]models.KepegawaianKualifikasi, int64, error) {
	var items []models.KepegawaianKualifikasi
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.KepegawaianKualifikasi{}).
		Preload("Tipe").
		Where("deleted_at IS NULL")

	if filter != nil {
		if filter.Nama != "" {
			query = query.Where("nama ILIKE ?", "%"+filter.Nama+"%")
		}
		if filter.TipeID != nil {
			query = query.Where("tipe_id = ?", filter.TipeID)
		}
		if filter.Penyelenggara != "" {
			query = query.Where("penyelenggara ILIKE ?", "%"+filter.Penyelenggara+"%")
		}
		if filter.IsAktif != nil {
			query = query.Where("is_aktif = ?", *filter.IsAktif)
		}
		if filter.IsExpired != nil {
			if *filter.IsExpired {
				query = query.Where("tanggal_expired IS NOT NULL AND tanggal_expired < ?", time.Now())
			} else {
				query = query.Where("tanggal_expired IS NULL OR tanggal_expired >= ?", time.Now())
			}
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// ── Update ────────────────────────────────────────────────────────────────────
func (r *repository) UpdateKualifikasi(ctx context.Context, m *models.KepegawaianKualifikasi) error {
	return r.db.WithContext(ctx).
		Model(m).
		Where("id = ? AND deleted_at IS NULL", m.ID).
		Updates(m).Error
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (r *repository) DeleteKualifikasi(ctx context.Context, id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.KepegawaianKualifikasi{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
		}).Error
}

// ── Exists ────────────────────────────────────────────────────────────────────
func (r *repository) ExistsByNomorSertifikatAndTipe(
	ctx context.Context,
	tipeID int64,
	NomorSertifikat string,
	excludeID int64,
) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&models.KepegawaianKualifikasi{}).
		Where("nomor_sertifikat = ? AND tipe_id = ? AND deleted_at IS NULL", NomorSertifikat, tipeID)
	if excludeID != 0 {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
