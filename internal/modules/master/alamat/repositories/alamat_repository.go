package repositories

import (
	"context"
	"errors"
	"time"

	"neosim_go/internal/modules/master/alamat/dto"
	"neosim_go/internal/modules/master/alamat/models"

	"gorm.io/gorm"
)

// =====================================================================
// Desa / Kelurahan
// =====================================================================
func (r *repository) CreateKelurahanDesa(ctx context.Context, m *models.MasterAlamatKelurahanDesa) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *repository) GetByIDKelurahanDesa(ctx context.Context, id int64) (*models.MasterAlamatKelurahanDesa, error) {
	var m models.MasterAlamatKelurahanDesa
	result := r.db.WithContext(ctx).Preload("Kecamatan").Where("id = ?", id).
		Where("master_alamat_kelurahan_desa.deleted_at IS NULL").First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

func (r *repository) ListKelurahanDesa(ctx context.Context, page, pageSize int, kecamatanID *int64, filter *dto.FilterKelurahanDesaRequest) ([]models.MasterAlamatKelurahanDesa, int64, error) {
	var items []models.MasterAlamatKelurahanDesa
	var total int64

	query := r.db.WithContext(ctx).Model(&models.MasterAlamatKelurahanDesa{}).
		Where("master_alamat_kelurahan_desa.deleted_at IS NULL")

	if kecamatanID != nil {
		query = query.Where("kecamatan_id = ?", *kecamatanID)
	}
	if filter != nil && filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if filter != nil && filter.FhirCode != nil {
		query = query.Where("fhir_code ILIKE ?", "%"+*filter.FhirCode+"%")
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

func (r *repository) UpdateKelurahanDesa(ctx context.Context, m *models.MasterAlamatKelurahanDesa) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *repository) DeleteKelurahanDesa(ctx context.Context, id int64, deletedBy int64) error {
	return r.db.WithContext(ctx).
		Model(&models.MasterAlamatKelurahanDesa{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_by": deletedBy,
		}).Error
}

func (r *repository) ExistsKelurahanDesaByCode(ctx context.Context, code string, excludeID *int64) (bool, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&models.MasterAlamatKelurahanDesa{}).Where("code = ?", code)
	if excludeID != nil {
		q = q.Where("id != ?", *excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}
