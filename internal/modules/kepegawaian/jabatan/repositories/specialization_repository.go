package repositories

import (
	"context"
	"errors"
	"time"

	"neosim_go/internal/modules/kepegawaian/jabatan/contracts"
	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// NewSpecializationRepository mengembalikan struct repository yang SAMA
// dengan repository entitas utama, dilihat sebagai contracts.SpecializationRepository.
// Berguna untuk test Specialization yang berdiri sendiri; di production cukup
// pakai repo yang sudah dibuat lewat NewKepegawaianJabatanRepository(db).
func NewSpecializationRepository(db *gorm.DB) contracts.SpecializationRepository {
	return &repository{db: db}
}

// ── Create ────────────────────────────────────────────────────────────────────
func (r *repository) CreateSpecialization(ctx context.Context, m *models.Specialization) error {
	return r.db.WithContext(ctx).
		Omit(clause.Associations).
		Create(m).Error
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (r *repository) GetSpecializationByID(ctx context.Context, id int64) (*models.Specialization, error) {
	var m models.Specialization
	err := r.db.WithContext(ctx).
		Preload("JobTitle").
		Preload("Kategori").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// ── GetByCode ─────────────────────────────────────────────────────────────────
func (r *repository) GetSpecializationByCode(ctx context.Context, code string) (*models.Specialization, error) {
	var m models.Specialization
	err := r.db.WithContext(ctx).
		Preload("JobTitle").
		Preload("Kategori").
		Where("code = ? AND deleted_at IS NULL", code).
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// ── List ──────────────────────────────────────────────────────────────────────
func (r *repository) ListSpecialization(ctx context.Context, page, pageSize int, filter *dto.FilterSpecializationRequest) ([]models.Specialization, int64, error) {
	var items []models.Specialization
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Specialization{}).Where("deleted_at IS NULL")

	if filter != nil {
		if filter.Code != "" {
			query = query.Where("code ILIKE ?", "%"+filter.Code+"%")
		}
		if filter.Label != "" {
			query = query.Where("label ILIKE ?", "%"+filter.Label+"%")
		}
		if filter.JobTitleID != nil {
			query = query.Where("job_title_id = ?", *filter.JobTitleID)
		}
		if filter.KategoriID != nil {
			query = query.Where("kategori_id = ?", *filter.KategoriID)
		}
		if filter.IsAktif != nil {
			query = query.Where("is_aktif = ?", *filter.IsAktif)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.
		Preload("JobTitle").
		Preload("Kategori").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// ── ListSelect ────────────────────────────────────────────────────────────────
func (r *repository) ListSelectSpecialization(ctx context.Context, search string) ([]models.Specialization, error) {
	var items []models.Specialization

	query := r.db.WithContext(ctx).Model(&models.Specialization{}).
		Select("id, code, label, gelar").
		Where("kepegawaian_jabatan_specializations.deleted_at IS NULL").
		Where("is_aktif = ?", true)

	if search != "" {
		like := "%" + search + "%"
		query = query.Where("label ILIKE ? OR code ILIKE ? OR gelar ILIKE ?", like, like, like)
	}

	if err := query.Order("label ASC").Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// ── Update ────────────────────────────────────────────────────────────────────
func (r *repository) UpdateSpecialization(ctx context.Context, m *models.Specialization) error {
	return r.db.WithContext(ctx).
		Omit(clause.Associations).
		Save(m).Error
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (r *repository) DeleteSpecialization(ctx context.Context, id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.Specialization{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
		}).Error
}
