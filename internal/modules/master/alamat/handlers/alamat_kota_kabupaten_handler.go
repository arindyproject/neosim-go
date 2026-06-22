package handlers

import (
	"neosim_go/internal/modules/master/alamat/dto"
	he "neosim_go/internal/shared/httputil"
	"neosim_go/internal/shared/response"
	"neosim_go/internal/shared/validator"
	"net/http"

	"github.com/labstack/echo/v5"
)

// Kota/Kabupaten =====================================================================
// ─────────────── List ─────────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Get list of Kota/Kabupaten
//	@Description	Get paginated list of Kota/Kabupaten
//	@Tags			master/alamat/kota
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			provinsi_id	query		int		false	"Filter by provinsi_id"
//	@Param			code		query		string	false	"Filter by code"
//	@Param			name		query		string	false	"Filter by name (partial match)"
//	@Param			page		query		int		false	"Page number"
//	@Param			page_size	query		int		false	"Page size"
//	@Success		200			{object}	response.MyGoResponse{data=[]dto.KotaKabupatenResponse}
//	@Router			/master/alamat/kota [get]
//
// ListKotaKabupaten handles GET /api/v1/master/alamat/kota
func (h *MasterAlamatHandler) ListKotaKabupaten(c *echo.Context) error {
	page, pageSize := he.ParsePagination(c)
	provinsiID := he.ParseOptionalInt64Query(c, "provinsi_id")
	filter := dto.FilterKotaKabupatenRequest{
		Code: c.QueryParam("code"),
		Name: c.QueryParam("name"),
	}

	items, total, err := h.service.ListKotaKabupaten(page, pageSize, provinsiID, &filter)
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, false, "Gagal mengambil data", nil, nil)
	}

	return response.Paginated(c, http.StatusOK, true, "Berhasil mengambil data", items, total, page, pageSize)
} // ───────────── List ───────────────────────────────────────────────────────────

// ─────────────── GetByID ────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Get Kota/Kabupaten
//	@Description	Get Kota/Kabupaten detail by :id, termasuk statistik turunan
//	@Tags			master/alamat/kota
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Kota/Kabupaten ID"
//	@Success		200	{object}	response.MyGoResponse{data=dto.KotaKabupatenDetailResponse}
//	@Router			/master/alamat/kota/{id} [get]
//
// GetByIDKotaKabupaten handles GET /api/v1/master/alamat/kota/:id
func (h *MasterAlamatHandler) GetByIDKotaKabupaten(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	item, err := h.service.GetByIDKotaKabupaten(id)
	if err != nil {
		return response.Response(c, http.StatusNotFound, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Berhasil mengambil data", item, nil)
} // ───────────── GetByID ────────────────────────────────────────────────────────

// ─────────────── Create ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Create Kota/Kabupaten
//	@Description	Create New Kota/Kabupaten
//	@Tags			master/alamat/kota
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateKotaKabupatenRequest	true	"Create Request"
//	@Success		201		{object}	response.MyGoResponse{data=dto.KotaKabupatenResponse}
//	@Router			/master/alamat/kota [post]
//
// CreateKotaKabupaten handles POST /api/v1/master/alamat/kota
func (h *MasterAlamatHandler) CreateKotaKabupaten(c *echo.Context) error {
	var req dto.CreateKotaKabupatenRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.CreateKotaKabupaten(&req, actor)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusCreated, true, "Data berhasil dibuat", item, nil)
} // ───────────── Create ─────────────────────────────────────────────────────────

// ─────────────── Update ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Update Kota/Kabupaten
//	@Description	Update Kota/Kabupaten by :id
//	@Tags			master/alamat/kota
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int								true	"Kota/Kabupaten ID"
//	@Param			body	body		dto.UpdateKotaKabupatenRequest	true	"Update Request"
//	@Success		200		{object}	response.MyGoResponse{data=dto.KotaKabupatenResponse}
//	@Router			/master/alamat/kota/{id} [put]
//
// UpdateKotaKabupaten handles PUT /api/v1/master/alamat/kota/:id
func (h *MasterAlamatHandler) UpdateKotaKabupaten(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	var req dto.UpdateKotaKabupatenRequest
	if err := c.Bind(&req); err != nil {
		return response.Response(c, http.StatusBadRequest, false, "Request tidak valid", nil, nil)
	}
	if errs := validator.Validate(req); errs != nil {
		return response.Response(c, http.StatusUnprocessableEntity, false, "Validasi gagal", nil, errs)
	}

	actor := he.BuildAuthContext(c)
	item, err := h.service.UpdateKotaKabupaten(id, &req, actor)
	if err != nil {
		return response.Response(c, he.NotFoundStatus(err, "Kota/Kabupaten tidak ditemukan"), false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil diupdate", item, nil)
} // ───────────── Update ─────────────────────────────────────────────────────────

// ─────────────── Delete ─────────────────────────────────────────────────────────
// MasterAlamatHandler godoc
//
//	@Summary		Delete Kota/Kabupaten
//	@Description	Delete Kota/Kabupaten by :id
//	@Tags			master/alamat/kota
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Kota/Kabupaten ID"
//	@Success		200	{object}	response.MyGoResponse{}
//	@Router			/master/alamat/kota/{id} [delete]
//
// DeleteKotaKabupaten handles DELETE /api/v1/master/alamat/kota/:id
func (h *MasterAlamatHandler) DeleteKotaKabupaten(c *echo.Context) error {
	id, err := he.ParseID(c)
	if err != nil {
		return response.Response(c, http.StatusBadRequest, false, "ID tidak valid", nil, nil)
	}

	actor := he.BuildAuthContext(c)
	if err := h.service.DeleteKotaKabupaten(id, actor); err != nil {
		return response.Response(c, he.NotFoundStatus(err, "Kota/Kabupaten tidak ditemukan"), false, err.Error(), nil, nil)
	}

	return response.Response(c, http.StatusOK, true, "Data berhasil dihapus", nil, nil)
} // ───────────── Delete ─────────────────────────────────────────────────────────
// Kota/Kabupaten =====================================================================
