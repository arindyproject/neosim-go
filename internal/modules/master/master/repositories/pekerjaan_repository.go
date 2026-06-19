package repositories

import (
	"errors"

	"neosim_go/internal/modules/master/master/dto"
	"neosim_go/internal/modules/master/master/models"

	"gorm.io/gorm"
)

// =====================================================================
// Pekerjaan
// =====================================================================
// ------------------Create---------------------------------------------
func (r *repository) CreatePekerjaan(m *models.MasterPekerjaan) error {
	return r.db.Create(m).Error
}

// ------------------GetByID--------------------------------------------
func (r *repository) GetByIDPekerjaan(id int64) (*models.MasterPekerjaan, error) {
	var m models.MasterPekerjaan
	result := r.db.Where("id = ?", id).
		Where("master_pekerjaan.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ------------------GetByName------------------------------------------
func (r *repository) GetByNamePekerjaan(name string) (*models.MasterPekerjaan, error) {
	var m models.MasterPekerjaan
	result := r.db.Where("name = ?", name).Where("master_pekerjaan.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ------------------List-----------------------------------------------
func (r *repository) ListPekerjaan(page, pageSize int, filter *dto.FilterMasterPekerjaanRequest) ([]models.MasterPekerjaan, int64, error) {
	var items []models.MasterPekerjaan
	var total int64

	query := r.db.Model(&models.MasterPekerjaan{}).Where("master_pekerjaan.deleted_at IS NULL")

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
func (r *repository) UpdatePekerjaan(m *models.MasterPekerjaan) error {
	return r.db.Save(m).Error
}

// ------------------Delete---------------------------------------------
func (r *repository) DeletePekerjaan(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.MasterPekerjaan{}).Error
} // ===================================================================
