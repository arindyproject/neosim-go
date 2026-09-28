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

// NewSpecializationKategoriRepository mengembalikan struct repository yang SAMA
// dengan repository entitas utama, dilihat sebagai contracts.SpecializationKategoriRepository.
// Berguna untuk test SpecializationKategori yang berdiri sendiri; di production cukup
// pakai repo yang sudah dibuat lewat NewKepegawaianJabatanRepository(db).
func NewSpecializationKategoriRepository(db *gorm.DB) contracts.SpecializationKategoriRepository {
	return &repository{db: db}
}

// ── Create ────────────────────────────────────────────────────────────────────
func (r *repository) CreateSpecializationKategori(ctx context.Context, m *models.SpecializationKategori) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (r *repository) GetSpecializationKategoriByID(ctx context.Context, id int64) (*models.SpecializationKategori, error) {
	var m models.SpecializationKategori
	result := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── GetByCode ─────────────────────────────────────────────────────────────────
func (r *repository) GetSpecializationKategoriByCode(ctx context.Context, code string) (*models.SpecializationKategori, error) {
	var m models.SpecializationKategori
	result := r.db.WithContext(ctx).Where("code = ?", code).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── GetByLabel ────────────────────────────────────────────────────────────────
func (r *repository) GetSpecializationKategoriByLabel(ctx context.Context, label string) (*models.SpecializationKategori, error) {
	var m models.SpecializationKategori
	result := r.db.WithContext(ctx).Where("label = ?", label).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── ListSelect ────────────────────────────────────────────────────────────────
func (r *repository) ListSelectSpecializationKategori(ctx context.Context, search string) ([]models.SpecializationKategori, error) {
	var items []models.SpecializationKategori

	query := r.db.WithContext(ctx).Model(&models.SpecializationKategori{}).
		Select("id, code, label").
		Where("kepegawaian_jabatan_specialization_kategoris.deleted_at IS NULL")

	if search != "" {
		query = query.Where("label ILIKE ? OR code ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Order("label ASC").Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// ── List ──────────────────────────────────────────────────────────────────────
func (r *repository) ListSpecializationKategori(ctx context.Context, page, pageSize int, filter *dto.FilterSpecializationKategoriRequest) ([]models.SpecializationKategori, int64, error) {
	var items []models.SpecializationKategori
	var total int64

	// GORM otomatis menangani soft delete
	query := r.db.WithContext(ctx).Model(&models.SpecializationKategori{})

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
func (r *repository) UpdateSpecializationKategori(ctx context.Context, m *models.SpecializationKategori) error {
	return r.db.WithContext(ctx).Save(m).Error
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (r *repository) DeleteSpecializationKategori(ctx context.Context, id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.SpecializationKategori{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
		}).Error
}
