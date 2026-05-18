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
	IsSuperuser       bool                      `json:"is_superuser"`
	IsActive          bool                      `json:"is_active"`
	IsStaff           bool                      `json:"is_staff"`
	IsVerified        bool                      `json:"is_verified"`
	PasswordChangedAt *time.Time                `json:"password_changed_at"`
	LastLoginAt       *time.Time                `json:"last_login_at"`
	Settings          []models.UserSetting      `json:"settings"`
	Histories         []authModels.LoginHistory `json:"histories"`
	CreatedAt         time.Time                 `json:"created_at"`
	UpdatedAt         time.Time                 `json:"updated_at"`
}

// ToUserResponse mengubah models.User menjadi UserResponse
func ToUserResponse(u *models.User) *UserResponse {
	settings, _ := u.GetSettings()
	return &UserResponse{
		ID:                u.ID,
		Photo:             u.Photo,
		PhotoThumbnail:    u.PhotoThumbnail,
		Username:          u.Username,
		Email:             u.Email,
		Name:              u.Name,
		IsSuperuser:       u.IsSuperuser,
		IsActive:          u.IsActive,
		IsStaff:           u.IsStaff,
		IsVerified:        u.IsVerified,
		PasswordChangedAt: u.PasswordChangedAt,
		LastLoginAt:       u.LastLoginAt,
		Settings:          settings,
		CreatedAt:         u.CreatedAt,
		UpdatedAt:         u.UpdatedAt,
	}
}

// ToUserListResponse mengubah slice models.User menjadi slice UserResponse
func ToUserListResponse(users []models.User) []UserResponse {
	var responses []UserResponse
	for _, u := range users {
		responses = append(responses, *ToUserResponse(&u))
	}
	return responses
}
