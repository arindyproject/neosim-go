package repositories

import (
	"errors"

	"neosim_go/internal/modules/master/alamat/dto"
	"neosim_go/internal/modules/master/alamat/models"

	"gorm.io/gorm"
)

// =====================================================================
// Kecamatan
// =====================================================================
func (r *repository) CreateKecamatan(m *models.MasterAlamatKecamatan) error {
	return r.db.Create(m).Error
}

func (r *repository) GetByIDKecamatan(id int64) (*models.MasterAlamatKecamatan, error) {
	var m models.MasterAlamatKecamatan
	result := r.db.Preload("KotaKabupaten").Where("id = ?", id).
		Where("master_alamat_kecamatan.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

func (r *repository) ListKecamatan(page, pageSize int, kotaKabupatenID *int64, filter *dto.FilterKecamatanRequest) ([]models.MasterAlamatKecamatan, int64, error) {
	var items []models.MasterAlamatKecamatan
	var total int64

	query := r.db.Model(&models.MasterAlamatKecamatan{}).
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

func (r *repository) UpdateKecamatan(m *models.MasterAlamatKecamatan) error {
	return r.db.Save(m).Error
}

func (r *repository) DeleteKecamatan(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.MasterAlamatKecamatan{}).Error
}

func (r *repository) ExistsKecamatanByCode(code string, excludeID *int64) (bool, error) {
	var count int64
	q := r.db.Model(&models.MasterAlamatKecamatan{}).Where("code = ?", code)
	if excludeID != nil {
		q = q.Where("id != ?", *excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}

// tambahan-------------------------------------------------------------
// CountDesaByKecamatanID menghitung jumlah desa/kelurahan di satu kecamatan
func (r *repository) CountDesaByKecamatanID(kecamatanID int64) (int64, error) {
	var count int64
	err := r.db.Model(&models.MasterAlamatKelurahanDesa{}).
		Where("kecamatan_id = ?", kecamatanID).
		// Tambahkan kondisi deleted_at IS NULL untuk tabel utama
		Where("master_alamat_kelurahan_desa.deleted_at IS NULL").
		Count(&count).Error
	return count, err
} //--------------------------------------------------------------------
