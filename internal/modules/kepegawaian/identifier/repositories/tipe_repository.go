package repositories

import (
	"errors"

	"neosim_go/internal/modules/kepegawaian/identifier/contracts"
	"neosim_go/internal/modules/kepegawaian/identifier/dto"
	"neosim_go/internal/modules/kepegawaian/identifier/models"

	"gorm.io/gorm"
)

// NewTipeRepository mengembalikan struct repository yang SAMA
// dengan repository entitas utama, dilihat sebagai contracts.TipeRepository.
// Berguna untuk test Tipe yang berdiri sendiri; di production cukup
// pakai repo yang sudah dibuat lewat NewKepegawaianIdentifierRepository(db).
func NewTipeRepository(db *gorm.DB) contracts.TipeRepository {
	return &repository{db: db}
}

func (r *repository) CreateTipe(m *models.Tipe) error {
	return r.db.Create(m).Error
}

func (r *repository) GetTipeByID(id int64) (*models.Tipe, error) {
	var m models.Tipe
	// GORM otomatis menambahkan 'deleted_at IS NULL'
	result := r.db.First(&m, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

func (r *repository) GetTipeByCode(code string) (*models.Tipe, error) {
	var m models.Tipe
	result := r.db.Where("code = ?", code).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

func (r *repository) GetTipeByLabel(label string) (*models.Tipe, error) {
	var m models.Tipe
	result := r.db.Where("label = ?", label).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

func (r *repository) ListTipe(page, pageSize int, filter *dto.FilterTipeRequest) ([]models.Tipe, int64, error) {
	var items []models.Tipe
	var total int64

	// GORM otomatis menangani soft delete
	query := r.db.Model(&models.Tipe{})

	if filter != nil {
		// Disesuaikan dengan field Label atau Code (misal pencarian kata kunci)
		if filter.Search != "" {
			query = query.Where("label ILIKE ? OR code ILIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
		} else {
			if filter.Label != "" {
				query = query.Where("label ILIKE ?", "%"+filter.Label+"%")
			}
			if filter.Code != "" {
				query = query.Where("code ILIKE ?", "%"+filter.Code+"%")
			}
		}

		if filter.IsNakes != nil {
			query = query.Where("is_nakes = ?", *filter.IsNakes)
		}
		if filter.IsRequired != nil {
			query = query.Where("is_required = ?", *filter.IsRequired)
		}
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

func (r *repository) UpdateTipe(m *models.Tipe) error {
	return r.db.Save(m).Error
}

func (r *repository) DeleteTipe(id int64) error {
	return r.db.Delete(&models.Tipe{}, id).Error
}
