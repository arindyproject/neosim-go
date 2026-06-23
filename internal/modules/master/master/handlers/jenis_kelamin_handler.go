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
//		@Summary		Get list of JenisKelamin
//		@Description	Get paginated list of JenisKelamin
//		@Tags			master/jenis_kelamin
//		@Accept			json
//		@Produce		json
//		@Security		BearerAuth
//		@Param			kode_kemenkes	query		string	false	"Filter by kode_kemenkes"
//		@Param			name			query		string	false	"Filter by name (partial match)"
//		@Param			page			query		int		false	"Page number"
//		@Param			page_size		query		int		false	"Page size"
//		@Success		200				{object}	response.MyGoResponse{data=[]dto.MasterJenisKelaminResponse}
//		@Router			/master/jenis_kelamin [get]
//
// ListJenisKelamin handles GET /api/v1/master/jenis_kelamin
func (h *MasterHandler) ListJenisKelamin(c *echo.Context) error {
	page, pageSize := he.ParsePagination(c)
	filter := dto.FilterMasterJenisKelaminRequest{
		KodeKemenkes: c.QueryParam("kode_kemenkes"),
		Name:         c.QueryParam("name"),
	}

	items, total, err := h.service.ListJenisKelamin(page, pageSize, &filter)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}

	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
} // ───────────── List ───────────────────────────────────────────────────────────

// ─────────────── GetByID ────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary		Get JenisKelamin
//	@Description	Get JenisKelamin detail by :id
//	@Tags			master/jenis_kelamin
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"JenisKelamin ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.MasterJenisKelaminResponse}
//	@Router			/master/jenis_kelamin/{id}	[get]
//
// GetByIDJenisKelamin handles GET /api/v1/master/jenis_kelamin/:id
func (h *MasterHandler) GetByIDJenisKelamin(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	item, err := h.service.GetByIDJenisKelamin(id)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── Create ─────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary	Create JenisKelamin
//	@Description	Create new JenisKelamin
//	@Tags		master/jenis_kelamin
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		body 		body		dto.CreateMasterJenisKelaminRequest	true	"JenisKelamin Request"
//	@Success	201	{object}	response.MyGoResponse{data=dto.MasterJenisKelaminResponse}
//	@Router		/master/jenis_kelamin				[post]
//
// CreateJenisKelamin handles POST /api/v1/master/jenis_kelamin
func (h *MasterHandler) CreateJenisKelamin(c *echo.Context) error {
	var req dto.CreateMasterJenisKelaminRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.CreateJenisKelamin(&req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
} // ───────────── Create ─────────────────────────────────────────────────────────

// ─────────────── Update ─────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary	Update JenisKelamin
//	@Description	Update JenisKelamin by :id
//	@Tags	master/jenis_kelamin
//	@Accept	json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		id		path		int			true	"JenisKelamin ID"
//	@Param	body	body	dto.UpdateMasterJenisKelaminRequest	true	"Update Request"
//	@Success	200	{object}	response.MyGoResponse{data=dto.MasterJenisKelaminResponse}
//	@Router		/master/jenis_kelamin/{id}	[put]
//
// UpdateJenisKelamin handles PUT /api/v1/master/jenis_kelamin/:id
func (h *MasterHandler) UpdateJenisKelamin(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateMasterJenisKelaminRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.UpdateJenisKelamin(id, &req, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
} // ───────────── Update ─────────────────────────────────────────────────────────

// ─────────────── Delete ─────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary	Delete JenisKelamin
//	@Description	Delete JenisKelamin by :id
//	@Tags	master/jenis_kelamin
//	@Accept	json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		id	path	int	true	"JenisKelamin ID"
//	@Success	200	{object}	response.MyGoResponse
//	@Router	/master/jenis_kelamin/{id}	[delete]
//
// DeleteJenisKelamin handles DELETE /api/v1/master/jenis_kelamin/:id
func (h *MasterHandler) DeleteJenisKelamin(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := he.BuildAuthContext(c)
	err = h.service.DeleteJenisKelamin(id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
} // ───────────── Delete ─────────────────────────────────────────────────────────
