package middlewares

import (
	"net/http"
	"strconv"

	"neosim_go/internal/modules/rbac/contracts"
	"neosim_go/internal/shared/response"

	"github.com/labstack/echo/v5"
)

// ─── Middleware ────────────────────────────────────────────────────────────────

// RequirePermission memastikan user memiliki permission tertentu
func RequirePermission(repo contracts.RBACRepository, permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if IsSuperuser(c) {
				return next(c)
			}
			userID, ok := GetUserIDFromContext(c)
			if !ok {
				return response.Response(c, http.StatusUnauthorized, false, "Autentikasi diperlukan", nil, nil)
			}
			has, err := repo.HasPermission(userID, permission)
			if err != nil {
				return response.Response(c, http.StatusInternalServerError, false, "Gagal cek permission", nil, nil)
			}
			if !has {
				return response.Response(c, http.StatusForbidden, false,
					"Akses ditolak. Permission '"+permission+"' diperlukan.", nil, nil)
			}
			return next(c)
		}
	}
}

// RequireAnyPermission memastikan user memiliki minimal satu permission
func RequireAnyPermission(repo contracts.RBACRepository, permissions ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if IsSuperuser(c) {
				return next(c)
			}
			userID, ok := GetUserIDFromContext(c)
			if !ok {
				return response.Response(c, http.StatusUnauthorized, false, "Autentikasi diperlukan", nil, nil)
			}
			userPerms, err := repo.GetUserAllPermissions(userID)
			if err != nil {
				return response.Response(c, http.StatusInternalServerError, false, "Gagal cek permission", nil, nil)
			}
			permSet := toSet(userPerms)
			for _, p := range permissions {
				if permSet[p] {
					return next(c)
				}
			}
			return response.Response(c, http.StatusForbidden, false, "Akses ditolak. Permission tidak mencukupi.", nil, nil)
		}
	}
}

// RequireRole memastikan user memiliki role tertentu
func RequireRole(repo contracts.RBACRepository, roleName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if IsSuperuser(c) {
				return next(c)
			}
			userID, ok := GetUserIDFromContext(c)
			if !ok {
				return response.Response(c, http.StatusUnauthorized, false, "Autentikasi diperlukan", nil, nil)
			}
			has, err := HasRole(repo, userID, roleName)
			if err != nil {
				return response.Response(c, http.StatusInternalServerError, false, "Gagal cek role", nil, nil)
			}
			if !has {
				return response.Response(c, http.StatusForbidden, false,
					"Akses ditolak. Role '"+roleName+"' diperlukan.", nil, nil)
			}
			return next(c)
		}
	}
}

// ─── Context Helpers ───────────────────────────────────────────────────────────

func IsSuperuser(c *echo.Context) bool {
	v, _ := c.Get("isSuperuser").(bool)
	return v
}

func GetUserIDFromContext(c *echo.Context) (int64, bool) {
	id, ok := c.Get("userID").(int64)
	return id, ok
}

func GetTargetUserID(c *echo.Context) (int64, error) {
	return strconv.ParseInt(c.Param("id"), 10, 64)
}

// ─── Programmatic Helpers (untuk service/handler) ─────────────────────────────

func HasPermission(repo contracts.RBACRepository, userID int64, permission string) (bool, error) {
	return repo.HasPermission(userID, permission)
}

func HasRole(repo contracts.RBACRepository, userID int64, roleName string) (bool, error) {
	roles, err := repo.GetUserRoles(userID)
	if err != nil {
		return false, err
	}
	for _, r := range roles {
		if r.Name == roleName {
			return true, nil
		}
	}
	return false, nil
}

func HasAnyRole(repo contracts.RBACRepository, userID int64, roleNames ...string) (bool, error) {
	roles, err := repo.GetUserRoles(userID)
	if err != nil {
		return false, err
	}
	roleSet := make(map[string]bool, len(roles))
	for _, r := range roles {
		roleSet[r.Name] = true
	}
	for _, name := range roleNames {
		if roleSet[name] {
			return true, nil
		}
	}
	return false, nil
}

func toSet(items []string) map[string]bool {
	set := make(map[string]bool, len(items))
	for _, v := range items {
		set[v] = true
	}
	return set
}
