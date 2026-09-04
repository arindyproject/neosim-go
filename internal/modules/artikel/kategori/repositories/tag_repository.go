package repositories

import (
	"context"
	"errors"

	"neosim_go/internal/modules/artikel/kategori/contracts"
	"neosim_go/internal/modules/artikel/kategori/dto"
	"neosim_go/internal/modules/artikel/kategori/models"

	"gorm.io/gorm"
)

// NewTagRepository mengembalikan struct repository yang SAMA
// dengan repository entitas utama, dilihat sebagai contracts.TagRepository.
// Berguna untuk test Tag yang berdiri sendiri; di production cukup
// pakai repo yang sudah dibuat lewat NewArtikelKategoriRepository(db).
func NewTagRepository(db *gorm.DB) contracts.TagRepository {
	return &repository{db: db}
}

// ── Create ────────────────────────────────────────────────────────────────────
func (r *repository) CreateTag(ctx context.Context,m *models.Tag) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (r *repository) GetTagByID(ctx context.Context,id int64) (*models.Tag, error) {
	var m models.Tag
	result := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── List ──────────────────────────────────────────────────────────────────────
func (r *repository) ListTag(ctx context.Context,page, pageSize int, filter *dto.FilterTagRequest) ([]models.Tag, int64, error) {
	var items []models.Tag
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Tag{}).Where("deleted_at IS NULL")
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
func (r *repository) UpdateTag(ctx context.Context,m *models.Tag) error {
	return r.db.WithContext(ctx).Save(m).Error
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (r *repository) DeleteTag(ctx context.Context,id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.Tag{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
	}).Error
}
