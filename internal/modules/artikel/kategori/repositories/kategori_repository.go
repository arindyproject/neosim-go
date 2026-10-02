package repositories

import (
	"context"
	"errors"
	"time"

	"neosim_go/internal/modules/artikel/kategori/dto"
	"neosim_go/internal/modules/artikel/kategori/models"

	"gorm.io/gorm"

	"neosim_go/internal/shared/sorting"

)

// allowedSortColumns memetakan nilai sort_by dari client ke ekspresi kolom DB.
// Hanya key di map ini yang boleh dipakai (mencegah SQL injection lewat Order()).
var KategoriSort = sorting.Config{
	Allowed: map[string]string{
		"name":          "name",
		"created_at":    "created_at",
		"updated_at":    "updated_at",
	},
	DefaultColumn: "created_at",
	DefaultDesc:   true,
	TieBreaker:    "id",
}

// ── Create ────────────────────────────────────────────────────────────────────
func (r *repository) CreateKategori(ctx context.Context, m *models.ArtikelKategori) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (r *repository) GetKategoriByID(ctx context.Context, id int64) (*models.ArtikelKategori, error) {
	var m models.ArtikelKategori
	result := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── List ──────────────────────────────────────────────────────────────────────
func (r *repository) ListKategori(ctx context.Context, page, pageSize int, filter *dto.FilterArtikelKategoriRequest) ([]models.ArtikelKategori, int64, error) {
	var items []models.ArtikelKategori
	var total int64

	query := r.db.WithContext(ctx).Model(&models.ArtikelKategori{}).Where("deleted_at IS NULL")

	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.
		Scopes(KategoriSort.Scope(filter.SortBy, filter.SortOrder)).
		Offset(offset).
		Limit(pageSize).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}


// ── Update ────────────────────────────────────────────────────────────────────
func (r *repository) UpdateKategori(ctx context.Context, m *models.ArtikelKategori) error {
	return r.db.WithContext(ctx).Save(m).Error
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (r *repository) DeleteKategori(ctx context.Context, id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.ArtikelKategori{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
		}).Error
}
