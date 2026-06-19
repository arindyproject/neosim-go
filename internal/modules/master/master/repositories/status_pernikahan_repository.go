package repositories

import (
	"errors"

	"neosim_go/internal/modules/master/master/dto"
	"neosim_go/internal/modules/master/master/models"

	"gorm.io/gorm"
)

// =====================================================================
// StatusPernikahan
// =====================================================================
// ------------------Create---------------------------------------------
func (r *repository) CreateStatusPernikahan(m *models.MasterStatusPernikahan) error {
	return r.db.Create(m).Error
}

// ------------------GetByID--------------------------------------------
func (r *repository) GetByIDStatusPernikahan(id int64) (*models.MasterStatusPernikahan, error) {
	var m models.MasterStatusPernikahan
	result := r.db.Where("id = ?", id).
		Where("master_status_pernikahan.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ------------------GetByName------------------------------------------
func (r *repository) GetByNameStatusPernikahan(name string) (*models.MasterStatusPernikahan, error) {
	var m models.MasterStatusPernikahan
	result := r.db.Where("name = ?", name).Where("master_status_pernikahan.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ------------------List-----------------------------------------------
// PERBAIKAN TYPO: LisStatusPernikahan -> ListStatusPernikahan
func (r *repository) ListStatusPernikahan(page, pageSize int, filter *dto.FilterMasterStatusPernikahanRequest) ([]models.MasterStatusPernikahan, int64, error) {
	var items []models.MasterStatusPernikahan
	var total int64

	query := r.db.Model(&models.MasterStatusPernikahan{}).Where("master_status_pernikahan.deleted_at IS NULL")

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
func (r *repository) UpdateStatusPernikahan(m *models.MasterStatusPernikahan) error {
	return r.db.Save(m).Error
}

// ------------------Delete---------------------------------------------
func (r *repository) DeleteStatusPernikahan(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.MasterStatusPernikahan{}).Error
} // ===================================================================
