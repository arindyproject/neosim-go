package handlers

import (
	"net/http"
	"strconv"

	"neosim_go/internal/modules/artikel/contracts"
	"neosim_go/internal/modules/artikel/dto"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"

	rbacMiddlewares "neosim_go/internal/modules/rbac/middlewares"

	"github.com/labstack/echo/v5"
)

// ArtikelHandler defines HTTP handlers
type ArtikelHandler struct {
	service contracts.Service
}

// NewArtikelHandler membuat instance handler baru
func NewArtikelHandler(service contracts.Service) *ArtikelHandler {
	return &ArtikelHandler{service: service}
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
// ArtikelHandler godoc
//
//	@Summary		Get list of users
//	@Description	Get paginated list of users
//	@Tags			artikel
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page			query		int		false	"Page number"
//	@Param			page_size		query		int		false	"Page size"
//	@Success		200				{object}	response.MyGoResponse{data=[]dto.ArtikelResponse}
//	@Router			/artikels [get]
//
// List handles GET /api/v1/artikels
func (h *ArtikelHandler) List(c *echo.Context) error {
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

	actor := buildAuthContext(c)
	items, total, err := h.service.List(page, pageSize, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}

	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetByID ───────────────────────────────────────────────────────────────────
// ArtikelHandler godoc
//
//	@Summary		Get artikel
//	@Description	Get artikel by :id
//	@Tags			artikel
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"artikel ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.ArtikelResponse}
//	@Router			/artikels/{id} [get]
//
// GetByID handles GET /api/v1/artikels/:id
func (h *ArtikelHandler) GetByID(c *echo.Context) error {
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
// ArtikelHandler godoc
//
//	@Summary		Create artikel
//	@Description	Create New artikel
//	@Tags			artikel
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateArtikelRequest	true	"Login Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.ArtikelResponse}
//	@Router			/artikels [post]
//
// Create handles POST /api/v1/artikels
func (h *ArtikelHandler) Create(c *echo.Context) error {
	var req dto.CreateArtikelRequest
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
// ArtikelHandler godoc
//
//	@Summary		Update artikel
//	@Description	Update artikel by :id
//	@Tags			artikel
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"artikel ID"
//	@Param			body	body		dto.UpdateArtikelRequest	true	"Login Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.ArtikelResponse}
//	@Router			/artikels/{id} [put]
//
// Update handles PUT /api/v1/artikels/:id
func (h *ArtikelHandler) Update(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateArtikelRequest
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
		if err.Error() == "artikel tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── Delete ────────────────────────────────────────────────────────────────────
// ArtikelHandler godoc
//
//	@Summary		Update artikel
//	@Description	Update artikel by :id
//	@Tags			artikel
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"artikel ID"
//	@Success		200		{object}	response.MyGoResponse{}
//	@Router			/artikels/{id} [delete]
//
// Delete handles DELETE /api/v1/artikels/:id
func (h *ArtikelHandler) Delete(c *echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := buildAuthContext(c)
	if err := h.service.Delete(id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "artikel tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
