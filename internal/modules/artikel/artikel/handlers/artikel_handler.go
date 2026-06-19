package handlers

import (
	"net/http"
	"strconv"

	"neosim_go/internal/modules/artikel/artikel/dto"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"

	he "neosim_go/internal/shared/httputil"

	"github.com/labstack/echo/v5"
)

// ─── Handlers ──────────────────────────────────────────────────────────────────

// ─── List ──────────────────────────────────────────────────────────────────────
// ArtikelHandler godoc
//
//	@Summary		Get list of Artikel
//	@Description	Get paginated list of Artikel
//	@Tags			artikel/artikel
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name			query		string	false	"Filter by name (partial match)"
//	@Param			page			query		int		false	"Page number"
//	@Param			page_size		query		int		false	"Page size"
//	@Success		200				{object}	response.MyGoResponse{data=[]dto.ArtikelResponse}
//	@Router			/api/v1/artikel [get]
//
// List handles GET /api/v1/artikel
func (h *ArtikelHandler) List(c *echo.Context) error {
	page, pageSize := 1, 10

	// Mengambil query parameter untuk filter
	filter := dto.FilterArtikelRequest{
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

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.List(page, pageSize,&filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}

	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetByID ───────────────────────────────────────────────────────────────────
// ArtikelHandler godoc
//
//	@Summary		Get Artikel
//	@Description	Get Artikel by :id
//	@Tags			artikel/artikel
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Artikel ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.ArtikelResponse}
//	@Router			/api/v1/artikel/{id} [get]
//
// GetByID handles GET /api/v1/artikel/:id
func (h *ArtikelHandler) GetByID(c *echo.Context) error {
	id, err :=  he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.GetByID(id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── Create ────────────────────────────────────────────────────────────────────
// ArtikelHandler godoc
//
//	@Summary		Create Artikel
//	@Description	Create New Artikel
//	@Tags			artikel/artikel
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateArtikelRequest	true	"Create Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.ArtikelResponse}
//	@Router			/api/v1/artikel [post]
//
// Create handles POST /api/v1/artikel
func (h *ArtikelHandler) Create(c *echo.Context) error {
	var req dto.CreateArtikelRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}

	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.Create(&req, getActorID(c), actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}


// ─── Update ────────────────────────────────────────────────────────────────────
// ArtikelHandler godoc
//
//	@Summary		Update Artikel
//	@Description	Update Artikel by :id
//	@Tags			artikel/artikel
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"Artikel ID"
//	@Param			body	body		dto.UpdateArtikelRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.ArtikelResponse}
//	@Router			/api/v1/artikel/{id} [put]
//
// Update handles PUT /api/v1/artikel/:id
func (h *ArtikelHandler) Update(c *echo.Context) error {
	id, err :=  he.ParseID(c)
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

	actor := he.BuildAuthContext(c)
	item, err := h.service.Update(id, &req, getActorID(c), actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "Artikel tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── Delete ────────────────────────────────────────────────────────────────────
// ArtikelHandler godoc
//
//	@Summary		Delete Artikel
//	@Description	Delete Artikel by :id
//	@Tags			artikel/artikel
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"Artikel ID"
//	@Success		200		{object}	response.MyGoResponse{}
//	@Router			/api/v1/artikel/{id} [delete]
//
// Delete handles DELETE /api/v1/artikel/:id
func (h *ArtikelHandler) Delete(c *echo.Context) error {
	id, err :=  he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := he.BuildAuthContext(c)
	if err := h.service.Delete(id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "Artikel tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
