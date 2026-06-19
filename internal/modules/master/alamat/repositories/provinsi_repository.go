package repositories

import (
	"errors"

	"neosim_go/internal/modules/master/alamat/dto"
	"neosim_go/internal/modules/master/alamat/models"

	"gorm.io/gorm"
)

// =====================================================================
// Provinsi
// =====================================================================
func (r *repository) CreateProvinsi(m *models.MasterAlamatProvinsi) error {
	return r.db.Create(m).Error
}

func (r *repository) GetByIDProvinsi(id int64) (*models.MasterAlamatProvinsi, error) {
	var m models.MasterAlamatProvinsi
	result := r.db.Preload("Negara").Where("id = ?", id).Where("master_alamat_provinsi.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

func (r *repository) ListProvinsi(page, pageSize int, negaraID *int64, filter *dto.FilterProvinsiRequest) ([]models.MasterAlamatProvinsi, int64, error) {
	var items []models.MasterAlamatProvinsi
	var total int64

	query := r.db.Model(&models.MasterAlamatProvinsi{}).
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

func (r *repository) UpdateProvinsi(m *models.MasterAlamatProvinsi) error {
	return r.db.Save(m).Error
}

func (r *repository) DeleteProvinsi(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.MasterAlamatProvinsi{}).Error
}

func (r *repository) ExistsProvinsiByCode(code string, excludeID *int64) (bool, error) {
	var count int64
	q := r.db.Model(&models.MasterAlamatProvinsi{}).Where("code = ?", code)
	if excludeID != nil {
		q = q.Where("id != ?", *excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}

// tambahan--------------------------------------------------------------
// CountKotaByProvinsiID menghitung jumlah kota/kabupaten di satu provinsi
func (r *repository) CountKotaByProvinsiID(provinsiID int64) (int64, error) {
	var count int64
	err := r.db.Model(&models.MasterAlamatKotaKabupaten{}).
		Where("provinsi_id = ?", provinsiID).
		// Tambahkan kondisi deleted_at IS NULL untuk tabel utama
		Where("master_alamat_kota_kabupaten.deleted_at IS NULL").
		Count(&count).Error
	return count, err
}

// CountKecamatanByProvinsiID menghitung jumlah kecamatan di satu provinsi
// (lintas 1 level: kota -> provinsi)
func (r *repository) CountKecamatanByProvinsiID(provinsiID int64) (int64, error) {
	var count int64
	err := r.db.Model(&models.MasterAlamatKecamatan{}).
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
func (r *repository) CountDesaByProvinsiID(provinsiID int64) (int64, error) {
	var count int64
	err := r.db.Model(&models.MasterAlamatKelurahanDesa{}).
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

// =====================================================================
// Kota / Kabupaten
// =====================================================================
func (r *repository) CreateKotaKabupaten(m *models.MasterAlamatKotaKabupaten) error {
	return r.db.Create(m).Error
}

func (r *repository) GetByIDKotaKabupaten(id int64) (*models.MasterAlamatKotaKabupaten, error) {
	var m models.MasterAlamatKotaKabupaten
	result := r.db.Preload("Provinsi").Where("id = ?", id).
		Where("master_alamat_kota_kabupaten.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

func (r *repository) ListKotaKabupaten(page, pageSize int, provinsiID *int64, filter *dto.FilterKotaKabupatenRequest) ([]models.MasterAlamatKotaKabupaten, int64, error) {
	var items []models.MasterAlamatKotaKabupaten
	var total int64

	query := r.db.Model(&models.MasterAlamatKotaKabupaten{}).
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

func (r *repository) UpdateKotaKabupaten(m *models.MasterAlamatKotaKabupaten) error {
	return r.db.Save(m).Error
}

func (r *repository) DeleteKotaKabupaten(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.MasterAlamatKotaKabupaten{}).Error
}

func (r *repository) ExistsKotaKabupatenByCode(code string, excludeID *int64) (bool, error) {
	var count int64
	q := r.db.Model(&models.MasterAlamatKotaKabupaten{}).Where("code = ?", code)
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
func (r *repository) CountKecamatanByKotaID(kotaID int64) (int64, error) {
	var count int64
	err := r.db.Model(&models.MasterAlamatKecamatan{}).
		Where("kota_kabupaten_id = ?", kotaID).
		// Tambahkan kondisi deleted_at IS NULL untuk tabel utama
		Where("master_alamat_kecamatan.deleted_at IS NULL").
		Count(&count).Error
	return count, err
}

// CountDesaByKotaID menghitung jumlah desa/kelurahan di satu kota/kabupaten
// (lintas 1 level: kecamatan -> kota)
func (r *repository) CountDesaByKotaID(kotaID int64) (int64, error) {
	var count int64
	err := r.db.Model(&models.MasterAlamatKelurahanDesa{}).
		Joins("JOIN master_alamat_kecamatan ON master_alamat_kecamatan.id = master_alamat_kelurahan_desa.kecamatan_id").
		Where("master_alamat_kecamatan.kota_kabupaten_id = ?", kotaID).
		// Tambahkan kondisi deleted_at IS NULL untuk tabel utama dan tabel yang di-join
		Where("master_alamat_kelurahan_desa.deleted_at IS NULL").
		Where("master_alamat_kecamatan.deleted_at IS NULL").
		Count(&count).Error
	return count, err
} //--------------------------------------------------------------------

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

// =====================================================================
// Desa / Kelurahan
// =====================================================================
func (r *repository) CreateKelurahanDesa(m *models.MasterAlamatKelurahanDesa) error {
	return r.db.Create(m).Error
}

func (r *repository) GetByIDKelurahanDesa(id int64) (*models.MasterAlamatKelurahanDesa, error) {
	var m models.MasterAlamatKelurahanDesa
	result := r.db.Preload("Kecamatan").Where("id = ?", id).
		Where("master_alamat_kelurahan_desa.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

func (r *repository) ListKelurahanDesa(page, pageSize int, kecamatanID *int64, filter *dto.FilterKelurahanDesaRequest) ([]models.MasterAlamatKelurahanDesa, int64, error) {
	var items []models.MasterAlamatKelurahanDesa
	var total int64

	query := r.db.Model(&models.MasterAlamatKelurahanDesa{}).
		Where("master_alamat_kelurahan_desa.deleted_at IS NULL")

	if kecamatanID != nil {
		query = query.Where("kecamatan_id = ?", *kecamatanID)
	}
	if filter != nil && filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("Kecamatan").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *repository) UpdateKelurahanDesa(m *models.MasterAlamatKelurahanDesa) error {
	return r.db.Save(m).Error
}

func (r *repository) DeleteKelurahanDesa(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.MasterAlamatKelurahanDesa{}).Error
}

func (r *repository) ExistsKelurahanDesaByCode(code string, excludeID *int64) (bool, error) {
	var count int64
	q := r.db.Model(&models.MasterAlamatKelurahanDesa{}).Where("code = ?", code)
	if excludeID != nil {
		q = q.Where("id != ?", *excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}
