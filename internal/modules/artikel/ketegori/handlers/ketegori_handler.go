package handlers

import (
	"net/http"
	"strconv"

	"neosim_go/internal/modules/artikel/ketegori/dto"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"

	he "neosim_go/internal/shared/httputil"

	"github.com/labstack/echo/v5"
)

// ─── Handlers ──────────────────────────────────────────────────────────────────

// ─── List ──────────────────────────────────────────────────────────────────────
// ArtikelKetegoriHandler godoc
//
//	@Summary		Get list of ArtikelKetegori
//	@Description	Get paginated list of ArtikelKetegori
//	@Tags			artikel/ketegori
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name			query		string	false	"Filter by name (partial match)"
//	@Param			page			query		int		false	"Page number"
//	@Param			page_size		query		int		false	"Page size"
//	@Success		200				{object}	response.MyGoResponse{data=[]dto.ArtikelKetegoriResponse}
//	@Router			/api/v1/artikel/ketegori [get]
//
// List handles GET /api/v1/artikel/ketegori
func (h *ArtikelKetegoriHandler) List(c *echo.Context) error {
	page, pageSize := 1, 10

	// Mengambil query parameter untuk filter
	filter := dto.FilterArtikelKetegoriRequest{
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
// ArtikelKetegoriHandler godoc
//
//	@Summary		Get ArtikelKetegori
//	@Description	Get ArtikelKetegori by :id
//	@Tags			artikel/ketegori
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"ArtikelKetegori ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.ArtikelKetegoriResponse}
//	@Router			/api/v1/artikel/ketegori/{id} [get]
//
// GetByID handles GET /api/v1/artikel/ketegori/:id
func (h *ArtikelKetegoriHandler) GetByID(c *echo.Context) error {
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
// ArtikelKetegoriHandler godoc
//
//	@Summary		Create ArtikelKetegori
//	@Description	Create New ArtikelKetegori
//	@Tags			artikel/ketegori
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateArtikelKetegoriRequest	true	"Create Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.ArtikelKetegoriResponse}
//	@Router			/api/v1/artikel/ketegori [post]
//
// Create handles POST /api/v1/artikel/ketegori
func (h *ArtikelKetegoriHandler) Create(c *echo.Context) error {
	var req dto.CreateArtikelKetegoriRequest
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
// ArtikelKetegoriHandler godoc
//
//	@Summary		Update ArtikelKetegori
//	@Description	Update ArtikelKetegori by :id
//	@Tags			artikel/ketegori
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"ArtikelKetegori ID"
//	@Param			body	body		dto.UpdateArtikelKetegoriRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.ArtikelKetegoriResponse}
//	@Router			/api/v1/artikel/ketegori/{id} [put]
//
// Update handles PUT /api/v1/artikel/ketegori/:id
func (h *ArtikelKetegoriHandler) Update(c *echo.Context) error {
	id, err :=  he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateArtikelKetegoriRequest
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
		if err.Error() == "ArtikelKetegori tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── Delete ────────────────────────────────────────────────────────────────────
// ArtikelKetegoriHandler godoc
//
//	@Summary		Delete ArtikelKetegori
//	@Description	Delete ArtikelKetegori by :id
//	@Tags			artikel/ketegori
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"ArtikelKetegori ID"
//	@Success		200		{object}	response.MyGoResponse{}
//	@Router			/api/v1/artikel/ketegori/{id} [delete]
//
// Delete handles DELETE /api/v1/artikel/ketegori/:id
func (h *ArtikelKetegoriHandler) Delete(c *echo.Context) error {
	id, err :=  he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := he.BuildAuthContext(c)
	if err := h.service.Delete(id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "ArtikelKetegori tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
