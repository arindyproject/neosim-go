package repositories

import (
	"context"
	"errors"
	"time"

	"neosim_go/internal/modules/kepegawaian/jabatan/contracts"
	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"

	"gorm.io/gorm"
)

// NewSpecializationRepository mengembalikan struct repository yang SAMA
// dengan repository entitas utama, dilihat sebagai contracts.SpecializationRepository.
// Berguna untuk test Specialization yang berdiri sendiri; di production cukup
// pakai repo yang sudah dibuat lewat NewKepegawaianJabatanRepository(db).
func NewSpecializationRepository(db *gorm.DB) contracts.SpecializationRepository {
	return &repository{db: db}
}

// ── Create ────────────────────────────────────────────────────────────────────
func (r *repository) CreateSpecialization(ctx context.Context,m *models.Specialization) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (r *repository) GetSpecializationByID(ctx context.Context,id int64) (*models.Specialization, error) {
	var m models.Specialization
	result := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── List ──────────────────────────────────────────────────────────────────────
func (r *repository) ListSpecialization(ctx context.Context,page, pageSize int, filter *dto.FilterSpecializationRequest) ([]models.Specialization, int64, error) {
	var items []models.Specialization
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Specialization{}).Where("deleted_at IS NULL")
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
func (r *repository) UpdateSpecialization(ctx context.Context,m *models.Specialization) error {
	return r.db.WithContext(ctx).Save(m).Error
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (r *repository) DeleteSpecialization(ctx context.Context,id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.Specialization{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
	}).Error
}
