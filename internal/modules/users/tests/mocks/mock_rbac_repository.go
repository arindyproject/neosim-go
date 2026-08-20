package mocks

import (
	"context"
	"neosim_go/internal/modules/rbac/models"

	"github.com/stretchr/testify/mock"
)

// MockRBACRepository adalah mock untuk rbacContracts.RBACRepository
type MockRBACRepository struct {
	mock.Mock
}

func (m *MockRBACRepository) IsSuperadmin(ctx context.Context, userID int64) (bool, error) {
	args := m.Called(userID)
	return args.Bool(0), args.Error(1)
}

// ─── Permission ────────────────────────────────────────────────────────────────

func (m *MockRBACRepository) CreatePermission(ctx context.Context, p *models.Permission) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockRBACRepository) GetPermissionByID(ctx context.Context, id int64) (*models.Permission, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Permission), args.Error(1)
}

func (m *MockRBACRepository) GetPermissionByName(ctx context.Context, name string) (*models.Permission, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Permission), args.Error(1)
}

func (m *MockRBACRepository) ListPermissions(ctx context.Context, page, pageSize int) ([]models.Permission, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]models.Permission), args.Get(1).(int64), args.Error(2)
}

func (m *MockRBACRepository) UpdatePermission(ctx context.Context, p *models.Permission) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockRBACRepository) DeletePermission(ctx context.Context, id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

// ─── Role ──────────────────────────────────────────────────────────────────────

func (m *MockRBACRepository) CreateRole(ctx context.Context, r *models.Role) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *MockRBACRepository) GetRoleByID(ctx context.Context, id int64) (*models.Role, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRBACRepository) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRBACRepository) ListRoles(ctx context.Context, page, pageSize int) ([]models.Role, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]models.Role), args.Get(1).(int64), args.Error(2)
}

func (m *MockRBACRepository) UpdateRole(ctx context.Context, r *models.Role) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *MockRBACRepository) DeleteRole(ctx context.Context, id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

// ─── Role ↔ Permission ─────────────────────────────────────────────────────────

func (m *MockRBACRepository) AssignPermissionsToRole(ctx context.Context, roleID int64, permissionIDs []int64) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

func (m *MockRBACRepository) RevokePermissionsFromRole(ctx context.Context, roleID int64, permissionIDs []int64) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

func (m *MockRBACRepository) GetRolePermissions(ctx context.Context, roleID int64) ([]models.Permission, error) {
	args := m.Called(roleID)
	return args.Get(0).([]models.Permission), args.Error(1)
}

func (m *MockRBACRepository) SyncRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

// ─── User ↔ Role ───────────────────────────────────────────────────────────────

func (m *MockRBACRepository) AssignRolesToUser(ctx context.Context, userID int64, roleIDs []int64, assignedBy *int64) error {
	args := m.Called(userID, roleIDs, assignedBy)
	return args.Error(0)
}

func (m *MockRBACRepository) RevokeRolesFromUser(ctx context.Context, userID int64, roleIDs []int64) error {
	args := m.Called(userID, roleIDs)
	return args.Error(0)
}

func (m *MockRBACRepository) GetUserRoles(ctx context.Context, userID int64) ([]models.Role, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Role), args.Error(1)
}

func (m *MockRBACRepository) GetUsersRoles(ctx context.Context, userIDs []int64) (map[int64][]models.Role, error) {
	args := m.Called(userIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[int64][]models.Role), args.Error(1)
}

func (m *MockRBACRepository) SyncUserRoles(ctx context.Context, userID int64, roleIDs []int64, assignedBy *int64) error {
	args := m.Called(userID, roleIDs, assignedBy)
	return args.Error(0)
}

// ─── User ↔ Permission (direct) ───────────────────────────────────────────────

func (m *MockRBACRepository) AssignDirectPermission(ctx context.Context, userID, permissionID int64, isGranted bool, assignedBy *int64) error {
	args := m.Called(userID, permissionID, isGranted, assignedBy)
	return args.Error(0)
}

func (m *MockRBACRepository) RevokeDirectPermission(ctx context.Context, userID, permissionID int64) error {
	args := m.Called(userID, permissionID)
	return args.Error(0)
}

func (m *MockRBACRepository) GetUserDirectPermissions(ctx context.Context, userID int64) ([]models.UserPermission, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.UserPermission), args.Error(1)
}

// ─── Check ─────────────────────────────────────────────────────────────────────

func (m *MockRBACRepository) GetUserAllPermissions(ctx context.Context, userID int64) ([]string, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRBACRepository) HasPermission(ctx context.Context, userID int64, permission string) (bool, error) {
	args := m.Called(userID, permission)
	return args.Bool(0), args.Error(1)
}
