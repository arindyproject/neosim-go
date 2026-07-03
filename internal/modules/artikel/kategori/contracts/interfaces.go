package contracts

import (
	"neosim_go/internal/modules/artikel/kategori/dto"
	"neosim_go/internal/modules/artikel/kategori/models"
	he "neosim_go/internal/shared/httputil"
)

// Repository defines database operations
type Repository interface {
	Create(m *models.ArtikelKategori) error
	GetByID(id int64) (*models.ArtikelKategori, error)
	List(page, pageSize int, filter *dto.FilterArtikelKategoriRequest) ([]models.ArtikelKategori, int64, error)
	Update(m *models.ArtikelKategori) error
	Delete(id int64) error
}

// Service defines business logic operations
type Service interface {
	Create(req *dto.CreateArtikelKategoriRequest, actor he.AuthContext) (*dto.ArtikelKategoriResponse, error)
	GetByID(id int64, actor he.AuthContext) (*dto.ArtikelKategoriResponse, error)
	List(page, pageSize int, filter *dto.FilterArtikelKategoriRequest, actor he.AuthContext) ([]dto.ArtikelKategoriResponse, int64, error)
	Update(id int64, req *dto.UpdateArtikelKategoriRequest, actor he.AuthContext) (*dto.ArtikelKategoriResponse, error)
	Delete(id int64, actor he.AuthContext) error
}
