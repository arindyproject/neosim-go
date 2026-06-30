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
//		@Summary		Get list of GolonganDarah
//		@Description	Get paginated list of GolonganDarah
//		@Tags			master/golongan_darah
//		@Accept			json
//		@Produce		json
//		@Security		BearerAuth
//		@Param			kode_kemenkes	query		string	false	"Filter by kode_kemenkes"
//		@Param			name			query		string	false	"Filter by name (partial match)"
//		@Param			page			query		int		false	"Page number"
//		@Param			page_size		query		int		false	"Page size"
//		@Success		200				{object}	response.MyGoResponse{data=[]dto.MasterGolonganDarahResponse}
//		@Router			/master/golongan_darah [get]
//
// ListGolonganDarah handles GET /api/v1/master/golongan_darah
func (h *MasterHandler) ListGolonganDarah(c *echo.Context) error {
	page, pageSize := he.ParsePagination(c, h.cfg)
	filter := dto.FilterMasterGolonganDarahRequest{
		KodeKemenkes: c.QueryParam("kode_kemenkes"),
		Name:         c.QueryParam("name"),
	}

	items, total, err := h.service.ListGolonganDarah(page, pageSize, &filter)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}

	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
} // ───────────── List ───────────────────────────────────────────────────────────

// ─────────────── GetByID ────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary		Get GolonganDarah
//	@Description	Get GolonganDarah detail by :id
//	@Tags			master/golongan_darah
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"GolonganDarah ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.MasterGolonganDarahResponse}
//	@Router			/master/golongan_darah/{id}	[get]
//
// GetByIDGolonganDarah handles GET /api/v1/master/golongan_darah/:id
func (h *MasterHandler) GetByIDGolonganDarah(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	item, err := h.service.GetByIDGolonganDarah(id)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── Create ─────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary	Create GolonganDarah
//	@Description	Create new GolonganDarah
//	@Tags		master/golongan_darah
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		body 		body		dto.CreateMasterGolonganDarahRequest	true	"GolonganDarah Request"
//	@Success	201	{object}	response.MyGoResponse{data=dto.MasterGolonganDarahResponse}
//	@Router		/master/golongan_darah				[post]
//
// CreateGolonganDarah handles POST /api/v1/master/golongan_darah
func (h *MasterHandler) CreateGolonganDarah(c *echo.Context) error {
	var req dto.CreateMasterGolonganDarahRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.CreateGolonganDarah(&req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
} // ───────────── Create ─────────────────────────────────────────────────────────

// ─────────────── Update ─────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary	Update GolonganDarah
//	@Description	Update GolonganDarah by :id
//	@Tags	master/golongan_darah
//	@Accept	json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		id		path		int			true	"GolonganDarah ID"
//	@Param	body	body	dto.UpdateMasterGolonganDarahRequest	true	"Update Request"
//	@Success	200	{object}	response.MyGoResponse{data=dto.MasterGolonganDarahResponse}
//	@Router		/master/golongan_darah/{id}	[put]
//
// UpdateGolonganDarah handles PUT /api/v1/master/golongan_darah/:id
func (h *MasterHandler) UpdateGolonganDarah(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateMasterGolonganDarahRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.UpdateGolonganDarah(id, &req, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
} // ───────────── Update ─────────────────────────────────────────────────────────

// ─────────────── Delete ─────────────────────────────────────────────────────────
// MasterHandler godoc
//
//	@Summary	Delete GolonganDarah
//	@Description	Delete GolonganDarah by :id
//	@Tags	master/golongan_darah
//	@Accept	json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		id	path	int	true	"GolonganDarah ID"
//	@Success	200	{object}	response.MyGoResponse
//	@Router	/master/golongan_darah/{id}	[delete]
//
// DeleteGolonganDarah handles DELETE /api/v1/master/golongan_darah/:id
func (h *MasterHandler) DeleteGolonganDarah(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := he.BuildAuthContext(c)
	err = h.service.DeleteGolonganDarah(id, actor)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
} // ───────────── Delete ─────────────────────────────────────────────────────────
