package mocks

import (
	"context"
	rbacModels "neosim_go/internal/modules/rbac/models"
	"github.com/stretchr/testify/mock"
)

// RBACRepositoryMock is a mock implementation of rbacContracts.RBACRepository
type RBACRepositoryMock struct {
	mock.Mock
}

func (m *RBACRepositoryMock) IsSuperadmin(ctx context.Context,userID int64) (bool, error) {
	args := m.Called(userID)
	return args.Bool(0), args.Error(1)
}

func (m *RBACRepositoryMock) CreatePermission(ctx context.Context,p *rbacModels.Permission) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *RBACRepositoryMock) GetPermissionByID(ctx context.Context,id int64) (*rbacModels.Permission, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rbacModels.Permission), args.Error(1)
}

func (m *RBACRepositoryMock) GetPermissionByName(ctx context.Context,name string) (*rbacModels.Permission, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rbacModels.Permission), args.Error(1)
}

func (m *RBACRepositoryMock) ListPermissions(ctx context.Context,page, pageSize int) ([]rbacModels.Permission, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]rbacModels.Permission), args.Get(1).(int64), args.Error(2)
}

func (m *RBACRepositoryMock) UpdatePermission(ctx context.Context,p *rbacModels.Permission) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *RBACRepositoryMock) DeletePermission(ctx context.Context,id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *RBACRepositoryMock) CreateRole(ctx context.Context,r *rbacModels.Role) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *RBACRepositoryMock) GetRoleByID(ctx context.Context,id int64) (*rbacModels.Role, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rbacModels.Role), args.Error(1)
}

func (m *RBACRepositoryMock) GetRoleByName(ctx context.Context,name string) (*rbacModels.Role, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rbacModels.Role), args.Error(1)
}

func (m *RBACRepositoryMock) ListRoles(ctx context.Context,page, pageSize int) ([]rbacModels.Role, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]rbacModels.Role), args.Get(1).(int64), args.Error(2)
}

func (m *RBACRepositoryMock) UpdateRole(ctx context.Context,r *rbacModels.Role) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *RBACRepositoryMock) DeleteRole(ctx context.Context,id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *RBACRepositoryMock) GetUsersRoles(ctx context.Context,userIDs []int64) (map[int64][]rbacModels.Role, error) {
	args := m.Called(userIDs)
	return args.Get(0).(map[int64][]rbacModels.Role), args.Error(1)
}

func (m *RBACRepositoryMock) AssignPermissionsToRole(ctx context.Context,roleID int64, permissionIDs []int64) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

func (m *RBACRepositoryMock) RevokePermissionsFromRole(ctx context.Context,roleID int64, permissionIDs []int64) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

func (m *RBACRepositoryMock) GetRolePermissions(ctx context.Context,roleID int64) ([]rbacModels.Permission, error) {
	args := m.Called(roleID)
	return args.Get(0).([]rbacModels.Permission), args.Error(1)
}

func (m *RBACRepositoryMock) SyncRolePermissions(ctx context.Context,roleID int64, permissionIDs []int64) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

func (m *RBACRepositoryMock) AssignRolesToUser(ctx context.Context,userID int64, roleIDs []int64, assignedBy *int64) error {
	args := m.Called(userID, roleIDs, assignedBy)
	return args.Error(0)
}

func (m *RBACRepositoryMock) RevokeRolesFromUser(ctx context.Context,userID int64, roleIDs []int64) error {
	args := m.Called(userID, roleIDs)
	return args.Error(0)
}

func (m *RBACRepositoryMock) GetUserRoles(ctx context.Context,userID int64) ([]rbacModels.Role, error) {
	args := m.Called(userID)
	return args.Get(0).([]rbacModels.Role), args.Error(1)
}

func (m *RBACRepositoryMock) SyncUserRoles(ctx context.Context,userID int64, roleIDs []int64, assignedBy *int64) error {
	args := m.Called(userID, roleIDs, assignedBy)
	return args.Error(0)
}

func (m *RBACRepositoryMock) AssignDirectPermission(ctx context.Context,userID, permissionID int64, isGranted bool, assignedBy *int64) error {
	args := m.Called(userID, permissionID, isGranted, assignedBy)
	return args.Error(0)
}

func (m *RBACRepositoryMock) RevokeDirectPermission(ctx context.Context,userID, permissionID int64) error {
	args := m.Called(userID, permissionID)
	return args.Error(0)
}

func (m *RBACRepositoryMock) GetUserDirectPermissions(ctx context.Context,userID int64) ([]rbacModels.UserPermission, error) {
	args := m.Called(userID)
	return args.Get(0).([]rbacModels.UserPermission), args.Error(1)
}

func (m *RBACRepositoryMock) GetUserAllPermissions(ctx context.Context,userID int64) ([]string, error) {
	args := m.Called(userID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *RBACRepositoryMock) HasPermission(ctx context.Context,userID int64, permission string) (bool, error) {
	args := m.Called(userID, permission)
	return args.Bool(0), args.Error(1)
}
