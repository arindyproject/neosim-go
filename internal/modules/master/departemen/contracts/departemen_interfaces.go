package contracts

import (
	"context"
	"neosim_go/internal/modules/master/departemen/dto"
	"neosim_go/internal/modules/master/departemen/models"
	he "neosim_go/internal/shared/httputil"
)

// MasterDepartemenRepository defines database operations for MasterDepartemen.
// Diimplementasikan oleh struct 'repository' (lihat repositories/repository.go
// & repositories/departemen_repository.go).
type MasterDepartemenRepository interface {
	CreateDepartemen(ctx context.Context,m *models.MasterDepartemen) error
	GetDepartemenByID(ctx context.Context,id int64) (*models.MasterDepartemen, error)
	ListDepartemen(ctx context.Context,page, pageSize int, filter *dto.FilterMasterDepartemenRequest) ([]models.MasterDepartemen, int64, error)
	UpdateDepartemen(ctx context.Context,m *models.MasterDepartemen) error
	DeleteDepartemen(ctx context.Context,id int64) error
}

// MasterDepartemenService defines business logic operations for MasterDepartemen.
// Diimplementasikan oleh struct 'service' (lihat services/service.go
// & services/departemen_service.go).
type MasterDepartemenService interface {
	CreateDepartemen(ctx context.Context,req *dto.CreateMasterDepartemenRequest, actor he.AuthContext) (*dto.MasterDepartemenResponse, error)
	GetDepartemenByID(ctx context.Context,id int64, actor he.AuthContext) (*dto.MasterDepartemenResponse, error)
	ListDepartemen(ctx context.Context,page, pageSize int, filter *dto.FilterMasterDepartemenRequest, actor he.AuthContext) ([]dto.MasterDepartemenResponse, int64, error)
	UpdateDepartemen(ctx context.Context,id int64, req *dto.UpdateMasterDepartemenRequest, actor he.AuthContext) (*dto.MasterDepartemenResponse, error)
	DeleteDepartemen(ctx context.Context,id int64, actor he.AuthContext) error
}
