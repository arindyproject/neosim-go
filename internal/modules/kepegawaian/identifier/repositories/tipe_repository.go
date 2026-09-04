package repositories

import (
	"context"
	"errors"
	"time"

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

func (r *repository) CreateTipe(ctx context.Context, m *models.Tipe) error {
	return r.db.Create(m).Error
}

func (r *repository) GetTipeByID(ctx context.Context, id int64) (*models.Tipe, error) {
	var m models.Tipe
	// GORM otomatis menambahkan 'deleted_at IS NULL'
	result := r.db.WithContext(ctx).First(&m, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

func (r *repository) GetTipeByCode(ctx context.Context, code string) (*models.Tipe, error) {
	var m models.Tipe
	result := r.db.WithContext(ctx).Where("code = ?", code).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

func (r *repository) GetTipeByLabel(ctx context.Context, label string) (*models.Tipe, error) {
	var m models.Tipe
	result := r.db.WithContext(ctx).Where("label = ?", label).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

func (r *repository) ListTipe(ctx context.Context, page, pageSize int, filter *dto.FilterTipeRequest) ([]models.Tipe, int64, error) {
	var items []models.Tipe
	var total int64

	// GORM otomatis menangani soft delete
	query := r.db.WithContext(ctx).Model(&models.Tipe{})

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

func (r *repository) UpdateTipe(ctx context.Context, m *models.Tipe) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *repository) DeleteTipe(ctx context.Context, id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.Tipe{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
		}).Error
}
