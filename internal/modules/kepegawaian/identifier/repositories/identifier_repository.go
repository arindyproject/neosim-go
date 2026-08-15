package repositories

import (
	"context"
	"errors"
	"time"

	"neosim_go/internal/modules/kepegawaian/identifier/dto"
	"neosim_go/internal/modules/kepegawaian/identifier/models"

	"gorm.io/gorm"
)

// ── Create ────────────────────────────────────────────────────────────────────

func (r *repository) CreateIdentifier(ctx context.Context, m *models.KepegawaianIdentifier) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ── Update ────────────────────────────────────────────────────────────────────

func (r *repository) UpdateIdentifier(ctx context.Context, m *models.KepegawaianIdentifier) error {
	return r.db.WithContext(ctx).
		Model(m).
		Where("id = ? AND deleted_at IS NULL", m.ID).
		Updates(m).Error
}

// ── Delete (soft delete) ─────────────────────────────────────────────────────

func (r *repository) DeleteIdentifier(ctx context.Context, id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.KepegawaianIdentifier{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
		}).Error
}

// ── Find ──────────────────────────────────────────────────────────────────────

func (r *repository) GetIdentifierByID(ctx context.Context, id int64) (*models.KepegawaianIdentifier, error) {
	var m models.KepegawaianIdentifier
	err := r.db.WithContext(ctx).
		Preload("Tipe").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, err
}

func (r *repository) ListIdentifier(
	ctx context.Context,
	page, pageSize int,
	filter *dto.FilterKepegawaianIdentifierRequest,
) ([]models.KepegawaianIdentifier, int64, error) {
	var items []models.KepegawaianIdentifier
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.KepegawaianIdentifier{}).
		Preload("Tipe").
		Where("deleted_at IS NULL")

	if filter != nil {
		if filter.PegawaiID != nil {
			query = query.Where("pegawai_id = ?", *filter.PegawaiID)
		}
		if filter.TipeID != nil {
			query = query.Where("tipe_id = ?", *filter.TipeID)
		}
		if filter.Nilai != "" {
			query = query.Where("nilai ILIKE ?", "%"+filter.Nilai+"%")
		}
		if filter.IsPrimary != nil {
			query = query.Where("is_primary = ?", *filter.IsPrimary)
		}
		if filter.IsAktif != nil {
			query = query.Where("is_aktif = ?", *filter.IsAktif)
		}
		// true = sudah expired, false = belum expired / tidak ada tanggal
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
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&items).Error

	return items, total, err
}

func (r *repository) FindByPegawaiID(ctx context.Context, pegawaiID int64) ([]models.KepegawaianIdentifier, error) {
	var items []models.KepegawaianIdentifier
	err := r.db.WithContext(ctx).
		Preload("Tipe").
		Where("pegawai_id = ? AND deleted_at IS NULL", pegawaiID).
		Order("tipe_id ASC, is_primary DESC, created_at DESC").
		Find(&items).Error
	return items, err
}

func (r *repository) FindByPegawaiIDAndTipe(ctx context.Context, pegawaiID, tipeID int64) ([]models.KepegawaianIdentifier, error) {
	var items []models.KepegawaianIdentifier
	err := r.db.WithContext(ctx).
		Preload("Tipe").
		Where("pegawai_id = ? AND tipe_id = ? AND deleted_at IS NULL", pegawaiID, tipeID).
		Order("is_primary DESC, created_at DESC").
		Find(&items).Error
	return items, err
}

func (r *repository) FindPrimaryByTipe(ctx context.Context, pegawaiID, tipeID int64) (*models.KepegawaianIdentifier, error) {
	var m models.KepegawaianIdentifier
	err := r.db.WithContext(ctx).
		Preload("Tipe").
		Where("pegawai_id = ? AND tipe_id = ? AND is_primary = true AND deleted_at IS NULL", pegawaiID, tipeID).
		First(&m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, err
}

func (r *repository) FindExpiringSoonIdentifier(ctx context.Context, days int) ([]models.KepegawaianIdentifier, error) {
	var items []models.KepegawaianIdentifier
	deadline := time.Now().AddDate(0, 0, days)

	err := r.db.WithContext(ctx).
		Preload("Tipe").
		Where(
			"tanggal_expired IS NOT NULL AND tanggal_expired <= ? AND tanggal_expired >= ? AND is_aktif = true AND deleted_at IS NULL",
			deadline,
			time.Now(),
		).
		Order("tanggal_expired ASC").
		Find(&items).Error

	return items, err
}

func (r *repository) FindExpiredIdentifier(ctx context.Context) ([]models.KepegawaianIdentifier, error) {
	var items []models.KepegawaianIdentifier
	err := r.db.WithContext(ctx).
		Preload("Tipe").
		Where(
			"tanggal_expired IS NOT NULL AND tanggal_expired < ? AND is_aktif = true AND deleted_at IS NULL",
			time.Now(),
		).
		Order("tanggal_expired ASC").
		Find(&items).Error
	return items, err
}

// ── Exists ────────────────────────────────────────────────────────────────────

func (r *repository) ExistsIdentifierByID(ctx context.Context, id int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.KepegawaianIdentifier{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	return count > 0, err
}

func (r *repository) ExistsByNilaiAndTipe(
	ctx context.Context,
	tipeID int64,
	nilai string,
	excludeID int64,
) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&models.KepegawaianIdentifier{}).
		Where("tipe_id = ? AND nilai = ? AND deleted_at IS NULL", tipeID, nilai)

	// excludeID > 0 saat update — exclude record sendiri
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// ── Helper ────────────────────────────────────────────────────────────────────

func (r *repository) UnsetPrimaryByPegawaiIDAndTipe(
	ctx context.Context,
	pegawaiID, tipeID int64,
	updatedBy int64,
) error {
	return r.db.WithContext(ctx).
		Model(&models.KepegawaianIdentifier{}).
		Where("pegawai_id = ? AND tipe_id = ? AND is_primary = true AND deleted_at IS NULL", pegawaiID, tipeID).
		Updates(map[string]any{
			"is_primary": false,
			"updated_by": updatedBy,
			"updated_at": time.Now(),
		}).Error
}
