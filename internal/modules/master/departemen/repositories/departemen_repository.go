package repositories

import (
	"context"
	"errors"
	"time"

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

// ── ListSelect ────────────────────────────────────────────────────────────────
func (r *repository) ListSelectDepartemen(ctx context.Context, search string) ([]models.MasterDepartemen, error) {
	var items []models.MasterDepartemen
	query := r.db.WithContext(ctx).Model(&models.MasterDepartemen{}).
		Select("id, name").
		Where("deleted_at IS NULL")
	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// ── List ──────────────────────────────────────────────────────────────────────
func (r *repository) ListDepartemen(ctx context.Context, page, pageSize int, filter *dto.FilterMasterDepartemenRequest) ([]models.MasterDepartemen, int64, error) {
	var items []models.MasterDepartemen
	var total int64

	query := r.db.WithContext(ctx).Model(&models.MasterDepartemen{}).Where("deleted_at IS NULL")

	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if filter.FhirCode != nil {
		query = query.Where("fhir_code ILIKE ?", "%"+*filter.FhirCode+"%")
	}

	if filter.FhirSystem != nil {
		query = query.Where("fhir_system ILIKE ?", "%"+*filter.FhirSystem+"%")
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
func (r *repository) DeleteDepartemen(ctx context.Context, id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.MasterDepartemen{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
		}).Error
}
