package repositories

import (
	"errors"

	"neosim_go/internal/modules/kepegawaian/kontak/dto"
	"neosim_go/internal/modules/kepegawaian/kontak/models"

	"gorm.io/gorm"
)

// ── Create ────────────────────────────────────────────────────────────────────
func (r *repository) CreateKontak(m *models.KepegawaianKontak) error {
	return r.db.Create(m).Error
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (r *repository) GetKontakByID(id int64) (*models.KepegawaianKontak, error) {
	var m models.KepegawaianKontak
	result := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── GetByPegawaiID ───────────────────────────────────────────────────────────────
func (r *repository) GetKontakByPegawaiID(pegawaiID int64) ([]models.KepegawaianKontak, error) {
	var items []models.KepegawaianKontak
	result := r.db.Where("pegawai_id = ? AND deleted_at IS NULL", pegawaiID).Find(&items)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return items, result.Error
}

// ── List ──────────────────────────────────────────────────────────────────────
func (r *repository) ListKontak(page, pageSize int, filter *dto.FilterKepegawaianKontakRequest) ([]models.KepegawaianKontak, int64, error) {
	var items []models.KepegawaianKontak
	var total int64

	query := r.db.Model(&models.KepegawaianKontak{}).Where("deleted_at IS NULL")

	if filter.PegawaiID != nil {
		query = query.Where("pegawai_id = ?", *filter.PegawaiID)
	}
	if filter.TipeID != nil {
		query = query.Where("tipe_id = ?", *filter.TipeID)
	}
	if filter.Nilai != nil {
		query = query.Where("nilai ILIKE ?", "%"+*filter.Nilai+"%")
	}
	if filter.IsPrimary != nil {
		query = query.Where("is_primary = ?", *filter.IsPrimary)
	}
	if filter.IsAktif != nil {
		query = query.Where("is_aktif = ?", *filter.IsAktif)
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

// ── Update ───────────────────────────────────────────────────────────────────
func (r *repository) UpdateKontak(m *models.KepegawaianKontak) error {
	return r.db.Save(m).Error
}

// ── Delete ───────────────────────────────────────────────────────────────────
func (r *repository) DeleteKontak(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.KepegawaianKontak{}).Error
}
