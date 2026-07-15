package handlers

import (
	"net/http"

	"neosim_go/internal/modules/artikel/kategori/dto"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"

	"github.com/labstack/echo/v5"
)

// Method di bawah ini ditempelkan ke struct ArtikelKategoriHandler yang
// sama dengan handler entitas utama (lihat handlers/handler.go). Nama method
// diberi suffix Tag agar tidak bentrok dengan List/GetByID/dst milik
// entitas utama pada struct handler yang sama.

// ─── ListTag ──────────────────────────────────────────────────────
//
//	@Summary		Get list of Tag
//	@Description	Get paginated list of Tag
//	@Tags			artikel/kategori
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.TagResponse}
//	@Router			/artikel/kategori/tags [get]
func (h *ArtikelKategoriHandler) ListTag(c *echo.Context) error {
	filter := dto.FilterTagRequest{Name: c.QueryParam("name")}
	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.ListTag(page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetTagByID ───────────────────────────────────────────────────
//
//	@Summary		Get Tag
//	@Description	Get Tag by :id
//	@Tags			artikel/kategori
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Tag ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.TagResponse}
//	@Router			/artikel/kategori/tags/{id} [get]
func (h *ArtikelKategoriHandler) GetTagByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.GetTagByID(id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── CreateTag ────────────────────────────────────────────────────
//
//	@Summary		Create Tag
//	@Description	Create New Tag
//	@Tags			artikel/kategori
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateTagRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.TagResponse}
//	@Router			/artikel/kategori/tags [post]
func (h *ArtikelKategoriHandler) CreateTag(c *echo.Context) error {
	var req dto.CreateTagRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.CreateTag(&req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── UpdateTag ────────────────────────────────────────────────────
//
//	@Summary		Update Tag
//	@Description	Update Tag by :id
//	@Tags			artikel/kategori
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"Tag ID"
//	@Param			body	body		dto.UpdateTagRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.TagResponse}
//	@Router			/artikel/kategori/tags/{id} [put]
func (h *ArtikelKategoriHandler) UpdateTag(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	var req dto.UpdateTagRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.UpdateTag(id, &req, actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "Tag tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── DeleteTag ────────────────────────────────────────────────────
//
//	@Summary		Delete Tag
//	@Description	Delete Tag by :id
//	@Tags			artikel/kategori
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Tag ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/artikel/kategori/tags/{id} [delete]
func (h *ArtikelKategoriHandler) DeleteTag(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.DeleteTag(id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "Tag tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
