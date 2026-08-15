package handlers

import (
	"io"
	"net/http"

	"neosim_go/internal/modules/artikel/artikel/dto"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"
	"neosim_go/internal/shared/binding"
	he "neosim_go/internal/shared/httputil"

	"github.com/labstack/echo/v5"
)

// ─── ListArtikel ─────────────────────────────────────────────────────
//
//	@Summary		Get list of Artikel
//	@Description	Get paginated list of Artikel
//	@Tags			artikel/artikel
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.ArtikelResponse}
//	@Router			/artikel [get]
func (h *ArtikelHandler) ListArtikel(c *echo.Context) error {

	filter := dto.FilterArtikelRequest{
		Name: c.QueryParam("name"),
	}
	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.ListArtikel(page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetArtikelByID ───────────────────────────────────────────────────
//
//	@Summary		Get Artikel
//	@Description	Get Artikel by :id
//	@Tags			artikel/artikel
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Artikel ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.ArtikelResponse}
//	@Router			/artikel/{id} [get]
func (h *ArtikelHandler) GetArtikelByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.GetArtikelByID(id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── CreateArtikel ────────────────────────────────────────────────────
//
//	@Summary		Create Artikel
//	@Description	Create New Artikel
//	@Tags			artikel/artikel
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateArtikelRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.ArtikelResponse}
//	@Router			/artikel [post]
func (h *ArtikelHandler) CreateArtikel(c *echo.Context) error {
	var req dto.CreateArtikelRequest
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
	item, err := h.service.CreateArtikel(&req,  actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── UpdateArtikel ────────────────────────────────────────────────────
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
//	@Router			/artikel/{id} [put]
func (h *ArtikelHandler) UpdateArtikel(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	var req dto.UpdateArtikelRequest
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
	item, err := h.service.UpdateArtikel(id, &req, actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "Artikel tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── DeleteArtikel ────────────────────────────────────────────────────
//
//	@Summary		Delete Artikel
//	@Description	Delete Artikel by :id
//	@Tags			artikel/artikel
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Artikel ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/artikel/{id} [delete]
func (h *ArtikelHandler) DeleteArtikel(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.DeleteArtikel(id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "Artikel tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
