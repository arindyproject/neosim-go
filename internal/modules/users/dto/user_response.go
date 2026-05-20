package dto

import (
	"time"

	authModels "neosim_go/internal/modules/auth/models"
	"neosim_go/internal/modules/users/models"
)

// UserResponse response untuk single user
type UserResponse struct {
	ID                int64                     `json:"id"`
	Photo             *string                   `json:"photo"`
	PhotoThumbnail    *string                   `json:"photo_thumbnail"`
	Username          string                    `json:"username"`
	Email             string                    `json:"email"`
	Name              string                    `json:"name"`
	IsSuperadmin      bool                      `json:"is_superadmin"`
	IsActive          bool                      `json:"is_active"`
	IsStaff           bool                      `json:"is_staff"`
	IsVerified        bool                      `json:"is_verified"`
	PasswordChangedAt *time.Time                `json:"password_changed_at"`
	LastLoginAt       *time.Time                `json:"last_login_at"`
	Settings          []models.UserSetting      `json:"settings"`
	Histories         []authModels.LoginHistory `json:"histories"`
	Creator           *models.UserCreator       `json:"creator"`
	CreatedAt         time.Time                 `json:"created_at"`
	UpdatedAt         time.Time                 `json:"updated_at"`
}

type UserSimpleResponse struct {
	ID             int64     `json:"id"`
	PhotoThumbnail *string   `json:"photo_thumbnail"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	IsSuperadmin   bool      `json:"is_superadmin"`
	IsActive       bool      `json:"is_active"`
	IsStaff        bool      `json:"is_staff"`
	IsVerified     bool      `json:"is_verified"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ToUserResponse mengubah models.User menjadi UserResponse
func ToUserResponse(u *models.User, histories []authModels.LoginHistory, creator *models.UserCreator) *UserResponse {

	settings, _ := u.GetSettings()
	return &UserResponse{
		ID:                u.ID,
		Photo:             u.Photo,
		PhotoThumbnail:    u.PhotoThumbnail,
		Username:          u.Username,
		Email:             u.Email,
		Name:              u.Name,
		IsSuperadmin:      u.IsSuperadmin,
		IsActive:          u.IsActive,
		IsStaff:           u.IsStaff,
		IsVerified:        u.IsVerified,
		PasswordChangedAt: u.PasswordChangedAt,
		LastLoginAt:       u.LastLoginAt,
		Settings:          settings,
		Creator:           creator,
		CreatedAt:         u.CreatedAt,
		UpdatedAt:         u.UpdatedAt,
	}
}

func ToUserSimpleResponse(u *models.User) *UserSimpleResponse {
	return &UserSimpleResponse{
		ID:             u.ID,
		PhotoThumbnail: u.PhotoThumbnail,
		Username:       u.Username,
		Email:          u.Email,
		Name:           u.Name,
		IsSuperadmin:   u.IsSuperadmin,
		IsActive:       u.IsActive,
		IsStaff:        u.IsStaff,
		IsVerified:     u.IsVerified,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}

// ToUserListResponse mengubah slice models.User menjadi slice UserResponse
func ToUserListResponse(users []models.User) []UserSimpleResponse {
	var responses []UserSimpleResponse
	for _, u := range users {
		responses = append(responses, *ToUserSimpleResponse(&u))
	}
	return responses
}
