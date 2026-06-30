package handlers

import (
	"neosim_go/internal/modules/master/alamat/dto"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"
	"net/http"

	"github.com/labstack/echo/v5"
)

// Kelurahan/Desa ========================================================================
// ─────────────── List ─────────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Get list of Kelurahan/Desa
//	@Description	Get paginated list of Kelurahan/Desa
//	@Tags			master/alamat/desa
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			kecamatan_id	query		int		false	"Filter by kecamatan_id"
//	@Param			code			query		string	false	"Filter by code"
//	@Param			name			query		string	false	"Filter by name (partial match)"
//	@Param			postal_code		query		string	false	"Filter by postal code"
//	@Param			page			query		int		false	"Page number"
//	@Param			page_size		query		int		false	"Page size"
//	@Success		200				{object}	response.MyGoResponse{data=[]dto.KelurahanDesaResponse}
//	@Router			/master/alamat/desa [get]
//
// ListKelurahanDesa handles GET /api/v1/master/alamat/desa
func (h *MasterAlamatHandler) ListKelurahanDesa(c *echo.Context) error {
	page, pageSize := he.ParsePagination(c, h.cfg)
	kecamatanID := he.ParseOptionalInt64Query(c, "kecamatan_id")
	filter := dto.FilterKelurahanDesaRequest{
		Code:       c.QueryParam("code"),
		Name:       c.QueryParam("name"),
		PostalCode: c.QueryParam("postal_code"),
	}

	items, total, err := h.service.ListKelurahanDesa(page, pageSize, kecamatanID, &filter)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}

	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
} // ───────────── List ───────────────────────────────────────────────────────────

// ─────────────── GetByID ────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Get Kelurahan/Desa
//	@Description	Get Kelurahan/Desa detail by :id, termasuk jalur hierarki lengkap
//	@Tags			master/alamat/desa
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Kelurahan/Desa ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.KelurahanDesaDetailResponse}
//	@Router			/master/alamat/desa/{id} [get]
//
// GetByIDKelurahanDesa handles GET /api/v1/master/alamat/desa/:id
func (h *MasterAlamatHandler) GetByIDKelurahanDesa(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	item, err := h.service.GetByIDKelurahanDesa(id)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── Create ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Create Kelurahan/Desa
//	@Description	Create New Kelurahan/Desa
//	@Tags			master/alamat/desa
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateKelurahanDesaRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.KelurahanDesaResponse}
//	@Router			/master/alamat/desa [post]
//
// CreateKelurahanDesa handles POST /api/v1/master/alamat/desa
func (h *MasterAlamatHandler) CreateKelurahanDesa(c *echo.Context) error {
	var req dto.CreateKelurahanDesaRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.CreateKelurahanDesa(&req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
} // ───────────── Create ─────────────────────────────────────────────────────────

// ─────────────── Update ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Update Kelurahan/Desa
//	@Description	Update Kelurahan/Desa by :id
//	@Tags			master/alamat/desa
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int								true	"Kelurahan/Desa ID"
//	@Param			body	body		dto.UpdateKelurahanDesaRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.KelurahanDesaResponse}
//	@Router			/master/alamat/desa/{id} [put]
//
// UpdateKelurahanDesa handles PUT /api/v1/master/alamat/desa/:id
func (h *MasterAlamatHandler) UpdateKelurahanDesa(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateKelurahanDesaRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.UpdateKelurahanDesa(id, &req, actor)
	if err != nil {
		return response.Response(c, he.NotFoundStatus(err, "Kelurahan/Desa tidak ditemukan"), false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
} // ───────────── Update ─────────────────────────────────────────────────────────

// ─────────────── Delete ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Delete Kelurahan/Desa
//	@Description	Delete Kelurahan/Desa by :id
//	@Tags			master/alamat/desa
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Kelurahan/Desa ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/master/alamat/desa/{id} [delete]
//
// DeleteKelurahanDesa handles DELETE /api/v1/master/alamat/desa/:id
func (h *MasterAlamatHandler) DeleteKelurahanDesa(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := he.BuildAuthContext(c)
	if err := h.service.DeleteKelurahanDesa(id, actor); err != nil {
		return response.Response(c, he.NotFoundStatus(err, "Kelurahan/Desa tidak ditemukan"), false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
} // ───────────── Delete ─────────────────────────────────────────────────────────
// Kelurahan/Desa ========================================================================
