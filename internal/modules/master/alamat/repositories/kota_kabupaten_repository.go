package repositories

import (
	"context"
	"errors"

	"neosim_go/internal/modules/master/alamat/dto"
	"neosim_go/internal/modules/master/alamat/models"

	"gorm.io/gorm"
)

// =====================================================================
// Kota / Kabupaten
// =====================================================================
func (r *repository) CreateKotaKabupaten(ctx context.Context, m *models.MasterAlamatKotaKabupaten) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *repository) GetByIDKotaKabupaten(ctx context.Context, id int64) (*models.MasterAlamatKotaKabupaten, error) {
	var m models.MasterAlamatKotaKabupaten
	result := r.db.WithContext(ctx).Preload("Provinsi").Where("id = ?", id).
		Where("master_alamat_kota_kabupaten.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

func (r *repository) ListKotaKabupaten(ctx context.Context, page, pageSize int, provinsiID *int64, filter *dto.FilterKotaKabupatenRequest) ([]models.MasterAlamatKotaKabupaten, int64, error) {
	var items []models.MasterAlamatKotaKabupaten
	var total int64

	query := r.db.WithContext(ctx).Model(&models.MasterAlamatKotaKabupaten{}).
		Where("master_alamat_kota_kabupaten.deleted_at IS NULL")

	if provinsiID != nil {
		query = query.Where("provinsi_id = ?", *provinsiID)
	}
	if filter != nil && filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("Provinsi").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *repository) UpdateKotaKabupaten(ctx context.Context, m *models.MasterAlamatKotaKabupaten) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *repository) DeleteKotaKabupaten(ctx context.Context, id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.MasterAlamatKotaKabupaten{}).Error
}

func (r *repository) ExistsKotaKabupatenByCode(ctx context.Context, code string, excludeID *int64) (bool, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&models.MasterAlamatKotaKabupaten{}).Where("code = ?", code)
	if excludeID != nil {
		q = q.Where("id != ?", *excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}

// tambahan--------------------------------------------------------------
// CountDesaByKotaID menghitung jumlah desa/kelurahan di satu kota/kabupaten
// (lintas 1 level: kecamatan -> kota)
// CountKecamatanByKotaID menghitung jumlah kecamatan di satu kota/kabupaten
func (r *repository) CountKecamatanByKotaID(ctx context.Context, kotaID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.MasterAlamatKecamatan{}).
		Where("kota_kabupaten_id = ?", kotaID).
		// Tambahkan kondisi deleted_at IS NULL untuk tabel utama
		Where("master_alamat_kecamatan.deleted_at IS NULL").
		Count(&count).Error
	return count, err
}

// CountDesaByKotaID menghitung jumlah desa/kelurahan di satu kota/kabupaten
// (lintas 1 level: kecamatan -> kota)
func (r *repository) CountDesaByKotaID(ctx context.Context, kotaID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.MasterAlamatKelurahanDesa{}).
		Joins("JOIN master_alamat_kecamatan ON master_alamat_kecamatan.id = master_alamat_kelurahan_desa.kecamatan_id").
		Where("master_alamat_kecamatan.kota_kabupaten_id = ?", kotaID).
		// Tambahkan kondisi deleted_at IS NULL untuk tabel utama dan tabel yang di-join
		Where("master_alamat_kelurahan_desa.deleted_at IS NULL").
		Where("master_alamat_kecamatan.deleted_at IS NULL").
		Count(&count).Error
	return count, err
} //--------------------------------------------------------------------
