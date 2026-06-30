package contracts

import (
	"neosim_go/internal/modules/artikel/artikel/dto"
	"neosim_go/internal/modules/artikel/artikel/models"
	he "neosim_go/internal/shared/httputil"
)

// Repository defines database operations
type Repository interface {
	Create(m *models.Artikel) error
	GetByID(id int64) (*models.Artikel, error)
	List(page, pageSize int, filter *dto.FilterArtikelRequest) ([]models.Artikel, int64, error)
	Update(m *models.Artikel) error
	Delete(id int64) error
}

// Service defines business logic operations
type Service interface {
	Create(req *dto.CreateArtikelRequest, actor he.AuthContext) (*dto.ArtikelResponse, error)
	GetByID(id int64, actor he.AuthContext) (*dto.ArtikelResponse, error)
	List(page, pageSize int, filter *dto.FilterArtikelRequest, actor he.AuthContext) ([]dto.ArtikelResponse, int64, error)
	Update(id int64, req *dto.UpdateArtikelRequest, actor he.AuthContext) (*dto.ArtikelResponse, error)
	Delete(id int64, actor he.AuthContext) error
}
