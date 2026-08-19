package handlers

import (
	"io"
	"neosim_go/internal/modules/master/alamat/dto"
	"neosim_go/internal/shared/binding"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"
	"net/http"

	"github.com/labstack/echo/v5"
)

// Negara ===========================================================================
// ─────────────── List ─────────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Get list of Negara
//	@Description	Get paginated list of Negara
//	@Tags			master/alamat/negara
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			code		query		string	false	"Filter by code"
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.NegaraResponse}
//	@Router			/master/alamat/negara [get]
//
// ListNegara handles GET /api/v1/master/alamat/alamat/negara
func (h *MasterAlamatHandler) ListNegara(c *echo.Context) error {
	page, pageSize := he.ParsePagination(c, h.cfg)
	filter := dto.FilterNegaraRequest{
		Code: c.QueryParam("code"),
		Name: c.QueryParam("name"),
	}

	items, total, err := h.service.ListNegara(c.Request().Context(), page, pageSize, &filter)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}

	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
} // ───────────── List ───────────────────────────────────────────────────────────

// ─────────────── GetByID ────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Get Negara
//	@Description	Get Negara by :id
//	@Tags			master/alamat/negara
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Negara ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.NegaraResponse}
//	@Router			/master/alamat/negara/{id} [get]
//
// GetByIDNegara handles GET /api/v1/master/alamat/alamat/negara/:id
func (h *MasterAlamatHandler) GetByIDNegara(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	item, err := h.service.GetByIDNegara(c.Request().Context(), id)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── Create ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Create Negara
//	@Description	Create New Negara
//	@Tags			master/alamat/negara
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateNegaraRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.NegaraResponse}
//	@Router			/master/alamat/negara [post]
//
// CreateNegara handles POST /api/v1/master/alamat/alamat/negara
func (h *MasterAlamatHandler) CreateNegara(c *echo.Context) error {
	var req dto.CreateNegaraRequest
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Gagal membaca request body", nil, err.Error())
	}

	if errs := binding.BindErrors(body, &req); len(errs) > 0 {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.CreateNegara(c.Request().Context(), &req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
} // ───────────── Create ─────────────────────────────────────────────────────────

// ─────────────── Update ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Update Negara
//	@Description	Update Negara by :id
//	@Tags			master/alamat/negara
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"Negara ID"
//	@Param			body	body		dto.UpdateNegaraRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.NegaraResponse}
//	@Router			/master/alamat/negara/{id} [put]
//
// UpdateNegara handles PUT /api/v1/master/alamat/alamat/negara/:id
func (h *MasterAlamatHandler) UpdateNegara(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateNegaraRequest
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Gagal membaca request body", nil, err.Error())
	}

	if errs := binding.BindErrors(body, &req); len(errs) > 0 {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.UpdateNegara(c.Request().Context(), id, &req, actor)
	if err != nil {
		return response.Response(c, he.NotFoundStatus(err, "Negara tidak ditemukan"), false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
} // ───────────── Update ─────────────────────────────────────────────────────────

// ─────────────── Delete ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Delete Negara
//	@Description	Delete Negara by :id
//	@Tags			master/alamat/negara
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Negara ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/master/alamat/negara/{id} [delete]
//
// DeleteNegara handles DELETE /api/v1/master/alamat/alamat/negara/:id
func (h *MasterAlamatHandler) DeleteNegara(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := he.BuildAuthContext(c)
	if err := h.service.DeleteNegara(c.Request().Context(), id, actor); err != nil {
		return response.Response(c, he.NotFoundStatus(err, "Negara tidak ditemukan"), false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
} // ───────────── Delete ─────────────────────────────────────────────────────────
// Negara ===========================================================================
