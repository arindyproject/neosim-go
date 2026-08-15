package handlers

import (
	"io"
	"net/http"

	"neosim_go/internal/modules/artikel/kategori/dto"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"
	"neosim_go/internal/shared/binding"
	he "neosim_go/internal/shared/httputil"

	"github.com/labstack/echo/v5"
)

// ─── ListKategori ─────────────────────────────────────────────────────
//
//	@Summary		Get list of ArtikelKategori
//	@Description	Get paginated list of ArtikelKategori
//	@Tags			artikel/kategori
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.ArtikelKategoriResponse}
//	@Router			/artikel/kategori [get]
func (h *ArtikelKategoriHandler) ListKategori(c *echo.Context) error {

	filter := dto.FilterArtikelKategoriRequest{
		Name: c.QueryParam("name"),
	}
	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.ListKategori(page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetKategoriByID ───────────────────────────────────────────────────
//
//	@Summary		Get ArtikelKategori
//	@Description	Get ArtikelKategori by :id
//	@Tags			artikel/kategori
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"ArtikelKategori ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.ArtikelKategoriResponse}
//	@Router			/artikel/kategori/{id} [get]
func (h *ArtikelKategoriHandler) GetKategoriByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.GetKategoriByID(id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── CreateKategori ────────────────────────────────────────────────────
//
//	@Summary		Create ArtikelKategori
//	@Description	Create New ArtikelKategori
//	@Tags			artikel/kategori
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateArtikelKategoriRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.ArtikelKategoriResponse}
//	@Router			/artikel/kategori [post]
func (h *ArtikelKategoriHandler) CreateKategori(c *echo.Context) error {
	var req dto.CreateArtikelKategoriRequest
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Gagal membaca request body", nil, err.Error())
	}

	if errs := binding.BindErrors(body, &req); len(errs) > 0 {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal (binding)", nil, errs)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal (validator)", nil, errs)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.CreateKategori(&req,  actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── UpdateKategori ────────────────────────────────────────────────────
//
//	@Summary		Update ArtikelKategori
//	@Description	Update ArtikelKategori by :id
//	@Tags			artikel/kategori
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"ArtikelKategori ID"
//	@Param			body	body		dto.UpdateArtikelKategoriRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.ArtikelKategoriResponse}
//	@Router			/artikel/kategori/{id} [put]
func (h *ArtikelKategoriHandler) UpdateKategori(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	var req dto.UpdateArtikelKategoriRequest
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Gagal membaca request body", nil, err.Error())
	}

	if errs := binding.BindErrors(body, &req); len(errs) > 0 {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal (binding)", nil, errs)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal (validator)", nil, errs)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.UpdateKategori(id, &req, actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "ArtikelKategori tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── DeleteKategori ────────────────────────────────────────────────────
//
//	@Summary		Delete ArtikelKategori
//	@Description	Delete ArtikelKategori by :id
//	@Tags			artikel/kategori
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"ArtikelKategori ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/artikel/kategori/{id} [delete]
func (h *ArtikelKategoriHandler) DeleteKategori(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.DeleteKategori(id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "ArtikelKategori tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
