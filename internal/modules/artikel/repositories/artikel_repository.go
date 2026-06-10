package repositories

import (
	"errors"

	"neosim_go/internal/modules/artikel/contracts"
	"neosim_go/internal/modules/artikel/models"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewArtikelRepository membuat instance repository baru
func NewArtikelRepository(db *gorm.DB) contracts.Repository {
	return &repository{db: db}
}

func (r *repository) Create(m *models.Artikel) error {
	return r.db.Create(m).Error
}

func (r *repository) GetByID(id int64) (*models.Artikel, error) {
	var m models.Artikel
	result := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

func (r *repository) List(page, pageSize int) ([]models.Artikel, int64, error) {
	var items []models.Artikel
	var total int64
	offset := (page - 1) * pageSize

	if err := r.db.Model(&models.Artikel{}).Where("deleted_at IS NULL").Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Where("deleted_at IS NULL").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *repository) Update(m *models.Artikel) error {
	return r.db.Save(m).Error
}

func (r *repository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.Artikel{}).Error
}
