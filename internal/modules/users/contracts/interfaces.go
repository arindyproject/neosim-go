package contracts

import (
	"context"
	"io"
	"neosim_go/internal/modules/users/dto"
	"neosim_go/internal/modules/users/models"
	he "neosim_go/internal/shared/httputil"
)

// ─── Repository ────────────────────────────────────────────────────────────────

type Repository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByIDs(ctx context.Context, ids []int64) ([]models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	List(ctx context.Context, page, pageSize int, filter *dto.UserFilter) ([]models.User, int64, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int64, deletedBy int64, reason string) error
	DeletedList(ctx context.Context, page, pageSize int, filter *dto.UserDeletedFilter) ([]models.User, int64, error)
	GetSettings(ctx context.Context, id int64) ([]models.UserSetting, error)
	UpdateSettings(ctx context.Context, id int64, settings []models.UserSetting) error
}

// ─── Service ───────────────────────────────────────────────────────────────────

type Service interface {
	// CRUD — operasi yang butuh auth context
	CreateUser(ctx context.Context, req *dto.CreateUserRequest, actor he.AuthContext) (*dto.UserSimpleResponse, error)
	GetUserByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.UserResponse, error)
	GetUserByUsername(ctx context.Context, username string, actor he.AuthContext) (*dto.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string, actor he.AuthContext) (*dto.UserResponse, error)
	ListUsers(ctx context.Context, page, pageSize int, filter *dto.UserFilter) ([]dto.UserSimpleResponse, int64, error)
	ListDeletedUsers(ctx context.Context, page, pageSize int, filter *dto.UserDeletedFilter, actor he.AuthContext) ([]dto.UserDeletedResponse, int64, error)
	UpdateUser(ctx context.Context, id int64, req *dto.UpdateUserRequest, actor he.AuthContext) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id int64, reason string, actor he.AuthContext) error

	// Password
	ChangePassword(ctx context.Context, id int64, req *dto.ChangePasswordRequest, actor he.AuthContext) (*dto.UserResponse, error)
	ResetPassword(ctx context.Context, id int64, actor he.AuthContext) error
	UpdateLastLogin(ctx context.Context, id int64) error

	// Settings
	GetSettings(ctx context.Context, id int64, actor he.AuthContext) ([]models.UserSetting, error)
	UpdateSettings(ctx context.Context, id int64, req *dto.UpdateSettingsRequest, actor he.AuthContext) (*dto.UserResponse, error)

	//Photo
	UploadPhoto(ctx context.Context, id int64, filename string, reader io.Reader, actor he.AuthContext) (*dto.UserResponse, error)
	DeletePhoto(ctx context.Context, id int64, actor he.AuthContext) (*dto.UserResponse, error)
}
