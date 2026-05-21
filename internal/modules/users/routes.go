package users

import (
	authMiddlewares "neosim_go/internal/modules/auth/middlewares"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
	rbacMiddlewares "neosim_go/internal/modules/rbac/middlewares"
	rbacModels "neosim_go/internal/modules/rbac/models"
	"neosim_go/internal/modules/users/handlers"
	"neosim_go/internal/shared/utils"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// RegisterRoutes mendaftarkan semua routes users dengan RBAC
// db dibutuhkan oleh JWTMiddleware untuk cek isSuperadmin secara realtime
func RegisterRoutes(e *echo.Echo, h *handlers.Handler, rbacRepo rbacContracts.RBACRepository, jwtManager *utils.JWTManager, db *gorm.DB) {
	jwt := authMiddlewares.JWTMiddleware(jwtManager, db)

	// ─── Public ────────────────────────────────────────────────
	public := e.Group("/api/v1/users")
	public.GET("/username/:username", h.GetByUsernameHandler)

	// ─── Protected ─────────────────────────────────────────────
	protected := e.Group("/api/v1/users", jwt)

	protected.GET("", h.ListUsersHandler,
		rbacMiddlewares.RequirePermission(rbacRepo, rbacModels.PermUsersRead),
	)
	protected.GET("/:id", h.GetUserHandler)
	protected.POST("", h.CreateUserHandler,
		rbacMiddlewares.RequirePermission(rbacRepo, rbacModels.PermUsersCreate),
	)
	protected.PUT("/:id", h.UpdateUserHandler)
	protected.DELETE("/:id", h.DeleteUserHandler)
	protected.PUT("/:id/change-password", h.ChangePasswordHandler)
	protected.GET("/:id/settings", h.GetSettingsHandler)
	protected.PUT("/:id/settings", h.UpdateSettingsHandler)
}
