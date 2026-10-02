package repositories

import (
	"context"
	"errors"
	"time"

	"neosim_go/internal/modules/kepegawaian/jabatan/contracts"
	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	"neosim_go/internal/shared/sorting"

	"gorm.io/gorm"
)

// NewJobTitleRepository mengembalikan struct repository yang SAMA
// dengan repository entitas utama, dilihat sebagai contracts.JobTitleRepository.
// Berguna untuk test JobTitle yang berdiri sendiri; di production cukup
// pakai repo yang sudah dibuat lewat NewKepegawaianJabatanRepository(db).
func NewJobTitleRepository(db *gorm.DB) contracts.JobTitleRepository {
	return &repository{db: db}
}

// allowedSortColumns memetakan nilai sort_by dari client ke ekspresi kolom DB.
// Hanya key di map ini yang boleh dipakai (mencegah SQL injection lewat Order()).
var JobTitleSort = sorting.Config{
	Allowed: map[string]string{
		"code":              "code",
		"label":             "label",
		"kategori_id":       "kategori_id",
		"rumpun_profesi_id": "rumpun_profesi_id",
		"point":             "point",
		"memerlukan_str":    "memerlukan_str",
		"memerlukan_sip":    "memerlukan_sip",
		"is_aktif":          "is_aktif",
		"jenjang_min":       "jenjang_min",
		"fhir_code":         "fhir_code",

		"created_at": "created_at",
		"updated_at": "updated_at",
	},
	DefaultColumn: "created_at",
	DefaultDesc:   true,
	TieBreaker:    "id",
}

// ── ListSelect ────────────────────────────────────────────────────────────────
func (r *repository) ListSelectJobTitle(ctx context.Context, search string) ([]models.JobTitle, error) {
	var items []models.JobTitle

	query := r.db.WithContext(ctx).Model(&models.JobTitle{}).
		Select("id, code, label").
		Where("kepegawaian_jabatan_job_titles.deleted_at IS NULL")

	if search != "" {
		query = query.Where("label ILIKE ? OR code ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Order("label ASC").Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// ── Create ────────────────────────────────────────────────────────────────────

func (r *repository) CreateJobTitle(ctx context.Context, m *models.JobTitle) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ── Update ────────────────────────────────────────────────────────────────────

func (r *repository) UpdateJobTitle(ctx context.Context, m *models.JobTitle) error {
	return r.db.WithContext(ctx).
		Model(m).
		Where("id = ? AND deleted_at IS NULL", m.ID).
		Updates(m).Error
}

// ── Delete (soft delete) ─────────────────────────────────────────────────────

func (r *repository) DeleteJobTitle(ctx context.Context, id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.JobTitle{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
		}).Error
}

// ── Find ──────────────────────────────────────────────────────────────────────

func (r *repository) GetJobTitleByID(ctx context.Context, id int64) (*models.JobTitle, error) {
	var m models.JobTitle
	err := r.db.WithContext(ctx).
		Preload("Kategori").
		Preload("RumpunProfesi").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, err
}

func (r *repository) ListJobTitle(
	ctx context.Context,
	page, pageSize int,
	filter *dto.FilterJobTitleRequest,
) ([]models.JobTitle, int64, error) {
	var items []models.JobTitle
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.JobTitle{}).
		Preload("Kategori").
		Preload("RumpunProfesi").
		Where("deleted_at IS NULL")

	if filter != nil {
		if filter.Label != "" {
			query = query.Where("label ILIKE ?", "%"+filter.Label+"%")
		}
		if filter.Code != "" {
			query = query.Where("code ILIKE ?", "%"+filter.Code+"%")
		}
		if filter.KategoriID != nil {
			query = query.Where("kategori_id = ?", *filter.KategoriID)
		}
		if filter.RumpunProfesiID != nil {
			query = query.Where("rumpun_profesi_id = ?", *filter.RumpunProfesiID)
		}
		if filter.MemerlukanSTR != nil {
			query = query.Where("memerlukan_str = ?", *filter.MemerlukanSTR)
		}
		if filter.MemerlukanSIP != nil {
			query = query.Where("memerlukan_sip = ?", *filter.MemerlukanSIP)
		}
		if filter.IsAktif != nil {
			query = query.Where("is_aktif = ?", *filter.IsAktif)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Scopes(JobTitleSort.Scope(filter.SortBy, filter.SortOrder)).
		Offset(offset).
		Limit(pageSize).
		Find(&items).Error

	if err != nil {
		return nil, 0, err
	}

	return items, total, err
}

// ── Exists ────────────────────────────────────────────────────────────────────

func (r *repository) ExistsJobTitleByID(ctx context.Context, id int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.JobTitle{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	return count > 0, err
}

// ExistsByCode mengecek duplikat 'code' — kolom ini sudah unique di level DB,
// tapi dicek dulu di sini supaya errornya jadi pesan yang jelas (422), bukan
// constraint violation mentah dari Postgres.
func (r *repository) ExistsByCode(ctx context.Context, code string, excludeID int64) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&models.JobTitle{}).
		Where("code = ? AND deleted_at IS NULL", code)

	// excludeID > 0 saat update — exclude record sendiri
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}
