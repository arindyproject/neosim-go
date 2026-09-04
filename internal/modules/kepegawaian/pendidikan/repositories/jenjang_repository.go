package repositories

import (
	"context"
	"errors"
	"time"

	"neosim_go/internal/modules/kepegawaian/pendidikan/contracts"
	"neosim_go/internal/modules/kepegawaian/pendidikan/dto"
	"neosim_go/internal/modules/kepegawaian/pendidikan/models"

	"gorm.io/gorm"
)

// NewJenjangRepository mengembalikan struct repository yang SAMA
// dengan repository entitas utama, dilihat sebagai contracts.JenjangRepository.
// Berguna untuk test Jenjang yang berdiri sendiri; di production cukup
// pakai repo yang sudah dibuat lewat NewKepegawaianPendidikanRepository(db).
func NewJenjangRepository(db *gorm.DB) contracts.JenjangRepository {
	return &repository{db: db}
}

// ── Create ────────────────────────────────────────────────────────────────────
func (r *repository) CreateJenjang(ctx context.Context, m *models.Jenjang) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (r *repository) GetJenjangByID(ctx context.Context, id int64) (*models.Jenjang, error) {
	var m models.Jenjang
	result := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── List ──────────────────────────────────────────────────────────────────────
func (r *repository) ListJenjang(ctx context.Context, page, pageSize int, filter *dto.FilterJenjangRequest) ([]models.Jenjang, int64, error) {
	var items []models.Jenjang
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Jenjang{}).Where("deleted_at IS NULL")
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
func (r *repository) UpdateJenjang(ctx context.Context, m *models.Jenjang) error {
	return r.db.WithContext(ctx).Save(m).Error
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (r *repository) DeleteJenjang(ctx context.Context, id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.Jenjang{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
		}).Error
}
