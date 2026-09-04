package repositories

import (
	"context"
	"errors"
	"time"

	"neosim_go/internal/modules/kepegawaian/kontak/contracts"
	"neosim_go/internal/modules/kepegawaian/kontak/dto"
	"neosim_go/internal/modules/kepegawaian/kontak/models"

	"gorm.io/gorm"
)

// NewTipeRepository mengembalikan struct repository yang SAMA
// dengan repository entitas utama, dilihat sebagai contracts.TipeRepository.
// Berguna untuk test Tipe yang berdiri sendiri; di production cukup
// pakai repo yang sudah dibuat lewat NewKepegawaianKontakRepository(db).
func NewTipeRepository(db *gorm.DB) contracts.TipeRepository {
	return &repository{db: db}
}

// ── Create ────────────────────────────────────────────────────────────────────
func (r *repository) CreateTipe(ctx context.Context, m *models.Tipe) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (r *repository) GetTipeByID(ctx context.Context, id int64) (*models.Tipe, error) {
	var m models.Tipe
	result := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── GetByCode ─────────────────────────────────────────────────────────────────
func (r *repository) GetTipeByCode(ctx context.Context, code string) (*models.Tipe, error) {
	var m models.Tipe
	result := r.db.WithContext(ctx).Where("code = ? AND deleted_at IS NULL", code).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── GetByLabel ────────────────────────────────────────────────────────────────
func (r *repository) GetTipeByLabel(ctx context.Context, label string) (*models.Tipe, error) {
	var m models.Tipe
	result := r.db.WithContext(ctx).Where("label = ? AND deleted_at IS NULL", label).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── List ──────────────────────────────────────────────────────────────────────
func (r *repository) ListTipe(ctx context.Context, page, pageSize int, filter *dto.FilterTipeRequest) ([]models.Tipe, int64, error) {
	var items []models.Tipe
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Tipe{}).Where("deleted_at IS NULL")
	if filter.Code != "" {
		query = query.Where("code ILIKE ?", "%"+filter.Code+"%")
	}
	if filter.Label != "" {
		query = query.Where("label ILIKE ?", "%"+filter.Label+"%")
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
func (r *repository) UpdateTipe(ctx context.Context, m *models.Tipe) error {
	return r.db.WithContext(ctx).Save(m).Error
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (r *repository) DeleteTipe(ctx context.Context, id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.Tipe{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
		}).Error
}
