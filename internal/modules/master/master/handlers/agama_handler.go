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
//		@Summary		Get list of Agama
//		@Description	Get paginated list of Agama
//		@Tags			master/agama
//		@Accept			json
//		@Produce		json
//		@Security		BearerAuth
//		@Param			kode_kemenkes	query		string	false	"Filter by kode_kemenkes"
//		@Param			name			query		string	false	"Filter by name (partial match)"
//		@Param			page			query		int		false	"Page number"
//		@Param			page_size		query		int		false	"Page size"
//		@Success		200				{object}	response.MyGoResponse{data=[]dto.MasterAgamaResponse}
//		@Router			/master/agama [get]
//
// ListAgama handles GET /api/v1/master/Agama
func (h *MasterHandler) ListAgama(c *echo.Context) error {
	page, pageSize := he.ParsePagination(c, h.cfg)
	filter := dto.FilterMasterAgamaRequest{
		KodeKemenkes: c.QueryParam("kode_kemenkes"),
		Name:         c.QueryParam("name"),
	}

	items, total, err := h.service.ListAgama(c.Request().Context(), page, pageSize, &filter)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}

	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
} // ───────────── List ───────────────────────────────────────────────────────────

// ─────────────── GetByID ────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary		Get Agama
//	@Description	Get Agama detail by :id
//	@Tags			master/agama
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Agama ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.MasterAgamaResponse}
//	@Router			/master/agama/{id}	[get]
//
// GetByIDAgama handles GET /api/v1/master/agama/:id
func (h *MasterHandler) GetByIDAgama(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	item, err := h.service.GetByIDAgama(c.Request().Context(), id)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── Create ─────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary	Create Agama
//	@Description	Create new Agama
//	@Tags		master/agama
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		body 		body		dto.CreateMasterAgamaRequest	true	"Agama Request"
//	@Success	201	{object}	response.MyGoResponse{data=dto.MasterAgamaResponse}
//	@Router		/master/agama				[post]
//
// CreateAgama handles POST /api/v1/master/agama
func (h *MasterHandler) CreateAgama(c *echo.Context) error {
	var req dto.CreateMasterAgamaRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.CreateAgama(c.Request().Context(), &req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
} // ───────────── Create ─────────────────────────────────────────────────────────

// ─────────────── Update ─────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary	Update Agama
//	@Description	Update Agama by :id
//	@Tags	master/agama
//	@Accept	json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		id		path		int			true	"Agama ID"
//	@Param	body	body	dto.UpdateMasterAgamaRequest	true	"Update Request"
//	@Success	200	{object}	response.MyGoResponse{data=dto.MasterAgamaResponse}
//	@Router		/master/agama/{id}	[put]
//
// UpdateAgama handles PUT /api/v1/master/agama/:id
func (h *MasterHandler) UpdateAgama(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateMasterAgamaRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.UpdateAgama(c.Request().Context(), id, &req, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
} // ───────────── Update ─────────────────────────────────────────────────────────

// ─────────────── Delete ─────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary	Delete Agama
//	@Description	Delete Agama by :id
//	@Tags	master/agama
//	@Accept	json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		id	path	int	true	"Agama ID"
//	@Success	200	{object}	response.MyGoResponse
//	@Router	/master/agama/{id}	[delete]
//
// DeleteAgama handles DELETE /api/v1/master/Agama/:id
func (h *MasterHandler) DeleteAgama(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := he.BuildAuthContext(c)
	err = h.service.DeleteAgama(c.Request().Context(), id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
} // ───────────── Delete ─────────────────────────────────────────────────────────
