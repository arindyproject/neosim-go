package repositories

import (
	"context"
	"errors"

	"neosim_go/internal/modules/master/master/dto"
	"neosim_go/internal/modules/master/master/models"

	"gorm.io/gorm"
)

// =====================================================================
// JenisKelamin
// =====================================================================
// ------------------Create---------------------------------------------
func (r *repository) CreateJenisKelamin(ctx context.Context, m *models.MasterJenisKelamin) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ------------------GetByID--------------------------------------------
func (r *repository) GetByIDJenisKelamin(ctx context.Context, id int64) (*models.MasterJenisKelamin, error) {
	var m models.MasterJenisKelamin
	result := r.db.WithContext(ctx).Where("id = ?", id).
		Where("master_jenis_kelamin.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ------------------GetByName------------------------------------------
func (r *repository) GetByNameJenisKelamin(ctx context.Context, name string) (*models.MasterJenisKelamin, error) {
	var m models.MasterJenisKelamin
	result := r.db.WithContext(ctx).Where("name = ?", name).Where("master_jenis_kelamin.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ------------------List-----------------------------------------------
func (r *repository) ListJenisKelamin(ctx context.Context, page, pageSize int, filter *dto.FilterMasterJenisKelaminRequest) ([]models.MasterJenisKelamin, int64, error) {
	var items []models.MasterJenisKelamin
	var total int64

	query := r.db.WithContext(ctx).Model(&models.MasterJenisKelamin{}).Where("master_jenis_kelamin.deleted_at IS NULL")

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
func (r *repository) UpdateJenisKelamin(ctx context.Context, m *models.MasterJenisKelamin) error {
	return r.db.WithContext(ctx).Save(m).Error
}

// ------------------Delete---------------------------------------------
func (r *repository) DeleteJenisKelamin(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.MasterJenisKelamin{}).Error
} // ===================================================================
