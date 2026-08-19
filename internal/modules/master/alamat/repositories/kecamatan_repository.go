package repositories

import (
	"context"
	"errors"

	"neosim_go/internal/modules/master/alamat/dto"
	"neosim_go/internal/modules/master/alamat/models"

	"gorm.io/gorm"
)

// =====================================================================
// Kecamatan
// =====================================================================
func (r *repository) CreateKecamatan(ctx context.Context, m *models.MasterAlamatKecamatan) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *repository) GetByIDKecamatan(ctx context.Context, id int64) (*models.MasterAlamatKecamatan, error) {
	var m models.MasterAlamatKecamatan
	result := r.db.WithContext(ctx).Preload("KotaKabupaten").Where("id = ?", id).
		Where("master_alamat_kecamatan.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

func (r *repository) ListKecamatan(ctx context.Context, page, pageSize int, kotaKabupatenID *int64, filter *dto.FilterKecamatanRequest) ([]models.MasterAlamatKecamatan, int64, error) {
	var items []models.MasterAlamatKecamatan
	var total int64

	query := r.db.WithContext(ctx).Model(&models.MasterAlamatKecamatan{}).
		Where("master_alamat_kecamatan.deleted_at IS NULL")

	if kotaKabupatenID != nil {
		query = query.Where("kota_kabupaten_id = ?", *kotaKabupatenID)
	}
	if filter != nil && filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("KotaKabupaten").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *repository) UpdateKecamatan(ctx context.Context, m *models.MasterAlamatKecamatan) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *repository) DeleteKecamatan(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.MasterAlamatKecamatan{}).Error
}

func (r *repository) ExistsKecamatanByCode(ctx context.Context, code string, excludeID *int64) (bool, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&models.MasterAlamatKecamatan{}).Where("code = ?", code)
	if excludeID != nil {
		q = q.Where("id != ?", *excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}

// tambahan-------------------------------------------------------------
// CountDesaByKecamatanID menghitung jumlah desa/kelurahan di satu kecamatan
func (r *repository) CountDesaByKecamatanID(ctx context.Context, kecamatanID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.MasterAlamatKelurahanDesa{}).
		Where("kecamatan_id = ?", kecamatanID).
		// Tambahkan kondisi deleted_at IS NULL untuk tabel utama
		Where("master_alamat_kelurahan_desa.deleted_at IS NULL").
		Count(&count).Error
	return count, err
} //--------------------------------------------------------------------
