package contracts

import (
	"neosim_go/internal/modules/kepegawaian/identifier/dto"
	"neosim_go/internal/modules/kepegawaian/identifier/models"
	he "neosim_go/internal/shared/httputil"
)

// KepegawaianIdentifierRepository defines database operations for KepegawaianIdentifier.
// Diimplementasikan oleh struct 'repository' (lihat repositories/repository.go
// & repositories/identifier_repository.go).
type KepegawaianIdentifierRepository interface {
	CreateIdentifier(m *models.KepegawaianIdentifier) error
	GetIdentifierByID(id int64) (*models.KepegawaianIdentifier, error)
	ListIdentifier(page, pageSize int, filter *dto.FilterKepegawaianIdentifierRequest) ([]models.KepegawaianIdentifier, int64, error)
	UpdateIdentifier(m *models.KepegawaianIdentifier) error
	DeleteIdentifier(id int64) error
}

// KepegawaianIdentifierService defines business logic operations for KepegawaianIdentifier.
// Diimplementasikan oleh struct 'service' (lihat services/service.go
// & services/identifier_service.go).
type KepegawaianIdentifierService interface {
	CreateIdentifier(req *dto.CreateKepegawaianIdentifierRequest, actor he.AuthContext) (*dto.KepegawaianIdentifierResponse, error)
	GetIdentifierByID(id int64, actor he.AuthContext) (*dto.KepegawaianIdentifierResponse, error)
	ListIdentifier(page, pageSize int, filter *dto.FilterKepegawaianIdentifierRequest, actor he.AuthContext) ([]dto.KepegawaianIdentifierResponse, int64, error)
	UpdateIdentifier(id int64, req *dto.UpdateKepegawaianIdentifierRequest, actor he.AuthContext) (*dto.KepegawaianIdentifierResponse, error)
	DeleteIdentifier(id int64, actor he.AuthContext) error
}
