package repositories

import (
	"context"
	"errors"
	"time"

	"neosim_go/internal/modules/kepegawaian/pegawai/contracts"
	"neosim_go/internal/modules/kepegawaian/pegawai/dto"
	"neosim_go/internal/modules/kepegawaian/pegawai/models"
	"neosim_go/internal/shared/sorting"

	"gorm.io/gorm"
)

// NewStatusRepository mengembalikan struct repository yang SAMA
// dengan repository entitas utama, dilihat sebagai contracts.StatusRepository.
// Berguna untuk test Status yang berdiri sendiri; di production cukup
// pakai repo yang sudah dibuat lewat NewKepegawaianPegawaiRepository(db).
func NewStatusRepository(db *gorm.DB) contracts.StatusRepository {
	return &repository{db: db}
}

// allowedSortColumns memetakan nilai sort_by dari client ke ekspresi kolom DB.
// Hanya key di map ini yang boleh dipakai (mencegah SQL injection lewat Order()).
var StatusSort = sorting.Config{
	Allowed: map[string]string{
		"code":       "code",
		"label":      "label",
		"fhir_code":  "fhir_code",
		"created_at": "created_at",
		"updated_at": "updated_at",
	},
	DefaultColumn: "created_at",
	DefaultDesc:   true,
	TieBreaker:    "id",
}

// ── Create ────────────────────────────────────────────────────────────────────
func (r *repository) CreateStatus(ctx context.Context, m *models.Status) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (r *repository) GetStatusByID(ctx context.Context, id int64) (*models.Status, error) {
	var m models.Status
	result := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── GetByCode ─────────────────────────────────────────────────────────────────
func (r *repository) GetStatusByCode(ctx context.Context, code string) (*models.Status, error) {
	var m models.Status
	result := r.db.WithContext(ctx).Where("code = ? AND deleted_at IS NULL", code).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── GetByLabel ────────────────────────────────────────────────────────────────
func (r *repository) GetStatusByLabel(ctx context.Context, label string) (*models.Status, error) {
	var m models.Status
	result := r.db.WithContext(ctx).Where("label = ? AND deleted_at IS NULL", label).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── ListSelect ────────────────────────────────────────────────────────────────
func (r *repository) ListSelectStatus(ctx context.Context, search string) ([]models.Status, error) {
	var items []models.Status

	query := r.db.WithContext(ctx).Model(&models.Status{}).
		Select("id, code, label").
		Where("kepegawaian_alamat_Statuss.deleted_at IS NULL")

	if search != "" {
		query = query.Where("label ILIKE ? OR code ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Order("label ASC").Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// ── List ──────────────────────────────────────────────────────────────────────
func (r *repository) ListStatus(ctx context.Context, page, pageSize int, filter *dto.FilterStatusRequest) ([]models.Status, int64, error) {
	var items []models.Status
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Status{}).Where("deleted_at IS NULL")
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
	if err := query.
		Offset(offset).
		Limit(pageSize).
		Scopes(StatusSort.Scope(filter.SortBy, filter.SortOrder)).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// ── Update ────────────────────────────────────────────────────────────────────
func (r *repository) UpdateStatus(ctx context.Context, m *models.Status) error {
	return r.db.WithContext(ctx).Save(m).Error
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (r *repository) DeleteStatus(ctx context.Context, id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.Status{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
		}).Error
}
