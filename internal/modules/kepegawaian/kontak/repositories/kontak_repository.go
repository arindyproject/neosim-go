package repositories

import (
	"context"
	"errors"
	"time"

	"neosim_go/internal/modules/kepegawaian/kontak/dto"
	"neosim_go/internal/modules/kepegawaian/kontak/models"

	"gorm.io/gorm"
)

// ── Create ────────────────────────────────────────────────────────────────────
func (r *repository) CreateKontak(ctx context.Context, m *models.KepegawaianKontak) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (r *repository) GetKontakByID(ctx context.Context, id int64) (*models.KepegawaianKontak, error) {
	var m models.KepegawaianKontak
	result := r.db.WithContext(ctx).Preload("Tipe").Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── GetByPegawaiID ───────────────────────────────────────────────────────────
func (r *repository) GetKontakByPegawaiID(ctx context.Context, pegawaiID int64, page, pageSize int) ([]models.KepegawaianKontak, int64, error) {
	var items []models.KepegawaianKontak
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.KepegawaianKontak{}).
		Preload("Tipe").
		Where("pegawai_id = ? AND deleted_at IS NULL", pegawaiID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Order("tipe_id ASC, is_primary DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&items).Error

	return items, total, err
}

// ── GetByPegawaiIDAndTipe ────────────────────────────────────────────────────
func (r *repository) GetByPegawaiIDAndTipe(ctx context.Context, pegawaiID, tipeID int64) ([]models.KepegawaianKontak, error) {
	var items []models.KepegawaianKontak
	err := r.db.WithContext(ctx).
		Preload("Tipe").
		Where("pegawai_id = ? AND tipe_id = ? AND deleted_at IS NULL", pegawaiID, tipeID).
		Order("is_primary DESC, created_at DESC").
		Find(&items).Error
	return items, err
}

// ── GetPrimaryByTipe ──────────────────────────────────────────────────────────
func (r *repository) GetPrimaryByTipe(ctx context.Context, pegawaiID, tipeID int64) (*models.KepegawaianKontak, error) {
	var m models.KepegawaianKontak
	err := r.db.WithContext(ctx).
		Preload("Tipe").
		Where("pegawai_id = ? AND tipe_id = ? AND is_primary = true AND deleted_at IS NULL", pegawaiID, tipeID).
		First(&m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, err
}

// ── List ──────────────────────────────────────────────────────────────────────
func (r *repository) ListKontak(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianKontakRequest) ([]models.KepegawaianKontak, int64, error) {
	var items []models.KepegawaianKontak
	var total int64

	query := r.db.WithContext(ctx).Model(&models.KepegawaianKontak{}).Preload("Tipe").Where("deleted_at IS NULL")

	if filter.PegawaiID != nil {
		query = query.Where("pegawai_id = ?", *filter.PegawaiID)
	}
	if filter.TipeID != nil {
		query = query.Where("tipe_id = ?", *filter.TipeID)
	}
	if filter.Nilai != nil {
		query = query.Where("nilai ILIKE ?", "%"+*filter.Nilai+"%")
	}
	if filter.IsPrimary != nil {
		query = query.Where("is_primary = ?", *filter.IsPrimary)
	}
	if filter.IsAktif != nil {
		query = query.Where("is_aktif = ?", *filter.IsAktif)
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

// ── Update ───────────────────────────────────────────────────────────────────
func (r *repository) UpdateKontak(ctx context.Context, m *models.KepegawaianKontak) error {
	return r.db.WithContext(ctx).Save(m).Error
}

// ── Delete ───────────────────────────────────────────────────────────────────
func (r *repository) DeleteKontak(ctx context.Context, id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.KepegawaianKontak{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
		}).Error
}

// ── ExistsByNilaiAndTipe ─────────────────────────────────────────────────────
func (r *repository) ExistsByNilaiAndTipe(
	ctx context.Context,
	tipeID int64,
	nilai string,
	excludeID int64,
) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&models.KepegawaianKontak{}).
		Where("tipe_id = ? AND nilai = ? AND deleted_at IS NULL", tipeID, nilai)

	// excludeID > 0 saat update — exclude record sendiri
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// ── UnsetPrimaryByPegawaiIDAndTipe ───────────────────────────────────────────
func (r *repository) UnsetPrimaryByPegawaiIDAndTipe(
	ctx context.Context,
	pegawaiID, tipeID int64,
	updatedBy int64,
) error {
	return r.db.WithContext(ctx).
		Model(&models.KepegawaianKontak{}).
		Where("pegawai_id = ? AND tipe_id = ? AND is_primary = true AND deleted_at IS NULL", pegawaiID, tipeID).
		Updates(map[string]any{
			"is_primary": false,
			"updated_by": updatedBy,
			"updated_at": time.Now(),
		}).Error
}
