package repositories

import (
	"errors"

	"neosim_go/internal/modules/kepegawaian/kontak/dto"
	"neosim_go/internal/modules/kepegawaian/kontak/models"

	"gorm.io/gorm"
)

func (r *repository) CreateKontak(m *models.KepegawaianKontak) error {
	return r.db.Create(m).Error
}

func (r *repository) GetKontakByID(id int64) (*models.KepegawaianKontak, error) {
	var m models.KepegawaianKontak
	result := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

func (r *repository) ListKontak(page, pageSize int, filter *dto.FilterKepegawaianKontakRequest) ([]models.KepegawaianKontak, int64, error) {
	var items []models.KepegawaianKontak
	var total int64

	query := r.db.Model(&models.KepegawaianKontak{}).Where("deleted_at IS NULL")

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

func (r *repository) UpdateKontak(m *models.KepegawaianKontak) error {
	return r.db.Save(m).Error
}

func (r *repository) DeleteKontak(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.KepegawaianKontak{}).Error
}
