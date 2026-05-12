package users

import (
	authMiddlewares "neosim_go/internal/modules/auth/middlewares"
	"neosim_go/internal/modules/auth/utils"
	"neosim_go/internal/modules/rbac/contracts"
	rbacMiddlewares "neosim_go/internal/modules/rbac/middlewares"
	rbacModels "neosim_go/internal/modules/rbac/models"
	"neosim_go/internal/modules/users/handlers"

	"github.com/labstack/echo/v5"
)

// RegisterRoutes mendaftarkan semua routes users dengan RBAC
//
// Aturan akses UpdateUser (dicek di service):
// - Superadmin         → boleh
// - Diri sendiri       → boleh
// - Permission users:update → boleh
// - Role hrd           → boleh
func RegisterRoutes(e *echo.Echo, h *handlers.Handler, rbacRepo contracts.RBACRepository, jwtManager *utils.JWTManager) {
	jwt := authMiddlewares.JWTMiddleware(jwtManager)

	// ─── Public (tidak butuh login) ────────────────────────────
	public := e.Group("/api/v1/users")
	public.GET("/username/:username", h.GetByUsernameHandler)

	// ─── Protected (butuh login) ───────────────────────────────
	protected := e.Group("/api/v1/users", jwt)

	// List — butuh permission users:read
	protected.GET("", h.ListUsersHandler,
		rbacMiddlewares.RequirePermission(rbacRepo, rbacModels.PermUsersRead),
	)

	// Get by ID — authorization di service (diri sendiri atau punya permission)
	protected.GET("/:id", h.GetUserHandler)

	// Create — butuh permission users:create
	protected.POST("", h.CreateUserHandler,
		rbacMiddlewares.RequirePermission(rbacRepo, rbacModels.PermUsersCreate),
	)

	// Update — route tidak ada middleware khusus karena authorization
	// dilakukan di service (superadmin | diri sendiri | users:update | role hrd)
	protected.PUT("/:id", h.UpdateUserHandler)

	// Delete — butuh permission users:delete atau superadmin (dicek di service)
	protected.DELETE("/:id", h.DeleteUserHandler)

	// Change Password — authorization di service (diri sendiri atau superadmin)
	protected.PUT("/:id/change-password", h.ChangePasswordHandler)

	// Settings — diri sendiri atau superadmin
	protected.GET("/:id/settings", h.GetSettingsHandler)
	protected.PUT("/:id/settings", h.UpdateSettingsHandler)
}
