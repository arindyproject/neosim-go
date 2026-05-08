package dto

import "neosim_go/internal/modules/users/models"

// ─── Request DTOs ──────────────────────────────────────────────────────────────

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=150"`
	Email    string `json:"email"    validate:"required,email"`
	Name     string `json:"name"     validate:"required,min=1,max=255"`
	Password string `json:"password" validate:"required,min=8"`
}

type UpdateUserRequest struct {
	Name        *string `json:"name"         validate:"omitempty,min=1,max=255"`
	Email       *string `json:"email"        validate:"omitempty,email"`
	Photo       *string `json:"photo"        validate:"omitempty,max=500"`
	IsActive    *bool   `json:"is_active"`
	IsStaff     *bool   `json:"is_staff"`
	IsSuperuser *bool   `json:"is_superuser"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// UpdateSettingsRequest request untuk update settings user
type UpdateSettingsRequest struct {
	Settings []models.UserSetting `json:"settings" validate:"required"`
}

// ─── Default Settings ──────────────────────────────────────────────────────────

// DefaultUserSettings mengembalikan default settings untuk user baru
// Dipanggil dari service saat CreateUser
func DefaultUserSettings() []models.UserSetting {
	return []models.UserSetting{
		{
			Key:         "is_dark_mode",
			Type:        "boolean",
			Value:       false,
			Description: "Aktifkan tema gelap pada antarmuka",
		},
	}
}
