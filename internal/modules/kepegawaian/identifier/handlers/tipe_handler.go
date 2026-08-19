package handlers

import (
	"net/http"
	"strconv"

	"neosim_go/internal/modules/kepegawaian/identifier/dto"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"

	"github.com/labstack/echo/v5"
)

// Method di bawah ini ditempelkan ke struct KepegawaianIdentifierHandler yang
// sama dengan handler entitas utama (lihat handlers/handler.go). Nama method
// diberi suffix Tipe agar tidak bentrok dengan method entitas utama
// pada struct handler yang sama.

// ─── ListTipe ──────────────────────────────────────────────────────
//
//	@Summary		Get list of Tipe
//	@Description	Get paginated list of Tipe with optional filters
//	@Tags			kepegawaian/identifier/tipe
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			search		query		string	false	"Search by label or code"
//	@Param			code		query		string	false	"Filter by code"
//	@Param			label		query		string	false	"Filter by label"
//	@Param			is_nakes	query		bool	false	"Filter by is_nakes"
//	@Param			is_required	query		bool	false	"Filter by is_required"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.TipeResponse}
//	@Router			/kepegawaian/identifier/tipes [get]
func (h *KepegawaianIdentifierHandler) ListTipe(c *echo.Context) error {
	filter := dto.FilterTipeRequest{
		Search: c.QueryParam("search"),
		Code:   c.QueryParam("code"),
		Label:  c.QueryParam("label"),
	}

	if isNakesStr := c.QueryParam("is_nakes"); isNakesStr != "" {
		if val, err := strconv.ParseBool(isNakesStr); err == nil {
			filter.IsNakes = &val
		}
	}

	if isReqStr := c.QueryParam("is_required"); isReqStr != "" {
		if val, err := strconv.ParseBool(isReqStr); err == nil {
			filter.IsRequired = &val
		}
	}

	page, pageSize := he.ParsePagination(c, h.cfg)

	actor := he.BuildAuthContext(c)
	items, total, err := h.service.ListTipe(c.Request().Context(), page, pageSize, &filter, actor)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}
	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
}

// ─── GetTipeByID ───────────────────────────────────────────────────
//
//	@Summary		Get Tipe
//	@Description	Get Tipe by :id
//	@Tags			kepegawaian/identifier/tipe
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Tipe ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.TipeResponse}
//	@Router			/kepegawaian/identifier/tipes/{id} [get]
func (h *KepegawaianIdentifierHandler) GetTipeByID(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.GetTipeByID(c.Request().Context(), id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
}

// ─── CreateTipe ────────────────────────────────────────────────────
//
//	@Summary		Create Tipe
//	@Description	Create New Tipe
//	@Tags			kepegawaian/identifier/tipe
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateTipeRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.TipeResponse}
//	@Router			/kepegawaian/identifier/tipes [post]
func (h *KepegawaianIdentifierHandler) CreateTipe(c *echo.Context) error {
	var req dto.CreateTipeRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.CreateTipe(c.Request().Context(), &req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
}

// ─── UpdateTipe ────────────────────────────────────────────────────
//
//	@Summary		Update Tipe
//	@Description	Update Tipe by :id
//	@Tags			kepegawaian/identifier/tipe
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int						true	"Tipe ID"
//	@Param			body	body		dto.UpdateTipeRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.TipeResponse}
//	@Router			/kepegawaian/identifier/tipes/{id} [put]
func (h *KepegawaianIdentifierHandler) UpdateTipe(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	var req dto.UpdateTipeRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}
	actor := he.BuildAuthContext(c)
	item, err := h.service.UpdateTipe(c.Request().Context(), id, &req, actor)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "Tipe tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
}

// ─── DeleteTipe ────────────────────────────────────────────────────
//
//	@Summary		Delete Tipe
//	@Description	Delete Tipe by :id
//	@Tags			kepegawaian/identifier/tipe
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Tipe ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/kepegawaian/identifier/tipes/{id} [delete]
func (h *KepegawaianIdentifierHandler) DeleteTipe(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}
	actor := he.BuildAuthContext(c)
	if err := h.service.DeleteTipe(c.Request().Context(), id, actor); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "Tipe tidak ditemukan" {
			status = http.StatusNotFound
		}
		return response.Response(c, status, false, err.Error(), nil, nil)
	}
	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
}
