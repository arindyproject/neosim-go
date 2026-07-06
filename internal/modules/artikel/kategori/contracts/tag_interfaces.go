package contracts

import (
	"neosim_go/internal/modules/artikel/kategori/dto"
	"neosim_go/internal/modules/artikel/kategori/models"
	he "neosim_go/internal/shared/httputil"
)

// TagRepository defines database operations for Tag
type TagRepository interface {
	Create(m *models.Tag) error
	GetByID(id int64) (*models.Tag, error)
	List(page, pageSize int, filter *dto.FilterTagRequest) ([]models.Tag, int64, error)
	Update(m *models.Tag) error
	Delete(id int64) error
}

// TagService defines business logic operations for Tag
type TagService interface {
	Create(req *dto.CreateTagRequest, createdBy *int64, actor he.AuthContext) (*dto.TagResponse, error)
	GetByID(id int64, actor he.AuthContext) (*dto.TagResponse, error)
	List(page, pageSize int, filter *dto.FilterTagRequest, actor he.AuthContext) ([]dto.TagResponse, int64, error)
	Update(id int64, req *dto.UpdateTagRequest, updatedBy *int64, actor he.AuthContext) (*dto.TagResponse, error)
	Delete(id int64, actor he.AuthContext) error
}
