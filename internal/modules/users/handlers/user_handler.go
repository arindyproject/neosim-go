package handlers

import (
	"net/http"
	"strconv"

	"neosim_go/internal/modules/users/contracts"
	"neosim_go/internal/modules/users/dto"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"

	"github.com/labstack/echo/v5"
)

// Handler defines HTTP handlers for user operations
type Handler struct {
	service contracts.Service
}

// NewHandler creates a new handler instance
func NewHandler(service contracts.Service) *Handler {
	return &Handler{service: service}
}

// ─── Private Helpers ───────────────────────────────────────────────────────────

func parseID(c *echo.Context) (int64, error) {
	return strconv.ParseInt(c.Param("id"), 10, 64)
}

func getActorID(c *echo.Context) *int64 {
	if userID, ok := c.Get("userID").(int64); ok {
		return &userID
	}
	return nil
}

// ─── User CRUD ─────────────────────────────────────────────────────────────────

// CreateUserHandler handles POST /api/v1/users
func (h *Handler) CreateUserHandler(c *echo.Context) error {
	var req dto.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	user, err := h.service.CreateUser(&req, getActorID(c))
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusCreated, true, "User berhasil dibuat", user, nil)
}

// GetUserHandler handles GET /api/v1/users/:id
func (h *Handler) GetUserHandler(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data user", user, nil)
}

// ListUsersHandler handles GET /api/v1/users
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

// UpdateUserHandler handles PUT /api/v1/users/:id
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

	user, err := h.service.UpdateUser(id, &req, getActorID(c))
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "user tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "User berhasil diupdate", user, nil)
}

// DeleteUserHandler handles DELETE /api/v1/users/:id
func (h *Handler) DeleteUserHandler(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	if err := h.service.DeleteUser(id); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "User berhasil dihapus", nil, nil)
}

// ─── Password ──────────────────────────────────────────────────────────────────

// ChangePasswordHandler handles POST /api/v1/users/:id/change-password
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

	if err := h.service.ChangePassword(id, &req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Password berhasil diubah", nil, nil)
}

// ─── Settings ──────────────────────────────────────────────────────────────────

// GetSettingsHandler handles GET /api/v1/users/:id/settings
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

// UpdateSettingsHandler handles PUT /api/v1/users/:id/settings
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

// GetByUsernameHandler handles GET /api/v1/users/username/:username
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
