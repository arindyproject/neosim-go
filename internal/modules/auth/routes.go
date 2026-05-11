// ════════════════════════════════════════════════════════════
// FILE: internal/modules/auth/routes.go
// ════════════════════════════════════════════════════════════
package auth

import (
	"neosim_go/internal/modules/auth/handlers"

	"github.com/labstack/echo/v5"
)

// RegisterRoutes mendaftarkan semua routes untuk module auth
func RegisterRoutes(e *echo.Echo, h *handlers.AuthHandler) {
	g := e.Group("/api/v1/auth")
	g.POST("/login", h.Login)
	g.POST("/register", h.Register)
	g.POST("/refresh", h.RefreshToken)
	g.POST("/forgot-password", h.ForgotPassword)
	g.POST("/reset-password", h.ResetPassword)
}
