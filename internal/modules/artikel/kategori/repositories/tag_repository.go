package repositories

import (
	"errors"

	"neosim_go/internal/modules/artikel/kategori/contracts"
	"neosim_go/internal/modules/artikel/kategori/dto"
	"neosim_go/internal/modules/artikel/kategori/models"

	"gorm.io/gorm"
)

type tagRepository struct {
	db *gorm.DB
}

// NewTagRepository membuat instance repository baru
func NewTagRepository(db *gorm.DB) contracts.TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Create(m *models.Tag) error {
	return r.db.Create(m).Error
}

func (r *tagRepository) GetByID(id int64) (*models.Tag, error) {
	var m models.Tag
	result := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

func (r *tagRepository) List(page, pageSize int, filter *dto.FilterTagRequest) ([]models.Tag, int64, error) {
	var items []models.Tag
	var total int64

	query := r.db.Model(&models.Tag{}).Where("deleted_at IS NULL")
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

func (r *tagRepository) Update(m *models.Tag) error {
	return r.db.Save(m).Error
}

func (r *tagRepository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.Tag{}).Error
}
