package repositories

import (
	"context"
	"errors"
	"time"

	"neosim_go/internal/modules/master/master/dto"
	"neosim_go/internal/modules/master/master/models"

	"gorm.io/gorm"
)

// =====================================================================
// GolonganDarah
// =====================================================================
// ------------------Create---------------------------------------------
func (r *repository) CreateGolonganDarah(ctx context.Context, m *models.MasterGolonganDarah) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ------------------GetByID--------------------------------------------
func (r *repository) GetByIDGolonganDarah(ctx context.Context, id int64) (*models.MasterGolonganDarah, error) {
	var m models.MasterGolonganDarah
	result := r.db.WithContext(ctx).Where("id = ?", id).
		Where("master_golongan_darah.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ------------------GetByName------------------------------------------
func (r *repository) GetByNameGolonganDarah(ctx context.Context, name string) (*models.MasterGolonganDarah, error) {
	var m models.MasterGolonganDarah
	result := r.db.WithContext(ctx).Where("name = ?", name).Where("master_golongan_darah.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ------------------ListSelect-----------------------------------------
func (r *repository) ListSelectGolonganDarah(ctx context.Context, search string) ([]models.MasterGolonganDarah, error) {
	var items []models.MasterGolonganDarah

	query := r.db.WithContext(ctx).Model(&models.MasterGolonganDarah{}).
		Select("id, name, kode_kemenkes").
		Where("master_golongan_darah.deleted_at IS NULL")

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := query.Order("name ASC").Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// ------------------List-----------------------------------------------
func (r *repository) ListGolonganDarah(ctx context.Context, page, pageSize int, filter *dto.FilterMasterGolonganDarahRequest) ([]models.MasterGolonganDarah, int64, error) {
	var items []models.MasterGolonganDarah
	var total int64

	query := r.db.WithContext(ctx).Model(&models.MasterGolonganDarah{}).Where("master_golongan_darah.deleted_at IS NULL")

	if filter != nil && filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if filter != nil && filter.KodeKemenkes != "" {
		query = query.Where("kode_kemenkes ILIKE ?", "%"+filter.KodeKemenkes+"%")
	}

	if filter != nil && filter.FhirCode != "" {
		query = query.Where("fhir_code ILIKE ?", "%"+filter.FhirCode+"%")
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

// ------------------Update---------------------------------------------
func (r *repository) UpdateGolonganDarah(ctx context.Context, m *models.MasterGolonganDarah) error {
	return r.db.WithContext(ctx).Save(m).Error
}

// ------------------Delete---------------------------------------------
func (r *repository) DeleteGolonganDarah(ctx context.Context, id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.MasterGolonganDarah{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
		}).Error
}

// ===================================================================
