package contracts

import (
	"neosim_go/internal/modules/users/dto"
	"neosim_go/internal/modules/users/models"
)

// Service interface defines the business logic operations
type Service interface {
	// CRUD
	CreateUser(req *dto.CreateUserRequest, createdBy *int64) (*dto.UserResponse, error)
	GetUserByID(id int64) (*dto.UserResponse, error)
	GetUserByUsername(username string) (*dto.UserResponse, error)
	GetUserByEmail(email string) (*dto.UserResponse, error)
	ListUsers(page, pageSize int) ([]dto.UserResponse, int64, error)
	UpdateUser(id int64, req *dto.UpdateUserRequest, updatedBy *int64) (*dto.UserResponse, error)
	DeleteUser(id int64) error

	// Password
	ChangePassword(id int64, req *dto.ChangePasswordRequest) error
	UpdateLastLogin(id int64) error

	// Settings
	GetSettings(id int64) ([]models.UserSetting, error)
	UpdateSettings(id int64, req *dto.UpdateSettingsRequest) error
}

// Repository interface defines the database operations
type Repository interface {
	Create(user *models.User) error
	GetByID(id int64) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	List(page, pageSize int) ([]models.User, int64, error)
	Update(user *models.User) error
	Delete(id int64) error
	GetSettings(id int64) ([]models.UserSetting, error)
	UpdateSettings(id int64, settings []models.UserSetting) error
}
