// repositories/pegawai_identifier_repository.go
package repositories

import (
	"context"
	"errors"
	"sync"
	"time"

	"gorm.io/gorm"

	"neosim_go/internal/modules/kepegawaian/identifier/dto"
	"neosim_go/internal/modules/kepegawaian/identifier/models"
)

// ── Create ────────────────────────────────────────────────────────────────────

func (r *repository) Create(ctx context.Context, identifier *models.KepegawaianIdentifier) error {
	return r.db.WithContext(ctx).Create(identifier).Error
}

// ── Update ────────────────────────────────────────────────────────────────────

func (r *repository) Update(ctx context.Context, identifier *models.KepegawaianIdentifier) error {
	return r.db.WithContext(ctx).
		Model(identifier).
		Where("id = ? AND deleted_at IS NULL", identifier.ID).
		Updates(identifier).Error
}

// ── Delete ────────────────────────────────────────────────────────────────────

func (r *repository) Delete(ctx context.Context, id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.KepegawaianIdentifier{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
		}).Error
}

// ── Find ──────────────────────────────────────────────────────────────────────

func (r *repository) FindByID(ctx context.Context, id int64) (*models.KepegawaianIdentifier, error) {
	var identifier models.KepegawaianIdentifier
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&identifier).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &identifier, err
}

func (r *repository) FindAll2(
	ctx context.Context,
	filter dto.FilterKepegawaianIdentifierRequest,
	page, limit int,
) ([]models.KepegawaianIdentifier, int64, error) {
	var (
		items    []models.KepegawaianIdentifier
		total    int64
		countErr error
		findErr  error
	)

	baseQuery := func() *gorm.DB {
		q := r.db.WithContext(ctx).
			Model(&models.KepegawaianIdentifier{}).
			Where("deleted_at IS NULL")
		return applyIdentifierFilter(q, filter)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		countErr = baseQuery().Count(&total).Error
	}()

	go func() {
		defer wg.Done()
		offset := (page - 1) * limit
		findErr = baseQuery().
			Order("created_at DESC").
			Offset(offset).
			Limit(limit).
			Find(&items).Error
	}()

	wg.Wait()

	if countErr != nil {
		return nil, 0, countErr
	}
	if findErr != nil {
		return nil, 0, findErr
	}

	return items, total, nil
}

func (r *repository) FindAll(
	ctx context.Context,
	filter dto.FilterKepegawaianIdentifierRequest,
	page, limit int,
) ([]models.KepegawaianIdentifier, int64, error) {
	var (
		items []models.KepegawaianIdentifier
		total int64
	)

	query := r.db.WithContext(ctx).
		Model(&models.KepegawaianIdentifier{}).
		Where("deleted_at IS NULL")

	// terapkan filter
	query = applyIdentifierFilter(query, filter)

	// hitung total sebelum pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// pagination
	offset := (page - 1) * limit
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&items).Error

	return items, total, err
}

func (r *repository) FindByPegawaiD(ctx context.Context, kepegawaianID int64) ([]models.KepegawaianIdentifier, error) {
	var items []models.KepegawaianIdentifier
	err := r.db.WithContext(ctx).
		Where("pegawai_id = ? AND deleted_at IS NULL", kepegawaianID).
		Order("tipe ASC, is_primary DESC, created_at DESC").
		Find(&items).Error
	return items, err
}

func (r *repository) FindByPegawaiDAndTipe(
	ctx context.Context,
	kepegawaianID int64,
	tipe models.IdentifierType,
) ([]models.KepegawaianIdentifier, error) {
	var items []models.KepegawaianIdentifier
	err := r.db.WithContext(ctx).
		Where("pegawai_id = ? AND tipe = ? AND deleted_at IS NULL", kepegawaianID, tipe).
		Order("is_primary DESC, created_at DESC").
		Find(&items).Error
	return items, err
}

func (r *repository) FindPrimaryByTipe(
	ctx context.Context,
	kepegawaianID int64,
	tipe models.IdentifierType,
) (*models.KepegawaianIdentifier, error) {
	var identifier models.KepegawaianIdentifier
	err := r.db.WithContext(ctx).
		Where("pegawai_id = ? AND tipe = ? AND is_primary = true AND deleted_at IS NULL", kepegawaianID, tipe).
		First(&identifier).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &identifier, err
}

func (r *repository) FindExpiringSoon(ctx context.Context, days int) ([]models.KepegawaianIdentifier, error) {
	var items []models.KepegawaianIdentifier
	deadline := time.Now().AddDate(0, 0, days)

	err := r.db.WithContext(ctx).
		Where(
			"tanggal_expired IS NOT NULL AND tanggal_expired <= ? AND tanggal_expired >= ? AND is_aktif = true AND deleted_at IS NULL",
			deadline,
			time.Now(),
		).
		Order("tanggal_expired ASC").
		Find(&items).Error

	return items, err
}

func (r *repository) FindExpired(ctx context.Context) ([]models.KepegawaianIdentifier, error) {
	var items []models.KepegawaianIdentifier
	err := r.db.WithContext(ctx).
		Where(
			"tanggal_expired IS NOT NULL AND tanggal_expired < ? AND is_aktif = true AND deleted_at IS NULL",
			time.Now(),
		).
		Order("tanggal_expired ASC").
		Find(&items).Error
	return items, err
}

// ── Exists ────────────────────────────────────────────────────────────────────

func (r *repository) ExistsByID(ctx context.Context, id int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.KepegawaianIdentifier{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	return count > 0, err
}

func (r *repository) ExistsByNilaiAndTipe(
	ctx context.Context,
	tipe models.IdentifierType,
	nilai string,
	excludeID int64,
) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&models.KepegawaianIdentifier{}).
		Where("tipe = ? AND nilai = ? AND deleted_at IS NULL", tipe, nilai)

	// excludeID > 0 saat update — exclude record sendiri
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// ── Helper ────────────────────────────────────────────────────────────────────

func (r *repository) UnsetPrimaryByPegawaiDAndTipe(
	ctx context.Context,
	kepegawaianID int64,
	tipe models.IdentifierType,
	updatedBy int64,
) error {
	return r.db.WithContext(ctx).
		Model(&models.KepegawaianIdentifier{}).
		Where("pegawai_id = ? AND tipe = ? AND is_primary = true AND deleted_at IS NULL", kepegawaianID, tipe).
		Updates(map[string]any{
			"is_primary": false,
			"updated_by": updatedBy,
			"updated_at": time.Now(),
		}).Error
}

// applyIdentifierFilter menerapkan kondisi filter ke query GORM
func applyIdentifierFilter(query *gorm.DB, filter dto.FilterKepegawaianIdentifierRequest) *gorm.DB {
	if filter.PegawaiID != nil {
		query = query.Where("pegawai_id = ?", *filter.PegawaiID)
	}
	if filter.Tipe != nil {
		query = query.Where("tipe = ?", *filter.Tipe)
	}
	if filter.IsPrimary != nil {
		query = query.Where("is_primary = ?", *filter.IsPrimary)
	}
	if filter.IsAktif != nil {
		query = query.Where("is_aktif = ?", *filter.IsAktif)
	}

	// filter is_expired: true = tanggal_expired sudah lewat, false = belum lewat atau tidak ada
	if filter.IsExpired != nil {
		if *filter.IsExpired {
			query = query.Where("tanggal_expired IS NOT NULL AND tanggal_expired < ?", time.Now())
		} else {
			query = query.Where("tanggal_expired IS NULL OR tanggal_expired >= ?", time.Now())
		}
	}

	return query
}
