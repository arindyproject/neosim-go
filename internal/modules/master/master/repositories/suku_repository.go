package repositories

import (
	"errors"

	"neosim_go/internal/modules/master/master/dto"
	"neosim_go/internal/modules/master/master/models"

	"gorm.io/gorm"
)

// =====================================================================
// Suku
// =====================================================================
// ------------------Create---------------------------------------------
func (r *repository) CreateSuku(m *models.MasterSuku) error {
	return r.db.Create(m).Error
}

// ------------------GetByID--------------------------------------------
func (r *repository) GetByIDSuku(id int64) (*models.MasterSuku, error) {
	var m models.MasterSuku
	result := r.db.Where("id = ?", id).
		Where("master_suku.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ------------------GetByName------------------------------------------
func (r *repository) GetByNameSuku(name string) (*models.MasterSuku, error) {
	var m models.MasterSuku
	result := r.db.Where("name = ?", name).Where("master_suku.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ------------------List-----------------------------------------------
func (r *repository) ListSuku(page, pageSize int, filter *dto.FilterMasterSukuRequest) ([]models.MasterSuku, int64, error) {
	var items []models.MasterSuku
	var total int64

	query := r.db.Model(&models.MasterSuku{}).Where("master_suku.deleted_at IS NULL")

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
func (r *repository) UpdateSuku(m *models.MasterSuku) error {
	return r.db.Save(m).Error
}

// ------------------Delete---------------------------------------------
func (r *repository) DeleteSuku(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.MasterSuku{}).Error
} // ===================================================================
