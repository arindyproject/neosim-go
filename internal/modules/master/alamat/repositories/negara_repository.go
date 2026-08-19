package repositories

import (
	"context"
	"errors"

	"neosim_go/internal/modules/master/alamat/dto"
	"neosim_go/internal/modules/master/alamat/models"

	"gorm.io/gorm"
)

// =====================================================================
// NEGARA
// =====================================================================
func (r *repository) CreateNegara(ctx context.Context, m *models.MasterAlamatNegara) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *repository) GetByIDNegara(ctx context.Context, id int64) (*models.MasterAlamatNegara, error) {
	var m models.MasterAlamatNegara
	result := r.db.WithContext(ctx).Where("id = ?", id).
		Where("master_alamat_negara.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

func (r *repository) ListNegara(ctx context.Context, page, pageSize int, filter *dto.FilterNegaraRequest) ([]models.MasterAlamatNegara, int64, error) {
	var items []models.MasterAlamatNegara
	var total int64

	query := r.db.WithContext(ctx).Model(&models.MasterAlamatNegara{}).Where("master_alamat_negara.deleted_at IS NULL")

	if filter != nil && filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
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

func (r *repository) UpdateNegara(ctx context.Context, m *models.MasterAlamatNegara) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *repository) DeleteNegara(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.MasterAlamatNegara{}).Error
}

func (r *repository) ExistsNegaraByCode(ctx context.Context, code string, excludeID *int64) (bool, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&models.MasterAlamatNegara{}).Where("code = ?", code)
	if excludeID != nil {
		q = q.Where("id != ?", *excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}
