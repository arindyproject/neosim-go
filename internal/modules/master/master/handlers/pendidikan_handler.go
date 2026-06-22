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
//		@Summary		Get list of Pendidikan
//		@Description	Get paginated list of Pendidikan
//		@Tags			master/pendidikan
//		@Accept			json
//		@Produce		json
//		@Security		BearerAuth
//		@Param			kode_kemenkes	query		string	false	"Filter by kode_kemenkes"
//		@Param			name			query		string	false	"Filter by name (partial match)"
//		@Param			page			query		int		false	"Page number"
//		@Param			page_size		query		int		false	"Page size"
//		@Success		200				{object}	response.MyGoResponse{data=[]dto.MasterPendidikanResponse}
//		@Router			/master/pendidikan [get]
//
// ListPendidikan handles GET /api/v1/master/pendidikan
func (h *MasterHandler) ListPendidikan(c *echo.Context) error {
	page, pageSize := he.ParsePagination(c)
	filter := dto.FilterMasterPendidikanRequest{
		KodeKemenkes: c.QueryParam("kode_kemenkes"),
		Name:         c.QueryParam("name"),
	}

	items, total, err := h.service.ListPendidikan(page, pageSize, &filter)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}

	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
} // ───────────── List ───────────────────────────────────────────────────────────

// ─────────────── GetByID ────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary		Get Pendidikan
//	@Description	Get Pendidikan detail by :id
//	@Tags			master/pendidikan
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Pendidikan ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.MasterPendidikanResponse}
//	@Router			/master/pendidikan/{id}	[get]
//
// GetByIDPendidikan handles GET /api/v1/master/pendidikan/:id
func (h *MasterHandler) GetByIDPendidikan(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	item, err := h.service.GetByIDPendidikan(id)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── Create ─────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary	Create Pendidikan
//	@Description	Create new Pendidikan
//	@Tags		master/pendidikan
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		body	body	dto.CreateMasterPendidikanRequest	true	"Pendidikan data"
//	@Success	201	{object}	response.MyGoResponse{data=dto.MasterPendidikanResponse}
//	@Router		/master/pendidikan				[post]
//
// CreatePendidikan handles POST /api/v1/master/pendidikan
func (h *MasterHandler) CreatePendidikan(c *echo.Context) error {
	var req dto.CreateMasterPendidikanRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.CreatePendidikan(&req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
} // ───────────── Create ─────────────────────────────────────────────────────────

// ─────────────── Update ─────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary	Update Pendidikan
//	@Description	Update Pendidikan by :id
//	@Tags	master/pendidikan
//	@Accept	json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		id		path		int			true	"Pendidikan ID"
//	@Param	body	body	dto.UpdateMasterPendidikanRequest	true	"Update Request"
//	@Success	200	{object}	response.MyGoResponse{data=dto.MasterPendidikanResponse}
//	@Router		/master/pendidikan/{id}	[put]
//
// UpdatePendidikan handles PUT /api/v1/master/pendidikan/:id
func (h *MasterHandler) UpdatePendidikan(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateMasterPendidikanRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.UpdatePendidikan(id, &req, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
} // ───────────── Update ─────────────────────────────────────────────────────────

// ─────────────── Delete ─────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary	Delete Pendidikan
//	@Description	Delete Pendidikan by :id
//	@Tags	master/pendidikan
//	@Accept	json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		id	path	int	true	"Pendidikan ID"
//	@Success	200	{object}	response.MyGoResponse
//	@Router	/master/pendidikan/{id}	[delete]
//
// DeletePendidikan handles DELETE /api/v1/master/pendidikan/:id
func (h *MasterHandler) DeletePendidikan(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := he.BuildAuthContext(c)
	err = h.service.DeletePendidikan(id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
} // ───────────── Delete ─────────────────────────────────────────────────────────
