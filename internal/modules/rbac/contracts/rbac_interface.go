package contracts

import (
	"context"
	"neosim_go/internal/modules/rbac/dto"
	"neosim_go/internal/modules/rbac/models"
)

// ─── Repository ────────────────────────────────────────────────────────────────

type RBACRepository interface {
	IsSuperadmin(ctx context.Context, userID int64) (bool, error)
	// Permission
	CreatePermission(ctx context.Context, p *models.Permission) error
	GetPermissionByID(ctx context.Context, id int64) (*models.Permission, error)
	GetPermissionByName(ctx context.Context, name string) (*models.Permission, error)
	ListPermissions(ctx context.Context, page, pageSize int) ([]models.Permission, int64, error)
	UpdatePermission(ctx context.Context, p *models.Permission) error
	DeletePermission(ctx context.Context, id int64) error

	// Role
	CreateRole(ctx context.Context, r *models.Role) error
	GetRoleByID(ctx context.Context, id int64) (*models.Role, error)
	GetRoleByName(ctx context.Context, name string) (*models.Role, error)
	ListRoles(ctx context.Context, page, pageSize int) ([]models.Role, int64, error)
	UpdateRole(ctx context.Context, r *models.Role) error
	DeleteRole(ctx context.Context, id int64) error
	GetUsersRoles(ctx context.Context, userIDs []int64) (map[int64][]models.Role, error)

	// Role ↔ Permission
	AssignPermissionsToRole(ctx context.Context, roleID int64, permissionIDs []int64) error
	RevokePermissionsFromRole(ctx context.Context, roleID int64, permissionIDs []int64) error
	GetRolePermissions(ctx context.Context, roleID int64) ([]models.Permission, error)
	SyncRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error

	// User ↔ Role
	AssignRolesToUser(ctx context.Context, userID int64, roleIDs []int64, assignedBy *int64) error
	RevokeRolesFromUser(ctx context.Context, userID int64, roleIDs []int64) error
	GetUserRoles(ctx context.Context, userID int64) ([]models.Role, error)
	SyncUserRoles(ctx context.Context, userID int64, roleIDs []int64, assignedBy *int64) error

	// User ↔ Permission (direct)
	AssignDirectPermission(ctx context.Context, userID, permissionID int64, isGranted bool, assignedBy *int64) error
	RevokeDirectPermission(ctx context.Context, userID, permissionID int64) error
	GetUserDirectPermissions(ctx context.Context, userID int64) ([]models.UserPermission, error)

	// Check
	GetUserAllPermissions(ctx context.Context, userID int64) ([]string, error) // gabungan dari role + direct
	HasPermission(ctx context.Context, userID int64, permission string) (bool, error)
}

// ─── Service ───────────────────────────────────────────────────────────────────

type RBACService interface {
	// Permission CRUD
	CreatePermission(ctx context.Context, req *dto.CreatePermissionRequest, createdBy *int64) (*dto.PermissionResponse, error)
	GetPermissionByID(ctx context.Context, id int64) (*dto.PermissionResponse, error)
	ListPermissions(ctx context.Context, page, pageSize int) ([]dto.PermissionResponse, int64, error)
	UpdatePermission(ctx context.Context, id int64, req *dto.UpdatePermissionRequest, updatedBy *int64) (*dto.PermissionResponse, error)
	DeletePermission(ctx context.Context, id int64) error

	// Role CRUD
	CreateRole(ctx context.Context, req *dto.CreateRoleRequest, createdBy *int64) (*dto.RoleResponse, error)
	GetRoleByID(ctx context.Context, id int64) (*dto.RoleResponse, error)
	ListRoles(ctx context.Context, page, pageSize int) ([]dto.RoleResponse, int64, error)
	UpdateRole(ctx context.Context, id int64, req *dto.UpdateRoleRequest, updatedBy *int64) (*dto.RoleResponse, error)
	DeleteRole(ctx context.Context, id int64) error

	// Role ↔ Permission
	AssignPermissionsToRole(ctx context.Context, roleID int64, req *dto.AssignPermissionsRequest) error
	RevokePermissionsFromRole(ctx context.Context, roleID int64, req *dto.AssignPermissionsRequest) error
	SyncRolePermissions(ctx context.Context, roleID int64, req *dto.AssignPermissionsRequest) error

	// User ↔ Role
	AssignRolesToUser(ctx context.Context, userID int64, req *dto.AssignRolesRequest, assignedBy *int64) error
	RevokeRolesFromUser(ctx context.Context, userID int64, req *dto.AssignRolesRequest) error
	SyncUserRoles(ctx context.Context, userID int64, req *dto.AssignRolesRequest, assignedBy *int64) error
	GetUserRoles(ctx context.Context, userID int64) ([]dto.RoleResponse, error)

	// User ↔ Permission (direct)
	AssignDirectPermission(ctx context.Context, userID int64, req *dto.AssignDirectPermissionRequest, assignedBy *int64) error
	RevokeDirectPermission(ctx context.Context, userID, permissionID int64) error
	GetUserDirectPermissions(ctx context.Context, userID int64) ([]dto.DirectPermissionResponse, error)

	// Check
	GetUserAllPermissions(ctx context.Context, userID int64) ([]string, error)
	HasPermission(ctx context.Context, userID int64, permission string) (bool, error)
}
