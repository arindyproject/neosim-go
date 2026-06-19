package handlers

import (
	"neosim_go/internal/modules/master/alamat/dto"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"
	"net/http"

	"github.com/labstack/echo/v5"
)

// Provinsi ==========================================================================
// ─────────────── List ─────────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Get list of Provinsi
//	@Description	Get paginated list of Provinsi
//	@Tags			master/alamat/provinsi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			negara_id	query		int		false	"Filter by negara_id"
//	@Param			code		query		string	false	"Filter by code"
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.ProvinsiResponse}
//	@Router			/master/alamat/provinsi [get]
//
// ListProvinsi handles GET /api/v1/master/alamat/provinsi
func (h *MasterAlamatHandler) ListProvinsi(c *echo.Context) error {
	page, pageSize := he.ParsePagination(c)
	negaraID := he.ParseOptionalInt64Query(c, "negara_id")
	filter := dto.FilterProvinsiRequest{
		Code: c.QueryParam("code"),
		Name: c.QueryParam("name"),
	}

	items, total, err := h.service.ListProvinsi(page, pageSize, negaraID, &filter)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}

	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
} // ───────────── List ───────────────────────────────────────────────────────────

// ─────────────── GetByID ────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Get Provinsi
//	@Description	Get Provinsi detail by :id, termasuk statistik turunan
//	@Tags			master/alamat/provinsi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Provinsi ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.ProvinsiDetailResponse}
//	@Router			/master/alamat/provinsi/{id} [get]
//
// GetByIDProvinsi handles GET /api/v1/master/alamat/provinsi/:id
func (h *MasterAlamatHandler) GetByIDProvinsi(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	item, err := h.service.GetByIDProvinsi(id)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── Create ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Create Provinsi
//	@Description	Create New Provinsi
//	@Tags			master/alamat/provinsi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateProvinsiRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.ProvinsiResponse}
//	@Router			/master/alamat/provinsi [post]
//
// CreateProvinsi handles POST /api/v1/master/alamat/provinsi
func (h *MasterAlamatHandler) CreateProvinsi(c *echo.Context) error {
	var req dto.CreateProvinsiRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.CreateProvinsi(&req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
} // ───────────── Create ─────────────────────────────────────────────────────────

// ─────────────── Update ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Update Provinsi
//	@Description	Update Provinsi by :id
//	@Tags			master/alamat/provinsi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int							true	"Provinsi ID"
//	@Param			body	body		dto.UpdateProvinsiRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.ProvinsiResponse}
//	@Router			/master/alamat/provinsi/{id} [put]
//
// UpdateProvinsi handles PUT /api/v1/master/alamat/provinsi/:id
func (h *MasterAlamatHandler) UpdateProvinsi(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateProvinsiRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.UpdateProvinsi(id, &req, actor)
	if err != nil {
		return response.Response(c, he.NotFoundStatus(err, "Provinsi tidak ditemukan"), false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
} // ───────────── Update ─────────────────────────────────────────────────────────

// ─────────────── Delete ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Delete Provinsi
//	@Description	Delete Provinsi by :id
//	@Tags			master/alamat/provinsi
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Provinsi ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/master/alamat/provinsi/{id} [delete]
//
// DeleteProvinsi handles DELETE /api/v1/master/alamat/provinsi/:id
func (h *MasterAlamatHandler) DeleteProvinsi(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := he.BuildAuthContext(c)
	if err := h.service.DeleteProvinsi(id, actor); err != nil {
		return response.Response(c, he.NotFoundStatus(err, "Provinsi tidak ditemukan"), false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
} // ───────────── Delete ─────────────────────────────────────────────────────────
// Provinsi ==========================================================================
