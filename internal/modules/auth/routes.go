// ════════════════════════════════════════════════════════════
// FILE: internal/modules/auth/routes.go
// ════════════════════════════════════════════════════════════
package auth

import (
	"neosim_go/internal/modules/auth/handlers"
	authMiddlewares "neosim_go/internal/modules/auth/middlewares"
	"neosim_go/internal/modules/auth/utils"

	"github.com/labstack/echo/v5"
)

// RegisterRoutes mendaftarkan semua routes untuk module auth
func RegisterRoutes(e *echo.Echo, h *handlers.AuthHandler, jwtManager *utils.JWTManager) {
	// ─── Public routes (tidak butuh login) ────────────────────
	public := e.Group("/api/v1/auth")
	public.POST("/login", h.Login)
	public.POST("/register", h.Register)
	public.POST("/refresh", h.RefreshToken)
	public.POST("/forgot-password", h.ForgotPassword)
	public.POST("/reset-password", h.ResetPassword)

	// ─── Protected routes (butuh login) ───────────────────────
	// Logout butuh JWT untuk:
	// 1. Verifikasi user yang request adalah pemilik sesi
	// 2. LogoutAll butuh userID dari claims
	protected := e.Group("/api/v1/auth", authMiddlewares.JWTMiddleware(jwtManager))
	protected.POST("/logout", h.Logout)
	protected.POST("/logout-all", h.LogoutAll)
}
