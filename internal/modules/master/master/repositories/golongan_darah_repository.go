package repositories

import (
	"errors"

	"neosim_go/internal/modules/master/master/dto"
	"neosim_go/internal/modules/master/master/models"

	"gorm.io/gorm"
)

// =====================================================================
// GolonganDarah
// =====================================================================
// ------------------Create---------------------------------------------
func (r *repository) CreateGolonganDarah(m *models.MasterGolonganDarah) error {
	return r.db.Create(m).Error
}

// ------------------GetByID--------------------------------------------
func (r *repository) GetByIDGolonganDarah(id int64) (*models.MasterGolonganDarah, error) {
	var m models.MasterGolonganDarah
	result := r.db.Where("id = ?", id).
		Where("master_golongan_darah.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ------------------GetByName------------------------------------------
func (r *repository) GetByNameGolonganDarah(name string) (*models.MasterGolonganDarah, error) {
	var m models.MasterGolonganDarah
	result := r.db.Where("name = ?", name).Where("master_golongan_darah.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ------------------List-----------------------------------------------
func (r *repository) ListGolonganDarah(page, pageSize int, filter *dto.FilterMasterGolonganDarahRequest) ([]models.MasterGolonganDarah, int64, error) {
	var items []models.MasterGolonganDarah
	var total int64

	query := r.db.Model(&models.MasterGolonganDarah{}).Where("master_golongan_darah.deleted_at IS NULL")

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
func (r *repository) UpdateGolonganDarah(m *models.MasterGolonganDarah) error {
	return r.db.Save(m).Error
}

// ------------------Delete---------------------------------------------
func (r *repository) DeleteGolonganDarah(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.MasterGolonganDarah{}).Error
} // ===================================================================
