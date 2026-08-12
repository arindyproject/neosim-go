package contracts

import (
	"neosim_go/internal/modules/artikel/artikel/dto"
	"neosim_go/internal/modules/artikel/artikel/models"
	he "neosim_go/internal/shared/httputil"
)

// ArtikelRepository defines database operations for Artikel.
// Diimplementasikan oleh struct 'repository' (lihat repositories/repository.go
// & repositories/artikel_repository.go).
type ArtikelRepository interface {
	CreateArtikel(m *models.Artikel) error
	GetArtikelByID(id int64) (*models.Artikel, error)
	ListArtikel(page, pageSize int, filter *dto.FilterArtikelRequest) ([]models.Artikel, int64, error)
	UpdateArtikel(m *models.Artikel) error
	DeleteArtikel(id int64) error
}

// ArtikelService defines business logic operations for Artikel.
// Diimplementasikan oleh struct 'service' (lihat services/service.go
// & services/artikel_service.go).
type ArtikelService interface {
	CreateArtikel(req *dto.CreateArtikelRequest, actor he.AuthContext) (*dto.ArtikelResponse, error)
	GetArtikelByID(id int64, actor he.AuthContext) (*dto.ArtikelResponse, error)
	ListArtikel(page, pageSize int, filter *dto.FilterArtikelRequest, actor he.AuthContext) ([]dto.ArtikelResponse, int64, error)
	UpdateArtikel(id int64, req *dto.UpdateArtikelRequest, actor he.AuthContext) (*dto.ArtikelResponse, error)
	DeleteArtikel(id int64, actor he.AuthContext) error
}
