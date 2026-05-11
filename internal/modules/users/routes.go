package users

import (
	"neosim_go/internal/modules/users/handlers"

	"neosim_go/internal/modules/auth/middlewares"
	"neosim_go/internal/modules/auth/utils"

	"github.com/labstack/echo/v5"
)

// RegisterRoutes registers user routes to the echo instance
func RegisterRoutes(e *echo.Echo, handler *handlers.Handler, jwtManager *utils.JWTManager) {
	// Public routes (no auth required)
	userGroup := e.Group("/api/v1/users") // Apply JWT auth middleware to all user routes
	userGroup.Use(middlewares.JWTMiddleware(jwtManager))

	// Create user
	userGroup.POST("", handler.CreateUserHandler)

	// List users
	userGroup.GET("", handler.ListUsersHandler)

	// Get user by ID
	userGroup.GET("/:id", handler.GetUserHandler)

	// Get user by username
	userGroup.GET("/username/:username", handler.GetByUsernameHandler)

	// Update user
	userGroup.PUT("/:id", handler.UpdateUserHandler)

	// Delete user
	userGroup.DELETE("/:id", handler.DeleteUserHandler)

	// Change password
	userGroup.POST("/:id/change-password", handler.ChangePasswordHandler)

	// Get user settings
	userGroup.GET("/:id/settings", handler.GetSettingsHandler)

	// Update user settings
	userGroup.PUT("/:id/settings", handler.UpdateSettingsHandler)
}
