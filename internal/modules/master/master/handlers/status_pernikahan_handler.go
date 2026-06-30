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
//		@Summary		Get list of StatusPernikahan
//		@Description	Get paginated list of StatusPernikahan
//		@Tags			master/status_pernikahan
//		@Accept			json
//		@Produce		json
//		@Security		BearerAuth
//		@Param			kode_kemenkes	query		string	false	"Filter by kode_kemenkes"
//		@Param			name			query		string	false	"Filter by name (partial match)"
//		@Param			page			query		int		false	"Page number"
//		@Param			page_size		query		int		false	"Page size"
//		@Success		200				{object}	response.MyGoResponse{data=[]dto.MasterStatusPernikahanResponse}
//		@Router			/master/status_pernikahan [get]
//
// ListStatusPernikahan handles GET /api/v1/master/status_pernikahan
func (h *MasterHandler) ListStatusPernikahan(c *echo.Context) error {
	page, pageSize := he.ParsePagination(c, h.cfg)
	filter := dto.FilterMasterStatusPernikahanRequest{
		KodeKemenkes: c.QueryParam("kode_kemenkes"),
		Name:         c.QueryParam("name"),
	}

	items, total, err := h.service.ListStatusPernikahan(page, pageSize, &filter)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}

	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
} // ───────────── List ───────────────────────────────────────────────────────────

// ─────────────── GetByID ────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary		Get StatusPernikahan
//	@Description	Get StatusPernikahan detail by :id
//	@Tags			master/status_pernikahan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"StatusPernikahan ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.MasterStatusPernikahanResponse}
//	@Router			/master/status_pernikahan/{id}	[get]
//
// GetByIDStatusPernikahan handles GET /api/v1/master/status_pernikahan/:id
func (h *MasterHandler) GetByIDStatusPernikahan(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	item, err := h.service.GetByIDStatusPernikahan(id)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── Create ─────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary	Create StatusPernikahan
//	@Description	Create new StatusPernikahan
//	@Tags		master/status_pernikahan
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		body	body	dto.CreateMasterStatusPernikahanRequest	true	"StatusPernikahan data"
//	@Success	201	{object}	response.MyGoResponse{data=dto.MasterStatusPernikahanResponse}
//	@Router		/master/status_pernikahan				[post]
//
// CreateStatusPernikahan handles POST /api/v1/master/status_pernikahan
func (h *MasterHandler) CreateStatusPernikahan(c *echo.Context) error {
	var req dto.CreateMasterStatusPernikahanRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.CreateStatusPernikahan(&req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
} // ───────────── Create ─────────────────────────────────────────────────────────

// ─────────────── Update ─────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary	Update StatusPernikahan
//	@Description	Update StatusPernikahan by :id
//	@Tags	master/status_pernikahan
//	@Accept	json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		id		path		int			true	"StatusPernikahan ID"
//	@Param	body	body	dto.UpdateMasterStatusPernikahanRequest	true	"Update Request"
//	@Success	200	{object}	response.MyGoResponse{data=dto.MasterStatusPernikahanResponse}
//	@Router		/master/status_pernikahan/{id}	[put]
//
// UpdateStatusPernikahan handles PUT /api/v1/master/status_pernikahan/:id
func (h *MasterHandler) UpdateStatusPernikahan(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateMasterStatusPernikahanRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.UpdateStatusPernikahan(id, &req, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
} // ───────────── Update ─────────────────────────────────────────────────────────

// ─────────────── Delete ─────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary	Delete StatusPernikahan
//	@Description	Delete StatusPernikahan by :id
//	@Tags	master/status_pernikahan
//	@Accept	json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		id	path	int	true	"StatusPernikahan ID"
//	@Success	200	{object}	response.MyGoResponse
//	@Router	/master/status_pernikahan/{id}	[delete]
//
// DeleteStatusPernikahan handles DELETE /api/v1/master/status_pernikahan/:id
func (h *MasterHandler) DeleteStatusPernikahan(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := he.BuildAuthContext(c)
	err = h.service.DeleteStatusPernikahan(id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
} // ───────────── Delete ─────────────────────────────────────────────────────────
