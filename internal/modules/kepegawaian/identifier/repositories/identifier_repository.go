package repositories

import (
	"errors"

	"neosim_go/internal/modules/kepegawaian/identifier/dto"
	"neosim_go/internal/modules/kepegawaian/identifier/models"

	"gorm.io/gorm"
)

func (r *repository) CreateIdentifier(m *models.KepegawaianIdentifier) error {
	return r.db.Create(m).Error
}

func (r *repository) GetIdentifierByID(id int64) (*models.KepegawaianIdentifier, error) {
	var m models.KepegawaianIdentifier
	result := r.db.
		Preload("Tipe").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&m)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

func (r *repository) ListIdentifier(page, pageSize int, filter *dto.FilterKepegawaianIdentifierRequest) ([]models.KepegawaianIdentifier, int64, error) {
	var items []models.KepegawaianIdentifier
	var total int64

	query := r.db.Model(&models.KepegawaianIdentifier{}).
		Preload("Tipe").
		Where("deleted_at IS NULL")

	if filter != nil {
		if filter.PegawaiID != nil {
			query = query.Where("pegawai_id = ?", *filter.PegawaiID)
		}
		if filter.TipeID != nil {
			query = query.Where("tipe_id = ?", *filter.TipeID)
		}
		if filter.Nilai != "" {
			query = query.Where("nilai ILIKE ?", "%"+filter.Nilai+"%")
		}
		if filter.IsPrimary != nil {
			query = query.Where("is_primary = ?", *filter.IsPrimary)
		}
		if filter.IsAktif != nil {
			query = query.Where("is_aktif = ?", *filter.IsAktif)
		}
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

func (r *repository) UpdateIdentifier(m *models.KepegawaianIdentifier) error {
	return r.db.Save(m).Error
}

func (r *repository) DeleteIdentifier(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.KepegawaianIdentifier{}).Error
}
