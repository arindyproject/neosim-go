package dto

import (
	"time"

	"neosim_go/internal/modules/rbac/models"
)

// ═══════════════════════════════════════════════════════════════
// RESPONSE DTOs
// ═══════════════════════════════════════════════════════════════

// ─── Permission ────────────────────────────────────────────────

type PermissionResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description *string   `json:"description"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToPermissionResponse(p *models.Permission) *PermissionResponse {
	return &PermissionResponse{
		ID:          p.ID,
		Name:        p.Name,
		DisplayName: p.DisplayName,
		Description: p.Description,
		Resource:    p.Resource,
		Action:      p.Action,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func ToPermissionListResponse(items []models.Permission) []PermissionResponse {
	var result []PermissionResponse
	for _, p := range items {
		result = append(result, *ToPermissionResponse(&p))
	}
	return result
}

// ─── Role ──────────────────────────────────────────────────────

type RoleResponse struct {
	ID          int64                `json:"id"`
	Name        string               `json:"name"`
	DisplayName string               `json:"display_name"`
	Description *string              `json:"description"`
	IsSystem    bool                 `json:"is_system"`
	Permissions []PermissionResponse `json:"permissions,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

func ToRoleResponse(r *models.Role) *RoleResponse {
	resp := &RoleResponse{
		ID:          r.ID,
		Name:        r.Name,
		DisplayName: r.DisplayName,
		Description: r.Description,
		IsSystem:    r.IsSystem,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
	if r.Permissions != nil {
		resp.Permissions = ToPermissionListResponse(r.Permissions)
	}
	return resp
}

func ToRoleListResponse(items []models.Role) []RoleResponse {
	var result []RoleResponse
	for _, r := range items {
		result = append(result, *ToRoleResponse(&r))
	}
	return result
}

// ─── Direct Permission ─────────────────────────────────────────

type DirectPermissionResponse struct {
	Permission PermissionResponse `json:"permission"`
	IsGranted  bool               `json:"is_granted"`
}

// ─── User RBAC Summary ─────────────────────────────────────────

type UserRBACResponse struct {
	UserID         int64                      `json:"user_id"`
	Roles          []RoleResponse             `json:"roles"`
	Permissions    []DirectPermissionResponse `json:"direct_permissions"`
	AllPermissions []string                   `json:"all_permissions"`
}
