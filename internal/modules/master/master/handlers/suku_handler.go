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
//		@Summary		Get list of Suku
//		@Description	Get paginated list of Suku
//		@Tags			master/suku
//		@Accept			json
//		@Produce		json
//		@Security		BearerAuth
//		@Param			kode_kemenkes	query		string	false	"Filter by kode_kemenkes"
//		@Param			name			query		string	false	"Filter by name (partial match)"
//		@Param			page			query		int		false	"Page number"
//		@Param			page_size		query		int		false	"Page size"
//		@Success		200				{object}	response.MyGoResponse{data=[]dto.MasterSukuResponse}
//		@Router			/master/suku [get]
//
// ListSuku handles GET /api/v1/master/suku
func (h *MasterHandler) ListSuku(c *echo.Context) error {
	page, pageSize := he.ParsePagination(c)
	filter := dto.FilterMasterSukuRequest{
		KodeKemenkes: c.QueryParam("kode_kemenkes"),
		Name:         c.QueryParam("name"),
	}

	items, total, err := h.service.ListSuku(page, pageSize, &filter)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}

	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
} // ───────────── List ───────────────────────────────────────────────────────────

// ─────────────── GetByID ────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary		Get Suku
//	@Description	Get Suku detail by :id
//	@Tags			master/suku
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Suku ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.MasterSukuResponse}
//	@Router			/master/suku/{id}	[get]
//
// GetByIDSuku handles GET /api/v1/master/suku/:id
func (h *MasterHandler) GetByIDSuku(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	item, err := h.service.GetByIDSuku(id)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── Create ─────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary	Create Suku
//	@Description	Create new Suku
//	@Tags		master/suku
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		body 		body		dto.CreateMasterSukuRequest	true	"Suku Request"
//	@Success	201	{object}	response.MyGoResponse{data=dto.MasterSukuResponse}
//	@Router		/master/suku				[post]
//
// CreateSuku handles POST /api/v1/master/suku
func (h *MasterHandler) CreateSuku(c *echo.Context) error {
	var req dto.CreateMasterSukuRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.CreateSuku(&req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
} // ───────────── Create ─────────────────────────────────────────────────────────

// ─────────────── Update ─────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary	Update Suku
//	@Description	Update Suku by :id
//	@Tags	master/suku
//	@Accept	json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		id		path		int			true	"Suku ID"
//	@Param	body	body	dto.UpdateMasterSukuRequest	true	"Update Request"
//	@Success	200	{object}	response.MyGoResponse{data=dto.MasterSukuResponse}
//	@Router		/master/suku/{id}	[put]
//
// UpdateSuku handles PUT /api/v1/master/suku/:id
func (h *MasterHandler) UpdateSuku(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateMasterSukuRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.UpdateSuku(id, &req, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
} // ───────────── Update ─────────────────────────────────────────────────────────

// ─────────────── Delete ─────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary	Delete Suku
//	@Description	Delete Suku by :id
//	@Tags	master/suku
//	@Accept	json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		id	path	int	true	"Suku ID"
//	@Success	200	{object}	response.MyGoResponse
//	@Router	/master/suku/{id}	[delete]
//
// DeleteSuku handles DELETE /api/v1/master/suku/:id
func (h *MasterHandler) DeleteSuku(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := he.BuildAuthContext(c)
	err = h.service.DeleteSuku(id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
} // ───────────── Delete ─────────────────────────────────────────────────────────
