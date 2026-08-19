package repositories

import (
	"context"
	"errors"

	"neosim_go/internal/modules/master/alamat/dto"
	"neosim_go/internal/modules/master/alamat/models"

	"gorm.io/gorm"
)

// =====================================================================
// Provinsi
// =====================================================================
func (r *repository) CreateProvinsi(ctx context.Context, m *models.MasterAlamatProvinsi) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *repository) GetByIDProvinsi(ctx context.Context, id int64) (*models.MasterAlamatProvinsi, error) {
	var m models.MasterAlamatProvinsi
	result := r.db.WithContext(ctx).Preload("Negara").Where("id = ?", id).Where("master_alamat_provinsi.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

func (r *repository) ListProvinsi(ctx context.Context, page, pageSize int, negaraID *int64, filter *dto.FilterProvinsiRequest) ([]models.MasterAlamatProvinsi, int64, error) {
	var items []models.MasterAlamatProvinsi
	var total int64

	query := r.db.WithContext(ctx).Model(&models.MasterAlamatProvinsi{}).
		Where("master_alamat_provinsi.deleted_at IS NULL")

	if negaraID != nil {
		query = query.Where("negara_id = ?", *negaraID)
	}
	if filter != nil && filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("Negara").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *repository) UpdateProvinsi(ctx context.Context, m *models.MasterAlamatProvinsi) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *repository) DeleteProvinsi(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.MasterAlamatProvinsi{}).Error
}

func (r *repository) ExistsProvinsiByCode(ctx context.Context, code string, excludeID *int64) (bool, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&models.MasterAlamatProvinsi{}).Where("code = ?", code)
	if excludeID != nil {
		q = q.Where("id != ?", *excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}

// tambahan--------------------------------------------------------------
// CountKotaByProvinsiID menghitung jumlah kota/kabupaten di satu provinsi
func (r *repository) CountKotaByProvinsiID(ctx context.Context, provinsiID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.MasterAlamatKotaKabupaten{}).
		Where("provinsi_id = ?", provinsiID).
		// Tambahkan kondisi deleted_at IS NULL untuk tabel utama
		Where("master_alamat_kota_kabupaten.deleted_at IS NULL").
		Count(&count).Error
	return count, err
}

// CountKecamatanByProvinsiID menghitung jumlah kecamatan di satu provinsi
// (lintas 1 level: kota -> provinsi)
func (r *repository) CountKecamatanByProvinsiID(ctx context.Context, provinsiID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.MasterAlamatKecamatan{}).
		Joins("JOIN master_alamat_kota_kabupaten ON master_alamat_kota_kabupaten.id = master_alamat_kecamatan.kota_kabupaten_id").
		Where("master_alamat_kota_kabupaten.provinsi_id = ?", provinsiID).
		// Tambahkan kondisi deleted_at IS NULL untuk semua tabel yang di-join
		Where("master_alamat_kecamatan.deleted_at IS NULL").
		Where("master_alamat_kota_kabupaten.deleted_at IS NULL").
		Count(&count).Error
	return count, err
}

// CountDesaByProvinsiID menghitung jumlah kelurahan/desa di satu provinsi
// (lintas 2 level: kecamatan -> kota -> provinsi, jadi pakai subquery/join)
func (r *repository) CountDesaByProvinsiID(ctx context.Context, provinsiID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.MasterAlamatKelurahanDesa{}).
		Joins("JOIN master_alamat_kecamatan ON master_alamat_kecamatan.id = master_alamat_kelurahan_desa.kecamatan_id").
		Joins("JOIN master_alamat_kota_kabupaten ON master_alamat_kota_kabupaten.id = master_alamat_kecamatan.kota_kabupaten_id").
		Where("master_alamat_kota_kabupaten.provinsi_id = ?", provinsiID).
		// Tambahkan kondisi deleted_at IS NULL untuk semua tabel yang di-join
		Where("master_alamat_kelurahan_desa.deleted_at IS NULL").
		Where("master_alamat_kecamatan.deleted_at IS NULL").
		Where("master_alamat_kota_kabupaten.deleted_at IS NULL").
		Count(&count).Error
	return count, err
} //--------------------------------------------------------------------
