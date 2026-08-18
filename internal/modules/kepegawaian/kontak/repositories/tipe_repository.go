package repositories

import (
	"errors"

	"neosim_go/internal/modules/kepegawaian/kontak/contracts"
	"neosim_go/internal/modules/kepegawaian/kontak/dto"
	"neosim_go/internal/modules/kepegawaian/kontak/models"

	"gorm.io/gorm"
)

// NewTipeRepository mengembalikan struct repository yang SAMA
// dengan repository entitas utama, dilihat sebagai contracts.TipeRepository.
// Berguna untuk test Tipe yang berdiri sendiri; di production cukup
// pakai repo yang sudah dibuat lewat NewKepegawaianKontakRepository(db).
func NewTipeRepository(db *gorm.DB) contracts.TipeRepository {
	return &repository{db: db}
}

// ── Create ────────────────────────────────────────────────────────────────────
func (r *repository) CreateTipe(m *models.Tipe) error {
	return r.db.Create(m).Error
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (r *repository) GetTipeByID(id int64) (*models.Tipe, error) {
	var m models.Tipe
	result := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── GetByCode ─────────────────────────────────────────────────────────────────
func (r *repository) GetTipeByCode(code string) (*models.Tipe, error) {
	var m models.Tipe
	result := r.db.Where("code = ? AND deleted_at IS NULL", code).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── GetByLabel ────────────────────────────────────────────────────────────────
func (r *repository) GetTipeByLabel(label string) (*models.Tipe, error) {
	var m models.Tipe
	result := r.db.Where("label = ? AND deleted_at IS NULL", label).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ── List ──────────────────────────────────────────────────────────────────────
func (r *repository) ListTipe(page, pageSize int, filter *dto.FilterTipeRequest) ([]models.Tipe, int64, error) {
	var items []models.Tipe
	var total int64

	query := r.db.Model(&models.Tipe{}).Where("deleted_at IS NULL")
	if filter.Code != "" {
		query = query.Where("code ILIKE ?", "%"+filter.Code+"%")
	}
	if filter.Label != "" {
		query = query.Where("label ILIKE ?", "%"+filter.Label+"%")
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
func (r *repository) UpdateTipe(m *models.Tipe) error {
	return r.db.Save(m).Error
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (r *repository) DeleteTipe(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.Tipe{}).Error
}
