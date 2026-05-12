package handlers

import (
	"net/http"
	"strconv"

	rbacMiddlewares "neosim_go/internal/modules/rbac/middlewares"
	userContracts "neosim_go/internal/modules/users/contracts"
	"neosim_go/internal/modules/users/dto"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	service userContracts.Service
}

func NewHandler(service userContracts.Service) *Handler {
	return &Handler{service: service}
}

// ─── Helpers ───────────────────────────────────────────────────────────────────

func parseID(c *echo.Context) (int64, error) {
	return strconv.ParseInt(c.Param("id"), 10, 64)
}

// buildAuthContext membuat AuthContext dari JWT claims di context
func buildAuthContext(c *echo.Context) userContracts.AuthContext {
	userID, _ := rbacMiddlewares.GetUserIDFromContext(c)
	isSuperuser := rbacMiddlewares.IsSuperuser(c)
	return userContracts.AuthContext{
		UserID:      userID,
		IsSuperuser: isSuperuser,
	}
}

// ─── User CRUD ─────────────────────────────────────────────────────────────────

// ListUsersHandler handles GET /api/v1/users
// Siapa yang bisa: semua yang login (data dirinya sendiri disaring di service)
func (h *Handler) ListUsersHandler(c *echo.Context) error {
	page, pageSize := 1, 10
	if p := c.QueryParam("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := c.QueryParam("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 && v <= 100 {
			pageSize = v
		}
	}

	users, total, err := h.service.ListUsers(page, pageSize)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data user", nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data user", users, total, page, pageSize)
}

// GetUserHandler handles GET /api/v1/users/:id
// Siapa yang bisa: superadmin, diri sendiri, atau yang punya permission users:read
func (h *Handler) GetUserHandler(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := buildAuthContext(c)
	user, err := h.service.GetUserByID(id, actor)
	if err != nil {
		if appErr, ok := err.(interface{ StatusCode() int }); ok {
			return response.Response(c, appErr.StatusCode(), false, err.Error(), nil, nil)
		}
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data user", user, nil)
}

// CreateUserHandler handles POST /api/v1/users
func (h *Handler) CreateUserHandler(c *echo.Context) error {
	var req dto.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actorID, _ := rbacMiddlewares.GetUserIDFromContext(c)
	user, err := h.service.CreateUser(&req, &actorID)
	if err != nil {
		if appErr, ok := err.(interface{ StatusCode() int }); ok {
			return response.Response(c, appErr.StatusCode(), false, err.Error(), nil, nil)
		}
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "User berhasil dibuat", user, nil)
}

// UpdateUserHandler handles PUT /api/v1/users/:id
//
// Authorization (dicek di service):
// - Superadmin → boleh
// - Diri sendiri → boleh (tapi tidak bisa ubah is_active, is_staff, is_superuser)
// - Punya permission "users:update" → boleh
// - Punya role "hrd" → boleh
func (h *Handler) UpdateUserHandler(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := buildAuthContext(c)
	user, err := h.service.UpdateUser(id, &req, actor)
	if err != nil {
		if appErr, ok := err.(interface{ StatusCode() int }); ok {
			return response.Response(c, appErr.StatusCode(), false, err.Error(), nil, nil)
		}
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "User berhasil diupdate", user, nil)
}

// DeleteUserHandler handles DELETE /api/v1/users/:id
// Siapa yang bisa: superadmin atau punya permission users:delete
func (h *Handler) DeleteUserHandler(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := buildAuthContext(c)
	if err := h.service.DeleteUser(id, actor); err != nil {
		if appErr, ok := err.(interface{ StatusCode() int }); ok {
			return response.Response(c, appErr.StatusCode(), false, err.Error(), nil, nil)
		}
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "User berhasil dihapus", nil, nil)
}

// ─── Password ──────────────────────────────────────────────────────────────────

// ChangePasswordHandler handles PUT /api/v1/users/:id/change-password
// Siapa yang bisa: diri sendiri atau superadmin
func (h *Handler) ChangePasswordHandler(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := buildAuthContext(c)
	if err := h.service.ChangePassword(id, &req, actor); err != nil {
		if appErr, ok := err.(interface{ StatusCode() int }); ok {
			return response.Response(c, appErr.StatusCode(), false, err.Error(), nil, nil)
		}
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Password berhasil diubah", nil, nil)
}

// ─── Settings ──────────────────────────────────────────────────────────────────

func (h *Handler) GetSettingsHandler(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	settings, err := h.service.GetSettings(id)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil settings", settings, nil)
}

func (h *Handler) UpdateSettingsHandler(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	var req dto.UpdateSettingsRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}
	if err := h.service.UpdateSettings(id, &req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Settings berhasil diupdate", nil, nil)
}

// ─── By Username ───────────────────────────────────────────────────────────────

func (h *Handler) GetByUsernameHandler(c *echo.Context) error {
	username := c.Param("username")
	if username == "" {
		return response.Response(c, http.StatusBadRequest, false, "Username tidak boleh kosong", nil, nil)
	}
	user, err := h.service.GetUserByUsername(username)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data user", user, nil)
}
