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

// NewPositionRepository mengembalikan struct repository yang SAMA
// dengan repository entitas utama, dilihat sebagai contracts.PositionRepository.
// Berguna untuk test Position yang berdiri sendiri; di production cukup
// pakai repo yang sudah dibuat lewat NewKepegawaianJabatanRepository(db).
func NewPositionRepository(db *gorm.DB) contracts.PositionRepository {
	return &repository{db: db}
}

// ── Create ────────────────────────────────────────────────────────────────────

func (r *repository) CreatePosition(ctx context.Context, m *models.Position) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ── Update ────────────────────────────────────────────────────────────────────

func (r *repository) UpdatePosition(ctx context.Context, m *models.Position) error {
	return r.db.WithContext(ctx).
		Model(m).
		Where("id = ? AND deleted_at IS NULL", m.ID).
		Updates(m).Error
}

// ── Delete (soft delete) ─────────────────────────────────────────────────────

func (r *repository) DeletePosition(ctx context.Context, id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.Position{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
		}).Error
}

// ── Find ──────────────────────────────────────────────────────────────────────

func (r *repository) GetPositionByID(ctx context.Context, id int64) (*models.Position, error) {
	var m models.Position
	err := r.db.WithContext(ctx).
		Preload("PositionKategori").
		Preload("Parent").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, err
}

// GetPositionWithChildrenByID sama seperti GetPositionByID, ditambah preload
// satu level Children — dipakai saat menampilkan node beserta anak langsungnya
// (mis. buka satu cabang bagan organisasi tanpa memuat seluruh tree).
func (r *repository) GetPositionWithChildrenByID(ctx context.Context, id int64) (*models.Position, error) {
	var m models.Position
	err := r.db.WithContext(ctx).
		Preload("PositionKategori").
		Preload("Parent").
		Preload("Children", "deleted_at IS NULL").
		Preload("Children.PositionKategori").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, err
}

func (r *repository) ListPosition(
	ctx context.Context,
	page, pageSize int,
	filter *dto.FilterPositionRequest,
) ([]models.Position, int64, error) {
	var items []models.Position
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.Position{}).
		Preload("PositionKategori").
		Where("deleted_at IS NULL")

	if filter != nil {
		if filter.Name != "" {
			query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
		}
		if filter.PositionKategoriID != nil {
			query = query.Where("position_kategori_id = ?", *filter.PositionKategoriID)
		}
		if filter.ParentID != nil {
			query = query.Where("parent_id = ?", *filter.ParentID)
		}
		if filter.DepartmentID != nil {
			query = query.Where("department_id = ?", *filter.DepartmentID)
		}
		if filter.IsAktif != nil {
			query = query.Where("is_aktif = ?", *filter.IsAktif)
		}
		// true = hanya puncak hierarki (parent_id NULL)
		if filter.IsRoot != nil && *filter.IsRoot {
			query = query.Where("parent_id IS NULL")
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Order("level_hierarki ASC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&items).Error

	return items, total, err
}

// FindChildrenByParentID mengambil semua anak langsung suatu position
// (dipakai untuk build tree bagan organisasi secara bertahap per level).
func (r *repository) FindChildrenByParentID(ctx context.Context, parentID int64) ([]models.Position, error) {
	var items []models.Position
	err := r.db.WithContext(ctx).
		Preload("PositionKategori").
		Where("parent_id = ? AND deleted_at IS NULL", parentID).
		Order("level_hierarki ASC, name ASC").
		Find(&items).Error
	return items, err
}

// FindRootPositions mengambil semua position puncak hierarki (parent_id NULL) —
// titik awal untuk merender bagan organisasi dari atas.
func (r *repository) FindRootPositions(ctx context.Context) ([]models.Position, error) {
	var items []models.Position
	err := r.db.WithContext(ctx).
		Preload("PositionKategori").
		Where("parent_id IS NULL AND deleted_at IS NULL").
		Order("level_hierarki ASC, name ASC").
		Find(&items).Error
	return items, err
}

// FindAllPositions mengambil seluruh Position (flat, dengan PositionKategori
// ter-preload) dalam satu query — dipakai service untuk merakit tree bagan
// organisasi di memory, bukan rekursi FindChildrenByParentID per level yang
// N+1 query untuk hierarki dalam.
func (r *repository) FindAllPositions(ctx context.Context, onlyAktif bool) ([]models.Position, error) {
	var items []models.Position

	query := r.db.WithContext(ctx).
		Preload("PositionKategori").
		Where("deleted_at IS NULL")

	if onlyAktif {
		query = query.Where("is_aktif = true")
	}

	err := query.
		Order("level_hierarki ASC, name ASC").
		Find(&items).Error

	return items, err
}

// ── Exists ────────────────────────────────────────────────────────────────────

func (r *repository) ExistsPositionByID(ctx context.Context, id int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Position{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	return count > 0, err
}

func (r *repository) ExistsByNameAndKategori(
	ctx context.Context,
	kategoriID int64,
	name string,
	excludeID int64,
) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&models.Position{}).
		Where("position_kategori_id = ? AND name = ? AND deleted_at IS NULL", kategoriID, name)

	// excludeID > 0 saat update — exclude record sendiri
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// HasChildren mengecek apakah suatu position masih punya anak aktif —
// dipakai service sebelum delete supaya tidak meninggalkan child dengan
// parent_id menggantung/yatim.
func (r *repository) HasChildren(ctx context.Context, id int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Position{}).
		Where("parent_id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	return count > 0, err
}
