package contracts

import (
	"context"
	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	he "neosim_go/internal/shared/httputil"
)

// PositionRepository defines database operations for Position.
// Diimplementasikan oleh struct 'repository' yang sama dengan entitas utama
// sub-module ini (lihat repositories/repository.go) — TIDAK ADA struct baru.
// Method diberi suffix nama item agar tidak bentrok saat di-embed ke
// contracts.Repository.
type PositionRepository interface {
	// Create
	CreatePosition(ctx context.Context, m *models.Position) error

	// Update
	UpdatePosition(ctx context.Context, m *models.Position) error

	// Delete (soft delete)
	DeletePosition(ctx context.Context, id int64, deletedBy int64) error

	// Find
	GetPositionByID(ctx context.Context, id int64) (*models.Position, error)
	GetPositionWithChildrenByID(ctx context.Context, id int64) (*models.Position, error)
	ListPosition(ctx context.Context, page, pageSize int, filter *dto.FilterPositionRequest) ([]models.Position, int64, error)
	FindChildrenByParentID(ctx context.Context, parentID int64) ([]models.Position, error)
	FindRootPositions(ctx context.Context) ([]models.Position, error)

	// FindAllPositions mengambil seluruh Position sekaligus (flat, dengan
	// PositionKategori ter-preload) untuk dirakit jadi tree di service layer
	// dalam satu pass O(n) — lebih murah daripada rekursi N+1 query.
	FindAllPositions(ctx context.Context, onlyAktif bool) ([]models.Position, error)

	// Exists
	ExistsPositionByID(ctx context.Context, id int64) (bool, error)
	ExistsByNameAndKategori(ctx context.Context, kategoriID int64, name string, excludeID int64) (bool, error)

	// Helper
	HasChildren(ctx context.Context, id int64) (bool, error)
}

// PositionService defines business logic operations for Position.
// Diimplementasikan oleh struct 'service' yang sama dengan entitas utama.
type PositionService interface {
	CreatePosition(ctx context.Context, req *dto.CreatePositionRequest, actor he.AuthContext) (*dto.PositionResponse, error)
	GetPositionByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.PositionResponse, error)
	GetPositionTree(ctx context.Context, onlyAktif bool, actor he.AuthContext) ([]dto.PositionTreeNode, error)
	ListPosition(ctx context.Context, page, pageSize int, filter *dto.FilterPositionRequest, actor he.AuthContext) ([]dto.PositionResponse, int64, error)
	UpdatePosition(ctx context.Context, id int64, req *dto.UpdatePositionRequest, actor he.AuthContext) (*dto.PositionResponse, error)
	DeletePosition(ctx context.Context, id int64, actor he.AuthContext) error
}
