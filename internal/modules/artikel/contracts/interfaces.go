package contracts

import (
	"neosim_go/internal/modules/artikel/dto"
	"neosim_go/internal/modules/artikel/models"
)

// AuthContext berisi informasi user yang sedang login untuk authorization
type AuthContext struct {
	UserID       int64
	IsSuperadmin bool
}

// Repository defines database operations
type Repository interface {
	Create(m *models.Artikel) error
	GetByID(id int64) (*models.Artikel, error)
	List(page, pageSize int) ([]models.Artikel, int64, error)
	Update(m *models.Artikel) error
	Delete(id int64) error
}

// Service defines business logic operations
type Service interface {
	Create(req *dto.CreateArtikelRequest, createdBy *int64, actor AuthContext) (*dto.ArtikelResponse, error)
	GetByID(id int64 , actor AuthContext) (*dto.ArtikelResponse, error)
	List(page, pageSize int, actor AuthContext) ([]dto.ArtikelResponse, int64, error)
	Update(id int64, req *dto.UpdateArtikelRequest, updatedBy *int64, actor AuthContext) (*dto.ArtikelResponse, error)
	Delete(id int64, actor AuthContext) error
}
