package repositories

import (
	"errors"

	"neosim_go/internal/modules/artikel/artikel/dto"
	"neosim_go/internal/modules/artikel/artikel/models"

	"gorm.io/gorm"
)

// ── Create ────────────────────────────────────────────────────────────────────
func (r *repository) CreateArtikel(m *models.Artikel) error {
	return r.db.Create(m).Error
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (r *repository) GetArtikelByID(id int64) (*models.Artikel, error) {
	var m models.Artikel
	result := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── List ──────────────────────────────────────────────────────────────────────
func (r *repository) ListArtikel(page, pageSize int, filter *dto.FilterArtikelRequest) ([]models.Artikel, int64, error) {
	var items []models.Artikel
	var total int64

	query := r.db.Model(&models.Artikel{}).Where("deleted_at IS NULL")

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


// ── Update ────────────────────────────────────────────────────────────────────
func (r *repository) UpdateArtikel(m *models.Artikel) error {
	return r.db.Save(m).Error
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (r *repository) DeleteArtikel(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.Artikel{}).Error
}
