package services

import (
	"context"
	rbacMiddlewares "neosim_go/internal/modules/rbac/middlewares"
	rbacModels "neosim_go/internal/modules/rbac/models"
	he "neosim_go/internal/shared/httputil"
)

// ─── Permission ─────────────────────────────────────────────────────────────────
func (s *service) canCreateMaster(ctx context.Context, actor he.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(ctx, s.rbacRepo, actor.UserID, rbacModels.PermMasterCreate); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(ctx, s.rbacRepo, actor.UserID, rbacModels.PermMasterManage); err != nil || has {
		return has, err
	}
	return false, nil
}

func (s *service) canUpdateMaster(ctx context.Context, actor he.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(ctx, s.rbacRepo, actor.UserID, rbacModels.PermMasterUpdate); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(ctx, s.rbacRepo, actor.UserID, rbacModels.PermMasterManage); err != nil || has {
		return has, err
	}
	return false, nil
}

func (s *service) canDeleteMaster(ctx context.Context, actor he.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(ctx, s.rbacRepo, actor.UserID, rbacModels.PermMasterDelete); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(ctx, s.rbacRepo, actor.UserID, rbacModels.PermMasterManage); err != nil || has {
		return has, err
	}
	return false, nil
}
