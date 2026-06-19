package repositories

import (
	"errors"

	"neosim_go/internal/modules/master/master/dto"
	"neosim_go/internal/modules/master/master/models"

	"gorm.io/gorm"
)

// =====================================================================
// Agama
// =====================================================================
// ------------------Create---------------------------------------------
func (r *repository) CreateAgama(m *models.MasterAgama) error {
	return r.db.Create(m).Error
}

// ------------------GetByID--------------------------------------------
func (r *repository) GetByIDAgama(id int64) (*models.MasterAgama, error) {
	var m models.MasterAgama
	result := r.db.Where("id = ?", id).
		Where("master_agama.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ------------------GetByName------------------------------------------
func (r *repository) GetByNameAgama(name string) (*models.MasterAgama, error) {
	var m models.MasterAgama
	result := r.db.Where("name = ?", name).Where("master_agama.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ------------------List-----------------------------------------------
func (r *repository) ListAgama(page, pageSize int, filter *dto.FilterMasterAgamaRequest) ([]models.MasterAgama, int64, error) {
	var items []models.MasterAgama
	var total int64

	query := r.db.Model(&models.MasterAgama{}).Where("master_agama.deleted_at IS NULL")

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
func (r *repository) UpdateAgama(m *models.MasterAgama) error {
	return r.db.Save(m).Error
}

// ------------------Delete---------------------------------------------
func (r *repository) DeleteAgama(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.MasterAgama{}).Error
} // ===================================================================
