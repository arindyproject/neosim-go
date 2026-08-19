package repositories

import (
	"context"
	"errors"

	"neosim_go/internal/modules/master/master/dto"
	"neosim_go/internal/modules/master/master/models"

	"gorm.io/gorm"
)

// =====================================================================
// Pendidikan
// =====================================================================
// ------------------Create---------------------------------------------
func (r *repository) CreatePendidikan(ctx context.Context, m *models.MasterPendidikan) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ------------------GetByID--------------------------------------------
func (r *repository) GetByIDPendidikan(ctx context.Context, id int64) (*models.MasterPendidikan, error) {
	var m models.MasterPendidikan
	result := r.db.WithContext(ctx).Where("id = ?", id).
		Where("master_pendidikan.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ------------------GetByName------------------------------------------
func (r *repository) GetByNamePendidikan(ctx context.Context, name string) (*models.MasterPendidikan, error) {
	var m models.MasterPendidikan
	result := r.db.WithContext(ctx).Where("name = ?", name).Where("master_pendidikan.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ------------------List-----------------------------------------------
func (r *repository) ListPendidikan(ctx context.Context, page, pageSize int, filter *dto.FilterMasterPendidikanRequest) ([]models.MasterPendidikan, int64, error) {
	var items []models.MasterPendidikan
	var total int64

	query := r.db.WithContext(ctx).Model(&models.MasterPendidikan{}).Where("master_pendidikan.deleted_at IS NULL")

	if filter != nil && filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if filter != nil && filter.KodeKemenkes != "" {
		query = query.Where("kode_kemenkes ILIKE ?", "%"+filter.KodeKemenkes+"%")
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
func (r *repository) UpdatePendidikan(ctx context.Context, m *models.MasterPendidikan) error {
	return r.db.WithContext(ctx).Save(m).Error
}

// ------------------Delete---------------------------------------------
func (r *repository) DeletePendidikan(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.MasterPendidikan{}).Error
} // ===================================================================
