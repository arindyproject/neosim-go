package services

import (
	"context"
	rbacMiddlewares "neosim_go/internal/modules/rbac/middlewares"
	rbacModels "neosim_go/internal/modules/rbac/models"
	he "neosim_go/internal/shared/httputil"
)

func (s *service) canCreateUser(ctx context.Context, actor he.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(ctx, s.rbacRepo, actor.UserID, rbacModels.PermUsersCreate); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(ctx, s.rbacRepo, actor.UserID, rbacModels.PermUsersManage); err != nil || has {
		return has, err
	}
	return rbacMiddlewares.HasAnyRole(ctx, s.rbacRepo, actor.UserID, "admin", "superadmin", "hrd")
}

func (s *service) canUpdateUser(ctx context.Context, actor he.AuthContext, targetUserID int64) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if actor.UserID == targetUserID {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(ctx, s.rbacRepo, actor.UserID, rbacModels.PermUsersUpdate); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(ctx, s.rbacRepo, actor.UserID, rbacModels.PermUsersManage); err != nil || has {
		return has, err
	}
	return rbacMiddlewares.HasAnyRole(ctx, s.rbacRepo, actor.UserID, "admin", "superadmin", "hrd")
}

func (s *service) canDeleteUser(ctx context.Context, actor he.AuthContext) (bool, error) {
	return actor.IsSuperadmin, nil
}

func (s *service) canReadUser(ctx context.Context, actor he.AuthContext, targetUserID int64) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if actor.UserID == targetUserID {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(ctx, s.rbacRepo, actor.UserID, rbacModels.PermUsersRead); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(ctx, s.rbacRepo, actor.UserID, rbacModels.PermUsersManage); err != nil || has {
		return has, err
	}
	return false, nil
}
