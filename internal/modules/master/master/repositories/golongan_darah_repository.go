package repositories

import (
	"context"
	"errors"

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
func (r *repository) DeleteGolonganDarah(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.MasterGolonganDarah{}).Error
} // ===================================================================
