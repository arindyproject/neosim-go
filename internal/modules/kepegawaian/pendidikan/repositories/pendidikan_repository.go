package repositories

import (
	"context"
	"errors"

	"neosim_go/internal/modules/kepegawaian/pendidikan/dto"
	"neosim_go/internal/modules/kepegawaian/pendidikan/models"

	"gorm.io/gorm"
)

// ── Create ────────────────────────────────────────────────────────────────────
func (r *repository) CreatePendidikan(ctx context.Context, m *models.KepegawaianPendidikan) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (r *repository) GetPendidikanByID(ctx context.Context, id int64) (*models.KepegawaianPendidikan, error) {
	var m models.KepegawaianPendidikan
	result := r.db.WithContext(ctx).Preload("Jenjang").Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── GetByPegawaiID ────────────────────────────────────────────────────────────
func (r *repository) GetPendidikanByPegawaiID(ctx context.Context,
	pegawaiID int64,
	page, pageSize int,
) ([]models.KepegawaianPendidikan, int64, error) {
	var items []models.KepegawaianPendidikan
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.KepegawaianPendidikan{}).
		Preload("Jenjang").
		Where("pegawai_id = ? AND deleted_at IS NULL", pegawaiID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Order("jenjang_id ASC, is_primary DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&items).Error

	return items, total, err
}

// ── GetByPegawaiIDAndTipe ─────────────────────────────────────────────────────
func (r *repository) GetByPegawaiIDAndTipe(ctx context.Context, pegawaiID, jenjangID int64) ([]models.KepegawaianPendidikan, error) {
	var items []models.KepegawaianPendidikan
	err := r.db.WithContext(ctx).
		Preload("Jenjang").
		Where("pegawai_id = ? AND jenjang_id = ? AND deleted_at IS NULL", pegawaiID, jenjangID).
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}

// ── List ──────────────────────────────────────────────────────────────────────
func (r *repository) ListPendidikan(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianPendidikanRequest) ([]models.KepegawaianPendidikan, int64, error) {
	var items []models.KepegawaianPendidikan
	var total int64

	query := r.db.WithContext(ctx).Model(&models.KepegawaianPendidikan{}).Preload("Jenjang").Where("deleted_at IS NULL")

	if filter != nil {
		if filter.PegawaiID != nil {
			query = query.Where("pegawai_id = ?", *filter.PegawaiID)
		}
		if filter.JenjangID != nil {
			query = query.Where("jenjang_id = ?", *filter.JenjangID)
		}
		if filter.NamaInstitusi != "" {
			query = query.Where("nama_institusi ILIKE ?", "%"+filter.NamaInstitusi+"%")
		}
		if filter.AlamatInstitusi != "" {
			query = query.Where("alamat_institusi ILIKE ?", "%"+filter.NamaInstitusi+"%")
		}
		if filter.BidangStudi != "" {
			query = query.Where("bidang_studi ILIKE ?", "%"+filter.NamaInstitusi+"%")
		}
		if filter.NomorIjazah != "" {
			query = query.Where("nomor_ijazah = ?", filter.NomorIjazah)
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
func (r *repository) UpdatePendidikan(ctx context.Context, m *models.KepegawaianPendidikan) error {
	return r.db.WithContext(ctx).Save(m).Error
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (r *repository) DeletePendidikan(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.KepegawaianPendidikan{}).Error
}

// Exists ───────────────────────────────────────────────────────────────────────

// ── ExistsPendidikanByID ──────────────────────────────────────────────────────
func (r *repository) ExistsPendidikanByID(ctx context.Context, id int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.KepegawaianPendidikan{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	return count > 0, err
}

// ── ExistsByNomorIjazah ───────────────────────────────────────────────────────
func (r *repository) ExistsByNomorIjazah(
	ctx context.Context,
	jenjangID int64,
	nomorIjazah string,
	excludeID int64,
) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&models.KepegawaianPendidikan{}).
		Where("jenjang_id = ? AND nomor_ijazah = ? AND deleted_at IS NULL", jenjangID, nomorIjazah)

	// excludeID > 0 saat update — exclude record sendiri
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}
