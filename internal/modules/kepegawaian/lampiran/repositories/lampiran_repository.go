package repositories

import (
	"errors"

	"neosim_go/internal/modules/kepegawaian/lampiran/dto"
	"neosim_go/internal/modules/kepegawaian/lampiran/models"

	"gorm.io/gorm"
)

func (r *repository) Create(m *models.KepegawaianLampiran) error {
	return r.db.Create(m).Error
}

func (r *repository) GetByID(id int64) (*models.KepegawaianLampiran, error) {
	var m models.KepegawaianLampiran
	result := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

func (r *repository) List(page, pageSize int, filter *dto.FilterKepegawaianLampiranRequest) ([]models.KepegawaianLampiran, int64, error) {
	var items []models.KepegawaianLampiran
	var total int64

	query := r.db.Model(&models.KepegawaianLampiran{}).Where("deleted_at IS NULL")

	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
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

func (r *repository) Update(m *models.KepegawaianLampiran) error {
	return r.db.Save(m).Error
}

func (r *repository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.KepegawaianLampiran{}).Error
}
