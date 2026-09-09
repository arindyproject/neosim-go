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
// Suku
// =====================================================================
// ------------------Create---------------------------------------------
func (r *repository) CreateSuku(ctx context.Context, m *models.MasterSuku) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ------------------GetByID--------------------------------------------
func (r *repository) GetByIDSuku(ctx context.Context, id int64) (*models.MasterSuku, error) {
	var m models.MasterSuku
	result := r.db.WithContext(ctx).Where("id = ?", id).
		Where("master_suku.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ------------------GetByName------------------------------------------
func (r *repository) GetByNameSuku(ctx context.Context, name string) (*models.MasterSuku, error) {
	var m models.MasterSuku
	result := r.db.WithContext(ctx).Where("name = ?", name).Where("master_suku.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ------------------ListSelect-----------------------------------------
func (r *repository) ListSelectSuku(ctx context.Context, search string) ([]models.MasterSuku, error) {
	var items []models.MasterSuku

	query := r.db.WithContext(ctx).Model(&models.MasterSuku{}).
		Select("id, name, kode_kemenkes").
		Where("master_suku.deleted_at IS NULL")

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := query.Order("name ASC").Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// ------------------List-----------------------------------------------
func (r *repository) ListSuku(ctx context.Context, page, pageSize int, filter *dto.FilterMasterSukuRequest) ([]models.MasterSuku, int64, error) {
	var items []models.MasterSuku
	var total int64

	query := r.db.WithContext(ctx).Model(&models.MasterSuku{}).Where("master_suku.deleted_at IS NULL")

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
func (r *repository) UpdateSuku(ctx context.Context, m *models.MasterSuku) error {
	return r.db.WithContext(ctx).Save(m).Error
}

// ------------------Delete---------------------------------------------
func (r *repository) DeleteSuku(ctx context.Context, id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.MasterSuku{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
		}).Error
}

// ===================================================================
