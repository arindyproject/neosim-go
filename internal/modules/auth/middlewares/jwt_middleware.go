package middlewares

import (
	"net/http"
	"strings"

	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/utils"

	"github.com/labstack/echo/v5"
)

// ─── JWT Middleware ────────────────────────────────────────────────────────────

// JWTMiddleware memvalidasi access token dari header Authorization
func JWTMiddleware(jwtManager *utils.JWTManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// Ambil token dari header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Response(c, http.StatusUnauthorized, false, "Token tidak ditemukan.", nil, nil)
			}

			// Format: "Bearer <token>"
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return response.Response(c, http.StatusUnauthorized, false, "Format token tidak valid.", nil, nil)
			}

			tokenStr := parts[1]

			// Parse & validasi token
			claims, err := jwtManager.ParseToken(tokenStr)
			if err != nil {
				return response.Response(c, http.StatusUnauthorized, false, "Token tidak valid atau sudah kedaluwarsa.", nil, nil)
			}

			// Pastikan hanya access token yang diterima
			if claims.TokenType != "access" {
				return response.Response(c, http.StatusUnauthorized, false, "Tipe token tidak valid.", nil, nil)
			}

			// Set claims ke context — bisa diambil di handler
			c.Set("userID", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("isSuperuser", claims.IsSuperuser)
			c.Set("isStaff", claims.IsStaff)

			return next(c)
		}
	}
}

// ─── Authorization Middlewares ─────────────────────────────────────────────────

// RequireSuperuser memastikan user adalah superuser
func RequireSuperuser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			isSuperuser, _ := c.Get("isSuperuser").(bool)
			if !isSuperuser {
				return response.Response(c, http.StatusForbidden, false, "Akses ditolak. Hanya superuser.", nil, nil)
			}
			return next(c)
		}
	}
}

// RequireStaff memastikan user adalah staff atau superuser
func RequireStaff() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			isStaff, _ := c.Get("isStaff").(bool)
			isSuperuser, _ := c.Get("isSuperuser").(bool)
			if !isStaff && !isSuperuser {
				return response.Response(c, http.StatusForbidden, false, "Akses ditolak. Hanya staff.", nil, nil)
			}
			return next(c)
		}
	}
}

// ─── Helper untuk Handler ──────────────────────────────────────────────────────

// GetUserID mengambil userID dari context (set oleh JWTMiddleware)
func GetUserID(c *echo.Context) (int64, bool) {
	userID, ok := c.Get("userID").(int64)
	return userID, ok
}

// GetUsername mengambil username dari context
func GetUsername(c *echo.Context) (string, bool) {
	username, ok := c.Get("username").(string)
	return username, ok
}
