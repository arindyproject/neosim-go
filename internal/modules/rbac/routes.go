// ════════════════════════════════════════════════════════════════
// FILE: internal/modules/rbac/routes.go
// ════════════════════════════════════════════════════════════════
package rbac

import (
	"neosim_go/internal/modules/auth/middlewares"
	"neosim_go/internal/modules/rbac/contracts"
	"neosim_go/internal/modules/rbac/handlers"
	rbacMiddlewares "neosim_go/internal/modules/rbac/middlewares"
	rbacModels "neosim_go/internal/modules/rbac/models"
	"neosim_go/internal/shared/utils"

	"github.com/labstack/echo/v5"
)

func RegisterRoutes(e *echo.Echo, h *handlers.RBACHandler, repo contracts.RBACRepository, jwtManager *utils.JWTManager) {
	jwt := middlewares.JWTMiddleware(jwtManager)

	// ─── Permissions (superadmin only) ─────────────────────────
	perms := e.Group("/api/v1/permissions",
		jwt,
		rbacMiddlewares.RequirePermission(repo, rbacModels.PermPermissionsManage),
	)
	perms.GET("", h.ListPermissions)
	perms.GET("/:id", h.GetPermission)
	perms.POST("", h.CreatePermission)
	perms.PUT("/:id", h.UpdatePermission)
	perms.DELETE("/:id", h.DeletePermission)

	// ─── Roles ─────────────────────────────────────────────────
	roles := e.Group("/api/v1/roles", jwt)
	roles.GET("", h.ListRoles,
		rbacMiddlewares.RequirePermission(repo, rbacModels.PermRolesRead),
	)
	roles.GET("/:id", h.GetRole,
		rbacMiddlewares.RequirePermission(repo, rbacModels.PermRolesRead),
	)
	roles.POST("", h.CreateRole,
		rbacMiddlewares.RequirePermission(repo, rbacModels.PermRolesCreate),
	)
	roles.PUT("/:id", h.UpdateRole,
		rbacMiddlewares.RequirePermission(repo, rbacModels.PermRolesUpdate),
	)
	roles.DELETE("/:id", h.DeleteRole,
		rbacMiddlewares.RequirePermission(repo, rbacModels.PermRolesDelete),
	)

	// ─── Role ↔ Permission ─────────────────────────────────────
	roles.POST("/:id/permissions", h.AssignPermissionsToRole,
		rbacMiddlewares.RequirePermission(repo, rbacModels.PermRolesManage),
	)
	roles.PUT("/:id/permissions", h.SyncRolePermissions,
		rbacMiddlewares.RequirePermission(repo, rbacModels.PermRolesManage),
	)
	roles.DELETE("/:id/permissions", h.RevokePermissionsFromRole,
		rbacMiddlewares.RequirePermission(repo, rbacModels.PermRolesManage),
	)

	// ─── User ↔ Role ───────────────────────────────────────────
	userRoles := e.Group("/api/v1/users/:user_id/roles",
		jwt,
		rbacMiddlewares.RequirePermission(repo, rbacModels.PermUsersManage),
	)
	userRoles.GET("", h.GetUserRoles)
	userRoles.POST("", h.AssignRolesToUser)
	userRoles.PUT("", h.SyncUserRoles)
	userRoles.DELETE("", h.RevokeRolesFromUser)

	// ─── User Permissions ──────────────────────────────────────
	userPerms := e.Group("/api/v1/users/:user_id/permissions",
		jwt,
		rbacMiddlewares.RequirePermission(repo, rbacModels.PermUsersManage),
	)
	userPerms.GET("", h.GetUserAllPermissions)
	userPerms.POST("", h.AssignDirectPermission)
}
