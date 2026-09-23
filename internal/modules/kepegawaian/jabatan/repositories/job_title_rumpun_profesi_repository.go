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

// NewJobTitleRumpunProfesiRepository mengembalikan struct repository yang SAMA
// dengan repository entitas utama, dilihat sebagai contracts.JobTitleRumpunProfesiRepository.
// Berguna untuk test JobTitleRumpunProfesi yang berdiri sendiri; di production cukup
// pakai repo yang sudah dibuat lewat NewKepegawaianJabatanRepository(db).
func NewJobTitleRumpunProfesiRepository(db *gorm.DB) contracts.JobTitleRumpunProfesiRepository {
	return &repository{db: db}
}

// ── Check ────────────────────────────────────────────────────────────────────
func (r *repository) CheckJobTitleRumpunProfesi(ctx context.Context, id int64) (bool, error) {
	var exists bool
	err := r.db.WithContext(ctx).
		Model(&models.JobTitleRumpunProfesi{}).
		Select("1").
		Where("id = ?", id).
		Limit(1).
		Scan(&exists).Error

	if err != nil {
		return false, err
	}

	return exists, nil
}

// ── Create ────────────────────────────────────────────────────────────────────
func (r *repository) CreateJobTitleRumpunProfesi(ctx context.Context, m *models.JobTitleRumpunProfesi) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (r *repository) GetJobTitleRumpunProfesiByID(ctx context.Context, id int64) (*models.JobTitleRumpunProfesi, error) {
	var m models.JobTitleRumpunProfesi
	result := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── GetByCode ─────────────────────────────────────────────────────────────────
func (r *repository) GetJobTitleRumpunProfesiByCode(ctx context.Context, code string) (*models.JobTitleRumpunProfesi, error) {
	var m models.JobTitleRumpunProfesi
	result := r.db.WithContext(ctx).Where("code = ? AND deleted_at IS NULL", code).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── GetByLabel ────────────────────────────────────────────────────────────────
func (r *repository) GetJobTitleRumpunProfesiByLabel(ctx context.Context, label string) (*models.JobTitleRumpunProfesi, error) {
	var m models.JobTitleRumpunProfesi
	result := r.db.WithContext(ctx).Where("label = ? AND deleted_at IS NULL", label).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── ListSelect ────────────────────────────────────────────────────────────────
func (r *repository) ListSelectJobTitleRumpunProfesi(ctx context.Context, search string) ([]models.JobTitleRumpunProfesi, error) {
	var items []models.JobTitleRumpunProfesi

	query := r.db.WithContext(ctx).Model(&models.JobTitleRumpunProfesi{}).
		Select("id, code, label").
		Where("kepegawaian_jabatan_job_title_rumpun_profesis.deleted_at IS NULL")

	if search != "" {
		query = query.Where("label ILIKE ? OR code ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Order("label ASC").Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// ── List ──────────────────────────────────────────────────────────────────────
func (r *repository) ListJobTitleRumpunProfesi(ctx context.Context, page, pageSize int, filter *dto.FilterJobTitleRumpunProfesiRequest) ([]models.JobTitleRumpunProfesi, int64, error) {
	var items []models.JobTitleRumpunProfesi
	var total int64

	query := r.db.WithContext(ctx).Model(&models.JobTitleRumpunProfesi{}).Where("deleted_at IS NULL")
	if filter.Code != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Code+"%")
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
func (r *repository) UpdateJobTitleRumpunProfesi(ctx context.Context, m *models.JobTitleRumpunProfesi) error {
	return r.db.WithContext(ctx).Save(m).Error
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (r *repository) DeleteJobTitleRumpunProfesi(ctx context.Context, id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.JobTitleRumpunProfesi{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
		}).Error
}
