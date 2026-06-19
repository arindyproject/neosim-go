package mocks

import (
	rbacModels "neosim_go/internal/modules/rbac/models"
	"github.com/stretchr/testify/mock"
)

// RBACRepositoryMock is a mock implementation of rbacContracts.RBACRepository
type RBACRepositoryMock struct {
	mock.Mock
}

func (m *RBACRepositoryMock) IsSuperadmin(userID int64) (bool, error) {
	args := m.Called(userID)
	return args.Bool(0), args.Error(1)
}

func (m *RBACRepositoryMock) CreatePermission(p *rbacModels.Permission) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *RBACRepositoryMock) GetPermissionByID(id int64) (*rbacModels.Permission, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rbacModels.Permission), args.Error(1)
}

func (m *RBACRepositoryMock) GetPermissionByName(name string) (*rbacModels.Permission, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rbacModels.Permission), args.Error(1)
}

func (m *RBACRepositoryMock) ListPermissions(page, pageSize int) ([]rbacModels.Permission, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]rbacModels.Permission), args.Get(1).(int64), args.Error(2)
}

func (m *RBACRepositoryMock) UpdatePermission(p *rbacModels.Permission) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *RBACRepositoryMock) DeletePermission(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *RBACRepositoryMock) CreateRole(r *rbacModels.Role) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *RBACRepositoryMock) GetRoleByID(id int64) (*rbacModels.Role, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rbacModels.Role), args.Error(1)
}

func (m *RBACRepositoryMock) GetRoleByName(name string) (*rbacModels.Role, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rbacModels.Role), args.Error(1)
}

func (m *RBACRepositoryMock) ListRoles(page, pageSize int) ([]rbacModels.Role, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]rbacModels.Role), args.Get(1).(int64), args.Error(2)
}

func (m *RBACRepositoryMock) UpdateRole(r *rbacModels.Role) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *RBACRepositoryMock) DeleteRole(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *RBACRepositoryMock) GetUsersRoles(userIDs []int64) (map[int64][]rbacModels.Role, error) {
	args := m.Called(userIDs)
	return args.Get(0).(map[int64][]rbacModels.Role), args.Error(1)
}

func (m *RBACRepositoryMock) AssignPermissionsToRole(roleID int64, permissionIDs []int64) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

func (m *RBACRepositoryMock) RevokePermissionsFromRole(roleID int64, permissionIDs []int64) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

func (m *RBACRepositoryMock) GetRolePermissions(roleID int64) ([]rbacModels.Permission, error) {
	args := m.Called(roleID)
	return args.Get(0).([]rbacModels.Permission), args.Error(1)
}

func (m *RBACRepositoryMock) SyncRolePermissions(roleID int64, permissionIDs []int64) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

func (m *RBACRepositoryMock) AssignRolesToUser(userID int64, roleIDs []int64, assignedBy *int64) error {
	args := m.Called(userID, roleIDs, assignedBy)
	return args.Error(0)
}

func (m *RBACRepositoryMock) RevokeRolesFromUser(userID int64, roleIDs []int64) error {
	args := m.Called(userID, roleIDs)
	return args.Error(0)
}

func (m *RBACRepositoryMock) GetUserRoles(userID int64) ([]rbacModels.Role, error) {
	args := m.Called(userID)
	return args.Get(0).([]rbacModels.Role), args.Error(1)
}

func (m *RBACRepositoryMock) SyncUserRoles(userID int64, roleIDs []int64, assignedBy *int64) error {
	args := m.Called(userID, roleIDs, assignedBy)
	return args.Error(0)
}

func (m *RBACRepositoryMock) AssignDirectPermission(userID, permissionID int64, isGranted bool, assignedBy *int64) error {
	args := m.Called(userID, permissionID, isGranted, assignedBy)
	return args.Error(0)
}

func (m *RBACRepositoryMock) RevokeDirectPermission(userID, permissionID int64) error {
	args := m.Called(userID, permissionID)
	return args.Error(0)
}

func (m *RBACRepositoryMock) GetUserDirectPermissions(userID int64) ([]rbacModels.UserPermission, error) {
	args := m.Called(userID)
	return args.Get(0).([]rbacModels.UserPermission), args.Error(1)
}

func (m *RBACRepositoryMock) GetUserAllPermissions(userID int64) ([]string, error) {
	args := m.Called(userID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *RBACRepositoryMock) HasPermission(userID int64, permission string) (bool, error) {
	args := m.Called(userID, permission)
	return args.Bool(0), args.Error(1)
}
