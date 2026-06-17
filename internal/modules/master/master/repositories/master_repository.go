package repositories

import (
	"errors"

	"neosim_go/internal/modules/master/master/contracts"
	"neosim_go/internal/modules/master/master/models"
	"neosim_go/internal/modules/master/master/dto"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewMasterRepository membuat instance repository baru
func NewMasterRepository(db *gorm.DB) contracts.Repository {
	return &repository{db: db}
}

func (r *repository) Create(m *models.Master) error {
	return r.db.Create(m).Error
}

func (r *repository) GetByID(id int64) (*models.Master, error) {
	var m models.Master
	result := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

func (r *repository) List(page, pageSize int, filter *dto.FilterMasterRequest) ([]models.Master, int64, error) {
	var items []models.Master
	var total int64

	//------------------------------------------------------------
	// 1. Inisialisasi basis query & pastikan record yang di-soft delete tidak ikut terbawa
	query := r.db.Model(&models.Master{}).Where("deleted_at IS NULL")

	// 2. Filter Teks (Menggunakan ILIKE untuk case-insensitive)
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	
	// 3. Hitung total data berdasarkan filter yang aktif
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 5. Ambil data dengan paginasi dan sorting
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *repository) Update(m *models.Master) error {
	return r.db.Save(m).Error
}

func (r *repository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.Master{}).Error
}
