package repositories

import (
	"context"
	"errors"

	"neosim_go/internal/modules/master/departemen/dto"
	"neosim_go/internal/modules/master/departemen/models"

	"gorm.io/gorm"
)

// ── Create ────────────────────────────────────────────────────────────────────
func (r *repository) CreateDepartemen(ctx context.Context, m *models.MasterDepartemen) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (r *repository) GetDepartemenByID(ctx context.Context, id int64) (*models.MasterDepartemen, error) {
	var m models.MasterDepartemen
	result := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── List ──────────────────────────────────────────────────────────────────────
func (r *repository) ListDepartemen(ctx context.Context, page, pageSize int, filter *dto.FilterMasterDepartemenRequest) ([]models.MasterDepartemen, int64, error) {
	var items []models.MasterDepartemen
	var total int64

	query := r.db.WithContext(ctx).Model(&models.MasterDepartemen{}).Where("deleted_at IS NULL")

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
func (r *repository) UpdateDepartemen(ctx context.Context, m *models.MasterDepartemen) error {
	return r.db.WithContext(ctx).Save(m).Error
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (r *repository) DeleteDepartemen(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.MasterDepartemen{}).Error
}
