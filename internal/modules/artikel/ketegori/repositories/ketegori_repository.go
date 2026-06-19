
package repositories

import (
	"errors"

	"neosim_go/internal/modules/artikel/ketegori/models"
	"neosim_go/internal/modules/artikel/ketegori/dto"

	"gorm.io/gorm"
)



func (r *repository) Create(m *models.ArtikelKetegori) error {
	return r.db.Create(m).Error
}

func (r *repository) GetByID(id int64) (*models.ArtikelKetegori, error) {
	var m models.ArtikelKetegori
	result := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, result.Error
}

func (r *repository) List(page, pageSize int, filter *dto.FilterArtikelKetegoriRequest) ([]models.ArtikelKetegori, int64, error) {
	var items []models.ArtikelKetegori
	var total int64

	//------------------------------------------------------------
	// 1. Inisialisasi basis query & pastikan record yang di-soft delete tidak ikut terbawa
	query := r.db.Model(&models.ArtikelKetegori{}).Where("deleted_at IS NULL")

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

func (r *repository) Update(m *models.ArtikelKetegori) error {
	return r.db.Save(m).Error
}

func (r *repository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.ArtikelKetegori{}).Error
}
