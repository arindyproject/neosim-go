package handlers

import (
	"neosim_go/internal/modules/master/master/dto"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"
	"net/http"

	"github.com/labstack/echo/v5"
)

// ─────────────── List ─────────────────────────────────────────────────────────────
//
//	 MasterHandler godoc
//
//		@Summary		Get list of Pekerjaan
//		@Description	Get paginated list of Pekerjaan
//		@Tags			master/pekerjaan
//		@Accept			json
//		@Produce		json
//		@Security		BearerAuth
//		@Param			kode_kemenkes	query		string	false	"Filter by kode_kemenkes"
//		@Param			name			query		string	false	"Filter by name (partial match)"
//		@Param			page			query		int		false	"Page number"
//		@Param			page_size		query		int		false	"Page size"
//		@Success		200				{object}	response.MyGoResponse{data=[]dto.MasterPekerjaanResponse}
//		@Router			/master/pekerjaan [get]
//
// ListPekerjaan handles GET /api/v1/master/pekerjaan
func (h *MasterHandler) ListPekerjaan(c *echo.Context) error {
	page, pageSize := he.ParsePagination(c, h.cfg)
	filter := dto.FilterMasterPekerjaanRequest{
		KodeKemenkes: c.QueryParam("kode_kemenkes"),
		Name:         c.QueryParam("name"),
	}

	items, total, err := h.service.ListPekerjaan(c.Request().Context(), page, pageSize, &filter)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}

	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
} // ───────────── List ───────────────────────────────────────────────────────────

// ─────────────── GetByID ────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary		Get Pekerjaan
//	@Description	Get Pekerjaan detail by :id
//	@Tags			master/pekerjaan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Pekerjaan ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.MasterPekerjaanResponse}
//	@Router			/master/pekerjaan/{id}	[get]
//
// GetByIDPekerjaan handles GET /api/v1/master/pekerjaan/:id
func (h *MasterHandler) GetByIDPekerjaan(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	item, err := h.service.GetByIDPekerjaan(c.Request().Context(), id)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── Create ─────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary	Create Pekerjaan
//	@Description	Create new Pekerjaan
//	@Tags		master/pekerjaan
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		body	body	dto.CreateMasterPekerjaanRequest	true	"Pekerjaan data"
//	@Success	201	{object}	response.MyGoResponse{data=dto.MasterPekerjaanResponse}
//	@Router		/master/pekerjaan				[post]
//
// CreatePekerjaan handles POST /api/v1/master/pekerjaan
func (h *MasterHandler) CreatePekerjaan(c *echo.Context) error {
	var req dto.CreateMasterPekerjaanRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.CreatePekerjaan(c.Request().Context(), &req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
} // ───────────── Create ─────────────────────────────────────────────────────────

// ─────────────── Update ─────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary	Update Pekerjaan
//	@Description	Update Pekerjaan by :id
//	@Tags	master/pekerjaan
//	@Accept	json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		id		path		int			true	"Pekerjaan ID"
//	@Param	body	body	dto.UpdateMasterPekerjaanRequest	true	"Update Request"
//	@Success	200	{object}	response.MyGoResponse{data=dto.MasterPekerjaanResponse}
//	@Router		/master/pekerjaan/{id}	[put]
//
// UpdatePekerjaan handles PUT /api/v1/master/pekerjaan/:id
func (h *MasterHandler) UpdatePekerjaan(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateMasterPekerjaanRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.UpdatePekerjaan(c.Request().Context(), id, &req, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
} // ───────────── Update ─────────────────────────────────────────────────────────

// ─────────────── Delete ─────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary	Delete Pekerjaan
//	@Description	Delete Pekerjaan by :id
//	@Tags	master/pekerjaan
//	@Accept	json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		id	path	int	true	"Pekerjaan ID"
//	@Success	200	{object}	response.MyGoResponse
//	@Router	/master/pekerjaan/{id}	[delete]
//
// DeletePekerjaan handles DELETE /api/v1/master/pekerjaan/:id
func (h *MasterHandler) DeletePekerjaan(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := he.BuildAuthContext(c)
	err = h.service.DeletePekerjaan(c.Request().Context(), id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
} // ───────────── Delete ─────────────────────────────────────────────────────────
