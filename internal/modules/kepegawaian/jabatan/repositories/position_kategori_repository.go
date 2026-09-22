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

// NewPositionKategoriRepository mengembalikan struct repository yang SAMA
// dengan repository entitas utama, dilihat sebagai contracts.PositionKategoriRepository.
// Berguna untuk test PositionKategori yang berdiri sendiri; di production cukup
// pakai repo yang sudah dibuat lewat NewKepegawaianJabatanRepository(db).
func NewPositionKategoriRepository(db *gorm.DB) contracts.PositionKategoriRepository {
	return &repository{db: db}
}

// ── Create ────────────────────────────────────────────────────────────────────
func (r *repository) CreatePositionKategori(ctx context.Context, m *models.PositionKategori) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (r *repository) GetPositionKategoriByID(ctx context.Context, id int64) (*models.PositionKategori, error) {
	var m models.PositionKategori
	result := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── GetByCode ─────────────────────────────────────────────────────────────────
func (r *repository) GetPositionKategoriByCode(ctx context.Context, code string) (*models.PositionKategori, error) {
	var m models.PositionKategori
	result := r.db.WithContext(ctx).Where("code = ?", code).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── GetByLabel ────────────────────────────────────────────────────────────────
func (r *repository) GetPositionKategoriByLabel(ctx context.Context, label string) (*models.PositionKategori, error) {
	var m models.PositionKategori
	result := r.db.WithContext(ctx).Where("label = ?", label).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── ListSelect ────────────────────────────────────────────────────────────────
func (r *repository) ListSelectPositionKategori(ctx context.Context, search string) ([]models.PositionKategori, error) {
	var items []models.PositionKategori

	query := r.db.WithContext(ctx).Model(&models.PositionKategori{}).
		Select("id, code, label").
		Where("kepegawaian_jabatan_position_kategoris.deleted_at IS NULL")

	if search != "" {
		query = query.Where("label ILIKE ? OR code ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Order("label ASC").Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// ── List ──────────────────────────────────────────────────────────────────────
func (r *repository) ListPositionKategori(ctx context.Context, page, pageSize int, filter *dto.FilterPositionKategoriRequest) ([]models.PositionKategori, int64, error) {
	var items []models.PositionKategori
	var total int64

	// GORM otomatis menangani soft delete
	query := r.db.WithContext(ctx).Model(&models.PositionKategori{})

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
func (r *repository) UpdatePositionKategori(ctx context.Context, m *models.PositionKategori) error {
	return r.db.WithContext(ctx).Save(m).Error
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (r *repository) DeletePositionKategori(ctx context.Context, id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.PositionKategori{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
		}).Error
}
