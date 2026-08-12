package contracts

import (
	"neosim_go/internal/modules/artikel/kategori/dto"
	"neosim_go/internal/modules/artikel/kategori/models"
	he "neosim_go/internal/shared/httputil"
)

// ArtikelKategoriRepository defines database operations for ArtikelKategori.
// Diimplementasikan oleh struct 'repository' (lihat repositories/repository.go
// & repositories/kategori_repository.go).
type ArtikelKategoriRepository interface {
	CreateKategori(m *models.ArtikelKategori) error
	GetKategoriByID(id int64) (*models.ArtikelKategori, error)
	ListKategori(page, pageSize int, filter *dto.FilterArtikelKategoriRequest) ([]models.ArtikelKategori, int64, error)
	UpdateKategori(m *models.ArtikelKategori) error
	DeleteKategori(id int64) error
}

// ArtikelKategoriService defines business logic operations for ArtikelKategori.
// Diimplementasikan oleh struct 'service' (lihat services/service.go
// & services/kategori_service.go).
type ArtikelKategoriService interface {
	CreateKategori(req *dto.CreateArtikelKategoriRequest, actor he.AuthContext) (*dto.ArtikelKategoriResponse, error)
	GetKategoriByID(id int64, actor he.AuthContext) (*dto.ArtikelKategoriResponse, error)
	ListKategori(page, pageSize int, filter *dto.FilterArtikelKategoriRequest, actor he.AuthContext) ([]dto.ArtikelKategoriResponse, int64, error)
	UpdateKategori(id int64, req *dto.UpdateArtikelKategoriRequest, actor he.AuthContext) (*dto.ArtikelKategoriResponse, error)
	DeleteKategori(id int64, actor he.AuthContext) error
}
