package handlers

import (
	"net/http"
	"strconv"

	"neosim_go/internal/modules/rbac/contracts"
	"neosim_go/internal/modules/rbac/dto"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"

	"github.com/labstack/echo/v5"
)

type RBACHandler struct {
	service contracts.RBACService
}

func NewRBACHandler(service contracts.RBACService) *RBACHandler {
	return &RBACHandler{service: service}
}

func parseID(c *echo.Context) (int64, error) {
	return strconv.ParseInt(c.Param("id"), 10, 64)
}

func getActorID(c *echo.Context) *int64 {
	if userID, ok := c.Get("userID").(int64); ok {
		return &userID
	}
	return nil
}

// ─── Permission Handlers ───────────────────────────────────────────────────────

func (h *RBACHandler) ListPermissions(c *echo.Context) error {
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

	items, total, err := h.service.ListPermissions(page, pageSize)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data permission", nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data permission", items, total, page, pageSize)
}

func (h *RBACHandler) GetPermission(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	item, err := h.service.GetPermissionByID(id)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil permission", item, nil)
}

func (h *RBACHandler) CreatePermission(c *echo.Context) error {
	var req dto.CreatePermissionRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}
	item, err := h.service.CreatePermission(&req, getActorID(c))
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Permission berhasil dibuat", item, nil)
}

func (h *RBACHandler) UpdatePermission(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	var req dto.UpdatePermissionRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}
	item, err := h.service.UpdatePermission(id, &req, getActorID(c))
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Permission berhasil diupdate", item, nil)
}

func (h *RBACHandler) DeletePermission(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	if err := h.service.DeletePermission(id); err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Permission berhasil dihapus", nil, nil)
}

// ─── Role Handlers ─────────────────────────────────────────────────────────────

func (h *RBACHandler) ListRoles(c *echo.Context) error {
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
	items, total, err := h.service.ListRoles(page, pageSize)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data role", nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data role", items, total, page, pageSize)
}

func (h *RBACHandler) GetRole(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	item, err := h.service.GetRoleByID(id)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil role", item, nil)
}

func (h *RBACHandler) CreateRole(c *echo.Context) error {
	var req dto.CreateRoleRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}
	item, err := h.service.CreateRole(&req, getActorID(c))
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Role berhasil dibuat", item, nil)
}

func (h *RBACHandler) UpdateRole(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	var req dto.UpdateRoleRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}
	item, err := h.service.UpdateRole(id, &req, getActorID(c))
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Role berhasil diupdate", item, nil)
}

func (h *RBACHandler) DeleteRole(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	if err := h.service.DeleteRole(id); err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Role berhasil dihapus", nil, nil)
}

// ─── Role ↔ Permission Handlers ────────────────────────────────────────────────

func (h *RBACHandler) AssignPermissionsToRole(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	var req dto.AssignPermissionsRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}
	if err := h.service.AssignPermissionsToRole(id, &req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Permission berhasil ditambahkan ke role", nil, nil)
}

func (h *RBACHandler) SyncRolePermissions(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	var req dto.AssignPermissionsRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if err := h.service.SyncRolePermissions(id, &req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Permission role berhasil disinkronkan", nil, nil)
}

func (h *RBACHandler) RevokePermissionsFromRole(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	var req dto.AssignPermissionsRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if err := h.service.RevokePermissionsFromRole(id, &req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Permission berhasil dicabut dari role", nil, nil)
}

// ─── User ↔ Role Handlers ──────────────────────────────────────────────────────

func (h *RBACHandler) GetUserRoles(c *echo.Context) error {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "User ID tidak valid", nil, nil)
	}
	roles, err := h.service.GetUserRoles(userID)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil role user", nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil role user", roles, nil)
}

func (h *RBACHandler) AssignRolesToUser(c *echo.Context) error {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "User ID tidak valid", nil, nil)
	}
	var req dto.AssignRolesRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}
	if err := h.service.AssignRolesToUser(userID, &req, getActorID(c)); err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Role berhasil ditambahkan ke user", nil, nil)
}

func (h *RBACHandler) SyncUserRoles(c *echo.Context) error {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "User ID tidak valid", nil, nil)
	}
	var req dto.AssignRolesRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if err := h.service.SyncUserRoles(userID, &req, getActorID(c)); err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Role user berhasil disinkronkan", nil, nil)
}

func (h *RBACHandler) RevokeRolesFromUser(c *echo.Context) error {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "User ID tidak valid", nil, nil)
	}
	var req dto.AssignRolesRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if err := h.service.RevokeRolesFromUser(userID, &req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Role berhasil dicabut dari user", nil, nil)
}

// ─── User Permissions ──────────────────────────────────────────────────────────

func (h *RBACHandler) GetUserAllPermissions(c *echo.Context) error {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "User ID tidak valid", nil, nil)
	}
	perms, err := h.service.GetUserAllPermissions(userID)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil permission", nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil semua permission user", perms, nil)
}

func (h *RBACHandler) AssignDirectPermission(c *echo.Context) error {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "User ID tidak valid", nil, nil)
	}
	var req dto.AssignDirectPermissionRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}
	if err := h.service.AssignDirectPermission(userID, &req, getActorID(c)); err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Direct permission berhasil ditetapkan", nil, nil)
}
