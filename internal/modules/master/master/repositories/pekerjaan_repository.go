package repositories

import (
	"context"
	"errors"
	"time"

	"neosim_go/internal/modules/master/master/dto"
	"neosim_go/internal/modules/master/master/models"

	"gorm.io/gorm"
)

// =====================================================================
// Pekerjaan
// =====================================================================
// ------------------Create---------------------------------------------
func (r *repository) CreatePekerjaan(ctx context.Context, m *models.MasterPekerjaan) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// ------------------GetByID--------------------------------------------
func (r *repository) GetByIDPekerjaan(ctx context.Context, id int64) (*models.MasterPekerjaan, error) {
	var m models.MasterPekerjaan
	result := r.db.WithContext(ctx).Where("id = ?", id).
		Where("master_pekerjaan.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ------------------GetByName------------------------------------------
func (r *repository) GetByNamePekerjaan(ctx context.Context, name string) (*models.MasterPekerjaan, error) {
	var m models.MasterPekerjaan
	result := r.db.WithContext(ctx).Where("name = ?", name).Where("master_pekerjaan.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

// ------------------ListSelect-----------------------------------------
func (r *repository) ListSelectPekerjaan(ctx context.Context, search string) ([]models.MasterPekerjaan, error) {
	var items []models.MasterPekerjaan

	query := r.db.WithContext(ctx).Model(&models.MasterPekerjaan{}).
		Select("id, name, kode_kemenkes").
		Where("master_pekerjaan.deleted_at IS NULL")

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := query.Order("name ASC").Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// ------------------List-----------------------------------------------
func (r *repository) ListPekerjaan(ctx context.Context, page, pageSize int, filter *dto.FilterMasterPekerjaanRequest) ([]models.MasterPekerjaan, int64, error) {
	var items []models.MasterPekerjaan
	var total int64

	query := r.db.WithContext(ctx).Model(&models.MasterPekerjaan{}).Where("master_pekerjaan.deleted_at IS NULL")

	if filter != nil && filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if filter != nil && filter.KodeKemenkes != "" {
		query = query.Where("kode_kemenkes ILIKE ?", "%"+filter.KodeKemenkes+"%")
	}

	if filter != nil && filter.FhirCode != "" {
		query = query.Where("fhir_code ILIKE ?", "%"+filter.FhirCode+"%")
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

// ------------------Update---------------------------------------------
func (r *repository) UpdatePekerjaan(ctx context.Context, m *models.MasterPekerjaan) error {
	return r.db.WithContext(ctx).Save(m).Error
}

// ------------------Delete---------------------------------------------
func (r *repository) DeletePekerjaan(ctx context.Context, id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.MasterPekerjaan{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
		}).Error
}

// ===================================================================
