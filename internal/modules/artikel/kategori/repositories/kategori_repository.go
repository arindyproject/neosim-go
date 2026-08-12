package repositories

import (
	"errors"

	"neosim_go/internal/modules/artikel/kategori/dto"
	"neosim_go/internal/modules/artikel/kategori/models"

	"gorm.io/gorm"
)

func (r *repository) CreateKategori(m *models.ArtikelKategori) error {
	return r.db.Create(m).Error
}

func (r *repository) GetKategoriByID(id int64) (*models.ArtikelKategori, error) {
	var m models.ArtikelKategori
	result := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

func (r *repository) ListKategori(page, pageSize int, filter *dto.FilterArtikelKategoriRequest) ([]models.ArtikelKategori, int64, error) {
	var items []models.ArtikelKategori
	var total int64

	query := r.db.Model(&models.ArtikelKategori{}).Where("deleted_at IS NULL")

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

func (r *repository) UpdateKategori(m *models.ArtikelKategori) error {
	return r.db.Save(m).Error
}

func (r *repository) DeleteKategori(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.ArtikelKategori{}).Error
}
