package services

import (
	"context"
	rbacMiddlewares "neosim_go/internal/modules/rbac/middlewares"
	rbacModels "neosim_go/internal/modules/rbac/models"
	he "neosim_go/internal/shared/httputil"
)


// ── canRead ───────────────────────────────────────────────────────────────────
func (s *service) canReadArtikelKategori(ctx context.Context,actor he.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(ctx,s.rbacRepo, actor.UserID, rbacModels.PermAnyRead); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(ctx,s.rbacRepo, actor.UserID, rbacModels.PermAnyManage); err != nil || has {
		return has, err
	}
	return false, nil
}


// ── canCreate ─────────────────────────────────────────────────────────────────
func (s *service) canCreateArtikelKategori(ctx context.Context,actor he.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(ctx,s.rbacRepo, actor.UserID, rbacModels.PermAnyCreate); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(ctx,s.rbacRepo, actor.UserID, rbacModels.PermAnyManage); err != nil || has {
		return has, err
	}
	return false, nil
}


// ── canUpdate ─────────────────────────────────────────────────────────────────
func (s *service) canUpdateArtikelKategori(ctx context.Context,actor he.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(ctx,s.rbacRepo, actor.UserID, rbacModels.PermAnyUpdate); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(ctx,s.rbacRepo, actor.UserID, rbacModels.PermAnyManage); err != nil || has {
		return has, err
	}
	return false, nil
}


// ── canDelete ─────────────────────────────────────────────────────────────────
func (s *service) canDeleteArtikelKategori(ctx context.Context,actor he.AuthContext) (bool, error) {
	if actor.IsSuperadmin {
		return true, nil
	}
	if has, err := rbacMiddlewares.HasPermission(ctx,s.rbacRepo, actor.UserID, rbacModels.PermAnyDelete); err != nil || has {
		return has, err
	}
	if has, err := rbacMiddlewares.HasPermission(ctx,s.rbacRepo, actor.UserID, rbacModels.PermAnyManage); err != nil || has {
		return has, err
	}
	return false, nil
}
