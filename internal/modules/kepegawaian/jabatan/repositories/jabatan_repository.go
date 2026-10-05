package repositories

import (
	"context"
	"errors"
	"time"

	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	"neosim_go/internal/shared/sorting"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// allowedSortColumns memetakan nilai sort_by dari client ke ekspresi kolom DB.
// Hanya key di map ini yang boleh dipakai (mencegah SQL injection lewat Order()).
var JabatanSort = sorting.Config{
	Allowed: map[string]string{
		"pegawai_id":        "pegawai_id",
		"department_id":     "department_id",
		"position_id":       "position_id",
		"job_title_id":      "job_title_id",
		"specialization_id": "specialization_id",
		"is_primary":        "is_primary",
		"is_aktif":          "is_aktif",
		"tanggal_mulai":     "tanggal_mulai",
		"tanggal_selesai":   "tanggal_selesai",
		"created_at":        "created_at",
		"updated_at":        "updated_at",
	},
	DefaultColumn: "created_at",
	DefaultDesc:   true,
	TieBreaker:    "id",
}

// ── Create ────────────────────────────────────────────────────────────────────
func (r *repository) CreateJabatan(ctx context.Context, m *models.KepegawaianJabatan) error {
	return r.db.WithContext(ctx).
		Omit(clause.Associations).
		Create(m).Error
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (r *repository) GetJabatanByID(ctx context.Context, id int64) (*models.KepegawaianJabatan, error) {
	var m models.KepegawaianJabatan
	err := r.db.WithContext(ctx).
		Preload("Position").
		Preload("JobTitle").
		Preload("Specialization").
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

// ── GetByPegawaiID ────────────────────────────────────────────────────────────
// Riwayat penugasan satu pegawai (aktif dan yang sudah ditutup), dengan paginasi.
func (r *repository) GetJabatanByPegawaiID(ctx context.Context, pegawaiID int64, page, pageSize int, filter *dto.FilterKepegawaianJabatanRequest) ([]models.KepegawaianJabatan, int64, error) {
	var items []models.KepegawaianJabatan
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.KepegawaianJabatan{}).
		Where("pegawai_id = ? AND deleted_at IS NULL", pegawaiID)

	if filter != nil {
		if filter.PegawaiID != nil {
			query = query.Where("pegawai_id = ?", *filter.PegawaiID)
		}
		if filter.DepartmentID != nil {
			query = query.Where("department_id = ?", *filter.DepartmentID)
		}
		if filter.PositionID != nil {
			query = query.Where("position_id = ?", *filter.PositionID)
		}
		if filter.JobTitleID != nil {
			query = query.Where("job_title_id = ?", *filter.JobTitleID)
		}
		if filter.SpecializationID != nil {
			query = query.Where("specialization_id = ?", *filter.SpecializationID)
		}
		if filter.IsPrimary != nil {
			query = query.Where("is_primary = ?", *filter.IsPrimary)
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
		Preload("Position").
		Preload("JobTitle").
		Preload("Specialization").
		Scopes(JabatanSort.Scope(filter.SortBy, filter.SortOrder)).
		Offset(offset).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// ── List ──────────────────────────────────────────────────────────────────────
func (r *repository) ListJabatan(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianJabatanRequest) ([]models.KepegawaianJabatan, int64, error) {
	var items []models.KepegawaianJabatan
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.KepegawaianJabatan{}).
		Where("deleted_at IS NULL")

	if filter != nil {
		if filter.PegawaiID != nil {
			query = query.Where("pegawai_id = ?", *filter.PegawaiID)
		}
		if filter.DepartmentID != nil {
			query = query.Where("department_id = ?", *filter.DepartmentID)
		}
		if filter.PositionID != nil {
			query = query.Where("position_id = ?", *filter.PositionID)
		}
		if filter.JobTitleID != nil {
			query = query.Where("job_title_id = ?", *filter.JobTitleID)
		}
		if filter.SpecializationID != nil {
			query = query.Where("specialization_id = ?", *filter.SpecializationID)
		}
		if filter.IsPrimary != nil {
			query = query.Where("is_primary = ?", *filter.IsPrimary)
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
		Preload("Position").
		Preload("JobTitle").
		Preload("Specialization").
		Scopes(JabatanSort.Scope(filter.SortBy, filter.SortOrder)).
		Offset(offset).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// ── ExistsPrimaryAktif ────────────────────────────────────────────────────────
// Cek apakah pegawai sudah punya jabatan primer yang aktif.
// excludeID diisi saat update (id yang sedang diedit), 0 saat create.
func (r *repository) ExistsPrimaryAktifByPegawai(ctx context.Context, pegawaiID, excludeID int64) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&models.KepegawaianJabatan{}).
		Where("pegawai_id = ? AND is_primary = true AND is_aktif = true AND deleted_at IS NULL", pegawaiID)
	if excludeID != 0 {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ── Update ────────────────────────────────────────────────────────────────────
func (r *repository) UpdateJabatan(ctx context.Context, m *models.KepegawaianJabatan) error {
	return r.db.WithContext(ctx).
		Omit(clause.Associations).
		Save(m).Error
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (r *repository) DeleteJabatan(ctx context.Context, id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.KepegawaianJabatan{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
		}).Error
}
