package repositories

import (
	"context"
	"errors"
	"time"

	"neosim_go/internal/modules/kepegawaian/alamat/dto"
	"neosim_go/internal/modules/kepegawaian/alamat/models"

	"gorm.io/gorm"
)

// ── Create ────────────────────────────────────────────────────────────────────
func (r *repository) CreateAlamat(ctx context.Context, m *models.KepegawaianAlamat) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (r *repository) GetAlamatByID(ctx context.Context, id int64) (*models.KepegawaianAlamat, error) {
	var m models.KepegawaianAlamat
	result := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── GetByPegawaiID ────────────────────────────────────────────────────────────
func (r *repository) GetAlamatByPegawaiID(ctx context.Context, pegawaiID int64, page, pageSize int) ([]models.KepegawaianAlamat, int64, error) {
	var items []models.KepegawaianAlamat
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.KepegawaianAlamat{}).
		Preload("Tipe").
		Where("pegawai_id = ? AND deleted_at IS NULL", pegawaiID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Order("tipe_id ASC, is_primary DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&items).Error

	return items, total, err
}

// ── List ──────────────────────────────────────────────────────────────────────
func (r *repository) ListAlamat(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianAlamatRequest) ([]models.KepegawaianAlamat, int64, error) {
	var items []models.KepegawaianAlamat
	var total int64

	query := r.db.WithContext(ctx).Model(&models.KepegawaianAlamat{}).Where("deleted_at IS NULL")

	if filter.Jalan != nil {
		// Ubah "name ILIKE ?" menjadi "jalan ILIKE ?"
		query = query.Where("jalan ILIKE ?", "%"+*filter.Jalan+"%")
	}

	// Ubah operator '==' menjadi '='
	if filter.NegaraID != nil {
		query = query.Where("negara_id = ?", filter.NegaraID)
	}

	if filter.ProvinsiID != nil {
		query = query.Where("provinsi_id = ?", filter.ProvinsiID)
	}

	if filter.KotaKabupatenID != nil {
		query = query.Where("kota_kabupaten_id = ?", filter.KotaKabupatenID)
	}

	if filter.KecamatanID != nil {
		query = query.Where("kecamatan_id = ?", filter.KecamatanID)
	}

	if filter.KelurahanDesaID != nil {
		query = query.Where("kelurahan_desa_id = ?", filter.KelurahanDesaID)
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
func (r *repository) UpdateAlamat(ctx context.Context, m *models.KepegawaianAlamat) error {
	return r.db.WithContext(ctx).Save(m).Error
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (r *repository) DeleteAlamat(ctx context.Context, id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.KepegawaianAlamat{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
		}).Error
}
