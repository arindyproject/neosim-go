package dto

import (
	"time"

	authModels "neosim_go/internal/modules/auth/models"
	rbacDto "neosim_go/internal/modules/rbac/dto"
	"neosim_go/internal/modules/users/models"
)

// ─── User Response (detail) ────────────────────────────────────────────────────

// UserResponse response lengkap untuk single user
// - roles      : tanpa permissions di dalamnya (flat)
// - permissions: semua permission dari role + direct, ditampilkan sebagai object lengkap
type UserResponse struct {
	ID                int64                        `json:"id"`
	Photo             *string                      `json:"photo"`
	PhotoThumbnail    *string                      `json:"photo_thumbnail"`
	Username          string                       `json:"username"`
	Email             string                       `json:"email"`
	Name              string                       `json:"name"`
	IsSuperadmin      bool                         `json:"is_superadmin"`
	IsActive          bool                         `json:"is_active"`
	IsStaff           bool                         `json:"is_staff"`
	IsVerified        bool                         `json:"is_verified"`
	PasswordChangedAt *time.Time                   `json:"password_changed_at"`
	LastLoginAt       *time.Time                   `json:"last_login_at"`
	Settings          []models.UserSetting         `json:"settings"`
	Roles             []rbacDto.RoleSimpleResponse `json:"roles"`       // ← tanpa permissions
	Permissions       []rbacDto.PermissionResponse `json:"permissions"` // ← object lengkap, deduplicated
	Histories         []authModels.LoginHistory    `json:"histories"`
	Creator           *models.UserCreator          `json:"creator"`
	CreatedAt         time.Time                    `json:"created_at"`
	UpdatedAt         time.Time                    `json:"updated_at"`
}

// ─── User Simple Response (list) ──────────────────────────────────────────────

// UserSimpleResponse response ringkas untuk list
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

// ─── Builders ──────────────────────────────────────────────────────────────────

// UserResponseParams parameter untuk ToUserResponse
type UserResponseParams struct {
	User        *models.User
	Roles       []rbacDto.RoleSimpleResponse // ← simple (tanpa permissions)
	Permissions []rbacDto.PermissionResponse // ← object lengkap
	Histories   []authModels.LoginHistory
	Creator     *models.UserCreator
}

// ToUserResponse mengubah User + RBAC data menjadi UserResponse lengkap
func ToUserResponse(p UserResponseParams) *UserResponse {
	settings, _ := p.User.GetSettings()

	roles := p.Roles
	if roles == nil {
		roles = []rbacDto.RoleSimpleResponse{}
	}

	permissions := p.Permissions
	if permissions == nil {
		permissions = []rbacDto.PermissionResponse{}
	}

	histories := p.Histories
	if histories == nil {
		histories = []authModels.LoginHistory{}
	}

	return &UserResponse{
		ID:                p.User.ID,
		Photo:             p.User.Photo,
		PhotoThumbnail:    p.User.PhotoThumbnail,
		Username:          p.User.Username,
		Email:             p.User.Email,
		Name:              p.User.Name,
		IsSuperadmin:      p.User.IsSuperadmin,
		IsActive:          p.User.IsActive,
		IsStaff:           p.User.IsStaff,
		IsVerified:        p.User.IsVerified,
		PasswordChangedAt: p.User.PasswordChangedAt,
		LastLoginAt:       p.User.LastLoginAt,
		Settings:          settings,
		Roles:             roles,
		Permissions:       permissions,
		Histories:         histories,
		Creator:           p.Creator,
		CreatedAt:         p.User.CreatedAt,
		UpdatedAt:         p.User.UpdatedAt,
	}
}

// ToUserSimpleResponse mengubah models.User menjadi UserSimpleResponse
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

// ToUserListResponse mengubah slice models.User menjadi slice UserSimpleResponse
func ToUserListResponse(users []models.User) []UserSimpleResponse {
	responses := make([]UserSimpleResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, *ToUserSimpleResponse(&u))
	}
	return responses
}
