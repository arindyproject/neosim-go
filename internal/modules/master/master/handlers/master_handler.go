package handlers

import (
	"net/http"
	"strconv"

	"neosim_go/internal/modules/master/master/contracts"
	"neosim_go/internal/modules/master/master/dto"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"

	rbacMiddlewares "neosim_go/internal/modules/rbac/middlewares"

	"github.com/labstack/echo/v5"
)

// MasterHandler defines HTTP handlers
type MasterHandler struct {
	service contracts.Service
}

// NewMasterHandler membuat instance handler baru
func NewMasterHandler(service contracts.Service) *MasterHandler {
	return &MasterHandler{service: service}
}

// buildAuthContext membuat AuthContext dari JWT claims di context
func buildAuthContext(c *echo.Context) contracts.AuthContext {
	userID, _ := rbacMiddlewares.GetUserIDFromContext(c)
	isSuperadmin := rbacMiddlewares.IsSuperadmin(c)
	return contracts.AuthContext{
		UserID:       userID,
		IsSuperadmin: isSuperadmin,
	}
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

// ─── Handlers ──────────────────────────────────────────────────────────────────

// ─── List ──────────────────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary		Get list of Master
//	@Description	Get paginated list of Master
//	@Tags			master/master
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name			query		string	false	"Filter by name (partial match)"
//	@Param			page			query		int		false	"Page number"
//	@Param			page_size		query		int		false	"Page size"
//	@Success		200				{object}	response.MyGoResponse{data=[]dto.MasterResponse}
//	@Router			/api/v1/master [get]
//
// List handles GET /api/v1/master
func (h *MasterHandler) List(c *echo.Context) error {
	page, pageSize := 1, 10

	// Mengambil query parameter untuk filter
	filter := dto.FilterMasterRequest{
		Name:     c.QueryParam("name"),
	}

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

	actor := buildAuthContext(c)
	items, total, err := h.service.List(page, pageSize,&filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}

	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetByID ───────────────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary		Get Master
//	@Description	Get Master by :id
//	@Tags			master/master
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Master ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.MasterResponse}
//	@Router			/api/v1/master/{id} [get]
//
// GetByID handles GET /api/v1/master/:id
func (h *MasterHandler) GetByID(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := buildAuthContext(c)
	item, err := h.service.GetByID(id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── Create ────────────────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary		Create Master
//	@Description	Create New Master
//	@Tags			master/master
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateMasterRequest	true	"Create Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.MasterResponse}
//	@Router			/api/v1/master [post]
//
// Create handles POST /api/v1/master
func (h *MasterHandler) Create(c *echo.Context) error {
	var req dto.CreateMasterRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := buildAuthContext(c)
	item, err := h.service.Create(&req, getActorID(c), actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}


// ─── Update ────────────────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary		Update Master
//	@Description	Update Master by :id
//	@Tags			master/master
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"Master ID"
//	@Param			body	body		dto.UpdateMasterRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.MasterResponse}
//	@Router			/api/v1/master/{id} [put]
//
// Update handles PUT /api/v1/master/:id
func (h *MasterHandler) Update(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateMasterRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := buildAuthContext(c)
	item, err := h.service.Update(id, &req, getActorID(c), actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "Master tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── Delete ────────────────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary		Delete Master
//	@Description	Delete Master by :id
//	@Tags			master/master
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"Master ID"
//	@Success		200		{object}	response.MyGoResponse{}
//	@Router			/api/v1/master/{id} [delete]
//
// Delete handles DELETE /api/v1/master/:id
func (h *MasterHandler) Delete(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := buildAuthContext(c)
	if err := h.service.Delete(id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "Master tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
