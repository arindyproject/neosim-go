package repositories

import (
	"errors"

	"neosim_go/internal/modules/master/master/dto"
	"neosim_go/internal/modules/master/master/models"

	"gorm.io/gorm"
)

// =====================================================================
// JenisKelamin
// =====================================================================
// ------------------Create---------------------------------------------
func (r *repository) CreateJenisKelamin(m *models.MasterJenisKelamin) error {
	return r.db.Create(m).Error
}

// ------------------GetByID--------------------------------------------
func (r *repository) GetByIDJenisKelamin(id int64) (*models.MasterJenisKelamin, error) {
	var m models.MasterJenisKelamin
	result := r.db.Where("id = ?", id).
		Where("master_jenis_kelamin.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ------------------GetByName------------------------------------------
func (r *repository) GetByNameJenisKelamin(name string) (*models.MasterJenisKelamin, error) {
	var m models.MasterJenisKelamin
	result := r.db.Where("name = ?", name).Where("master_jenis_kelamin.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ------------------List-----------------------------------------------
func (r *repository) ListJenisKelamin(page, pageSize int, filter *dto.FilterMasterJenisKelaminRequest) ([]models.MasterJenisKelamin, int64, error) {
	var items []models.MasterJenisKelamin
	var total int64

	query := r.db.Model(&models.MasterJenisKelamin{}).Where("master_jenis_kelamin.deleted_at IS NULL")

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
func (r *repository) UpdateJenisKelamin(m *models.MasterJenisKelamin) error {
	return r.db.Save(m).Error
}

// ------------------Delete---------------------------------------------
func (r *repository) DeleteJenisKelamin(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.MasterJenisKelamin{}).Error
} // ===================================================================
