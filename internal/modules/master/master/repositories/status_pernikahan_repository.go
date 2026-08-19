package repositories

import (
	"context"
	"errors"

	"neosim_go/internal/modules/master/master/dto"
	"neosim_go/internal/modules/master/master/models"

	"gorm.io/gorm"
)

// =====================================================================
// StatusPernikahan
// =====================================================================
// ------------------Create---------------------------------------------
func (r *repository) CreateStatusPernikahan(ctx context.Context, m *models.MasterStatusPernikahan) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ------------------GetByID--------------------------------------------
func (r *repository) GetByIDStatusPernikahan(ctx context.Context, id int64) (*models.MasterStatusPernikahan, error) {
	var m models.MasterStatusPernikahan
	result := r.db.WithContext(ctx).Where("id = ?", id).
		Where("master_status_pernikahan.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ------------------GetByName------------------------------------------
func (r *repository) GetByNameStatusPernikahan(ctx context.Context, name string) (*models.MasterStatusPernikahan, error) {
	var m models.MasterStatusPernikahan
	result := r.db.WithContext(ctx).Where("name = ?", name).Where("master_status_pernikahan.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ------------------List-----------------------------------------------
// PERBAIKAN TYPO: LisStatusPernikahan -> ListStatusPernikahan
func (r *repository) ListStatusPernikahan(ctx context.Context, page, pageSize int, filter *dto.FilterMasterStatusPernikahanRequest) ([]models.MasterStatusPernikahan, int64, error) {
	var items []models.MasterStatusPernikahan
	var total int64

	query := r.db.WithContext(ctx).Model(&models.MasterStatusPernikahan{}).Where("master_status_pernikahan.deleted_at IS NULL")

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
func (r *repository) UpdateStatusPernikahan(ctx context.Context, m *models.MasterStatusPernikahan) error {
	return r.db.WithContext(ctx).Save(m).Error
}

// ------------------Delete---------------------------------------------
func (r *repository) DeleteStatusPernikahan(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.MasterStatusPernikahan{}).Error
} // ===================================================================
