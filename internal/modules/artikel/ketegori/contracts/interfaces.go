package contracts

import (
	"neosim_go/internal/modules/artikel/ketegori/dto"
	"neosim_go/internal/modules/artikel/ketegori/models"
	he "neosim_go/internal/shared/httputil"
)



// Repository defines database operations
type Repository interface {
	Create(m *models.ArtikelKetegori) error
	GetByID(id int64) (*models.ArtikelKetegori, error)
	List(page, pageSize int, filter *dto.FilterArtikelKetegoriRequest) ([]models.ArtikelKetegori, int64, error)
	Update(m *models.ArtikelKetegori) error
	Delete(id int64) error
}

// Service defines business logic operations
type Service interface {
	Create(req *dto.CreateArtikelKetegoriRequest, createdBy *int64, actor he.AuthContext) (*dto.ArtikelKetegoriResponse, error)
	GetByID(id int64 , actor he.AuthContext) (*dto.ArtikelKetegoriResponse, error)
	List(page, pageSize int,filter *dto.FilterArtikelKetegoriRequest, actor he.AuthContext) ([]dto.ArtikelKetegoriResponse, int64, error)
	Update(id int64, req *dto.UpdateArtikelKetegoriRequest, updatedBy *int64, actor he.AuthContext) (*dto.ArtikelKetegoriResponse, error)
	Delete(id int64, actor he.AuthContext) error
}
