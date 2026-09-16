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

// ── FindByPegawaiID ───────────────────────────────────────────────────────────
// FindAlamatByPegawaiID mengambil SEMUA alamat aktif (belum dihapus) milik satu
// pegawai tanpa pagination. Dipakai untuk pengecekan sebelum menghapus alamat
// primary — memastikan masih ada alamat lain milik pegawai yang sama.
func (r *repository) FindAlamatByPegawaiID(ctx context.Context, pegawaiID int64) ([]models.KepegawaianAlamat, error) {
	var items []models.KepegawaianAlamat
	err := r.db.WithContext(ctx).
		Where("pegawai_id = ? AND deleted_at IS NULL", pegawaiID).
		Find(&items).Error
	return items, err
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

// ── UnsetPrimary ──────────────────────────────────────────────────────────────
// UnsetPrimaryAlamatByPegawaiID meng-unset semua alamat primary milik seorang
// pegawai (lintas tipe), karena hanya boleh ada SATU alamat primary per pegawai.
func (r *repository) UnsetPrimaryAlamatByPegawaiID(
	ctx context.Context,
	pegawaiID int64,
	updatedBy int64,
) error {
	return r.db.WithContext(ctx).
		Model(&models.KepegawaianAlamat{}).
		Where("pegawai_id = ? AND is_primary = true AND deleted_at IS NULL", pegawaiID).
		Updates(map[string]any{
			"is_primary": false,
			"updated_by": updatedBy,
			"updated_at": time.Now(),
		}).Error
}

// ── CheckDuplicate ───────────────────────────────────────────────────────────
// CheckDuplicateAlamat memeriksa apakah pegawai yang bersangkutan (m.PegawaiID)
// sudah pernah menginput alamat yang identik sebelumnya. Pengecekan discope
// per pegawai, sehingga pegawai lain tetap boleh memiliki alamat yang sama persis.
//
// excludeID diisi saat proses Update, agar record yang sedang diedit tidak
// dianggap duplikat terhadap dirinya sendiri.
func (r *repository) CheckDuplicateAlamat(ctx context.Context, m *models.KepegawaianAlamat, excludeID *int64) (bool, error) {
	query := r.db.WithContext(ctx).
		Model(&models.KepegawaianAlamat{}).
		Where("deleted_at IS NULL").
		Where("pegawai_id = ?", m.PegawaiID).
		Where("tipe_id = ?", m.TipeID).
		Where("LOWER(TRIM(jalan)) = LOWER(TRIM(?))", m.Jalan)

	// String nullable: RT, RW, KodePos
	query = whereNullableString(query, "rt", m.RT)
	query = whereNullableString(query, "rw", m.RW)
	query = whereNullableString(query, "kode_pos", m.KodePos)

	// Hirarki wilayah nullable: NegaraID, ProvinsiID, KotaKabupatenID, KecamatanID, KelurahanDesaID
	query = whereNullableID(query, "negara_id", m.NegaraID)
	query = whereNullableID(query, "provinsi_id", m.ProvinsiID)
	query = whereNullableID(query, "kota_kabupaten_id", m.KotaKabupatenID)
	query = whereNullableID(query, "kecamatan_id", m.KecamatanID)
	query = whereNullableID(query, "kelurahan_desa_id", m.KelurahanDesaID)

	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// whereNullableString menambahkan kondisi "kolom IS NULL" jika val nil,
// atau "kolom = val" (case-insensitive, trimmed) jika val terisi.
func whereNullableString(q *gorm.DB, column string, val *string) *gorm.DB {
	if val == nil {
		return q.Where(column + " IS NULL")
	}
	return q.Where("LOWER(TRIM("+column+")) = LOWER(TRIM(?))", *val)
}

// whereNullableID menambahkan kondisi "kolom IS NULL" jika val nil,
// atau "kolom = val" jika val terisi.
func whereNullableID(q *gorm.DB, column string, val *int64) *gorm.DB {
	if val == nil {
		return q.Where(column + " IS NULL")
	}
	return q.Where(column+" = ?", *val)
}
